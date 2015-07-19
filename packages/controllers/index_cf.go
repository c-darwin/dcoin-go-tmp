package controllers
import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"github.com/c-darwin/dcoin-go-tmp/packages/static"
	"html/template"
	"bytes"
	"encoding/json"
)

type indexCf struct {
	CfUrl string
	Lang string
	Nav template.JS
	CfLang map[string]string
}

func IndexCf(w http.ResponseWriter, r *http.Request) {

	nav := ""
	if len(r.URL.RawQuery) > 0 {
		re, _ := regexp.Compile(`category\-([0-9]+)`)
		match := re.FindStringSubmatch(r.URL.RawQuery)
		if len(match) > 0 {
			nav = "fc_navigate ('cfCatalog', {'category_id':"+match[1]+"})\n"
		} else {
			re, _ := regexp.Compile(`([A-Z0-9]{7}|id-[0-9]+)\-?([0-9]+)?\-?(funders|comments|news|home|payment)?`)
			match0 := re.FindStringSubmatch(r.URL.RawQuery)
			if len(match0) > 1 {
				// $m[1] - название валюты или id валюты
				// $m[2] - id языка
				// $m[3] - тип страницы (funders|comments|news)
				addNav := ""
				re, _ := regexp.Compile(`id\-([0-9]+)`)
				match := re.FindStringSubmatch(match0[1])
				if len(match) > 1 {
					addNav+= "'onlyProjectId':'"+match[1]+"',";
				} else {
					addNav+= "'onlyCfCurrencyName':'"+match[1]+"',";
				}
				if len(match0) > 2 {
					addNav+="'lang_id':'"+match0[2]+"',";
				}
				if len(match0) > 3 {
					addNav+="'page':'"+match0[3]+"',";
				}
				addNav = addNav[:len(addNav)-1]
				nav = "fc_navigate ('cfPagePreview', {"+addNav+"})\n";
			}
		}
	} else {
		nav = "fc_navigate ('cfCatalog')\n";
	}

	log.Debug(nav)

	c := new(Controller)
	c.r = r
	dbInit := false;
	if len(configIni["db_user"]) > 0 || (configIni["db_type"]=="sqlite") {
		dbInit = true
	}
	if dbInit {
		var err error
		c.DCDB, err = utils.NewDbConnect(configIni)
		if err != nil {
			log.Error("%v", err)
			dbInit = false
		} else {
			defer utils.DbClose(c.DCDB)
		}
		// отсутвие таблы выдаст ошибку, значит процесс инсталяции еще не пройден и надо выдать 0-й шаг
		_, err = c.DCDB.Single("SELECT progress FROM install").String()
		if err != nil {
			log.Error("%v", err)
			dbInit = false
		}

		cfUrl, err := c.GetCfUrl()
		cfLang, err := c.GetAllCfLng()

		r.ParseForm()
		parameters_ := make(map[string]interface {})
		err = json.Unmarshal([]byte(c.r.PostFormValue("parameters")), &parameters_)
		if err != nil {
			log.Error("%v", err)
		}
		parameters := make(map[string]string)
		for k, v := range parameters_ {
			parameters[k] = utils.InterfaceToStr(v)
		}
		lang := GetLang(w, r, parameters)

		data, err := static.Asset("static/templates/index_cf.html")
		t := template.New("template")
		t, err = t.Parse(string(data))
		if err != nil {
			fmt.Println(err)
		}
		b := new(bytes.Buffer)
		t.Execute(b, &indexCf{CfUrl: cfUrl, Lang: utils.IntToStr(lang), Nav: template.JS(nav), CfLang: cfLang })
		w.Write(b.Bytes())
	}
}
