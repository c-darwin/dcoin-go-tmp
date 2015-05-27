package dcparser

import (
	"fmt"
	"utils"
	"encoding/json"
	//"regexp"
	//"math"
	//"strings"
//	"os"
	//"time"
	//"strings"
	//"sort"
//	"time"
	//"consts"
	"log"
)


// Эту транзакцию имеет право генерить только нод, который генерит данный блок
// подписана нодовским ключом
func (p *Parser) NewPctInit() (error) {
	var err error
	var fields []string
	fields = []string {"new_pct", "sign"}
	p.TxMap, err = p.GetTxMap(fields);
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}

type newPctType struct {
	Currency map[string]map[string]string
	Referral map[string]string
}

func (p *Parser) NewPctFront() (error) {

	err := p.generalCheck()
	if err != nil {
		return p.ErrInfo(err)
	}

	// является ли данный юзер майнером
	err = p.checkMiner(p.TxUserID)
	if err != nil {
		return p.ErrInfo(err)
	}

	nodePublicKey, err := p.GetNodePublicKey(p.TxUserID)
	if err != nil {
		return p.ErrInfo(err)
	}
	if len(nodePublicKey) == 0 {
		return p.ErrInfo("incorrect user_id")
	}

	newPctCurrency := make(map[string]map[string]string)
	// раньше не было рефских
	if p.BlockData!=nil && p.BlockData.BlockId<=77951 {
		err = json.Unmarshal([]byte(p.TxMap["new_pct"]), &newPctCurrency)
		if err != nil {
			return p.ErrInfo(err)
		}
	} else {
		newPctTx:=new(newPctType)
		err = json.Unmarshal([]byte(p.TxMap["new_pct"]), &newPctTx)
		if err != nil {
			return p.ErrInfo(err)
		}
		if newPctTx.Referral==nil {
			return p.ErrInfo("!Referral")
		}
		newPctCurrency = newPctTx.Currency
	}
	if len(newPctCurrency) == 0 {
		return p.ErrInfo("!newPctCurrency")
	}

	// проверим, верно ли указаны ID валют
	currencyIdsSql := ""
	countCurrency := 0
	for id := range newPctCurrency {
		currencyIdsSql+=id+","
		countCurrency++
	}
	currencyIdsSql = currencyIdsSql[0:len(currencyIdsSql)-1]
	count, err := p.Single("SELECT count(id) FROM currency WHERE id IN ("+currencyIdsSql+")").Int()
	if err != nil {
		return p.ErrInfo(err)
	}
	if count != countCurrency {
		return p.ErrInfo("count_currency")
	}

	forSign := fmt.Sprintf("%s,%s,%s,%s", p.TxMap["type"], p.TxMap["time"], p.TxMap["user_id"], p.TxMap["new_pct"])
	CheckSignResult, err := utils.CheckSign([][]byte{nodePublicKey}, forSign, p.TxMap["sign"], true);
	if err != nil {
		return p.ErrInfo(err)
	}
	if !CheckSignResult {
		return p.ErrInfo("incorrect sign")
	}

	// проверим, прошло ли 2 недели с момента последнего обновления pct
	pctTime, err := p.Single("SELECT max(time) FROM pct").Int64()
	if err != nil {
		return p.ErrInfo(err)
	}
	if p.TxTime - pctTime <= p.Variables.Int64["new_pct_period"] {
		return p.ErrInfo(fmt.Sprintf("14 days error %d - %d <= %d", p.TxTime, pctTime, p.Variables.Int64["new_pct_period"] ))
	}
	// берем все голоса miner_pct
	pctVotes := make(map[int64]map[string]map[string]int64)
	rows, err := p.Query("SELECT currency_id, pct, count(user_id) as votes FROM votes_miner_pct GROUP BY currency_id, pct")
	if err != nil {
		return p.ErrInfo(err)
	}
	defer rows.Close()
	for rows.Next() {
		var currency_id, votes int64
		var pct string
		err = rows.Scan(&currency_id, &pct, &votes)
		if err!= nil {
			return p.ErrInfo(err)
		}
		log.Println("currency_id", currency_id, "pct", pct, "votes", votes)
		if len(pctVotes[currency_id]) == 0 {
			pctVotes[currency_id] = make(map[string]map[string]int64)
		}
		if len(pctVotes[currency_id]["miner_pct"]) == 0 {
			pctVotes[currency_id]["miner_pct"] = make(map[string]int64)
		}
		pctVotes[currency_id]["miner_pct"][pct] = votes
	}

	// берем все голоса user_pct
	rows, err = p.Query("SELECT currency_id, pct, count(user_id) as votes FROM votes_user_pct GROUP BY currency_id, pct")
	if err != nil {
		return p.ErrInfo(err)
	}
	defer rows.Close()
	for rows.Next() {
		var currency_id, votes int64
		var pct string
		err = rows.Scan(&currency_id, &pct, &votes)
		if err!= nil {
			return p.ErrInfo(err)
		}
		log.Println("currency_id", currency_id, "pct", pct, "votes", votes)
		if len(pctVotes[currency_id]) == 0 {
			pctVotes[currency_id] = make(map[string]map[string]int64)
		}
		if len(pctVotes[currency_id]["user_pct"]) == 0 {
			pctVotes[currency_id]["user_pct"] = make(map[string]int64)
		}
		pctVotes[currency_id]["user_pct"][pct] = votes
	}

	newPct := make(map[string]map[string]map[string]string)
	newPct["currency"] = make(map[string]map[string]string)
	var userMaxKey int64
	PctArray := utils.GetPctArray()

	log.Println("pctVotes", pctVotes)
	for currencyId, data := range pctVotes {

		currencyIdStr := utils.Int64ToStr(currencyId)
		// определяем % для майнеров
		pctArr := utils.MakePctArray(data["miner_pct"])
		key := utils.GetMaxVote(pctArr, 0, 390, 100)
		if len(newPct["currency"][currencyIdStr]) == 0{
			newPct["currency"][currencyIdStr] = make(map[string]string)
		}
		newPct["currency"][currencyIdStr]["miner_pct"] = utils.GetPctValue(key)

		// определяем % для юзеров
		pctArr = utils.MakePctArray(data["user_pct"])
		// раньше не было завимости юзерского % от майнерского
		if p.BlockData!=nil && p.BlockData.BlockId<=95263 {
			userMaxKey = 390
		} else {
			pctY := utils.ArraySearch(newPct["currency"][currencyIdStr]["miner_pct"], PctArray)
			maxUserPctY := utils.Round(utils.StrToFloat64(pctY)/2, 2)
			userMaxKey = utils.FindUserPct(int(maxUserPctY))
			// отрезаем лишнее, т.к. поиск идет ровно до макимального возможного, т.е. до miner_pct/2
			pctArr = utils.DelUserPct(pctArr, userMaxKey);
		}
		key = utils.GetMaxVote(pctArr, 0, userMaxKey, 100)
		log.Println("data[user_pct]", data["user_pct"])
		log.Println("pctArr", pctArr)
		log.Println("userMaxKey", userMaxKey)
		log.Println("key", key)
		newPct["currency"][currencyIdStr]["user_pct"] = utils.GetPctValue(key)
		log.Println("user_pct", newPct["currency"][currencyIdStr]["user_pct"])
	}

	var jsonData []byte
	// раньше не было рефских
	if p.BlockData != nil && p.BlockData.BlockId <= 77951 {

		newPct_ := newPct["currency"];
		jsonData, err = json.Marshal(newPct_)
		if err!= nil {
			return p.ErrInfo(err)
		}
	} else {

		newPct_ := new(newPctType)
		newPct_.Currency = make(map[string]map[string]string)
		newPct_.Currency = newPct["currency"]
		newPct_.Referral = make(map[string]string)
		refLevels := []string{"first", "second", "third"}
		for i:=0; i<len(refLevels); i++ {
			level := refLevels[i]
			var votesReferral []map[int64]int64

			// берем все голоса
			//pctVotes := make(map[int64]map[string]map[string]int64)
			rows, err := p.Query("SELECT ?, count(user_id) as votes FROM votes_referral GROUP BY ?", level, level)
			if err != nil {
				return p.ErrInfo(err)
			}
			defer rows.Close()
			for rows.Next() {
				var level, votes int64
				err = rows.Scan(&level, &votes)
				if err!= nil {
					return p.ErrInfo(err)
				}
				votesReferral = append(votesReferral, map[int64]int64{level:votes})
			}
			newPct_.Referral[level] = utils.Int64ToStr(utils.GetMaxVote(votesReferral, 0, 30, 10))
		}
		jsonData, err = json.Marshal(newPct_)
		if err!= nil {
			return p.ErrInfo(err)
		}
	}

	if string(p.TxMap["new_pct"]) != string(jsonData) {
		return p.ErrInfo("p.TxMap[new_pct] != jsonData "+string(p.TxMap["new_pct"])+"!="+string(jsonData))
	}
	log.Println(string(jsonData))

	return nil
}


