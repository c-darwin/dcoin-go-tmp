package daemons

import (
	"github.com/astaxie/beego/config"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"github.com/op/go-logging"
	"flag"
)

var (
	log = logging.MustGetLogger("daemons")
	DaemonCh chan bool
	AnswerDaemonCh chan bool
	configIni map[string]string
)

func init() {

	flag.Parse()

  	// мониторим config.ini на наличие изменений
	go func() {
		for {
			/*utils.Dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				log.Info("%v", utils.ErrInfo(err))
			}*/
			/*pwd, err := os.Getwd()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}*/

			configIni_, err := config.NewConfig("ini", *utils.Dir+"/config.ini")
			if err != nil {
				log.Info("%v", utils.ErrInfo(err))
			}
			configIni, err = configIni_.GetSection("default")
			if err != nil {
				log.Info("%v", utils.ErrInfo(err))
			}
			if len(configIni["db_type"]) > 0 {
				break
			}
			utils.Sleep(3)
		}
	}()
}

func CheckDaemonsRestart() bool {
	select {
	case <-DaemonCh:
		AnswerDaemonCh<-true
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