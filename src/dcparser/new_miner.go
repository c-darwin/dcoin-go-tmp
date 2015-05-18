package dcparser

import (
	"fmt"
	"utils"
	"encoding/json"
	"regexp"
)


type exampleSpots struct {
	Face []interface {} `json:"face"`
	Profile []interface {} `json:"profile"`
}
func (p *Parser) NewMinerInit() (error) {
	fields := []string {"race", "country", "latitude", "longitude", "host", "face_coords", "profile_coords", "face_hash", "profile_hash", "video_type", "video_url_id", "node_public_key", "sign"}
	TxMap := make(map[string][]byte)
	TxMap, err := p.GetTxMap(fields);
	p.TxMap = TxMap;
	fmt.Println("TxMap", p.TxMap)
	if err != nil {
		return utils.ErrInfo(err)
	}
	TxMap["node_public_key"] = utils.BinToHex(TxMap["node_public_key"]);
	return nil
}

func (p *Parser) NewMinerFront() (error) {
	err := p.generalCheck()
	if err != nil {
		return utils.ErrInfo(err)
	}
	// получим кол-во точек для face и profile
	exampleSpots_, err := p.DCDB.Single("SELECT example_spots FROM spots_compatibility")
	if err != nil {
		return utils.ErrInfo(err)
	}

	exampleSpots := new(exampleSpots)
	err = json.Unmarshal([]byte(exampleSpots_), &exampleSpots)
	if err != nil {
		return utils.ErrInfo(err)
	}

	if !utils.CheckInputData(p.TxMap["race"], "race") {
		return utils.ErrInfoFmt("race")
	}
	if !utils.CheckInputData(p.TxMap["country"], "country") {
		return utils.ErrInfoFmt("country")
	}
	if !utils.CheckInputData(p.TxMap["latitude"], "coordinate") {
		return utils.ErrInfoFmt("latitude")
	}
	if !utils.CheckInputData(p.TxMap["longitude"], "coordinate") {
		return utils.ErrInfoFmt("longitude")
	}
	if !utils.CheckInputData(p.TxMap["host"], "host") {
		return utils.ErrInfoFmt("host")
	}
	if !utils.CheckInputData_(p.TxMap["face_coords"], "coords", utils.IntToStr(len(exampleSpots.Face))) {
		return utils.ErrInfoFmt("face_coords")
	}
	if !utils.CheckInputData_(p.TxMap["profile_coords"], "coords", utils.IntToStr(len(exampleSpots.Profile))) {
		return utils.ErrInfoFmt("profile_coords")
	}
	if !utils.CheckInputData(p.TxMap["face_hash"], "photo_hash") {
		return utils.ErrInfoFmt("face_hash")
	}
	if !utils.CheckInputData(p.TxMap["profile_hash"], "photo_hash") {
		return utils.ErrInfoFmt("profile_hash")
	}
	if !utils.CheckInputData(p.TxMap["video_type"], "video_type") {
		return utils.ErrInfoFmt("video_type")
	}
	if !utils.CheckInputData(p.TxMap["video_url_id"], "video_url_id") {
		return utils.ErrInfoFmt("video_url_id")
	}
	if !utils.CheckInputData(p.TxMap["node_public_key"], "public_key") {
		return utils.ErrInfoFmt("node_public_key")
	}
	forSign := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s,%s", p.TxMap["type"], p.TxMap["time"], p.TxMap["user_id"], p.TxMap["race"], p.TxMap["country"], p.TxMap["latitude"], p.TxMap["longitude"], p.TxMap["host"], p.TxMap["face_hash"], p.TxMap["profile_hash"], p.TxMap["face_coords"], p.TxMap["profile_coords"], p.TxMap["video_type"], p.TxMap["video_url_id"], p.TxMap["node_public_key"])
	fmt.Println("forSign", forSign)
	CheckSignResult, err := utils.CheckSign(p.PublicKeys, forSign, p.TxMap["sign"], false);
	if err != nil {
		return utils.ErrInfo(err)
	}
	if !CheckSignResult {
		return utils.ErrInfoFmt("incorrect sign")
	}

	// проверим, не кончились ли попытки стать майнером у данного юзера
	count, err := p.countMinerAttempt(string(p.TxMap["user_id"]), "user_voting")
	if count >= utils.BytesToInt64(p.Variables["miner_votes_attempt"]) {
		return utils.ErrInfoFmt("miner_votes_attempt")
	}
	if err != nil {
		return utils.ErrInfo(err)
	}
	//  на всякий случай не даем начать нодовское, если идет юзерское голосование
	userVoting, err := p.DCDB.Single("SELECT id FROM votes_miners WHERE user_id = ? AND type = 'user_voting' AND votes_end = 0", p.TxMap["user_id"])
	if err != nil {
		return utils.ErrInfo(err)
	}
	if len(userVoting) > 0 {
		return utils.ErrInfoFmt("existing $user_voting")
	}

	// проверим, не является ли юзер майнером и  не разжалованный ли это бывший майнер
	minerStatus, err := p.DCDB.Single("SELECT status FROM miners_data WHERE user_id = ? AND status IN ('miner','passive_miner','suspended_miner')", p.TxMap["user_id"])
	if err != nil {
		return utils.ErrInfo(err)
	}
	if len(minerStatus) > 0 {
		return utils.ErrInfoFmt("incorrect miner status")
	}

	// разрешен 1 запрос за сутки
	err = p.limitRequest(p.Variables["limit_new_miner"], "new_miner", p.Variables["limit_new_miner_period"])
	if err != nil {
		return utils.ErrInfo(err)
	}

	return nil
}


