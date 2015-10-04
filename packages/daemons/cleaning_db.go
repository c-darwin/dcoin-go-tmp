package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"os"
	"regexp"
)

func CleaningDb() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("daemon Recovered", r)
			panic(r)
		}
	}()

	const GoroutineName = "CleaningDb"
	d := new(daemon)
	d.DCDB = DbConnect()
	if d.DCDB == nil {
		return
	}
	d.goRoutineName = GoroutineName
	if utils.Mobile() {
		d.sleepTime = 1800
	} else {
		d.sleepTime = 60
	}
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

		curBlockId, err := d.GetBlockId()
		if err != nil {
			if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}

		// пишем свежие блоки в резервный блокчейн
		endBlockId, err := utils.GetEndBlockId()
		if err != nil {
			if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		if curBlockId-30 > endBlockId {
			blocks, err := d.GetMap(`
					SELECT id, data
					FROM block_chain
					WHERE id > ? AND id <= ?
					ORDER BY id
					`, "id", "data", endBlockId, curBlockId-30)
			if err != nil {
				if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
			file, err := os.OpenFile(*utils.Dir+"/public/blockchain", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
			if err != nil {
				if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
			for id, data := range blocks {
				blockData := append(utils.DecToBin(id, 5), utils.EncodeLengthPlusData(data)...)
				sizeAndData := append(utils.DecToBin(len(data), 5), blockData...)
				//err := ioutil.WriteFile(*utils.Dir+"/public/blockchain", append(sizeAndData, utils.DecToBin(len(sizeAndData), 5)...), 0644)
				if _, err = file.Write(append(sizeAndData, utils.DecToBin(len(sizeAndData), 5)...)); err != nil {
					file.Close()
					if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
					continue BEGIN
				}
				if err != nil {
					file.Close()
					if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
					continue BEGIN
				}
			}
			file.Close()
		}

		autoReload, err := d.Single("SELECT auto_reload FROM config").Int64()
		if err != nil {
			if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		log.Debug("autoReload: %v", autoReload)
		if autoReload < 60 {
			if d.dPrintSleep(utils.ErrInfo("autoReload < 60"), d.sleepTime)  {	break BEGIN }
			continue BEGIN
		}

		// если main_lock висит более x минут, значит был какой-то сбой
		mainLock, err := d.Single("SELECT lock_time FROM main_lock WHERE script_name NOT IN ('my_lock', 'cleaning_db')").Int64()
		if err != nil {
			if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
			continue BEGIN
		}
		log.Debug("mainLock: %v", mainLock)
		log.Debug("utils.Time(): %v", utils.Time())
		if mainLock > 0 && utils.Time()-autoReload > mainLock {
			// на всякий случай пометим, что работаем
			err = d.ExecSql("UPDATE main_lock SET script_name = 'cleaning_db'")
			if err != nil {
				if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
			err = d.ExecSql("UPDATE config SET pool_tech_works = 1")
			if err != nil {
				if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
			allTables, err := d.GetAllTables()
			if err != nil {
				if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
			for _, table := range allTables {
				log.Debug("table: %s", table)
				if ok, _ := regexp.MatchString(`my_|install|config|daemons|payment_systems|community|cf_lang`, table); !ok {
					log.Debug("DELETE FROM %s", table)
					err = d.ExecSql("DELETE FROM " + table)
					if err != nil {
						if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
						continue BEGIN
					}
					if table == "cf_currency" {
						err = d.SetAI("cf_currency", 999)
						if err != nil {
							if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
							continue BEGIN
						}
					} else if table == "admin" {
						err = d.ExecSql("INSERT INTO admin (user_id) VALUES (1)")
						if err != nil {
							if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
							continue BEGIN
						}
					} else {
						log.Debug("SET AI %s", table)
						if d.ConfigIni["db_type"] == "postgresql" {
							err = d.SetAI(table, 1)
						} else {
							err = d.SetAI(table, 0)
						}
						if err != nil {
							log.Error("%v", err)
						}
					}
				}
			}
		}

		if d.dSleep(d.sleepTime) {
			break BEGIN
		}
	}
}
