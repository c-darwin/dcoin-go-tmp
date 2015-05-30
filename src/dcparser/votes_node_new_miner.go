package dcparser

import (
	"fmt"
	"utils"
//	"encoding/json"
	//"regexp"
	//"math"
	"strings"
	"os"
)

// голосования нодов, которые должны сохранить фото у себя.
// если смог загрузить фото к себе и хэш сошелся - 1, если нет - 0
// эту транзакцию генерит нод со своим ключом

func (p *Parser) VotesNodeNewMinerInit() (error) {
	fields := []string {"vote_id", "result", "sign"}
	TxMap := make(map[string][]byte)
	TxMap, err := p.GetTxMap(fields);
	p.TxMap = TxMap;
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}

func (p *Parser) VotesNodeNewMinerFront() (error) {
	err := p.generalCheck()
	if err != nil {
		return p.ErrInfo(err)
	}
	// является ли данный юзер майнером
	err = p.checkMiner(p.TxUserID)
	if err != nil {
		return err
	}

	if !utils.CheckInputData(p.TxMap["result"], "vote") {
		return utils.ErrInfoFmt("incorrect vote")
	}
	// получим public_key
	p.nodePublicKey, err = p.GetNodePublicKey(p.TxUserID)
	if len(p.nodePublicKey)==0 {
		return utils.ErrInfoFmt("incorrect user_id len(nodePublicKey) = 0")
	}

	if !utils.CheckInputData(p.TxMap["vote_id"], "bigint") {
		return utils.ErrInfoFmt("incorrect bigint")
	}

	forSign := fmt.Sprintf("%s,%s,%s,%s,%s", p.TxMap["type"], p.TxMap["time"], p.TxMap["user_id"], p.TxMap["vote_id"], p.TxMap["result"])
	fmt.Println("forSign", forSign)
	CheckSignResult, err := utils.CheckSign([][]byte{p.nodePublicKey}, forSign, p.TxMap["sign"], true);
	if err != nil {
		return p.ErrInfo(err)
	}
	if !CheckSignResult {
		return utils.ErrInfoFmt("incorrect sign")
	}

	// проверим, верно ли указан ID и не закончилось ли голосование
	id, err := p.Single("SELECT id FROM votes_miners WHERE id = ? AND type = 'node_voting' AND votes_end = 0", p.TxMap["vote_id"]).Int64()
	if err != nil {
		return p.ErrInfo(err)
	}
	if id == 0 {
		return p.ErrInfo(fmt.Errorf("voting is over"))
	}

	// проверим, не повторное ли это голосование данного юзера
	num, err := p.Single("SELECT count(user_id) FROM log_votes WHERE user_id = ? AND voting_id = ? AND type = 'votes_miners'", p.TxMap["user_id"], p.TxMap["vote_id"]).Int64()
	if err != nil {
		return p.ErrInfo(err)
	}
	if num > 0 {
		return utils.ErrInfoFmt("double voting")
	}

	// нод не должен голосовать более X раз за сутки, чтобы не было доса
	err = p.limitRequest(p.Variables.Int64["node_voting"], "votes_nodes", p.Variables.Int64["node_voting_period"])
	if err != nil {
		return p.ErrInfo(err)
	}

	return nil
}


