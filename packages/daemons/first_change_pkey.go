package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"bytes"
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
			if err != nil {
				d.PrintSleep(err, 1)
				continue BEGIN
			}
			log.Debug("status: %v / myPrefix: %v", status, myPrefix)
			if status == "waiting_accept_new_key" {

				// если ключ кто-то сменил
				userPublicKey, err := d.GetUserPublicKey(uid)
				if err != nil {
					d.PrintSleep(err, 1)
					continue BEGIN
				}
				myUserPublicKey, err := d.GetMyPublicKey(myPrefix)
				if err != nil {
					d.PrintSleep(err, 1)
					continue BEGIN
				}
				if !bytes.Equal(myUserPublicKey, []byte(userPublicKey)) {
					log.Debug("myUserPublicKey:%s != userPublicKey:%s", utils.BinToHex(myUserPublicKey), utils.BinToHex(userPublicKey))
					// удаляем старый ключ
					err = d.ExecSql(`DELETE FROM `+myPrefix+`my_keys`)
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					// и user_id
					q := `UPDATE `+myPrefix+`my_table SET status="my_pending", user_id=0`
					if len(community) > 0 {
						q = `DELETE FROM `+myPrefix+`my_table`
					}
					err = d.ExecSql(q)
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					if len(community) > 0 {
						err = d.ExecSql(`DELETE FROM community WHERE user_id = ?`, uid)
						if err != nil {
							d.PrintSleep(err, 1)
							continue BEGIN
						}
					}
					// и пробуем взять новый
					userId, _, err := d.GetAvailableKey()
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					if userId > 0 {
						if len(community) > 0 {
							err = d.ExecSql(`INSERT INTO community (user_id) VALUES (?)`, userId)
							if err != nil {
								d.PrintSleep(err, 1)
								continue BEGIN
							}
						}
						// генерим и шлем новую тр-ию
						err = d.SendTxChangePkey(userId);
						if err != nil {
							d.PrintSleep(err, 1)
							continue BEGIN
						}
					} else {
						// если userId == 0, значит ключей в паблике больше нет и юзеру придется искать ключ самому
						continue
					}
				}

				// проверим, не прошла ли тр-ия и не сменился ли уже ключ
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
					// также, если это пул, то надо удалить приватный ключ из базы данных, т.к. взлом пула будет означать угон всех акков
					// хранение ключа на мобильном - безопасно, хранение ключа на ПК, пока Dcoin не стал популярен, тоже допустимо
					if len(community) > 0 {
						err = d.ExecSql(`DELETE private_key FROM `+myPrefix+`my_keys`)
						if err != nil {
							d.PrintSleep(err, 1)
							continue BEGIN
						}
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


