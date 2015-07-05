package controllers
import (
	"utils"
	"errors"
	"encoding/json"
	"log"
)
func (c *Controller) SaveNotifications() (string, error) {

	if c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	c.r.ParseForm()

	var data []map[string]interface {}
	err := json.Unmarshal([]byte(c.r.PostFormValue("data")), &data)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	log.Println("data:",data)

	for k, _ :=  range data {
		err := c.ExecSql(`
				UPDATE `+c.MyPrefix+`my_notifications
				SET  email = ?,
						sms =  ?
				WHERE name = ?
				`, data[k]["email"].(float64), data[k]["sms"].(float64), data[k]["name"].(string))
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	}

	return `{"error":0}`, nil

}

