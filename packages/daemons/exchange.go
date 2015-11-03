package daemons

import (
	//"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"crypto/rand"
	"regexp"
)

func Exchange() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("daemon Recovered", r)
			panic(r)
		}
	}()


	const GoroutineName = "Exchange"
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


BEGIN:
	for {
		log.Info(GoroutineName)
		MonitorDaemonCh <- []string{GoroutineName, utils.Int64ToStr(utils.Time())}

		// проверим, не нужно ли нам выйти из цикла
		if CheckDaemonsRestart() {
			break BEGIN
		}


		blockId, err := d.GetConfirmedBlockId()
		if err != nil {
			if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}

		var myPrefix string
		community, err := d.GetCommunityUsers()
		if len(community) > 0 {
			adminUserId, err := d.GetPoolAdminUserId()
			if err != nil {
				if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
			myPrefix = utils.Int64ToStr(adminUserId)+"_"
		} else {
			userId, err := d.Single(`SELECT user_id FROM my_table`).String()
			if err != nil {
				if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
			myPrefix = userId+"_"
		}

		eConfig, err := d.GetMap(`SELECT * FROM e_config`, "name", "value")
		if err != nil {
			if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		confirmations := utils.StrToInt64(eConfig["confirmations"])

		// все валюты, с которыми работаем
		currencyList, err := utils.EGetCurrencyList()
		if err != nil {
			if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}


		// ++++++++++++ reduction ++++++++++++

		// максимальный номер блока для процентов. Чтобы брать только новые
		maxReductionBlock, err := d.Single(`SELECT max(block_id) FROM e_reduction`).Int64()
		if err != nil {
			if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}

		// если есть свежая reduction, то нужно остановить торги
		reduction, err := d.Single(`SELECT block_id FROM reduction WHERE block_id > ? and pct > 0`, maxReductionBlock).Int64()
		if err != nil {
			if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		if reduction > 0 {
			err = d.ExecSql(`INSERT INTO e_reduction_lock (time) VALUES (?)`, utils.Time())
			if err != nil {
				log.Error("%v", utils.ErrInfo(err))
			}
		}

		// если уже прошло 10 блоков с момента обнаружения reduction, то производим сокращение объема монет
		rows, err := d.Query(d.FormatQuery(`SELECT pct, currency_id, time, block_id FROM reduction WHERE block_id > ? AND
	block_id < ?`), maxReductionBlock, blockId-confirmations)
		if err != nil {
			if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		for rows.Next() {
			var pct float64
			var currencyId, rTime, blockId int64
			err = rows.Scan(&pct, &currencyId, &rTime, &blockId)
			if err != nil {
				rows.Close()
				if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
			//	$k = (100-$row['pct'])/100;
			k := (100-pct)/100
			// уменьшаем все средства на счетах
			err = d.ExecSql(`UPDATE e_wallets SET amount = amount * ? WHERE currency_id = ?`, k, currencyId)

			// уменьшаем все средства на вывод
			err = d.ExecSql(`UPDATE e_withdraw SET amount = amount * ? WHERE currency_id = ? AND close_time = 0`, k, currencyId)

			// уменьшаем все ордеры на продажу
			err = d.ExecSql(`UPDATE e_orders SET amount = amount * ?, begin_amount = begin_amount * ? WHERE sell_currency_id = ?`, k, k, currencyId)

			err = d.ExecSql(`INSERT INTO e_reduction (
					time,
					block_id,
					currency_id,
					pct
				)
				VALUES (
					?,
					?,
					?,
					?
				)`, rTime, blockId, currencyId, pct)
			if err != nil {
				rows.Close()
				if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
		}
		rows.Close()

		// ++++++++++++ DC ++++++++++++
		/*
		 * Важно! отключать в кроне при обнулении данных в БД
		 *
		 * По крону даем 5 минут, чтобы скрипт успел всё сделать
		 * 1. Получаем инфу о входящих переводах и начисляем их на счета юзеров
		 * 2. Обновляем проценты
		 * 3. Чистим кол-во отправленных смс-ок
		 * */
		nodePrivateKey, err := d.GetNodePrivateKey(myPrefix)
		if err != nil {
			if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		/*
		 * Получаем инфу о входящих переводах и начисляем их на счета юзеров
		 * */

		// если всё остановлено из-за найденного блока с reduction, то входящие переводы не обрабатываем
		reductionLock, err := utils.EGetReductionLock()
		if err != nil {
			if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		if reductionLock > 0 {
			if d.dPrintSleep(utils.ErrInfo("reductionLock"), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}

		rows, err = d.Query(d.FormatQuery(`SELECT amount, id, block_id, type_id, currency_id, to_user_id, tx_time, comment FROM my_dc_transactions WHERE type = 'from_user' AND block_id < ? AND merchant_checked = 0 AND status = 'approved' ORDER BY id DESC`), blockId-confirmations)
		if err != nil {
			if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		for rows.Next() {
			var amount float64
			var id, blockId, typeId, currencyId, toUserId, txTime int64
			var comment string
			err = rows.Scan(&amount, &id, &blockId, &typeId, &currencyId, &toUserId, &txTime, &comment)
			if err != nil {
				rows.Close()
				if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {    break BEGIN }
				continue BEGIN
			}
			// отметим merchant_checked=1, чтобы больше не брать эту тр-ию
			err = d.ExecSql(`UPDATE my_dc_transactions SET merchant_checked = 1 WHERE id = ?`, id)
			if err != nil {
				rows.Close()
				if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {    break BEGIN }
				continue BEGIN
			}

			// вначале нужно проверить, точно ли есть такой перевод в блоке
			binaryData, err := d.Single(`SELECT data FROM block_chain WHERE id = ?`, blockId).Bytes()
			if err != nil {
				rows.Close()
				if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {    break BEGIN }
				continue BEGIN
			}
			p := new(dcparser.Parser)
			p.BinaryData = binaryData
			p.ParseDataLite()
			for _, txMap := range p.TxMapsArr {
				// пропускаем все ненужные тр-ии
				if txMap.Int64["type"] != utils.TypeInt("SendDc") {
					continue
				}

				// сравнение данных из таблы my_dc_transactions с тем, что в блоке
				if txMap.Int64["user_id"] == typeId && txMap.Int64["currency_id"] == currencyId && txMap.Money["amount"] == amount && txMap.String["comment"] == comment && txMap.Int64["to_user_id"] == toUserId {

					// расшифруем коммент
					block, _ := pem.Decode([]byte(nodePrivateKey))
					if block == nil || block.Type != "RSA PRIVATE KEY" {
						rows.Close()
						if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
						continue BEGIN
					}
					private_key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
					if err != nil {
						rows.Close()
						if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
						continue BEGIN
					}
					decryptedComment_, err := rsa.DecryptPKCS1v15(rand.Reader, private_key, []byte(comment))
					if err != nil {
						rows.Close()
						if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
						continue BEGIN
					}
					decryptedComment := string(decryptedComment_)
					// запишем расшифрованный коммент, чтобы потом можно было найти перевод в ручном режиме
					err = d.ExecSql("UPDATE "+myPrefix+"my_dc_transactions SET comment = ?, comment_status = 'decrypted' WHERE id = ?", decryptedComment, id)
					if err != nil {
						rows.Close()
						if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
						continue BEGIN
					}

					// возможно юзер сделал перевод в той валюте, которая у нас на бирже еще не торгуется
					if len(currencyList[currencyId]) == 0 {
						continue
					}

					// возможно, что чуть раньше было reduction, а это значит, что все тр-ии,
					// которые мы ещё не обработали и которые были До блока с reduction нужно принимать с учетом reduction
					// т.к. средства на нашем счете уже урезались, а  вот те, что после reduction - остались в том виде, в котором пришли
					lastReduction, err := d.OneRow("SELECT block_id, pct FROM reduction WHERE currency_id  = ? ORDER BY block_id", currencyId).Int64()
					if err != nil {
						rows.Close()
						if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
						continue BEGIN
					}
					if blockId <= lastReduction["block_id"] {
						// сумму с учетом reduction
						k0 := (100 - lastReduction["pct"]) / 100
						amount = amount * float64(k0)
					}
					// начисляем средства на счет того, чей id указан в комменте
					r, _ := regexp.Compile(`(?i)\s*#\s*([0-9]+)\s*`)
					user := r.FindStringSubmatch(decryptedComment)
					if len(user) == 0 {
						continue
					}
					uid := utils.StrToInt64(user[1])
					userAmountAndProfit := utils.EUserAmountAndProfit(uid, currencyId)
					newAmount_ := userAmountAndProfit + amount

					utils.UpdEWallet(uid, currencyId, utils.Time(), newAmount_)

					// для отчетности запишем в историю
					err = d.ExecSql(`
						INSERT INTO e_adding_funds (
							user_id,
							currency_id,
							time,
							amount
						)
						VALUES (
							?,
							?,
							?,
							?
						)`, uid, currencyId, txTime, amount)
					if err != nil {
						rows.Close()
						if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
						continue BEGIN
					}
				}
			}
		}
		rows.Close()

		/*
		 * Обновляем проценты
		 * */
		// максимальный номер блока для процентов. Чтобы брать только новые
		maxPctBlock, err := d.Single(`SELECT max(block_id) FROM e_user_pct`).Int64()

		rows, err = d.Query(d.FormatQuery(`SELECT time, block_id, currency_id, user FROM pct WHERE  block_id < ? AND block_id > ?`), blockId-confirmations, maxPctBlock)
		if err != nil {
			if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		for rows.Next() {
			var pct float64
			var pTime, blockId, currencyId int64
			err = rows.Scan(&pTime, &blockId, &currencyId, &pct)
			if err != nil {
				rows.Close()
				if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {    break BEGIN }
				continue BEGIN
			}
			d.ExecSql(`
				INSERT INTO e_user_pct (
					time,
					block_id,
					currency_id,
					pct
				)
				VALUES (
					?,
					?,
					?,
					?
				)`, pTime, blockId, currencyId, pct)
		}
		rows.Close()

		if d.dSleep(d.sleepTime) {
			break BEGIN
		}
	}
}
