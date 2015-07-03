package controllers
import (
	"utils"
	"log"
	"strings"
)

type upgrade4Page struct {
	Alert string
	UserId int64
	Lang map[string]string
	VideoUrl string
	SaveAndGotoStep string
	UpgradeMenu string
}

func (c *Controller) Upgrade4() (string, error) {

	log.Println("Upgrade4")

	videoUrl := ""

	// есть ли загруженное видео.
	data, err := c.Single("SELECT video_url_id, video_type FROM "+c.MyPrefix+"my_table").String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	switch data["video_type"] {
	case "youtube":
		videoUrl = "http://www.youtube.com/embed/"+data["video_url_id"]
	case "vimeo":
		videoUrl = "http://www.vimeo.com/embed/"+data["video_url_id"]
	case "youku":
		videoUrl = "http://www.youku.com/embed/"+data["video_url_id"]
	}

	saveAndGotoStep := strings.Replace(c.Lang["save_and_goto_step"], "[num]", "5", -1)
	upgradeMenu := utils.MakeUpgradeMenu(4)


	TemplateStr, err := makeTemplate("upgrade_4", "upgrade4", &upgrade4Page{
		Alert: c.Alert,
		Lang: c.Lang,
		SaveAndGotoStep: saveAndGotoStep,
		UpgradeMenu: upgradeMenu,
		VideoUrl: videoUrl,
		UserId: c.SessUserId})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

