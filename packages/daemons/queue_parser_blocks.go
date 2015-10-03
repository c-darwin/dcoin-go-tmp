package daemons

import (
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"math/big"
)

/* Берем блок. Если блок имеет лучший хэш, то ищем, в каком блоке у нас пошла вилка
 * Если вилка пошла менее чем variables->rollback_blocks блоков назад, то
 *  - получаем всю цепочку блоков,
 *  - откатываем фронтальные данные от наших блоков,
 *  - заносим фронт. данные из новой цепочки
 *  - если нет ошибок, то откатываем наши данные из блоков
 *  - и заносим новые данные
 *  - если где-то есть ошибки, то откатываемся к нашим прежним данным
 * Если вилка была давно, то ничего не трогаем, и оставлеяем скрипту blocks_collection.php
 * Ограничение variables->rollback_blocks нужно для защиты от подставных блоков
 *
 * */

func QueueParserBlocks() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("daemon Recovered", r)
			panic(r)
		}
	}()

	if utils.Mobile() {
		sleepTime = 1800
	} else {
		sleepTime = 10
	}
	const GoroutineName = "QueueParserBlocks"
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

		err, restart := d.dbLock()
		if restart {
			break BEGIN
		}
		if err != nil {
			if d.dPrintSleep(err, sleepTime) {	break BEGIN }
			continue BEGIN
		}

		prevBlockData, err := d.OneRow("SELECT * FROM info_block").String()
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), sleepTime) {	break BEGIN }
			continue BEGIN
		}
		newBlockData, err := d.OneRow("SELECT * FROM queue_blocks").String()
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), sleepTime) {	break BEGIN }
			continue BEGIN
		}
		if len(newBlockData) == 0 {
			if d.unlockPrintSleep(utils.ErrInfo(err), sleepTime) {	break BEGIN }
			continue BEGIN
		}
		newBlockData["head_hash_hex"] = string(utils.BinToHex(newBlockData["head_hash"]))
		prevBlockData["head_hash_hex"] = string(utils.BinToHex(prevBlockData["head_hash"]))
		newBlockData["hash_hex"] = string(utils.BinToHex(newBlockData["hash"]))
		prevBlockData["hash_hex"] = string(utils.BinToHex(prevBlockData["hash"]))

		variables, err := d.GetAllVariables()
		if err != nil {
			if d.unlockPrintSleep(utils.ErrInfo(err), sleepTime) {	break BEGIN }
			continue BEGIN
		}

		/*
		 * Базовая проверка
		 */

		// проверим, укладывается ли блок в лимит rollback_blocks_1
		if utils.StrToInt64(newBlockData["block_id"]) > utils.StrToInt64(prevBlockData["block_id"])+variables.Int64["rollback_blocks_1"] {
			d.DeleteQueueBlock(newBlockData["head_hash_hex"], newBlockData["hash_hex"])
			if d.unlockPrintSleep(utils.ErrInfo("rollback_blocks_1"), 1) {	break BEGIN }
			continue BEGIN
		}

		// проверим не старый ли блок в очереди
		if utils.StrToInt64(newBlockData["block_id"]) < utils.StrToInt64(prevBlockData["block_id"]) {
			d.DeleteQueueBlock(newBlockData["head_hash_hex"], newBlockData["hash_hex"])
			if d.unlockPrintSleep(utils.ErrInfo("old block"), 1) {	break BEGIN }
			continue BEGIN
		}

		if utils.StrToInt64(newBlockData["block_id"]) == utils.StrToInt64(prevBlockData["block_id"]) {
			// сравним хэши
			hash1 := big.NewInt(0)
			hash1.SetString(string(newBlockData["head_hash_hex"]), 16)
			hash2 := big.NewInt(0)
			hash2.SetString(string(prevBlockData["head_hash_hex"]), 16)
			// newBlockData["head_hash_hex"]) <= prevBlockData["head_hash_hex"]
			if hash1.Cmp(hash2) < 1 {
				// если это тотже блок и его генерил тот же юзер, то могут быть равные head_hash
				if hash1.Cmp(hash2) == 0 {
					// в этом случае проверяем вторые хэши. Если новый блок имеет больший хэш, то нам он не нужен
					// или если тот же хэш, значит блоки одинаковые

					hash1 := big.NewInt(0)
					hash1.SetString(string(newBlockData["hash_hex"]), 16)
					hash2 := big.NewInt(0)
					hash2.SetString(string(prevBlockData["hash_hex"]), 16)
					// newBlockData["head_hash_hex"]) >= prevBlockData["head_hash_hex"]
					if hash1.Cmp(hash2) >= 0 {
						d.DeleteQueueBlock(newBlockData["head_hash_hex"], newBlockData["hash_hex"])
						if d.unlockPrintSleep(utils.ErrInfo("newBlockData hash_hex == prevBlockData hash_hex"), 1) {	break BEGIN }
						continue BEGIN
					}
				}
			} else {
				d.DeleteQueueBlock(newBlockData["head_hash_hex"], newBlockData["hash_hex"])
				if d.unlockPrintSleep(utils.ErrInfo("newBlockData head_hash_hex >  prevBlockData head_hash_hex"), 1) {	break BEGIN }
				continue BEGIN
			}
		}

		/*
		 * Загрузка блоков для детальной проверки
		 */
		host, err := d.Single("SELECT tcp_host FROM miners_data WHERE user_id  =  ?", newBlockData["user_id"]).String()
		if err != nil {
			d.DeleteQueueBlock(newBlockData["head_hash_hex"], newBlockData["hash_hex"])
			if d.unlockPrintSleep(utils.ErrInfo(err), sleepTime) {	break BEGIN }
			continue BEGIN
		}
		blockId := utils.StrToInt64(newBlockData["block_id"])

		p := new(dcparser.Parser)
		p.DCDB = d.DCDB
		p.GoroutineName = GoroutineName
		err = p.GetBlocks(blockId, host, utils.StrToInt64(newBlockData["user_id"]), "rollback_blocks_1", GoroutineName, 7, "")
		if err != nil {
			d.DeleteQueueBlock(newBlockData["head_hash_hex"], newBlockData["hash_hex"])
			d.NodesBan(utils.StrToInt64(newBlockData["user_id"]), fmt.Sprintf("%v", err))
			if d.unlockPrintSleep(utils.ErrInfo(err), 1) {	break BEGIN }
			continue BEGIN
		}

		d.dbUnlock()

		if d.dSleep(sleepTime) {
			break BEGIN
		}
	}

}
