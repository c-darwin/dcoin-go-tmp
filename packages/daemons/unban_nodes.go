package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)


func UnbanNodes() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("daemon Recovered", r)
			panic(r)
		}
	}()

	GoroutineName := "UnbanNodes"
	d := new(daemon)
	d.DCDB = DbConnect(GoroutineName)
	if d.DCDB == nil {
		return
	}
	d.goRoutineName = GoroutineName
	d.sleepTime = 3600
	if !d.CheckInstall(DaemonCh, AnswerDaemonCh, GoroutineName) {
		return
	}
	d.DCDB = DbConnect(GoroutineName)
	if d.DCDB == nil {
		return
	}

BEGIN:
	for {
		log.Info(GoroutineName)
		MonitorDaemonCh <- []string{GoroutineName, utils.Int64ToStr(utils.Time())}

		// проверим, не нужно ли нам выйти из цикла
		if CheckDaemonsRestart(GoroutineName) {
			break BEGIN
		}

		err = d.ExecSql("DELETE FROM nodes_ban")
		if err != nil {
			if (d.dPrintSleep(err, d.sleepTime)) {
				break BEGIN
			}
			continue BEGIN
		}

		if d.dSleep(d.sleepTime) {
			break BEGIN
		}
	}
	log.Debug("break BEGIN %v", GoroutineName)
}
