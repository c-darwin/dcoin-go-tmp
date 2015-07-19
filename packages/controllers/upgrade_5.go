package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"

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

	log.Debug("Upgrade5")

	data, err := c.OneRow("SELECT http_host, tcp_host FROM "+c.MyPrefix+"my_table").String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	log.Debug("c.Community: %v", c.Community)
	log.Debug("c.PoolAdminUserId: %v", c.PoolAdminUserId)
	if c.Community && len(data["http_host"]) == 0 && len(data["tcp_host"]) == 0 {
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

