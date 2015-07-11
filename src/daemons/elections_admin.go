package daemons

import (
	"utils"
	"encoding/json"
	"fmt"
	"dcparser"
)

func elections_admin(configIni map[string]string) string {

	const GoroutineName = "elections_admin"

	db := utils.DbConnect(configIni)
	db.GoroutineName = GoroutineName
	BEGIN:
	for {

		err := db.DbLock()
		if err != nil {
			db.PrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		blockId, err := db.GetBlockId()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		if blockId == 0 {
			db.UnlockPrintSleep(utils.ErrInfo("blockId == 0"), 60)
			continue BEGIN
		}

		_, _, myMinerId, _, _, _, err := db.TestBlock();
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		// а майнер ли я ?
		if myMinerId == 0 {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		variables, err := db.GetAllVariables()
		curTime := utils.Time()

		// проверим, прошло ли 2 недели с момента последнего обновления
		adminTime, err := db.Single("SELECT time FROM admin").Int64()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		if curTime - adminTime <= variables.Int64["new_pct_period"] {
			db.UnlockPrintSleep(utils.ErrInfo("14 day error"), 60)
			continue BEGIN
		}

		// сколько всего майнеров
		countMiners, err := db.Single("SELECT count(miner_id) FROM miners WHERE active  =  1").Int64()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		if countMiners < 1000 {
			db.UnlockPrintSleep(utils.ErrInfo("countMiners < 1000"), 60)
			continue BEGIN
		}

		// берем все голоса
		var newAdmin int64
		data, err := db.GetMap(`
				SELECT	 admin_user_id,
							  count(user_id) as votes
				FROM votes_admin
				WHERE time > ?
				GROUP BY  admin_user_id
				`, curTime - variables.Int64["new_pct_period"], "admin_user_id", "votes")
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		for admin_user_id, votes := range data {
			// если более 50% майнеров проголосовали
			if utils.StrToInt64(votes) > countMiners/2 {
				newAdmin = utils.StrToInt64(admin_user_id)
			}
		}
		if newAdmin == 0 {
			db.UnlockPrintSleep(utils.ErrInfo("newAdmin == 0"), 60)
			continue BEGIN
		}

		_, myUserId, _, _, _, _, err := db.TestBlock();
		forSign := fmt.Sprintf("%v,%v,%v,%v", utils.TypeInt("NewAdmin"), curTime, myUserId, newAdmin)
		binSign, err := db.GetBinSign(forSign, myUserId)
		if err!= nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		data := utils.DecToBin(utils.TypeInt("NewAdmin"), 1)
		data = append(data, utils.DecToBin(curTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(utils.Int64ToByte(myUserId))...)
		data = append(data, utils.EncodeLengthPlusData(utils.Int64ToByte(newAdmin))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(binSign))...)

		err = db.InsertReplaceTxInQueue(data)
		if err!= nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}

		p := new(dcparser.Parser)
		err = db.TxParser(data, utils.HexToBin(utils.Md5(data)), true)
		if err != nil {
			db.PrintSleep(err, 60)
			continue BEGIN
		}

		db.DbUnlock()
		utils.Sleep(60)
	}
	return ""
}


