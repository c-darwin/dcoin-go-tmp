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

type updatingBlockchainStruct struct {
	Lang map[string]string
	WaitText string
}

func (c *Controller) Updating_blockchain() (string, error) {
	fmt.Println("updating_blockchain")
	data, err := bindatastatic.Asset("static/templates/updating_blockchain.html")
	t := template.New("template")
	t, err = t.Parse(string(data))
	if err != nil {
		return "", err
	}
	var waitText string
	firstLoadBlockchain, err := c.DCDB.Single("SELECT first_load_blockchain FROM config")
	if err != nil {
		return "", err
	}
	if firstLoadBlockchain=="file" {
		waitText = c.Lang["loading_blockchain_please_wait"]
	} else {
		waitText = c.Lang["is_synchronized_with_the_dc_network"]
	}
	b := new(bytes.Buffer)
	t.Execute(b, &updatingBlockchainStruct{Lang: c.Lang, WaitText: waitText})
	return b.String(), nil
}
