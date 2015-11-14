package dcparser

import (
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/consts"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

/**
фронт. проверка + занесение данных из блока в таблицы и info_block
*/
func (p *Parser) ParseDataFull() error {

	p.TxIds = []string{}
	p.dataPre()
	if p.dataType != 0 { // парсим только блоки
		return utils.ErrInfo(fmt.Errorf("incorrect dataType"))
	}
	var err error

	p.Variables, err = p.GetAllVariables()
	if err != nil {
		return utils.ErrInfo(err)
	}

	//if len(p.BinaryData) > 500000 {
	//	ioutil.WriteFile("block-"+string(utils.DSha256(p.BinaryData)), p.BinaryData, 0644)
	//}

	err = p.ParseBlock()
	if err != nil {
		return utils.ErrInfo(err)
	}

	// проверим данные, указанные в заголовке блока
	err = p.CheckBlockHeader()
	if err != nil {
		return utils.ErrInfo(err)
	}


	utils.WriteSelectiveLog("DELETE FROM transactions WHERE used = 1")
	afect, err := p.ExecSqlGetAffect("DELETE FROM transactions WHERE used = 1")
	if err != nil {
		utils.WriteSelectiveLog(err)
		return utils.ErrInfo(err)
	}
	utils.WriteSelectiveLog("afect: "+utils.Int64ToStr(afect))

	txCounter := make(map[int64]int64)
	p.fullTxBinaryData = p.BinaryData
	var txForRollbackTo []byte
	if len(p.BinaryData) > 0 {
		for {
			// обработка тр-ий может занять много времени, нужно отметиться
			p.UpdDaemonTime(p.GoroutineName)
			p.halfRollback = false
			log.Debug("&p.BinaryData", p.BinaryData)
			transactionSize := utils.DecodeLength(&p.BinaryData)
			if len(p.BinaryData) == 0 {
				return utils.ErrInfo(fmt.Errorf("empty BinaryData"))
			}

			// отчекрыжим одну транзакцию от списка транзакций
			//log.Debug("++p.BinaryData=%x\n", p.BinaryData)
			//log.Debug("transactionSize", transactionSize)
			transactionBinaryData := utils.BytesShift(&p.BinaryData, transactionSize)
			transactionBinaryDataFull := transactionBinaryData
			//ioutil.WriteFile("/tmp/dctx", transactionBinaryDataFull, 0644)
			//ioutil.WriteFile("/tmp/dctxhash", utils.Md5(transactionBinaryDataFull), 0644)
			// добавляем взятую тр-ию в набор тр-ий для RollbackTo, в котором пойдем в обратном порядке
			txForRollbackTo = append(txForRollbackTo, utils.EncodeLengthPlusData(transactionBinaryData)...)
			//log.Debug("transactionBinaryData: %x\n", transactionBinaryData)
			//log.Debug("txForRollbackTo: %x\n", txForRollbackTo)

			err = p.CheckLogTx(transactionBinaryDataFull)
			if err != nil {
				p.RollbackTo(txForRollbackTo, true, false)
				return utils.ErrInfo(err)
			}

			utils.WriteSelectiveLog("UPDATE transactions SET used=1 WHERE hex(hash) = "+string(utils.Md5(transactionBinaryDataFull)))
			affect, err := p.ExecSqlGetAffect("UPDATE transactions SET used=1 WHERE hex(hash) = ?", utils.Md5(transactionBinaryDataFull))
			if err != nil {
				utils.WriteSelectiveLog(err)
				utils.WriteSelectiveLog("RollbackTo")
				p.RollbackTo(txForRollbackTo, true, false)
				return utils.ErrInfo(err)
			}
			utils.WriteSelectiveLog("affect: "+utils.Int64ToStr(affect))
			//log.Debug("transactionBinaryData", transactionBinaryData)
			p.TxHash = utils.Md5(transactionBinaryData)
			log.Debug("p.TxHash %s", p.TxHash)
			p.TxSlice, err = p.ParseTransaction(&transactionBinaryData)
			log.Debug("p.TxSlice %v", p.TxSlice)
			if err != nil {
				p.RollbackTo(txForRollbackTo, true, false)
				return err
			}

			if p.BlockData.BlockId > 1 {
				var userId int64
				// txSlice[3] могут подсунуть пустой
				if len(p.TxSlice) > 3 {
					if !utils.CheckInputData(p.TxSlice[3], "int64") {
						return utils.ErrInfo(fmt.Errorf("empty user_id"))
					} else {
						userId = utils.BytesToInt64(p.TxSlice[3])
					}
				} else {
					return utils.ErrInfo(fmt.Errorf("empty user_id"))
				}

				// считаем по каждому юзеру, сколько в блоке от него транзакций
				txCounter[userId]++

				// чтобы 1 юзер не смог прислать дос-блок размером в 10гб, который заполнит своими же транзакциями
				if txCounter[userId] > p.Variables.Int64["max_block_user_transactions"] {
					p.RollbackTo(txForRollbackTo, true, false)
					return utils.ErrInfo(fmt.Errorf("max_block_user_transactions"))
				}
			}

			// время в транзакции не может быть больше, чем на MAX_TX_FORW сек времени блока
			// и  время в транзакции не может быть меньше времени блока -24ч.
			if utils.BytesToInt64(p.TxSlice[2])-consts.MAX_TX_FORW > p.BlockData.Time || utils.BytesToInt64(p.TxSlice[2]) < p.BlockData.Time-consts.MAX_TX_BACK {
				p.RollbackTo(txForRollbackTo, true, false)
				return utils.ErrInfo(fmt.Errorf("incorrect transaction time"))
			}

			// проверим, есть ли такой тип тр-ий
			_, ok := consts.TxTypes[utils.BytesToInt(p.TxSlice[1])]
			if !ok {
				return utils.ErrInfo(fmt.Errorf("nonexistent type"))
			}

			p.TxMap = map[string][]byte{}

			// для статы
			p.TxIds = append(p.TxIds, string(p.TxSlice[1]))

			MethodName := consts.TxTypes[utils.BytesToInt(p.TxSlice[1])]
			log.Debug("MethodName", MethodName+"Init")
			err_ := utils.CallMethod(p, MethodName+"Init")
			if _, ok := err_.(error); ok {
				log.Error("error: %v", err)
				return utils.ErrInfo(err_.(error))
			}

			log.Debug("MethodName", MethodName+"Front")
			err_ = utils.CallMethod(p, MethodName+"Front")
			if _, ok := err_.(error); ok {
				log.Error("error: %v", err_)
				p.RollbackTo(txForRollbackTo, true, false)
				return utils.ErrInfo(err_.(error))
			}

			log.Debug("MethodName", MethodName)
			err_ = utils.CallMethod(p, MethodName)
			if _, ok := err_.(error); ok {
				log.Error("error: %v", err)
				return utils.ErrInfo(err_.(error))
			}

			// даем юзеру понять, что его тр-ия попала в блок
			p.ExecSql("UPDATE transactions_status SET block_id = ? WHERE hex(hash) = ?", p.BlockData.BlockId, utils.Md5(transactionBinaryDataFull))

			// Тут было time(). А значит если бы в цепочке блоков были блоки в которых были бы одинаковые хэши тр-ий, то ParseDataFull вернул бы error
			err = p.InsertInLogTx(transactionBinaryDataFull, utils.BytesToInt64(p.TxMap["time"]))
			if err != nil {
				return utils.ErrInfo(err)
			}

			if len(p.BinaryData) == 0 {
				break
			}
		}
	}

	p.UpdBlockInfo()

	return nil
}