package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"strings"
	"sort"
)

type eMyFinancePage struct {
	Lang          map[string]string
	CurrencyList  map[int64]string
	UserId	int64
	MyFinanceHistory []*EmyFinanceType
	Collapse string
	Currency map[string]map[string]string
}

type EmyFinanceType struct {
	Ftype, Status, AddType, Method string
	Amount, WdAmount float64
	CurrencyId, AddTime, CloseTime, OpenTime int64
}

func (c *Controller) EMyFinance() (string, error) {

	var err error

	if c.SessUserId == 0 {
		return `<script language="javascript"> window.location.href = "/e/"</script>If you are not redirected automatically, follow the <a href="/e/">/e/</a>`, nil
	}

	confirmations := c.EConfig["confirmations"]

	currencyList, err := utils.EGetCurrencyList()

	// счет, куда юзеры должны слать DC
	mainDcAccount := c.EConfig["main_dc_account"]

	currency := make(map[string]map[string]string)

	// валюты, по которым идут торги на бирже
	//var myWallets []map[string]string
	eCurrency, err := c.GetAll(`SELECT name, id FROM e_currency ORDER BY sort_id ASC`, -1)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	for _, data := range eCurrency {
		wallet, err := c.OneRow("SELECT * FROM e_wallets WHERE user_id  =  ? AND currency_id  =  ?", c.SessUserId, data["id"]).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if len(wallet) > 0 {
			amount := utils.StrToFloat64(wallet["amount"])
			profit, err := utils.DB.CalcProfitGen(utils.StrToInt64(wallet["currency_id"]), amount, 0, utils.StrToInt64(wallet["last_update"]), utils.Time(), "wallet")
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			wallet["amount"] = utils.Float64ToStr(amount + profit)
		} else {
			wallet["amount"] = "0"
		}

		currency[data["id"]] = make(map[string]string)
		currency[data["id"]]["amount"] = wallet["amount"]
		currency[data["id"]]["name"] = data["name"]
		if utils.StrToInt64(data["id"]) < 1000 { //DC
			currency[data["id"]]["input"] = strings.Replace(c.Lang["dc_deposit_text"], "[dc_currency]", data["name"],  -1)
			currency[data["id"]]["input"] = strings.Replace(currency[data["id"]]["input"], "[account]", mainDcAccount, -1)
			currency[data["id"]]["input"] = strings.Replace(currency[data["id"]]["input"], "[user_id]",utils.Int64ToStr(c.SessUserId),  -1)
			currency[data["id"]]["input"] = strings.Replace(currency[data["id"]]["input"], "[confirmations]",  confirmations,-1)
		}

		currency[data["id"]]["output"] = `<div class="pull-left"><h4>`+c.Lang["withdraw0"]+` `+data["name"]+`</h4>
			<table class="table_out">
			<tbody>
			<tr>
			<td>`+c.Lang["your_dcoin_account"]+`:</td>
			<td class="form-inline"><input id="account-`+data["id"]+`" class="form-control col-xs-3" type="text"></td>
			</tr>
			<tr>
			<td>`+c.Lang["amount_to_withdrawal"]+`:</td>
			<td class="form-inline" style="line-height: 35px"><input id="amount-`+data["id"]+`" class="form-control col-xs-3" maxlength="15" type="text"  onkeyup="calc_withdraw_amount(`+data["id"]+`, '0.1')" onchange="calc_withdraw_amount(`+data["id"]+`, '0.1')" style="margin-right:5px"> `+data["name"]+`</td>
			</tr>
			<tr>
			<td>`+c.Lang["you_will_receive"]+`:</td>
			<td class="form-inline" style="line-height: 35px"><input  disabled="" id="withdraw_amount-`+data["id"]+`" class="form-control col-xs-3" maxlength="15" type="text" style="margin-right:5px"> `+data["name"]+`</td>
			</tr>
			</tbody></table><div id="alerts-`+data["id"]+`"></div><button class="btn btn-outline btn-primary" onclick="withdraw(`+data["id"]+`, 'Dcoin')">`+c.Lang["withdrawal"]+`</button>
			</div><div class="pull-left" style="margin-left:30px; margin-top:43px; border-left: 4px solid #ccc; padding:7px 7px; width:400px">`
		dcWithdrawText := strings.Replace(c.Lang["dc_withdraw_text"], "[min_amount]", "5",  -1)
		dcWithdrawText = strings.Replace(dcWithdrawText, "[currency]", data["name"],  -1)
		currency[data["id"]]["output"] += dcWithdrawText + `</div>`
	}


	currency["1001"] = make(map[string]string)
	currency["1001"]["name"] = "USD"
	currency["1001"]["input"] = `<div class="pull-left"><h4>`+c.Lang["deposit0"]+` USD</h4>
		<select id="ps_select" class="form-control">
		  <option value="pm">Perfect Money</option>
		  <option value="ik">МТС, Магафон, Терминалы, W1, Paxum</option>
		</select>
			<div style="display:block" id="pm_form">
				<form action="https://perfectmoney.is/api/step1.asp" method="POST">
					<input type="hidden" name="PAYEE_ACCOUNT" value="U8777532">
					<input type="hidden" name="PAYEE_NAME" value="DcoinSimple">
					<input type="hidden" name="PAYMENT_ID" value="`+utils.Int64ToStr(c.SessUserId)+`">
					<input type="hidden" name="PAYMENT_UNITS" value="USD">
					<input type="hidden" name="STATUS_URL" value="http://DcoinSimple.com/pm.php">
					<input type="hidden" name="PAYMENT_URL" value="http://DcoinSimple.com">
					<input type="hidden" name="PAYMENT_URL_METHOD" value="LINK">
					<input type="hidden" name="NOPAYMENT_URL" value="http://DcoinSimple.com">
					<input type="hidden" name="NOPAYMENT_URL_METHOD" value="LINK">
					<input type="hidden" name="SUGGESTED_MEMO" value="Dcoins">
					<input type="hidden" name="BAGGAGE_FIELDS" value="">
					<table class="table_out">
					<tbody>
						<tr>
						<td>`+c.Lang["amount_to_pay"]+`</td>
						<td class="form-inline" style="line-height: 35px;"><input name="PAYMENT_AMOUNT" class="form-control" type="text" style="margin-right:5px; width:120px"><input type="submit" value="`+c.Lang["deposit"]+`" class="btn btn-outline btn-success" name="PAYMENT_METHOD"></td>
						</tr>
						<tr>
					 </tbody>
					 </table>
				</form>
			</div>
			<div style="display:none" id="ik_form">
				<form id="payment" name="payment" method="post" action="https://sci.interkassa.com/" enctype="utf-8">
				    <input type="hidden" name="ik_co_id" value="53cfd5e2bf4efc831c9fd661" />
					<input type="hidden" name="ik_pm_no" value="ID_4233" />
					<input type="hidden" name="ik_cur" value="USD" />
					<input type="hidden" name="ik_desc" value="`+utils.Int64ToStr(c.SessUserId)+`" />
				<table class="table_out">
				<tbody>
					<tr>
					<td>`+c.Lang["amount_to_pay"]+`</td>
					<td class="form-inline" style="line-height: 35px;"><input name="ik_am" class="form-control" type="text" style="margin-right:5px; width:120px"><input type="submit" value="`+c.Lang["deposit"]+`" class="btn btn-outline btn-success"></td>
					</tr>
					<tr>
				 </tbody>
				 </table>

					</form>
			</div>
			</div>`

	currency["1001"]["output"] = `<div class="pull-left"><h4>`+c.Lang["withdraw0"]+` USD</h4>
		<table class="table_out">
			<tbody>
			<tr>
			<td>`+c.Lang["withdrawal_on_the_purse"]+`:</td>
			<td class="form-inline"><div class="form-group"><select class="form-control" style="width:300px"><option>Perfect Money [1.5%] [min 10 USD]</option></select></div></td>
			</tr>
			<tr>
			<td>`+c.Lang["purse"]+`:</td>
			<td class="form-inline" style="line-height: 35px;"><input id="account-1001" class="form-control" type="text" style="margin-right:5px; width:300px"></td>
			</tr>
			<tr>
			<td>`+c.Lang["amount_to_withdrawal"]+`:</td>
			<td class="form-inline" style="line-height: 35px;"><input id="amount-1001" class="form-control" type="text"  onkeyup="calc_withdraw_amount(1001, '1.5')" onchange="calc_withdraw_amount(1001, '1.5')" style="margin-right:5px; width:300px"></td>
			</tr>
			<tr>
			<td>`+c.Lang["you_will_receive"]+`:</td>
			<td class="form-inline" style="line-height: 35px"><input  disabled="" id="withdraw_amount-1001" class="form-control" type="text" style="margin-right:5px; width:300px"> </td>
			</tr>
			</tbody></table><div id="alerts-1001"></div><button class="btn btn-outline btn-primary" onclick="withdraw(1001, 'Perfect-money')">`+c.Lang["withdrawal"]+`</button>
			</div><div class="pull-left" style="margin-left:30px; margin-top:43px; border-left: 4px solid #ccc; padding:7px 7px; width:350px">`+c.Lang["withdrawal_within_hours"]+`</div>`

	types := map[string]string{"withdraw": c.Lang["withdraw0"], "adding_funds": c.Lang["deposit0"]}

	// история вывода средств
	myFinanceHistory_ := make(map[int64][]*EmyFinanceType)
	rows, err := c.Query(c.FormatQuery(`
			SELECT amount, wd_amount, close_time, currency_id, method, open_time
			FROM e_withdraw
			WHERE user_id = ?
			ORDER BY open_time DESC
			LIMIT 40
			`), c.SessUserId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	for rows.Next() {
		Finance := new(EmyFinanceType)
		err = rows.Scan(&Finance.Amount, &Finance.WdAmount, &Finance.CloseTime, &Finance.CurrencyId, &Finance.Method, &Finance.OpenTime)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		Finance.Ftype = types["withdraw"]
		Finance.Amount = Finance.WdAmount
		if Finance.CloseTime == 0 {
			Finance.Status = c.Lang["in_process"]
		} else {
			Finance.Status = `<span class="text-success"><strong>`+c.Lang["ready"]+`</strong></span> (`+utils.Int64ToStr(Finance.CloseTime)+`)`
		}
		Finance.Method = Finance.Method + ` (`+currencyList[Finance.CurrencyId]+`)`
		myFinanceHistory_[Finance.OpenTime] = append(myFinanceHistory_[Finance.OpenTime], Finance)
	}


	// история ввода средств
	rows, err = c.Query(c.FormatQuery(`
			SELECT amount, time, currency_id
			FROM e_adding_funds
			WHERE user_id = ?
			ORDER BY time DESC
			LIMIT 40
			`), c.SessUserId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	for rows.Next() {
		Finance := new(EmyFinanceType)
		err = rows.Scan(&Finance.Amount,  &Finance.AddTime, &Finance.CurrencyId)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		Finance.AddType = types["adding_funds"]
		Finance.Status = `<span class="text-success"><strong>`+c.Lang["ready"]+`</strong></span>`
		Finance.Method = `Dcoin (`+currencyList[Finance.CurrencyId]+`)`
		myFinanceHistory_[Finance.AddTime] = append(myFinanceHistory_[Finance.AddTime], Finance)
	}

	// история ввода средств IK
	rows, err = c.Query(c.FormatQuery(`
			SELECT amount, time, currency_id
			FROM e_adding_funds_ik
			WHERE user_id = ?
			ORDER BY time DESC
			LIMIT 40
			`), c.SessUserId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	for rows.Next() {
		Finance := new(EmyFinanceType)
		err = rows.Scan(&Finance.Amount, &Finance.AddTime, &Finance.CurrencyId)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		Finance.AddType = types["adding_funds"]
		Finance.Status = `<span class="text-success"><strong>`+c.Lang["ready"]+`</strong></span>`
		Finance.Method = `Interkassa (`+currencyList[Finance.CurrencyId]+`)`
		myFinanceHistory_[Finance.AddTime] = append(myFinanceHistory_[Finance.AddTime], Finance)
	}


	// история ввода средств PM
	rows, err = c.Query(c.FormatQuery(`
			SELECT amount, time, currency_id
			FROM e_adding_funds_pm
			WHERE user_id = ?
			ORDER BY time DESC
			LIMIT 40
			`), c.SessUserId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	for rows.Next() {
		Finance := new(EmyFinanceType)
		err = rows.Scan(&Finance.Amount, &Finance.AddTime, &Finance.CurrencyId)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		Finance.AddType = types["adding_funds"]
		Finance.Status = `<span class="text-success"><strong>`+c.Lang["ready"]+`</strong></span>`
		Finance.Method = `PerfectMoney (`+currencyList[Finance.CurrencyId]+`)`
		myFinanceHistory_[Finance.AddTime] = append(myFinanceHistory_[Finance.AddTime], Finance)
	}

	//map[int64][]*EmyFinanceType
	var keys []int
	for k := range myFinanceHistory_ {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	var my_finance_history []*EmyFinanceType
	for _, k := range keys {
		for _, data := range myFinanceHistory_[int64(k)] {
			my_finance_history = append(my_finance_history, data)
		}
	}
	///home/z/go-projects/src/github.com/c-darwin/dcoin-go-tmp/packages/controllers/e_my_finance.go:275: cannot use myFinanceHistory_[k] (type []*EmyFinanceType) as type *EmyFinanceType in append



	collapse := c.Parameters["collapse"]

	TemplateStr, err := makeTemplate("e_my_finance", "eMyFinance", &eMyFinancePage {
		Lang:         c.Lang,
		UserId: c.SessUserId,
		MyFinanceHistory : my_finance_history,
		Collapse: collapse,
		Currency: currency,
		CurrencyList: currencyList})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
