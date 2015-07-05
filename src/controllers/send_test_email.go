package controllers
import (
	"utils"
	"fmt"
	"errors"
	//"log"
)
func (c *Controller) SendTestEmail() (string, error) {

	if c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	c.r.ParseForm()

	mailData, err := c.OneRow("SELECT * FROM "+c.MyPrefix+"my_table").String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	To:=mailData["email"]
	if len(mailData["use_smtp"]) > 0 && len(mailData["smtp_server"]) > 0 {
		err = c.SendMail("Test", "Test", To, mailData)
		if err != nil {
			return fmt.Sprintf(`{"error":"%s"}`, err), nil
		}
	} else if c.Community {
		// в пуле пробуем послать с смтп-ешника админа пула
		prefix := utils.Int64ToStr(c.PoolAdminUserId)+"_"
		mailData, err := c.OneRow("SELECT * FROM "+prefix+"my_table").String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		err = c.SendMail("Test", "Test",To, mailData)
		if err != nil {
			return fmt.Sprintf(`{"error":"%s"}`, err), nil
		}
	} else {
		return `{"error":"Incorrect mail data"}`, nil
	}

	return `{"error":"null"}`, nil
}

