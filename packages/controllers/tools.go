package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"net/http"
	"regexp"
)

func Tools(w http.ResponseWriter, r *http.Request) {

	var err error
	log.Debug("Tools")
	w.Header().Set("Content-type", "text/html")

	c := new(Controller)
	c.r = r
	c.w = w
	dbInit := false
	if len(configIni["db_user"]) > 0 || configIni["db_type"] == "sqlite" {
		dbInit = true
	}

	if dbInit {
		c.DCDB = utils.DB
		if c.DCDB == nil {
			log.Error("utils.DB == nil")
			dbInit = false
		}
		c.Variables, err = c.GetAllVariables()

	}

	r.ParseForm()
	controllerName := r.FormValue("controllerName")
	log.Debug("controllerName=", controllerName)

	html := ""
	if ok, _ := regexp.MatchString(`^(?i)GetBlock|AvailableKeys$`, controllerName); !ok {
		html = "Access denied"
	} else {
		// вызываем контроллер в зависимости от шаблона
		html, err = CallController(c, controllerName)
		if err != nil {
			log.Error("%v", err)
		}
	}
	w.Write([]byte(html))
}
