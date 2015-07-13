package controllers
import (
	"html/template"
	"bytes"
	"github.com/c-darwin/dcoin-tmp/packages/static"
	"github.com/c-darwin/dcoin-tmp/packages/utils"
	"time"
	"log"
	"encoding/json"
)

type walletsListPage struct {
	SignData string
	CfProjectId int64
	Alert string
	Lang map[string]string
	CurrencyList map[int64]string
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
	ConfigCommission map[int64][]float64
	LastTxFormatted string
	ArbitrationTrustList map[int64]map[int64][]string
	ShowSignData bool
	Names map[string]string
	CountSignArr []int
}

func (c *Controller) WalletsList() (string, error) {

	var err error

	// валюты
	currencyList:=c.CurrencyListCf

	confirmedBlockId:=c.ConfirmedBlockId

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
				myDcTransactions[id]["numBlocks"] = "0"
				blockId := utils.StrToInt64(data["block_id"])
				if blockId > 0 {
					myDcTransactions[id]["numBlocks"] = utils.Int64ToStr(confirmedBlockId - blockId)
				}
			}
		}
	}
	userType := "SendDc";
	projectType := "CfSendDc";
	userTypeId := utils.TypeInt(userType)
	projectTypeId := utils.TypeInt(projectType)
	timeNow := time.Now().Unix()
	currentBlockId, err := c.GetBlockId()
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
	cfProjectId := int64(utils.StrToFloat64(c.Parameters["projectId"]))

	last_tx, err := c.GetLastTx(c.SessUserId, utils.TypesToIds([]string{"send_dc"}), 1, c.TimeFormat)
	lastTxFormatted := ""
	if len(last_tx)>0 {
		lastTxFormatted, _ = utils.MakeLastTx(last_tx, c.Lang);
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
	arbitrationTrustList :=make(map[int64]map[int64][]string)
	var jsonMap map[string][]string
	for arbitrator_user_id, conditions := range arbitrationTrustList_ {
		err = json.Unmarshal([]byte(conditions), &jsonMap)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		uidInt := utils.StrToInt64(arbitrator_user_id)
		arbitrationTrustList[uidInt] = make(map[int64][]string)
		for currenycId, data:= range jsonMap{
			arbitrationTrustList[uidInt][utils.StrToInt64(currenycId)] = data
		}
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
	funcMap := template.FuncMap{
		"strToInt64": func(text string) int64 {
			return utils.StrToInt64(text)
		},
		"makeCurrencyName": func(currencyId int64) string {
			if currencyId >= 1000 {
				return ""
			} else {
				return "d"
			}
		},
	}
	t := template.Must(template.New("template").Funcs(funcMap).Parse(string(data)))
	t = template.Must(t.Parse(string(alert_success)))
	t = template.Must(t.Parse(string(signatures)))
	b := new(bytes.Buffer)
	err = t.ExecuteTemplate(b, "walletsList", &walletsListPage{
		CountSignArr: c.CountSignArr,
		CfProjectId: cfProjectId,
		Names: names,
		UserIdStr: utils.Int64ToStr(c.SessUserId),
		Alert: c.Alert,
		Community: c.Community,
		ConfigCommission: c.ConfigCommission,
		ProjectType: projectType,
		UserType: userType,
		UserId: c.SessUserId,
		Lang: c.Lang,
		CurrencyList: currencyList,
		Wallets: wallets,
		MyDcTransactions: myDcTransactions,
		UserTypeId: userTypeId,
		ProjectTypeId: projectTypeId,
		Time: timeNow,
		CurrentBlockId: currentBlockId,
		ConfirmedBlockId: confirmedBlockId,
		MinerId: minerId,
		Config: c.NodeConfig,
		LastTxFormatted: 	lastTxFormatted,
		ArbitrationTrustList: arbitrationTrustList,
		ShowSignData: c.ShowSignData,
		SignData: ""})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return b.String(), nil
}
