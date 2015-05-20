package dcparser

import (
	"fmt"
	"utils"
	"encoding/json"
	//"regexp"
	//"math"
	//"strings"
//	"os"
	//"time"
	//"strings"
	"time"
)


func (p *Parser) VotesComplexInit() (error) {
	var err error
	var fields []string
	fields = []string {"json_data", "sign"}
	p.TxMap, err = p.GetTxMap(fields);
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}
type vComplex struct {
	Currency map[string][]float64 `json:"currency"`
	Referral map[string]uint `json:"referral"`
	Admin uint64 `json:"admin"`
}
func (p *Parser) VotesComplexFront() (error) {

	err := p.generalCheck()
	if err != nil {
		return p.ErrInfo(err)
	}

	// является ли данный юзер майнером
	err = p.checkMiner(p.TxMap["user_id"])
	if err != nil {
		return p.ErrInfo(err)
	}

	var txTime int64
	if p.BlockData!=nil {
		txTime = p.BlockData.Time
	} else {
		txTime = time.Now().Unix() - 30
	}

	// прошло ли 30 дней с момента регистрации майнера
	err = p.checkMinerNewbie()
	if err != nil {
		return p.ErrInfo(err)
	}

	forSign := fmt.Sprintf("%s,%s,%s,%s", p.TxMap["type"], p.TxMap["time"], p.TxMap["user_id"], p.TxMap["json_data"])
	CheckSignResult, err := utils.CheckSign(p.PublicKeys, forSign, p.TxMap["sign"], false);
	if err != nil {
		return p.ErrInfo(err)
	}
	if !CheckSignResult {
		return p.ErrInfo("incorrect sign")
	}


	currencyVotes := make(map[string][]float64)
	var doubleCheck []int64
 	 // раньше не было рефских
	if p.BlockData==nil || p.BlockData.BlockId > 77951 {

		vComplex := new(vComplex)
		err = json.Unmarshal(p.TxMap["json_data"], &vComplex)
		if err != nil {
			return p.ErrInfo(err)
		}

		if vComplex.Referral == nil {
			return p.ErrInfo("!Referral")
		}
		if vComplex.Currency == nil {
			return p.ErrInfo("!Currency")
		}
		if p.BlockData==nil || p.BlockData.BlockId > 153750 {
			if vComplex.Admin > 0 {
				adminUserId, err := p.Single("SELECT user_id FROM users WHERE user_id  =  ?",vComplex.Admin).Int64()
				if err != nil {
					return p.ErrInfo(err)
				}
				if adminUserId == 0 {
					return p.ErrInfo("incorrect admin user_id")
				}
			}
		}
		if !utils.CheckInputData(vComplex.Referral["first"], "referral") || !utils.CheckInputData(vComplex.Referral["second"], "referral") || !utils.CheckInputData(vComplex.Referral["third"], "referral") {
			return p.ErrInfo("incorrect referral")
		}
		currencyVotes = vComplex.Currency
	} else {
		vComplex := make(map[string][]float64)
		err = json.Unmarshal(p.TxMap["json_data"], &vComplex)
		if err != nil {
			return p.ErrInfo(err)
		}
		currencyVotes = vComplex
	}
	for currencyId, data := range currencyVotes {
		if !utils.CheckInputData(currencyId, "int") {
			return p.ErrInfo("incorrect currencyId")
		}

		// проверим, что нет дублей
		if (utils.InSliceInt64(int64(currencyId), doubleCheck)) {
			return p.ErrInfo("double currencyId")
		}
		doubleCheck = append(doubleCheck, int64(currencyId))

		// есть ли такая валюта
		currencyId, err = p.Single("SELECT id FROM currency WHERE id  =  ?", currencyId).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		if currencyId == 0 {
			return p.ErrInfo("incorrect currencyId")
		}
		// у юзера по данной валюте должна быть обещанная сумма, которая имеет статус mining/repaid и находится с таким статусом >90 дней
		id, err := p.Single(`
			SELECT id FROM promised_amount
			WHERE currency_id  =  ? AND user_id  =  ? AND status IN ('mining', 'repaid') AND start_time < ? AND start_time > 0 AND del_block_id  =  0 AND del_mining_block_id  =  0`, currencyId, p.TxMap["user_id"], (txTime - p.Variables.Int64["min_hold_time_promise_amount"])).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		if id == 0 {
			return p.ErrInfo("incorrect currencyId")
		}

		// если по данной валюте еще не набралось >1000 майнеров, то за неё голосовать нельзя.
		countMiners, err := p.Single(`
			SELECT count (*) FROM
			(SELECT user_id
			FROM promised_amount
			WHERE start_time < ? AND del_block_id  =  0 AND status IN ('mining', 'repaid') AND currency_id  =  ? AND del_block_id  =  0 AND del_mining_block_id  =  0
			GROUP BY user_id) as t1`, currencyId, (txTime - p.Variables.Int64["min_hold_time_promise_amount"])).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		if countMiners < p.Variables.Int64["min_miners_of_voting"] {
			return p.ErrInfo("countMiners")
		}
	}

	return nil
}

