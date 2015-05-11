package controllers
import (
	//"database/sql"
	//_ "github.com/lib/pq"
	//"reflect"
	"fmt"
	"html/template"
	//"bufio"
	"bytes"
	"static"
	//"utils"
	//"runtime"
	//"consts"
	//"schema"
)

type installStep1Struct struct {
	Lang map[string]string
}

// Шаг 1 - выбор либо стандартных настроек (sqlite и блокчейн с сервера) либо расширенных - pg/mysql и загрузка с нодов
func (c *Controller) Install_step_1() (string, error) {

	fmt.Println("Install_step_1")
	data, err := static.Asset("static/templates/install_step_1.html")
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println(data)
	t := template.New("template")
	t, _ = t.Parse(string(data))

	b := new(bytes.Buffer)
	t.Execute(b, &installStep0Struct{Lang: c.Lang})
	return b.String(), nil
}
