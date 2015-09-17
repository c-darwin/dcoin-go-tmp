package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

func FirstChangePkey() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("daemon Recovered", r)
			panic(r)
		}
	}()

	const GoroutineName = "FirstChangePkey"
	d := new(daemon)
	d.DCDB = DbConnect()
	if d.DCDB == nil {
		return
	}
	d.goRoutineName = GoroutineName
	if !d.CheckInstall(DaemonCh, AnswerDaemonCh) {
		return
	}
	d.DCDB = DbConnect()
	if d.DCDB == nil {
		return
	}

	BEGIN:
	for {
		log.Info(GoroutineName)
		MonitorDaemonCh <- []string{GoroutineName, utils.Int64ToStr(utils.Time())}

		// проверим, не нужно ли нам выйти из цикла
		if CheckDaemonsRestart() {
			break BEGIN
		}

		community, err := d.GetCommunityUsers()
		if err != nil {
			d.PrintSleep(err, 1)
			continue BEGIN
		}
		var uids []int64
		if len(community) > 0 {
			uids = community
		} else {
			myuid, err := d.GetMyUserId("")
			if err != nil {
				d.PrintSleep(err, 1)
				continue BEGIN
			}
			uids = append(uids, myuid)
		}
		log.Debug("uids %v", uids)
		var status, myPrefix string
		for _, uid := range uids {
			if len(community) > 0 {
				myPrefix = utils.Int64ToStr(uid)+"_"
			}
			status, err = d.Single(`SELECT status FROM `+myPrefix+`my_table`).String()

			log.Debug("status: %v / myPrefix: %v", status, myPrefix)
			if status == "waiting_accept_new_key" {

				// проверим, не прошла тр-ия и не сменился ли уже ключ
				userPubKeyCount, err := d.Single(`SELECT count(*) FROM `+myPrefix+`my_keys WHERE status='approved'`).Int64()
				if err != nil {
					d.PrintSleep(err, 1)
					continue BEGIN
				}
				log.Debug("userPubKey: %v", userPubKeyCount)
				if userPubKeyCount > 1 {
					err = d.ExecSql(`UPDATE `+myPrefix+`my_table SET status='user'`)
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					d.PrintSleep(err, 1)
					continue BEGIN
				}

				lastTx, err := d.GetLastTx(uid, utils.TypesToIds([]string{"ChangePrimaryKey"}), 1, "2006-02-01 15:04:05")
				if err != nil {
					d.PrintSleep(err, 1)
					continue BEGIN
				}
				log.Debug("lastTx: %v", lastTx)
				if len(lastTx) > 0 {
					if len(lastTx[0]["error"]) > 0 || utils.Time() - utils.StrToInt64(lastTx[0]["time_int"]) > 1800 {
						// генерим и шлем новую тр-ию
						err = d.SendTxChangePkey(uid);
						if err != nil {
							d.PrintSleep(err, 1)
							continue BEGIN
						}
					}
				}
			}
		}

		for i:=0; i < 60; i++ {
			utils.Sleep(1)
			// проверим, не нужно ли нам выйти из цикла
			if CheckDaemonsRestart() {
				break BEGIN
			}
		}
	}

}


