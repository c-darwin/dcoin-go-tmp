package controllers
import (
	"html/template"
	"bytes"
	"github.com/c-darwin/dcoin-go-tmp/packages/static"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"strings"
	"time"
	"log"
	"encoding/json"
)

type currencyExchangePage struct {
	SignData string
	ShowSignData bool
	BuyCurrencyName string
	SellCurrencyName string
	TxType string
	TxTypeId int64
	TimeNow int64
	UserId int64
	Alert string
	Lang map[string]string
	CurrencyList map[int64]string
	CurrencyListFullName map[int64]string
	CurrencyListName map[int64]string
	SellOrders []map[string]string
	BuyOrders []map[string]string
	MyOrders []map[string]string
	ConfigCommission map[int64][]float64
	BuyCurrencyId int64
	SellCurrencyId int64
	WalletsAmounts map[int64]float64
	CountSignArr []int
}

func (c *Controller) CurrencyExchange() (string, error) {

	txType := "NewForexOrder";
	txTypeId := utils.TypeInt(txType)
	timeNow := time.Now().Unix()

	c.r.ParseForm()
	parameters := make(map[string]string)
	err := json.Unmarshal([]byte(c.r.PostFormValue("parameters")), &parameters)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	log.Println("parameters", parameters)
	addSql := ""
	if len(parameters["all_currencies"]) == 0 {
		// по умолчанию выдаем только те валюты, которые есть хоть у кого-то на кошельках
		actualCurrencies, err := c.GetList("SELECT currency_id FROM wallets GROUP BY currency_id").String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if len(actualCurrencies) > 0 {
			addSql = " WHERE id IN ("+strings.Join(actualCurrencies, ",")+")";
		}
	}
	currencyListName := make(map[int64]string)
	currency, err := c.GetMap("SELECT id, name FROM currency "+addSql+" ORDER BY name", "id", "name")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	for k, v := range currency {
		currencyListName[utils.StrToInt64(k)] = v
	}

	var sellCurrencyId, buyCurrencyId int64
	if len(parameters["buy_currency_id"]) > 0 {
		buyCurrencyId = utils.StrToInt64(parameters["buy_currency_id"])
		c.sess.Set("buy_currency_id", buyCurrencyId)
	}
	if len(parameters["sell_currency_id"]) > 0 {
		sellCurrencyId = utils.StrToInt64(parameters["sell_currency_id"])
		c.sess.Set("sell_currency_id", sellCurrencyId)
	}
	if buyCurrencyId == 0 {
		buyCurrencyId = GetSessInt64("buy_currency_id", c.sess)
	}
	if sellCurrencyId == 0 {
		sellCurrencyId = GetSessInt64("sell_currency_id", c.sess)
	}
	if buyCurrencyId == 0 {
		buyCurrencyId = 1
	}
	if sellCurrencyId == 0 {
		sellCurrencyId = 72
	}

	buyCurrencyName := currencyListName[buyCurrencyId]
	sellCurrencyName := currencyListName[sellCurrencyId]


	// валюты
	currencyListFullName, err := c.GetCurrencyListFullName()
	if err != nil {
		return "", utils.ErrInfo(err)
	}




	buyOrders, err := c.GetAll(`
			SELECT *
			FROM forex_orders
			WHERE  buy_currency_id =  ? AND
						 sell_currency_id = ? AND
						 empty_block_id = 0 AND
						 del_block_id = 0
						 `, 100, buyCurrencyId, sellCurrencyId)

	sellOrders, err := c.GetAll(`
			SELECT *
			FROM forex_orders
			WHERE  buy_currency_id =  ? AND
						 sell_currency_id = ? AND
						 empty_block_id = 0 AND
						 del_block_id = 0
						 `, 100, sellCurrencyId, buyCurrencyId)


	myOrders, err := c.GetAll(`
			SELECT *
			FROM forex_orders
			WHERE user_id =  ? AND
						 empty_block_id = 0 AND
						 del_block_id = 0
						 `, 100, c.SessRestricted)

	rows, err := c.Query(c.FormatQuery(`
			SELECT amount, currency_id, last_update
			FROM wallets
			WHERE user_id = ? AND
						currency_id IN (?, ?)
			`), c.SessUserId, sellCurrencyId, buyCurrencyId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	defer rows.Close()
	walletsAmounts := make(map[int64]float64)
	for rows.Next() {
		var amount float64
		var currency_id, last_update int64
		err = rows.Scan(&amount, &currency_id, &last_update)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		profit, err := c.CalcProfitGen(currency_id, amount, c.SessUserId, last_update, timeNow, "wallets")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		amount+=profit
		amount = utils.Round(amount, 2)
		forex_orders_amount, err := c.Single("SELECT sum(amount) FROM forex_orders WHERE user_id = ? AND sell_currency_id = ? AND del_block_id = 0", c.SessUserId, currency_id).Float64()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		amount-=forex_orders_amount
		walletsAmounts[currency_id] = amount
	}


	data, err := static.Asset("static/templates/currency_exchange.html")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	alert_success, err := static.Asset("static/templates/alert_success.html")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	signatures, err := static.Asset("static/templates/signatures.html")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	funcMap := template.FuncMap{
		"mult": func(a, b string) float64 {
			return utils.StrToFloat64(a)*utils.StrToFloat64(b)
		},
		"div": func(a, b interface{}) float64 {
			return utils.InterfaceToFloat64(a)/utils.InterfaceToFloat64(b)
		},
		"round": func(a float64, num int) float64 {
			return utils.Round(a, num)
		},
	}
	t := template.Must(template.New("template").Funcs(funcMap).Parse(string(data)))
	t = template.Must(t.Parse(string(alert_success)))
	t = template.Must(t.Parse(string(signatures)))
	b := new(bytes.Buffer)
	err = t.ExecuteTemplate(b, "currencyExchange", &currencyExchangePage{
		Lang: c.Lang,
		CountSignArr: c.CountSignArr,
		ShowSignData: c.ShowSignData,
		WalletsAmounts: walletsAmounts,
		CurrencyListName: currencyListName,
		BuyCurrencyId: buyCurrencyId,
		SellCurrencyId: sellCurrencyId,
		BuyCurrencyName: buyCurrencyName,
		SellCurrencyName: sellCurrencyName,
		CurrencyListFullName : currencyListFullName,
		ConfigCommission: c.ConfigCommission,
		TimeNow: timeNow,
		SellOrders: sellOrders,
		BuyOrders: buyOrders,
		MyOrders: myOrders,
		UserId: c.SessUserId,
		TxType: txType,
		TxTypeId: txTypeId,
		SignData: ""})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return b.String(), nil
}
