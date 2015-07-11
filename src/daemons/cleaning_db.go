package daemons

import (
	"utils"
)

func cleaning_db(configIni map[string]string) {

	const GoroutineName = "cleaning_db"

	db := utils.DbConnect(configIni)
	db.GoroutineName = GoroutineName
	BEGIN:
	for {

		err := db.DbLock()
		if err != nil {
			db.PrintSleep(utils.ErrInfo(err), 10)
			continue BEGIN
		}
		blockId, err := db.GetBlockId()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 10)
			continue BEGIN
		}

		// пишем свежие блоки в резервный блокчейн
		endBlockId, err :=

		db.DbUnlock()

		utils.Sleep(10)
	}
}
