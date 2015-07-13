package controllers
import (
	"github.com/c-darwin/dcoin-tmp/packages/utils"
	"time"
	"fmt"
)

type PromisedAmountActualizationPage struct {
	SignData string
	ShowSignData bool
	TxType string
	TxTypeId int64
	TimeNow int64
	UserId int64
	Alert string
	Lang map[string]string
	CountSignArr []int
}

func (c *Controller) PromisedAmountActualization() (string, error) {

	txType := "ActualizationPromisedAmounts";
	txTypeId := utils.TypeInt(txType)
	timeNow := time.Now().Unix()

	TemplateStr, err := makeTemplate("promised_amount_actualization", "promisedAmountsActualization", &PromisedAmountActualizationPage{
		Alert: c.Alert,
		Lang: c.Lang,
		CountSignArr: c.CountSignArr,
		ShowSignData: c.ShowSignData,
		UserId: c.SessUserId,
		TimeNow: timeNow,
		TxType: txType,
		TxTypeId: txTypeId,
		SignData: fmt.Sprintf("%v,%v,%v", txTypeId, timeNow, c.SessUserId)})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

