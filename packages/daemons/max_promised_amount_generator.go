package daemons

import (
	"encoding/json"
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

/*
 * Каждые 2 недели собираем инфу о голосах за max_promised_amount и создаем тр-ию, которая
 * попадет в DC сеть только, если мы окажемся генератором блока
 * */

func MaxPromisedAmountGenerator() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("daemon Recovered", r)
			panic(r)
		}
	}()


	const GoroutineName = "MaxPromisedAmountGenerator"
	d := new(daemon)
	d.DCDB = DbConnect()
	if d.DCDB == nil {
		return
	}
	d.goRoutineName = GoroutineName
	if utils.Mobile() {
		d.sleepTime = 3600
	} else {
		d.sleepTime = 60
	}
	if !d.CheckInstall(DaemonCh, AnswerDaemonCh) {
		return
	}
	d.DCDB = DbConnect()
	if d.DCDB == nil {
		return
	}

	err = d.notMinerSetSleepTime(1800)
	if err != nil {
		log.Error("%v", err)
		return
	}

BEGIN:
	for {
		log.Info(GoroutineName)
		MonitorDaemonCh <- []string{GoroutineName, utils.Int64ToStr(utils.Time())}

		// проверим, не нужно ли нам выйти из цикла
		if CheckDaemonsRestart() {
			break BEGIN
		}

		err, restart := d.dbLock()
		if restart {
			break BEGIN
		}
		if err != nil {
			if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
			continue BEGIN
		}

		blockId, err := d.GetBlockId()
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		if blockId == 0 {
			d.unlockPrintSleep(utils.ErrInfo("blockId == 0"), 1)
			continue BEGIN
		}

		_, _, myMinerId, _, _, _, err := d.TestBlock()
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		// а майнер ли я ?
		if myMinerId == 0 {
			if d.unlockPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		variables, err := d.GetAllVariables()
		curTime := utils.Time()

		totalCountCurrencies, err := d.GetCountCurrencies()
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		// проверим, прошло ли 2 недели с момента последнего обновления
		pctTime, err := d.Single("SELECT max(time) FROM max_promised_amounts").Int64()
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		if curTime-pctTime <= variables.Int64["new_max_promised_amount"] {
			d.unlockPrintSleep(utils.ErrInfo("14 day error"), 1)
			continue BEGIN
		}

		// берем все голоса
		maxPromisedAmountVotes := make(map[int64][]map[int64]int64)
		rows, err := d.Query("SELECT currency_id, amount, count(user_id) as votes FROM votes_max_promised_amount GROUP BY currency_id, amount ORDER BY currency_id, amount ASC")
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		for rows.Next() {
			var currency_id, amount, votes int64
			err = rows.Scan(&currency_id, &amount, &votes)
			if err != nil {
				rows.Close()
				if d.unlockPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
			maxPromisedAmountVotes[currency_id] = append(maxPromisedAmountVotes[currency_id], map[int64]int64{amount: votes})
			//fmt.Println("currency_id", currency_id)
		}
		rows.Close()

		NewMaxPromisedAmountsVotes := make(map[string]int64)
		for currencyId, amountsAndVotes := range maxPromisedAmountVotes {
			NewMaxPromisedAmountsVotes[utils.Int64ToStr(currencyId)] = utils.GetMaxVote(amountsAndVotes, 0, totalCountCurrencies, 10)
		}

		jsonData, err := json.Marshal(NewMaxPromisedAmountsVotes)
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}

		_, myUserId, _, _, _, _, err := d.TestBlock()
		forSign := fmt.Sprintf("%v,%v,%v,%v,%v,%v", utils.TypeInt("NewMaxPromisedAmounts"), curTime, myUserId, jsonData)
		binSign, err := d.GetBinSign(forSign, myUserId)
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		data := utils.DecToBin(utils.TypeInt("NewMaxPromisedAmounts"), 1)
		data = append(data, utils.DecToBin(curTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(utils.Int64ToByte(myUserId))...)
		data = append(data, utils.EncodeLengthPlusData(jsonData)...)
		data = append(data, utils.EncodeLengthPlusData([]byte(binSign))...)

		err = d.InsertReplaceTxInQueue(data)
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}

		p := new(dcparser.Parser)
		p.DCDB = d.DCDB
		err = p.TxParser(utils.HexToBin(utils.Md5(data)), data, true)
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}

		d.dbUnlock()

		if d.dSleep(d.sleepTime) {
			break BEGIN
		}
	}

}
