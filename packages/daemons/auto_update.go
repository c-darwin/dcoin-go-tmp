package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)


func AutoUpdate() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("daemon Recovered", r)
			panic(r)
		}
	}()

	GoroutineName := "AutoUpdate"
	d := new(daemon)
	d.DCDB = DbConnect()
	if d.DCDB == nil {
		return
	}
	d.goRoutineName = GoroutineName
	d.sleepTime = 3600
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

		config, err := d.GetNodeConfig()
		if err != nil {
			if (d.dPrintSleep(err, d.sleepTime)) {
				break BEGIN
			}
			continue BEGIN
		}

		if config["auto_update"] == "1" {
			_, url, err := utils.GetUpdVerAndUrl(config["auto_update_url"])
			if err != nil {
				if (d.dPrintSleep(err, d.sleepTime)) {
					break BEGIN
				}
				continue BEGIN
			}
			err = utils.DcoinUpd(url)
			if err != nil {
				if (d.dPrintSleep(err, d.sleepTime)) {
					break BEGIN
				}
				continue BEGIN
			}
		}

		if d.dSleep(d.sleepTime) {
			break BEGIN
		}
	}
}
