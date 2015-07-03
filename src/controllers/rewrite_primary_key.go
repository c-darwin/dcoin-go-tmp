package controllers
import (
	"utils"
	"log"
	"errors"
)

type rewritePrimaryKeyPage struct {
	Alert string
	Lang map[string]string
}

func (c *Controller) RewritePrimaryKey() (string, error) {

	log.Println("PoolAdminControl")

	if !c.PoolAdmin {
		return "", utils.ErrInfo(errors.New("access denied"))
	}

	if len(c.r.FormValue("n")) > 0 {

		c.r.ParseForm()
		n := []byte(c.r.FormValue("n"))
		e := []byte(c.r.FormValue("e"))
		if !utils.CheckInputData(n, "hex") {
			return "", utils.ErrInfo(errors.New("incorrect n"))
		}
		if !utils.CheckInputData(e, "hex") {
			return "", utils.ErrInfo(errors.New("incorrect e"))
		}
		publicKey := utils.MakeAsn1(n, e)

		// проверим, есть ли вообще такой публичный ключ
		userId, err := c.Single("SELECT user_id FROM users WHERE public_key_0 = [hex]", publicKey).Int64()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if userId == 0 {
			return "", utils.ErrInfo(errors.New("incorrect public_key"))
		}

		// может быть юзер уже майнер?
		minerId, err := c.GetMyMinerId(userId)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		status := ""
		if minerId > 0 {
			status = "miner"
		} else {
			status = "user"
		}

		err = c.ExecSql(`DELETE FROM my_keys`)
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		err = c.ExecSql(`INSERT INTO `+c.MyPrefix+`my_keys (public_key, status) VALUES ([hex], ?)`, publicKey, "approved")
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		err = c.ExecSql(`UPDATE `+c.MyPrefix+`my_table SET user_id = ?, miner_id = ?, status = ?`, userId, minerId, status)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	}

	TemplateStr, err := makeTemplate("rewrite_primary_key", "rewritePrimaryKey", &rewritePrimaryKeyPage {
		Alert: c.Alert,
		Lang: c.Lang})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}