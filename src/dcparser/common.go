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
	"reflect"
	"math"
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
	Variables *utils.Variables
	CurrentBlockId int64
	fullTxBinaryData []byte
	halfRollback bool // уже не актуально, т.к. нет ни одной половинной фронт-проверки
	TxHash []byte
	TxSlice [][]byte
	MerkleRoot []byte
	GoroutineName string
	CurrentVersion string
	MrklRoot []byte
	PublicKeys [][]byte
	AdminUserId int64
	TxUserID int64
	TxTime int64
	nodePublicKey []byte
}

func (p *Parser) limitRequest(limit_ interface{}, txType string, period_ interface{}) error {

	var limit int
	switch limit_.(type) {
	case string:
		limit = utils.StrToInt(limit_.(string))
	case int:
		limit = limit_.(int)
	case int64:
		limit = int(limit_.(int64))
	}

	var period int
	switch period_.(type) {
	case string:
		period = utils.StrToInt(period_.(string))
	case int:
		period = period_.(int)
	}

	time := utils.BytesToInt(p.TxMap["time"])
	num, err := p.DCDB.Single("SELECT count(time) FROM log_time_"+txType+" WHERE user_id = ? AND time > ?", p.TxMap["user_id"], (time-period)).Int()
	if err != nil {
		return err
	}
	if num >= limit {
		return utils.ErrInfo(fmt.Errorf("[limit_requests] log_time_%v %v >= %v", txType, num, limit))
	} else {
		err := p.DCDB.ExecSql("INSERT INTO log_time_"+txType+" (user_id, time) VALUES (?, ?)", p.TxMap["user_id"], time)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) getAdminUserId() error {
	AdminUserId, err := p.DCDB.Single("SELECT user_id FROM admin").Int64()
	if err != nil {
		return utils.ErrInfo(err)
	}
	p.AdminUserId = AdminUserId
	return nil
}
func (p *Parser) checkMinerNewbie() error {
	var time int64
	if p.BlockData != nil {
		time = p.BlockData.Time
	} else {
		time = utils.BytesToInt64(p.TxMap["time"])
	}
	regTime, err := p.DCDB.Single("SELECT reg_time FROM miners_data WHERE user_id = ?", p.TxMap["user_id"]).Int64()
	err = p.getAdminUserId()
	if err != nil {
		return utils.ErrInfo(err)
	}
	if (p.BlockData==nil) || (p.BlockData!=nil && p.BlockData.BlockId > 29047) {
		if regTime > (time - p.Variables.Int64["miner_newbie_time"]) && utils.BytesToInt64(p.TxMap["user_id"]) != p.AdminUserId {
			return utils.ErrInfo(fmt.Errorf("error miner_newbie (%v > %v - %v)", regTime, time, p.Variables.Int64["miner_newbie_time"]))
		}
	}
	return nil
}


func (p *Parser) checkMiner(userId []byte) error {
	// в cash_request_out передается to_user_id
	var blockId int64
	addSql := ""
	// если разжаловали в этом блоке, то считаем всё еще майнером
	if p.BlockData!=nil {
		blockId = p.BlockData.BlockId
		addSql = " OR `ban_block_id`= "+utils.Int64ToStr(blockId)
	}

	// когда админ разжаловывает майнера, у него пропадет miner_id
	minerId, err := p.DCDB.Single("SELECT miner_id FROM miners_data WHERE user_id = ? AND (miner_id>0 "+addSql+")", userId).Int64()
	if err != nil {
		return err
	}
	// если есть бан в этом же блоке, то будет miner_id = 0, но условно считаем, что проверка пройдена
	if (minerId > 0) || (minerId == 0 && blockId > 0) {
		return nil
	} else {
		return utils.ErrInfoFmt("incorrect miner id")
	}
}

// общая проверка для всех _front
func (p *Parser) generalCheck() error {
	if !utils.CheckInputData(p.TxMap["user_id"], "int64") {
		return utils.ErrInfoFmt("incorrect user_id")
	}
	if !utils.CheckInputData(p.TxMap["time"], "int") {
		return utils.ErrInfoFmt("incorrect time")
	}
	// проверим, есть ли такой юзер и заодно получим public_key
	data, err := p.DCDB.OneRow("SELECT public_key_0,public_key_1,public_key_2	FROM users WHERE user_id = ?", p.TxMap["user_id"]).String()
	if err != nil {
		return utils.ErrInfo(err)
	}
	if len(data["public_key_0"])==0 {
		return utils.ErrInfoFmt("incorrect user_id")
	}
	p.PublicKeys = append(p.PublicKeys, []byte(data["public_key_0"]))
	if len(data["public_key_1"]) > 0 {
		p.PublicKeys = append(p.PublicKeys, []byte(data["public_key_1"]))
	}
	if len(data["public_key_2"]) > 0 {
		p.PublicKeys = append(p.PublicKeys, []byte(data["public_key_2"]))
	}
	// чтобы не записали слишком длинную подпись
	// 128 - это нод-ключ
	if len(p.TxMap["sign"]) < 128 || len(p.TxMap["sign"]) > 5000  {
		return utils.ErrInfoFmt("incorrect sign size")
	}
	return nil
}

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
		p.Variables.Int64["max_tx_size"] = 1048576
		first = true
	} else {
		first = false
	}
	fmt.Println(first)

	// меркель рут нужен для проверки подписи блока, а также проверки лимитов MAX_TX_SIZE и MAX_TX_COUNT
	fmt.Println(p.Variables)
	p.MrklRoot, err = utils.GetMrklroot(p.BinaryData, p.Variables, first)
	fmt.Println(p.MrklRoot)
	if err != nil {
		return utils.ErrInfo(err)
	}

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
		if p.PrevBlock.Time + isReadySleep + generatorSleep + isReadySleep2 - p.BlockData.Time > p.Variables.Int64["error_time"] {
			return utils.ErrInfo(fmt.Errorf("incorrect block time %d + %d + %d+ %d - %d > %d", p.PrevBlock.Time, isReadySleep, generatorSleep, isReadySleep2, p.BlockData.Time, p.Variables.Int64["error_time"]))
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
		forSign := fmt.Sprintf("0,%d,%s,%d,%d,%d,%s", p.BlockData.BlockId, p.PrevBlock.Hash, p.BlockData.Time, p.BlockData.UserId, p.BlockData.Level, p.MrklRoot)
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
	hash, err := p.DCDB.Single(`SELECT hash FROM log_transactions WHERE hash = [hex]`, utils.Md5(tx_binary)).String()
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
				err = p.DCDB.ExecSql("UPDATE transactions SET verified = 0 WHERE hash = [hex]", p.TxHash)
				if err != nil {
					return utils.ErrInfo(err)
				}
			} else { // ====================================
				err = p.DCDB.ExecSql("UPDATE transactions SET used = 0 WHERE hash = [hex]", p.TxHash)
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
			fmt.Printf("length%d\n", length)
			if length > 0 && length < p.Variables.Int64["max_tx_size"] {
				data := utils.BytesShift(transactionBinaryData, length)
				returnSlice = append(returnSlice, data)
				merkleSlice = append(merkleSlice, utils.DSha256(data))
				fmt.Printf("utils.DSha256(data) %s\n", utils.DSha256(data))
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
	fmt.Println("merkleSlice", merkleSlice)
	p.MerkleRoot = utils.MerkleTreeRoot(merkleSlice)
	fmt.Printf("MerkleRoot %s\n", p.MerkleRoot)
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
	fmt.Println("p.Variables", p.Variables)
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
			//fmt.Printf("++p.BinaryData=%x\n", p.BinaryData)
			//fmt.Println("transactionSize", transactionSize)
			transactionBinaryData := utils.BytesShift(&p.BinaryData, transactionSize)
			transactionBinaryDataFull := transactionBinaryData

			// добавляем взятую тр-ию в набор тр-ий для RollbackTo, в котором пойдем в обратном порядке
			txForRollbackTo = append(txForRollbackTo, utils.EncodeLengthPlusData(transactionBinaryData)...)
			//fmt.Printf("transactionBinaryData: %x\n", transactionBinaryData)
			//fmt.Printf("txForRollbackTo: %x\n", txForRollbackTo)

			err = p.CheckLogTx(transactionBinaryDataFull)
			if err != nil {
				//fmt.Println("err", err)
				//fmt.Println("RollbackTo")
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
				if txCounter[userId] > p.Variables.Int64["max_block_user_transactions"]  {
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
			fmt.Println("MethodName", MethodName+"Init")
			err_ := utils.CallMethod(p,MethodName+"Init")
			if _, ok := err_.(error); ok {
				fmt.Println(err)
				return utils.ErrInfo(err_.(error))
			}

			fmt.Println("MethodName", MethodName+"Front")
			err_ = utils.CallMethod(p,MethodName+"Front")
			if _, ok := err_.(error); ok {
				fmt.Println(err)
				p.RollbackTo(txForRollbackTo, true, false);
				return utils.ErrInfo(err_.(error))
			}

			fmt.Println("MethodName", MethodName)
			err_ = utils.CallMethod(p,MethodName)
			if _, ok := err_.(error); ok {
				fmt.Println(err)
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
	forSha := fmt.Sprintf("%d,%s,%s,%d,%d,%d", blockId, p.PrevBlock.Hash, p.MrklRoot, p.BlockData.Time, p.BlockData.UserId, p.BlockData.Level)
	fmt.Println("forSha", forSha)
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
	p.TxUserID = utils.BytesToInt64(TxMap["user_id"])
	p.TxTime = utils.BytesToInt64(TxMap["time"])
	p.PublicKeys =nil
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
			myUserId, err = p.Single("SELECT user_id FROM "+myPrefix+"my_table").Int64()
			if err != nil {
				return myUserId, myBlockId, myPrefix, myUserIds, err
			}
		}
	} else {
		myUserId, err = p.Single("SELECT user_id FROM my_table").Int64()
		if err != nil {
			return myUserId, myBlockId, myPrefix, myUserIds, err
		}
		myUserIds = append(myUserIds, myUserId)
	}
	myBlockId, err = p.Single("SELECT my_block_id FROM config").Int64()
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

func (p *Parser) CheckInputData(data map[string]string) (error) {
	for k, v := range data {
		if !utils.CheckInputData(p.TxMap[k], v) {
			return fmt.Errorf("incorrect "+k)
		}
	}
	return nil
}

func (p *Parser) limitRequestsRollback(txType string) error {
	time := p.TxMap["time"]
	return p.DCDB.ExecSql("DELETE FROM log_time_"+txType+" WHERE user_id = ? AND time = ?", p.TxMap["user_id"], time)
}

func (p *Parser) countMinerAttempt(userId, vType string) (int64, error) {
	count, err := p.DCDB.Single("SELECT count(user_id) FROM votes_miners WHERE user_id = ? AND type = ?", userId, vType).Int64()
	if err != nil {
		return 0, utils.ErrInfo(err)
	}
	return count, nil
}
// откатываем ID на кол-во затронутых строк
func (p *Parser) rollbackAI(table string, num int64) (error) {

	if num == 0 {
		return nil
	}

	current, err := p.Single("SELECT id FROM "+table+" ORDER BY id DESC LIMIT 1", ).Int64()
	if err != nil {
		return utils.ErrInfo(err)
	}
	NewAi := current + num

	if p.ConfigIni["db_type"] == "postgresql" {
		pg_get_serial_sequence, err := p.Single("SELECT pg_get_serial_sequence('"+table+"', 'id')").String()
		if err != nil {
			return utils.ErrInfo(err)
		}
		err = p.ExecSql("ALTER SEQUENCE "+pg_get_serial_sequence+" RESTART WITH "+utils.Int64ToStr(NewAi))
		if err != nil {
			return utils.ErrInfo(err)
		}
	} else if p.ConfigIni["db_type"] == "mysql" {
		err := p.DCDB.ExecSql("ALTER TABLE "+table+" AUTO_INCREMENT = ?", NewAi)
		if err != nil {
			return utils.ErrInfo(err)
		}
	} else if p.ConfigIni["db_type"] == "sqlite" {
		err := p.DCDB.ExecSql("UPDATE SQLITE_SEQUENCE SET seq = ? WHERE name = ?", NewAi, table)
		if err != nil {
			return utils.ErrInfo(err)
		}
	}
	return nil
}

func (p *Parser) generalCheckAdmin() error {
	if !utils.CheckInputData(p.TxMap["user_id"], "int") {
		return utils.ErrInfoFmt("user_id")
	}
	// точно ли это текущий админ
	err := p.getAdminUserId()
	if err != nil {
		return utils.ErrInfo(err)
	}
	// точно ли это текущий админ
	if p.AdminUserId != utils.BytesToInt64(p.TxMap["user_id"]) {
		return utils.ErrInfoFmt("user_id (%d!=%d)", p.AdminUserId, p.TxMap["user_id"])
	}
	// проверим, есть ли такой юзер и заодно получим public_key
	data, err := p.DCDB.OneRow("SELECT public_key_0, public_key_1, public_key_2 FROM  users WHERE user_id = ?", p.TxMap["user_id"]).String()
	if err != nil {
		return utils.ErrInfo(err)
	}
	if len(data["public_key_0"])==0 {
		return utils.ErrInfoFmt("incorrect user_id")
	}
	p.PublicKeys = append(p.PublicKeys, []byte(data["public_key_0"]))
	if len(data["public_key_1"]) > 0 {
		p.PublicKeys = append(p.PublicKeys, []byte(data["public_key_1"]))
	}
	if len(data["public_key_2"]) > 0 {
		p.PublicKeys = append(p.PublicKeys, []byte(data["public_key_2"]))
	}
	// чтобы не записали слишком длинную подпись
	// 128 - это нод-ключ
	if len(p.TxMap["sign"]) < 128 || len(p.TxMap["sign"]) > 5000  {
		return utils.ErrInfoFmt("incorrect sign size")
	}
	return nil
}

func (p *Parser) generalRollback(table string, whereUserId_ interface {}, addWhere string, AI bool) error {
	var whereUserId string
	switch whereUserId_.(type) {
	case string:
		whereUserId = whereUserId_.(string)
	case []byte:
		whereUserId = string(whereUserId_.([]byte))
	}

	where := ""
	if len(whereUserId) > 0 {
		where = fmt.Sprintf("WHERE user_id = %d", whereUserId)
	}
	// получим log_id, по которому можно найти данные, которые были до этого
	logId, err := p.Single("SELECT log_id FROM ? "+where+addWhere, table).Int64()
	if err != nil {
		return utils.ErrInfo(err)
	}
	// если $log_id = 0, значит восстанавливать нечего и нужно просто удалить запись
	if logId == 0 {
		err = p.ExecSql("DELETE FROM ? "+where+addWhere, table)
		if err != nil {
			return utils.ErrInfo(err)
		}
	} else {
		// данные, которые восстановим
		data, err := p.OneRow("SELECT * FROM log_"+table+" WHERE log_id = ?", logId).String()
		if err != nil {
			return utils.ErrInfo(err)
		}
		addSql := ""
		for k, v := range data {
			// block_id т.к. в log_ он нужен для удаления старых данных, а в обычной табле не нужен
			if k == "log_id" || k == "prev_log_id"  || k == "block_id" {
				continue
			}
			if k == "node_public_key" {
				switch p.ConfigIni["db_type"] {
				case "sqlite":
					addSql += fmt.Sprintf("%v='%x',", k, v)
				case "postgresql":
					addSql += fmt.Sprintf("%v=decode(%x,'HEX'),", k, v)
				case "mysql":
					addSql += fmt.Sprintf("%v=UNHEX(%x),", k, v)
				}
			} else {
				addSql += fmt.Sprintf("%v = '%v',", k, v)
			}
		}
		// всегда пишем предыдущий log_id
		addSql += fmt.Sprintf("log_id = %d,", data["prev_log_id"])
		addSql = addSql[0:len(addSql)-1]
		err = p.ExecSql("UPDATE ? SET "+addSql+where+addWhere, table)
		if err != nil {
			return utils.ErrInfo(err)
		}
		// подчищаем _log
		err = p.ExecSql("DELETE FROM ? WHERE log_id= ?", table, logId)
		if err != nil {
			return utils.ErrInfo(err)
		}
		err = p.rollbackAI("log_"+table, 1)
		if err != nil {
			return utils.ErrInfo(err)
		}
	}
	return nil
}

func arrayIntersect(arr1, arr2 map[int]int) bool {
	for _, v := range arr1 {
		for _, v2 := range arr2 {
			if v == v2 {
				return true
			}
		}
	}
	return false
}

func  (p *Parser) minersCheckMyMinerIdAndVotes0(data *MinerData) bool {
	if (arrayIntersect(data.myMinersIds, data.minersIds)) && (int64(data.votes0) > data.minMinersKeepers || data.votes0 == len(data.minersIds)) {
		return true
	} else {
		return false
	}
}

func  (p *Parser) minersCheckVotes1(data *MinerData) bool {
	fmt.Println("data.votes1",data.votes1)
	fmt.Println("data.minMinersKeepers",data.minMinersKeepers)
	fmt.Println("data.minersIds",data.minersIds)
	if int64(data.votes1) >= data.minMinersKeepers || data.votes1 == len(data.minersIds) /*|| data.adminUiserId == p.TxUserID Админская нода не решающая*/ {
		return true
	} else {
		return false
	}
}

func getMinersKeepers(ctx0, maxMinerId0, minersKeepers0 string, arr0 bool) map[int]int {
	ctx:=utils.StrToInt(ctx0)
	maxMinerId:=utils.StrToInt(maxMinerId0)
	minersKeepers:=utils.StrToInt(minersKeepers0)
	result := make(map[int]int)
	newResult := make(map[int]int)
	var ctx_ float64
	ctx_ = float64(ctx)
	for i:=0; i<minersKeepers; i++ {
		//fmt.Println("ctx", ctx)
		//var hi float34
		hi := ctx_ / float64(127773)
		//fmt.Println("hi", hi)
		lo := int(ctx_) % 127773
		//fmt.Println("lo", lo)
		x := (float64(16807) * float64(lo)) - (float64(2836) * hi)
		//fmt.Println("x", x, float64(16807), float64(lo), float64(2836), hi)
		if x <= 0 {
			x += 0x7fffffff
		}
		ctx_ = x
		rez := int(ctx_) % (maxMinerId+1)
		//fmt.Println("rez", rez)
		if rez == 0 {
			rez = 1
		}
		result[rez] = 1
	}
	if arr0 {
		i:=0
		for k, _ := range result {
			newResult[i] = k
			i++
		}
	} else {
		newResult = result
	}
	return newResult
}

func (p *Parser) FormatBlockData() string {
	result := ""
	v := reflect.ValueOf(*p.BlockData)
	typeOfT := v.Type()
	if typeOfT.Kind() == reflect.Ptr {
		typeOfT = typeOfT.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		name := typeOfT.Field(i).Name
		switch name {
		case "BlockId", "Time", "UserId", "Level":
			result += "["+name+"] = "+fmt.Sprintf("%d\n", v.Field(i).Interface())
		case "Sign", "Hash", "HeadHash":
			result += "["+name+"] = "+fmt.Sprintf("%x\n", v.Field(i).Interface())
		default :
			result += "["+name+"] = "+fmt.Sprintf("%s\n", v.Field(i).Interface())
		}
	}
	return result
}

func (p *Parser) FormatTxMap() string {
	result := ""
	for k, v := range p.TxMap {
		switch k {
		case "sign":
			result += "["+k+"] = "+fmt.Sprintf("%x\n", v)
		default :
			result += "["+k+"] = "+fmt.Sprintf("%s\n", v)
		}
	}
	return result
}

func (p *Parser) ErrInfo(err_ interface{}) error {
	var err error
	switch err_.(type) {
	case error:
		err = err_.(error)
	case string:
		err = fmt.Errorf(err_.(string))
	}
	return fmt.Errorf("[ERROR] %s (%s)\n%s\n%s", err, utils.Caller(1), p.FormatBlockData(), p.FormatTxMap())
}

func (p *Parser) maxDayVotesRollback() (error) {
	err := p.ExecSql("DELETE FROM log_time_votes WHERE user_id = ? AND time = ?", p.TxMap["user_id"], p.TxTime)
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}

func (p *Parser) maxDayVotes() (error) {
	// нельзя за сутки голосовать более max_day_votes раз
	num, err := p.Single("SELECT count(time) FROM log_time_votes WHERE user_id = ?", p.TxMap["user_id"]).Int64()
	if err != nil {
		return p.ErrInfo(err)
	}
	if num >= p.Variables.Int64["max_day_votes"] {
		return p.ErrInfo(fmt.Sprintf("[limit_requests] max_day_votes log_time_votes limits %d >=%d", num, p.Variables.Int64["max_day_votes"]))
	} else {
		err = p.ExecSql("INSERT INTO log_time_votes ( user_id, time ) VALUES ( ?, ? )", p.TxMap["user_id"], p.TxTime)
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	return nil
}


// начисление баллов
func (p *Parser) points(points int64) (error) {
	data, err := p.OneRow("SELECT time_start, points, log_id FROM points WHERE user_id = ?", p.TxMap["user_id"]).String()
	if err != nil {
		return p.ErrInfo(err)
	}
	pointsStatusTimeStart, err := p.Single("SELECT time_start FROM points_status WHERE user_id = ? ORDER BY time_start DESC", p.TxMap["user_id"]).Int64()
	if err != nil {
		return p.ErrInfo(err)
	}

	timeStart := data["time_start"]
	prevLogId := data["log_id"]

	// если $time_start = 0, значит это первый голос юзера
	if len(timeStart) == 0 {
		err = p.ExecSql("INSERT INTO points ( user_id, time_start, points ) VALUES ( ?, ?, ? )", p.TxMap["user_id"], p.BlockData.Time, points)
		if err != nil {
			return p.ErrInfo(err)
		}
		// первый месяц в любом случае будет юзером
		err = p.ExecSql("INSERT INTO points_status ( user_id, time_start, status, block_id ) VALUES ( ?, ?, 'user', ? )", p.TxMap["user_id"], p.BlockData.Time, p.BlockData.BlockId)
		if err != nil {
			return p.ErrInfo(err)
		}
	} else if p.BlockData.Time - pointsStatusTimeStart > p.Variables.Int64["points_update_time"] { // если прошел месяц
		err = p.pointsUpdate(utils.StrToInt64(data["points"]), prevLogId, timeStart, pointsStatusTimeStart, p.TxUserID, points)
		if err != nil {
			return p.ErrInfo(err)
		}
	} else { // прошло меньше месяца
		// прибавляем баллы
		err = p.ExecSql("UPDATE points SET points = points+1 WHERE user_id = ?", points, p.TxMap["user_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
		/*// просто для вывода в лог
		err = p.ExecSql("SELECT * FROM points WHERE user_id = ?", p.TxMap["user_id"])
		if err != nil {
			return p.ErrInfo(err)
		}*/
	}
	return nil
}



// добавляем новые points_status
// $points - текущие points юзера из таблы points
// $new_points - новые баллы, если это вызов из тр-ии, где идет головование
func (p *Parser) pointsUpdate(points int64, prevLogId string, timeStart string, pointsStatusTimeStart int64, userId int64, newPoints int64) (error) {

	// среднее значение баллов
	mean, err := p.Single("SELECT sum(points)/count(points) FROM points WHERE points > 0").Float64()
	if err != nil {
		return p.ErrInfo(err)
	}

	// есть ли тр-ия с голосованием votes_complex за последние 4 недели
	count, err := p.Single("SELECT count(user_id) FROM votes_miner_pct WHERE user_id = ? AND time > ?", userId, (p.BlockData.Time - p.Variables.Int64["limit_votes_complex_period"]*2)).Int64()
	if err != nil {
		return p.ErrInfo(err)
	}

	// и хватает ли наших баллов для получения статуса майнера
	if count > 0 && float64(points+newPoints) >= mean * float64(p.Variables.Int64["points_factor"]) {
		// от $time_start до текущего времени могло пройти несколько месяцев. 1-й месяц будет майнер, остальные - юзер
		minerStartTime := pointsStatusTimeStart + p.Variables.Int64["points_update_time"]
		err = p.ExecSql("INSERT INTO points_status ( user_id, time_start, status, block_id ) VALUES ( ?, ?, 'miner', ? )", userId, minerStartTime, p.BlockData.BlockId)
		if err != nil {
			return p.ErrInfo(err)
		}

		// сколько прошло месяцев после $miner_start_time
		remainingTime := p.BlockData.Time - minerStartTime
		if remainingTime > 0 {
			remainingMonths := math.Floor(float64(remainingTime / p.Variables.Int64["points_update_time"]))
			if remainingMonths > 0 {
				// следующая запись должна быть ровно через 1 месяц после $miner_start_time
				userStartTime := minerStartTime + p.Variables.Int64["points_update_time"]
				err = p.ExecSql("INSERT INTO points_status ( user_id, time_start, status, block_id ) VALUES ( ?, ?, 'user', ? )", userId, userStartTime, p.BlockData.BlockId)
				if err != nil {
					return p.ErrInfo(err)
				}
				// и если что-то осталось
				if remainingMonths > 1 {
					userStartTime = minerStartTime + int64(remainingMonths) * p.Variables.Int64["points_update_time"]
					err = p.ExecSql("INSERT INTO points_status ( user_id, time_start, status, block_id ) VALUES ( ?, ?, 'user', ? )", userId, userStartTime, p.BlockData.BlockId)
					if err != nil {
						return p.ErrInfo(err)
					}
				}
			}
		}
	} else {
		// следующая запись должна быть ровно через 1 месяц после предыдущего статуса
		userStartTime := pointsStatusTimeStart + p.Variables.Int64["points_update_time"]
		err = p.ExecSql("INSERT INTO points_status ( user_id, time_start, status, block_id ) VALUES ( ?, ?, 'user', ? )", userId, userStartTime, p.BlockData.BlockId)
		if err != nil {
			return p.ErrInfo(err)
		}
		// сколько прошло месяцев после $miner_start_time
		remainingTime :=  p.BlockData.Time - userStartTime

		if remainingTime > 0 {

			remainingMonths := math.Floor(float64(remainingTime / p.Variables.Int64["points_update_time"]))
			if remainingMonths > 0 {
				userStartTime = userStartTime + int64(remainingMonths) * p.Variables.Int64["points_update_time"]
				err = p.ExecSql("INSERT INTO points_status ( user_id, time_start, status, block_id ) VALUES ( ?, ?, 'user', ? )", userId, userStartTime, p.BlockData.BlockId)
				if err != nil {
					return p.ErrInfo(err)
				}
			}
		}
	}

	// перед тем, как обновить time_start, нужно его залогировать
	logId, err := p.ExecSqlGetLastInsertId("INSERT INTO log_points ( time_start, points, block_id, prev_log_id ) VALUES ( ?, ?, ?, ? )", timeStart, points, p.BlockData.BlockId, prevLogId)
	if err != nil {
		return p.ErrInfo(err)
	}

	// начисляем баллы с чистого листа и обновляем время
	err = p.ExecSql("UPDATE points SET points = 0, time_start = ?, log_id = ? WHERE user_id = ?", p.BlockData.Time, logId, userId)
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}

func (p *Parser) checkTrueVotes(data map[string]int64) bool {
	if 	data["votes_1"] >= data["votes_1_min"] ||
			(p.TxUserID == p.AdminUserId && string(p.TxMap["result"]) == "1" && data["count_miners"] < 1000)	|| data["votes_1"] == data["count_miners"] {
		return true
	} else {
		return false
	}
}

func (p *Parser) insOrUpdMiners(userId int64) (int64, error) {
	miners, err := p.OneRow("SELECT miner_id, log_id FROM miners WHERE active = 0").Int64()
	if err != nil {
		return 0, p.ErrInfo(err)
	}
	minerId := miners["miner_id"]
	if minerId == 0 {
		minerId, err = p.ExecSqlGetLastInsertId("INSERT INTO miners (active) VALUES (1)")
		if err != nil {
			return 0, p.ErrInfo(err)
		}
	} else {
		logId, err := p.ExecSqlGetLastInsertId("INSERT INTO log_miners ( block_id, prev_log_id ) VALUES ( ?, ?)", p.BlockData.BlockId, miners["log_id"])
		if err != nil {
			return 0, p.ErrInfo(err)
		}
		err = p.ExecSql("UPDATE miners SET active = 1, log_id = ? WHERE miner_id = ?", logId, minerId)
		if err != nil {
			return 0, p.ErrInfo(err)
		}
	}
	err = p.ExecSql("UPDATE miners_data SET status = 'miner', miner_id = ?, reg_time = ? WHERE user_id = ?", minerId, p.BlockData.Time, userId)
	if err != nil {
		return 0, p.ErrInfo(err)
	}
	return minerId, nil
}

func (p *Parser) check24hOrAdminVote(data map[string]int64) bool {

	if (/*прошло > 24h от начала голосования ?*/(p.BlockData.Time - data["votes_period"]) > data["votes_start_time"] &&
			// преодолен ли один из лимитов, либо проголосовали все майнеры
			(data["votes_0"] >= data["votes_0_min"] ||
				data["votes_1"] >= data["votes_1_min"] ||
				data["votes_0"] == data["count_miners"] ||
				data["votes_1"] == data["count_miners"])) ||
			/*голос админа решающий в любое время, если <1000 майнеров в системе*/
			(p.TxUserID == p.AdminUserId && data["count_miners"] < 1000) {
		return true
	} else {
		return false
	}
}


func (p *Parser) insOrUpdMinersRollback(minerId int64) error {

	// нужно проверить, был ли получен наш miner_id в результате замены забаненного майнера
	logId, err := p.Single("SELECT log_id FROM miners WHERE miner_id  =  ?", minerId).Int64()
	if err != nil {
		return p.ErrInfo(err)
	}
	if logId > 0 {

		// данные, которые восстановим
		prevLogId, err := p.Single("SELECT prev_log_id FROM log_miners WHERE log_id  =  ?", logId).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		// $log_data['prev_log_id'] может быть = 0
		err = p.ExecSql("UPDATE miners SET active = 0, log_id = ? WHERE miner_id = ?", prevLogId, minerId)
		if err != nil {
			return p.ErrInfo(err)
		}
		// подчищаем _log
		err = p.ExecSql("DELETE FROM log_miners WHERE log_id = ?", logId)
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.rollbackAI("log_miners", 1)
		if err != nil {
			return p.ErrInfo(err)
		}
	} else {
		err = p.ExecSql("DELETE FROM miners WHERE miner_id = ?", minerId)
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.rollbackAI("miners", 1)
		if err != nil {
			return p.ErrInfo(err)
		}
	}

	return nil
}


// $points - баллы, которые были начислены за голос
func (p *Parser) pointsRollback(points int64) error {
	data, err := p.OneRow("SELECT time_start, points, log_id FROM points WHERE user_id  =  ?", p.TxMap["user_id"]).Int64()
	if err != nil {
		return p.ErrInfo(err)
	}
	if len(data) == 0 {
		return nil
	}
	// если time_start=времени в блоке, points=$points и log_id=0, значит это самая первая запись
	if data["time_start"] == p.BlockData.Time && data["points"] == points && data["log_id"] == 0 {
		err = p.ExecSql("DELETE FROM points WHERE user_id = ?", p.TxMap["user_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.ExecSql("DELETE FROM points_status WHERE user_id = ?", p.TxMap["user_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
	} else if data["time_start"] == p.BlockData.Time { // если прошел месяц и запись в табле points была обновлена в этой тр-ии, т.е. time_start = block_data['time']
		err = p.pointsUpdateRollback(data["log_id"], p.TxUserID)
		if err != nil {
			return p.ErrInfo(err)
		}
	} else { // прошло меньше месяца
		// отнимаем баллы
		err = p.ExecSql("UPDATE points SET points = points - "+utils.Int64ToStr(points)+" WHERE user_id = ?", points, p.TxMap["user_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	return nil
}

func (p *Parser) pointsUpdateRollback(logId, userId int64) error {
	err := p.ExecSql("DELETE FROM points_status WHERE block_id = ?", p.BlockData.BlockId)
	if err != nil {
		return p.ErrInfo(err)
	}
	if logId > 0 {
		// данные, которые восстановим
		logData, err := p.OneRow("SELECT time_start, prev_log_id, points FROM log_points WHERE log_id  =  ?", logId).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.ExecSql("UPDATE points SET time_start = ?, points = ?, log_id = ? WHERE user_id = ?", logData["time_start"], logData["points"], logData["prev_log_id"], userId)
		if err != nil {
			return p.ErrInfo(err)
		}
		// подчищаем _log
		err = p.ExecSql("DELETE FROM log_points WHERE log_id = ?", logId)
		if err != nil {
			return p.ErrInfo(err)
		}
		p.rollbackAI("log_points", 1)
	} else {
		err = p.ExecSql("DELETE FROM points WHERE user_id = ?", userId)
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	return nil
}
