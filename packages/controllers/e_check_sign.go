package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"regexp"
)

func (c *Controller) ECheckSign() (string, error) {

	c.r.ParseForm()
	userId := utils.StrToInt64(c.r.FormValue("user_id"))
	sign := []byte(c.r.FormValue("sign"))
	if !utils.CheckInputData(string(sign), "hex_sign") {
		return `{"result":"incorrect sign"}`, nil
	}
	if !utils.CheckInputData(userId, "int") {
		return `{"result":"incorrect user_id"}`, nil
	}

	RemoteAddr := utils.RemoteAddrFix(c.r.RemoteAddr)
	re := regexp.MustCompile(`(.*?):[0-9]+$`)
	match := re.FindStringSubmatch(RemoteAddr)
	if len(match) != 0 {
		RemoteAddr = match[1]
	}
	log.Debug("RemoteAddr %s", RemoteAddr)
	hash := utils.Md5(c.r.Header.Get("User-Agent") + RemoteAddr)
	log.Debug("hash %s", hash)
	forSign, err := c.Single(`SELECT data FROM e_authorization WHERE hex(hash) = ?`, hash).String()
	if err != nil {
		return "{\"result\":0}", err
	}

	publicKey, err := c.GetMyPublicKey("")
	if err != nil {
		return "{\"result\":0}", err
	}

	// проверим подпись
	resultCheckSign, err := utils.CheckSign([][]byte{publicKey}, forSign, utils.HexToBin(sign), true)
	if resultCheckSign {
		// если это первый запрос, то создаем запись в табле users
		eUserId, err := c.Single(`SELECT id	FROM e_users WHERE user_id = ?`, userId).Int64()
		if err != nil {
			return "{\"result\":0}", err
		}
		if eUserId == 0 {
			eUserId, err = c.ExecSqlGetLastInsertId(`INSERT INTO e_users (user_id) VALUES (?)`, "user_id", userId)
		}
		token := utils.RandSeq(30)
		c.ExecSql(`INSERT INTO e_tokens (token, user_id) VALUES (?, ?)`, token, eUserId)
		return "{\"result\":1, \"token\":"+token+"}", nil

	} else {
		return "{\"result\":0}", nil
	}
}