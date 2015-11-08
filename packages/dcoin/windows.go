// +build windows

package dcoin

import (
	"os/exec"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

func KillPid(pid int) error {
	exec.Command("taskkill","/pid", pid)
	return nil
}
