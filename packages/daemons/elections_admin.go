package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
)

func ElectionsAdmin() {

	const GoroutineName = "ElectionsAdmin"
	d := new(daemon)
	d.DCDB = DbConnect()
	if d.DCDB == nil {
		return
	}
	d.goRoutineName = GoroutineName
	if !d.CheckInstall(DaemonCh, AnswerDaemonCh) {
		return
	}

	BEGIN:
	for {
		log.Info(GoroutineName)
		// проверим, не нужно ли нам выйти из цикла
		if CheckDaemonsRestart() {
			break BEGIN
		}

		err, restart := d.dbLock()
		if restart {
			break BEGIN
		}
		if err != nil {
			d.PrintSleep(err, 1)
			continue BEGIN
		}

		blockId, err := d.GetBlockId()
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		if blockId == 0 {
			d.unlockPrintSleep(utils.ErrInfo("blockId == 0"), 1)
			continue BEGIN
		}

		_, _, myMinerId, _, _, _, err := d.TestBlock();
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		// а майнер ли я ?
		if myMinerId == 0 {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		variables, err := d.GetAllVariables()
		curTime := utils.Time()

		// проверим, прошло ли 2 недели с момента последнего обновления
		adminTime, err := d.Single("SELECT time FROM admin").Int64()
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		if curTime - adminTime <= variables.Int64["new_pct_period"] {
			d.unlockPrintSleep(utils.ErrInfo("14 day error"), 1)
			continue BEGIN
		}

		// сколько всего майнеров
		countMiners, err := d.Single("SELECT count(miner_id) FROM miners WHERE active  =  1").Int64()
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		if countMiners < 1000 {
			d.unlockPrintSleep(utils.ErrInfo("countMiners < 1000"), 1)
			continue BEGIN
		}

		// берем все голоса
		var newAdmin int64
		votes_admin, err := d.GetMap(`
				SELECT	 admin_user_id,
							  count(user_id) as votes
				FROM votes_admin
				WHERE time > ?
				GROUP BY  admin_user_id
				`, "admin_user_id", "votes", curTime - variables.Int64["new_pct_period"])
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		for admin_user_id, votes := range votes_admin {
			// если более 50% майнеров проголосовали
			if utils.StrToInt64(votes) > countMiners/2 {
				newAdmin = utils.StrToInt64(admin_user_id)
			}
		}
		if newAdmin == 0 {
			d.unlockPrintSleep(utils.ErrInfo("newAdmin == 0"), 1)
			continue BEGIN
		}

		_, myUserId, _, _, _, _, err := d.TestBlock();
		forSign := fmt.Sprintf("%v,%v,%v,%v", utils.TypeInt("NewAdmin"), curTime, myUserId, newAdmin)
		binSign, err := d.GetBinSign(forSign, myUserId)
		if err!= nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		data := utils.DecToBin(utils.TypeInt("NewAdmin"), 1)
		data = append(data, utils.DecToBin(curTime, 4)...)
		data = append(data, utils.EncodeLengthPlusData(utils.Int64ToByte(myUserId))...)
		data = append(data, utils.EncodeLengthPlusData(utils.Int64ToByte(newAdmin))...)
		data = append(data, utils.EncodeLengthPlusData([]byte(binSign))...)

		err = d.InsertReplaceTxInQueue(data)
		if err!= nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		p := new(dcparser.Parser)
		err = p.TxParser(utils.HexToBin(utils.Md5(data)), data, true)
		if err != nil {
			d.unlockPrintSleep(err, 1)
			continue BEGIN
		}

		d.dbUnlock()
		for i:=0; i < 60; i++ {
			utils.Sleep(1)
			// проверим, не нужно ли нам выйти из цикла
			if CheckDaemonsRestart() {
				break BEGIN
			}
		}
	}

}


