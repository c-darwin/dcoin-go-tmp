package controllers
import (
	"fmt"
	"html/template"
	"bytes"
	"github.com/c-darwin/dcoin-tmp/packages/static"
)

type installStep0Struct struct {
	Lang map[string]string
}

// Шаг 0 - выбор языка
func (c *Controller) InstallStep0() (string, error) {

	data, err := static.Asset("static/templates/install_step_0.html")
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
