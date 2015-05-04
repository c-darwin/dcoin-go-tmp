package dcparser

import (
	//"fmt"
	"utils"
)

func (p *Parser) NewHolidaysInit() (error) {
	fields := []string {"start_time", "end_time", "sign"}
	TxMap := make(map[string]string)
	TxMap, err := p.GetTxMap(fields);
	//fmt.Println("0 TxMapp", TxMap)
	p.TxMap = TxMap;
	if err != nil {
		return err
	}
	//fmt.Println("1 TxMap", p.TxMap)
	return nil
}

func (p *Parser) NewHolidays() (error) {
	//fmt.Println("TxMap", p.TxMap)
	//var myUserIds []int64;
	_, err := p.ExecSql(`INSERT INTO holidays (user_id,	start_time,end_time) VALUES ($1, $2, $3)`,
		p.TxMap["user_id"], p.TxMap["start_time"], p.TxMap["end_time"])
	if err != nil {
		return err
	}
	// проверим, не наш ли это user_id
	myUserId, myBlockId, myPrefix, _ := p.GetMyUserId(utils.StrToInt64(p.TxMap["user_id"]))
	//fmt.Println(myUserIds)
	if utils.StrToInt64(p.TxMap["user_id"]) == myUserId && myBlockId <= utils.StrToInt64(p.BlockData["block_id"]) {
		// обновим статус в нашей локальной табле
		_, err := p.ExecSql("DELETE FROM "+myPrefix+"my_holidays WHERE start_time=$1 AND end_time=$2", p.TxMap["start_time"], p.TxMap["end_time"])
		if err != nil {
			return err
		}
	}
	return nil
}
func (p *Parser) NewHolidaysRollback() (error) {
	//fmt.Println(p.TxMap)
	_, err := p.ExecSql("DELETE FROM holidays WHERE user_id=$1 AND start_time=$2 AND end_time=$3", p.TxMap["user_id"], p.TxMap["start_time"], p.TxMap["end_time"])

	if err != nil {
		return utils.ErrInfo(err)
	}

	err = p.rollbackAI("holidays", 1)
	if err != nil {
		return utils.ErrInfo(err)
	}
	return err
}
