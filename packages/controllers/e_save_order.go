package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"errors"
	"strings"
)

func (c *Controller) ESaveOrder() (string, error) {

	c.r.ParseForm()
	sellCurrencyId := utils.StrToInt64(c.r.FormValue("sell_currency_id"))
	buyCurrencyId := utils.StrToInt64(c.r.FormValue("buy_currency_id"))
	amount := utils.StrToFloat64(c.r.FormValue("amount"))
	sellRate := utils.StrToFloat64(c.r.FormValue("sell_rate"))
	orderType := utils.StrToFloat64(c.r.FormValue("type"))
	// можно ли торговать такими валютами
	checkCurrency, err := c.Single("SELECT count(id) FROM e_currency WHERE id IN (?, ?)", sellCurrencyId, buyCurrencyId).Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	if checkCurrency!=2 {
		return "", errors.New("Currency error")
	}
	if orderType!="sell" || orderType!="buy" {
		return "", errors.New("Type error")
	}
	if amount == 0 {
		return "", errors.New(c.Lang["amount_error"])
	}
	if amount < 0.001 && sellCurrencyId<1000 {
		return "", errors.New(strings.Replace(c.Lang["save_order_min_amount"], "[amount]", "0.001", -1))
	}
	if sellRate < 0.0001 {
		return "", errors.New(strings.Replace(c.Lang["save_order_min_price"], "[price]", "0.0001", -1))
	}
	geReductionLock, err = geReductionLock()
	if len(geReductionLock) > 0 {
		return "", errors.New(strings.Replace(c.Lang["creating_orders_unavailable"], "[minutes]", "30", -1))
	}

	// нужно проверить, есть ли нужная сумма на счету юзера
	userAmountAndProfit := userAmountAndProfit(c.SessUserId, sellCurrencyId)
	if userAmountAndProfit < amount {
		return "", errors.New(c.Lang["not_enough_money"]+" ("+utils.StrToFloat64(userAmountAndProfit)+"<"+utils.StrToFloat64(amount)+")"+strings.Replace(c.Lang["add_funds_link"], "[currency]", "USD", -1))
	}

	err = c.newForexOrder(c.SessUserId, amount, sellRate, sellCurrencyId, buyCurrencyId, orderType)
	if err!=nil {
		return nil, err
	} else {
		utils.JsonAnswer("success", c.Lang["order_created"])
	}

	return ``, nil
}
