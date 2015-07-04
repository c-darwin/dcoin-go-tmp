package controllers
import (
	"errors"
	"utils"
	"os"
)

func (c *Controller) DeleteVideo() (string, error) {

	if c.SessUserId == 0 || c.SessRestricted != 0 {
		return "", errors.New("Permission denied")
	}

	if c.r.FormValue("type") == "mp4" {
		err := os.Remove("public/"+utils.Int64ToStr(c.SessUserId)+"_user_video.mp4")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	} else if c.r.FormValue("type") == "webm_ogg" {
		err := os.Remove("public/"+utils.Int64ToStr(c.SessUserId)+"_user_video.ogv")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		err = os.Remove("public/"+utils.Int64ToStr(c.SessUserId)+"_user_video.webm")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	}

	return ``, nil
}
