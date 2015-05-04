package controllers
import (
	//"database/sql"
	//_ "github.com/lib/pq"
	//"reflect"
	"fmt"
	"html/template"
	//"bufio"
	"bytes"
	"bindatastatic"
)

type page struct {
	Title string
	Msg string
}

func (p *Controller) Home() (string, error) {
	fmt.Println("Home1")
	data, err := bindatastatic.Asset("static/templates/home.html")
	t := template.New("template")
	t, err = t.Parse(string(data))
	if err!=nil{
		return "", err
	}
	b := new(bytes.Buffer)
	t.Execute(b, &page{Title: p.Lang["geolocation"]})
	return b.String(), nil
}
