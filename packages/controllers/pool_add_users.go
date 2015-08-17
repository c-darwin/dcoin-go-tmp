package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"errors"
	"encoding/json"
	"regexp"
	"bytes"
	"io"
	"github.com/c-darwin/dcoin-go-tmp/packages/schema"
)
type JsonBackup struct {
	Community []string `json:"community"`
	Data map[string][]map[string]string `json:"data"`
}

func (c *Controller) PoolAddUsers() (string, error) {

	if !c.NodeAdmin || c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	if !c.Community {
		return "", utils.ErrInfo(errors.New("Single mode"))
	}

	c.r.ParseMultipartForm(32 << 20)
	file, _, err := c.r.FormFile("file")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	defer file.Close()
	//log.Debug("", buffer.String())

	var mainMap JsonBackup
	err = json.Unmarshal(buffer.Bytes(), &mainMap)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	log.Debug("Unmarshal ok")

	schema_ := &schema.SchemaStruct{}
	schema_.DCDB = c.DCDB
	schema_.DbType = c.ConfigIni["db_type"]
	for i:=0; i <len(mainMap.Community); i++ {
		schema_.PrefixUserId = utils.StrToInt(mainMap.Community[i])
		schema_.GetSchema()
		c.ExecSql(`INSERT INTO community (user_id) VALUES (?)`, mainMap.Community[i])
	}

	allTables, err := c.GetAllTables()

	for table, arr := range mainMap.Data {
		if !utils.InSliceString(table, allTables) {
			continue
		}
		//_ = c.ExecSql(`DROP TABLE `+table)
		//if err != nil {
		//	return "", utils.ErrInfo(err)
		//}
		log.Debug(table)
		for i, data := range arr {
			log.Debug("%v", i)
			colNames := ""
			values := []interface {} {}
			qq := ""
			for name, value := range data {

				if ok, _ := regexp.MatchString("my_table", table); ok{
					if name == "host" {
						name = "http_host"
					}
				}
				if name == "show_progressbar" {
					name = "show_progress_bar"
				}

				colNames += name+","
				values = append(values, value)
				if ok, _ := regexp.MatchString("(hash_code|public_key|encrypted)", name); ok{
					qq+="[hex],"
				} else {
					qq+="?,"
				}
			}
			colNames = colNames[0:len(colNames)-1]
			qq = qq[0:len(qq)-1]
			query := `INSERT INTO `+table+` (`+colNames+`) VALUES (`+qq+`)`
			log.Debug("%v", query)
			log.Debug("%v", values)
			err = c.ExecSql(query, values...)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}
	}


	return "", nil
}
