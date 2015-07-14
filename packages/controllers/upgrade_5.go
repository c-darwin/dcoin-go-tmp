package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"log"
	"strings"
)

type upgrade5Page struct {
	Alert string
	UserId int64
	Lang map[string]string
	HttpHost string
	TcpHost string
	SaveAndGotoStep string
	UpgradeMenu string
	Community bool
}

func (c *Controller) Upgrade5() (string, error) {

	log.Println("Upgrade5")

	data, err := c.OneRow("SELECT http_host, tcp_host FROM "+c.MyPrefix+"my_table").String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	if c.Community && len(data) == 0 {
		data, err = c.OneRow("SELECT http_host, tcp_host FROM miners_data WHERE user_id  =  ?", c.PoolAdminUserId).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	}

	saveAndGotoStep := strings.Replace(c.Lang["save_and_goto_step"], "[num]", "6", -1)
	upgradeMenu := utils.MakeUpgradeMenu(5)

	TemplateStr, err := makeTemplate("upgrade_5", "upgrade5", &upgrade5Page {
		Alert: c.Alert,
		Lang: c.Lang,
		SaveAndGotoStep: saveAndGotoStep,
		UpgradeMenu: upgradeMenu,
		HttpHost: data["http_host"],
		TcpHost: data["tcp_host"],
		Community: c.Community,
		UserId: c.SessUserId})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

