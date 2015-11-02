package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/c-darwin/dcoin-go-tmp/packages/static"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

type index struct {
	DbOk        bool
	Lang        map[string]string
	Key         string
	SetLang     string
	IOS         bool
	Upgrade3    string
	Upgrade4    string
	Android     bool
	Mobile      bool
	ShowIOSMenu bool
}

func Index(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	parameters_ := make(map[string]interface{})
	err := json.Unmarshal([]byte(r.PostFormValue("parameters")), &parameters_)
	if err != nil {
		log.Error("%v", err)
	}
	log.Debug("parameters_=%", parameters_)
	parameters := make(map[string]string)
	for k, v := range parameters_ {
		parameters[k] = utils.InterfaceToStr(v)
	}

	lang := GetLang(w, r, parameters)

	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		log.Error("%v", err)
	}
	defer sess.SessionRelease(w)

	sessUserId := GetSessUserId(sess)

	var key, myPrefix, status string
	var communityUsers []int64

	if utils.DB != nil && utils.DB.DB != nil {
		communityUsers, err = utils.DB.GetCommunityUsers()
		if err != nil {
			log.Error("%v", err)
		}
		if len(communityUsers) > 0 {
			myPrefix = utils.Int64ToStr(sessUserId) + "_"
		}
		status, err := utils.DB.Single("SELECT status FROM " + myPrefix + "my_table").String()

		// чтобы нельзя было зайти по локалке
		// :: - для маков
		if ok, _ := regexp.MatchString(`(\:\:)|(127\.0\.0\.1)`, r.RemoteAddr); ok {
			if status != "waiting_accept_new_key" && status != "waiting_set_new_key" {
				key, err = utils.DB.Single("SELECT private_key FROM " + myPrefix + "my_keys WHERE block_id = (SELECT max(block_id) FROM " + myPrefix + "my_keys)").String()
				if err != nil {
					log.Error("%v", err)
				}
			}
		}
	}

	showIOSMenu := true
	// Когда меню не выдаем
	if utils.DB == nil || utils.DB.DB == nil {
		showIOSMenu = false
	} else {
		if status == "my_pending" {
			showIOSMenu = false
		}
	}

	if sessUserId == 0 {
		showIOSMenu = false
	}

	if showIOSMenu && utils.DB != nil && utils.DB.DB != nil {
		blockData, err := utils.DB.GetInfoBlock()
		if err != nil {
			log.Error("%v", err)
		}
		wTime := int64(12)
		wTimeReady := int64(2)
		log.Debug("wTime: %v / utils.Time(): %v / blockData[time]: %v", wTime, utils.Time(), utils.StrToInt64(blockData["time"]))
		// если время менее 12 часов от текущего, то выдаем не подвержденные, а просто те, что есть в блокчейне
		if utils.Time()-utils.StrToInt64(blockData["time"]) < 3600*wTime {
			lastBlockData, err := utils.DB.GetLastBlockData()
			if err != nil {
				log.Error("%v", err)
			}
			log.Debug("lastBlockData[lastBlockTime]: %v", lastBlockData["lastBlockTime"])
			log.Debug("time.Now().Unix(): %v", utils.Time())
			if utils.Time()-lastBlockData["lastBlockTime"] >= 3600*wTimeReady {
				showIOSMenu = false
			}
		} else {
			showIOSMenu = false
		}
	}
	if showIOSMenu && !utils.Mobile() {
		showIOSMenu = false
	}

	mobile := utils.Mobile()
	if ok, _ := regexp.MatchString("(?i)(iPod|iPhone|iPad|Android)", r.UserAgent()); ok {
		mobile = true
	}

	var upgrade3 string
	if len(r.FormValue("upgrade3")) > 0 {
		upgrade3 = "1"
	}
	var upgrade4 string
	if len(r.FormValue("upgrade4")) > 0 {
		upgrade4 = "1"
	}
	formKey := r.FormValue("key")
	if len(formKey) > 0 {
		key = formKey
		// пишем в сессию, что бы ctrl+F5 не сбрасывал ключ (для авто-входа с dcoin.club)
		sess.Set("private_key", key)
	} else if len(key) == 0 {
		key = GetSessPrivateKey(w, r)
	}
	key = strings.Replace(key, "\r", "\n", -1)
	key = strings.Replace(key, "\n\n", "\n", -1)
	key = strings.Replace(key, "\n", "\\\n", -1)

	setLang := r.FormValue("lang")

	data, err := static.Asset("static/templates/index.html")
	t := template.New("template")
	t, err = t.Parse(string(data))
	if err != nil {
		log.Error("%v", err)
	}
	b := new(bytes.Buffer)
	err = t.Execute(b, &index{
		Upgrade3:    upgrade3,
		Upgrade4:    upgrade4,
		DbOk:        true,
		Lang:        globalLangReadOnly[lang],
		Key:         key,
		SetLang:     setLang,
		ShowIOSMenu: showIOSMenu,
		/*IOS: true,
		Android: false,
		Mobile: true})*/
		IOS:     utils.IOS(),
		Android: utils.Android(),
		Mobile:  mobile})
	if err != nil {
		log.Error("%v", err)
	}
	w.Write(b.Bytes())
}
