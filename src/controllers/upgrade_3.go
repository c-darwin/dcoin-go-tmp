package controllers
import (
	"utils"
	"log"
	"strings"
	"os"
)

type upgrade3Page struct {
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

func (c *Controller) Upgrade3() (string, error) {

	log.Println("Upgrade3")

	userProfile := "public/"+utils.Int64ToStr(c.SessUserId)+"_user_profile.jpg"
	userFace := "public/"+utils.Int64ToStr(c.SessUserId)+"_user_face.jpg"

	if _, err := os.Stat(userProfile); os.IsNotExist(err) {
		userProfile = ""
	}
	if _, err := os.Stat(userFace); os.IsNotExist(err) {
		userFace = ""
	}

	// текущий набор точек для шаблонов
	examplePoints := c.GetPoints(c.Lang)

	// точки, которые юзер уже отмечал


	step := "2"
	nextStep := "3"
	photoType := "profile"
	photo := userProfile

	saveAndGotoStep := strings.Replace(c.Lang["save_and_goto_step"], "[num]", "4", -1)
	upgradeMenu := utils.MakeUpgradeMenu(3)

	TemplateStr, err := makeTemplate("upgrade_3", "upgrade3", &upgrade3Page{
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

