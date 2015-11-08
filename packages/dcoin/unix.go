// +build linux darwin freebsd
// +build 386 amd64

package dcoin

import (
	"syscall"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

func KillPid(pid int) error {
	err := syscall.Kill(pid, syscall.SIGTERM)
	if err != nil {
		return utils.ErrInfo(err)
	}
	return nil
}