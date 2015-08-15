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

	const GoroutineName = "QueueParserTx"
	d := new(daemon)
	d.DCDB = DbConnect()
	if d.DCDB == nil {
		return
	}
	d.goRoutineName = GoroutineName
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

		err, restart := d.dbLock()
		if restart {
			break BEGIN
		}
		if err != nil {
			d.PrintSleep(err, 1)
			continue BEGIN
		}

		blockId, err := d.GetBlockId()
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		if blockId == 0 {
			d.unlockPrintSleep(utils.ErrInfo("blockId == 0"), 1)
			continue BEGIN
		}

		// чистим зацикленные
		err = d.ExecSql("DELETE FROM transactions WHERE verified = 0 AND used = 0 AND counter > 10")
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		p := new(dcparser.Parser)
		p.DCDB = d.DCDB
		err = p.AllTxParser()
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		d.dbUnlock()

		utils.Sleep(1)

		log.Info("%v", "Happy end")
	}


}
