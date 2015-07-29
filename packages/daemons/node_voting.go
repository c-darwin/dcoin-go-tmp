package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/consts"
	"io/ioutil"
	//"log"
	"os"
)

/*
 * Если наш miner_id есть среди тех, кто должен скачать фото нового майнера к себе, то качаем
 */


func NodeVoting() {

	const GoroutineName = "NodeVoting"

	db := DbConnect()
	if db == nil {
		return
	}
	db.GoroutineName = GoroutineName
	if !db.CheckInstall(DaemonCh, AnswerDaemonCh) {
		return
	}

	dir, err := utils.GetCurrentDir()
	if err != nil {
		return
	}

	BEGIN:
	for {
		log.Info(GoroutineName)
		// проверим, не нужно ли нам выйти из цикла
		if CheckDaemonsRestart() {
			break BEGIN
		}

		err, restart := db.DbLock(DaemonCh, AnswerDaemonCh)
		if restart {
			break BEGIN
		}
		if err != nil {
			db.PrintSleep(err, 1)
			continue BEGIN
		}

		// берем данные, которые находятся на голосовании нодов
		rows, err := db.Query(`
				SELECT miners_data.user_id,
							 http_host as host,
							 face_hash,
							 profile_hash,
							 photo_block_id,
							 photo_max_miner_id,
							 miners_keepers,
							 id as vote_id,
							 miner_id
				FROM votes_miners
				LEFT JOIN miners_data
						 ON votes_miners.user_id = miners_data.user_id
				WHERE cron_checked_time < ? AND
							 votes_end = 0 AND
							 type = 'node_voting'
				`, utils.Time() - consts.CRON_CHECKED_TIME_SEC)
		if err != nil {
			db.PrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		defer rows.Close()
		if  ok := rows.Next(); ok {
			var vote_id, miner_id int64;
			var user_id, host, row_face_hash, row_profile_hash, photo_block_id, photo_max_miner_id, miners_keepers string;
			err = rows.Scan(&user_id, &host, &row_face_hash, &row_profile_hash, &photo_block_id, &photo_max_miner_id, &miners_keepers, &vote_id, &miner_id)
			if err != nil {
				db.PrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}

			// проверим, не нужно нам выйти, т.к. обновилась версия софта
			if CheckDaemonsRestart() {
				utils.Sleep(1)
				break
			}
			minersIds := utils.GetMinersKeepers(photo_block_id, photo_max_miner_id, miners_keepers, true);
			myUsersIds, err := db.GetMyUsersIds(true, true)
			myMinersIds, err := db.GetMyMinersIds(myUsersIds)

			// нет ли нас среди тех, кто должен скачать фото к себе и проголосовать
			var intersectMyMiners []int64
			for _, id := range minersIds {
				if utils.InSliceInt64(int64(id), myMinersIds)	{
					intersectMyMiners = append(intersectMyMiners, int64(id))
				}
			}
			var vote int64
			if len(intersectMyMiners) > 0 {
				// копируем фото  к себе
				profilePath := dir+"/public/profile_"+user_id+".jpg";
				_, err = utils.DownloadToFile(host+"/public/"+user_id+"_user_profile.jpg", profilePath, 60, DaemonCh, AnswerDaemonCh)
				if err != nil {
					db.PrintSleep(utils.ErrInfo(err), 1)
					continue BEGIN
				}
				facePath := dir+"/public/face_"+user_id+".jpg";
				_, err = utils.DownloadToFile(host+"/public/"+user_id+"_user_face.jpg", facePath, 60, DaemonCh, AnswerDaemonCh)
				if err != nil {
					db.PrintSleep(utils.ErrInfo(err), 1)
					continue BEGIN
				}
				// хэши скопированных фото
				profileFile, err := ioutil.ReadFile(profilePath)
				if err != nil {
					db.PrintSleep(utils.ErrInfo(err), 1)
					continue BEGIN
				}
				profileHash := string(utils.DSha256(profileFile))
				log.Info("%v", "profileHash", profileHash)
				faceFile, err := ioutil.ReadFile(facePath)
				if err != nil {
					db.PrintSleep(utils.ErrInfo(err), 1)
					continue BEGIN
				}
				faceHash := string(utils.DSha256(faceFile))
				log.Info("%v", "faceHash", faceHash)

				// проверяем хэш. Если сходится, то голосуем за, если нет - против и размер не должен быть более 200 Kb.
				if profileHash == row_face_hash && faceHash == row_profile_hash && len(profileFile) < 204800 && len(faceFile) < 204800 {
					vote = 1
				} else {
					vote = 0 // если хэш не сходится, то удаляем только что скаченное фото
					os.Remove(profilePath)
					os.Remove(facePath)
				}

				// проходимся по всем нашим майнерам, если это пул и по одному, если это сингл-мод
				for _, myMinerId := range intersectMyMiners{

					myUserId, err := db.Single("SELECT user_id FROM miners_data WHERE miner_id  =  ?", myMinerId).Int64()
					if err != nil {
						db.PrintSleep(utils.ErrInfo(err), 1)
						continue BEGIN
					}

					curTime := utils.Time()

					forSign := fmt.Sprintf("%v,%v,%v,%v,%v", utils.TypeInt("VotesNodeNewMiner"), curTime, myUserId, vote_id, vote)
					binSign, err := db.GetBinSign(forSign, myUserId)
					if err!= nil {
						db.UnlockPrintSleep(utils.ErrInfo(err), 1)
						continue BEGIN
					}
					data := utils.DecToBin(utils.TypeInt("VotesNodeNewMiner"), 1)
					data = append(data, utils.DecToBin(curTime, 4)...)
					data = append(data, utils.EncodeLengthPlusData(utils.Int64ToByte(myUserId))...)
					data = append(data, utils.EncodeLengthPlusData(utils.Int64ToByte(vote_id))...)
					data = append(data, utils.EncodeLengthPlusData(utils.Int64ToByte(vote))...)
					data = append(data, utils.EncodeLengthPlusData([]byte(binSign))...)


					err = db.InsertReplaceTxInQueue(data)
					if err!= nil {
						db.UnlockPrintSleep(utils.ErrInfo(err), 1)
						continue BEGIN
					}

				}
			}

			// отмечаем, чтобы больше не брать эту строку
			err = db.ExecSql("UPDATE votes_miners SET cron_checked_time = ? WHERE id = ?", utils.Time(), vote_id)
			if err != nil {
				db.UnlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
		}
		db.DbUnlock()
		for i:=0; i < 60; i++ {
			utils.Sleep(1)
			// проверим, не нужно ли нам выйти из цикла
			if CheckDaemonsRestart() {
				break BEGIN
			}
		}
	}

}
