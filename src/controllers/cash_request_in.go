package controllers
import (
	"utils"
)

type CashRequestInPage struct {
	Alert string
	SignData string
	ShowSignData bool
	TxType string
	TxTypeId int64
	TimeNow int64
	UserId int64
	Lang map[string]string
	CountSignArr []int
	CurrencyList map[int64]string
	CashRequestsStatus map[string]string
	MyCashRequests []map[string]string
	ActualData map[string]string
}

func (c *Controller) CashRequestIn() (string, error) {

	txType := "CashRequestIn";
	txTypeId := utils.TypeInt(txType)
	timeNow := utils.Time()

	cashRequestsStatus := map[string]string{"my_pending": c.Lang["local_pending"], "pending": c.Lang["pending"], "approved": c.Lang["approved"], "rejected": c.Lang["rejected"]}

	// Узнаем свой user_id
	userId, err:=c.GetMyUserId(c.MyPrefix)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// актуальный запрос к нам на получение налички. Может быть только 1.
	actualData, err := c.OneRow(`
		SELECT `+c.MyPrefix+`my_cash_requests.cash_request_id,
					 `+c.MyPrefix+`my_cash_requests.id,
					 `+c.MyPrefix+`my_cash_requests.comment_status,
					 `+c.MyPrefix+`my_cash_requests.comment,
					 cash_requests.amount,
					 cash_requests.currency_id,
					 cash_requests.from_user_id,
					 LOWER(HEX(cash_requests.hash_code)) as hash_code
		FROM `+c.MyPrefix+`my_cash_requests
		LEFT JOIN cash_requests ON cash_requests.id = `+c.MyPrefix+`my_cash_requests.cash_request_id
		WHERE cash_requests.to_user_id = ? AND
					 cash_requests.status = 'pending' AND
					 cash_requests.time > ? AND
					 cash_requests.del_block_id = 0 AND
					 cash_requests.for_repaid_del_block_id = 0
		ORDER BY cash_request_id DESC
		LIMIT 1`, userId, utils.Time() - c.Variables.Int64["cash_request_time"]).String()

	// список ранее отправленных ответов на запросы.
	myCashRequests, err := c.GetAll("SELECT * FROM "+c.MyPrefix+"my_cash_requests WHERE to_user_id = ?", -1, userId)

	TemplateStr, err := makeTemplate("cash_request_in", "cashRequestIn", &CashRequestInPage {
		Alert: c.Alert,
		Lang: c.Lang,
		CountSignArr: c.CountSignArr,
		ShowSignData: c.ShowSignData,
		UserId: userId,
		TimeNow: timeNow,
		TxType: txType,
		TxTypeId: txTypeId,
		SignData: "",
		CurrencyList : c.CurrencyList,
		CashRequestsStatus: cashRequestsStatus,
		MyCashRequests: myCashRequests,
		ActualData: actualData})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

