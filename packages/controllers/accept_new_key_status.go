package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"errors"
)

func (c *Controller) AcceptNewKeyStatus() (string, error) {

	if c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	result := ""
	status, err := c.DCDB.Single("SELECT status FROM "+c.MyPrefix+"my_table").String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	if status == "approved" {
		result = "ok"
	}

	return result, nil
}

