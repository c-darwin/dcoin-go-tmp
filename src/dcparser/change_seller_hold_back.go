package dcparser

import (
	"fmt"
	"utils"
	"consts"
)

//  продавец меняет % и кол-во дней для новых сделок
func (p *Parser) TmpInit() (error) {

	fields := []map[string]string {{"arbitration_days_refund":"int64"}, {"hold_back_pct":"money"}, {"sign":"bytes"}}
	err := p.GetTxMaps(fields);
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}


func (p *Parser) TmpFront() (error) {

	err := p.generalCheck()
	if err != nil {
		return p.ErrInfo(err)
	}

	verifyData := map[string]string {"arbitration_days_refund":"smallint", "hold_back_pct":"pct"}
	err = p.CheckInputData(verifyData)
	if err != nil {
		return p.ErrInfo(err)
	}

	if p.TxMaps.Money["hold_back_pct"] < 0.01 || p.TxMaps.Money["hold_back_pct"] > 100 {
		return p.ErrInfo("incorrect hold_back_pct")
	}

	forSign := fmt.Sprintf("%s,%s,%s,%s,%s", p.TxMap["type"], p.TxMap["time"], p.TxMap["user_id"], p.TxMap["arbitration_days_refund"], p.TxMap["hold_back_pct"])
	CheckSignResult, err := utils.CheckSign(p.PublicKeys, forSign, p.TxMap["sign"], false);
	if err != nil {
		return p.ErrInfo(err)
	}
	if !CheckSignResult {
		return p.ErrInfo("incorrect sign")
	}

	err = p.limitRequest(consts.LIMIT_CHANGE_SELLER_HOLD_BACK, "change_seller_hold_back", consts.LIMIT_CHANGE_SELLER_HOLD_BACK_PERIOD)
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}

func (p *Parser) Tmp() (error) {
	return p.selectiveLoggingAndUpd([]string{"arbitration_days_refund", "seller_hold_back_pct"}, []string{p.TxMaps.Int64["arbitration_days_refund"], p.TxMaps.Money["hold_back_pct"]}, "users", []string{"user_id"}, []string{utils.Int64ToStr(p.TxUserID)})
}

func (p *Parser) TmpRollback() (error) {
	return p.selectiveRollback([]string{"arbitration_days_refund","seller_hold_back_pct"}, "users", "user_id="+utils.Int64ToStr(p.TxUserID), false)
}

func (p *Parser) TmpRollbackFront() error {
	return p.limitRequestsRollback("change_seller_hold_back")
}