func (p *Parser) NewPct() (error) {

	newPctCurrency := make(map[string]map[string]string)
	newPctTx:=new(newPctType)
	// раньше не было рефских
	if p.BlockData.BlockId<=77951 {
		err := json.Unmarshal([]byte(p.TxMap["new_pct"]), &newPctCurrency)
		if err != nil {
			return p.ErrInfo(err)
		}
	} else {
		err := json.Unmarshal([]byte(p.TxMap["new_pct"]), &newPctTx)
		if err != nil {
			return p.ErrInfo(err)
		}
		if newPctTx.Referral==nil {
			return p.ErrInfo("!Referral")
		}
		newPctCurrency = newPctTx.Currency
	}
	for currencyId, data := range(newPctCurrency) {
		err := p.ExecSql("INSERT INTO pct ( time, currency_id, miner, user, block_id ) VALUES ( ?, ?, ?, ?, ? )", p.BlockData.Time, currencyId, data["miner_pct"], data["user_pct"], p.BlockData.BlockId)
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	if p.BlockData.BlockId > 77951 {
		err := p.selectiveLoggingAndUpd([]string{"first", "second", "third"}, []string{newPctTx.Referral["first"], newPctTx.Referral["second"], newPctTx.Referral["third"]}, "referral", []string{}, []string{})
		if err != nil {
			return p.ErrInfo(err)
		}
	}

	return nil
}

func (p *Parser) NewPctRollback() (error) {
	if p.BlockData.BlockId < 77951 {
		err := p.selectiveRollback([]string{"first", "second", "third"}, "referral", "", false)
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	return nil
}

func (p *Parser) NewPctRollbackFront() error {

	return nil
}
