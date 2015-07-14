package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

type installStep0Struct struct {
	Lang map[string]string
}

// Шаг 0 - выбор языка
func (c *Controller) InstallStep0() (string, error) {

	TemplateStr, err := makeTemplate("install_step_0", "installStep0", &installStep0Struct{
		Lang: c.Lang})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}