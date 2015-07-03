package controllers
import (
	"errors"
	"utils"
)

func (c *Controller) ClearVideo() (string, error) {

	if c.SessUserId == 0 || c.SessRestricted != 0 {
		return "", errors.New("Permission denied")
	}

	err := c.ExecSql("UPDATE "+c.MyPrefix+"my_table SET video_url_id = ?, video_type = ?", "", "")
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	return ``, nil
}
