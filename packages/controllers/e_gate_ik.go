package controllers

import (
	"fmt"
	"strings"
	//"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"sort"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"encoding/base64"
)

func (c *Controller) EGateIk() (string, error) {

	c.r.ParseForm()
	var body []byte
	fmt.Println(c.r.Body.Read(body))
	fmt.Println(c.r.Form)
	var ikNames []string
	for name, _ := range c.r.Form {
		if name[:2] == "ik" && name!="ik_sign" {
			ikNames = append(ikNames, name)
		}
	}
	sort.Strings(ikNames)
	fmt.Println(ikNames)

	var ikValues []string
	for _, names := range ikNames {
		ikValues = append(ikValues, c.r.Form[names][0])
	}
	ikValues = append(ikValues, "aZNwrkXK1eptH3qM")
	fmt.Println(ikValues)
	sign := strings.Join(ikValues, ":")
	fmt.Println(sign)
	sign = base64.StdEncoding.EncodeToString(utils.HexToBin(utils.Md5(sign)))
	fmt.Println(sign)


	//map[ik_co_rfn:[100] ik_cur:[RUB] ik_co_prs_id:[406295558666] ik_inv_crt:[2015-11-03 19:35:46] ik_inv_prc:[2015-11-03 19:35:46] ik_inv_st:[success] ik_pm_no:[ID_4233] ik_am:[100] ik_trn_id:[] ik_pw_via:[test_interkassa_test_xts] ik_ps_price:[103] ik_desc:[Event Description] ik_sign:[FKIE3rf+ULdnh1AZAYzjKw==] controllerName:[EGateIk] ik_co_id:[560ecc4e3d1eaf52348b4567] ik_inv_id:[41977100]]

/*
	sign := strings.ToUpper(utils.Md5(c.r.FormValue("ik_co_rfn")+":"+c.r.FormValue("ik_cur")+":"+c.r.FormValue("ik_co_prs_id")+":"+c.r.FormValue("ik_inv_crt")+":"+c.r.FormValue("ik_inv_prc")+":"+c.r.FormValue("ik_inv_st")+":"+c.r.FormValue("ik_pm_no")+":"+c.r.FormValue("ik_am")+":"+c.r.FormValue("ik_trn_id")+":"+c.r.FormValue("ik_pw_via")+":"+c.r.FormValue("ik_ps_price")+":"+c.r.FormValue("ik_desc")+":"+strings.ToUpper(utils.Md5(c.EConfig["pm_s_key"]))+":"+c.r.FormValue("TIMESTAMPGMT")))
/*
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
