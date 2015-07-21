package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	//"log"
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
)

/*
 * Берем тр-ии из очереди и обрабатываем
 * */

func QueueParserTx() {

	GoroutineName := "QueueParserTx"

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

		// чистим зацикленные
		err = db.ExecSql("DELETE FROM transactions WHERE verified = 0 AND used = 0 AND counter > 10")
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		p := new(dcparser.Parser)
		p.DCDB = db
		err = p.AllTxParser()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		db.DbUnlock()

		utils.Sleep(1)

		log.Info("%v", "Happy end")
	}


}
