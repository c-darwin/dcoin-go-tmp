package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"time"
)

type eMainPage struct {
	AlertMessages []map[string]string
	Lang          map[string]string
	CurrencyList  map[int64]string
	Commission    string
	Members       int64
	SellMax       float64
	BuyMin        float64
	Orders        eOrders
	DcCurrency    string
	Currency      string
	DcCurrencyId  int64
	CurrencyId    int64
}

func eGetCurrencyList() (map[int64]string, error) {
	rez := make(map[int64]string)
	list, err := utils.DB.GetMap("SELECT id, name FROM e_currency ORDER BY name", "id", "name")
	if err != nil {
		return rez, utils.ErrInfo(err)
	}
	for id, name := range list {
		rez[utils.StrToInt64(id)] = name
	}
	return rez, nil
}

func (c *Controller) EMain() (string, error) {

	var err error

	dcCurrencyId := utils.StrToInt64(c.Parameters["dc_currency_id"])
	currencyId := utils.StrToInt64(c.Parameters["currency_id"])
	if dcCurrencyId == 0 {
		dcCurrencyId = 72
	}
	if currencyId == 0 {
		currencyId = 1001
	}

	// все валюты, с которыми работаем
	currencyList, err := eGetCurrencyList()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	if len(currencyList[dcCurrencyId]) == 0 || len(currencyList[currencyId]) == 0 {
		return "", utils.ErrInfo("incorrect currency")
	}

	// пары валют для меню
	dcCurrency := currencyList[dcCurrencyId]
	currency := currencyList[currencyId]

	// история сделок
	var tradeHistory []map[string]string

	// откатываем наши блоки до начала вилки
	rows, err := c.Query(c.FormatQuery(`
			SELECT sell_currency_id, sell_rate, amount, time
			FROM e_trade
			WHERE ((sell_currency_id = ? AND buy_currency_id = ?) OR (sell_currency_id = ? AND buy_currency_id = ?)) AND main = 1
			ORDER BY time DESC
			LIMIT 40
			`), dcCurrencyId, currencyId, currencyId, dcCurrencyId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	for rows.Next() {
		var sellCurrencyId, eTime int64
		var sellRate, amount float64
		err = rows.Scan(&sellCurrencyId, sellRate, amount, eTime)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		var eType string
		var eAmount float64
		var eTotal float64
		if sellCurrencyId == dcCurrencyId {
			eType = "sell"
			sellRate = 1 / sellRate
			eAmount = amount
			eTotal = amount * sellRate
		} else {
			eType = "buy"
			eAmount = amount * (1 / sellRate)
			eTotal = amount
		}
		t := time.Unix(eTime, 0)
		tradeHistory = append(tradeHistory, map[string]string{"Time": t.Format(c.TimeFormat), "Type": eType, "SellRate": utils.Float64ToStr(sellRate), "Amount": utils.Float64ToStr(eAmount), "Total": utils.Float64ToStr(eTotal)})
	}

	// активные ордеры на продажу
	var orders eOrders
	rows, err = c.Query(c.FormatQuery(`
			SELECT sell_rate, amount
			FROM e_orders
			WHERE (sell_currency_id = ? AND buy_currency_id = ? AND
						empty_time = 0 AND
						del_time = 0 AND
						amount > 0
			ORDER BY sell_rate DESC
			LIMIT 100
			`), dcCurrencyId, currencyId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	// мин. цена покупки
	var buyMin float64
	for rows.Next() {
		var sellRate, amount float64
		err = rows.Scan(&sellRate, amount)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		orders.Sell = map[float64]float64{sellRate: orders.Sell[sellRate] + amount}
		if buyMin == 0 {
			buyMin = sellRate
		} else if sellRate < buyMin {
			buyMin = sellRate
		}
	}

	// активные ордеры на покупку
	rows, err = c.Query(c.FormatQuery(`
			SELECT *
			FROM e_orders
			WHERE (sell_currency_id = ? AND buy_currency_id = ?) AND
					empty_time = 0 AND
					del_time = 0 AND
					amount > 0
			ORDER BY sell_rate DESC
			LIMIT 100
			`), currencyId, dcCurrencyId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	// мин. цена продажи
	var sellMax float64
	for rows.Next() {
		var sellRate, amount float64
		err = rows.Scan(&sellRate, amount)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		orders.Buy = map[float64]float64{sellRate: orders.Buy[sellRate] + amount*(1/sellRate)}
		if sellMax == 0 {
			sellMax = sellRate
		} else if sellRate < sellMax {
			sellMax = sellRate
		}
	}

	// комиссия
	commission := c.EConfig["commission"]

	// кол-во юзеров
	members, err := c.Single(`SELECT count(*) FROM e_users`).Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	TemplateStr, err := makeTemplate("emain", "eMain", &eMainPage{
		Lang:         c.Lang,
		Commission:   commission,
		Members:      members,
		SellMax:      sellMax,
		BuyMin:       buyMin,
		Orders:       orders,
		DcCurrency:   dcCurrency,
		Currency:     currency,
		DcCurrencyId: dcCurrencyId,
		CurrencyId:   currencyId,
		CurrencyList: currencyList})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

type eOrders struct {
	Sell map[float64]float64
	Buy  map[float64]float64
}
