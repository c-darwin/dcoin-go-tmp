package dcparser
import (
	//"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"utils"
	"os"
	"time"
	"consts"
	"sync"
)

type Parser struct {
	*utils.DCDB
	TxMap map[string][]byte
	BlockData *utils.BlockData
	PrevBlock *utils.BlockData
	BinaryData []byte
	blockHashHex []byte
	dataType int64
	blockHex []byte
	Variables map[string]string
	CurrentBlockId int64
	fullTxBinaryData []byte
	halfRollback bool // уже не актуально, т.к. нет ни одной половинной фронт-проверки
	TxHash []byte
	TxSlice [][]byte
	MerkleRoot []byte
	GoroutineName string
	CurrentVersion string
}
/*
func TypeArray (txType string) int32 {
	for k, v := range x {
		if v == txType {
			return int32(k)
		}
	}
	return 0
}*/


func (p *Parser) dataPre() {
	p.blockHashHex = utils.DSha256(p.BinaryData)
	p.blockHex = utils.BinToHex(p.BinaryData)
	// определим тип данных
	p.dataType =  utils.BinToDec(utils.BytesShift(&p.BinaryData, 1))
	fmt.Println("dataType", p.dataType)
}



func (p *Parser) ParseBlock() error {
	/*
	Заголовок (от 143 до 527 байт )
	TYPE (0-блок, 1-тр-я)     1
	BLOCK_ID   				       4
	TIME       					       4
	USER_ID                         5
	LEVEL                              1
	SIGN                               от 128 до 512 байт. Подпись от TYPE, BLOCK_ID, PREV_BLOCK_HASH, TIME, USER_ID, LEVEL, MRKL_ROOT
	Далее - тело блока (Тр-ии)
	*/
	p.BlockData = utils.ParseBlockHeader(&p.BinaryData)

	p.CurrentBlockId = p.BlockData.BlockId
	fmt.Println(p.BlockData)

	return nil
}



