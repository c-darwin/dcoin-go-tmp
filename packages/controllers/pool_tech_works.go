package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

type PoolTechWorksPage struct {
	Alert string
	Lang map[string]string
	MyModalIdName string
	Ver string
}

func (c *Controller) PoolTechWorks() (string, error) {

	info, err := c.GetInfoBlock()

	TemplateStr, err := makeTemplate("pool_tech_works", "poolTechWorks", &PoolTechWorksPage {
		Alert: c.Alert,
		Lang: c.Lang,
		Ver: info["current_version"],
		MyModalIdName: "myModalLogin"})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

