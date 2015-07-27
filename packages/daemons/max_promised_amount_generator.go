package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"encoding/json"
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
)

/*
 * Каждые 2 недели собираем инфу о голосах за max_promised_amount и создаем тр-ию, которая
 * попадет в DC сеть только, если мы окажемся генератором блока
 * */

func MaxPromisedAmountGenerator() {

	const GoroutineName = "MaxPromisedAmountGenerator"

	db := DbConnect()
	if db == nil {
		return
	}
	db.GoroutineName = GoroutineName
	if !db.CheckInstall(DaemonCh, AnswerDaemonCh) {
		return
	}

	BEGIN:
	for {
		log.Info(GoroutineName)
		// проверим, не нужно ли нам выйти из цикла
		if CheckDaemonsRestart() {
			break BEGIN
		}

		err, restart := db.DbLock(DaemonCh, AnswerDaemonCh)
		if restart {
			break BEGIN
		}
		if err != nil {
			db.PrintSleep(err, 1)
			continue BEGIN
		}

		blockId, err := db.GetBlockId()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		if blockId == 0 {
			db.UnlockPrintSleep(utils.ErrInfo("blockId == 0"), 1)
			continue BEGIN
		}

		_, _, myMinerId, _, _, _, err := db.TestBlock();
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		// а майнер ли я ?
		if myMinerId == 0 {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		variables, err := db.GetAllVariables()
		curTime := utils.Time()

		totalCountCurrencies, err := db.GetCountCurrencies()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		// проверим, прошло ли 2 недели с момента последнего обновления
		pctTime, err := db.Single("SELECT max(time) FROM max_promised_amounts").Int64()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		if curTime - pctTime <= variables.Int64["new_max_promised_amount"] {
			db.UnlockPrintSleep(utils.ErrInfo("14 day error"), 1)
			continue BEGIN
		}

		// берем все голоса
		maxPromisedAmountVotes := make(map[int64][]map[int64]int64)
		rows, err := db.Query("SELECT currency_id, amount, count(user_id) as votes FROM votes_max_promised_amount GROUP BY currency_id, amount ORDER BY currency_id, amount ASC")
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		defer rows.Close()
		for rows.Next() {
			var currency_id, amount, votes int64
			err = rows.Scan(&currency_id, &amount, &votes)
			if err!= nil {
				db.UnlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			maxPromisedAmountVotes[currency_id] = append(maxPromisedAmountVotes[currency_id], map[int64]int64{amount:votes})
			fmt.Println("currency_id", currency_id)
		}

		NewMaxPromisedAmountsVotes := make(map[string]int64)
		for currencyId, amountsAndVotes := range maxPromisedAmountVotes {
			NewMaxPromisedAmountsVotes[utils.Int64ToStr(currencyId)] = utils.GetMaxVote(amountsAndVotes, 0, totalCountCurrencies, 10)
		}

		jsonData, err := json.Marshal(NewMaxPromisedAmountsVotes)
		if err!= nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		_, myUserId, _, _, _, _, err := db.TestBlock();
		forSign := fmt.Sprintf("%v,%v,%v,%v,%v,%v", utils.TypeInt("NewMaxPromisedAmounts"), curTime, myUserId, jsonData)
		binSign, err := db.GetBinSign(forSign, myUserId)
		if err!= nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		data := utils.DecToBin(utils.TypeInt("NewMaxPromisedAmounts"), 1)
		data = append(data, utils.DecToBin(curTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(utils.Int64ToByte(myUserId))...)
		data = append(data, utils.EncodeLengthPlusData(jsonData)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(binSign))...)

		err = db.InsertReplaceTxInQueue(data)
		if err!= nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		p := new(dcparser.Parser)
		p.DCDB = db
		err = p.TxParser(utils.HexToBin(utils.Md5(data)), data, true)
		if err != nil {
			db.PrintSleep(err, 1)
			continue BEGIN
		}

		db.DbUnlock()
		utils.Sleep(1)
	}

}


