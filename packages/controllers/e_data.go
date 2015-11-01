package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"errors"
	"encoding/json"
)

func (c *Controller) EData() (string, error) {

	c.w.Header().Set("Access-Control-Allow-Origin", "*")

	c.r.ParseForm()
	token := c.r.FormValue("token")
	if !utils.CheckInputData(token, "string") {
		return "", errors.New("incorrect token")
	}
	// сколько всего продается dUSD
	tokenMap, err := c.OneRow(`SELECT * FROM e_tokens WHERE token = ?`, token).String()
	if err!=nil {
		return "", utils.ErrInfo(err)
	}
	wallets, err := c.GetAll(`SELECT * FROM e_wallets WHERE user_id = ?`, 100, tokenMap["user_id"])
	if err!=nil {
		return "", utils.ErrInfo(err)
	}
	orders, err := c.GetAll(`SELECT * FROM e_orders WHERE user_id = ? ORDER BY time DESC LIMIT 10`, 100, tokenMap["user_id"])
	if err!=nil {
		return "", utils.ErrInfo(err)
	}
	withdraw, err := c.GetAll(`SELECT * FROM e_withdraw WHERE user_id = ? ORDER BY open_time DESC LIMIT 10`, 100, tokenMap["user_id"])
	if err!=nil {
		return "", utils.ErrInfo(err)
	}

	//print json_encode(array('token'=>$token, 'wallets'=>$wallets, 'orders'=>$orders, 'withdraw'=>$withdraw));
	jsonData, err := json.Marshal(&EInfoResult{token: tokenMap, wallets: wallets, orders: orders, withdraw: withdraw})
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	return string(jsonData), nil
}
