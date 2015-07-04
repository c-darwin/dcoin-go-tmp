package controllers
import (
	"errors"
	"utils"
)

func (c *Controller) SaveRaceCountry() (string, error) {

	if c.SessUserId == 0 || c.SessRestricted != 0 {
		return "", errors.New("Permission denied")
	}

	c.r.ParseForm()
	race := int(utils.StrToFloat64(c.r.FormValue("race")))
	country := int(utils.StrToFloat64(c.r.FormValue("country")))
	err := c.ExecSql("UPDATE "+c.MyPrefix+"my_table SET race = ?, country = ?", race, country)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return `{"error":"0"}`, nil
}
