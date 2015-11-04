package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

type eMyOrdersPage struct {
	Lang          map[string]string
	CurrencyList  map[int64]string
	UserId	int64
	MyOrders []*EmyOrders
}

func (c *Controller) EMyOrders() (string, error) {

	var err error

	if c.SessUserId == 0 {
		return `<script language="javascript"> window.location.href = "/e/"</script>If you are not redirected automatically, follow the <a href="/e/">/e/</a>`, nil
	}

	currencyList, err := utils.EGetCurrencyList()

	var myOrders []*EmyOrders

	rows, err := c.Query(c.FormatQuery(`
			SELECT time, empty_time, amount, begin_amount, sell_currency_id, buy_currency_id, sell_rate, begin_amount
			FROM e_orders
			WHERE user_id = ? AND
						 del_time = 0
			ORDER BY time DESC
			LIMIT 40
			`), c.SessUserId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	defer rows.Close()
	for rows.Next() {
		//var time, emptyTime, sellCurrencyId, buyCurrencyId int64
		//var amount, beginAmount, sellRate, total float64
		//var status, orderType, pair string
		myOrder := new(EmyOrders)
		//err = rows.Scan(&emyOrder.time, &emptyTime, &amount, &beginAmount, &sellCurrencyId, &buyCurrencyId, &sellRate, &beginAmount)
		err = rows.Scan(&myOrder.Time, &myOrder.EmptyTime, &myOrder.Amount, &myOrder.BeginAmount, &myOrder.SellCurrencyId, &myOrder.BuyCurrencyId, &myOrder.SellRate, &myOrder.BeginAmount)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if myOrder.EmptyTime == 0 {
			myOrder.Status = c.Lang["active"]
		} else {
			myOrder.Status = c.Lang["executed"]
		}

		// на сколько % выполнен ордер
		myOrder.OrderComplete = utils.Round(100-(myOrder.Amount/myOrder.BeginAmount)*100, 1)

		// определим тип ордера и пару
		if myOrder.SellCurrencyId < 1000 {
			myOrder.OrderType = "sell"
			myOrder.SellRate = 1/myOrder.SellRate
			myOrder.Total = myOrder.BeginAmount*myOrder.SellRate
			myOrder.BeginAmount = myOrder.BeginAmount
			myOrder.Pair = currencyList[myOrder.SellCurrencyId]+"/"+currencyList[myOrder.BuyCurrencyId]
		} else {
			myOrder.OrderType = "buy"
			myOrder.Total = myOrder.BeginAmount
			myOrder.Amount = myOrder.BeginAmount * (1/myOrder.SellRate)
			myOrder.Pair = currencyList[myOrder.BuyCurrencyId]+"/"+currencyList[myOrder.SellCurrencyId]
		}

		myOrders = append(myOrders, myOrder)
	}

	TemplateStr, err := makeTemplate("e_my_orders", "eMyOrders", &eMyOrdersPage{
		Lang:         c.Lang,
		UserId: c.SessUserId,
		MyOrders : myOrders,
		CurrencyList: currencyList})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

type EmyOrders struct {
	Time, EmptyTime, SellCurrencyId, BuyCurrencyId int64
	Amount, BeginAmount, SellRate, Total, OrderComplete float64
	Status, OrderType, Pair string
}
