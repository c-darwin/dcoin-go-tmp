package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"errors"
	"encoding/base64"
)

type ERedirectPage struct {
	Lang           map[string]string
	EConfig map[string]string
	TokenId string
	Amount string
	EURL string
	MDesc string
}

func (c *Controller) ERedirect() (string, error) {

	c.r.ParseForm()
	token := c.r.FormValue("FormToken")
	amount := c.r.FormValue("FormExAmount")
	buyCurrencyId := utils.StrToInt64(c.r.FormValue("FormDC"))

	if !utils.CheckInputData(token, "string") {
		return "", errors.New("incorrect data")
	}

	// order_id занесем когда поуступят деньги в платежной системе
	err := c.ExecSql(`UPDATE e_tokens SET buy_currency_id = ?, amount_fiat = ? WHERE token = ?`, buyCurrencyId, utils.StrToFloat64(c.r.FormValue("FormExAmount")), token)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	tokenId, err := c.Single(`SELECT id FROM e_tokens WHERE token = ?`, token).String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	TemplateStr, err := makeTemplate("e_redirect", "eRedirect", &ERedirectPage{
		Lang:           c.Lang,
		EConfig : c.EConfig,
		TokenId: tokenId,
		EURL: c.EURL,
		MDesc:  base64.StdEncoding.EncodeToString(utils.Int64ToByte(c.SessUserId)),
		Amount: amount})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
