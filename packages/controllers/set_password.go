package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

type setPasswordPage struct {
	Lang map[string]string
}

func (c *Controller) SetPassword() (string, error) {

	TemplateStr, err := makeTemplate("set_password", "setPassword", &setPasswordPage {
		Lang: c.Lang})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}


