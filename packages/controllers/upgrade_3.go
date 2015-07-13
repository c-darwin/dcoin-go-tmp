package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
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
	FaceCoords string
	ProfileCoords string
	UserProfile string
	UserFace string
	ExamplePoints map[string]string
}

func (c *Controller) Upgrade3() (string, error) {

	log.Println("Upgrade3")

	userProfile := "public/"+utils.Int64ToStr(c.SessUserId)+"_user_profile.jpg"
	userFace := "public/"+utils.Int64ToStr(c.SessUserId)+"_user_face.jpg"

	if _, err := os.Stat(userProfile); os.IsNotExist(err) {
		userProfile = ""
	} else {
		userProfile = "public/"+utils.Int64ToStr(c.SessUserId)+"_user_profile.jpg?r="+utils.IntToStr(utils.RandInt(0, 99999))
	}
	if _, err := os.Stat(userFace); os.IsNotExist(err) {
		userFace = ""
	} else {
		userFace = "public/"+utils.Int64ToStr(c.SessUserId)+"_user_face.jpg?r="+utils.IntToStr(utils.RandInt(0, 99999))
	}

	// текущий набор точек для шаблонов
	examplePoints, err := c.GetPoints(c.Lang)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// точки, которые юзер уже отмечал
	data, err := c.OneRow("SELECT face_coords, profile_coords FROM "+c.MyPrefix+"my_table").String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	faceCoords := ""
	profileCoords := ""
	if len(data["face_coords"]) > 0 {
		faceCoords = data["face_coords"]
		profileCoords = data["profile_coords"]
	}

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
		FaceCoords: faceCoords,
		ProfileCoords: profileCoords,
		UserProfile: userProfile,
		UserFace: userFace,
		ExamplePoints: examplePoints})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

