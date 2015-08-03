package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

/*
 * Генерим код, который юзер должен подписать своим ключем, доказав тем самым, что именно он хочет войти в аккаунт
 * */

func (c *Controller) SignLogin() (string, error) {

	var hash []byte
	loginCode := utils.RandSeq(20)
	if c.ConfigIni["sign_hash"] == "ip" {
		hash = utils.Md5(c.r.RemoteAddr);
	} else {
		hash = utils.Md5(c.r.Header.Get("User-Agent")+c.r.Header.Get("REMOTE_ADDR"));
	}
	log.Debug("hash", hash)
	err := c.DCDB.ExecSql(`DELETE FROM authorization WHERE hash = [hex]`, hash)
	if err != nil {
		return "", err
	}
	err = c.DCDB.ExecSql(`INSERT INTO authorization (hash, data) VALUES ([hex], ?)`, hash, loginCode)
	if err != nil {
		return "", err
	}
	return "\""+loginCode+"\"", nil
}
