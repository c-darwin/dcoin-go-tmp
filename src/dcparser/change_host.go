package dcparser

import (
	"fmt"
	"utils"
	//"log"
	//"encoding/json"
	//"regexp"
	//"math"
	//"strings"
//	"os"
//	"time"
	//"strings"
	//"bytes"
	//"consts"
//	"math"
//	"database/sql"
//	"bytes"
)

func (p *Parser) ChangeHostInit() (error) {

	fields := []map[string]string {{"host":"string"}, {"sign":"bytes"}}
	err := p.GetTxMaps(fields);
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}


func (p *Parser) ChangeHostFront() (error) {

	err := p.generalCheck()
	if err != nil {
		return p.ErrInfo(err)
	}

	verifyData := map[string]string {"host":"host"}
	err = p.CheckInputData(verifyData)
	if err != nil {
		return p.ErrInfo(err)
	}

	// является ли данный юзер майнером
	err = p.checkMiner(p.TxUserID)
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

	var CheckSignResult bool
	forSign := fmt.Sprintf("%s,%s,%s,%s", p.TxMap["type"], p.TxMap["time"], p.TxMap["user_id"], p.TxMap["host"])
	if p.BlockData!=nil && p.BlockData.BlockId <= 240240 {
		CheckSignResult, err = utils.CheckSign(p.PublicKeys, forSign, p.TxMap["sign"], false);
	} else {
		CheckSignResult, err = utils.CheckSign([][]byte{p.nodePublicKey}, forSign, p.TxMap["sign"], true);
	}
	if err != nil {
		return p.ErrInfo(err)
	}
	if !CheckSignResult {
		return p.ErrInfo("incorrect sign")
	}

	err = p.limitRequest(p.Variables.Int64["limit_change_host"], "change_host", p.Variables.Int64["limit_change_host_period"])
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}

func (p *Parser) ChangeHost() (error) {
	err := p.selectiveLoggingAndUpd([]string{"host"}, []interface {}{p.TxMaps.String["host"]}, "miners_data", []string{"user_id"}, []string{utils.Int64ToStr(p.TxUserID)})
	if err != nil {
		return p.ErrInfo(err)
	}
	myUserId, myBlockId, myPrefix, _ , err := p.GetMyUserId(p.TxUserID)
	if err != nil {
		return p.ErrInfo(err)
	}
	if p.TxUserID == myUserId && myBlockId <= p.BlockData.BlockId {
		// обновим статус в нашей локальной табле.
		err = p.ExecSql("UPDATE "+myPrefix+"my_table SET host_status = 'approved'")
		if err != nil {
			return p.ErrInfo(err)
		}
	}

	return nil
}

func (p *Parser) ChangeHostRollback() (error) {
	err := p.selectiveRollback([]string{"host"}, "miners_data", "user_id="+utils.Int64ToStr(p.TxUserID), false)
	if err != nil {
		return p.ErrInfo(err)
	}
	myUserId, _, myPrefix, _ , err := p.GetMyUserId(p.TxUserID)
	if err != nil {
		return p.ErrInfo(err)
	}
	if p.TxUserID == myUserId {
		// обновим статус в нашей локальной табле.
		err = p.ExecSql("UPDATE "+myPrefix+"my_table SET host_status = 'my_pending'")
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	return nil
}

func (p *Parser) ChangeHostRollbackFront() error {
	return p.limitRequestsRollback("change_host")
}
