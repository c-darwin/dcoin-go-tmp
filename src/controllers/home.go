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
)

type page struct {
	Lang map[string]string
	Title string
	Msg string
	Alert string
	MyNotice map[string]string
	PoolAdmin bool
	UserId int64
	CashRequests bool
	ShowMap bool
	BlockId int64
	ConfirmedBlockId int64
	CurrencyList map[string]string
	Assignments int64
	SumPromisedAmount map[string]string
	RandMiners []int64
	Points int64
	SessRestricted int64
}

func (c *Controller) Home() (string, error) {
	fmt.Println("Home1")
	data, err := static.Asset("static/templates/home.html")
	if err != nil {
		return "", err
	}
	alert_success, err := static.Asset("static/templates/alert_success.html")
	if err != nil {
		return "", err
	}
	funcMap := template.FuncMap{
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	t := template.Must(template.New("template").Funcs(funcMap).Parse(string(data)))
	t = template.Must(t.Parse(string(alert_success)))

	var publicKey []byte
	var poolAdmin bool
	var cashRequests bool
	var showMap bool
	if c.SessRestricted == 0 {
		var err error
		publicKey, err = c.GetMyPublicKey(c.MyPrefix)
		if err != nil {
			return "", err
		}
		publicKey = utils.BinToHex(publicKey)
		count, err := c.Single("SELECT count(id) FROM cash_requests WHERE to_user_id  =  ? AND status  =  'pending' AND for_repaid_del_block_id  =  0 AND del_block_id  =  0", c.SessUserId).Int64()
		if err != nil {
			return "", err
		}
		if count > 0 {
			cashRequests = true
		}
		show, err := c.Single("SELECT show_map FROM "+c.MyPrefix+"my_table").Int64()
		if err != nil {
			return "", err
		}
		if show > 0 {
			showMap = true
		}
	}
	if c.Community {
		poolAdminUserId, err := c.GetPoolAdminUserId()
		if err != nil {
			return "", err
		}
		if c.SessUserId == poolAdminUserId {
			poolAdmin = true
		}
	}

	wallets, err := c.GetBalances(c.SessUserId)
	if err != nil {
		return "", err
	}
	//var walletsByCurrency map[string]map[string]string
	walletsByCurrency:=make(map[string]map[string]string)
	for _, data := range wallets {
		walletsByCurrency[data["currency_id"]] = data
	}
	blockId, err := c.GetBlockId()
	if err != nil {
		return "", err
	}
	confirmedBlockId, err := c.GetConfirmedBlockId()
	if err != nil {
		return "", err
	}
	currencyList, err := c.GetCurrencyList(true)
	if err != nil {
		return "", err
	}

	// задания
	var assignments int64
	count, err := c.Single("SELECT count(id) FROM votes_miners WHERE votes_end  =  0 AND type  =  'user_voting'").Int64()
	if err != nil {
		return "", err
	}
	assignments+=count

	// вначале получим ID валют, которые мы можем проверять.
	currencyIds, err := c.GetList("SELECT currency_id FROM promised_amount WHERE status IN ('mining', 'repaid') AND user_id  =  ?", c.SessUserId).String()
	if len(currencyIds) > 0 || c.SessUserId == 1 {
		addSql := ""
		if c.SessUserId == 1 {
			addSql = ""
		} else {
			addSql = "AND currency_id IN ("+strings.Join(currencyIds, ",")+")"
		}
		count, err := c.Single("SELECT count(id) FROM promised_amount WHERE status  =  'pending' AND del_block_id  =  0 "+addSql).Int64()
		if err != nil {
			return "", err
		}
		assignments+=count
	}

	if c.SessRestricted == 0 {
		count, err := c.Single("SELECT count(id) FROM "+c.MyPrefix+"my_tasks WHERE time > ?", time.Now().Unix()).Int64()
		if err != nil {
			return "", err
		}
		assignments-=count
		if assignments < 0 {
			assignments = 0
		}
	}

	// баллы
	points, err := c.Single("SELECT points FROM points WHERE user_id  =  ?", c.SessUserId).Int64()
	if err != nil {
		return "", err
	}

	currency_pct := make(map[int64]map[string]string)
	// проценты
	listPct, err := c.GetMap("SELECT * FROM currency", "id", "name")
	for id, name := range(listPct) {
		pct, err := c.OneRow("SELECT * FROM pct WHERE currency_id  =  ? ORDER BY block_id DESC", id).Float64()
		if err != nil {
			return "", err
		}
		currency_pct[utils.StrToInt64(id)] = map[string]string{"name":name, "miner":utils.Float64ToStr(utils.Round((math.Pow(1+pct["miner"], 3600*24*365)-1)*100, 2)), "user":utils.Float64ToStr(utils.Round((math.Pow(1+pct["user"], 3600*24*365)-1)*100, 2)), "miner_block":utils.Float64ToStr(utils.Round((math.Pow(1+pct["miner"], 120)-1)*100, 4)),  "user_block":utils.Float64ToStr(utils.Round((math.Pow(1+pct["user"], 120)-1)*100, 4)), "miner_sec":utils.Float64ToStr(pct["miner"]), "user_sec":utils.Float64ToStr(pct["user"])}
	}
	// случайне майнеры для нанесения на карту
	maxMinerId, err := c.Single("SELECT max(miner_id) FROM miners_data").Int64()
	if err != nil {
		return "", err
	}
	randMiners, err := c.GetList("SELECT user_id FROM miners_data WHERE status  =  'miner' AND user_id > 7 AND user_id != 106 AND longitude > 0 AND miner_id IN ("+strings.Join(utils.RandSlice(1, maxMinerId, 3), ",")+") LIMIT 3").Int64()
	if err != nil {
		return "", err
	}

	// получаем кол-во DC на кошельках
	sumWallets_, err := c.GetMap("SELECT currency_id, sum(amount) as sum_amount FROM wallets GROUP BY currency_id", "currency_id", "sum_amount")
	if err != nil {
		return "", err
	}
	sumWallets := make(map[int64]float64)
	for currencyId, amount := range(sumWallets_) {
		sumWallets[utils.StrToInt64(currencyId)] = utils.StrToFloat64(amount)
	}

	// получаем кол-во TDC на обещанных суммах, плюсуем к тому, что на кошельках
	sumTdc, err := c.GetMap("SELECT currency_id, sum(tdc_amount) as sum_amount FROM promised_amount GROUP BY currency_id", "currency_id", "sum_amount")
	if err != nil {
		return "", err
	}

	for currencyId, amount := range(sumTdc) {
		currencyIdInt := utils.StrToInt64(currencyId)
		if sumWallets[currencyIdInt] == 0 {
			sumWallets[currencyIdInt] =  utils.StrToFloat64(amount)
		} else {
			sumWallets[currencyIdInt] += utils.StrToFloat64(amount)
		}
	}

	// получаем суммы обещанных сумм
	sumPromisedAmount, err := c.GetMap("SELECT currency_id, sum(amount) as sum_amount FROM promised_amount WHERE status = 'mining' AND del_block_id = 0 AND (cash_request_out_time = 0 OR cash_request_out_time > ?) GROUP BY currency_id", "currency_id", "sum_amount", time.Now().Unix() - c.Variables.Int64["cash_request_time"])
	if err != nil {
		return "", err
	}

	_, _, _, err = c.GetPromisedAmounts(c.SessUserId, c.Variables.Int64["cash_request_time"])


	b := new(bytes.Buffer)
	t.ExecuteTemplate(b, "home", &page{SessRestricted: c.SessRestricted, SumPromisedAmount: sumPromisedAmount, RandMiners: randMiners, Points: points, Assignments:assignments, CurrencyList:currencyList, ConfirmedBlockId: confirmedBlockId, CashRequests: cashRequests, ShowMap: showMap, BlockId: blockId, UserId: c.SessUserId, PoolAdmin: poolAdmin, Alert: "", MyNotice: c.MyNotice, Lang:  c.Lang, Title: c.Lang["geolocation"]})
	 return b.String(), nil
}