func (p *Parser) VotesComplex() (error) {

	// начисляем баллы
	p.points(p.Variables.Int64["promised_amount_points"])

	// логируем, чтобы юзер {$this->tx_data['user_id']} не смог повторно проголосовать
	err := p.ExecSql("INSERT INTO log_votes ( user_id, voting_id, type ) VALUES ( ?, ?, 'promised_amount' )", p.TxMap["user_id"], p.TxMap["promised_amount_id"])
	if err != nil {
		return p.ErrInfo(err)
	}

	// обновляем голоса
	err = p.ExecSql("UPDATE promised_amount SET votes_"+string(p.TxMap["result"])+" = votes_"+string(p.TxMap["result"])+" + 1 WHERE id = ?", p.TxMap["promised_amount_id"])
	if err != nil {
		return p.ErrInfo(err)
	}

	promisedAmountData, err := p.OneRow("SELECT log_id, status, start_time, tdc_amount_update, user_id, votes_start_time, votes_0, votes_1 FROM promised_amount WHERE id  =  ?", p.TxMap["promised_amount_id"]).String()
	if err != nil {
		return p.ErrInfo(err)
	}

	data := make(map[string]int64)
	data["count_miners"], err = p.Single("SELECT count(miner_id) FROM miners").Int64()
	if err != nil {
		return p.ErrInfo(err)
	}

	data["votes_0"] = utils.StrToInt64(promisedAmountData["votes_0"])
	data["votes_1"] =  utils.StrToInt64(promisedAmountData["votes_1"])
	data["votes_start_time"] =  utils.StrToInt64(promisedAmountData["votes_start_time"])
	data["votes_0_min"] = p.Variables.Int64["promised_amount_votes_0"]
	data["votes_1_min"] = p.Variables.Int64["promised_amount_votes_1"]
	data["votes_period"] = p.Variables.Int64["promised_amount_votes_period"]

	// -----------------------------------------------------------------------------
	// если голос решающий или голос админа
	// голос админа - решающий только при <1000 майнеров.
	// -----------------------------------------------------------------------------
	err = p.getAdminUserId()
	if err != nil {
		return p.ErrInfo(err)
	}

	if p.check24hOrAdminVote(data) {

		// нужно залогировать, т.к. не известно, какие были status и tdc_amount_update
		logId, err := p.ExecSqlGetLastInsertId("INSERT INTO log_promised_amount ( status, start_time, tdc_amount_update, block_id, prev_log_id ) VALUES ( ?, ?, ?, ?, ? )", data["status"], data["start_time"], data["tdc_amount_update"], p.BlockData.BlockId, data["log_id"])
		if err != nil {
			return p.ErrInfo(err)
		}

		// перевесили голоса "за" или 1 голос от админа
		if p.checkTrueVotes(data) {
			err = p.ExecSql("UPDATE promised_amount SET status = 'mining', start_time = ?, tdc_amount_update = ?, log_id = ? WHERE id = ?", p.BlockData.Time, p.BlockData.Time, logId, p.TxMap["promised_amount_id"])
			if err != nil {
				return p.ErrInfo(err)
			}
			// есть ли у данного юзера woc
			woc, err := p.Single("SELECT id FROM promised_amount WHERE currency_id  =  1 AND user_id  =  ?", data["user_id"]).Int64()
			if err != nil {
				return p.ErrInfo(err)
			}
			if woc == 0 {
				wocAmount, err := p.Single("SELECT amount FROM max_promised_amounts WHERE id  =  1 ORDER BY time DESC").String()
				if err != nil {
					return p.ErrInfo(err)
				}
				// добавляем WOC
				err = p.ExecSql("INSERT INTO promised_amount ( user_id, amount, currency_id, start_time, status, tdc_amount_update, woc_block_id ) VALUES ( ?, ?, 1, ?, 'mining', ?, ? )", data["user_id"], wocAmount, p.BlockData.Time, p.BlockData.Time, p.BlockData.BlockId)
				if err != nil {
					return p.ErrInfo(err)
				}
			}
		} else { // перевесили голоса "против"
			err = p.ExecSql("UPDATE promised_amount SET status = 'rejected', start_time = 0, tdc_amount_update = ?, log_id = ? WHERE id = ?", p.BlockData.Time, logId, p.TxMap["promised_amount_id"])
			if err != nil {
				return p.ErrInfo(err)
			}
		}
	}

	// возможно с голосом пришел коммент
	myUserId, _, myPrefix, _ , err:= p.GetMyUserId(utils.BytesToInt64(p.TxMap["user_id"]))
	if err != nil {
		return err
	}
	if p.TxUserID == myUserId {
		err = p.ExecSql("INSERT INTO "+myPrefix+"my_comments ( type, id, comment ) VALUES ( 'promised_amount', ?, ? )", p.TxMap["promised_amount_id"], p.TxMap["comment"])
		if err != nil {
			return p.ErrInfo(err)
		}
	}

	return nil
}