func (p *Parser) NewMiner() (error) {
	// пишем в БД нового юзера
	NewMinerId, err := p.DCDB.ExecSqlGetLastInsertId("INSERT INTO users (public_key_0, referral) VALUES ([hex], ?)", p.TxMap["public_key_hex"], p.TxMap["user_id"])
	if err != nil {
		return utils.ErrInfo(err)
	}

	// если работаем в режиме пула, то ищем тех, у кого еще нет user_id
	community, err := p.DCDB.GetCommunityUsers()
	if err != nil {
		return utils.ErrInfo(err)
	}
	if len(community) > 0 {
		for _, userId := range community {
			myPrefix := utils.Int64ToStr(userId)+"_"
			myUserId, err := p.DCDB.Single("SELECT user_id FROM "+myPrefix+"my_table")
			if err != nil {
				return utils.ErrInfo(err)
			}
			if myUserId == "" {
				// проверим, не наш ли это public_key, чтобы записать полученный user_id в my_table
				myPublicKey, err := p.DCDB.Single("SELECT public_key FROM "+myPrefix+"my_keys WHERE public_key = [hex]", p.TxMap["public_key_hex"])
				if err != nil {
					return utils.ErrInfo(err)
				}
				if myPublicKey != "" {
					// теперь у нас полноценный юзерский акк, и его можно апргрейдить до майнерского
					err = p.DCDB.ExecSql("UPDATE "+myPrefix+"my_table SET user_id = ?, status = 'user', notification_status = 0", NewMinerId)
					if err != nil {
						return utils.ErrInfo(err)
					}
					err = p.DCDB.ExecSql("UPDATE "+myPrefix+"my_keys SET block_id = ? WHERE public_key = [hex]", p.BlockData.BlockId, p.TxMap["public_key_hex"])
					if err != nil {
						return utils.ErrInfo(err)
					}
				}
			}
		}
	} else {
		myUserId, err := p.DCDB.Single("SELECT user_id FROM my_table")
		if err != nil {
			return utils.ErrInfo(err)
		}
		if myUserId == "" {

			// проверим, не наш ли это public_key, чтобы записать полученный user_id в my_table
			myPublicKey, err := p.DCDB.Single("SELECT public_key FROM my_keys WHERE public_key = [hex]", p.TxMap["public_key_hex"])
			if err != nil {
				return utils.ErrInfo(err)
			}
			if myPublicKey != "" {
				// теперь у нас полноценный юзерский акк, и его можно апргрейдить до майнерского
				err = p.DCDB.ExecSql("UPDATE my_table SET user_id = ?, status = 'user', notification_status = 0", NewMinerId)
				if err != nil {
					return utils.ErrInfo(err)
				}
				err = p.DCDB.ExecSql("UPDATE my_keys SET block_id = ? WHERE public_key = [hex]", p.BlockData.BlockId, p.TxMap["public_key_hex"])
				if err != nil {
					return utils.ErrInfo(err)
				}
			}
		}
	}
	// проверим, не наш ли это user_id
	myUserId, myBlockId, myPrefix, _ , err:= p.GetMyUserId(utils.BytesToInt64(p.TxMap["user_id"]))
	if err != nil {
		return utils.ErrInfo(err)
	}
	if utils.BytesToInt64(p.TxMap["user_id"]) == myUserId && myBlockId <= p.BlockData.BlockId {
		p.DCDB.ExecSql("UPDATE "+myPrefix+"my_new_users SET status ='approved', user_id = ? WHERE public_key = [hex]", NewMinerId, p.TxMap["public_key_hex"])
	}
	return nil
}

