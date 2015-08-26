package daemons

import (
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	//"log"
	"errors"
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
)

/**
 * Демон, который отсчитывает время, которые необходимо ждать после того,
 * как началось одноуровневое соревнование, у кого хэш меньше.
 * Когда время прошло, то берется блок из таблы testblock и заносится в
 * queue и queue_front для занесение данных к себе и отправки другим
 *
 */

func TestblockIsReady() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("daemon Recovered", r)
			panic(r)
		}
	}()

	const GoroutineName = "TestblockIsReady"
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
		MonitorDaemonCh <- []string{GoroutineName, utils.Int64ToStr(utils.Time())}

		// проверим, не нужно ли нам выйти из цикла
		if CheckDaemonsRestart() {
			break BEGIN
		}

		LocalGateIp, err := d.GetMyLocalGateIp()
		if err != nil {
			d.PrintSleep(utils.ErrInfo(err), 1)
			continue
		}
		if len(LocalGateIp) > 0 {
			d.PrintSleep(utils.ErrInfo(errors.New("len(LocalGateIp) > 0")), 5)
			continue
		}

		// сколько нужно спать
		prevBlock, myUserId, myMinerId, currentUserId, level, levelsRange, err := d.TestBlock()
		if err != nil {
			d.PrintSleep(utils.ErrInfo(err), 1)
			continue
		}
		log.Info("%v", prevBlock, myUserId, myMinerId, currentUserId, level, levelsRange)

		if myMinerId == 0 {
			d.PrintSleepInfo(utils.ErrInfo(errors.New("myMinerId == 0 ")), 1)
			continue
		}

		sleepData, err := d.GetSleepData();
		sleep := d.GetIsReadySleep(prevBlock.Level, sleepData["is_ready"])
		prevHeadHash := prevBlock.HeadHash

		// Если случится откат или придет новый блок, то testblock станет неактуален
		startSleep := utils.Time()
		for i:=0; i < int(sleep); i++ {
			err, restart := d.dbLock()
			if restart {
				break BEGIN
			}
			if err != nil {
				d.PrintSleep(err, 1)
				continue BEGIN
			}

			newHeadHash, err := d.Single("SELECT head_hash FROM info_block").String()
			if err != nil {
				d.PrintSleep(utils.ErrInfo(err), 1)
				continue
			}
			d.dbUnlock()
			newHeadHash = string(utils.BinToHex([]byte(newHeadHash)))
			if newHeadHash != prevHeadHash {
				d.PrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			log.Info("%v", "i", i, "time", utils.Time())
			if utils.Time() - startSleep > sleep {
				break
			}
			utils.Sleep(1)  // спим 1 сек. общее время = $sleep
		}


		/*
		Заголовок
		TYPE (0-блок, 1-тр-я)       FF (256)
		BLOCK_ID   				       FF FF FF FF (4 294 967 295)
		TIME       					       FF FF FF FF (4 294 967 295)
		USER_ID                          FF FF FF FF FF (1 099 511 627 775)
		LEVEL                              FF (256)
		SIGN                               от 128 байта до 512 байт. Подпись от TYPE, BLOCK_ID, PREV_BLOCK_HASH, TIME, USER_ID, LEVEL, MRKL_ROOT
		Далее - тело блока (Тр-ии)
		*/

		// блокируем изменения данных в тестблоке
		// также, нужно блокировать main, т.к. изменение в info_block и block_chain ведут к изменению подписи в testblock
		err, restart := d.dbLock()
		if restart {
			break BEGIN
		}
		if err != nil {
			d.PrintSleep(err, 1)
			continue BEGIN
		}
		// за промежуток в main_unlock и main_lock мог прийти новый блок
		prevBlock, myUserId, myMinerId, currentUserId, level, levelsRange, err = d.TestBlock()
		if err != nil {
			d.PrintSleep(utils.ErrInfo(err), 1)
			continue
		}
		log.Info("%v", prevBlock, myUserId, myMinerId, currentUserId, level, levelsRange)

		// на всякий случай убедимся, что блок не изменился
		if prevBlock.HeadHash != prevHeadHash {
			d.unlockPrintSleep(utils.ErrInfo(errors.New("prevBlock.HeadHash != prevHeadHash")), 1)
			continue
		}

		// составим блок. заголовок + тело + подпись
		testBlockData, err := d.OneRow("SELECT * FROM testblock WHERE status  =  'active'").String()
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(errors.New("prevBlock.HeadHash != prevHeadHash")), 1)
			continue
		}
		log.Debug("testBlockData: %v", testBlockData)
		if len(testBlockData) == 0 {
			d.unlockPrintSleep(utils.ErrInfo(errors.New("null $testblock_data")), 1)
			continue
		}
		// получим транзакции
		var testBlockDataTx []byte
		transactionsTestBlock, err := d.GetList("SELECT data FROM transactions_testblock ORDER BY id ASC").String()
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		for _, data := range transactionsTestBlock {
			testBlockDataTx = append(testBlockDataTx, utils.EncodeLengthPlusData([]byte(data))...)
		}

		// в промежутке межде тем, как блок был сгенерирован и запуском данного скрипта может измениться текущий блок
		// поэтому нужно проверять подпись блока из тестблока
		prevBlockHash, err := d.Single("SELECT hash FROM info_block").Bytes()
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		prevBlockHash = utils.BinToHex(prevBlockHash)
		nodePublicKey, err := d.GetNodePublicKey(utils.StrToInt64(testBlockData["user_id"]))
		if err != nil {
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		forSign := fmt.Sprintf("0,%v,%s,%v,%v,%v,%s", testBlockData["block_id"], prevBlockHash, testBlockData["time"], testBlockData["user_id"], testBlockData["level"],utils.BinToHex([]byte(testBlockData["mrkl_root"])))
		log.Debug("forSign %v", forSign)
		log.Debug("signature %x", testBlockData["signature"])

		p := new(dcparser.Parser)
		p.DCDB = d.DCDB
		// проверяем подпись
		_, err = utils.CheckSign([][]byte{nodePublicKey}, forSign, []byte(testBlockData["signature"]), true);
		if err != nil {
			log.Error("incorrect signature %v")
			p.RollbackTransactionsTestblock(true)
			err = d.ExecSql("DELETE FROM testblock")
			if err != nil {
				d.unlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		// БАГ
		if utils.StrToInt64(testBlockData["block_id"]) == prevBlock.BlockId {
			log.Error("testBlockData block_id =  prevBlock.BlockId (%v=%v)", testBlockData["block_id"], prevBlock.BlockId)

			err = p.RollbackTransactionsTestblock(true)
			if err != nil {
				d.unlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			err = d.ExecSql("DELETE FROM testblock")
			if err != nil {
				d.unlockPrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			d.unlockPrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		// готовим заголовок
		newBlockIdBinary := utils.DecToBin(utils.StrToInt64(testBlockData["block_id"]), 4 );
		timeBinary := utils.DecToBin(utils.StrToInt64(testBlockData["time"]), 4 );
		userIdBinary := utils.DecToBin(utils.StrToInt64(testBlockData["user_id"]), 5 );
		levelBinary := utils.DecToBin(utils.StrToInt64(testBlockData["level"]), 1 );
		//prevBlockHashBinary := prevBlock.Hash
		//merkleRootBinary := testBlockData["mrklRoot"];

		// заголовок
		blockHeader := utils.DecToBin(0, 1)
		blockHeader = append(blockHeader, newBlockIdBinary...)
		blockHeader = append(blockHeader, timeBinary...)
		blockHeader = append(blockHeader, userIdBinary...)
		blockHeader = append(blockHeader, levelBinary...)
		blockHeader = append(blockHeader, utils.EncodeLengthPlusData([]byte(testBlockData["signature"]))...)

		// сам блок
		block := append(blockHeader, testBlockDataTx...)
		log.Debug("block %x", block)

		// теперь нужно разнести блок по таблицам и после этого мы будем его слать всем нодам скриптом disseminator.php
		p.BinaryData = block
		err = p.ParseDataFront()
		if err != nil {
			d.PrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		// и можно удалять данные о тестблоке, т.к. они перешел в нормальный блок
		err = d.ExecSql("DELETE FROM transactions_testblock")
		if err != nil {
			d.PrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}
		err = d.ExecSql("DELETE FROM testblock")
		if err != nil {
			d.PrintSleep(utils.ErrInfo(err), 1)
			continue BEGIN
		}

		// между testblock_generator и testbock_is_ready
		p.RollbackTransactionsTestblock(false)

		d.dbUnlock()

		log.Info("%v", "Happy end")

		utils.Sleep(1)
	}


}
