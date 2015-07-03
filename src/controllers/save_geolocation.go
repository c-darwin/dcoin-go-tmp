package controllers
import (
	"errors"
	"utils"
	"strings"
)

func (c *Controller) SaveUserCoords() (string, error) {

	if c.SessUserId == 0 || c.SessRestricted != 0 {
		return "", errors.New("Permission denied")
	}

	c.r.ParseForm()
	geolocation := c.r.FormValue("geolocation")
	if len(geolocation) > 0 {
		x := strings.Split(geolocation, ",")
		if len(x) == 2 {
			geolocationLat := utils.Round(utils.StrToFloat64(x[0]), 5)
			geolocationLon := utils.Round(utils.StrToFloat64(x[1]), 5)
			err := c.ExecSql("UPDATE "+c.MyPrefix+"my_table SET geolocation = ?", geolocationLat+", "+geolocationLon)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}
	}
	return `{"error":"0"}`, nil
}