func (p *Parser) NewMinerRollback() (error) {
	// если работаем в режиме пула, то ищем тех, у кого записан такой ключ
	community, err := p.DCDB.GetCommunityUsers()
	if err != nil {
		return utils.ErrInfo(err)
	}
	if len(community) > 0 {
		for _, userId := range community {
				myPrefix := utils.Int64ToStr(userId)+"_"
				// проверим, не наш ли это public_key, чтобы записать полученный user_id в my_table
				myPublicKey, err := p.DCDB.Single("SELECT public_key FROM "+myPrefix+"my_keys WHERE public_key = [hex]", p.TxMap["public_key_hex"])
				if err != nil {
					return utils.ErrInfo(err)
				}
				if myPublicKey != "" {
					// теперь у нас полноценный юзерский акк, и его можно апргрейдить до майнерского
					err = p.DCDB.ExecSql("UPDATE "+myPrefix+"my_table SET user_id = 0, status = 'my_pending', notification_status = 0")
					if err != nil {
						return utils.ErrInfo(err)
					}
					err = p.DCDB.ExecSql("UPDATE "+myPrefix+"my_keys SET block_id = 0 WHERE block_id = ?", p.BlockData.BlockId)
					if err != nil {
						return utils.ErrInfo(err)
					}
				}
		}
	} else {
			// проверим, не наш ли это public_key
			myPublicKey, err := p.DCDB.Single("SELECT public_key FROM my_keys WHERE public_key = [hex]", p.TxMap["public_key_hex"])
			if err != nil {
				return utils.ErrInfo(err)
			}
			if myPublicKey != "" {
				err = p.DCDB.ExecSql("UPDATE my_table SET user_id = 0, status = 'my_pending', notification_status = 0")
				if err != nil {
					return utils.ErrInfo(err)
				}
				err = p.DCDB.ExecSql("UPDATE my_keys SET block_id = 0 WHERE block_id = ?", p.BlockData.BlockId)
				if err != nil {
					return utils.ErrInfo(err)
				}
			}
	}
	err = p.DCDB.ExecSql("DELETE FROM users WHERE public_key_0 = [hex]", p.TxMap["public_key_hex"])
	if err != nil {
		return utils.ErrInfo(err)
	}
	err = p.rollbackAI("users", 1)
	if err != nil {
		return utils.ErrInfo(err)
	}
	return nil
}

func (p *Parser) NewMinerRollbackFront() error {
	return p.limitRequestsRollback("new_user")
}
