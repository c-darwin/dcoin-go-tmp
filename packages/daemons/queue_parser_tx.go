package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

/*
 * Берем тр-ии из очереди и обрабатываем
 * */

func QueueParserTx() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("daemon Recovered", r)
			panic(r)
		}
	}()

	if utils.Mobile() {
		sleepTime = 180
	} else {
		sleepTime = 1
	}
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
			if d.dPrintSleep(err, sleepTime) {	break BEGIN }
			continue BEGIN
		}

		blockId, err := d.GetBlockId()
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), sleepTime) {	break BEGIN }
			continue BEGIN
		}
		if blockId == 0 {
			if d.unlockPrintSleep(utils.ErrInfo("blockId == 0"), sleepTime) {	break BEGIN }
			continue BEGIN
		}

		// чистим зацикленные
		err = d.ExecSql("DELETE FROM transactions WHERE verified = 0 AND used = 0 AND counter > 10")
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), sleepTime) {	break BEGIN }
			continue BEGIN
		}

		p := new(dcparser.Parser)
		p.DCDB = d.DCDB
		err = p.AllTxParser()
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), sleepTime) {	break BEGIN }
			continue BEGIN
		}

		d.dbUnlock()

		if d.dSleep(sleepTime) {
			break BEGIN
		}
	}

}
