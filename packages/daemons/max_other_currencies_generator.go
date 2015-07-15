package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"encoding/json"
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
)

/*
 * Каждые 2 недели собираем инфу о голосах за max_other_currencies и создаем тр-ию, которая
 * попадет в DC сеть только, если мы окажемся генератором блока
 * */

func MaxOtherCurrenciesGenerator() string {

	const GoroutineName = "MaxOtherCurrenciesGenerator"
	db := utils.DbConnect(configIni)
	db.GoroutineName = GoroutineName
	db.CheckInstall()
BEGIN:
	for {

		err := db.DbLock()
		if err != nil {
			db.PrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		blockId, err := db.GetBlockId()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		if blockId == 0 {
			db.UnlockPrintSleep(utils.ErrInfo("blockId == 0"), 60)
			continue BEGIN
		}

		_, _, myMinerId, _, _, _, err := db.TestBlock();
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		// а майнер ли я ?
		if myMinerId == 0 {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		variables, err := db.GetAllVariables()
		curTime := utils.Time()

		totalCountCurrencies, err := db.GetCountCurrencies()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		// проверим, прошло ли 2 недели с момента последнего обновления
		pctTime, err := db.Single("SELECT max(time) FROM max_other_currencies_time").Int64()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		if curTime - pctTime <= variables.Int64["new_max_other_currencies"] {
			db.UnlockPrintSleep(utils.ErrInfo("14 day error"), 60)
			continue BEGIN
		}

		// берем все голоса
		maxOtherCurrenciesVotes := make(map[int64][]map[int64]int64)
		rows, err := db.Query("SELECT currency_id, count, count(user_id) as votes FROM votes_max_other_currencies GROUP BY currency_id, count ORDER BY currency_id, count ASC")
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		defer rows.Close()
		for rows.Next() {
			var currency_id, count, votes int64
			err = rows.Scan(&currency_id, &count, &votes)
			if err!= nil {
				db.UnlockPrintSleep(utils.ErrInfo(err), 60)
				continue BEGIN
			}
			maxOtherCurrenciesVotes[currency_id] = append(maxOtherCurrenciesVotes[currency_id], map[int64]int64{count:votes})
		}

		newMaxOtherCurrenciesVotes := make(map[string]int64)
		for currencyId, countAndVotes := range maxOtherCurrenciesVotes {
			newMaxOtherCurrenciesVotes[utils.Int64ToStr(currencyId)] = utils.GetMaxVote(countAndVotes, 0, totalCountCurrencies, 10)
		}

		jsonData, err := json.Marshal(newMaxOtherCurrenciesVotes)

		_, myUserId, _, _, _, _, err := db.TestBlock();
		forSign := fmt.Sprintf("%v,%v,%v,%v,%v,%v", utils.TypeInt("NewMaxOtherCurrencies"), curTime, myUserId, jsonData)
		binSign, err := db.GetBinSign(forSign, myUserId)
		if err!= nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		data := utils.DecToBin(utils.TypeInt("NewMaxOtherCurrencies"), 1)
		data = append(data, utils.DecToBin(curTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(utils.Int64ToByte(myUserId))...)
		data = append(data, utils.EncodeLengthPlusData(jsonData)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(binSign))...)

		err = db.InsertReplaceTxInQueue(data)
		if err!= nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}

		p := new(dcparser.Parser)
		p.DCDB = db
		err = p.TxParser(data, utils.HexToBin(utils.Md5(data)), true)
		if err != nil {
			db.PrintSleep(err, 60)
			continue BEGIN
		}


		db.DbUnlock()
		utils.Sleep(60)
	}
	return ""
}


