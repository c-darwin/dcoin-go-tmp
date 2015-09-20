package controllers
import (
	"errors"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

type availableKeysPage struct {
	AutoLogin bool
	Key string
	LangId int
}

func checkAvailableKey(key string, db *utils.DCDB) (int64, string, error) {
	publicKeyAsn, err := utils.GetPublicFromPrivate(key)
	if err != nil {
		log.Debug("%v", err)
		return 0, "", utils.ErrInfo(err)
	}
	log.Debug("publicKeyAsn: %s", publicKeyAsn)
	userId, err := db.Single("SELECT user_id FROM users WHERE hex(public_key_0) = ?", publicKeyAsn).Int64()
	if err != nil {
		return 0, "", utils.ErrInfo(err)
	}
	log.Debug("userId: %s", userId)
	if userId == 0 {
		return 0, "", errors.New("null userId")
	}
	allTables, err := db.GetAllTables()
	if err != nil {
		return 0, "", utils.ErrInfo(err)
	}
	// может другой юзер уже начал смену ключа. актуально для пула
	if utils.InSliceString(utils.Int64ToStr(userId)+"_my_table", allTables) {
		return 0, "", errors.New("exists _my_table")
	}
	return userId, string(publicKeyAsn), nil
}

func (c *Controller) AvailableKeys() (string, error) {

	if c.Community {
		// если это пул, то будет прислан email
		email := c.r.FormValue("email")
		if !utils.ValidateEmail(email) {
			return utils.JsonAnswer("Incorrect email", "error").String(), nil
		}
		// если мест в пуле нет, то просто запишем юзера в очередь
		pool_max_users, err := c.Single("SELECT pool_max_users FROM config").Int()
		if err != nil {
			return "", utils.JsonAnswer(utils.ErrInfo(err), "error").Error()
		}
		if len(c.CommunityUsers) >= pool_max_users {
			err = c.ExecSql("INSERT INTO pool_waiting_list ( email, time, user_id ) VALUES ( ?, ?, ? )", email, utils.Time(), 0)
			if err != nil {
				return "", utils.JsonAnswer(utils.ErrInfo(err), "error").Error()
			}
			return utils.JsonAnswer(c.Lang["pool_is_full"], "error").String(), nil
		}
	}
	userId, publicKey, err := c.GetAvailableKey()
	if err != nil {
		return "", utils.JsonAnswer(utils.ErrInfo(err), "error").Error()
	}
	if userId > 0 {
		c.sess.Set("user_id", userId)
		c.sess.Set("public_key", publicKey)
		log.Debug("user_id: %d", userId)
		log.Debug("public_key: %s", publicKey)
		return utils.JsonAnswer("success", "success").String(), nil
	}
	/*
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
	*/
	return utils.JsonAnswer("no_available_keys", "error").String(), nil
}
