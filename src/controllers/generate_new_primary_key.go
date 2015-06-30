package controllers
import (
	"utils"
	"log"
	"errors"
	"encoding/json"
)

func (c *Controller) GenerateNewPrimaryKey() (string, error) {

	if c.SessRestricted!=0 {
		return "", errors.New("Permission denied")
	}

	c.r.ParseForm()
	password := c.r.FormValue("password")

	priv, pub := utils.GenKeys()
	if len(password) > 0 {
		encKey, err :=utils.Encrypt(utils.Md5(password), []byte(priv))
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		priv = string(encKey)
	}
	json, err := json.Marshal(map[string]string{"private_key": priv, "public_key": pub, "password_hash": string(utils.DSha256(password))})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	log.Println(json)
	return string(json), nil
}

