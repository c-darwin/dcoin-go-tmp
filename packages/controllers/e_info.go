package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"errors"
	"encoding/json"
	"fmt"
)

func (c *Controller) EInfo() (string, error) {

	c.w.Header().Set("Access-Control-Allow-Origin", "*")

	c.r.ParseForm()
	token := c.r.FormValue("token")
	if !utils.CheckInputData(token, "string") {
		return "", errors.New("incorrect token")
	}
	fmt.Println("token",token)
	tokenMap, err := c.OneRow(`SELECT * FROM e_tokens WHERE token = ?`, token).String()
	if err!=nil {
		return "", utils.ErrInfo(err)
	}
	fmt.Println("tokenMap",tokenMap)
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
	m := EInfoResult{
		Token: tokenMap,
		Wallets: wallets,
		Orders: orders,
		Withdraw: withdraw,
	}
	jsonData, err := json.Marshal(m)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	fmt.Println("jsonData", jsonData)
	return string(jsonData), nil
}

type EInfoResult struct {
	Token map[string]string
	Wallets []map[string]string
	Orders []map[string]string
	Withdraw []map[string]string
}