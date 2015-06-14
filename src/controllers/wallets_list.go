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
//	"strings"
	"time"
	//"math"
	"log"
	//"consts"
	"encoding/json"
)

type walletsListPage struct {
	Alert string
	Lang map[string]string
	CurrencyList map[int]string
	Wallets []utils.DCAmounts
	MyDcTransactions []map[string]string
	UserTypeId int64
	UserType string
	ProjectTypeId int64
	ProjectType string
	Time int64
	CurrentBlockId int64
	ConfirmedBlockId int64
	Community bool
	MinerId int64
	UserId int64
	UserIdStr string
	Config map[string]string
	ConfigCommission map[string][]float64
	LastTxFormatted string
	ArbitrationTrustList map[string]map[string][]string
	ShowSignData bool
}

func (c *Controller) WalletsList() (string, error) {

	fmt.Println("WalletsList")

	// валюты
	currencyList, err := c.GetCurrencyList(true)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	log.Println("currencyList", currencyList)

	var wallets []utils.DCAmounts
	var myDcTransactions []map[string]string
	if c.SessUserId > 0 {
		wallets, err = c.GetBalances(c.SessUserId)
		if c.SessRestricted == 0 {
			myDcTransactions, err = c.GetAll("SELECT * FROM "+c.MyPrefix+"my_dc_transactions ORDER BY id DESC LIMIT 0, 100", 100)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			for id, data := range myDcTransactions {
				t := time.Unix(utils.StrToInt64(data["time"]), 0)
				timeFormatted := t.Format(c.TimeFormat)
				log.Println("timeFormatted", utils.StrToInt64(data["time"]), timeFormatted, c.TimeFormat )
				myDcTransactions[id]["timeFormatted"] = timeFormatted
			}
		}
	}
	userType := "sendDc";
	projectType := "cfSendDc";
	userTypeId := utils.TypeInt(userType)
	projectTypeId := utils.TypeInt(projectType)
	timeNow := time.Now().Unix()
	currentBlockId, err := c.GetBlockId()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	confirmedBlockId, err := c.GetConfirmedBlockId()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

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

	minerId, err := c.GetMyMinerId(c.SessUserId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	c.r.ParseForm()
	// если юзер кликнул по кнопку "профинансировать" со страницы проекта
	//parameters := c.r.FormValue("parameters")
	//cfProjectId := utils.StrToInt64(parameters["project_id"])

	// нужна мин. комиссия на пуле для перевода монет
	config, err := c.GetNodeConfig()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	configCommission := make(map[string][]float64)
	if len(config["commission"]) > 0 {
		err = json.Unmarshal([]byte(config["commission"]), &configCommission)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	}

	last_tx, err := c.GetLastTx(c.SessUserId, utils.TypesToIds([]string{"send_dc"}), 1, c.TimeFormat)
	lastTxFormatted := ""
	if len(last_tx)>0 {
		lastTxFormatted = utils.MakeLastTx(last_tx, c.Lang);
	}
	arbitrationTrustList_, err := c.GetMap(`
			SELECT arbitrator_user_id,
					 	conditions
			FROM arbitration_trust_list
			LEFT JOIN arbitrator_conditions ON arbitrator_conditions.user_id = arbitration_trust_list.arbitrator_user_id
			WHERE arbitration_trust_list.user_id = ?
	`, "arbitrator_user_id", "conditions", c.SessUserId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	arbitrationTrustList :=make(map[string]map[string][]string)
	var jsonMap map[string][]string
	for arbitrator_user_id, conditions := range arbitrationTrustList_ {
		err = json.Unmarshal([]byte(conditions), &jsonMap)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		arbitrationTrustList[arbitrator_user_id] = make(map[string][]string)
		arbitrationTrustList[arbitrator_user_id] = jsonMap
	}
	log.Println("arbitrationTrustList", arbitrationTrustList)

	data, err := static.Asset("static/templates/wallets_list.html")
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

	t := template.Must(template.New("template").Parse(string(data)))
	t = template.Must(t.Parse(string(alert_success)))
	t = template.Must(t.Parse(string(signatures)))
	b := new(bytes.Buffer)
	t.ExecuteTemplate(b, "walletsList", &walletsListPage{UserIdStr: utils.Int64ToStr(c.SessUserId),Alert: "", Community: c.Community, ConfigCommission: configCommission, ProjectType: projectType, UserType: userType, UserId: c.SessUserId, Lang: c.Lang, CurrencyList: currencyList, Wallets: wallets, MyDcTransactions: myDcTransactions, UserTypeId: userTypeId, ProjectTypeId: projectTypeId, Time: timeNow, CurrentBlockId: currentBlockId, ConfirmedBlockId: confirmedBlockId, MinerId: minerId, Config: config, LastTxFormatted: lastTxFormatted, ArbitrationTrustList: arbitrationTrustList, ShowSignData: c.ShowSignData})
	 return b.String(), nil
}
