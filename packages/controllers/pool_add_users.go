package controllers
import (
	"github.com/c-darwin/dcoin-tmp/packages/utils"
	"errors"
	"log"
	"encoding/json"
	"regexp"
	"bytes"
	"io"
)

func (c *Controller) PoolAddUsers() (string, error) {

	if !c.NodeAdmin || c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
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
	log.Println(buffer.String())

	mainMap := make(map[string][]map[string]string)

	err = json.Unmarshal(buffer.Bytes(), &mainMap)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	for table, arr := range mainMap {
		err = c.ExecSql(`DELETE FROM `+table)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		log.Println(table)
		for i, data := range arr {
			log.Println(i)
			colNames := ""
			values := []interface {} {}
			qq := ""
			for name, value := range data {
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
			log.Println(query)
			log.Println(values)
			err = c.ExecSql(query, values...)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}
	}


	return "", nil
}
