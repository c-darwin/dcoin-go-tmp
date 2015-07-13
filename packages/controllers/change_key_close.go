package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

type changeKeyClosePage struct {
	Alert string
	SignData string
	ShowSignData bool
	CountSignArr []int
	Lang map[string]string
	UserId int64
	TxType string
	TxTypeId int64
	TimeNow int64
}

func (c *Controller) ChangeKeyClose() (string, error) {

	txType := "ChangeKeyClose";
	txTypeId := utils.TypeInt(txType)
	timeNow := utils.Time()

	TemplateStr, err := makeTemplate("change_key_close", "changeKeyClose", &changeKeyClosePage{
		Alert: c.Alert,
		Lang: c.Lang,
		ShowSignData: c.ShowSignData,
		SignData: "",
		UserId: c.SessUserId,
		CountSignArr: c.CountSignArr,
		TimeNow: timeNow,
		TxType: txType,
		TxTypeId: txTypeId})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
