package controllers
import (
	//"database/sql"
	//_ "github.com/lib/pq"
	//"reflect"
	"fmt"
	//"utils"
	//"runtime"
	//"consts"
	//"schema"
	"utils"
	"log"
)

type signLoginStruct struct {

}

/*
 * Генерим код, который юзер должен подписать своим ключем, доказав тем самым, что именно он хочет войти в аккаунт
 * */

func (c *Controller) Sign_login() (string, error) {

	log.Println("sign_login")
	var hash []byte
	loginCode := utils.RandSeq(20)
	fmt.Println(c.r.RemoteAddr)
	if c.ConfigIni["sign_hash"] == "ip" {
		hash = utils.Md5(c.r.RemoteAddr);
	} else {
		hash = utils.Md5(c.r.Header.Get("User-Agent")+c.r.Header.Get("REMOTE_ADDR"));
	}
	log.Println("hash", hash)
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
