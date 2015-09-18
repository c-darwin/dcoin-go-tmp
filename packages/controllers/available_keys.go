package controllers
import (
	"errors"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"regexp"
	"math/rand"
	//"strings"
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

	keysStr, err := utils.GetHttpTextAnswer("http://dcoin.club/keys")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	//keysStr = strings.Replace(keysStr, "\n", "", -1)
	r, _ := regexp.Compile("(?s)-----BEGIN RSA PRIVATE KEY-----(.*?)-----END RSA PRIVATE KEY-----")
	keys := r.FindAllString(keysStr, -1)
	for i := range keys {
		j := rand.Intn(i + 1)
		keys[i], keys[j] = keys[j], keys[i]
	}
	for _, key := range keys {
		userId, pubKey, err := checkAvailableKey(key, c.DCDB)
		if err != nil {
			log.Error("%v", utils.ErrInfo(err))
		}
		log.Debug("checkAvailableKey userId: %v", userId)
		if userId > 0 {
			// запишем приватный ключ в БД, чтобы можно было подписать тр-ию на смену ключа
			myPref := ""
			if c.Community {
				myPref = utils.Int64ToStr(userId)+"_"
				err = c.ExecSql("INSERT INTO "+myPref+"my_table (user_id, status) VALUES (?, ?)", userId, "waiting_set_new_key")
				if err != nil {
					return "", utils.ErrInfo(err)
				}
			} else {
				err = c.ExecSql("UPDATE my_table SET user_id = ?, status = ?", userId, "waiting_set_new_key")
				if err != nil {
					return "", utils.ErrInfo(err)
				}
			}

			// пишем приватный в my_keys т.к. им будем подписывать тр-ию на смену ключа
			err = c.ExecSql("INSERT INTO "+myPref+"my_keys (private_key, status) VALUES (?, ?)", key, "approved")
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			newPrivKey, newPubKey := utils.GenKeys()
			// сразу генерируем новый ключ и пишем приватный временно в my_keys, чтобы можно было выдавать юзеру для скачивания
			err = c.ExecSql("INSERT INTO "+myPref+"my_keys (private_key, public_key, status) VALUES (?, ?, ?)", newPrivKey, utils.HexToBin([]byte(newPubKey)), "my_pending")
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			c.sess.Set("user_id", userId)
			c.sess.Set("public_key", pubKey)
			log.Debug("user_id: %d", userId)
			log.Debug("public_key: %s", pubKey)

			/*c.w.Header().Set("Content-Type", "application/octet-stream")
			c.w.Header().Set("Content-Length", utils.IntToStr(len(key)))
			c.w.Header().Set("Content-Disposition", `attachment; filename="key.txt"`)
			c.w.Header().Set("Access-Control-Allow-Origin", "*")
			c.w.Write([]byte(key))*/
			return "ok", nil
		}
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

	return "", nil
}
