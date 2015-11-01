package controllers

import (
	"errors"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

func (c *Controller) ClearDb() (string, error) {

	if !c.NodeAdmin || c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	err := c.ExecSql(`UPDATE install SET progress = 0`)
	if err != nil {
		utils.Mutex.Unlock()
		return "", utils.ErrInfo(err)
	}

	return "", nil
}
