package controllers
import (
	"errors"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"github.com/c-darwin/dcoin-go-tmp/packages/static"
	"html/template"
	"bytes"
	"strings"
)

type availableKeysPage struct {
	AutoLogin bool
	Key string
	LangId int
}

func checkAvailableKey(key string, db *utils.DCDB) (error) {
	publicKeyAsn, err := utils.GetPublicFromPrivate(key)
	if err != nil {
		return utils.ErrInfo(err)
	}
	userId, err := db.Single("SELECT user_id FROM users WHERE public_key_0  =  [hex]", publicKeyAsn).Int64()
	if err != nil {
		return utils.ErrInfo(err)
	}
	if userId == 0 {
		return errors.New("null userId")
	}
	return nil
}

func (c *Controller) AvailableKeys() (string, error) {

	c.r.ParseForm()
	var key string
	var autoLogin bool
	if len(c.r.FormValue("auto_login")) > 0 {
		autoLogin = true
	}
	langId := utils.StrToInt(c.r.FormValue("langId"))
	if langId == 0 {
		langId = 1
	}
	for i:=0;;i++ {
		maxId, err := c.Single("SELECT max(id) FROM _my_refs", utils.Time()-1800).Int()
		randId := utils.RandInt(1, maxId+1)
		key, err = c.Single("SELECT private_key FROM _my_refs WHERE used_time < ? AND id = ?", utils.Time()-1800, randId).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		key = strings.Replace(key, `\r\n`, "\n", -1)
		if checkAvailableKey(key, c.DCDB) == nil || i > 10 {
			if i > 10 {
				key = ""
				log.Debug("%v", err)
			}
			break
		}
	}
	if len(key) > 0 {
		c.ExecSql("UPDATE _my_refs SET used_time = ? WHERE private_key = ?", utils.Time(), key)
	}

	if len(c.r.FormValue("download")) > 0 {
		c.w.Header().Set("Content-Type", "application/octet-stream")
		c.w.Header().Set("Content-Length", utils.IntToStr(len(key)))
		c.w.Header().Set("Content-Disposition", `attachment; filename="key.txt"`)
		c.w.Header().Set("Access-Control-Allow-Origin", "*")
		c.w.Write([]byte(key))
	} else {
		data, err := static.Asset("static/templates/available_keys.html")
		t := template.Must(template.New("template").Parse(string(data)))
		b := new(bytes.Buffer)
		err = t.ExecuteTemplate(b, "availableKeys", &availableKeysPage{AutoLogin: autoLogin, Key: key, LangId: langId})
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		return b.String(), nil
	}
	return "", nil
}
