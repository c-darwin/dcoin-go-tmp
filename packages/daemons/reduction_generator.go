package daemons

import (
	"errors"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	//"log"
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/consts"
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
)

/*
 * Каждые 2 недели собираем инфу о голосах за % и создаем тр-ию, которая
 * попадет в DC сеть только, если мы окажемся генератором блока
 * */
func ReductionGenerator() {
	/*defer func() {
		if r := recover(); r != nil {
			log.Error("daemon Recovered", r)
			panic(r)
		}
	}()*/


	const GoroutineName = "ReductionGenerator"
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
			if d.unlockPrintSleep(err, d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		if blockId == 0 {
			if d.unlockPrintSleep(errors.New("blockId == 0"), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}

		_, _, myMinerId, _, _, _, err := d.TestBlock()
		if err != nil {
			if d.unlockPrintSleep(err, d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		// а майнер ли я ?
		if myMinerId == 0 {
			if d.unlockPrintSleep(err, d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		variables, err := d.GetAllVariables()
		curTime := utils.Time()
		var reductionType string
		var reductionCurrencyId int
		var reductionPct float64

		// ===== ручное урезание денежной массы
		// получаем кол-во обещанных сумм у разных юзеров по каждой валюте. start_time есть только у тех, у кого статус mining/repaid
		promisedAmount, err := d.GetMap(`
				SELECT currency_id, count(user_id) as count
				FROM (
						SELECT currency_id, user_id
						FROM promised_amount
						WHERE start_time < ?  AND
									 del_block_id = 0 AND
									 del_mining_block_id = 0 AND
									 status IN ('mining', 'repaid')
						GROUP BY  user_id, currency_id
						) as t1
				GROUP BY  currency_id`, "currency_id", "count", (curTime - variables.Int64["min_hold_time_promise_amount"]))
		if err != nil {
			if d.unlockPrintSleep(err, d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		log.Info("%v", "promisedAmount", promisedAmount)

		// берем все голоса юзеров
		rows, err := d.Query(d.FormatQuery(`
				SELECT currency_id,
				  		  pct,
						  count(currency_id) as votes
				FROM votes_reduction
				WHERE time > ?
				GROUP BY  currency_id, pct
				`), curTime-variables.Int64["reduction_period"])
		if err != nil {
			if d.unlockPrintSleep(err, d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		for rows.Next() {
			var votes int64
			var pct float64
			var currency_id string
			err = rows.Scan(&currency_id, &pct, &votes)
			if err != nil {
				rows.Close()
				if d.unlockPrintSleep(err, d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
			if len(promisedAmount[currency_id]) == 0 || promisedAmount[currency_id] == "0" {
				continue
			}
			// если голосов за урезание > 50% от числа всех держателей данной валюты
			if votes >= utils.StrToInt64(promisedAmount[currency_id])/2 {
				// проверим, прошло ли 2 недели с последнего урезания
				reductionTime, err := d.Single("SELECT max(time) FROM reduction WHERE currency_id  =  ? AND type  =  'manual'", currency_id).Int64()
				if err != nil {
					rows.Close()
					if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
					continue BEGIN
				}
				if curTime-reductionTime > variables.Int64["reduction_period"] {
					reductionCurrencyId = utils.StrToInt(currency_id)
					reductionPct = pct
					reductionType = "manual"
					log.Info("%v", "reductionCurrencyId", reductionCurrencyId, "reductionPct", reductionPct, "reductionType", reductionType)
					break
				}
			}
		}
		rows.Close()

		// =======  авто-урезание денежной массы из-за малого объема обещанных сумм

		// получаем кол-во DC на кошельках
		sumWallets_, err := d.GetMap("SELECT currency_id, sum(amount) as sum_amount FROM wallets GROUP BY currency_id", "currency_id", "sum_amount")
		if err != nil {
			if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		sumWallets := make(map[int]float64)
		for currencyId, amount := range sumWallets_ {
			sumWallets[utils.StrToInt(currencyId)] = utils.StrToFloat64(amount)
		}

		// получаем кол-во TDC на обещанных суммах, плюсуем к тому, что на кошельках
		sumTdc, err := d.GetMap("SELECT currency_id, sum(tdc_amount) as sum_amount FROM promised_amount GROUP BY currency_id", "currency_id", "sum_amount")
		if err != nil {
			if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
			continue BEGIN
		}

		for currencyId, amount := range sumTdc {
			currencyIdInt := utils.StrToInt(currencyId)
			if sumWallets[currencyIdInt] == 0 {
				sumWallets[currencyIdInt] = utils.StrToFloat64(amount)
			} else {
				sumWallets[currencyIdInt] += utils.StrToFloat64(amount)
			}
		}

		// получаем суммы обещанных сумм
		sumPromisedAmount, err := d.GetMap(`
				SELECT currency_id,
					   		sum(amount) as sum_amount
				FROM promised_amount
				WHERE status = 'mining' AND
							 del_block_id = 0 AND
							 del_mining_block_id = 0 AND
							  (cash_request_out_time = 0 OR cash_request_out_time > ?)
				GROUP BY currency_id
				`, "currency_id", "sum_amount", curTime-variables.Int64["cash_request_time"])
		if err != nil {
			if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
			continue BEGIN
		}

		if len(sumWallets) > 0 {
			for currencyId, sumAmount := range sumWallets {
				//недопустимо для WOC
				if currencyId == 1 {
					continue
				}
				reductionTime, err := d.Single("SELECT max(time) FROM reduction WHERE currency_id  =  ? AND type  =  'auto'", currencyId).Int64()
				if err != nil {
					if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
					continue BEGIN
				}
				// прошло ли 48 часов
				if curTime-reductionTime <= consts.AUTO_REDUCTION_PERIOD {
					continue
				}

				// если обещанных сумм менее чем 100% от объема DC на кошельках, то запускаем урезание
				if utils.StrToFloat64(sumPromisedAmount[utils.IntToStr(currencyId)]) < sumAmount*consts.AUTO_REDUCTION_PROMISED_AMOUNT_PCT {

					// проверим, есть ли хотя бы 1000 юзеров, у которых на кошелках есть или была данная валюты
					countUsers, err := d.Single("SELECT count(user_id) FROM wallets WHERE currency_id  =  ?", currencyId).Int64()
					if err != nil {
						if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
						continue BEGIN
					}
					if countUsers >= consts.AUTO_REDUCTION_PROMISED_AMOUNT_MIN {
						reductionCurrencyId = currencyId
						reductionPct = consts.AUTO_REDUCTION_PCT
						reductionType = "promised_amount"
						break
					}
				}

			}
		}
		if reductionCurrencyId > 0 && reductionPct > 0 {

			_, myUserId, _, _, _, _, err := d.TestBlock()
			forSign := fmt.Sprintf("%v,%v,%v,%v,%v,%v", utils.TypeInt("NewReduction"), curTime, myUserId, reductionCurrencyId, reductionPct, reductionType)
			binSign, err := d.GetBinSign(forSign, myUserId)
			if err != nil {
				if d.unlockPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
			data := utils.DecToBin(utils.TypeInt("NewReduction"), 1)
			data = append(data, utils.DecToBin(curTime, 4)...)
			data = append(data, utils.EncodeLengthPlusData(utils.Int64ToByte(myUserId))...)
			data = append(data, utils.EncodeLengthPlusData(utils.Int64ToByte(int64(reductionCurrencyId)))...)
			data = append(data, utils.EncodeLengthPlusData(utils.Float64ToBytes(reductionPct))...)
			data = append(data, utils.EncodeLengthPlusData([]byte(reductionType))...)
			data = append(data, utils.EncodeLengthPlusData([]byte(binSign))...)

			err = d.InsertReplaceTxInQueue(data)
			if err != nil {
				if d.unlockPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
				continue BEGIN
			}

			// и не закрывая main_lock переводим нашу тр-ию в verified=1, откатив все несовместимые тр-ии
			// таким образом у нас будут в блоке только актуальные голоса.
			// а если придет другой блок и станет verified=0, то эта тр-ия просто удалится.

			p := new(dcparser.Parser)
			err = p.TxParser(utils.HexToBin(utils.Md5(data)), data, true)
			if err != nil {
				if d.unlockPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
		}
		d.dbUnlock()

		if d.dSleep(d.sleepTime) {
			break BEGIN
		}
	}

}
