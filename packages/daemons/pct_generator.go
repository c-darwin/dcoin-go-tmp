package daemons

import (
	"dcoin/packages/utils"
	"log"
	"encoding/json"
	"fmt"
	"dcoin/packages/dcparser"
)

/*
 * Каждые 2 недели собираем инфу о голосах за % и создаем тр-ию, которая
 * попадет в DC сеть только, если мы окажемся генератором блока
 * */
func PctGenerator(configIni map[string]string) string {

	const GoroutineName = "PctGenerator"
	db := utils.DbConnect(configIni)
	db.GoroutineName = GoroutineName
	db.CheckInstall()
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

		// проверим, прошло ли 2 недели с момента последнего обновления pct
		pctTime, err := db.Single("SELECT max(time) FROM pct").Int64()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 60)
			continue BEGIN
		}
		if curTime - pctTime > variables.Int64["new_pct_period"] {

			// берем все голоса miner_pct
			pctVotes := make(map[int64]map[string]map[string]int64)
			rows, err := db.Query("SELECT currency_id, pct, count(user_id) as votes FROM votes_miner_pct GROUP BY currency_id, pct ORDER BY currency_id, pct ASC")
			if err != nil {
				db.UnlockPrintSleep(utils.ErrInfo(err), 60) 			
				continue BEGIN
			}
			defer rows.Close()
			for rows.Next() {
				var currency_id, votes int64
				var pct string
				err = rows.Scan(&currency_id, &pct, &votes)
				if err!= nil {
					db.UnlockPrintSleep(utils.ErrInfo(err), 60)
					continue BEGIN
				}
				log.Println("newpctcurrency_id", currency_id, "pct", pct, "votes", votes)
				if len(pctVotes[currency_id]) == 0 {
					pctVotes[currency_id] = make(map[string]map[string]int64)
				}
				if len(pctVotes[currency_id]["miner_pct"]) == 0 {
					pctVotes[currency_id]["miner_pct"] = make(map[string]int64)
				}
				pctVotes[currency_id]["miner_pct"][pct] = votes
			}
	
			// берем все голоса user_pct
			rows, err = db.Query("SELECT currency_id, pct, count(user_id) as votes FROM votes_user_pct GROUP BY currency_id, pct ORDER BY currency_id, pct ASC")
			if err != nil {
				db.UnlockPrintSleep(utils.ErrInfo(err), 60)
				continue BEGIN
			}
			defer rows.Close()
			for rows.Next() {
				var currency_id, votes int64
				var pct string
				err = rows.Scan(&currency_id, &pct, &votes)
				if err!= nil {
					db.UnlockPrintSleep(utils.ErrInfo(err), 60)
					continue BEGIN
				}
				log.Println("currency_id", currency_id, "pct", pct, "votes", votes)
				if len(pctVotes[currency_id]) == 0 {
					pctVotes[currency_id] = make(map[string]map[string]int64)
				}
				if len(pctVotes[currency_id]["user_pct"]) == 0 {
					pctVotes[currency_id]["user_pct"] = make(map[string]int64)
				}
				pctVotes[currency_id]["user_pct"][pct] = votes
			}
	
			newPct := make(map[string]map[string]map[string]string)
			newPct["currency"] = make(map[string]map[string]string)
			var userMaxKey int64
			PctArray := utils.GetPctArray()
	
			log.Println("pctVotes", pctVotes)
			for currencyId, data := range pctVotes {
	
				currencyIdStr := utils.Int64ToStr(currencyId)
				// определяем % для майнеров
				pctArr := utils.MakePctArray(data["miner_pct"])
				log.Println("pctArrminer_pct", pctArr, currencyId)
				key := utils.GetMaxVote(pctArr, 0, 390, 100)
				log.Println("key", key)
				if len(newPct["currency"][currencyIdStr]) == 0{
					newPct["currency"][currencyIdStr] = make(map[string]string)
				}
				newPct["currency"][currencyIdStr]["miner_pct"] = utils.GetPctValue(key)
	
				// определяем % для юзеров
				pctArr = utils.MakePctArray(data["user_pct"])
				log.Println("pctArruser_pct", pctArr, currencyId)

				log.Println("newPct", newPct)
				pctY := utils.ArraySearch(newPct["currency"][currencyIdStr]["miner_pct"], PctArray)
				log.Println("newPct[currency][currencyIdStr][miner_pct]", newPct["currency"][currencyIdStr]["miner_pct"])
				log.Println("PctArray", PctArray)
				log.Println("miner_pct $pct_y=", pctY)
				maxUserPctY := utils.Round(utils.StrToFloat64(pctY)/2, 2)
				userMaxKey = utils.FindUserPct(int(maxUserPctY))
				log.Println("maxUserPctY", maxUserPctY, "userMaxKey", userMaxKey, "currencyIdStr", currencyIdStr)
				// отрезаем лишнее, т.к. поиск идет ровно до макимального возможного, т.е. до miner_pct/2
				pctArr = utils.DelUserPct(pctArr, userMaxKey);
				log.Println("pctArr", pctArr)

				key = utils.GetMaxVote(pctArr, 0, userMaxKey, 100)
				log.Println("data[user_pct]", data["user_pct"])
				log.Println("pctArr", pctArr)
				log.Println("userMaxKey", userMaxKey)
				log.Println("key", key)
				newPct["currency"][currencyIdStr]["user_pct"] = utils.GetPctValue(key)
				log.Println("user_pct", newPct["currency"][currencyIdStr]["user_pct"])
			}
	
			newPct_ := new(newPctType)
			newPct_.Currency = make(map[string]map[string]string)
			newPct_.Currency = newPct["currency"]
			newPct_.Referral = make(map[string]int64)
			refLevels := []string{"first", "second", "third"}
			for i:=0; i<len(refLevels); i++ {
				level := refLevels[i]
				var votesReferral []map[int64]int64
	        	// берем все голоса
				rows, err := db.Query("SELECT "+level+", count(user_id) as votes FROM votes_referral GROUP BY "+level+" ORDER BY "+level+" ASC ")
				if err != nil {
					db.UnlockPrintSleep(utils.ErrInfo(err), 60)
					continue BEGIN
				}
				defer rows.Close()
				for rows.Next() {
					var level_, votes int64
					err = rows.Scan(&level_, &votes)
					if err!= nil {
						db.UnlockPrintSleep(utils.ErrInfo(err), 60)
						continue BEGIN
					}
					votesReferral = append(votesReferral, map[int64]int64{level_:votes})
				}
				newPct_.Referral[level] = (utils.GetMaxVote(votesReferral, 0, 30, 10))
			}
			jsonData, err := json.Marshal(newPct_)
			if err!= nil {
				db.UnlockPrintSleep(utils.ErrInfo(err), 60)
				continue BEGIN
			}

			_, myUserId, _, _, _, _, err := db.TestBlock();
			forSign := fmt.Sprintf("%v,%v,%v,%v,%v,%v", utils.TypeInt("NewPct"), curTime, myUserId, jsonData)
			binSign, err := db.GetBinSign(forSign, myUserId)
			if err!= nil {
				db.UnlockPrintSleep(utils.ErrInfo(err), 60)
				continue BEGIN
			}
			data := utils.DecToBin(utils.TypeInt("NewPct"), 1)
			data = append(data, utils.DecToBin(curTime, 4)...)
			data = append(data, utils.EncodeLengthPlusData(utils.Int64ToByte(myUserId))...)
			data = append(data, utils.EncodeLengthPlusData(jsonData)...)
			data = append(data, utils.EncodeLengthPlusData([]byte(binSign))...)

			err = db.InsertReplaceTxInQueue(data)
			if err!= nil {
				db.UnlockPrintSleep(utils.ErrInfo(err), 60)
				continue BEGIN
			}

			// и не закрывая main_lock переводим нашу тр-ию в verified=1, откатив все несовместимые тр-ии
			// таким образом у нас будут в блоке только актуальные голоса.
			// а если придет другой блок и станет verified=0, то эта тр-ия просто удалится.

			p := new(dcparser.Parser)
			err = p.TxParser(data, utils.HexToBin(utils.Md5(data)), true)
			if err != nil {
				db.PrintSleep(err, 60)
				continue BEGIN
			}
		}
		db.DbUnlock()
		utils.Sleep(60)
	}
	return ""
}

type newPctType struct {
	Currency map[string]map[string]string `json:"currency"`
	Referral map[string]int64 `json:"referral"`
}
