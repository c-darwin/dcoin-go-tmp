package daemons

import (
	"github.com/astaxie/beego/config"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"github.com/op/go-logging"
)
var log = logging.MustGetLogger("daemons")
var DaemonCh chan bool
var AnswerDaemonCh chan bool
var configIni map[string]string

func init() {
  	// мониторим config.ini на наличие изменений
	go func() {
		for {
			configIni_, err := config.NewConfig("ini", "config.ini")
			if err != nil {
				log.Info("%v", utils.ErrInfo(err))
			}
			configIni, err = configIni_.GetSection("default")
			if err != nil {
				log.Info("%v", utils.ErrInfo(err))
			}
			utils.Sleep(1)
		}
	}()
}

func CallDeamons(name string) {
	//utils.CallMethod(p, name)
}

func DbConnect() *utils.DCDB {
	for {
		if len(configIni) == 0 {
			utils.Sleep(1)
			continue
		}
		db, err := utils.NewDbConnect(configIni)
		if err == nil {
			return db
		} else {
			utils.Sleep(1)
		}
	}
	return nil
}