func (p *Parser) CheckBlockHeader() error {
	var err error
	// инфа о предыдущем блоке (т.е. последнем занесенном)
	if p.PrevBlock == nil {
		p.PrevBlock, err = p.DCDB.GetBlockDataFromBlockChain(p.BlockData.BlockId-1)
		fmt.Println("PrevBlock 0",p.PrevBlock)
		if err != nil {
			return utils.ErrInfo(err)
		}
	}
	fmt.Println("PrevBlock",p.PrevBlock)
	fmt.Println("p.PrevBlock.BlockId",p.PrevBlock.BlockId)
	// для локальных тестов
	if p.PrevBlock.BlockId == 1 {
		if p.GetConfigIni("start_block_id") != "" {
				p.PrevBlock.BlockId = utils.StrToInt64(p.GetConfigIni("start_block_id"))
		}
	}

	var first bool
	if p.BlockData.BlockId == 1 {
		p.Variables["max_tx_size"] = "1048576"
		first = true
	} else {
		first = false
	}
	fmt.Println(first)

	// меркель рут нужен для проверки подписи блока, а также проверки лимитов MAX_TX_SIZE и MAX_TX_COUNT
	mrklRoot := utils.GetMrklroot(p.BinaryData, p.Variables, first)
	fmt.Println(mrklRoot)

	// проверим время
	if !utils.CheckInputData(p.BlockData.Time, "int") {
		fmt.Println("p.BlockData.Time",p.BlockData.Time)
		return utils.ErrInfo(fmt.Errorf("incorrect time"))
	}
	// проверим уровень
	if !utils.CheckInputData(p.BlockData.Level, "level") {
		return utils.ErrInfo(fmt.Errorf("incorrect level"))
	}

	// получим значения для сна
	sleepData, err:=p.DCDB.GetSleepData()
	if err != nil {
		return utils.ErrInfo(err)
	}

	// узнаем время, которые было затрачено в ожидании is_ready предыдущим блоком
	isReadySleep := p.DCDB.GetIsReadySleep(p.PrevBlock.Level, sleepData["is_ready"])
	fmt.Println("isReadySleep", isReadySleep)

	// сколько сек должен ждать нод, перед тем, как начать генерить блок, если нашел себя в одном из уровней.
	generatorSleep := utils.GetGeneratorSleep(p.BlockData.Level, sleepData["generator"])
	fmt.Println("generatorSleep", generatorSleep)

	// сумма is_ready всех предыдущих уровней, которые не успели сгенерить блок
	isReadySleep2 := utils.GetIsReadySleepSum(p.BlockData.Level , sleepData["is_ready"])
	fmt.Println("isReadySleep2", isReadySleep2)

	// не слишком ли рано прислан этот блок. допустима погрешность = error_time
	if !first {
		if p.PrevBlock.Time + isReadySleep + generatorSleep + isReadySleep2 - p.BlockData.Time > utils.StrToInt64(p.Variables["error_time"]) {
			return utils.ErrInfo(fmt.Errorf("incorrect block time %d + %d + %d+ %d - %d > %d", p.PrevBlock.Time, isReadySleep, generatorSleep, isReadySleep2, p.BlockData.Time, p.Variables["error_time"]))
		}
	}

	// исключим тех, кто сгенерил блок с бегущими часами
	if p.BlockData.Time > time.Now().Unix() {
		utils.ErrInfo(fmt.Errorf("incorrect block time"))
	}

	// проверим ID блока
	if !utils.CheckInputData(p.BlockData.BlockId, "int") {
		return utils.ErrInfo(fmt.Errorf("incorrect block_id"))
	}

	// проверим, верный ли ID блока
	if !first {
		if p.BlockData.BlockId != p.PrevBlock.BlockId+1 {
			return utils.ErrInfo(fmt.Errorf("incorrect block_id %d != %d +1", p.BlockData.BlockId, p.PrevBlock.BlockId))
		}
	}

	// проверим, есть ли такой майнер и заодно получим public_key
	nodePublicKey, err := p.DCDB.GetNodePublicKey(p.BlockData.UserId)
	if err != nil {
		return utils.ErrInfo(err)
	}

	if !first {
		if len(nodePublicKey)==0 {
			return utils.ErrInfo(fmt.Errorf("empty nodePublicKey"))
		}
		// SIGN от 128 байта до 512 байт. Подпись от TYPE, BLOCK_ID, PREV_BLOCK_HASH, TIME, USER_ID, LEVEL, MRKL_ROOT
		forSign := fmt.Sprintf("0,%d,%s,%d,%d,%d,%s", p.BlockData.BlockId, p.PrevBlock.Hash, p.BlockData.Time, p.BlockData.UserId, p.BlockData.Level, mrklRoot)
		fmt.Println(forSign)
		// проверим подпись
		resultCheckSign, err := utils.CheckSign([][]byte{nodePublicKey}, forSign, p.BlockData.Sign, true);
		if err != nil {
			return utils.ErrInfo(err)
		}
		if !resultCheckSign {
			return utils.ErrInfo(fmt.Errorf("incorrect signature"))
		}
	}
	return nil
}

// Это защита от dos, когда одну транзакцию можно было бы послать миллион раз,
// и она каждый раз успешно проходила бы фронтальную проверку
func (p *Parser) CheckLogTx(tx_binary []byte) error {
	hash, err := p.DCDB.Single(`SELECT hash FROM log_transactions WHERE hash = [hex]`, utils.Md5(tx_binary))
	if err != nil {
		return utils.ErrInfo(err)
	}
	if len(hash) > 0 {
		return utils.ErrInfo(fmt.Errorf("double log_transactions"))
	}
	return nil
}

