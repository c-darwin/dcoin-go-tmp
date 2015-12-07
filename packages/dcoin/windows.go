// +build windows

package dcoin

import (
	//"os/exec"
	//"fmt"
	//"os"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"os/exec"
)

func KillPid(pid string) error {
	err := utils.DB.ExecSql(`INSERT INTO stop_daemons(stop_time) VALUES (?)`, utils.Time())
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		return err
	}
	//var rez []byte
	/*file, err := os.OpenFile("kill", os.O_APPEND|os.O_WRONLY|os.O_CREATE,0600)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString("1")
	*/
	err := exec.Command("taskkill","/pid", pid).Start()
	if err!=nil {
		return err
	}
	//fmt.Printf("taskkill /pid %s: %s\n", pid, rez)
	return nil
}
