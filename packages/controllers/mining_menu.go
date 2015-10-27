package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"time"
)

type miningMenuPage struct {
	SignData          string
	ShowSignData      bool
	TxType            string
	TxTypeId          int64
	TimeNow           int64
	UserId            int64
	Alert             string
	Lang              map[string]string
	CountSignArr      []int
	CreditId          float64
	CurrencyList      map[int64]string
	LastTxFormatted   string
	MyComments        []map[string]string
	MinerVotesAttempt int64
	Host              string
	Result            string
	NodePrivateKey    string
	Mobile string
}

func (c *Controller) MiningMenu() (string, error) {

	var err error

	if len(c.Parameters["skip_promised_amount"]) > 0 {
		err = c.ExecSql("UPDATE " + c.MyPrefix + "my_table SET hide_first_promised_amount = 1")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	}
	if len(c.Parameters["skip_commission"]) > 0 {
		err = c.ExecSql("UPDATE " + c.MyPrefix + "my_table SET hide_first_commission = 1")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	}

	var result string
	checkCommission := func() error {
		// установлена ли комиссия
		commission, err := c.Single("SELECT commission FROM commission WHERE user_id  =  ?", c.SessUserId).String()
		if err != nil {
			return utils.ErrInfo(err)
		}
		if len(commission) == 0 {
			// возможно юзер уже отправил запрос на добавление комиссии
			last_tx, err := c.GetLastTx(c.SessUserId, utils.TypesToIds([]string{"change_commission"}), 1, c.TimeFormat)
			if err != nil {
				return utils.ErrInfo(err)
			}
			if len(last_tx) > 0 && (len(last_tx[0]["queue_tx"]) > 0 || len(last_tx[0]["tx"]) > 0) {
				// авансом выдаем полное майнерское меню
				result = "full_mining_menu"
			} else {
				// возможно юзер нажал кнопку "пропустить"
				hideFirstCommission, err := c.Single("SELECT hide_first_commission FROM " + c.MyPrefix + "my_table").Int64()
				if err != nil {
					return utils.ErrInfo(err)
				}
				if hideFirstCommission == 0 {
					result = "need_commission"
				} else {
					result = "full_mining_menu"
				}
			}
		} else {
			result = "full_mining_menu"
		}
		return nil
	}

	hostTpl := ""
	// чтобы при добавлении общенных сумм, смены комиссий редиректило сюда
	navigate := "miningMenu"
	if c.SessRestricted != 0 {
		result = "need_email"
	} else {
		myMinerId, err := c.Single("SELECT miner_id FROM miners_data WHERE user_id  =  ?", c.SessUserId).Int64()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if myMinerId == 0 {
			// проверим, послали ли мы запрос в DC-сеть
			data, err := c.OneRow("SELECT node_voting_send_request, http_host as host FROM " + c.MyPrefix + "my_table").String()
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			node_voting_send_request := utils.StrToInt64(data["node_voting_send_request"])
			host := data["host"]
			// если прошло менее 1 часа
			if time.Now().Unix()-node_voting_send_request < 3600 {
				result = "pending"
			} else if node_voting_send_request > 0 {
				// голосование нодов
				nodeVotesEnd, err := c.Single("SELECT votes_end FROM votes_miners WHERE user_id  =  ? AND type  =  'node_voting' ORDER BY id DESC", c.SessUserId).String()
				if err != nil {
					return "", utils.ErrInfo(err)
				}
				if nodeVotesEnd == "1" { // голосование нодов завершено
					userVotesEnd, err := c.Single("SELECT votes_end FROM votes_miners WHERE user_id  =  ? AND type  =  'user_voting' ORDER BY id DESC", c.SessUserId).String()
					if err != nil {
						return "", utils.ErrInfo(err)
					}
					if userVotesEnd == "1" { // юзерское голосование закончено
						result = "bad"
					} else if userVotesEnd == "0" { // идет юзерское голосование
						result = "users_pending"
					} else {
						result = "bad_photos_hash"
						hostTpl = host
					}
				} else if nodeVotesEnd == "0" && time.Now().Unix()-node_voting_send_request < 86400 { // голосование нодов началось, ждем.
					result = "nodes_pending"
				} else if nodeVotesEnd == "0" && time.Now().Unix()-node_voting_send_request >= 86400 { // голосование нодов удет более суток и еще не завершилось
					result = "resend"
				} else { // запрос в DC-сеть еще не дошел и голосования не начались
					// если прошло менее 1 часа
					if time.Now().Unix()-node_voting_send_request < 3600 {
						result = "pending"
					} else { // где-то проблема и запрос не ушел.
						result = "resend"
					}
				}
			} else { // запрос на получение статуса "майнер" мы еще не слали
				result = "null"
			}
		} else {
			// добавлена ли обещанная сумма
			promisedAmount, err := c.Single("SELECT id FROM promised_amount WHERE user_id  =  ?", c.SessUserId).Int64()
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			if promisedAmount == 0 {
				// возможно юзер уже отправил запрос на добавление обещенной суммы
				last_tx, err := c.GetLastTx(c.SessUserId, utils.TypesToIds([]string{"new_promised_amount"}), 1, c.TimeFormat)
				if len(last_tx) > 0 && (len(last_tx[0]["queue_tx"]) > 0 || len(last_tx[0]["tx"]) > 0) {
					// установлена ли комиссия
					err = checkCommission()
					if err != nil {
						return "", utils.ErrInfo(err)
					}
				} else {
					// возможно юзер нажал кнопку "пропустить"
					hideFirstPromisedAmount, err := c.Single("SELECT hide_first_promised_amount FROM " + c.MyPrefix + "my_table").Int64()
					if err != nil {
						return "", utils.ErrInfo(err)
					}
					if hideFirstPromisedAmount == 0 {
						result = "need_promised_amount"
					} else {
						err = checkCommission()
						if err != nil {
							return "", utils.ErrInfo(err)
						}
					}
				}
			} else {
				// установлена ли комиссия
				err = checkCommission()
				if err != nil {
					return "", utils.ErrInfo(err)
				}
			}
		}
	}

	var minerVotesAttempt int64
	var myComments []map[string]string
	c.Navigate = navigate
	lastTxFormatted := ""
	tplName := ""
	tplTitle := ""
	log.Debug("result:", result)
	var nodePrivateKey string
	if result == "null" {
		tplName = "upgrade_0"
		tplTitle = "upgrade0"
		return c.Upgrade0()
	} else if result == "need_email" {
		tplName = "sign_up_in_the_pool"
		tplTitle = "signUpInThePool"
		//  сгенерим ключ для нода
		nodePrivateKey, _ = utils.GenKeys()
	} else if result == "need_promised_amount" {
		tplName = "promised_amount_add"
		tplTitle = "upgrade"
		return c.NewPromisedAmount()
	} else if result == "need_commission" {
		tplName = "change_commission"
		tplTitle = "changeCommission"
		return c.ChangeCommission()
	} else if result == "full_mining_menu" {
		tplName = "mining_menu"
		tplTitle = "miningMenu"
		last_tx, err := c.GetLastTx(c.SessUserId, utils.TypesToIds([]string{"new_user", "new_miner", "new_promised_amount", "change_promised_amount", "votes_miner", "change_geolocation", "votes_promised_amount", "del_promised_amount", "cash_request_out", "cash_request_in", "votes_complex", "for_repaid_fix", "new_holidays", "actualization_promised_amounts", "mining", "new_miner_update", "change_host", "change_commission"}), 3, c.TimeFormat)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if len(last_tx) > 0 {
			lastTxFormatted, _ = utils.MakeLastTx(last_tx, c.Lang)
		}
	} else {
		// сколько у нас осталось попыток стать майнером.
		countAttempt, err := c.CountMinerAttempt(c.SessUserId, "user_voting")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		minerVotesAttempt = c.Variables.Int64["miner_votes_attempt"] - countAttempt

		// комментарии проголосовавших
		myComments, err = c.GetAll(`SELECT * FROM `+c.MyPrefix+`my_comments WHERE comment != 'null' AND type NOT IN ('arbitrator','seller')`, -1)
		tplName = "upgrade"
		tplTitle = "upgrade"
	}

	log.Debug("tplName, tplTitle %v, %v", tplName, tplTitle)
	TemplateStr, err := makeTemplate(tplName, tplTitle, &miningMenuPage{
		Alert:             c.Alert,
		Lang:              c.Lang,
		CountSignArr:      c.CountSignArr,
		ShowSignData:      c.ShowSignData,
		UserId:            c.SessUserId,
		SignData:          "",
		CurrencyList:      c.CurrencyList,
		LastTxFormatted:   lastTxFormatted,
		MyComments:        myComments,
		Result:            result,
		NodePrivateKey:    nodePrivateKey,
		MinerVotesAttempt: minerVotesAttempt,
		Mobile: utils.Mobile(),
		Host:              hostTpl})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
