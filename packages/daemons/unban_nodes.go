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

	sleepTime := 3600
	GoroutineName := "UnbanNodes"
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

		err = d.ExecSql("DELETE FROM nodes_ban")
		if err != nil {
			if (d.dPrintSleep(err, sleepTime)) {
				break BEGIN
			}
			continue BEGIN
		}

		if d.dSleep(sleepTime) {
			break BEGIN
		}
	}
}
