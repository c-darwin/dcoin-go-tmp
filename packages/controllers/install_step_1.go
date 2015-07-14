package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

type installStep1Struct struct {
	Lang map[string]string
}

// Шаг 1 - выбор либо стандартных настроек (sqlite и блокчейн с сервера) либо расширенных - pg/mysql и загрузка с нодов
func (c *Controller) InstallStep1() (string, error) {

	TemplateStr, err := makeTemplate("install_step_1", "installStep1", &installStep0Struct{
		Lang: c.Lang})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
