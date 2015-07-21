package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"io/ioutil"
	"regexp"
)

func CleaningDb() {

	const GoroutineName = "CleaningDb"

	db := DbConnect()
	if db == nil {
		return
	}
	db.GoroutineName = GoroutineName
	if !db.CheckInstall(DaemonCh, AnswerDaemonCh) {
		return
	}

	BEGIN:
	for {
		log.Info(GoroutineName)
		// проверим, не нужно ли нам выйти из цикла
		if CheckDaemonsRestart() {
			break BEGIN
		}
		err = db.DbLock(DaemonCh, AnswerDaemonCh)
		if err != nil {
			db.PrintSleep(utils.ErrInfo(err), 0)
			break BEGIN
		}

		curBlockId, err := db.GetBlockId()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		// пишем свежие блоки в резервный блокчейн
		endBlockId, err := utils.GetEndBlockId()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		if curBlockId - 30 > endBlockId {
			blocks, err := db.GetMap(`
					SELECT id, data
					FROM block_chain
					WHERE id > ? AND id < = ?
					ORDER BY id
					`, "id", "data", endBlockId, curBlockId-30)
			if err != nil {
				db.UnlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			for id, data := range blocks {
				blockData := append(utils.DecToBin(id, 5), utils.EncodeLengthPlusData(data)...)
				sizeAndData := append(utils.DecToBin(len(data), 5), blockData...)
				err := ioutil.WriteFile("public/blockchain", append(sizeAndData, utils.DecToBin(len(sizeAndData), 5)...), 0644)
				if err != nil {
					db.UnlockPrintSleep(utils.ErrInfo(err), 1)
					continue BEGIN
				}
			}
		}

		autoReload, err := db.Single("SELECT auto_reload FROM config").Int64()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		if autoReload < 60 {
			db.UnlockPrintSleep(utils.ErrInfo("autoReload < 60"), 1)
			continue BEGIN
		}

		// если main_lock висит более x минут, значит был какой-то сбой
		mainLock, err := db.Single("SELECT lock_time FROM main_lock WHERE script_name NOT IN ('my_lock', 'cleaning_db')").Int64()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		if mainLock > 0 && utils.Time() - autoReload > mainLock {
			// на всякий случай пометим, что работаем
			err = db.ExecSql("UPDATE main_lock SET script_name = 'cleaning_db'")
			if err != nil {
				db.UnlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			err = db.ExecSql("UPDATE config SET pool_tech_works = 1")
			if err != nil {
				db.UnlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			allTables, err := db.GetAllTables()
			if err != nil {
				db.UnlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			for _, table := range allTables {

				if ok, _ := regexp.MatchString(`my_|install|config|daemons|payment_systems|community|cf_lang`, table); !ok{
					err = db.ExecSql("DELETE FROM "+table)
					if err != nil {
						db.UnlockPrintSleep(utils.ErrInfo(err), 1)
						continue BEGIN
					}
					if table == "cf_currency" {
						err = db.ExecSql("ALTER TABLE cf_currency auto_increment = 1000")
						if err != nil {
							db.UnlockPrintSleep(utils.ErrInfo(err), 1)
							continue BEGIN
						}
					} else if table == "admin" {
						err = db.ExecSql("INSERT INTO admin (user_id) VALUES (1)")
						if err != nil {
							db.UnlockPrintSleep(utils.ErrInfo(err), 1)
							continue BEGIN
						}
					}
				}
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
