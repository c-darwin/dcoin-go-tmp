package controllers
import (
	//"database/sql"
	//_ "github.com/lib/pq"
	//"reflect"
	"fmt"
	"html/template"
	//"bufio"
	"bytes"
	"static"
	"utils"
	"strings"
	"time"
	"math"
//	"log"
	"consts"
)

type walletsListPage struct {
	Lang map[string]string
}

func (c *Controller) WalletsList() (string, error) {

	fmt.Println("WalletsList")

	// валюты
	currencyList, err := c.GetCurrencyList(false)

	var myDcTransactions []map[string]string
	if c.SessUserId > 0 {
		wallets, err := c.GetBalances(c.SessUserId)
		if c.SessRestricted == 0 {
			myDcTransactions, err := c.GetAll("SELECT * FROM "+c.MyPrefix+"my_dc_transactions ORDER BY id DESC LIMIT 0, 100", 100)
		}
	}
	userType := "sendDc";
	projectType := "cfSendDc";
	userTypeId := utils.TypeInt(userType)
	projectTypeId := utils.TypeInt(projectType)
	time := time.Now().Unix()
	currentBlockId, err := c.GetBlockId()
	confirmedBlockId, err := c.GetConfirmedBlockId()

	names := make(map[string]string)
	names["cash_request"] = c.Lang["cash"]
	names["from_mining_id"] = c.Lang["from_mining"]
	names["from_repaid"] = c.Lang["from_repaid_mining"]
	names["from_user"] = c.Lang["from_user"]
	names["node_commission"] = c.Lang["node_commission"]
	names["system_commission"] = c.Lang["system_commission"]
	names["referral"] = c.Lang["referral"]
	names["cf_project"] = "Crowd funding"
	names["cf_project_refund"] = "Crowd funding refund"

	minerId, err := c.GetMyMinerId()
	if err != nil {
		return "", err
	}


	c.r.ParseForm()
	// если юзер кликнул по кнопку "профинансировать" со страницы проекта
	parameters := c.r.FormValue("parameters")
	cfProjectId := utils.StrToInt64(parameters["project_id"])

	// нужна мин. комиссия на пуле для перевода монет
	config := c.GetNodeConfig()

	last_tx, err := c.GetLastTx(c.SessUserId, utils.TypesToIds([]string{"send_dc"}), 1, c.TimeFormat)
	if len(last_tx)>0 {
		last_tx_formatted := utils.MakeLastTx();
	}
	tpl["arbitrationTrustList"], err := c.GetMap("SELECT arbitrator_user_id, conditions FROM arbitration_trust_list LEFT JOIN arbitrator_conditions ON arbitrator_conditions?user_id  =  arbitration_trust_list?arbitrator_user_id WHERE arbitration_trust_list?user_id  =  ? ?, 'list', " []string("arbitrator_user_id", "conditions')", user_id).String()
	if err != nil {
		return p.ErrInfo(err)
	}
	data, err := static.Asset("static/templates/home.html")
	if err != nil {
		return "", err
	}
	alert_success, err := static.Asset("static/templates/alert_success.html")
	if err != nil {
		return "", err
	}
	signatures, err := static.Asset("static/templates/signatures.html")
	if err != nil {
		return "", err
	}
	funcMap := template.FuncMap{
		"ReplaceCurrency": func(text, name string) string { return strings.Replace(text, "[currency]", name, -1) },
	}
	t := template.Must(template.New("template").Funcs(funcMap).Parse(string(data)))
	t = template.Must(t.Parse(string(alert_success)))
	t = template.Must(t.Parse(string(signatures)))
	b := new(bytes.Buffer)
	t.ExecuteTemplate(b, "home", &page{CountSignArr: c.CountSignArr, CountSign: c.CountSign, CalcTotal:calcTotal, Admin: c.Admin, CurrencyPct:currency_pct, SumWallets:sumWallets, Wallets: walletsByCurrency, PromisedAmountListGen: promisedAmountListGen, SessRestricted: c.SessRestricted, SumPromisedAmount: sumPromisedAmount, RandMiners: randMiners, Points: points, Assignments:assignments, CurrencyList:currencyList, ConfirmedBlockId: confirmedBlockId, CashRequests: cashRequests, ShowMap: showMap, BlockId: blockId, UserId: c.SessUserId, PoolAdmin: poolAdmin, Alert: "", MyNotice: c.MyNotice, Lang:  c.Lang, Title: c.Lang["geolocation"]})
	 return b.String(), nil
}
