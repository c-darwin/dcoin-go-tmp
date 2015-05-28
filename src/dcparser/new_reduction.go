package dcparser

import (
	"fmt"
	"utils"
	//"encoding/json"
	//"regexp"
	//"math"
	//"strings"
//	"os"
	//"time"
	//"strings"
	//"sort"
//	"time"
	//"consts"
	"consts"
)


func (p *Parser) NewReductionInit() (error) {
	var err error
	var fields []string
	if p.BlockData!=nil && p.BlockData.BlockId < 85849 {
		fields = []string {"currency_id", "pct", "sign"}
	} else {
		fields = []string {"currency_id", "pct", "reduction_type", "sign"}
	}
	p.TxMap, err = p.GetTxMap(fields);
	p.TxMapS, err = p.GetTxMapStr(fields);
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}
func (p *Parser) NewReductionFront() (error) {

	err := p.generalCheck()
	if err != nil {
		return p.ErrInfo(err)
	}

	// является ли данный юзер майнером
	err = p.checkMiner(p.TxUserID)
	if err != nil {
		return p.ErrInfo(err)
	}

	verifyData := map[string]string {"currency_id":"int", "pct":"int"}
	err = p.CheckInputData(verifyData)
	if err != nil {
		return p.ErrInfo(err)
	}

	if !utils.InSliceInt64(utils.BytesToInt64(p.TxMap["pct"]), consts.ReductionDC) {
		return p.ErrInfo("incorrect pct")
	}

	if p.BlockData!=nil && p.BlockData.BlockId < 85849 {
		// для всех тр-ий из старых блоков просто присваем manual, т.к. там не было других типов
		p.TxMapS["reduction_type"]	= "manual"
	} else {
		verifyData := map[string]string {"reduction_type":"reduction_type"}
		err = p.CheckInputData(verifyData)
		if err != nil {
			return p.ErrInfo(err)
		}
	}

	nodePublicKey, err := p.GetNodePublicKey(p.TxUserID)
	if err != nil {
		return p.ErrInfo(err)
	}
	if len(nodePublicKey) == 0 {
		return p.ErrInfo("incorrect user_id")
	}

	currencyId, err := p.CheckCurrencyId(utils.BytesToInt64(p.TxMap["currency_id"]))
	if err != nil {
		return p.ErrInfo(err)
	}
	if currencyId == 0 {
		return p.ErrInfo("incorrect currency_id")
	}

	forSign:=""
	if p.BlockData!=nil && p.BlockData.BlockId < 85849 {
		forSign = fmt.Sprintf("%s,%s,%s,%s,%s", p.TxMapS["type"], p.TxMapS["time"], p.TxMapS["user_id"], p.TxMapS["currency_id"], p.TxMapS["pct"])
	} else {
		forSign = fmt.Sprintf("%s,%s,%s,%s,%s,%s", p.TxMapS["type"], p.TxMapS["time"], p.TxMapS["user_id"], p.TxMapS["currency_id"], p.TxMapS["pct"], p.TxMapS["reduction_type"])
	}
	CheckSignResult, err := utils.CheckSign([][]byte{nodePublicKey}, forSign, p.TxMap["sign"], true);
	if err != nil {
		return p.ErrInfo(err)
	}
	if !CheckSignResult {
		return p.ErrInfo("incorrect sign")
	}
	if p.TxMapS["reduction_type"] == "manual" {
		// проверим, прошло ли 2 недели с момента последнего reduction
		reductionTime, err := p.Single("SELECT max(time) FROM reduction WHERE currency_id  =  ? AND type  =  'manual'", p.TxMapS["currency_id"]).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		if p.TxTime - reductionTime <= p.Variables.Int64["reduction_period"] {
			return p.ErrInfo("reduction_period error")
		}
	} else {
		reductionTime, err := p.Single("SELECT max(time) FROM reduction WHERE currency_id  =  ? AND type  =  'auto'", p.TxMapS["currency_id"]).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		// или 48 часов, если это авто-урезание
		if p.TxTime - reductionTime <= consts.AUTO_REDUCTION_PERIOD {
			return p.ErrInfo("reduction_period error")
		}
	}

	if p.TxMapS["reduction_type"] == "manual" {

		// получаем кол-во обещанных сумм у разных юзеров по каждой валюте. start_time есть только у тех, у кого статус mining/repaid
		promisedAmount, err := p.DCDB.GetMap(`
					SELECT currency_id, count(user_id) as count
					FROM (
							SELECT currency_id, user_id
							FROM promised_amount
							WHERE start_time < ?  AND
										 del_block_id = 0 AND
										 del_mining_block_id = 0 AND
										 status IN ('mining', 'repaid')
							GROUP BY  user_id, currency_id
							) as t1
					GROUP BY  currency_id`, "currency_id", "count", (p.TxTime - p.Variables.Int64["min_hold_time_promise_amount"]))
		if err != nil {
			return p.ErrInfo(err)
		}
		if len(promisedAmount[p.TxMapS["currency_id"]]) == 0 {
			return p.ErrInfo("empty promised_amount")
		}
		// берем все голоса юзеров по данной валюте
		countVotes, err := p.Single("SELECT count(currency_id) as votes FROM votes_reduction WHERE time > ? AND currency_id  =  ? AND pct  =  ?", (p.TxTime - p.Variables.Int64["reduction_period"]), p.TxMapS["currency_id"], p.TxMapS["pct"]).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		if countVotes < utils.StrToInt64(promisedAmount[p.TxMapS["currency_id"]]) / 2 {
			return p.ErrInfo("incorrect count_votes")
		}
	} else if p.TxMapS["reduction_type"] == "promised_amount" {

		// и недопустимо для WOC
		if p.TxMapS["currency_id"] == "1" {
			return p.ErrInfo("WOC AUTO_REDUCTION_CASHs")
		}
		// проверим, есть ли хотябы 1000 юзеров, у которых на кошелках есть или была данная валюты
		countUsers, err := p.Single("SELECT count(user_id) FROM wallets WHERE currency_id  =  ?", p.TxMapS["currency_id"]).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		if countUsers < consts.AUTO_REDUCTION_PROMISED_AMOUNT_MIN {
			return p.ErrInfo("AUTO_REDUCTION_PROMISED_AMOUNT_MIN")
		}

		// получаем кол-во DC на кошельках
		sumWallets, err := p.Single("SELECT sum(amount) FROM wallets WHERE currency_id  =  ?", p.TxMapS["currency_id"]).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}

		// получаем кол-во TDC на обещанных суммах
		sumPromisedAmountTdc, err := p.Single("SELECT sum(tdc_amount) FROM promised_amount WHERE currency_id  =  ?", p.TxMapS["currency_id"]).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		sumWallets += sumPromisedAmountTdc;

		// получаем суммы обещанных сумм. при этом не берем те, что имеют просроченные cash_request_out
		sumPromisedAmount, err := p.Single("SELECT sum(amount) FROM promised_amount WHERE status  =  'mining' AND del_block_id  =  0 AND del_mining_block_id  =  0 AND currency_id  =  ? AND (cash_request_out_time  =  0 OR cash_request_out_time > ?)", p.TxMapS["currency_id"], (p.TxTime - p.Variables.Int64["cash_request_time"])).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		// если обещанных сумм менее чем 100% от объема DC на кошельках, то всё норм, если нет - ошибка
		if sumPromisedAmount >= sumWallets * consts.AUTO_REDUCTION_PROMISED_AMOUNT_PCT {
			return p.ErrInfo("error reduction $sum_promised_amount")
		}
	}

	return nil
}

