// +build windows

package dcoin

import (
	"os/exec"
)

func KillPid(pid string) error {
	exec.Command("taskkill","/pid", pid)
	return nil
}
