package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
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

	const GoroutineName = "QueueParserBlocks"
	db := DbConnect()
	db.GoroutineName = GoroutineName
	db.CheckInstall()
BEGIN:
	for {

		// проверим, не нужно нам выйти, т.к. обновилась версия софта
		if db.CheckDaemonRestart() {
			utils.Sleep(1)
			break
		}

		err := db.DbLock()
		if err != nil {
			db.PrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		prevBlockData, err := db.OneRow("SELECT * FROM info_block").String()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		newBlockData, err := db.OneRow("SELECT * FROM queue_blocks").String()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		if len(newBlockData) == 0 {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		newBlockData["head_hash_hex"] = string(utils.BinToHex(newBlockData["head_hash"]))
		prevBlockData["head_hash_hex"] = string(utils.BinToHex(prevBlockData["head_hash"]))
		newBlockData["hash_hex"] = string(utils.BinToHex(newBlockData["hash"]))
		prevBlockData["hash_hex"] = string(utils.BinToHex(prevBlockData["hash"]))

		variables, err := db.GetAllVariables()
		if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		/*
		 * Базовая проверка
		 */

		// проверим, укладывается ли блок в лимит rollback_blocks_1
		if utils.StrToInt64(newBlockData["block_id"]) > utils.StrToInt64(prevBlockData["block_id"]) + variables.Int64["rollback_blocks_1"] {
			db.DeleteQueueBlock(newBlockData["head_hash_hex"], newBlockData["hash_hex"])
			db.UnlockPrintSleep(utils.ErrInfo("rollback_blocks_1"), 1)
			continue BEGIN
		}

		// проверим не старый ли блок в очереди
		if utils.StrToInt64(newBlockData["block_id"]) < utils.StrToInt64(prevBlockData["block_id"]) {
			db.DeleteQueueBlock(newBlockData["head_hash_hex"], newBlockData["hash_hex"])
			db.UnlockPrintSleep(utils.ErrInfo("old block"), 1)
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
						db.DeleteQueueBlock(newBlockData["head_hash_hex"], newBlockData["hash_hex"])
						db.UnlockPrintSleep(utils.ErrInfo("newBlockData hash_hex == prevBlockData hash_hex"), 1)
						continue BEGIN
					}
				}
			} else {
				db.DeleteQueueBlock(newBlockData["head_hash_hex"], newBlockData["hash_hex"])
				db.UnlockPrintSleep(utils.ErrInfo("newBlockData head_hash_hex >  prevBlockData head_hash_hex"), 1)
				continue BEGIN
			}
		}

		/*
		 * Загрузка блоков для детальной проверки
		 */
		host, err := db.Single("SELECT tcp_host FROM miners_data WHERE user_id  =  ?", newBlockData["user_id"]).String()
		if err != nil {
			db.DeleteQueueBlock(newBlockData["head_hash_hex"], newBlockData["hash_hex"])
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		blockId := utils.StrToInt64(newBlockData["block_id"])

		p := new(dcparser.Parser)
		//func (p *Parser) GetBlocks (blockId int64, host string, userId int64, rollbackBlocks, goroutineName, getBlockScriptName, addNodeHost string) error {
		err = p.GetBlocks(blockId, host, utils.StrToInt64(newBlockData["user_id"]), "rollback_blocks_1", GoroutineName, 7, "")
		if err != nil {
			db.DeleteQueueBlock(newBlockData["head_hash_hex"], newBlockData["hash_hex"])
			db.NodesBan(utils.StrToInt64(newBlockData["user_id"]), fmt.Sprintf("%v", err))
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		db.DbUnlock()
		utils.Sleep(10)
	}

}
