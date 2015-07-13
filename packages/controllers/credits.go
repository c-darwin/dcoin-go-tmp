package controllers
import (
	"html/template"
	"bytes"
	"github.com/c-darwin/dcoin-tmp/packages/static"
	"github.com/c-darwin/dcoin-tmp/packages/utils"
	"time"
)

type creditsPage struct {
	SignData string
	ShowSignData bool
	TxType string
	TxTypeId int64
	TimeNow int64
	UserId int64
	Alert string
	Lang map[string]string
	CountSignArr []int
	I_debtor []*credit
	I_creditor []*credit
	CurrencyList map[int64]string
	CreditPart float64
}

type credit struct {
	Id int64
	Pct float64
	Time int64
	Amount float64
	Currency_id int64
	To_user_id int64
	From_user_id int64
}

func (c *Controller) Credits() (string, error) {

	txType := "ChangeCreditPart";
	txTypeId := utils.TypeInt(txType)
	timeNow := time.Now().Unix()

	var I_debtor, I_creditor []*credit

	rows, err := c.Query(c.FormatQuery("SELECT id, pct, time, amount, currency_id, from_user_id, to_user_id FROM credits WHERE (from_user_id = ? OR to_user_id = ?) AND del_block_id = 0"), c.SessUserId, c.SessUserId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	defer rows.Close()
	for rows.Next() {
		var pct, amount float64
		var id, currency_id, from_user_id, txtime, to_user_id int64
		err = rows.Scan(&id, &pct, &txtime, &amount, &currency_id, &from_user_id, &to_user_id)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		credit_:=&credit{Id: id, Pct:pct, Time:txtime, Amount:amount, Currency_id:currency_id, From_user_id:from_user_id, To_user_id:to_user_id}
		if c.SessUserId == from_user_id {
			I_debtor = append(I_debtor, credit_)
		} else {
			I_creditor = append(I_creditor, credit_)
		}
	}

	creditPart, err := c.Single("SELECT credit_part FROM users WHERE user_id  =  ?", c.SessUserId).Float64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	data, err := static.Asset("static/templates/credits.html")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	signatures, err := static.Asset("static/templates/signatures.html")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	alert_success, err := static.Asset("static/templates/alert_success.html")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	funcMap := template.FuncMap{
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	t := template.Must(template.New("template").Funcs(funcMap).Parse(string(data)))
	t = template.Must(t.Parse(string(alert_success)))
	t = template.Must(t.Parse(string(signatures)))
	b := new(bytes.Buffer)
	err = t.ExecuteTemplate(b, "credits", &creditsPage{
		Alert: c.Alert,
		Lang: c.Lang,
		CountSignArr: c.CountSignArr,
		ShowSignData: c.ShowSignData,
		UserId: c.SessUserId,
		TimeNow: timeNow,
		TxType: txType,
		TxTypeId: txTypeId,
		SignData: "",
		CurrencyList: c.CurrencyListCf,
		CreditPart: creditPart,
		I_debtor: I_debtor,
		I_creditor: I_creditor})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return b.String(), nil
}