//  если в ходе проверки тр-ий возникает ошибка, то вызываем откатчик всех занесенных тр-ий
func (p *Parser) RollbackTo (binaryData []byte, skipCurrent bool, onlyFront bool) error {
	var err error
	if len(binaryData) > 0 {
		// вначале нужно получить размеры всех тр-ий, чтобы пройтись по ним в обратном порядке
		binForSize := binaryData
		var sizesSlice []int64
		for {
			txSize := utils.DecodeLength(&binForSize)
			if (txSize == 0) {
				break
			}
			sizesSlice = append(sizesSlice, txSize)
			// удалим тр-ию
			fmt.Println("txSize", txSize)
			//fmt.Println("binForSize", binForSize)
			utils.BytesShift(&binForSize, txSize)
			if len(binForSize) == 0 {
				break
			}
		}
		sizesSlice = utils.SliceReverse(sizesSlice)
		for i:=0; i<len(sizesSlice); i++ {
			// обработка тр-ий может занять много времени, нужно отметиться
			p.DCDB.UpdDaemonTime(p.GoroutineName)
			// отчекрыжим одну транзакцию
			transactionBinaryData := utils.BytesShiftReverse(&binaryData, sizesSlice[i])
			transactionBinaryData_ := transactionBinaryData
			// узнаем кол-во байт, которое занимает размер
			size_ := len(utils.EncodeLength(sizesSlice[i]))
			// удалим размер
			utils.BytesShiftReverse(&binaryData, size_)
			p.TxHash = utils.Md5(transactionBinaryData)
			p.TxSlice, err = p.ParseTransaction(&transactionBinaryData)
			if err != nil {
				return utils.ErrInfo(err)
			}
			MethodName := consts.TxTypes[utils.BytesToInt(p.TxSlice[1])]
			p.TxMap = map[string][]byte{}
			err_ := utils.CallMethod(p, MethodName+"_init")
			if _, ok := err_.(error); ok {
				return utils.ErrInfo(err_.(error))
			}

			// если дошли до тр-ии, которая вызвала ошибку, то откатываем только фронтальную проверку
			if i == 0 {
				if skipCurrent { // тр-ия, которая вызвала ошибку закончилась еще до фронт. проверки, т.е. откатывать по ней вообще нечего
					continue
				}
				// если успели дойти только до половины фронтальной функции
				MethodNameRollbackFront := ""
				if p.halfRollback {
					MethodNameRollbackFront = MethodName+"_rollback_front_0"
				} else {
					MethodNameRollbackFront = MethodName+"_rollback_front"
				}
				// откатываем только фронтальную проверку
				err_ = utils.CallMethod(p, MethodNameRollbackFront)
				if _, ok := err_.(error); ok {
					return utils.ErrInfo(err_.(error))
				}
			} else if onlyFront {
				err_ = utils.CallMethod(p, MethodName+"_rollback_front")
				if _, ok := err_.(error); ok {
					return utils.ErrInfo(err_.(error))
				}
			} else {
				err_ = utils.CallMethod(p, MethodName+"_rollback_front")
				if _, ok := err_.(error); ok {
					return utils.ErrInfo(err_.(error))
				}
				err_ = utils.CallMethod(p, MethodName+"_rollback")
				if _, ok := err_.(error); ok {
					return utils.ErrInfo(err_.(error))
				}
			}
			p.DCDB.DelLogTx(transactionBinaryData_)

			// =================== ради эксперимента =========
			if onlyFront {
				_, err = p.DCDB.ExecSql("UPDATE transactions SET verified = 0 WHERE hash = [hex]", p.TxHash)
				if err != nil {
					return utils.ErrInfo(err)
				}
			} else { // ====================================
				_, err = p.DCDB.ExecSql("UPDATE transactions SET used = 0 WHERE hash = [hex]", p.TxHash)
				if err != nil {
					return utils.ErrInfo(err)
				}
			}
		}
	}
	return err
}

