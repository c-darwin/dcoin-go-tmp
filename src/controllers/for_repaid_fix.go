package controllers
import (
	"utils"
	"time"
	"log"
	"fmt"
)

type forRepaidFixPage struct {
	Alert string
	SignData string
	ShowSignData bool
	TxType string
	TxTypeId int64
	TimeNow int64
	UserId int64
	Lang map[string]string
	CountSignArr []int
}

func (c *Controller) ForRepaidFix() (string, error) {

	log.Println("ForRepaidFix")

	txType := "for_repaid_fix";
	txTypeId := utils.TypeInt(txType)
	timeNow := time.Now().Unix()

	TemplateStr, err := makeTemplate("for_repaid_fix", "forRepaidFix", &forRepaidFixPage{
		Alert: c.Alert,
		Lang: c.Lang,
		CountSignArr: c.CountSignArr,
		ShowSignData: c.ShowSignData,
		UserId: c.SessUserId,
		TimeNow: timeNow,
		TxType: txType,
		TxTypeId: txTypeId,
		SignData: fmt.Sprintf("%d,%d,%d", txTypeId, timeNow, c.SessUserId)})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

