package controllers

import (
	"errors"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"encoding/json"
	"regexp"
	"time"
)

func (c *Controller) EDataBaseDump() (string, error) {

	if !c.NodeAdmin || c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	allTables, err := c.GetAllTables()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	c.r.ParseForm()

	mainMap := make(map[string][]map[string]string)

	for _, table := range allTables {
		re := regexp.MustCompile("^e_")
		match := re.FindStringSubmatch(table)
		if len(match) > 0 {
			data, err := c.GetAll(`SELECT * FROM `+table, -1)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			mainMap[table] = data
		}
	}

	jsonData, _ := json.Marshal(mainMap)
	log.Debug(string(jsonData))

	c.w.Header().Set("Content-Type", "text/plain")
	c.w.Header().Set("Content-Length", utils.IntToStr(len(jsonData)))
	t := time.Unix(utils.Time(), 0)
	c.w.Header().Set("Content-Disposition", `attachment; filename="dcoin_e_backup-`+t.Format(c.TimeFormat)+`.txt`)
	if _, err := c.w.Write(jsonData); err != nil {
		return "", utils.ErrInfo(errors.New("unable to write text"))
	}

	return "", nil
}