func (p *Parser) NewReduction() (error) {
	d := 100 - (utils.BytesToFloat64(p.TxMap["pct"]) / 100)
	if utils.BytesToInt(p.TxMap["pct"]) > 0 {

		// т.к. невозможо 2 отката подряд из-за промежутка в 2 дня между reduction,
		// то можем использовать только бекап на 1 уровень назад вместо _log
		err := p.ExecSql("UPDATE wallets SET amount_backup = amount, amount = amount*(?) WHERE currency_id = ?", d, p.TxMapS["currency_id"])
		if err != nil {
			return p.ErrInfo(err)
		}

		// если бы не урезали amount, то пришлось бы делать пересчет tdc по всем, у кого есть данная валюта
		// после 87826 блока убрано amount_backup = amount, amount = amount*({$d}) т.к. теряется смысл в reduction c type=promised_amount
		err = p.ExecSql("UPDATE promised_amount SET tdc_amount_backup = tdc_amount, tdc_amount = tdc_amount*(?) WHERE currency_id = ?", d, p.TxMapS["currency_id"])
		if err != nil {
			return p.ErrInfo(err)
		}

		// все свежие cash_request_out_time отменяем
		err = p.ExecSql("UPDATE promised_amount SET cash_request_out_time_backup = cash_request_out_time, cash_request_out_time = 0 WHERE currency_id = ? AND cash_request_out_time > ?", p.TxMapS["currency_id"], (p.BlockData.Time - p.Variables.Int64["cash_request_time"]))
		if err != nil {
			return p.ErrInfo(err)
		}

		// все текущие cash_requests, т.е. по которым не прошло 2 суток
		err = p.ExecSql("UPDATE cash_requests SET del_block_id = ? WHERE currency_id = ? AND status = 'pending' AND time > ?", p.BlockData.BlockId, p.TxMapS["currency_id"], (p.BlockData.Time - p.Variables.Int64["cash_request_time"]))
		if err != nil {
			return p.ErrInfo(err)
		}

		// форeкс-ордеры
		err = p.ExecSql("UPDATE forex_orders SET amount_backup = amount, amount = amount*(?) WHERE sell_currency_id = ?", d, p.TxMapS["currency_id"])
		if err != nil {
			return p.ErrInfo(err)
		}

		// крауд-фандинг
		err = p.ExecSql("UPDATE cf_funding SET amount_backup = amount, amount = amount*(?) WHERE currency_id = ?", d, p.TxMapS["currency_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	rType := ""
	if p.TxMapS["reduction_type"] == "manual" {
		rType = "manual"
	} else {
		rType = "auto"
	}
	err := p.ExecSql("INSERT INTO reduction ( time, currency_id, type, pct, block_id ) VALUES ( ?, ?, ?, ?, ? )", p.BlockData.Time, p.TxMapS["currency_id"], rType, p.TxMapS["pct"], p.BlockData.BlockId)
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}

func (p *Parser) NewReductionRollback() (error) {
	if utils.StrToInt64(p.TxMapS["pct"])>0 {
		// крауд-фандинг
		err := p.ExecSql("UPDATE cf_funding SET amount = amount_backup, amount_backup = 0 WHERE currency_id = ?", p.TxMapS["currency_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
		// форекс-ордеры
		err = p.ExecSql("UPDATE forex_orders SET amount = amount_backup, amount_backup = 0 WHERE sell_currency_id = ?", p.TxMapS["currency_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.ExecSql("UPDATE cash_requests SET del_block_id = 0 WHERE del_block_id = ?", p.BlockData.BlockId)
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.ExecSql("UPDATE promised_amount SET cash_request_out_time = cash_request_out_time_backup WHERE currency_id = ? AND cash_request_out_time > ?", p.TxMapS["currency_id"], (p.BlockData.Time - p.Variables.Int64["cash_request_time"]))
		if err != nil {
			return p.ErrInfo(err)
		}
		// после 87826 блока убрано  amount = amount_backup т.к. теряется смысл в reduction c type=promised_amount
		err = p.ExecSql("UPDATE promised_amount SET tdc_amount = tdc_amount_backup WHERE currency_id = ?", p.TxMapS["currency_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.ExecSql("UPDATE wallets SET amount = amount_backup, amount_backup = 0 WHERE currency_id = ?", p.TxMapS["currency_id"])
		if err != nil {
			return p.ErrInfo(err)
		}

	}

	affect, err := p.ExecSqlGetAffect("DELETE FROM reduction WHERE block_id = ?", p.BlockData.BlockId)
	if err != nil {
		return p.ErrInfo(err)
	}
	p.rollbackAI("reduction", affect)

	return nil
}

func (p *Parser) NewReductionRollbackFront() error {

	return nil
}
