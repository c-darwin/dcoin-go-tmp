package daemons

import (
	"flag"
	"github.com/astaxie/beego/config"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"github.com/op/go-logging"
	"time"
)

var (
	log             = logging.MustGetLogger("daemons")
	DaemonCh        chan bool
	AnswerDaemonCh  chan bool
	MonitorDaemonCh chan []string = make(chan []string, 100)
	configIni       map[string]string
)

type daemon struct {
	*utils.DCDB
	goRoutineName string
}

func (d *daemon) dbLock() (error, bool) {
	return d.DbLock(DaemonCh, AnswerDaemonCh, d.goRoutineName)
}

func (d *daemon) dbUnlock() error {
	log.Debug("dbUnlock %v", utils.Caller(1))
	return d.DbUnlock(d.goRoutineName)
}

func (d *daemon) unlockPrintSleep(err error, sleep time.Duration) {
	if err != nil {
		log.Error("%v", err)
	}
	err = d.DbUnlock(d.goRoutineName)
	if err != nil {
		log.Error("%v", err)
	}
	utils.Sleep(sleep)
}

func (d *daemon) unlockPrintSleepInfo(err error, sleep time.Duration) {
	if err != nil {
		log.Error("%v", err)
	}
	err = d.DbUnlock(d.goRoutineName)
	if err != nil {
		log.Error("%v", err)
	}
	utils.Sleep(sleep)
}

func ConfigInit() {
	// мониторим config.ini на наличие изменений
	go func() {
		for {
			configIni_, err := config.NewConfig("ini", *utils.Dir+"/config.ini")
			if err != nil {
				log.Error("%v", utils.ErrInfo(err))
			}
			configIni, err = configIni_.GetSection("default")
			if err != nil {
				log.Error("%v", utils.ErrInfo(err))
			}
			if len(configIni["db_type"]) > 0 {
				break
			}
			utils.Sleep(3)
		}
	}()
}

func init() {
	flag.Parse()
}

func CheckDaemonsRestart() bool {
	select {
	case <-DaemonCh:
		AnswerDaemonCh <- true
		return true
	default:
	}
	return false
}

func DbConnect() *utils.DCDB {
	for {
		if CheckDaemonsRestart() {
			return nil
		}
		if utils.DB == nil {
			utils.Sleep(1)
		} else {
			return utils.DB
		}
	}
	return nil
}
