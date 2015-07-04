package controllers
import (
	//"database/sql"
	//_ "github.com/lib/pq"
	//"reflect"
	"fmt"
	"html/template"
	//"bufio"
	"bytes"
	//"net/http"
	//"strings"
	"utils"
	//"text/tabwriter"
	"static"
)

type loginStruct struct {
	Lang map[string]string
	MyModalIdName string
	UserID int64
	PoolTechWorks int
	SetupPassword bool
}

func (c *Controller) Login() (string, error) {
	var pool_tech_works int
	fmt.Println("login")


	funcMap := template.FuncMap{
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
/*
	t, err := template.New("Dcoin").Funcs(funcMap).ParseFiles("templates/login.html", "templates/modal.html")
	if err!=nil{
		fmt.Println(err)
	}*/
	data, err := static.Asset("static/templates/login.html")
	if err != nil {
		return "", err
	}
	modal, err := static.Asset("static/templates/modal.html")
	if err != nil {
		return "", err
	}

	t := template.Must(template.New("template").Funcs(funcMap).Parse(string(data)))
	t = template.Must(t.Parse(string(modal)))

	b := new(bytes.Buffer)

	// есть ли установочный пароль и был ли начально записан ключ
	var setupPassword bool
	if !c.Community {
		setupPassword_, err := c.Single("SELECT setup_password FROM config").String()
		if err != nil {
			return "", err
		}
		myKey, err := c.GetMyPublicKey(c.MyPrefix)
		if err != nil {
			return "", err
		}
		if len(myKey) == 0 && len(setupPassword_) > 0 {
			setupPassword = true
		}
	}
	//fmt.Println(c.Lang)
	// проверим, не идут ли тех. работы на пуле
	config, err := c.DCDB.OneRow("SELECT pool_admin_user_id, pool_tech_works FROM config").String()
	if err != nil {
		return "", err
	}
	if len(config["pool_admin_user_id"]) > 1 && config["pool_admin_user_id"] != utils.Int64ToStr(c.UserId) && config["pool_tech_works"] == "1" {
	 	pool_tech_works = 1
	} else {
		pool_tech_works = 0
	}
	fmt.Println(c.Lang["login_help_text"])
	t.ExecuteTemplate(b, "login", &loginStruct{Lang:  c.Lang, MyModalIdName: "myModalLogin", UserID: c.UserId, PoolTechWorks: pool_tech_works, SetupPassword: setupPassword})
	return b.String(), nil
}