func (p *Parser) VotesNodeNewMiner() (error) {
	var votes [2]int
	votesData, err := p.OneRow("SELECT user_id, votes_start_time, votes_0, votes_1 FROM votes_miners WHERE id = ?", p.TxMap["vote_id"]).Int()
	if err != nil {
		return p.ErrInfo(err)
	}
	minersData, err := p.OneRow("SELECT photo_block_id, photo_max_miner_id, miners_keepers, log_id FROM miners_data WHERE user_id = ?", votesData["user_id"]).String()
	// $votes_data['user_id'] - это юзер, за которого голосуют
	if err != nil {
		return p.ErrInfo(err)
	}

	votes[0] = votesData["votes_0"]
	votes[1] = votesData["votes_1"]
	// прибавим голос
	votes[utils.BytesToInt(p.TxMap["result"])]++

	// обновляем голоса. При откате просто вычитаем
	err = p.ExecSql("UPDATE votes_miners SET votes_"+string(p.TxMap["result"])+" = ? WHERE id = ?", votes[utils.BytesToInt(p.TxMap["result"])], p.TxMap["vote_id"] )
	if err != nil {
		return p.ErrInfo(err)
	}

	// логируем, чтобы юзер {$this->tx_data['user_id']} не смог повторно проголосовать
	err = p.ExecSql("INSERT INTO log_votes (user_id, voting_id, type) VALUES (?, ?, 'votes_miners')", p.TxMap["user_id"], p.TxMap["vote_id"])
	if err != nil {
		return p.ErrInfo(err)
	}

	// ID майнеров, у которых сохраняются фотки
	minersIds := getMinersKeepers( minersData["photo_block_id"], minersData["photo_max_miner_id"], minersData["miners_keepers"], true)

	// данные для проверки окончания голосования

	minerData := new(MinerData)
	minerData.myMinersIds, err = p.getMyMinersIds()
	if err != nil {
		return p.ErrInfo(err)
	}
	minerData.adminUiserId, err = p.GetAdminUserId()
	if err != nil {
		return p.ErrInfo(err)
	}
	minerData.minersIds = minersIds
	minerData.votes0 = votes[0]
	minerData.votes1 = votes[1]
	minerData.minMinersKeepers = p.Variables.Int64["min_miners_keepers"]
	if p.minersCheckVotes1(minerData) || p.minersCheckMyMinerIdAndVotes0(minerData) {
		// отмечаем, что голосование нодов закончено
		err = p.ExecSql("UPDATE votes_miners SET votes_end = 1, end_block_id = ? WHERE id = ?", p.BlockData.BlockId, p.TxMap["vote_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
		// отметим del_block_id всем, кто голосовал за данного юзера,
		// чтобы через N блоков по крону удалить бесполезные записи
		err = p.ExecSql("UPDATE log_votes SET del_block_id = ? WHERE voting_id = ? AND type = 'votes_miners'", p.BlockData.BlockId, p.TxMap["vote_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
	}

	// если набрано >=X голосов "за", то пишем в БД, что юзер готов к проверке людьми
	// либо если набранное кол-во голосов= кол-ву майнеров (актуально в самом начале запуска проекта)
	if p.minersCheckVotes1(minerData) {
		err = p.ExecSql("INSERT INTO votes_miners ( user_id, type, votes_start_time ) VALUES ( ?, 'user_voting', ? )", votesData["user_id"], p.BlockData.Time)
		if err != nil {
			return p.ErrInfo(err)
		}

		// и отмечаем лицо как готовое участвовать в поиске дублей
		err = p.ExecSql("UPDATE faces SET status = 'used' WHERE user_id = ?", votesData["user_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
	} else if (p.minersCheckMyMinerIdAndVotes0(minerData)) {
		// если набрано >5 голосов "против" и мы среди тех X майнеров, которые копировали фото к себе
		// либо если набранное кол-во голосов = кол-ву майнеров (актуально в самом начале запуска проекта)
		facePath := fmt.Sprintf("public/face_%v.jpg", votesData["user_id"])
		profilePath := fmt.Sprintf("public/profile_%v.jpg", votesData["user_id"])

		faceRandName := ""
		profileRandName := ""
		// возможно фото к нам не было скопировано, т.к. хост был недоступен.
		if _, err := os.Stat(profilePath); os.IsNotExist(err) {
			faceRandName = ""
			profileRandName = ""
		} else if _, err := os.Stat(facePath); os.IsNotExist(err) {
			faceRandName = ""
			profileRandName = ""
		} else {
			faceRandName = utils.RandSeq(30)
			profileRandName = utils.RandSeq(30)

			// перемещаем фото в корзину, откуда по крону будем удалять данные
			err = utils.CopyFileContents(facePath, faceRandName)
			if err != nil {
				return p.ErrInfo(err)
			}
			err =os.Remove(facePath)
			if err != nil {
				return p.ErrInfo(err)
			}
			err = utils.CopyFileContents(profilePath, profileRandName)
			if err != nil {
				return p.ErrInfo(err)
			}
			err =os.Remove(profilePath)
			if err != nil {
				return p.ErrInfo(err)
			}

			// если в корзине что-то есть, то логируем
			// отсутствие файлов также логируем, т.к. больше негде, а при откате эти данные очень важны.
			logData, err := p.OneRow("SELECT * FROM recycle_bin WHERE user_id = ?", votesData["user_id"]).String()
			if err != nil {
				return p.ErrInfo(err)
			}
			if len(logData) > 0 {
				logId, err := p.ExecSqlGetLastInsertId("INSERT INTO log_recycle_bin ( user_id, profile_file_name, face_file_name, block_id, prev_log_id ) VALUES ( ?, ?, ?, ?, ? )", logData["user_id"], logData["profile_file_name"], logData["face_file_name"], p.BlockData.BlockId, logData["log_id"])
				if err != nil {
					return p.ErrInfo(err)
				}
				err = p.ExecSql("UPDATE recycle_bin SET profile_file_name = ?, face_file_name = ?, log_id = ? WHERE user_id = ?", profileRandName, faceRandName, logId, logData["user_id"])
				if err != nil {
					return p.ErrInfo(err)
				}
			} else {
				err = p.ExecSql("INSERT INTO recycle_bin ( user_id, profile_file_name, face_file_name ) VALUES ( ?, ?, ? )", votesData["user_id"], profileRandName, faceRandName)
				if err != nil {
					return p.ErrInfo(err)
				}
			}
		}
	}

	return nil
}


type MinerData struct {
	adminUiserId int64
	myMinersIds map[int]int
	minersIds map[int]int
	votes0 int
	votes1 int
	minMinersKeepers int64
}

func (p *Parser) VotesNodeNewMinerRollback() (error) {

	votesData, err := p.OneRow("SELECT user_id, votes_start_time, votes_0, votes_1 FROM votes_miners WHERE id = ?", p.TxMap["vote_id"]).Int()
	if err != nil {
		return p.ErrInfo(err)
	}
	minersData, err := p.OneRow("SELECT photo_block_id, photo_max_miner_id, miners_keepers, log_id FROM miners_data WHERE user_id = ?", votesData["user_id"]).String()
	if err != nil {
		return p.ErrInfo(err)
	}

	minerData := new(MinerData)
	// запомним голоса- пригодится чуть ниже в minersCheckVotes1
	minerData.votes0 = votesData["votes_0"]
	minerData.votes1 = votesData["votes_1"]

	var votes [2]int
	votes[0] = votesData["votes_0"]
	votes[1] = votesData["votes_1"]
	// вычтем голос
	votes[utils.BytesToInt(p.TxMap["result"])]--

	// обновляем голоса
	err = p.ExecSql("UPDATE votes_miners SET votes_"+string(p.TxMap["result"])+" = ? WHERE id = ?", votes[utils.BytesToInt(p.TxMap["result"])], p.TxMap["vote_id"])
	if err != nil {
		return p.ErrInfo(err)
	}

	// удаляем нашу запись из log_votes
	err = p.ExecSql("DELETE FROM log_votes WHERE user_id = ? AND voting_id = ? AND type = 'votes_miners'", p.TxMap["user_id"], p.TxMap["vote_id"])
	if err != nil {
		return p.ErrInfo(err)
	}
	minersIds := getMinersKeepers( minersData["photo_block_id"], minersData["photo_max_miner_id"], minersData["miners_keepers"], true)
	minerData.myMinersIds, err = p.getMyMinersIds()
	if err != nil {
		return p.ErrInfo(err)
	}
	minerData.minersIds = minersIds
	minerData.minMinersKeepers = p.Variables.Int64["min_miners_keepers"]

	if p.minersCheckVotes1(minerData) || p.minersCheckMyMinerIdAndVotes0(minerData) {

		// отменяем отметку о том, что голосование нодов закончено
		err = p.ExecSql("UPDATE votes_miners SET votes_end = 0, end_block_id = 0 WHERE id = ?", p.TxMap["vote_id"])
		if err != nil {
			return p.ErrInfo(err)
		}

		// всем, кому ставили del_block_id, его убираем, т.е. отменяем будущее удаление по крону
		err = p.ExecSql("UPDATE log_votes SET del_block_id = 0 WHERE voting_id = ? AND type = 'votes_miners' AND del_block_id = ? ", p.TxMap["vote_id"], p.BlockData.BlockId)
		if err != nil {
			return p.ErrInfo(err)
		}
	}

	// если набрано >=5 голосов, то отменяем  в БД, что юзер готов к проверке людьми
	if p.minersCheckVotes1(minerData) {
		// отменяем созданное юзерское голосование
		err = p.ExecSql("DELETE FROM votes_miners WHERE user_id = ? AND votes_start_time = ? AND type = 'user_voting'", votesData["user_id"], p.BlockData.Time)
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.rollbackAI("votes_miners", 1)
		if err != nil {
			return p.ErrInfo(err)
		}

		// и отмечаем лицо как неучаствующее в поиске клонов
		err = p.ExecSql("UPDATE faces SET status = 'pending' WHERE user_id = ?", votesData["user_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
	} else if (p.minersCheckMyMinerIdAndVotes0(minerData)) {
		// если фото плохое и мы среди тех 10 майнеров, которые копировали (или нет) фото к себе,
		// а затем переместили фото в корзину

		// получаем rand_name из логов
		data, err := p.OneRow("SELECT profile_file_name, face_file_name FROM recycle_bin WHERE user_id = ?", votesData["user_id"]).String()
		if err != nil {
			return p.ErrInfo(err)
		}

		// перемещаем фото из корзины, если есть, что перемещать
		if len(data["profile_file_name"])>0 && len(data["face_file_name"])>0 {
			utils.CopyFileContents("recycle_bin/"+data["face_file_name"], "public/face_"+utils.IntToStr(votesData["user_id"])+".jpg")
			utils.CopyFileContents("recycle_bin/"+data["profile_file_name"], "public/profile_"+utils.IntToStr(votesData["user_id"])+".jpg")
		}
		p.generalRollback("recycle_bin", votesData["user_id"], "", false)
	}


	return nil
}

func (p *Parser) VotesNodeNewMinerRollbackFront() error {
	return p.limitRequestsRollback("votes_nodes")
}

func Int64SliceToStr(Int []int64) []string {
	var result []string
	for _, v := range Int {
		result = append(result, utils.Int64ToStr(v))
	}
	return result
}

func IntSliceToStr(Int []int) []string {
	var result []string
	for _, v := range Int {
		result = append(result, utils.IntToStr(v))
	}
	return result
}

func (p *Parser)  getMyMinersIds() (map[int]int, error) {
	myMinersIds := make(map[int]int)
	var err error
	collective, err := p.GetCommunityUsers()
	if err != nil {
		return myMinersIds, p.ErrInfo(err)
	}
	if len(collective) > 0 {
		myMinersIds, err = p.GetList("SELECT miner_id FROM miners_data WHERE user_id IN ("+strings.Join(Int64SliceToStr(collective), ",")+") AND miner_id > 0").MapInt()
		if err != nil {
			return myMinersIds, p.ErrInfo(err)
		}
	} else {
		minerId, err := p.Single("SELECT miner_id FROM my_table").Int()
		if err != nil {
			return myMinersIds, p.ErrInfo(err)
		}
		myMinersIds[0] = minerId
	}
	return myMinersIds, nil
}
