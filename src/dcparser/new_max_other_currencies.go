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
	"sort"
	"time"
	"consts"
)


func (p *Parser) NewMaxOtherCurrenciesInit() (error) {
	var err error
	var fields []string
	fields = []string {"new_max_other_currencies", "sign"}
	p.TxMap, err = p.GetTxMap(fields);
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}
type vComplex struct {
	Currency map[string][]float64 `json:"currency"`
	Referral map[string]int64 `json:"referral"`
	Admin int64 `json:"admin"`
}
func (p *Parser) NewMaxOtherCurrenciesFront() (error) {

	err := p.generalCheck()
	if err != nil {
		return p.ErrInfo(err)
	}

	// является ли данный юзер майнером
	err = p.checkMiner(p.TxMap["user_id"])
	if err != nil {
		return p.ErrInfo(err)
	}

	var txTime int64
	if p.BlockData!=nil {
		txTime = p.BlockData.Time
	} else {
		txTime = time.Now().Unix() - 30
	}

	nodePublicKey, err := p.GetNodePublicKey(p.TxUserID)
	if err != nil {
		return p.ErrInfo(err)
	}
	if len(nodePublicKey) == 0 {
		return p.ErrInfo("incorrect user_id")
	}

	totalCountCurrencies := p.GetCountCurrencies()

	// проверим, верно ли указаны ID валют
	currencyList := make(map[string]int64)
	err = json.Unmarshal(p.TxMap["new_max_other_currencies"], &currencyList)
	if err != nil {
		return p.ErrInfo(err)
	}
	currencyIdsSql := ""
	var countCurrency int64
	for currencyId, count := range currencyList {
		if !utils.CheckInputData(currencyId, "int") {
			return p.ErrInfo("currencyId")
		}
		currencyIdsSql += currencyId+",";
		countCurrency++
		if count > totalCountCurrencies {
			return p.ErrInfo("count > totalCountCurrencies")
		}
	}
	currencyIdsSql := currencyIdsSql[0:len(currencyIdsSql)-1]
	if countCurrency == 0 {
		return p.ErrInfo("countCurrency")
	}
	count, err := p.Single("SELECT count(id) FROM currency WHERE id IN ("+currencyIdsSql+")").Int64()
	if err != nil {
		return p.ErrInfo(err)
	}
	if count != countCurrency {
		return p.ErrInfo("count != countCurrency")
	}

	forSign := fmt.Sprintf("%s,%s,%s,%s", p.TxMap["type"], p.TxMap["time"], p.TxMap["user_id"], p.TxMap["new_max_other_currencies"])
	CheckSignResult, err := utils.CheckSign(p.PublicKeys, forSign, p.TxMap["sign"], false);
	if err != nil {
		return p.ErrInfo(err)
	}
	if !CheckSignResult {
		return p.ErrInfo("incorrect sign")
	}

	// проверим, прошло ли 2 недели с момента последнего обновления
	pctTime, err := p.Single("SELECT max(time) FROM max_other_currencies_time").Int64()
	if err != nil {
		return p.ErrInfo(err)
	}
	if p.TxTime - pctTime <= p.Variables.Int64["new_max_other_currencies"] {
		return p.ErrInfo("14 day error")
	}

	// берем все голоса
	maxOtherCurrenciesVotes := make(map[int64]map[int64]int64)
	rows, err := p.Query("SELECT currency_id, count, count(user_id) as votes FROM votes_max_other_currencies GROUP BY currency_id, count")
	if err != nil {
		return p.ErrInfo(err)
	}
	defer rows.Close()
	if  ok := rows.Next(); ok {
		var currency_id, count, votes int64
		err = rows.Scan(&currency_id, &count, &votes)
		if err!= nil {
			return p.ErrInfo(err)
		}
		maxOtherCurrenciesVotes[currency_id] = make(map[int64]int64)
		maxOtherCurrenciesVotes[currency_id][count] = votes
	}

	newMaxOtherCurrenciesVotes := make(map[int64]int64)
	for currencyId, countAndVotes := range maxOtherCurrenciesVotes {
		newMaxOtherCurrenciesVotes[currencyId] = p.getMaxVote(countAndVotes, 0, totalCountCurrencies, 10)
	}

	jsonData, err := json.Marshal(newMaxOtherCurrenciesVotes)
	if err!= nil {
		return p.ErrInfo(err)
	}
	if p.TxMap["new_max_other_currencies"] != jsonData {
		return p.ErrInfo("p.TxMap[new_max_other_currencies] != jsonData")
	}

	return nil
}

func (p *Parser) NewMaxOtherCurrencies() (error) {

	currencyList := make(map[string]int64)
	err := json.Unmarshal(p.TxMap["new_max_other_currencies"], &currencyList)
	if err != nil {
		return p.ErrInfo(err)
	}

	var currencyIds []int
	for k := range currencyList {
		currencyIds = append(currencyIds, utils.StrToInt(k))
	}
	sort.Ints(currencyIds)
	//sort.Sort(sort.Reverse(sort.IntSlice(keys)))

	for _, currencyId := range currencyIds {
		count:=currencyList[utils.IntToStr(currencyId)]

		logData, err := p.Single("SELECT max_other_currencies, log_id FROM currency WHERE id  =  ?", currencyId).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		logId, err = p.ExecSqlGetLastInsertId("INSERT INTO log_currency ( max_other_currencies, prev_log_id ) VALUES ( ?, ? )", logData["max_other_currencies"], logData["log_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.ExecSql("UPDATE currency SET max_other_currencies = ?, log_id = ? WHERE id = ?", count, logId, currencyId)
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	err = p.ExecSql("INSERT INTO max_other_currencies_time ( time ) VALUES ( ? )", p.BlockData.Time)
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}

func (p *Parser) NewMaxOtherCurrenciesRollback() (error) {
	currencyList := make(map[string]int64)
	err := json.Unmarshal(p.TxMap["new_max_other_currencies"], &currencyList)
	if err != nil {
		return p.ErrInfo(err)
	}

	var currencyIds []int
	for k := range currencyList {
		currencyIds = append(currencyIds, utils.StrToInt(k))
	}
	sort.Sort(sort.Reverse(sort.IntSlice(currencyIds)))

	for _, currencyId := range currencyIds {

		count := currencyList[utils.IntToStr(currencyId)]
		logId, err := p.Single("SELECT log_id FROM currency WHERE id  =  ?", currencyId).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}

		logData, err := p.OneRow("SELECT max_other_currencies, prev_log_id FROM log_currency WHERE log_id  =  ?", logId).String()
		if err != nil {
			return p.ErrInfo(err)
		}

		err = p.ExecSql("UPDATE currency SET max_other_currencies = ?, log_id = ? WHERE id = ?", logData["max_other_currencies"], logData["prev_log_id"], currencyId)
		if err != nil {
			return p.ErrInfo(err)
		}

		err = p.ExecSql("DELETE FROM log_currency WHERE log_id = ?", logId)
		if err != nil {
			return p.ErrInfo(err)
		}
		p.rollbackAI("log_currency", 1)

		err = p.ExecSql("DELETE FROM max_other_currencies_time WHERE time = ?", p.BlockData.Time)
		if err != nil {
			return p.ErrInfo(err)
		}

	}
	return nil
}

func (p *Parser) NewMaxOtherCurrenciesRollbackFront() error {

	return nil
}
