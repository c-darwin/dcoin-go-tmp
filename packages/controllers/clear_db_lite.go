package controllers
import (
	"errors"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

func (c *Controller) ClearDbLite() (string, error) {

	if !c.NodeAdmin || c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	err := c.ExecSql(`DELETE FROM main_lock`)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	err = c.ExecSql(`INSERT INTO main_lock (lock_time, script_name) VALUES (1, 'nulling')`)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	return "", nil
}
