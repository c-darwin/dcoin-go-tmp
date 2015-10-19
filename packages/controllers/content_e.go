package controllers

import (
	"encoding/json"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"html/template"
	"net/http"
	"regexp"
)

type contentE struct {
	CfUrl  string
	Lang   string
	Nav    template.JS
	CfLang map[string]string
}

func ContentE(w http.ResponseWriter, r *http.Request) {

	if utils.DB != nil && utils.DB.DB != nil {

		sess, _ := globalSessions.SessionStart(w, r)
		defer sess.SessionRelease(w)

		c := new(Controller)
		c.r = r
		c.SessUserId = GetSessUserId(sess)
		c.DCDB = utils.DB

		r.ParseForm()
		tplName := r.FormValue("controllerName")
		parameters_ := make(map[string]interface{})
		err := json.Unmarshal([]byte(c.r.PostFormValue("parameters")), &parameters_)
		if err != nil {
			log.Error("%v", err)
		}
		parameters := make(map[string]string)
		for k, v := range parameters_ {
			parameters[k] = utils.InterfaceToStr(v)
		}
		c.Parameters = parameters
		lang := GetLang(w, r, parameters)
		c.Lang = globalLangReadOnly[lang]
		log.Debug("c.Lang:", c.Lang)
		c.LangInt = int64(lang)
		if lang == 42 {
			c.TimeFormat = "2006-01-02 15:04:05"
		} else {
			c.TimeFormat = "2006-02-01 15:04:05"
		}
		// если в параметрах пришел язык, то установим его
		newLang := utils.StrToInt(parameters["lang"])
		if newLang > 0 {
			log.Debug("newLang", newLang)
			SetLang(w, r, newLang)
		}

		c.EConfig, err = c.GetMap(`SELECT * FROM e_config`, "name", "value")
		if err != nil {
			log.Error("%v", err)
		}
		html := ""
		if ok, _ := regexp.MatchString(`^(?i)emain`, tplName); !ok {
			html = "Access denied"
		} else {
			// вызываем контроллер в зависимости от шаблона
			html, err = CallController(c, tplName)
			if err != nil {
				log.Error("%v", err)
			}
		}
		w.Write([]byte(html))
	}
}
