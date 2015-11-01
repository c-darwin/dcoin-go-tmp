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

	// сколько всего продается DC
	eOrders, err := c.GetAll(`SELECT sell_currency_id, sum(amount) as amount FROM e_orders GROUP BY sell_currency_id WHERE sell_currency_id < 1000`, 100)
	if err!=nil {
		return "", utils.ErrInfo(err)
	}
	values := ""
	for _, data := range eOrders {
		values = eOrders["amount"]+` `+c.CurrencyList[eOrders["sell_currency_id"]]+`, `
	}
	if len(values) > 0 {
		values = values[:len(values)-2]
	}
	ps, err := c.Single(`SELECT value FROM e_config WHERE name = 'ps'`).String()
	if err!=nil {
		return "", utils.ErrInfo(err)
	}

	jsonData, err := json.Marshal(map[string]string{"values": values, "ps": ps})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return string(jsonData), nil

}
