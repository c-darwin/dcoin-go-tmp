package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

/*
 * Генерим код, который юзер должен подписать своим ключем, доказав тем самым, что именно он хочет войти в аккаунт
 * */

func (c *Controller) SignLogin() (string, error) {

	c.w.Header().Set("Access-Control-Allow-Origin", "*")

	var hash []byte
	loginCode := utils.RandSeq(20)
	log.Debug("configIni[sign_hash] %s", configIni["sign_hash"])
	log.Debug("c.r.RemoteAddr %s", c.r.RemoteAddr)
	log.Debug("c.r.Header.Get(User-Agent) %s", c.r.Header.Get("User-Agent"))
	if configIni["sign_hash"] == "ip" {
		hash = utils.Md5(c.r.RemoteAddr);
	} else {
		hash = utils.Md5(c.r.Header.Get("User-Agent")+c.r.RemoteAddr);
	}
	log.Debug("hash %s", hash)
	err := c.DCDB.ExecSql(`DELETE FROM authorization WHERE hash = [hex]`, hash)
	if err != nil {
		return "", err
	}
	err = c.DCDB.ExecSql(`INSERT INTO authorization (hash, data) VALUES ([hex], ?)`, hash, loginCode)
	if err != nil {
		return "", err
	}
	log.Debug("loginCode %v", loginCode)
	return "\""+loginCode+"\"", nil
}
