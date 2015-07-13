package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"time"
	"strings"
)

type newPromisedAmountPage struct {
	Alert string
	SignData string
	ShowSignData bool
	TxType string
	TxTypeId int64
	TimeNow int64
	UserId int64
	Lang map[string]string
	CountSignArr []int
	LastTxFormatted string
	ConfigCommission map[int64][]float64
	Navigate string
	CurrencyId int64
	CurrencyList map[int64]map[string]string
	CurrencyListName map[int64]string
	MaxPromisedAmounts map[string]string
	LimitsText string
	PaymentSystems map[string]string
	CountPs []int
}

func (c *Controller) NewPromisedAmount() (string, error) {

	txType := "NewPromisedAmount";
	txTypeId := utils.TypeInt(txType)
	timeNow := time.Now().Unix()

	navigate := "promisedAmountList"
	if len(c.Navigate) > 0 {
		navigate = c.Navigate
	}

	rows, err := c.Query(c.FormatQuery(`
		SELECT id,
					 name,
					 full_name,
					 max_other_currencies
		FROM currency
		ORDER BY full_name`))
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	currencyList := make(map[int64]map[string]string)
	currencyListName := make(map[int64]string)
	defer rows.Close()
	for rows.Next() {
		var id  int64
		var name, full_name, max_other_currencies string
		err = rows.Scan(&id, &name, &full_name, &max_other_currencies)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		currencyList[id] = map[string]string{"id":utils.Int64ToStr(id), "name":name, "full_name":full_name, "max_other_currencies":max_other_currencies}
		currencyListName[id] = name
	}

	paymentSystems, err := c.GetPaymentSystems()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	maxPromisedAmounts, err := c.GetMap(`SELECT currency_id, amount FROM max_promised_amounts WHERE block_id = 1`, "currency_id", "amount")
	maxPromisedAmountsMaxBlock, err := c.GetMap(`SELECT currency_id, amount FROM max_promised_amounts WHERE block_id = (SELECT max(block_id) FROM max_promised_amounts ) OR block_id = 0`, "currency_id", "amount")
	for k, v := range maxPromisedAmountsMaxBlock {
		maxPromisedAmounts[k] = v
	}

	// валюта, которая выбрана в селект-боксе
	currencyId := int64(72)

	limitsText := strings.Replace(c.Lang["limits_text"], "[limit]", utils.Int64ToStr(c.Variables.Int64["limit_promised_amount"]), -1)
	limitsText = strings.Replace(limitsText, "[period]", c.Periods[c.Variables.Int64["limit_promised_amount_period"]], -1)

	countPs := []int{1,2,3,4,5}

	TemplateStr, err := makeTemplate("new_promised_amount", "newPromisedAmount", &newPromisedAmountPage{
		Alert: c.Alert,
		Lang: c.Lang,
		CountSignArr: c.CountSignArr,
		ShowSignData: c.ShowSignData,
		UserId: c.SessUserId,
		TimeNow: timeNow,
		TxType: txType,
		TxTypeId: txTypeId,
		SignData: "",
		ConfigCommission: c.ConfigCommission,
		Navigate: navigate,
		CurrencyId: currencyId,
		CurrencyList: currencyList,
		CurrencyListName: currencyListName,
		MaxPromisedAmounts: maxPromisedAmounts,
		LimitsText: limitsText,
		PaymentSystems: paymentSystems,
		CountPs: countPs})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
