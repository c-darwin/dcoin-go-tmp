package controllers
import (
	"html/template"
	"bytes"
	"net/http"
	"strings"
	"github.com/c-darwin/dcoin-go-tmp/packages/static"
	"encoding/json"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

type index struct {
	DbOk bool
	Lang map[string]string
	Key string
	SetLang string
}

func Index(w http.ResponseWriter, r *http.Request) {

	parameters_ := make(map[string]interface {})
	err := json.Unmarshal([]byte(r.PostFormValue("parameters")), &parameters_)
	if err != nil {
		log.Error("%v", err)
	}
	log.Debug("parameters_=%",parameters_)
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

	var key string
	if utils.DB != nil {
		communityUsers, err := utils.DB.GetCommunityUsers()
		if err != nil {
			log.Error("%v", err)
		}
		myPrefix := ""
		if len(communityUsers) > 0 {
			myPrefix = utils.Int64ToStr(sessUserId)+"_";
		}
		key, err = utils.DB.Single("SELECT private_key FROM "+myPrefix+"my_keys WHERE block_id = (SELECT max(block_id) FROM "+myPrefix+"my_keys)").String()
		if err != nil {
			log.Error("%v", err)
		}
	}

	r.ParseForm()
	formKey := r.FormValue("key")
	if len(formKey) > 0 {
		key = formKey
		// пишем в сессию, что бы ctrl+F5 не сбрасывал ключ (для авто-входа с dcoin.club)
		sess.Set("private_key", key)
	} else if len(key) < 0 {
		key = GetSessPrivateKey(w, r)
	}
	key = strings.Replace(key,"\r","\n",-1)
	key = strings.Replace(key,"\n\n","\n",-1)
	key = strings.Replace(key,"\n","\\\n",-1)

	setLang := r.FormValue("lang")

	data, err := static.Asset("static/templates/index.html")
	t := template.New("template")
	t, err = t.Parse(string(data))
	if err != nil {
		log.Error("%v", err)
	}
	b := new(bytes.Buffer)
	err = t.Execute(b, &index{DbOk: true, Lang: globalLangReadOnly[lang], Key: key, SetLang: setLang})
	if err != nil {
		log.Error("%v", err)
	}
	w.Write(b.Bytes())
}
