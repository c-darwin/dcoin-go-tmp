package controllers
import (
	"utils"
	"log"
	"strings"
)


func (c *Controller) DcoinKey() (string, error) {

//	var err error
	log.Println("DcoinKey")

	param := utils.ParamType{X: 176, Y: 100, Width: 100, Bg_path: "static/img/k_bg.png"}

	privKey, _ := utils.GenKeys()
	privKey = strings.Replace(privKey, "-----BEGIN RSA PRIVATE KEY-----", "", -1)
	privKey = strings.Replace(privKey, "-----END RSA PRIVATE KEY-----", "", -1)
	buffer, err := utils.KeyToImg(privKey, "", c.SessUserId, c.TimeFormat, param)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	c.w.Header().Set("Content-Type", "image/png")
	c.w.Header().Set("Content-Length", utils.IntToStr(len(buffer.Bytes())))
	if _, err := c.w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}

	return "", nil
}


