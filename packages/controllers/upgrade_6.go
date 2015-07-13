package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"log"
	"strings"
)

type upgrade6Page struct {
	Alert string
	UserId int64
	Lang map[string]string
	GeolocationLat string
	GeolocationLon string
	SaveAndGotoStep string
	UpgradeMenu string
}

func (c *Controller) Upgrade6() (string, error) {

	log.Println("Upgrade6")

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


	saveAndGotoStep := strings.Replace(c.Lang["save_and_goto_step"], "[num]", "7", -1)
	upgradeMenu := utils.MakeUpgradeMenu(6)

	TemplateStr, err := makeTemplate("upgrade_6", "upgrade6", &upgrade6Page {
		Alert: c.Alert,
		Lang: c.Lang,
		SaveAndGotoStep: saveAndGotoStep,
		UpgradeMenu: upgradeMenu,
		GeolocationLat: geolocationLat,
		GeolocationLon: geolocationLon,
		UserId: c.SessUserId})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

