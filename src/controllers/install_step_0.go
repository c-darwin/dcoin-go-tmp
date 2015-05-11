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

type installStep0Struct struct {
	Lang map[string]string
}

// Шаг 0 - выбор языка
func (c *Controller) Install_step_0() (string, error) {

	fmt.Println("Install_step_0")

	//sql := schema.GetSchema("sqlite", 0)
	/*sql := `CREATE TABLE "555555_my_cf_funding" (
"id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
"from_user_id" bigint(20) NOT NULL,
"project_id" bigint(20) NOT NULL,
"amount" decimal(15,2) NOT NULL,
"del_block_id" int(11)  NOT NULL,
"time" int(10)  NOT NULL,
"block_id" int(10)  NOT NULL,
"comment" text NOT NULL,
"comment_status" varchar(100)  NOT NULL DEFAULT 'decrypted'
);
`*/
	/*fmt.Println(sql)
	res, err := c.DCDB.Exec(sql)
	if err != nil {
		fmt.Println("err", err)
	}
	fmt.Println("res", res)*/

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
