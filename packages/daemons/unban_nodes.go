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
			d.PrintSleep(err, 1)
			continue BEGIN
		}
		for i:=0; i < 3600; i++ {
			if CheckDaemonsRestart() {
				utils.Sleep(1)
				break BEGIN
			}
			utils.Sleep(1)
		}
	}
}
