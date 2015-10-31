package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"errors"
)

func (c *Controller) ERedirect() (string, error) {

	c.r.ParseForm()
	ps := c.r.FormValue("FormExPs")
	token := c.r.FormValue("FormToken")
	amount := utils.StrToFloat64(c.r.FormValue("FormExAmount"))
	buyCurrencyId := utils.StrToInt64(c.r.FormValue("FormDC"))

	if !utils.CheckInputData(ps, "string") || !utils.CheckInputData(token, "string") {
		return "", errors.New("incorrect data")
	}

	// order_id занесем когда поуступят деньги в платежной системе
	err := c.ExecSql(`UPDATE e_tokens SET ps = ?, buy_currency_id = ?, amount_fiat = ? WHERE token = ?`, ps, buyCurrencyId, amount, token)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	tokenId, err := c.Single(`SELECT id FROM e_tokens WHERE token = ?`, token).String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	result := ""
	if ps == "pm" {
		result = `<form action="https://perfectmoney.is/api/step1.asp" method="POST" id="pm">
		<input type="hidden" name="PAYEE_ACCOUNT" value="U6198385">
		<input type="hidden" name="PAYEE_NAME" value="DcoinSimple">
		<input type="hidden" name="PAYMENT_ID" value="token-`+tokenId+`">
		<input type="hidden" name="PAYMENT_UNITS" value="USD">
		<input type="hidden" name="STATUS_URL" value="http://DcoinSimple.com/pm.php">
		<input type="hidden" name="PAYMENT_URL" value="http://DcoinSimple.com">
		<input type="hidden" name="PAYMENT_URL_METHOD" value="LINK">
		<input type="hidden" name="NOPAYMENT_URL" value="http://DcoinSimple.com">
		<input type="hidden" name="NOPAYMENT_URL_METHOD" value="LINK">
		<input type="hidden" name="SUGGESTED_MEMO" value="Dcoins">
		<input type="hidden" name="BAGGAGE_FIELDS" value="">
		<input type="hidden" name="PAYMENT_AMOUNT" value="`+utils.Float64ToStr(amount)+`">
	       </form>`
	} else if ps == "mobile" {
		result = `<form name="payment" method="post" action="https://sci.interkassa.com/" enctype="utf-8" id="pm">
		<input type="hidden" name="ik_co_id" value="53cfd5e2bf4efc831c9fd661" />
		<input type="hidden" name="ik_pm_no" value="ID_4233" />
		<input type="hidden" name="ik_cur" value="USD" />
		<input type="hidden" name="ik_desc" value="token-`+tokenId+`" />
		<input type="hidden" name="ik_am" value="`+utils.Float64ToStr(amount)+`" >
		</form>`
	}
	return result, nil

}
