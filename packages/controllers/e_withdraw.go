package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"fmt"
	"errors"
)

func (c *Controller) EWithdraw() (string, error) {

	if c.SessUserId == 0 {
		return "", errors.New(c.Lang["sign_up_please"])
	}

	c.r.ParseForm()
	currencyId := utils.StrToInt64(c.r.FormValue("currency_id"))
	if !utils.CheckInputData(c.r.FormValue("amount"), "amount") {
		return "", fmt.Errorf("incorrect amount")
	}
	method := c.r.FormValue("method")
	if !utils.CheckInputData(method, "method") {
		return "", fmt.Errorf("incorrect method")
	}
	account := c.r.FormValue("account")
	if !utils.CheckInputData(account, "account") {
		return "", fmt.Errorf("incorrect account")
	}
	amount := utils.StrToFloat64(c.r.FormValue("amount"))

	curTime := utils.Time()

	// нужно проверить, есть ли нужная сумма на счету юзера
	userAmount := utils.EUserAmountAndProfit(c.SessUserId, currencyId)
	if userAmount < amount {
		return "", fmt.Errorf("%s (%f<%f)", c.Lang["not_enough_money"], userAmount, amount)
	}
	if method != "Dcoin" && currencyId < 1000 {
		return "", fmt.Errorf("incorrect method")
	}

	err := userLock(c.SessUserId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	err = c.ExecSql(`UPDATE e_wallets SET amount = ?, last_update = ? WHERE user_id = ? AND currency_id = ?`, userAmount-amount, curTime, c.SessUserId, currencyId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	var commission float64
	if method == "Dcoin" {
		commission = utils.StrToFloat64(c.EConfig["dc_commission"])
	} else if method == "Perfect-money" {
		commission = utils.StrToFloat64(c.EConfig["pm_commission"])
	}
	wdAmount := utils.ClearNull(utils.Float64ToStr(amount * (1-commission/100)), 2)

	c.ExecSql(`INSERT INTO e_withdraw(open_time, user_id, currency_id, account, amount, wd_amount, method) VALUES (?, ?, ?, ?, ?, ?, ?)`, curTime, c.SessUserId, account, amount, wdAmount, method)
	userUnlock(c.SessUserId)

	return utils.JsonAnswer("success", c.Lang["request_is_created"]).String(), nil
}
