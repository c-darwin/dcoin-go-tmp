package controllers
import (
	"errors"
	"utils"
)

func (c *Controller) SaveHost() (string, error) {

	if c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	c.r.ParseForm()
	host := c.r.FormValue("host")
	if host[len(host)-1:] != "/" {
		host += "/"
	}

	if !utils.CheckInputData(host, "host")  {
		return `{"error":"1"}`, nil
	} else {
		err := c.ExecSql("UPDATE "+c.MyPrefix+"my_table SET host = ?", host)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		return `{"error":"0"}`, nil
	}

	return `{"error":"0"}`, nil
}
