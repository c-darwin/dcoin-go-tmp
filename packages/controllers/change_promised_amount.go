package controllers
import (
	"github.com/c-darwin/dcoin-tmp/packages/utils"
	"time"
	"fmt"
)

type changePromisedAmountPage struct {
	Alert string
	SignData string
	ShowSignData bool
	TxType string
	TxTypeId int64
	TimeNow int64
	UserId int64
	Lang map[string]string
	CountSignArr []int
	LastTxFormatted string
	CurrencyList map[int64]string
	PromisedAmountId string
	Amount string
}

func (c *Controller) ChangePromisedAmount() (string, error) {

	txType := "ChangePromisedAmount";
	txTypeId := utils.TypeInt(txType)
	timeNow := time.Now().Unix()
	promisedAmountId := c.Parameters["promised_amount_id"]
	amount := c.Parameters["amount"]

	TemplateStr, err := makeTemplate("change_promised_amount", "changePromisedAmount", &changePromisedAmountPage {
		Alert: c.Alert,
		Lang: c.Lang,
		CountSignArr: c.CountSignArr,
		ShowSignData: c.ShowSignData,
		UserId: c.SessUserId,
		TimeNow: timeNow,
		TxType: txType,
		TxTypeId: txTypeId,
		SignData: fmt.Sprintf("%d,%d,%d,%s,%s", txTypeId, timeNow, c.SessUserId, promisedAmountId, amount),
		PromisedAmountId: promisedAmountId,
		Amount: amount})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

