// +build windows

package dcoin

import (
	//"os/exec"
	//"fmt"
	//"os"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"os/exec"
	"fmt"
)

func KillPid(pid string) error {
	if utils.DB != nil && utils.DB.DB != nil {
		err := utils.DB.ExecSql(`INSERT INTO stop_daemons(stop_time) VALUES (?)`, utils.Time())
		if err != nil {
			log.Error("%v", utils.ErrInfo(err))
			return err
		}
	}
	//var rez []byte
	/*file, err := os.OpenFile("kill", os.O_APPEND|os.O_WRONLY|os.O_CREATE,0600)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString("1")
	*/
	/*err := exec.Command("taskkill","/pid", pid).Start()
	if err!=nil {
		return err
	}*/
	rez, err := exec.Command("tasklist","/fi", "PID eq "+pid).Output()
	if err!=nil {
		return err
	}
	if string(rez) == "" {
		return fmt.Errorf("null")
	} else {
		log.Debug("%rez s", string(rez))
		fmt.Println("rez", string(rez))
	}
	//fmt.Printf("taskkill /pid %s: %s\n", pid, rez)
	return nil
}
