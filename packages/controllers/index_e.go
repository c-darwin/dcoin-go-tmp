package controllers

import (
	"bytes"
	"github.com/c-darwin/dcoin-go-tmp/packages/static"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"html/template"
	"net/http"
)

type indexE struct {
	MyWallets []map[string]string
	Lang      map[string]string
	Nav       template.JS
	UserId    int64
}

func IndexE(w http.ResponseWriter, r *http.Request) {

	var err error

	if utils.DB != nil && utils.DB.DB != nil {

		sess, _ := globalSessions.SessionStart(w, r)
		defer sess.SessionRelease(w)

		c := new(Controller)
		c.r = r
		c.SessUserId = GetSessEUserId(sess)
		c.DCDB = utils.DB

		r.ParseForm()

		c.Parameters, err = c.GetParameters()
		log.Debug("parameters=", c.Parameters)

		lang := GetLang(w, r, c.Parameters)
		log.Debug("lang", lang)
		c.Lang = globalLangReadOnly[lang]

		var myWallets []map[string]string
		if c.SessUserId > 0 {
			myWallets, err = c.getMyWallets()
			if err != nil {
				w.Write([]byte(utils.ErrInfo(err).Error()))
				log.Error("%v", err)
				return
			}
		}

		data, err := static.Asset("static/templates/index_e.html")
		if err != nil {
			w.Write([]byte(utils.ErrInfo(err).Error()))
			log.Error("%v", err)
			return
		}
		t := template.New("template")
		t, err = t.Parse(string(data))
		if err != nil {
			w.Write([]byte(utils.ErrInfo(err).Error()))
			log.Error("%v", err)
			return
		}
		b := new(bytes.Buffer)
		err = t.Execute(b, &indexE{MyWallets: myWallets, Lang: c.Lang, UserId: c.SessUserId})
		if err != nil {
			log.Error("%v", err)
			w.Write([]byte(utils.ErrInfo(err).Error()))
		} else {
			w.Write(b.Bytes())
		}

	}
}

func (c*Controller) getMyWallets() ([]map[string]string, error) {
	var myWallets []map[string]string
	eCurrency, err := c.GetAll(`SELECT name as currency_name, id FROM e_currency ORDER BY sort_id ASC`, -1)
	if err != nil {
		return myWallets, utils.ErrInfo(err)
	}
	for _, data := range eCurrency {
		wallet, err := c.OneRow("SELECT * FROM e_wallets WHERE user_id  =  ? AND currency_id  =  ?", c.SessUserId, data["id"]).String()
		if err != nil {
			return myWallets, utils.ErrInfo(err)
		}
		if len(wallet) > 0 {
			amount := utils.StrToFloat64(wallet["amount"])
			profit, err := utils.DB.CalcProfitGen(utils.StrToInt64(wallet["currency_id"]), amount, 0, utils.StrToInt64(wallet["last_update"]), utils.Time(), "wallet")
			if err != nil {
				return myWallets, utils.ErrInfo(err)
			}
			myWallets = append(myWallets, map[string]string{"amount": utils.Float64ToStr(amount + profit), "currency_id": data["currency_name"], "last_update": wallet["last_update"]})
		}
	}
	return myWallets, nil
}