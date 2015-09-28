package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"strings"
)

type upgrade5Page struct {
	Alert string
	UserId int64
	Lang map[string]string
	GeolocationLat string
	GeolocationLon string
	SaveAndGotoStep string
	UpgradeMenu string
	Mobile bool
}

func (c *Controller) Upgrade5() (string, error) {

	log.Debug("Upgrade5")

	geolocationLat := ""
	geolocationLon := ""
	geolocation, err := c.Single("SELECT geolocation FROM "+c.MyPrefix+"my_table").String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	if len(geolocation) > 0 {
		x := strings.Split(geolocation, ", ")
		if len(x) == 2 {
			geolocationLat = x[0]
			geolocationLon = x[1]
		}
	}

	saveAndGotoStep := strings.Replace(c.Lang["save_and_goto_step"], "[num]", "6", -1)
	upgradeMenu := utils.MakeUpgradeMenu(5)

	TemplateStr, err := makeTemplate("upgrade_5", "upgrade5", &upgrade5Page {
		Alert: c.Alert,
		Lang: c.Lang,
		SaveAndGotoStep: saveAndGotoStep,
		UpgradeMenu: upgradeMenu,
		GeolocationLat: geolocationLat,
		GeolocationLon: geolocationLon,
		UserId: c.SessUserId,
		Mobile: utils.Mobile()})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

