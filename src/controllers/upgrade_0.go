package controllers
import (
	"utils"
//	"time"
	"log"
	"consts"
)

type upgrade0Page struct {
	Alert string
	SignData string
	ShowSignData bool
	CountSignArr []int
	UserId int64
	Lang map[string]string
	Countries []string
	Country int
	Race int
}

func (c *Controller) Upgrade0() (string, error) {

	log.Println("Upgrade0")

	data, err := c.OneRow("SELECT race, country FROM "+c.MyPrefix+"my_table").Int()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	race := data["race"]
	country := 0
	if race > 0 {
		country = data["country"]
	}
	TemplateStr, err := makeTemplate("upgrade_0", "Upgrade0", &upgrade0Page{
		Alert: c.Alert,
		Lang: c.Lang,
		CountSignArr: c.CountSignArr,
		ShowSignData: c.ShowSignData,
		UserId: c.SessUserId,
		Countries: consts.Countries,
		Country: country,
		Race: race})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

