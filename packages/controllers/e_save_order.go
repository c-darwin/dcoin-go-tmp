package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"errors"
	"strings"
	"time"
	"sync"
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

func geReductionLock() {
	return utils.DB.ExecSql("SELECT time FROM e_reduction_lock")
}

func userLock(userId int64) error {
	var affect int64
	var err error
	// даем время, чтобы lock освободился от другого запроса
	for i:=0; i<4; i++ {
		affect, err = utils.DB.ExecSqlGetAffect(`UPDATE e_users SET lock = ? WHERE id = ? AND lock = 0`, utils.Time(), userId)
		if affect > 1 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if affect == 0 {
		return errors.New("queue error")
	}
	return nil
}

func userAmountAndProfit(userId, currencyId int64) float64 {
	var UserCurrencyId, UserLastUpdate int64
	var UserAmount float64
	err := utils.DB.QueryRow(utils.DB.FormatQuery("SELECT currency_id, amount, last_update FROM e_wallets WHERE user_id  =  ? AND currency_id  =  ?"), userId, currencyId).Scan(&UserCurrencyId, &UserAmount, &UserLastUpdate)
	if err != nil {
		return 0
	}
	if UserAmount <= 0 {
		return 0
	}
	profit, err := utils.DB.CalcProfitGen(UserCurrencyId, UserAmount, 0, UserLastUpdate, utils.Time(), "wallet")
	return UserAmount + profit
}

var eWallets = &sync.Mutex{}

func (c *Controller) newForexOrder(userId int64, amount, sellRate float64, sellCurrencyId, buyCurrencyId int64, orderType string) error {
	log.Debug("userId: %v / amount: %v / sellRate: %v / float64: %v / sellCurrencyId: %v / buyCurrencyId: %v / orderType: %v ", userId, amount, sellRate, sellCurrencyId, buyCurrencyId, orderType)
	curTime := utils.Time()
	err := userLock(userId)
	if err!=nil {
		return utils.ErrInfo(err)
	}
	commission := amount * (c.ECommission/100)
	newAmount := amount - commission
	if newAmount < commission {
		commission = 0
		newAmount = amount
	}
	log.Debug("newAmount: %v / commission: %v ", newAmount, commission)
	if commission {
		userAmount := userAmountAndProfit(1, sellCurrencyId)
		newAmount_ := userAmount + commission
		// наисляем комиссию системе
		err = updEWallet(1, sellCurrencyId, utils.Time, newAmount_)
		if err != nil {
			return utils.ErrInfo(err)
		}
		// и сразу вычитаем комиссию с кошелька юзера
		userAmount = userAmountAndProfit(userId, sellCurrencyId)
		err = utils.DB.ExecSql("UPDATE e_wallets SET amount = ?, last_update = ? WHERE user_id = ? AND currency_id = ?",userAmount-commission, utils.Time(), userId, sellCurrencyId)
		if err != nil {
			return utils.ErrInfo(err)
		}
	}
	// обратный курс. нужен для поиска по ордерам
	reverseRate := utils.Round(1/sellRate, 6)

	var totalBuyAmount, totalSellAmount int64
	if orderType == "buy" {
		totalBuyAmount = newAmount + reverseRate
	} else {
		totalSellAmount = newAmount
	}

	var debit float64
	var prevUserId int64
	// берем из БД только те ордеры, которые удовлетворяют нашим требованиям
	rows, err := utils.DB.Query(utils.DB.FormatQuery(`
				SELECT id, user_id, amount, sell_rate, buy_currency_id, sell_currency_id
				FROM e_orders
				WHERE buy_currency_id = ? AND
							 sell_rate >= ? AND
							 sell_currency_id = ?  AND
							 del_time = 0 AND
							 empty_time = 0
				ORDER BY sell_rate DESC
				`, sellCurrencyId, reverseRate, buyCurrencyId)
	if err != nil {
		return utils.ErrInfo(err)
	}

	for rows.Next() {
		var rowId, rowUserId, rowBuyCurrencyId, rowSellCurrencyId int64
		var rowAmount, rowSellRate float64
		err = rows.Scan(&rowId, &rowUserId, &rowAmount, &rowSellRate, &rowBuyCurrencyId, &rowSellCurrencyId)
		if err != nil {
			return utils.ErrInfo(err)
		}
		log.Debug("rowId: %v / rowUserId: %v / rowAmount: %v / rowSellRate: %v / rowBuyCurrencyId: %v / rowSellCurrencyId: %v", rowId, rowUserId, rowAmount, rowSellRate, rowBuyCurrencyId, rowSellCurrencyId)

		// чтобы ордеры одного и тоже юзера не вызывали стопор, пропускаем все его оредера
		if rowUserId == prevUserId {
			continue
		}
		// блокируем юзера, чей ордер взяли, кроме самого себя
		if rowUserId != userId {
			lockErr := userLock(rowUserId)
			if lockErr!=nil {
				log.Error("%v", utils.ErrInfo(err))
				prevUserId = rowUserId
				continue
			}
		}
		if orderType == "buy" {
			// удовлетворит ли данный ордер наш запрос целиком
			if rowAmount>= totalBuyAmount {
				debit = totalBuyAmount
				log.Debug("order ENDED")
			} else {
				debit = rowAmount
			}
		} else {
			// удовлетворит ли данный ордер наш запрос целиком
			if rowAmount/rowSellRate >= totalSellAmount {
				debit = totalSellAmount * rowSellRate
			} else {
				debit = rowAmount
			}
		}
		log.Debug("totalBuyAmount: %v / debit: %v", totalBuyAmount, debit)
		if rowAmount - debit < 0.01 { // ордер опустошили
			err = utils.DB.ExecSql("UPDATE e_orders SET amount = 0, empty_time = ? WHERE id = ?", curTime, rowId)
			if err != nil {
				return utils.ErrInfo(err)
			}
		} else {
			// вычитаем забранную сумму из ордера
			err = utils.DB.ExecSql("UPDATE e_orders SET amount = amount - ? WHERE id = ?", debit, rowId)
			if err != nil {
				return utils.ErrInfo(err)
			}
		}
		mySellRate := utils.ClearNull(1/rowSellRate, 6)
		myAmount := debit * mySellRate
		eTradeSellCurrencyId := sellCurrencyId
		eTradeBuyCurrencyId := buyCurrencyId

		// для истории сделок
		err = utils.DB.ExecSql("INSERT INTO e_trade ( user_id, sell_currency_id, sell_rate, amount, buy_currency_id, time, main ) VALUES ( ?, ?, ?, ?, ?, ?, 1 )", userId, eTradeSellCurrencyId, mySellRate, myAmount, eTradeBuyCurrencyId, curTime)
		if err != nil {
			return utils.ErrInfo(err)
		}

		// тот, чей ордер обрабатываем
		err = utils.DB.ExecSql("INSERT INTO e_trade ( user_id, sell_currency_id, sell_rate, amount, buy_currency_id, time ) VALUES ( ?, ?, ?, ?, ?, ? )", rowUserId, eTradeBuyCurrencyId, rowSellRate, debit, eTradeSellCurrencyId, curTime)
		if err != nil {
			return utils.ErrInfo(err)
		}

		// ==== Продавец валюты (тот, чей ордер обработали) ====
		// сколько продавец данного ордера продал валюты
		sellerSellAmount := debit

		// сколько продавец получил buy_currency_id с продажи суммы $seller_sell_amount по его курсу
		sellerBuyAmount := sellerSellAmount * (1/rowSellRate)

		// начисляем валюту, которую продавец получил (R)
		userAmount := userAmountAndProfit(rowUserId, rowBuyCurrencyId)
		newAmount_ := userAmount + sellerBuyAmount
		err = updEWallet(rowUserId, rowBuyCurrencyId, utils.Time(), newAmount_)
		if err != nil {
			return utils.ErrInfo(err)
		}

		// ====== Покупатель валюты (наш юзер) ======

		// списываем валюту, которую мы продали (R)
		userAmount = userAmountAndProfit(userId, rowBuyCurrencyId)
		newAmount_ = userAmount - sellerBuyAmount
		err = updEWallet(userId, rowBuyCurrencyId, utils.Time(), newAmount_)
		if err != nil {
			return utils.ErrInfo(err)
		}

		// начисляем валюту, которую мы получили (U)
		userAmount = userAmountAndProfit(userId, rowSellCurrencyId)
		newAmount_ = userAmount + sellerSellAmount
		err = updEWallet(userId, rowSellCurrencyId, utils.Time(), newAmount_)
		if err != nil {
			return utils.ErrInfo(err)
		}

		if orderType == "buy" {
			totalBuyAmount-=rowAmount
			if totalBuyAmount <= 0 {
				userUnlock(rowUserId)
				break; // проход по ордерам прекращаем, т.к. наш запрос удовлетворен
			}
		} else {
			totalSellAmount-=rowAmount/rowSellRate
			if totalSellAmount <= 0 {
				userUnlock(rowUserId)
				break; // проход по ордерам прекращаем, т.к. наш запрос удовлетворен
			}
		}
		if rowUserId != userId {
			userUnlock(rowUserId)
		}
	}

	log.Debug("totalBuyAmount: %v / orderType: %v / sellRate: %v / reverseRate: %v", totalBuyAmount, orderType, sellRate, reverseRate)

	// если после прохода по всем имеющимся ордерам мы не набрали нужную сумму, то создаем свой ордер
	if totalBuyAmount > 0 || totalSellAmount > 0 {
		var newOrderAmount float64
		if orderType == "buy" {
			newOrderAmount = totalBuyAmount * sellRate
		} else {
			newOrderAmount = totalSellAmount
		}
		log.Debug("newOrderAmount: %v", newOrderAmount)
		if newOrderAmount >= 0.000001 {
			err = utils.DB.ExecSql("INSERT INTO e_orders ( time, user_id, sell_currency_id, sell_rate, begin_amount, amount, buy_currency_id ) VALUES ( ?, ?, ?, ?, ?, ?, ? )", curTime, userId, sellCurrencyId, sellRate, newOrderAmount, newOrderAmount, buyCurrencyId)
			if err != nil {
				return utils.ErrInfo(err)
			}

			// вычитаем с баланса сумму созданного ордера
			userAmount := userAmountAndProfit(userId, sellCurrencyId)
			err = utils.DB.ExecSql("UPDATE e_wallets SET amount = ?, last_update = ? WHERE user_id = ? AND currency_id = ?", userAmount-newOrderAmount, utils.Time(), userId, sellCurrencyId)
			if err != nil {
				return utils.ErrInfo(err)
			}
		}
	}

	userUnlock(userId)

	return nil
}

func updEWallet(userId, currencyId, lastUpdate int64, amount float64) error {
	eWallets.Lock()
	exists, err := utils.DB.Single(`SELECT user_id FROM e_wallets WHERE user_id = ?`, userId).Int64()
	if err!=nil {
		eWallets.Unlock()
		return utils.ErrInfo(err)
	}
	if exists == 0 {
		err = utils.DB.ExecSql("INSERT INTO e_wallets ( user_id, currency_id, amount, last_update ) VALUES ( ?, ?, ?, ? )", userId, currencyId, amount, lastUpdate)
		if err != nil {
			eWallets.Unlock()
			return utils.ErrInfo(err)
		}
	} else {
		err = utils.DB.ExecSql("UPDATE e_wallets SET amount = ?, last_update = ? WHERE user_id = ?", amount, lastUpdate, userId)
		if err != nil {
			eWallets.Unlock()
			return utils.ErrInfo(err)
		}
	}
	eWallets.Unlock()
	return nil
}

func userUnlock(userId int64) {
	utils.DB.ExecSql("UPDATE e_users SET lock = 0 WHERE id = ?", userId)
}