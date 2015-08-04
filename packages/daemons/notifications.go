package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"fmt"
	"math"
)

func Notifications() {

	const GoroutineName = "Notifications"
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
		// валюты
		currencyList, err := d.GetCurrencyList(false)
		if err != nil {
			d.PrintSleep(err, 1)
			continue BEGIN
		}
		notificationsArray := make(map[string]map[int64]map[string]string)
		userEmailSmsData := make(map[int64]map[string]string)

		myUsersIds, err := d.GetCommunityUsers()
		if err != nil {
			d.PrintSleep(err, 1)
			continue BEGIN
		}
		var community bool
		if len(myUsersIds) == 0 {
			community = false
			myUserId, err := d.GetMyUserId("")
			if err != nil {
				d.PrintSleep(err, 1)
				continue BEGIN
			}
			myUsersIds = append(myUsersIds, myUserId)
		} else {
			community = true
		}
		/*myPrefix, err:= d.GetMyPrefix()
		if err != nil {
			d.PrintSleep(err, 1)
			continue BEGIN
		}*/
		myBlockId, err:= d.GetMyBlockId()
		if err != nil {
			d.PrintSleep(err, 1)
			continue BEGIN
		}
		blockId, err:= d.GetBlockId()
		if err != nil {
			d.PrintSleep(err, 1)
			continue BEGIN
		}
		if myBlockId > blockId {
			d.PrintSleep(err, 1)
			continue BEGIN
		}
		if len(myUsersIds) > 0 {
			for i:=0; i < len(myUsersIds); i++ {
				myPrefix := ""
				if community {
					myPrefix = utils.Int64ToStr(myUsersIds[i])+"_";
				}
				myData, err := d.OneRow("SELECT * FROM "+myPrefix+"my_table").String()
				if err != nil {
					d.PrintSleep(err, 1)
					continue BEGIN
				}
				// на пуле шлем уведомления только майнерам
				if community && myData["miner_id"] == "0" {
					continue
				}
				myNotifications, err := d.GetAll("SELECT * FROM "+myPrefix+"my_notifications", -1)
				if err != nil {
					d.PrintSleep(err, 1)
					continue BEGIN
				}
				for _, data := range myNotifications {
					notificationsArray[data["name"]] = make(map[int64]map[string]string)
					notificationsArray[data["name"]][myUsersIds[i]] = map[string]string{"email": data["email"], "sms": data["sms"]}
					userEmailSmsData[myUsersIds[i]] = myData
				}
			}
		}

		poolAdminUserId, err := d.GetPoolAdminUserId()
		if err != nil {
			d.PrintSleep(err, 1)
			continue BEGIN
		}
		subj := "DCoin notifications"
		for name, notificationInfo := range notificationsArray {
			switch name {
			case "admin_messages":
				data, err := d.OneRow("SELECT id, message FROM alert_messages WHERE notification  =  0").String()
				if err != nil {
					d.PrintSleep(err, 1)
					continue BEGIN
				}
				if len(data) > 0 {
					err = d.ExecSql("UPDATE alert_messages SET notification = 1 WHERE id = ?", data["id"])
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					if myBlockId > blockId {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					for userId, emailSms := range notificationInfo {
						if emailSms["email"] == "1" {
							err = d.SendMail("From Admin: "+data["message"], subj, userEmailSmsData[userId]["email"], userEmailSmsData[userId], community, poolAdminUserId)
							if err != nil {
								d.PrintSleep(err, 1)
								continue BEGIN
							}
						}
						if emailSms["sms"] == "1" {
							_, err = utils.SendSms(userEmailSmsData[userId]["sms_http_get_request"], userEmailSmsData[userId]["text"])
							if err != nil {
								d.PrintSleep(err, 1)
								continue BEGIN
							}
						}
					}
				}
			case "incoming_cash_requests":
				for i:=0; i<len(myUsersIds); i++ {
					myPrefix := ""
					if community {
						myPrefix = utils.Int64ToStr(myUsersIds[i])+"_"
					}
					userId := myUsersIds[i]
					data, err := d.OneRow("SELECT id, amount, currency_id FROM "+myPrefix+"my_cash_requests WHERE to_user_id  =  ? AND notification  =  0 AND status  =  'pending'", myPrefix, userId).String()
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					if len(data) > 0 {
						text := `You"ve got the request for `+data["amount"]+` `+currencyList[utils.StrToInt64(data["currencyId"])]+`. It has to be repaid within the next 48 hours.`
						if notificationsArray[name][userId]["email"] == "1" {
							err = d.SendMail(text, subj, userEmailSmsData[userId]["email"], userEmailSmsData[userId], community, poolAdminUserId)
							if err != nil {
								d.PrintSleep(err, 1)
								continue BEGIN
							}
						}
						if notificationsArray[name][userId]["sms"] == "1" {
							utils.SendSms(userEmailSmsData[userId]["sms_http_get_request"], text)
						}
						err = d.ExecSql("UPDATE "+myPrefix+"my_cash_requests SET notification = 1 WHERE id = ?", data["id"])
						if err != nil {
							d.PrintSleep(err, 1)
							continue BEGIN
						}
					}
				}
			case "change_in_status":

				for i:=0; i<len(myUsersIds); i++ {
					myPrefix := ""
					if community {
						myPrefix = utils.Int64ToStr(myUsersIds[i])+"_"
					}
					userId := myUsersIds[i]
					status, err := d.Single("SELECT status FROM "+myPrefix+"my_table WHERE notification_status = 0").String()
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					if len(status) > 0 {
						text := `New status: `+status
						if notificationsArray[name][userId]["email"] == "1" {
							err = d.SendMail(text, subj, userEmailSmsData[userId]["email"], userEmailSmsData[userId], community, poolAdminUserId)
							if err != nil {
								d.PrintSleep(err, 1)
								continue BEGIN
							}
						}
						if notificationsArray[name][userId]["sms"] == "1" {
							utils.SendSms(userEmailSmsData[userId]["sms_http_get_request"], text)
						}
						err = d.ExecSql("UPDATE "+myPrefix+"my_table SET notification_status = 1")
						if err != nil {
							d.PrintSleep(err, 1)
							continue BEGIN
						}
					}
				}
			case "dc_came_from" :

				for i:=0; i<len(myUsersIds); i++ {
					myPrefix := ""
					if community {
						myPrefix = utils.Int64ToStr(myUsersIds[i])+"_"
					}
					userId := myUsersIds[i]
					// откатываем наши блоки до начала вилки
					myDcTransactions, err := d.GetAll(`
							SELECT  id,
							               amount,
										 currency_id,
										 comment_status,
										 comment
							FROM `+myPrefix+`my_dc_transactions
							WHERE to_user_id = ? AND
									 	notification = 0 AND
									 	status = 'approved'` ,- 1, userId)
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					for _, data := range myDcTransactions {
						comment := ""
						if data["comment_status"] == "decrypted" {
							comment = `<br><span style="font-size:16px">`+data["comment"]+`</span>`
						}
						text := `You've got `+data["amount"]+` D`+currencyList[utils.StrToInt64(data["currency_id"])]+` `+comment
						if notificationsArray[name][userId]["email"] == "1" {
							err = d.SendMail(text, subj, userEmailSmsData[userId]["email"], userEmailSmsData[userId], community, poolAdminUserId)
							if err != nil {
								d.PrintSleep(err, 1)
								continue BEGIN
							}
						}
						if notificationsArray[name][userId]["sms"] == "1" {
							utils.SendSms(userEmailSmsData[userId]["sms_http_get_request"], text)
						}
						err = d.ExecSql("UPDATE "+myPrefix+"my_dc_transactions SET notification = 1 WHERE id = ?", data["id"])
						if err != nil {
							d.PrintSleep(err, 1)
							continue BEGIN
						}
					}
				}
			case "dc_sent" :

				for i:=0; i<len(myUsersIds); i++ {
					myPrefix := ""
					if community {
						myPrefix = utils.Int64ToStr(myUsersIds[i])+"_"
					}
					userId := myUsersIds[i]
					// откатываем наши блоки до начала вилки
					myDcTransactions, err := d.GetAll(`
							SELECT id,
									    amount,
									    currency_id
							FROM `+myPrefix+`my_dc_transactions
							WHERE to_user_id !=  ? AND
										 notification = 0 AND
										 status = 'approved'`, -1, userId)
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					for _, data := range myDcTransactions {
						text := `Debiting `+data["amount"]+` D`+currencyList[utils.StrToInt64(data["currency_id"])]
						if notificationsArray[name][userId]["email"] == "1" {
							err = d.SendMail(text, subj, userEmailSmsData[userId]["email"], userEmailSmsData[userId], community, poolAdminUserId)
							if err != nil {
								d.PrintSleep(err, 1)
								continue BEGIN
							}
						}
						if notificationsArray[name][userId]["sms"] == "1" {
							utils.SendSms(userEmailSmsData[userId]["sms_http_get_request"], text)
						}
						err = d.ExecSql("UPDATE "+myPrefix+"my_dc_transactions SET notification = 1 WHERE id = ?", data["id"])
						if err != nil {
							d.PrintSleep(err, 1)
							continue BEGIN
						}
					}
				}
			case "update_primary_key" :

				for i:=0; i<len(myUsersIds); i++ {
					myPrefix := ""
					if community {
						myPrefix = utils.Int64ToStr(myUsersIds[i])+"_"
					}
					userId := myUsersIds[i]
					data, err := d.OneRow("SELECT id FROM "+myPrefix+"my_keys WHERE notification = 0 AND status = 'approved'").String()
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					if len(data) > 0 {
						text := `Update primary key`
						if notificationsArray[name][userId]["email"] == "1" {
							err = d.SendMail(text, subj, userEmailSmsData[userId]["email"], userEmailSmsData[userId], community, poolAdminUserId)
							if err != nil {
								d.PrintSleep(err, 1)
								continue BEGIN
							}
						}
						if notificationsArray[name][userId]["sms"] == "1" {
							utils.SendSms(userEmailSmsData[userId]["sms_http_get_request"], text)
						}
						err = d.ExecSql("UPDATE "+myPrefix+"my_keys SET notification = 1 WHERE id = ?", data["id"])
						if err != nil {
							d.PrintSleep(err, 1)
							continue BEGIN
						}
					}
				}
			case "update_email":

				for i:=0; i<len(myUsersIds); i++ {
					myPrefix := ""
					if community {
						myPrefix = utils.Int64ToStr(myUsersIds[i])+"_"
					}
					userId := myUsersIds[i]
					myNewEmail, err := d.Single("SELECT status FROM "+myPrefix+"my_table WHERE notification_email  =  0").String()
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					if len(myNewEmail) > 0 {
						text := `New email: `+myNewEmail
						if notificationsArray[name][userId]["email"] == "1" {
							err = d.SendMail(text, subj, userEmailSmsData[userId]["email"], userEmailSmsData[userId], community, poolAdminUserId)
							if err != nil {
								d.PrintSleep(err, 1)
								continue BEGIN
							}
						}
						if notificationsArray[name][userId]["sms"] == "1" {
							utils.SendSms(userEmailSmsData[userId]["sms_http_get_request"], text)
						}
						err = d.ExecSql("UPDATE "+myPrefix+"my_table SET notification_email = 1")
						if err != nil {
							d.PrintSleep(err, 1)
							continue BEGIN
						}
					}
				}
			case "update_sms_request" :

				for i:=0; i<len(myUsersIds); i++ {
					myPrefix := ""
					if community {
						myPrefix = utils.Int64ToStr(myUsersIds[i])+"_"
					}
					userId := myUsersIds[i]
					smsHttpGetRequest, err := d.Single("SELECT sms_http_get_request FROM "+myPrefix+"my_table WHERE notification_sms_http_get_request  =  0").String()
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					if len(smsHttpGetRequest) > 0 {
						text := `New sms_http_get_request `+smsHttpGetRequest
						if notificationsArray[name][userId]["email"] == "1" {
							err = d.SendMail(text, subj, userEmailSmsData[userId]["email"], userEmailSmsData[userId], community, poolAdminUserId)
							if err != nil {
								d.PrintSleep(err, 1)
								continue BEGIN
							}
						}
						if notificationsArray[name][userId]["sms"] == "1" {
							utils.SendSms(userEmailSmsData[userId]["sms_http_get_request"], text)
						}
						err = d.ExecSql("UPDATE "+myPrefix+"my_table SET notification_sms_http_get_request = 1")
						if err != nil {
							d.PrintSleep(err, 1)
							continue BEGIN
						}
					}
				}
			case "voting_results":
				myDcTransactions, err := d.GetAll(`
						SELECT  id,
									 currency_id,
									 miner,
									 user,
									 block_id
						FROM pct
						WHERE notification = 0`, -1)
				if err != nil {
					d.PrintSleep(err, 1)
					continue BEGIN
				}
				text := ""
				pctUpd := false
				for _, data := range myDcTransactions {
					pctUpd = true
					text += fmt.Sprintf("New pct %v! miner: %v %/block, user: %v %/block", currencyList[utils.StrToInt64(data["currency_id"])], ((math.Pow(1+utils.StrToFloat64(data["miner"]), 120)-1)*100), ((math.Pow(1+utils.StrToFloat64(data["user"]), 120)-1)*100))
				}
				if pctUpd {
					err = d.ExecSql("UPDATE pct SET notification = 1 WHERE notification = 0")
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
				}
				if myBlockId > blockId {
					d.PrintSleep(err, 1)
					continue BEGIN
				}

				// шлется что-то не то, потом поправлю, пока отключил
				if community {
					d.PrintSleep(err, 1)
					continue BEGIN
				}

				if len(text) > 0 {
					for userId, emailSms := range notificationInfo {
						if emailSms["email"] == "1" {
							err = d.SendMail(text, subj, userEmailSmsData[userId]["email"], userEmailSmsData[userId], community, poolAdminUserId)
							if err != nil {
								d.PrintSleep(err, 1)
								continue BEGIN
							}
						}
						if emailSms["sms"] == "1" {
							utils.SendSms(userEmailSmsData[userId]["sms_http_get_request"], text)
						}
					}
				}
			case "voting_time": // Прошло 2 недели с момента Вашего голосования за %

				for i:=0; i<len(myUsersIds); i++ {
					myPrefix := ""
					if community {
						myPrefix = utils.Int64ToStr(myUsersIds[i])+"_"
					}
					userId := myUsersIds[i]
					lastVoting, err := d.Single("SELECT last_voting FROM "+myPrefix+"my_complex_votes WHERE notification  =  0").Int64()
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					if lastVoting > 0 && utils.Time()-lastVoting > 86400*14 {
						text := "It's 2 weeks from the moment you voted."
						if notificationsArray[name][userId]["email"] == "1" {
							err = d.SendMail(text, subj, userEmailSmsData[userId]["email"], userEmailSmsData[userId], community, poolAdminUserId)
							if err != nil {
								d.PrintSleep(err, 1)
								continue BEGIN
							}
						}
						if notificationsArray[name][userId]["sms"] == "1" {
							utils.SendSms(userEmailSmsData[userId]["sms_http_get_request"], text)
						}
						err = d.ExecSql("UPDATE "+myPrefix+"my_complex_votes SET notification = 1")
						if err != nil {
							d.PrintSleep(err, 1)
							continue BEGIN
						}
					}
				}

			case "new_version":

				newVersion, err := d.Single("SELECT version FROM new_version WHERE notification  =  0 AND alert  =  1").String()
				if err != nil {
					d.PrintSleep(err, 1)
					continue BEGIN
				}

				err = d.ExecSql("UPDATE new_version SET notification = 1 WHERE version = ?", newVersion)
				if err != nil {
					d.PrintSleep(err, 1)
					continue BEGIN
				}

				if myBlockId > blockId {
					d.PrintSleep(err, 1)
					continue BEGIN
				}

				// в пуле это лишняя инфа
				if community{
					d.PrintSleep(err, 1)
					continue BEGIN
				}
				if len(newVersion) > 0 {
					for userId, emailSms := range notificationInfo {
						text := "New version: "+newVersion
						if emailSms["email"] == "1" {
							err = d.SendMail(text, subj, userEmailSmsData[userId]["email"], userEmailSmsData[userId], community, poolAdminUserId)
							if err != nil {
								d.PrintSleep(err, 1)
								continue BEGIN
							}
						}
						if emailSms["sms"] == "1" {
							utils.SendSms(userEmailSmsData[userId]["sms_http_get_request"], text)
						}
					}
				}

			case "node_time": // Расхождение времени сервера более чем на 5 сек

				var adminUserId int64
				// если работаем в режиме пула, то нужно слать инфу админу пула
				if community {
					adminUserId = poolAdminUserId
				} else {
					// проверим, нода ли мы
					my_table, err := d.OneRow("SELECT user_id, miner_id FROM my_table").Int64()
					if err != nil {
						d.PrintSleep(err, 1)
						continue BEGIN
					}
					if my_table["miner_id"] == 0 {
						d.PrintSleep(err, 1)
						continue BEGIN
					} else {
						adminUserId = my_table["user_id"]
					}
					emailSms := notificationInfo[adminUserId];
					myData := userEmailSmsData[adminUserId];
					if len(myData) > 0 {
						networkTime, err := utils.GetNetworkTime()
						if err != nil {
							d.PrintSleep(err, 1)
							continue BEGIN
						}
						diff := int64(math.Abs(float64(utils.Time() - networkTime.Unix())))
						text := ""
						if diff > 5 {
							text = "Divergence time "+utils.Int64ToStr(diff)+" sec"
						}
						if emailSms["email"] == "1" {
							err = d.SendMail(text, subj, myData["email"], myData, community, poolAdminUserId)
							if err != nil {
								d.PrintSleep(err, 1)
								continue BEGIN
							}
						}
						if emailSms["sms"] == "1" {
							utils.SendSms(myData["sms_http_get_request"], text)
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
