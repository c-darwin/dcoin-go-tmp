package dcparser

import (
	"fmt"
	"utils"
	"encoding/json"
	//"regexp"
	//"math"
	"strings"
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
		return utils.ErrInfo(err)
	}
	return nil
}

func (p *Parser) VotesNodeNewMinerFront() (error) {
	err := p.generalCheck()
	if err != nil {
		return utils.ErrInfo(err)
	}
	// является ли данный юзер майнером
	if !p.checkMiner(p.TxMap["user_id"]) {
		return utils.ErrInfoFmt("incorrect user_id")
	}

	if !utils.CheckInputData(p.TxMap["result"], "vote") {
		return utils.ErrInfoFmt("incorrect vote")
	}
	// получим public_key
	p.nodePublicKey, err := p.GetNodePublicKey(p.TxUserID)
	if len(p.nodePublicKey)==0 {
		return utils.ErrInfoFmt("incorrect user_id len(nodePublicKey) = 0")
	}

	if !utils.CheckInputData(p.TxMap["vote_id"], "bigint") {
		return utils.ErrInfoFmt("incorrect bigint")
	}

	forSign := fmt.Sprintf("%s,%s,%s,%s,%s", p.TxMap["type"], p.TxMap["time"], p.TxMap["user_id"], p.TxMap["vote_id"], p.TxMap["result"])
	fmt.Println("forSign", forSign)
	CheckSignResult, err := utils.CheckSign([][]byte{p.nodePublicKey}, forSign, p.TxMap["sign"], false);
	if err != nil {
		return utils.ErrInfo(err)
	}
	if !CheckSignResult {
		return utils.ErrInfoFmt("incorrect sign")
	}

	// проверим, верно ли указан ID и не закончилось ли голосование
	id, err := p.Single("SELECT id FROM votes_miners WHERE id = ? AND type = 'node_voting' AND votes_end = 0", p.TxMap["vote_id"]).Int64()
	if err != nil {
		return utils.ErrInfo(err)
	}
	if id == 0 {
		return utils.ErrInfoFmt("voting is over")
	}

	// проверим, не повторное ли это голосование данного юзера
	num, err := p.Single("SELECT count(user_id) FROM log_votes WHERE user_id = ? AND voting_id = ? AND type = 'votes_miners'", p.TxMap["user_id"], p.TxMap["vote_id"]).Int64()
	if err != nil {
		return utils.ErrInfo(err)
	}
	if num > 0 {
		return utils.ErrInfoFmt("double voting")
	}

	// нод не должен голосовать более X раз за сутки, чтобы не было доса
	err = p.limitRequest(p.Variables.Int64["node_voting"], "votes_nodes", p.Variables.Int64["node_voting_period"])
	if err != nil {
		return utils.ErrInfo(err)
	}

	return nil
}


func (p *Parser) VotesNodeNewMiner() (error) {
	var votes [2]int
	votesData, err := p.OneRow("SELECT user_id, votes_start_time, votes_0, votes_1 FROM votes_miners WHERE id = ?", p.TxMap["vote_id"])
	if err != nil {
		return utils.ErrInfo(err)
	}
	minersData, err := p.OneRow("SELECT photo_block_id, photo_max_miner_id, miners_keepers, log_id FROM miners_data WHERE user_id = ?", votesData["user_id"])
	// $votes_data['user_id'] - это юзер, за которого голосуют
	if err != nil {
		return utils.ErrInfo(err)
	}

	votes[0] = utils.StrToInt(votesData["votes_0"])
	votes[1] = utils.StrToInt(votesData["votes_1"])
	// прибавим голос
	votes[utils.StrToInt(p.TxMap["result"])]++

	// обновляем голоса. При откате просто вычитаем
	err = p.ExecSql("UPDATE votes_miners SET votes_"+string(p.TxMap["result"])+"")
	if err != nil {
		return utils.ErrInfo(err)
	}

	// логируем, чтобы юзер {$this->tx_data['user_id']} не смог повторно проголосовать
	err = p.ExecSql("INSERT INTO log_votes (user_id, voting_id, type) VALUES (?, ?, 'votes_miners')", p.TxMap["user_id"], p.TxMap["vote_id"])
	if err != nil {
		return utils.ErrInfo(err)
	}

	// ID майнеров, у которых сохраняются фотки
	minersIds := getMinersKeepers( minersData["photo_block_id"], minersData["photo_max_miner_id"], minersData["miners_keepers"], true)

	// данные для проверки окончания голосования

	minerData := new(MinerData)
	minerData.myMinersIds, err = p.getMyMinersIds()
	if err != nil {
		return utils.ErrInfo(err)
	}
	minerData.minersIds = minersIds
	minerData.votes0 = votesData["votes_0"]
	minerData.votes1 = votesData["votes_1"]
	minerData.minMinersKeepers = p.Variables.Int64["min_miners_keepers"]
	if p.minersCheckVotes1(minerData) || p.minersCheckMyMinerIdAndVotes0(minerData) {
		// отмечаем, что голосование нодов закончено
		p.ExecSql("UPDATE votes_miners SET votes_end = 1, end_block_id = ")
	}
	return nil
}


type MinerData struct {
	myMinersIds map[int]int
	minersIds map[int]int
	votes0 int
	votes1 int
	minMinersKeepers int64
}

func (p *Parser) VotesNodeNewMinerRollback() (error) {

	return nil
}

func (p *Parser) VotesNodeNewMinerRollbackFront() error {
	return p.limitRequestsRollback("new_miner")
}

func (p *Parser)  getMyMinersIds() (map[int]int, error) {
	var myMinersIds make(map[int]int)
	var err error
	collective, err := p.GetCommunityUsers()
	if err != nil {
		return myMinersIds, utils.ErrInfo(err)
	}
	if len(collective) > 0 {
		myMinersIds, err = p.GetList("SELECT miner_id FROM miners_data WHERE user_id IN "+strings.Join(collective, ",")+" AND miner_id > 0").MapInt()
		if err != nil {
			return myMinersIds, utils.ErrInfo(err)
		}
	} else {
		minerId, err := p.Single("SELECT miner_id FROM my_table").String()
		if err != nil {
			return myMinersIds, utils.ErrInfo(err)
		}
		myMinersIds = append(myMinersIds, minerId)
	}
	return myMinersIds, nil
}
