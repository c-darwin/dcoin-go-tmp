package controllers

import (
	"errors"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

func (c *Controller) ClearDbLite() (string, error) {

	if !c.NodeAdmin || c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	// остановим демонов, иначе будет паника, когда таблы обнулятся

	utils.Mutex.Lock()
	err := c.ExecSql(`DELETE FROM main_lock`)
	if err != nil {
		utils.Mutex.Unlock()
		return "", utils.ErrInfo(err)
	}
	err = c.ExecSql(`INSERT INTO main_lock (lock_time, script_name) VALUES (1, 'nulling')`)
	if err != nil {
		utils.Mutex.Unlock()
		return "", utils.ErrInfo(err)
	}
	utils.Mutex.Unlock()

	return "", nil
}
