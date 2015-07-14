package daemons

import (
	"github.com/astaxie/beego/config"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"log"

)
var configIni map[string]string
func init() {
  	// мониторим config.ini на наличие изменений
	go func() {
		for {
			configIni_, err := config.NewConfig("ini", "config.ini")
			if err != nil {
				log.Println(utils.ErrInfo(err))
			}
			configIni, err = configIni_.GetSection("default")
			if err != nil {
				log.Println(utils.ErrInfo(err))
			}
			utils.Sleep(1)
		}
	}()
}