func (p *Parser) ParseTransaction (transactionBinaryData *[]byte) ([][]byte, error) {

	var returnSlice [][]byte
	var transSlice [][]byte
	var merkleSlice [][]byte

	if  len(*transactionBinaryData) > 0 {

		// хэш транзакции
		transSlice = append(transSlice, utils.DSha256(*transactionBinaryData))

		// первый байт - тип транзакции
		transSlice = append(transSlice, utils.Int64ToByte(utils.BinToDecBytesShift(transactionBinaryData, 1)))
		if len(*transactionBinaryData) == 0 {
			return transSlice, utils.ErrInfo(fmt.Errorf("incorrect tx"))
		}

		// следующие 4 байта - время транзакции
		transSlice = append(transSlice, utils.Int64ToByte(utils.BinToDecBytesShift(transactionBinaryData, 4)))
		if len(*transactionBinaryData) == 0 {
			return transSlice, utils.ErrInfo(fmt.Errorf("incorrect tx"))
		}

		// преобразуем бинарные данные транзакции в массив
		i:=0
		for {
			length := utils.DecodeLength(transactionBinaryData)
			if length > 0 && length < utils.StrToInt64(p.Variables["max_tx_size"]) {
				data := utils.BytesShift(transactionBinaryData, length)
				returnSlice = append(returnSlice, data)
				merkleSlice = append(merkleSlice, utils.DSha256(data))
			}
			i++
			if length == 0 || i >= 20 { // у нас нет тр-ий с более чем 20 элементами
				break
			}
		}
		if len(*transactionBinaryData) > 0 {
			return transSlice, utils.ErrInfo(fmt.Errorf("incorrect transactionBinaryData"))
		}
	} else {
		merkleSlice = append(merkleSlice, []byte("0"))
	}
	p.MerkleRoot = utils.MerkleTreeRoot(merkleSlice)
	return append(transSlice, returnSlice...), nil
}

func (p *Parser) InsertIntoBlockchain() {
	var mutex = &sync.Mutex{}
	// для локальных тестов
	if p.BlockData.BlockId == 1 {
		if p.GetConfigIni("start_block_id") != "" {
			p.BlockData.BlockId = utils.StrToInt64(p.GetConfigIni("start_block_id"))
		}
	}
	mutex.Lock()
	// пишем в цепочку блоков
	p.DCDB.ExecSql("DELETE FROM block_chain WHERE id = ?", p.BlockData.BlockId)
	p.DCDB.ExecSql("INSERT INTO block_chain (id, hash, head_hash, data) VALUES (?, [hex],[hex],[hex])",
		p.BlockData.BlockId, p.BlockData.Hash, p.BlockData.HeadHash, p.blockHex)
	mutex.Unlock()
}
/*public function insert_into_blockchain()
	{
		if ($AffectedRows<1) {

			debug_print(">>>>>>>>>>> BUG LOAD DATA LOCAL INFILE  '{$file}' IGNORE INTO TABLE block_chain", __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__);

			$row = $this->db->query( __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__, "
							SELECT *
							FROM `".DB_PREFIX."block_chain`
							WHERE `id` = {$this->block_data['block_id']}
							", 'fetch_array');

			print_r_hex($row);

			// ========================= временно для поиска бага: ====================================

			$this->db->query( __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__, "
			LOAD DATA LOCAL INFILE  '{$file}' REPLACE INTO TABLE `".DB_PREFIX."block_chain`
			FIELDS TERMINATED BY '\t'
			(`id`, @hash, @head_hash, @data)
			SET `hash` = UNHEX(@hash),
				   `head_hash` = UNHEX(@head_hash),
				   `data` = UNHEX(@data)
			");

			//print 'getAffectedRows='.$this->db->getAffectedRows()."\n";
			// =================================================================================
		}
		unlink($file);
	}*/

