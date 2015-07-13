package daemons

import (
	"github.com/c-darwin/dcoin-tmp/packages/utils"
	"log"
	"github.com/c-darwin/dcoin-tmp/packages/dcparser"
)

/*
 * Берем тр-ии из очереди и обрабатываем
 * */

func QueueParserTx(configIni map[string]string) string {

	GoroutineName := "QueueParserTx"
	db := utils.DbConnect(configIni)
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

		blockId, err := db.GetBlockId()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		if blockId == 0 {
			db.UnlockPrintSleep(utils.ErrInfo("blockId == 0"), 1)
			continue BEGIN
		}

		// чистим зацикленные
		err = db.ExecSql("DELETE FROM transactions WHERE verified = 0 AND used = 0 AND counter > 10")
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		p := new(dcparser.Parser)
		err = p.AllTxParser()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		db.DbUnlock()

		utils.Sleep(1)

		log.Println("Happy end")
	}

	return ""
}
