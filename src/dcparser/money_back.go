package dcparser

import (
	"fmt"
	"utils"
)

func (p *Parser) TmpInit() (error) {

	fields := []map[string]string {{"name":"int64"}, {"sign":"bytes"}}
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

	verifyData := map[string]string {"name":"bigint", "name2":"bigint"}
	err = p.CheckInputData(verifyData)
	if err != nil {
		return p.ErrInfo(err)
	}

	// является ли данный юзер майнером
	err = p.checkMiner(p.TxUserID)
	if err != nil {
		return p.ErrInfo(err)
	}

	// У юзера не должно быть cash_requests с pending
	err = p.CheckCashRequests(p.TxUserID)
	if err != nil {
		return p.ErrInfo(err)
	}

	// нодовский ключ
	nodePublicKey, err := p.GetNodePublicKey(p.TxUserID)
	if err != nil {
		return p.ErrInfo(err)
	}
	if len(nodePublicKey) == 0 {
		return p.ErrInfo("incorrect user_id")
	}

	forSign := fmt.Sprintf("%s,%s,%s,%s", p.TxMap["type"], p.TxMap["time"], p.TxMap["user_id"], p.TxMap["name"])
	CheckSignResult, err := utils.CheckSign(p.PublicKeys, forSign, p.TxMap["sign"], false);
	if err != nil {
		return p.ErrInfo(err)
	}
	if !CheckSignResult {
		return p.ErrInfo("incorrect sign")
	}

	err = p.limitRequest(p.Variables.Int64["limit_name"], "name", p.Variables.Int64["limit_name_period"])
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}

func (p *Parser) Tmp() (error) {

	err := p.selectiveLoggingAndUpd([]string{"host"}, []string{p.TxMaps.String["host"]}, "miners_data", []string{"user_id"}, []string{utils.Int64ToStr(p.TxUserID)})
	if err != nil {
		return p.ErrInfo(err)
	}

	myUserId, myBlockId, myPrefix, _ , err := p.GetMyUserId(p.TxUserID)
	if err != nil {
		return p.ErrInfo(err)
	}
	if p.TxUserID == myUserId && myBlockId <= p.BlockData.BlockId {
		fmt.Println(myPrefix)
	}

	return nil
}

func (p *Parser) TmpRollback() (error) {
	err := p.selectiveRollback([]string{"public_key_0","public_key_1","public_key_2","change_key_close"}, "users", "user_id="+utils.Int64ToStr(p.TxMaps.Int64["for_user_id"]), false)

	return nil
}

func (p *Parser) TmpRollbackFront() error {
	return p.limitRequestsRollback("name")
}
