package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	//"log"
	"errors"
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
)

/*
 * Парсим и разносим данные из queue_testblock
 * */

func QueueParserTestblock() {

	GoroutineName := "QueueParserTestblock"
	db := DbConnect()
	db.GoroutineName = GoroutineName
	db.CheckInstall()
BEGIN:
	for {

		// проверим, не нужно нам выйти, т.к. обновилась версия софта
		if db.CheckDaemonRestart() {
			utils.Sleep(1)
			break
		}

		err := db.DbLock()
		if err != nil {
			db.PrintSleep(utils.ErrInfo(err), 1)
			continue
		}

		data, err := db.OneRow("SELECT * FROM queue_testblock ORDER BY head_hash ASC").String()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue
		}
		if len(data) == 0 {
			db.UnlockPrintSleep(utils.ErrInfo(errors.New("len(data) == 0")), 1)
			continue
		}

		newBlock := []byte(data["data"])
		newHeaderHash := utils.BinToHex([]byte(data["head_hash"]))
		tx := utils.DeleteHeader(newBlock)

		// сразу можно удалять данные из таблы-очереди
		err = db.ExecSql("DELETE FROM queue_testblock WHERE head_hash = [hex]", newHeaderHash)
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(errors.New("len(data) == 0")), 1)
			continue
		}

		// прежде всего нужно проверить, а нет ли в этом блоке ошибок с несовметимыми тр-ми
		// при полной проверке, а не только фронтальной проблем с несовместимыми тр-ми не будет, т.к. там даные сразу пишутся в таблицы
		// а тут у нас данные пишутся только в log_time_
		// и сами тр-ии пишем в отдельную таблу
		p := new(dcparser.Parser)
		p.DCDB = db;
		if len(tx) > 0 {
			log.Debug("len(tx): %d", len(tx))
			for {
				log.Debug("tx: %x", tx)
				txSize := utils.DecodeLength(&tx)
				log.Debug("txSize: %d", txSize)
				// отделим одну транзакцию от списка транзакций
				txBinaryData := utils.BytesShift(&tx, txSize)
				log.Debug("txBinaryData: %x", txBinaryData)
				// проверим, нет ли несовместимых тр-ий
				fatalError, waitError, _, _, _, _ := p.ClearIncompatibleTx(txBinaryData, false)
				if len(fatalError) > 0 || len(waitError) > 0 {
					db.UnlockPrintSleep(utils.ErrInfo(errors.New(" len(fatalError) > 0 || len(waitError) > 0")), 1)
					continue BEGIN
				}

				if len(tx) == 0 {
					break
				}
			}
		}
		// откатим тр-ии тестблока, но не удаляя их, т.к. далее еще можем их вернуть
		p.RollbackTransactionsTestblock(false)

		// проверим блок, который получился с данными, которые прислал другой нод
		p.BinaryData = newBlock
		err = p.ParseDataGate(false)
		if err != nil {

			log.Info("%v", err)

			// т.к. мы откатили наши тр-ии из transactions_testblock, то теперь нужно обработать их по новой
			// получим наши транзакции в 1 бинарнике, просто для удобства

			var myTestBlockBody []byte
			transactionsTestblock, err := db.GetAll("SELECT data FROM transactions_testblock ORDER BY id ASC", -1)
			if err != nil {
				db.PrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			for _, data := range transactionsTestblock {
				myTestBlockBody = append(myTestBlockBody, []byte(data["data"])...)
			}

			if len(myTestBlockBody) > 0 {
				p.BinaryData = append(utils.DecToBin(0, 1), myTestBlockBody...)
				err = p.ParseDataGate(true)
				if err != nil {
					db.PrintSleep(utils.ErrInfo(err), 1)
					continue BEGIN
				}
			}
		} else {
			// наши тр-ии уже не актуальны, т.к. мы их откатили
			err = db.ExecSql("DELETE FROM transactions_testblock")
			if err != nil {
				db.PrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			exists, err := db.Single(`SELECT block_id FROM testblock`).Int64()
			if err != nil {
				db.PrintSleep(utils.ErrInfo(err), 1)
				return
			}
			if exists > 0 {
				// если всё нормально, то пишем в таблу testblock новые тр-ии и новые данные по юзеру их сгенерившему
				err = db.ExecSql(`
				UPDATE testblock
				SET  time = ?,
						user_id = ?,
						header_hash = [hex],
						signature = [hex],
						mrkl_root = [hex]
				`, p.BlockData.Time, p.BlockData.UserId, newHeaderHash, utils.BinToHex(p.BlockData.Sign), p.MrklRoot)
				if err != nil {
					db.PrintSleep(utils.ErrInfo(err), 1)
					continue BEGIN
				}
			}

			// и сами тр-ии пишем в отдельную таблу
			if len(tx) > 0 {
				for {
					txSize := utils.DecodeLength(&tx)
					// отчекрыжим одну транзакцию от списка транзакций
					txBinaryData := utils.BytesShift(&tx, txSize)
					// получим тип тр-ии и юзера
					txType, userId, toUserId := utils.GetTxTypeAndUserId(txBinaryData)
					md5 := utils.Md5(txBinaryData)
					dataHex := utils.BinToHex(txBinaryData)

					err = db.ExecSql("DELETE FROM transactions_testblock")
					if err != nil {
						db.PrintSleep(utils.ErrInfo(err), 1)
						continue BEGIN
					}

					err = db.ExecSql("INSERT INTO transactions_testblock (hash, data, type, user_id, third_var) VALUES ([hex], [hex], ?, ?, ?)", md5, dataHex, txType, userId, toUserId)
					if err != nil {
						db.PrintSleep(utils.ErrInfo(err), 1)
						continue BEGIN
					}

					if len(tx) == 0 {
						break
					}
				}
			}

			// удаляем всё, где хэш больше нашего
			err = db.ExecSql("DELETE FROM queue_testblock WHERE head_hash > [hex]", newHeaderHash)
			if err != nil {
				db.PrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}

			// возможно нужно откатить и тр-ии с verified=1 и used=0 из transactions
			// т.к. в transactions может быть тр-ия на удаление банкноты
			// и в transactions_testblock только что была залита такая же тр-ия
			// выходит, что блок, который будет сгенерен на основе transactions будет ошибочным
			// или при откате transactions будет сделан вычет из log_time_....
			// и выйдет что попавшая в блок тр-я из transactions_testblock попала минуя запись  log_time_....
			err = p.RollbackTransactions()
			if err != nil {
				db.PrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
		}

		db.DbUnlock()

		utils.Sleep(1)

		log.Info("%v", "Happy end")
	}


}
