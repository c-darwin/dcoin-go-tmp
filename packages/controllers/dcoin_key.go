package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"strings"
	"errors"
	"encoding/base64"
)

func (c *Controller) DcoinKey() (string, error) {

	c.r.ParseForm()

	paramNoPass := utils.ParamType{X: 176, Y: 100, Width: 100, Bg_path: "static/img/k_bg.png"}
	paramPass := utils.ParamType{X: 167, Y: 93, Width: 118, Bg_path: "static/img/k_bg_pass.png"}

	privKey, _ := utils.GenKeys()

	var param utils.ParamType
	var privateKey string
	if len(c.r.FormValue("password")) > 0 {
		privateKey_, err := utils.Encrypt(utils.Md5(c.r.FormValue("password")), []byte(privKey))
		privateKey = base64.StdEncoding.EncodeToString(privateKey_)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		param = paramPass
	} else {
		privateKey = strings.Replace(privKey, "-----BEGIN RSA PRIVATE KEY-----", "", -1)
		privateKey = strings.Replace(privateKey, "-----END RSA PRIVATE KEY-----", "", -1)
		param = paramNoPass
	}

	//if ok, _ := regexp.MatchString("(iPod|iPhone|iPad)", c.r.UserAgent()); ok{
	if len(c.r.FormValue("iPhone")) > 0 {
		buffer, err := utils.KeyToImg(privateKey, "", c.SessUserId, c.TimeFormat, param)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		c.w.Header().Set("Content-Type", "image/png")
		c.w.Header().Set("Content-Length", utils.IntToStr(len(buffer.Bytes())))
		c.w.Header().Set("Content-Disposition", `attachment; filename="Dcoin-private-key-`+utils.Int64ToStr(c.SessUserId)+`.png`)
		if _, err := c.w.Write(buffer.Bytes()); err != nil {
			return "", utils.ErrInfo(errors.New("unable to write image"))
		}
	} else {
		c.w.Header().Set("Content-Type", "text/plain")
		c.w.Header().Set("Content-Length", utils.IntToStr(len(privateKey)))
		c.w.Header().Set("Content-Disposition", `attachment; filename="Dcoin-private-key-`+utils.Int64ToStr(c.SessUserId)+`.txt`)
		if _, err := c.w.Write([]byte(privateKey)); err != nil {
			return "", utils.ErrInfo(errors.New("unable to write text"))
		}
	}

	return "", nil
}


