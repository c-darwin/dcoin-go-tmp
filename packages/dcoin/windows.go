// +build windows

package dcoin

import (
	"os/exec"
	"fmt"
)

func KillPid(pid string) error {
	var rez []byte
	rez, err := exec.Command("taskkill","/pid", pid).Start()
	if err!=nil {
		return err
	}
	fmt.Printf("taskkill /pid %s: %s\n", pid, rez)
	return nil
}
