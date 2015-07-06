package daemons

import (
	"utils"
)

func Notifications(configIni map[string]string) {

	const GoroutineName = "notifications"

	db := utils.DbConnect(configIni)
	BEGIN:
	for {
		// валюты
		/*currencyList, err := db.GetCurrencyList(false)
		if err != nil {
			db.PrintSleep(err, 60)
			continue BEGIN
		}*/
		notificationsArray := make(map[string]map[int64]map[string]string)
		userEmailSmsData := make(map[int64]map[string]string)

		myUsersIds, err := db.GetCommunityUsers()
		if err != nil {
			db.PrintSleep(err, 60)
			continue BEGIN
		}
		var community bool
		if len(myUsersIds) == 0 {
			community = false
			myUserId, err := db.GetMyUserId("")
			if err != nil {
				db.PrintSleep(err, 60)
				continue BEGIN
			}
			myUsersIds = append(myUsersIds, myUserId)
		} else {
			community = true
		}
		/*myPrefix, err:= db.GetMyPrefix()
		if err != nil {
			db.PrintSleep(err, 1)
			continue BEGIN
		}*/
		myBlockId, err:= db.GetMyBlockId()
		if err != nil {
			db.PrintSleep(err, 60)
			continue BEGIN
		}
		blockId, err:= db.GetBlockId()
		if err != nil {
			db.PrintSleep(err, 60)
			continue BEGIN
		}
		if myBlockId > blockId {
			db.PrintSleep(err, 60)
			continue BEGIN
		}
		if len(myUsersIds) > 0 {
			for i:=0; i < len(myUsersIds); i++ {
				myPrefix := ""
				if community {
					myPrefix = utils.Int64ToStr(myUsersIds[i])+"_";
				}
				myData, err := db.OneRow("SELECT * FROM "+myPrefix+"my_table").String()
				if err != nil {
					db.PrintSleep(err, 60)
					continue BEGIN
				}
				// на пуле шлем уведомления только майнерам
				if community && myData["miner_id"] == "0" {
					continue
				}
				myNotifications, err := db.GetAll("SELECT * FROM "+myPrefix+"my_notifications", -1)
				if err != nil {
					db.PrintSleep(err, 60)
					continue BEGIN
				}
				for _, data := range myNotifications {
					notificationsArray[data["name"]] = make(map[int64]map[string]string)
					notificationsArray[data["name"]][myUsersIds[i]] = map[string]string{"email": data["email"], "sms": data["sms"]}
					userEmailSmsData[myUsersIds[i]] = myData
				}
			}
		}

		poolAdminUserId, err := db.GetPoolAdminUserId()
		if err != nil {
			db.PrintSleep(err, 60)
			continue BEGIN
		}
		subj := "DCoin notifications"
		for name, notificationInfo := range notificationsArray {
			switch name {
			case "admin_messages":
				data, err := db.OneRow("SELECT id, message FROM alert_messages WHERE notification  =  0").String()
				if err != nil {
					db.PrintSleep(err, 60)
					continue BEGIN
				}
				if len(data) > 0 {
					err = db.ExecSql("UPDATE alert_messages SET notification = 1 WHERE id = ?", data["id"])
					if err != nil {
						db.PrintSleep(err, 60)
						continue BEGIN
					}
					if myBlockId > blockId {
						db.PrintSleep(err, 60)
						continue BEGIN
					}
					for userId, emailSms := range notificationInfo {
						if emailSms["email"] == "1" {
							err = db.SendMail("From Admin: "+data["message"], subj, userEmailSmsData[userId]["email"], userEmailSmsData[userId], community, poolAdminUserId)
							if err != nil {
								db.PrintSleep(err, 60)
								continue BEGIN
							}
						}
						if emailSms["sms"] == "1" {
							_, err = utils.SendSms(userEmailSmsData[userId]["sms_http_get_request"], userEmailSmsData[userId]["text"])
							if err != nil {
								db.PrintSleep(err, 60)
								continue BEGIN
							}
						}
					}
				}
			case "incoming_cash_requests":

			}
		}
	}
}
