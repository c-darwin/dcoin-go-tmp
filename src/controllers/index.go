package controllers
import (
	//"database/sql"
	//_ "github.com/lib/pq"
	//"reflect"
	"fmt"
	"html/template"
	//"bufio"
	"bytes"
	"net/http"
	"strings"
	"bindatastatic"
)

type index struct {
	DbOk bool
	Lang map[string]string
	Key string
	SetLang string
}

func Index(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Index")

	lang := GetLang(w, r)

	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)

	r.ParseForm()
	var key string
	formKey := r.FormValue("key")
	if len(formKey) > 0 {
		key = formKey
		// пишем в сессию, что бы ctrl+F5 не сбрасывал ключ (для авто-входа с dcoin.me)
		sess.Set("private_key", key)
	} else {
		key = GetSessPrivateKey(w, r)
	}
	key = strings.Replace(key,"\r","\n",-1)
	key = strings.Replace(key,"\n\n","\n",-1)
	key = strings.Replace(key,"\n","\\\n",-1)

	setLang := r.FormValue("lang")

	data, err := bindatastatic.Asset("static/templates/index.html")
	t := template.New("template")
	t, err = t.Parse(string(data))
	//t, err := template.Parse("templates/index.html")
	if err != nil {
		fmt.Println(err)
	}
	b := new(bytes.Buffer)
	t.Execute(b, &index{DbOk: true, Lang: globalLangReadOnly[lang], Key: key, SetLang: setLang})
	w.Write(b.Bytes())
}