func (p *Parser) VotesComplexRollback() (error) {

	// вычитаем баллы
	p.pointsRollback(p.Variables.Int64["promised_amount_points"])

	// удаляем логирование, чтобы юзер {$this->tx_data['user_id']} не смог повторно проголосовать
	err := p.ExecSql("DELETE FROM log_votes WHERE user_id = ? AND voting_id = ? AND type = 'promised_amount'", p.TxMap["user_id"], p.TxMap["promised_amount_id"])
	if err != nil {
		return p.ErrInfo(err)
	}

	// обновляем голоса
	err = p.ExecSql("UPDATE promised_amount SET votes_"+string(p.TxMap["result"])+" = votes_"+string(p.TxMap["result"])+" - 1 WHERE id = ?", p.TxMap["promised_amount_id"])
	if err != nil {
		return p.ErrInfo(err)
	}
	data, err := p.OneRow("SELECT status, user_id, log_id FROM promised_amount WHERE id  =  ?", p.TxMap["promised_amount_id"]).String()
	if err != nil {
		return p.ErrInfo(err)
	}

	// если статус mining или rejected, значит голос был решающим
	if data["status"] == "mining" || data["status"] == "rejected" {

		// восстановим из лога
		logData, err := p.OneRow("SELECT status, start_time, tdc_amount_update, prev_log_id FROM log_promised_amount WHERE log_id  =  ?", data["log_id"]).String()
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.ExecSql("UPDATE promised_amount SET status = ?, start_time = ?, tdc_amount_update = ?, log_id = ? WHERE id = ?", logData["status"], logData["start_time"], logData["tdc_amount_update"], logData["prev_log_id"], p.TxMap["promised_amount_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
		// подчищаем _log
		err = p.ExecSql("DELETE FROM log_promised_amount WHERE log_id = ?", data["log_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
		p.rollbackAI("log_promised_amount", 1)

		// был ли добавлен woc
		woc, err := p.Single("SELECT id FROM promised_amount WHERE currency_id  =  1 AND woc_block_id  =  ? AND user_id  =  ?", p.BlockData.BlockId, data["user_id"]).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		if woc > 0 {
			err = p.ExecSql("DELETE FROM promised_amount WHERE id = ?", woc)
			if err != nil {
				return p.ErrInfo(err)
			}
			p.rollbackAI("promised_amount", 1)
		}
	}



	return nil
}

func (p *Parser) VotesComplexRollbackFront() error {
	return p.maxDayVotesRollback()
}
