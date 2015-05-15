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
)

type signLoginStruct struct {
}

/*
 * Генерим код, который юзер должен подписать своим ключем, доказав тем самым, что именно он хочет войти в аккаунт
 * */

func (c *Controller) Sign_login() (string, error) {

	fmt.Println("sign_login")
	var hash []byte
	loginCode := utils.RandSeq(20)
	fmt.Println(c.r.RemoteAddr)
	if configIni["sign_hash"] == "ip" {
		hash = utils.Md5(c.r.RemoteAddr);
	} else {
		hash = utils.Md5(c.r.Header.Get("User-Agent")+c.r.Header.Get("REMOTE_ADDR"));
	}
	fmt.Println("hash", hash)
	var sql string
	switch configIni["db_type"] {
	case "sqlite":
		sql = `DELETE FROM "authorization" WHERE "hash" = $1`
		_, err := c.DCDB.ExecSql(sql, hash)
		if err != nil {
			return "", err
		}
		sql =`INSERT INTO "authorization" ("hash", data) VALUES ($1, $2)`
		_, err = c.DCDB.ExecSql(sql, hash, loginCode)
		if err != nil {
			return "", err
		}
	case "postgresql":

		sql = `DELETE FROM "authorization" WHERE "hash" = $1`
		_, err := c.DCDB.ExecSql(sql, hash)
		if err != nil {
			return "", err
		}
		sql =`INSERT INTO  "authorization" (hash, data) VALUES (decode($1,'HEX'), $2)`
		_, err = c.DCDB.ExecSql(sql, hash, loginCode)
		if err != nil {
			return "", err
		}
	case "mysql":

		sql = `DELETE FROM "authorization" WHERE "hash" = $1`
		_, err := c.DCDB.ExecSql(sql, hash)
		if err != nil {
			return "", err
		}
		sql =`INSERT INTO authorization (hash, data) VALUES (0x$1, $2)`
		_, err = c.DCDB.ExecSql(sql, hash, loginCode)
		if err != nil {
			return "", err
		}
	}
	return "\""+loginCode+"\"", nil
}
