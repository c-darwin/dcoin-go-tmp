package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"net/http"
	"regexp"
)

func AjaxE(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("ajax Recovered", r)
			panic(r)
		}
	}()
	log.Debug("AjaxE")
	w.Header().Set("Content-type", "text/html")

	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		log.Error("%v", err)
		return
	}
	defer sess.SessionRelease(w)
	sessUserId := GetSessEUserId(sess)
	log.Debug("sessUserId", sessUserId)

	c := new(Controller)
	c.r = r
	c.w = w
	c.sess = sess
	c.SessUserId = sessUserId

	if utils.DB == nil || utils.DB.DB == nil {
		log.Error("utils.DB == nil")
		w.Write([]byte("DB == nil"))
		return
	}
	c.DCDB = utils.DB

	c.Parameters, err = c.GetParameters()
	log.Debug("parameters=", c.Parameters)

	lang := GetLang(w, r, c.Parameters)
	log.Debug("lang", lang)
	c.Lang = globalLangReadOnly[lang]
	c.LangInt = int64(lang)
	if lang == 42 {
		c.TimeFormat = "2006-01-02 15:04:05"
	} else {
		c.TimeFormat = "2006-02-01 15:04:05"
	}

	r.ParseForm()
	controllerName := r.FormValue("controllerName")
	log.Debug("controllerName=", controllerName)

	html := ""

	c.EConfig, err = c.GetMap(`SELECT * FROM e_config`, "name", "value")
	if err != nil {
		log.Error("%v", err)
	}
	c.ECommission = utils.StrToFloat64(c.EConfig["commission"])

	if ok, _ := regexp.MatchString(`^(?i)ESaveOrder|ESignUp|ELogin|ELogout|ESignLogin|ECheckSign|ERedirect$`, controllerName); !ok {
		html = "Access denied 0"
	} else {
		if ok, _ := regexp.MatchString(`^(?i)ESaveOrder|ESignUp|ELogin|ESignLogin|ECheckSign|ERedirect$`, controllerName); !ok && c.SessUserId <= 0 {
			html = "Access denied 1"
		} else {
			// вызываем контроллер в зависимости от шаблона
			log.Debug("controllerName %s", controllerName)
			html, err = CallController(c, controllerName)
			log.Debug("html %s", html)
			if err != nil {
				log.Error("ajax error: %v", err)
			}
		}
	}
	w.Write([]byte(html))

}

