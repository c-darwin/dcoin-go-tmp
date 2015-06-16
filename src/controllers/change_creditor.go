package controllers
import (
	"utils"
	"time"
	"log"
)

type changeCreditorPage struct {
	SignData string
	ShowSignData bool
	TxType string
	TxTypeId int64
	TimeNow int64
	UserId int64
	Alert string
	Lang map[string]string
	CountSignArr []int
	CreditId float64
}

func (c *Controller) ChangeCreditor() (string, error) {

	log.Println("ChangeCreditor")

	txType := "change_creditor";
	txTypeId := utils.TypeInt(txType)
	timeNow := time.Now().Unix()

	log.Println("c.Parameters", c.Parameters)
	creditId := utils.Round(utils.StrToFloat64(c.Parameters["credit_id"]), 0)
	log.Println("creditId", creditId)

	TemplateStr, err := makeTemplate("change_creditor", "changeCreditor", &changeCreditorPage{
		Alert: c.Alert,
		Lang: c.Lang,
		CountSignArr: c.CountSignArr,
		ShowSignData: c.ShowSignData,
		UserId: c.SessUserId,
		TimeNow: timeNow,
		TxType: txType,
		TxTypeId: txTypeId,
		SignData: "",
		CreditId: creditId})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

