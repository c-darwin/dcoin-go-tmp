package controllers
import (
    "encoding/json"
	"utils"
	"log"
	"time"
)

func (c *Controller) GetSellerData() (string, error) {

	c.r.ParseForm()

	getUserId := utils.StrToInt64(c.r.FormValue("user_id"))
	if !utils.CheckInputData(getUserId, "int") {
		return `{"result":"incorrect userId"}`, nil
	}
	currencyId := utils.StrToInt64(c.r.FormValue("currency_id"))
	if !utils.CheckInputData(currencyId, "currency_id") {
		return `{"result":"incorrect currency_id"}`, nil
	}

	arbitrationTrustList, err := c.GetList("SELECT arbitrator_user_id FROM arbitration_trust_list WHERE user_id  =  ?", getUserId).Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	/*
	 * Статистика по продавцу
	 * */
	// оборот всего
	sellerTurnover, err := c.Single("SELECT sum(amount) FROM orders WHERE seller  =  ? AND currency_id  =  ?", getUserId, currencyId).Float64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// оборот за месяц
	sellerTurnoverM, err := c.Single("SELECT sum(amount) FROM orders WHERE seller  =  ? AND time > ? AND currency_id  =  ?", getUserId, time.Now().Unix()-3600*24*30, currencyId).Float64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// Кол-во покупателей за последний месяц
	buyersCountM, err := c.Single("SELECT count(id) FROM ( SELECT id FROM orders WHERE seller  =  ? AND time > ? AND currency_id  =  ? GROUP BY buyer ) as t1", getUserId, currencyId, time.Now().Unix()-3600*24*30).Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// Кол-во покупателей-майнеров за последний месяц
	buyersMinersCountM, err := c.Single(`
			SELECT count(id)
			FROM (
				SELECT orders.id
				FROM orders
				LEFT JOIN miners_data ON miners_data.user_id =  orders.buyer
				WHERE seller = ? AND
							 orders.time > ? AND
							 orders.currency_id =  ? AND
							 miner_id > 0
				GROUP BY buyer
			) as t1
		`, getUserId, time.Now().Unix()-3600*24*30, currencyId).Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// Кол-во покупателей всего
	buyersCount, err := c.Single("SELECT count(id) FROM ( SELECT id FROM orders WHERE seller  =  ? AND currency_id  =  ? GROUP BY buyer ) as t1", getUserId, currencyId).Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// Кол-во покупателей-майнеров всего
	buyersMinersCount, err := c.Single(`
		SELECT count(id)
		FROM (
			SELECT orders.id
			FROM orders
			LEFT JOIN miners_data ON miners_data.user_id =  orders.buyer
			WHERE seller = ? AND
						 orders.currency_id = ? AND
						 miner_id > 0
			GROUP BY buyer
		) as t1
	`, getUserId, currencyId).Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// Заморожено для манибека
	holdAmount, err := c.Single(`
		SELECT sum(hold_back_amount)
		FROM orders
		LEFT JOIN miners_data ON miners_data.user_id =  orders.buyer
		WHERE seller = ? AND
					 orders.currency_id = ? AND
					 miner_id > 0
		GROUP BY buyer
	`, getUserId, currencyId).Float64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// Холдбек % на 30 дней
	sellerData, err := c.OneRow("SELECT seller_hold_back_pct, arbitration_days_refund FROM users WHERE user_id  =  ?", getUserId).String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	sellerHoldBackPct := utils.StrToFloat64(sellerData["seller_hold_back_pct"])
	arbitrationDaysRefund := utils.StrToInt64(sellerData["arbitrationDaysRefund"])
	buyersCount = buyersCount - buyersMinersCount;
	buyersCountM = buyersCountM - buyersMinersCountM;

	rez := selleData{Trust_list: arbitrationTrustList, Seller_hold_back_pct: sellerHoldBackPct,  Arbitration_days_refund: arbitrationDaysRefund, Buyers_miners_count_m: buyersMinersCountM, Buyers_miners_count: buyersMinersCount, Buyers_count: buyersCount ,Buyers_count_m: buyersCountM ,Seller_turnover_m: sellerTurnoverM, Seller_turnover: sellerTurnover, Hold_amount: holdAmount}
	log.Println(rez)
	log.Println(arbitrationTrustList)
	result, err := json.Marshal(rez)
	if err != nil {
		log.Println(err)
		return "", utils.ErrInfo(err)
	}
	log.Println(string(result))

	return string(result), nil
}

type selleData struct {
	Trust_list []int64 `json:"trust_list"`
	Seller_hold_back_pct float64 `json:"seller_hold_back_pct"`
	Arbitration_days_refund int64 `json:"arbitration_days_refund"`
	Buyers_miners_count_m int64 `json:"buyers_miners_count_m"`
	Buyers_miners_count int64 `json:"buyers_miners_count"`
	Buyers_count int64 `json:"buyers_count"`
	Buyers_count_m int64 `json:"buyers_count_m"`
	Seller_turnover_m float64 `json:"seller_turnover_m"`
	Seller_turnover float64 `json:"seller_turnover"`
	Hold_amount float64 `json:"hold_amount"`
}
