package controllers
import (
	"utils"
	"log"
	"strings"
	"os"
)

type upgrade1Page struct {
	Alert string
	SignData string
	ShowSignData bool
	CountSignArr []int
	UserId int64
	Lang map[string]string
	SaveAndGotoStep string
	UpgradeMenu string
	Step string
	NextStep string
	PhotoType string
	Photo string
}

func (c *Controller) Upgrade1() (string, error) {

	log.Println("Upgrade1")

	userFace := ""
	/*userProfile := ""

	path := "public/"+utils.Int64ToStr(c.SessUserId)+"_user_profile.jpg"
	if _, err := os.Stat(path); err == nil {
		userProfile = path
	}*/

	path := "public/"+utils.Int64ToStr(c.SessUserId)+"_user_face.jpg"
	if _, err := os.Stat(path); err == nil {
		userFace = path
	}

	step := "1"
	nextStep := "2"
	photoType := "face"
	photo := userFace

	saveAndGotoStep := strings.Replace(c.Lang["save_and_goto_step"], "[num]", "2", -1)
	upgradeMenu := utils.MakeUpgradeMenu(1)

	TemplateStr, err := makeTemplate("upgrade_1_and_2", "upgrade1And2", &upgrade1Page{
		Alert: c.Alert,
		Lang: c.Lang,
		CountSignArr: c.CountSignArr,
		ShowSignData: c.ShowSignData,
		SaveAndGotoStep: saveAndGotoStep,
		UpgradeMenu: upgradeMenu,
		UserId: c.SessUserId,
		PhotoType: photoType,
		Photo: photo,
		Step: step,
		NextStep: nextStep})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

