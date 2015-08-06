package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"regexp"
	"os"
)

func CleaningDb() {

	const GoroutineName = "CleaningDb"
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

		curBlockId, err := d.GetBlockId()
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		// пишем свежие блоки в резервный блокчейн
		endBlockId, err := utils.GetEndBlockId()
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		if curBlockId - 30 > endBlockId {
			blocks, err := d.GetMap(`
					SELECT id, data
					FROM block_chain
					WHERE id > ? AND id <= ?
					ORDER BY id
					`, "id", "data", endBlockId, curBlockId-30)
			if err != nil {
				d.unlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			file, err := os.OpenFile(*utils.Dir+"/public/blockchain", os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				d.unlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			for id, data := range blocks {
				blockData := append(utils.DecToBin(id, 5), utils.EncodeLengthPlusData(data)...)
				sizeAndData := append(utils.DecToBin(len(data), 5), blockData...)
				//err := ioutil.WriteFile(*utils.Dir+"/public/blockchain", append(sizeAndData, utils.DecToBin(len(sizeAndData), 5)...), 0644)
				if _, err = file.Write(append(sizeAndData, utils.DecToBin(len(sizeAndData), 5)...)); err != nil {
					file.Close()
					d.unlockPrintSleep(utils.ErrInfo(err), 1)
					continue BEGIN
				}
				if err != nil {
					file.Close()
					d.unlockPrintSleep(utils.ErrInfo(err), 1)
					continue BEGIN
				}
			}
			file.Close()
		}

		autoReload, err := d.Single("SELECT auto_reload FROM config").Int64()
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		if autoReload < 60 {
			d.unlockPrintSleep(utils.ErrInfo("autoReload < 60"), 1)
			continue BEGIN
		}

		// если main_lock висит более x минут, значит был какой-то сбой
		mainLock, err := d.Single("SELECT lock_time FROM main_lock WHERE script_name NOT IN ('my_lock', 'cleaning_db')").Int64()
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		if mainLock > 0 && utils.Time() - autoReload > mainLock {
			// на всякий случай пометим, что работаем
			err = d.ExecSql("UPDATE main_lock SET script_name = 'cleaning_db'")
			if err != nil {
				d.unlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			err = d.ExecSql("UPDATE config SET pool_tech_works = 1")
			if err != nil {
				d.unlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			allTables, err := d.GetAllTables()
			if err != nil {
				d.unlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			for _, table := range allTables {

				if ok, _ := regexp.MatchString(`my_|install|config|daemons|payment_systems|community|cf_lang`, table); !ok{
					err = d.ExecSql("DELETE FROM "+table)
					if err != nil {
						d.unlockPrintSleep(utils.ErrInfo(err), 1)
						continue BEGIN
					}
					if table == "cf_currency" {
						err = d.ExecSql("ALTER TABLE cf_currency auto_increment = 1000")
						if err != nil {
							d.unlockPrintSleep(utils.ErrInfo(err), 1)
							continue BEGIN
						}
					} else if table == "admin" {
						err = d.ExecSql("INSERT INTO admin (user_id) VALUES (1)")
						if err != nil {
							d.unlockPrintSleep(utils.ErrInfo(err), 1)
							continue BEGIN
						}
					}
				}
			}

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
