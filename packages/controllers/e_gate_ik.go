package controllers

import (
	"fmt"
)

func (c *Controller) EGateIk() (string, error) {

	c.r.ParseForm()
	var body []byte
	fmt.Println(c.r.Body.Read(body))
	fmt.Println(c.r.Form)
	/*
	sign := strings.ToUpper(utils.Md5(c.r.FormValue("PAYMENT_ID")+":"+c.r.FormValue("PAYEE_ACCOUNT")+":"+c.r.FormValue("PAYMENT_AMOUNT")+":"+c.r.FormValue("PAYMENT_UNITS")+":"+c.r.FormValue("PAYMENT_BATCH_NUM")+":"+c.r.FormValue("PAYER_ACCOUNT")+":"+strings.ToUpper(utils.Md5(c.EConfig["pm_s_key"]))+":"+c.r.FormValue("TIMESTAMPGMT")))

	txTime := utils.StrToInt64(c.r.FormValue("TIMESTAMPGMT"));

	if sign != c.r.FormValue("V2_HASH") {
		return "", errors.New("Incorrect signature")
	}

	currencyId := int64(0)

	if c.r.FormValue("PAYMENT_UNITS") == "USD" {
		currencyId = 1001
	}
	amount := utils.StrToFloat64(c.r.FormValue("PAYMENT_AMOUNT"))
	pmId := utils.StrToInt64(c.r.FormValue("PAYMENT_BATCH_NUM"))
	// проверим, не зачисляли ли мы уже это платеж
	existsId, err := c.Single(`SELECT id FROM e_adding_funds_pm WHERE id = ?`, pmId).Int64()
	if err!=nil {
		return "", utils.ErrInfo(err)
	}
	if existsId != 0 {
		return "", errors.New("Incorrect PAYMENT_BATCH_NUM")
	}
	paymentInfo := c.r.FormValue("PAYMENT_ID")

	err = EPayment(paymentInfo, currencyId, txTime, amount, pmId, "pm", c.ECommission)
	if err!=nil {
		return "", utils.ErrInfo(err)
	}
*/
	return ``, nil
}