/**
	фронт. проверка + занесение данных из блока в таблицы и info_block
*/
func (p *Parser) ParseDataFull() error {

	p.dataPre()
	if p.dataType != 0  { // парсим только блоки
		return utils.ErrInfo(fmt.Errorf("incorrect dataType"))
	}
	var err error

	p.Variables, err = p.DCDB.GetAllVariables()
	if err != nil {
		return utils.ErrInfo(err)
	}
	err = p.ParseBlock()
	if err != nil {
		return utils.ErrInfo(err)
	}

	// проверим данные, указанные в заголовке блока
	err = p.CheckBlockHeader()
	if err != nil {
		return utils.ErrInfo(err)
	}


	p.DCDB.Single("DELETE FROM transactions WHERE used = 1")

	txCounter := make(map[int64]int64)
	p.fullTxBinaryData = p.BinaryData
	var txForRollbackTo []byte
	if len(p.BinaryData) > 0 {
		for {
			// обработка тр-ий может занять много времени, нужно отметиться
			p.DCDB.UpdDaemonTime(p.GoroutineName)
			p.halfRollback = false
			fmt.Println("&p.BinaryData", p.BinaryData)
			transactionSize := utils.DecodeLength(&p.BinaryData)
			if len(p.BinaryData) == 0 {
				return utils.ErrInfo(fmt.Errorf("empty BinaryData"))
			}


			// отчекрыжим одну транзакцию от списка транзакций
			fmt.Printf("++p.BinaryData=%x\n", p.BinaryData)
			fmt.Println("transactionSize", transactionSize)
			transactionBinaryData := utils.BytesShift(&p.BinaryData, transactionSize)
			transactionBinaryDataFull := transactionBinaryData

			// добавляем взятую тр-ию в набор тр-ий для RollbackTo, в котором пойдем в обратном порядке
			txForRollbackTo = append(txForRollbackTo, utils.EncodeLengthPlusData(transactionBinaryData)...)

			err = p.CheckLogTx(transactionBinaryDataFull)
			if err != nil {
				fmt.Println("err", err)
				fmt.Println("RollbackTo")
				p.RollbackTo(txForRollbackTo, true, false);
				return err
			}

			p.DCDB.ExecSql("UPDATE transactions SET used=1 WHERE hash = [hex]", utils.Md5(transactionBinaryDataFull))
			//fmt.Println("transactionBinaryData", transactionBinaryData)
			p.TxHash = utils.Md5(transactionBinaryData)
			fmt.Println("p.TxHash", p.TxHash)
			p.TxSlice, err = p.ParseTransaction(&transactionBinaryData)
			fmt.Println("p.TxSlice", p.TxSlice)
			if err !=nil {
				fmt.Println("err", err)
				fmt.Println("RollbackTo")
				p.RollbackTo (txForRollbackTo, true, false)
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
				if txCounter[userId] > utils.StrToInt64(p.Variables["max_block_user_transactions"])  {
					p.RollbackTo(txForRollbackTo, true, false)
					return utils.ErrInfo(fmt.Errorf("max_block_user_transactions"))
				}
			}

			// время в транзакции не может быть больше, чем на MAX_TX_FORW сек времени блока
			// и  время в транзакции не может быть меньше времени блока -24ч.
			if utils.BytesToInt64(p.TxSlice[2]) - consts.MAX_TX_FORW > p.BlockData.Time || utils.BytesToInt64(p.TxSlice[2]) < p.BlockData.Time - consts.MAX_TX_BACK {
				p.RollbackTo(txForRollbackTo, true, false)
				return utils.ErrInfo(fmt.Errorf("incorrect transaction time"))
			}

			// проверим, есть ли такой тип тр-ий
			_, ok := consts.TxTypes[utils.BytesToInt(p.TxSlice[1])]
			if (!ok) {
				return utils.ErrInfo(fmt.Errorf("nonexistent type"))
			}

			p.TxMap = map[string][]byte{}

			MethodName := consts.TxTypes[utils.BytesToInt(p.TxSlice[1])]

			err_ := utils.CallMethod(p,MethodName+"_init")
			if _, ok := err_.(error); ok {
				return utils.ErrInfo(err_.(error))
			}

			err_ = utils.CallMethod(p,MethodName+"_front")
			if _, ok := err_.(error); ok {
				p.RollbackTo(txForRollbackTo, true, false);
				return utils.ErrInfo(err_.(error))
			}

			err_ = utils.CallMethod(p,MethodName)
			if _, ok := err_.(error); ok {
				return utils.ErrInfo(err_.(error))
			}

			// даем юзеру понять, что его тр-ия попала в блок
			p.DCDB.ExecSql("UPDATE transactions_status SET block_id = ? WHERE hash = [hex]", p.BlockData.BlockId, utils.Md5(transactionBinaryDataFull))

			// Тут было time(). А значит если бы в цепочке блоков были блоки в которых были бы одинаковые хэши тр-ий, то ParseDataFull вернул бы error
			err = p.DCDB.InsertInLogTx(transactionBinaryDataFull, utils.BytesToInt64(p.TxMap["time"]))
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

func (p *Parser) UpdBlockInfo() {
	blockId := p.BlockData.BlockId
	// для локальных тестов
	if p.BlockData.BlockId == 1 {
		if p.GetConfigIni("start_block_id") != "" {
			blockId = utils.StrToInt64(p.GetConfigIni("start_block_id"))
		}
	}
	headHashData := fmt.Sprintf("%d,%d,%s", p.BlockData.UserId, blockId, p.PrevBlock.HeadHash)
	p.BlockData.HeadHash = utils.DSha256(headHashData)
	forSha := fmt.Sprintf("%d,%s,%s,%d,%d,%d", blockId, p.PrevBlock.Hash, p.MerkleRoot, p.BlockData.Time, p.BlockData.UserId, p.BlockData.Level)
	p.BlockData.Hash = utils.DSha256(forSha)

	if p.BlockData.BlockId == 1 {
		p.DCDB.ExecSql("INSERT INTO info_block (hash, head_hash, block_id, time, level, current_version) VALUES ([hex], [hex], ?, ?, ?, ?)",
			p.BlockData.Hash, p.BlockData.HeadHash, blockId, p.BlockData.Time, p.BlockData.Level, p.CurrentVersion)
	} else {
		p.DCDB.ExecSql("UPDATE info_block SET hash = [hex], head_hash = [hex], block_id = ?, time = ?, level = ?, sent = 0",
			p.BlockData.Hash, p.BlockData.HeadHash, blockId, p.BlockData.Time, p.BlockData.Level)
		p.DCDB.ExecSql("UPDATE config SET my_block_id = ? WHERE my_block_id < ?", blockId, blockId)
	}
}

func (p *Parser) GetTxMap(fields []string) (map[string][]byte, error) {
	//fmt.Println("p.TxSlice", p.TxSlice)
	//fmt.Println("fields", fields)
	if len(p.TxSlice) != len(fields)+4 {
		return nil, fmt.Errorf("bad transaction_array %d != %d (type=%d)",  len(p.TxSlice),  len(fields)+4, p.TxSlice[0])
	}
	TxMap := make(map[string][]byte)
	TxMap["hash"] = p.TxSlice[0]
	TxMap["type"] = p.TxSlice[1]
	TxMap["time"] = p.TxSlice[2]
	TxMap["user_id"] = p.TxSlice[3]
	for i, field := range fields {
		TxMap[field] = p.TxSlice[i+4]
	}
	//fmt.Println("TxMap", TxMap)
	//fmt.Println("TxMap[hash]", TxMap["hash"])
	//fmt.Println("p.TxSlice[0]", p.TxSlice[0])
	return TxMap, nil
}

func (p *Parser) GetMyUserId(userId int64) (int64, int64, string, []int64, error) {
	var myUserId int64
	var myPrefix string
	var myUserIds []int64
	var myBlockId int64
	collective, err := p.GetCommunityUsers()
	if len(collective) > 0 {// если работаем в пуле
		myUserIds = collective
		// есть ли юзер, который задействован среди юзеров нашего пула
		if utils.InSliceInt64(userId, collective) {
			myPrefix = fmt.Sprintf("%d_", userId)
			// чтобы не было проблем с change_primary_key нужно получить user_id только тогда, когда он был реально выдан
			// в будущем можно будет переделать, чтобы user_id можно было указывать всем и всегда заранее.
			// тогда при сбросе будут собираться более полные таблы my_, а не только те, что заполнятся в change_primary_key
			myUserId, err = p.SingleInt64("SELECT user_id FROM "+myPrefix+"my_table")
			if err != nil {
				return myUserId, myBlockId, myPrefix, myUserIds, err
			}
		}
	} else {
		myUserId, err = p.SingleInt64("SELECT user_id FROM my_table")
		if err != nil {
			return myUserId, myBlockId, myPrefix, myUserIds, err
		}
		myUserIds = append(myUserIds, myUserId)
	}
	myBlockId, err = p.SingleInt64("SELECT my_block_id FROM config")
	if err != nil {
		return myUserId, myBlockId, myPrefix, myUserIds, err
	}
	return myUserId, myBlockId, myPrefix, myUserIds, nil
}

func MakeTest(parser *Parser, txType string, hashesStart map[string]string) error {
	//fmt.Println("dcparser."+txType+"Init")
	err := utils.CallMethod(parser, txType+"Init")
	//fmt.Println(err)

	if i, ok := err.(error); ok {
		fmt.Println(err.(error), i)
		return err.(error)
	}

	if len(os.Args)==1 {
		err = utils.CallMethod(parser, txType)
		if i, ok := err.(error); ok {
			fmt.Println(err.(error), i)
			return err.(error)
		}

		//fmt.Println("-------------------")
		// узнаем, какие таблицы были затронуты в результате выполнения основного метода
		hashesMiddle, err := parser.AllHashes()
		if err != nil {
			return utils.ErrInfo(err)
		}
		var tables []string
		for table, hash := range hashesMiddle {
			if hash!=hashesStart[table] {
				tables = append(tables, table)
			}
		}
		fmt.Println(tables)

		// rollback
		err0 := utils.CallMethod(parser, txType+"Rollback")
		if i, ok := err0.(error); ok {
			fmt.Println(err0.(error), i)
			return err0.(error)
		}

		// сраниим хэши, которые были до начала и те, что получились после роллбэка
		hashesEnd, err := parser.AllHashes()
		if err != nil {
			return utils.ErrInfo(err)
		}
		for table, hash := range hashesEnd {
			if hash!=hashesStart[table] {
				fmt.Println("ERROR in table ", table)
			}
		}

	} else if os.Args[1] == "w" {
		err = utils.CallMethod(parser, txType)
		if i, ok := err.(error); ok {
			fmt.Println(err.(error), i)
			return err.(error)
		}
	} else if os.Args[1] == "r" {
		err = utils.CallMethod(parser, txType+"Rollback")
		if i, ok := err.(error); ok {
			fmt.Println(err.(error), i)
			return err.(error)
		}
	}
	return nil
}

// откатываем ID на кол-во затронутых строк
func (p *Parser) rollbackAI(table string, num int64) (error) {
	//fmt.Println("table", table)
	current, err := p.Single("SELECT id FROM "+table+" ORDER BY id DESC LIMIT 1", )
	if err != nil {
		return utils.ErrInfo(err)
	}

	pg_get_serial_sequence, err := p.Single("SELECT pg_get_serial_sequence('"+table+"', 'id')")
	if err != nil {
		return utils.ErrInfo(err)
	}

	_, err = p.ExecSql("ALTER SEQUENCE "+pg_get_serial_sequence+" RESTART WITH "+utils.Int64ToStr(utils.StrToInt64(current)+num))
	if err != nil {
		return utils.ErrInfo(err)
	}
	return err
}
