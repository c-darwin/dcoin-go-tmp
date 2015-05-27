package dcparser
import (
	//"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"utils"
	//"os"
	"time"
	"consts"
	"sync"
	"reflect"
	"math"
	"log"
	"encoding/json"
	"database/sql"
	"strings"
)

type txMapsType struct {
	Int64 map[string]int64
	String map[string]string
	Bytes map[string][]byte
	Float64 map[string]float64
	Money map[string]float64
}
type Parser struct {
	*utils.DCDB
	TxMaps *txMapsType
	TxMap map[string][]byte
	TxMapS map[string]string
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
	newPublicKeysHex [3][]byte
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


func (p *Parser) checkMiner(userId int64) error {
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


func (p *Parser) GetTxMaps(fields []map[string]string) (error) {
	if len(p.TxSlice) != len(fields)+4 {
		return fmt.Errorf("bad transaction_array %d != %d (type=%d)",  len(p.TxSlice),  len(fields)+4, p.TxSlice[0])
	}
	//log.Println("p.TxSlice", p.TxSlice)
	p.TxMap = make(map[string][]byte)
	p.TxMaps = new(txMapsType)
	p.TxMaps.Float64 = make(map[string]float64)
	p.TxMaps.Money = make(map[string]float64)
	p.TxMaps.Int64 = make(map[string]int64)
	p.TxMaps.Bytes = make(map[string][]byte)
	p.TxMaps.String = make(map[string]string)
	p.TxMaps.Bytes["hash"] = p.TxSlice[0]
	p.TxMaps.Int64["type"] = utils.BytesToInt64( p.TxSlice[1])
	p.TxMaps.Int64["time"] = utils.BytesToInt64(p.TxSlice[2])
	p.TxMaps.Int64["user_id"] = utils.BytesToInt64(p.TxSlice[3])
	p.TxMap["hash"] = p.TxSlice[0]
	p.TxMap["type"] = p.TxSlice[1]
	p.TxMap["time"] = p.TxSlice[2]
	p.TxMap["user_id"] = p.TxSlice[3]
	for i:=0; i<len(fields); i++ {
		for field, fType := range fields[i] {
			p.TxMap[field] = p.TxSlice[i+4]
			switch fType {
			case "int64":
				p.TxMaps.Int64[field] =  utils.BytesToInt64(p.TxSlice[i+4])
			case "float64":
				p.TxMaps.Float64[field] =  utils.BytesToFloat64(p.TxSlice[i+4])
			case "money":
				p.TxMaps.Money[field] = math.Floor(utils.BytesToFloat64(p.TxSlice[i+4])*100)/100
			case "bytes":
				p.TxMaps.Bytes[field] =  p.TxSlice[i+4]
			case "string":
				p.TxMaps.String[field] =  string(p.TxSlice[i+4])
			}
		}
	}
	p.TxUserID = p.TxMaps.Int64["user_id"]
	p.TxTime = p.TxMaps.Int64["time"]
	p.PublicKeys = nil
	//log.Println("p.TxMaps", p.TxMaps)
	//log.Println("p.TxMap", p.TxMap)
	return nil
}


// старое
func (p *Parser) GetTxMap(fields []string) (map[string][]byte, error) {
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
	fmt.Println("TxMap", TxMap)
	//fmt.Println("TxMap[hash]", TxMap["hash"])
	//fmt.Println("p.TxSlice[0]", p.TxSlice[0])
	return TxMap, nil
}

// старое
func (p *Parser) GetTxMapStr(fields []string) (map[string]string, error) {
	//fmt.Println("p.TxSlice", p.TxSlice)
	//fmt.Println("fields", fields)
	if len(p.TxSlice) != len(fields)+4 {
		return nil, fmt.Errorf("bad transaction_array %d != %d (type=%d)",  len(p.TxSlice),  len(fields)+4, p.TxSlice[0])
	}
	TxMapS := make(map[string]string)
	TxMapS["hash"] = string(p.TxSlice[0])
	TxMapS["type"] =string( p.TxSlice[1])
	TxMapS["time"] = string(p.TxSlice[2])
	TxMapS["user_id"] = string(p.TxSlice[3])
	for i, field := range fields {
		TxMapS[field] = string(p.TxSlice[i+4])
	}
	p.TxUserID = utils.StrToInt64(TxMapS["user_id"])
	p.TxTime = utils.StrToInt64(TxMapS["time"])
	p.PublicKeys =nil
	fmt.Println("TxMapS", TxMapS)
	//fmt.Println("TxMap[hash]", TxMap["hash"])
	//fmt.Println("p.TxSlice[0]", p.TxSlice[0])
	return TxMapS, nil
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
	case int64:
		whereUserId = utils.Int64ToStr(whereUserId_.(int64))
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
	data, err := p.OneRow("SELECT time_start, points, log_id FROM points WHERE user_id = ?", p.TxMap["user_id"]).Int64()
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
	if timeStart == 0 {
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
		err = p.pointsUpdate(data["points"], prevLogId, timeStart, pointsStatusTimeStart, p.TxUserID, points)
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


func (p *Parser) calcProfit_(amount float64, timeStart, timeFinish int64, pctArray []map[int64]map[string]float64, pointsStatusArray []map[int64]string, holidaysArray [][]int64, maxPromisedAmountArray []map[int64]string, currencyId int64, repaidAmount float64) (float64, error) {
	if p.BlockData!=nil && p.BlockData.BlockId<=24946 {
		return p.CalcProfit_24946(amount, timeStart, timeFinish, pctArray, pointsStatusArray, holidaysArray, maxPromisedAmountArray, currencyId, repaidAmount)
	} else {
		return p.CalcProfit(amount, timeStart, timeFinish, pctArray, pointsStatusArray, holidaysArray, maxPromisedAmountArray, currencyId, repaidAmount)
	}
}

// обновление points_status на основе points
// вызов данного метода безопасен для rollback методов, т.к. при rollback данные кошельков восстаналиваются из log_wallets не трогая points
func (p *Parser) pointsUpdateMain(userId int64) error {
	data, err := p.OneRow("SELECT time_start, points, log_id FROM points WHERE user_id  =  ?", userId).Int64()
	if err != nil {
		return p.ErrInfo(err)
	}
	pointsStatusTimeStart, err := p.Single("SELECT time_start FROM points_status WHERE user_id  =  ? ORDER BY time_start DESC", userId).Int64()
	if err != nil {
		return p.ErrInfo(err)
	}
	if len(data) > 0 && p.BlockData.Time - pointsStatusTimeStart > p.Variables.Int64["points_update_time"] {
		err = p.pointsUpdate(data["points"], data["log_id"], data["time_start"], pointsStatusTimeStart, userId, 0)
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	return nil
}

func (p *Parser) getPointsStatus(userId int64) ([]map[int64]string, error) {

	// т.к. перед вызовом этой функции всегда идет обновление points_status, значит при данном запросе у нас
	// всегда будут свежие данные, т.е. крайний элемент массива всегда будет относиться к текущим 30-и дням
	var result []map[int64]string
	rows, err := p.Query("SELECT time_start, status FROM points_status WHERE user_id= ? ORDER BY time_start ASC", userId)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var time_start int64
		var status string
		err = rows.Scan(&time_start, &status)
		if err!= nil {
			return result, err
		}
		result = append(result, map[int64]string{time_start:status})
	}

	// НО! При фронтальной проверке может получиться, что последний элемент miner и прошло более 30-и дней.
	// поэтому нужно добавлять последний элемент = user, если вызов происходит не в блоке
	pointsUpdateTime := p.Variables.Int64["points_update_time"]
	if p.BlockData!=nil && len(result)>0 {
		for time_start, _ := range result[len(result)-1] {
			if time_start < time.Now().Unix() - pointsUpdateTime {
				result = append(result, map[int64]string{time_start+pointsUpdateTime:"user"})
			}
		}
	}
	// для майнеров, которые не получили ни одного балла, а уже шлют кому-то DC, или для всех юзеров
	if len(result) == 0 {
		result = append(result, map[int64]string{0:"user"})
	}
	return result, nil
}

func (p *Parser) pointsUpdateRollbackMain(userId int64) error {
	data, err := p.OneRow("SELECT time_start, log_id FROM points WHERE user_id  =  ?", userId).Int64()
	if err != nil {
		return err
	}
	if p.BlockData.Time == data["time_start"] {
		err = p.pointsUpdateRollback(data["log_id"], userId)
		if err != nil {
			return err
		}
	}
	return nil
}

// добавляем новые points_status
// $points - текущие points юзера из таблы points
// $new_points - новые баллы, если это вызов из тр-ии, где идет головование
func (p *Parser) pointsUpdate(points, prevLogId, timeStart, pointsStatusTimeStart, userId, newPoints int64) (error) {

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

// не использовать для комментов
func (p *Parser) selectiveLoggingAndUpd(fields , values []string, table string, whereFields, whereValues []string) error {

	addSqlFields := ""
	for _, field := range fields {
		addSqlFields += field+",";
	}

	addSqlWhere := ""
	for i:=0; i < len(whereFields); i++ {
		addSqlWhere += whereFields[i]+"="+whereValues[i]+" AND "
	}
	if len(addSqlWhere) > 0 {
		addSqlWhere = " WHERE "+addSqlWhere[0:len(addSqlWhere)-5]
	}
	// если есть, что логировать
	logData, err := p.OneRow("SELECT "+addSqlFields+" log_id FROM "+table+" "+addSqlWhere).String()
	if err != nil {
		return err
	}
	if len(logData) > 0 {
		addSqlValues := ""
		for k, v := range logData {
			if utils.InSliceString(k, []string{"hash", "tx_hash", "public_key_0", "public_key_1", "public_key_2"}) && v!="" {
				query:=""
				switch p.ConfigIni["dbType"] {
				case "sqlite":
					query = `'`+v+`',`
				case "postgresql":
					query = `decode(`+v+`,'HEX'),`
				case "mysql":
					query = `UNHEX("`+v+`"),`
				}
				addSqlValues+=query
			} else {
				addSqlValues+=`'`+v+`',`
			}
		}
		addSqlValues = addSqlValues[0:len(addSqlValues)-1]

		logId, err := p.ExecSqlGetLastInsertId("INSERT INTO log_"+table+" ( "+addSqlFields+" prev_log_id, block_id ) VALUES ( "+addSqlValues+", ? )", p.BlockData.BlockId)
		if err != nil {
			return err
		}
		addSqlUpdate := ""
		for i:=0; i < len(fields); i++ {
			if utils.InSliceString(fields[i], []string{"hash", "tx_hash", "public_key_0", "public_key_1", "public_key_2"}) && len(values[i])!=0 {
				query:=""
				switch p.ConfigIni["dbType"] {
				case "sqlite":
					query = fields[i]+`='`+values[i]+`',`
				case "postgresql":
					query = fields[i]+`=decode(`+values[i]+`,'HEX'),`
				case "mysql":
					query = fields[i]+`=UNHEX("`+values[i]+`"),`
				}
				addSqlUpdate+=query
			} else {
				addSqlUpdate+= fields[i]+`='`+values[i]+`',`
			}
		}
		err = p.ExecSql("UPDATE "+table+" SET "+addSqlUpdate+" log_id = ? "+addSqlWhere, logId)
		if err != nil {
			return err
		}
	} else {
		addSqlIns0 := "";
		addSqlIns1 := "";
		for i:=0; i < len(fields); i++ {
			addSqlIns0 += ``+fields[i]+`,`
			if utils.InSliceString(fields[i], []string{"hash", "tx_hash", "public_key_0", "public_key_1", "public_key_2"}) && len(values[i])!=0 {
				query:=""
				switch p.ConfigIni["dbType"] {
				case "sqlite":
					query = `'`+values[i]+`',`
				case "postgresql":
					query = `decode(`+values[i]+`,'HEX'),`
				case "mysql":
					query = `UNHEX("`+values[i]+`"),`
				}
				addSqlIns1+=query
			} else {
				addSqlIns1+=`'`+values[i]+`',`
			}
		}
		for i:=0; i< len(whereFields); i++ {
			addSqlIns0+=``+whereFields[i]+`,`
			addSqlIns1+=`'`+whereValues[i]+`',`
		}
		addSqlIns0 = addSqlIns0[0:len(addSqlIns0)-1]
		addSqlIns1 = addSqlIns1[0:len(addSqlIns1)-1]
		err = p.ExecSql("INSERT INTO "+table+" ("+addSqlIns0+") VALUES ("+addSqlIns1+")")
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) loan_payments(toUserId int64, amount float64, currencyId int64) (float64, error) {
	amountForCredit := amount;

	// нужно узнать, какую часть от суммы заещик хочет оставить себе
	creditPart, err := p.Single("SELECT credit_part FROM users WHERE user_id  =  ?", toUserId).Float64()
	if err != nil {
		return 0, err
	}
	if creditPart > 0 {
		save := math.Floor( utils.Round( amount*(creditPart/100), 3) *100 ) / 100
		if save < 0.01 {
			save = 0
		}
		amountForCredit-=save
	}
	amountForCreditSave := amountForCredit;

	rows, err := p.Query("SELECT pct, amount, id, to_user_id FROM credits WHERE from_user_id = ? AND currency_id = ? AND amount > 0 AND del_block_id = 0 ORDER BY time", toUserId, currencyId)
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var pct, amount float64
		var id, to_user_id int64
		err = rows.Scan(&pct, &amount, &id, &to_user_id)
		if err!= nil {
			return 0, err
		}
		var sum float64
		var take float64
		if p.BlockData.BlockId > 169525 {
			sum = utils.Round(pct/100 * amountForCreditSave, 2);
		} else {
			sum = utils.Round(pct/100 * amount, 2);
		}
		if sum < 0.01 {
			sum = 0.01
		}
		if (sum > amountForCredit) {
			sum = amountForCredit;
		}
		if (sum - amount > 0) {
			take = amount;
		} else {
			take = sum
		}
		amountForCredit -= take;
		err := p.selectiveLoggingAndUpd([]string{"amount", "tx_hash", "tx_block_id"}, []string{utils.Float64ToStr(amount-take), string(p.TxHash), utils.Int64ToStr(p.BlockData.BlockId)}, "credits", []string{"id"}, []string{utils.Int64ToStr(id)})
		if err!= nil {
			return 0, err
		}
		err = p.updateRecipientWallet(toUserId, currencyId, take, "loan_payment", toUserId, "loan_payment", "decrypted", false)
		if err!= nil {
			return 0, err
		}
	}
	return amount - (amountForCreditSave - amountForCredit), nil
}

func (p *Parser) adaptQuery(query_ string) (string) {
	query:=query_
	if p.ConfigIni["db_type"]=="mysql" {
		query = strings.Replace(query, "delete", "`delete`", -1)
	}
	return query
}

/*
 * Начисляем новые DC юзеру, пересчитав ему % от того, что уже было на кошельке
 * */
func (p *Parser) updateRecipientWallet(toUserId, currencyId int64, amount float64, from string, fromId int64, comment, commentStatus string, credits bool) error {

	if currencyId == 0 {
		return p.ErrInfo("currencyId == 0")
	}
	walletWhere := "user_id = "+utils.Int64ToStr(toUserId)+" AND currency_id = "+utils.Int64ToStr(currencyId);
	walletData, err := p.OneRow("SELECT amount, amount_backup, last_update, log_id FROM wallets WHERE "+walletWhere).String()
	if err != nil {
		return p.ErrInfo(err)
	}
	// если кошелек получателя создан, то
	// начисляем DC на кошелек получателя.
	if len(walletData) > 0 {

		// возможно у юзера есть долги и по ним нужно рассчитаться.
		if credits != false && currencyId < 1000 {
			amount, err = p.loan_payments(toUserId, amount, currencyId);
			if err != nil {
				return p.ErrInfo(err)
			}
		}

		// нужно залогировать текущие значения для to_user_id
		logId, err := p.ExecSqlGetLastInsertId("INSERT INTO log_wallets ( amount, amount_backup, last_update, block_id, prev_log_id ) VALUES ( ?, ?, ?, ?, ? )", walletData["amount"], walletData["amount_backup"], walletData["last_update"], p.BlockData.BlockId, walletData["log_id"])
		if err != nil {
			return p.ErrInfo(err)
		}
		pointsStatus := []map[int64]string {{0:"user"}}
		// holidays не нужны, т.к. это не TDC, а DC
		// то, что выросло на кошельке
		var newDCSum float64
		if (currencyId>=1000)  {// >=1000 - это CF-валюты, которые не растут
			newDCSum = utils.StrToFloat64(walletData["amount"])
		} else {
			pct, err := p.GetPct()
			if err != nil {
				return p.ErrInfo(err)
			}
			profit, err := p.calcProfit_(utils.StrToFloat64(walletData["amount"]), utils.StrToInt64(walletData["lastUpdate"]), p.BlockData.Time, pct[currencyId], pointsStatus, [][]int64{}, []map[int64]string{}, 0, 0)
			newDCSum = utils.StrToFloat64(walletData["amount"])+profit
		}
		// итоговая сумма DC
		newDCSumEnd := newDCSum + amount;

		// Плюсуем на кошелек с соответствующей валютой.
		err = p.ExecSql("UPDATE wallets SET amount = ?, last_update = ?, log_id = ? WHERE "+walletWhere, newDCSumEnd, p.BlockData.Time, logId)
		if err != nil {
			return p.ErrInfo(err)
		}
	} else {

		// возможно у юзера есть долги и по ним нужно рассчитаться.
		if credits != false && currencyId < 1000 {
			amount, err = p.loan_payments(toUserId, amount, currencyId);
			if err != nil {
				return p.ErrInfo(err)
			}
		}

		// если кошелек получателя не создан, то создадим и запишем на него сумму перевода.
		err = p.ExecSql("INSERT INTO wallets ( user_id, currency_id, amount, last_update ) VALUES ( ?, ?, ?, ? )", toUserId, currencyId, amount, p.BlockData.Time)
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	myUserId, myBlockId, myPrefix, _ , err:= p.GetMyUserId(utils.BytesToInt64(p.TxMap["user_id"]))
	if err != nil {
		return p.ErrInfo(err)
	}
	if toUserId == myUserId && myBlockId <= p.BlockData.BlockId {
		if from == "from_user" && len(comment)>0 && commentStatus!="decrypted" { // Перевод между юзерами
			commentStatus = "encrypted"
			comment = string(utils.BinToHex([]byte(comment)))
		} else { // системные комменты (комиссия, майнинг и пр.)
			commentStatus = "decrypted"
		}
		// для отчетов и api пишем транзакцию
		err = p.ExecSql("INSERT INTO "+myPrefix+"my_dc_transactions ( type, type_id, to_user_id, amount, time, block_id, currency_id, comment, comment_status ) VALUES ( ?, ?, ?, ?, ?, ?, ?, ?, ? )", from, fromId, toUserId, amount, p.BlockData.Time, p.BlockData.BlockId, currencyId, comment, commentStatus)
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	return nil
}


func (p *Parser) updateSenderWallet(fromUserId, currencyId int64, amount,  commission float64, from string, fromId, toUserId int64, comment, commentStatus string) error {
	// получим инфу о текущих значениях таблицы wallets для юзера from_user_id
	walletWhere := "user_id = "+utils.Int64ToStr(fromUserId)+" AND currency_id = "+utils.Int64ToStr(currencyId)
	walletData, err := p.OneRow("SELECT amount, amount_backup, last_update, log_id FROM wallets WHERE "+walletWhere).String()
	if err != nil {
		return p.ErrInfo(err)
	}
	// перед тем, как менять значения на кошельках юзеров, нужно залогировать текущие значения для юзера from_user_id

	logId, err := p.ExecSqlGetLastInsertId("INSERT INTO log_wallets ( amount, amount_backup, last_update, block_id, prev_log_id ) VALUES ( ?, ?, ?, ?, ? )", walletData["amount"], walletData["amount_backup"], walletData["last_update"], p.BlockData.BlockId, walletData["log_id"])
	if err != nil {
		return p.ErrInfo(err)
	}

	pointsStatus := []map[int64]string {{0:"user"}}

	pct, err := p.GetPct()
	// пересчитаем DC на кошельке отправителя
	// обновим сумму и дату на кошельке отправителя.
	// holidays не нужны, т.к. это не TDC, а DC.
	var newDCSum float64
	walletDataAmountFloat64 := utils.StrToFloat64(walletData["amount"])
	if currencyId >= 1000 {// >=1000 - это CF-валюты, которые не растут
		newDCSum = walletDataAmountFloat64
	} else {
		profit, err := p.calcProfit_(walletDataAmountFloat64, utils.StrToInt64(walletData["lastUpdate"]), p.BlockData.Time, pct[currencyId], pointsStatus,[][]int64{}, []map[int64]string{}, 0, 0)
		if err != nil {
			return p.ErrInfo(err)
		}
		newDCSum = walletDataAmountFloat64 + profit - amount - commission;
	}
	err = p.ExecSql("UPDATE wallets SET amount = ?, last_update = ?, log_id = ? WHERE "+walletWhere, newDCSum, p.BlockData.Time, logId)
	if err != nil {
		return p.ErrInfo(err)
	}
	myUserId, myBlockId, myPrefix, _ , err:= p.GetMyUserId(utils.BytesToInt64(p.TxMap["user_id"]))
	if err != nil {
		return p.ErrInfo(err)
	}

	if fromUserId == myUserId && myBlockId <= p.BlockData.BlockId {
		var where0, set0 string
		if (from == "cfProject") {
			where0 = "";
			set0 = " toUserId = "+utils.Int64ToStr(toUserId)+", ";
		} else {
			where0 = " toUserId = "+utils.Int64ToStr(toUserId)+" AND ";
			set0 = "";
		}
		myId, err := p.Single("SELECT id FROM "+myPrefix+"my_dc_transactions WHERE status  =  'pending' AND type  =  '"+from+"' AND type_id  =  "+utils.Int64ToStr(fromUserId)+" AND "+where0+" amount  =  "+utils.Float64ToStr(amount)+" AND commission  =  "+utils.Float64ToStr(commission)+" AND currency_id  =  "+utils.Int64ToStr(currencyId)).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		if myId > 0 {
			err = p.ExecSql("UPDATE "+myPrefix+"my_dc_transactions SET status = 'approved', "+set0+" time = "+utils.Int64ToStr(p.BlockData.Time)+", block_id = "+utils.Int64ToStr(p.BlockData.BlockId)+" WHERE id = "+utils.Int64ToStr(myId))
			if err != nil {
				return p.ErrInfo(err)
			}
		} else {
			err = p.ExecSql("INSERT INTO "+myPrefix+"my_dc_transactions ( status, type, type_id, to_user_id, amount, commission, currency_id, comment, comment_status, time, block_id ) VALUES ( 'approved', ?, ?, ?, ?, ?, ?, ?, ?, ?, ? )", from, fromUserId, toUserId, amount, commission, currencyId, comment, commentStatus, p.BlockData.Time, p.BlockData.BlockId)
			if err != nil {
				return p.ErrInfo(err)
			}
		}
	}
	return nil
}

func (p*Parser) mydctxRollback () error {

	// если работаем в режиме пула
	community, err := p.GetCommunityUsers()
	if err != nil {
		return p.ErrInfo(err)
	}
	if len(community) > 0 {
		for i:=0; i<len(community); i++ {
			myPrefix := utils.Int64ToStr(community[i])+"_";
			// может захватиться несколько транзакций, но это не страшно, т.к. всё равно надо откатывать
			affect, err := p.ExecSqlGetAffect("DELETE FROM "+myPrefix+"my_dc_transactions WHERE block_id = ?", p.BlockData.BlockId)
			if err != nil {
				return p.ErrInfo(err)
			}
			err = p.rollbackAI(myPrefix+"my_dc_transactions", affect)
			if err != nil {
				return p.ErrInfo(err)
			}
		}
	} else {
		// может захватиться несколько транзакций, но это не страшно, т.к. всё равно надо откатывать
		affect, err := p.ExecSqlGetAffect("DELETE FROM my_dc_transactions WHERE block_id = ?", p.BlockData.BlockId)
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.rollbackAI("my_dc_transactions", affect)
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	return nil
}

func (p*Parser) limitRequestsMoneyOrdersRollback() error {
	err := p.ExecSql("DELETE FROM log_time_money_orders WHERE tx_hash = [hex]", p.TxHash)
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}

func (p*Parser) loanPaymentsRollback (userId, currencyId int64) error {
	// было `amount` > 0  в WHERE, из-за чего были проблемы с откатами, т.к. amount может быть равно 0, если кредит был погашен этой тр-ей
	rows, err := p.Query("SELECT id, to_user_id FROM credits WHERE from_user_id = ? AND currency_id = ? AND tx_block_id = ? AND tx_hash = [hex] AND del_block_id = 0 ORDER BY time DESC", userId, currencyId, p.BlockData.BlockId, p.TxHash)
	if err != nil {
		return  err
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		var to_user_id int64
		err = rows.Scan(&id, &to_user_id)
		if err!= nil {
			return err
		}
		err := p.selectiveRollback([]string{"amount", "tx_hash", "tx_block_id"}, "credits", "id="+id, false)
		if err != nil {
			return err
		}
		err = p.generalRollback("wallets", to_user_id, "AND currency_id = "+utils.Int64ToStr(currencyId), false)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p*Parser) getRefs(userId int64) ([3]int64, error) {
	var result [3]int64
	// получим рефов
	rez, err := p.Single("SELECT referral FROM users WHERE user_id  =  ?", userId).Int64()
	result[0] = rez
	if err != nil {
		return result, p.ErrInfo(err)
	}
	if result[0] > 0 {
		rez, err := p.Single("SELECT referral FROM users WHERE user_id  =  ?", result[0]).Int64()
		result[1] = rez
		if err != nil {
			return result, p.ErrInfo(err)
		}
		if result[1] > 0 {
			rez, err := p.Single("SELECT referral FROM users WHERE user_id  =  ?", result[1]).Int64()
			result[2] = rez
			if err != nil {
				return result, p.ErrInfo(err)
			}
		}
	}
	return result, nil
}

func (p *Parser) getTdc(promisedAmountId, userId int64) (float64, error) {
	// используем $this->tx_data['time'], оно всегда меньше времени блока, а значит TDC будет тут чуть меньше. В блоке (не фронт. проверке) уже будет использоваться time из блока
	var time int64
	if p.BlockData!=nil {
		time = p.BlockData.Time
	} else {
		time = p.TxTime
	}
	pct, err := p.GetPct()
	if err != nil {
		return 0, err
	}
	maxPromisedAmounts, err := p.GetMaxPromisedAmounts()
	if err != nil {
		return 0, err
	}
	log.Println("pct", pct)
	log.Println("maxPromisedAmounts", maxPromisedAmounts)

	var status string
	var amount, tdc_amount float64
	var currency_id, tdc_amount_update int64
	err = p.QueryRow("SELECT status, amount, currency_id, tdc_amount, tdc_amount_update FROM promised_amount WHERE id  =  ?", promisedAmountId).Scan(&status, &amount, &currency_id, &tdc_amount, &tdc_amount_update)
	if err != nil  && err!=sql.ErrNoRows {
		return 0, err
	}
	pointsStatus, err := p.getPointsStatus(userId);
	if err != nil {
		return 0, err
	}
	userHolidays, err := p.GetHolidays(userId);
	if err != nil {
		return 0, err
	}
	existsCashRequests := p.CheckCashRequests(userId)
	var newTdc float64
	// для WOC майнинг не зависит от неудовлетворенных cash_requests, т.к. WOC юзер никому не обещал отдавать. Также, WOC не бывает repaid
	if status == "mining" && (existsCashRequests==nil || currency_id==1) {
		repeadAmount, err:=p.GetRepaidAmount(currency_id, userId)
		if err != nil {
			return 0, err
		}
		profit, err := p.calcProfit_ ( amount+tdc_amount, tdc_amount_update, time, pct[currency_id], pointsStatus, userHolidays, maxPromisedAmounts[currency_id], currency_id, repeadAmount );
		if err != nil {
			return 0, err
		}
		newTdc = tdc_amount + profit
	} else if status == "repaid" && existsCashRequests==nil {
		profit, err := p.calcProfit_ ( tdc_amount, tdc_amount_update, time, pct[currency_id], pointsStatus, [][]int64{}, []map[int64]string{}, 0, 0 )
		if err != nil {
			return 0, err
		}
		newTdc = tdc_amount + profit
	} else { // rejected/change_geo/suspended
		newTdc = tdc_amount
	}
	return newTdc, nil
}

// откат не всех полей, а только указанных, либо 1 строку, если нет where
func (p *Parser) selectiveRollback(fields []string, table string, where string, rollback bool) error {
	if len(where) > 0 {
		where = " WHERE "+where
	}
	addSqlFields := ""
	for _, field := range fields {
		addSqlFields+=field+","
	}
	// получим log_id, по которому можно найти данные, которые были до этого
	logId, err := p.Single("SELECT log_id FROM ? "+where, table).Int64()
	if err != nil {
		return err
	}
	if logId > 0 {
		// данные, которые восстановим
		logData, err := p.OneRow("SELECT "+addSqlFields+" prev_log_id FROM log_"+table+" WHERE log_id  =  ?", logId).String()
		if err != nil {
			return err
		}
		addSqlUpdate:=""
		for _, field := range fields {
			if utils.InSliceString(field, []string{"hash", "tx_hash", "public_key_0", "public_key_1", "public_key_2"}) && len(logData[field])!=0 {
				query:=""
				switch p.ConfigIni["dbType"] {
				case "sqlite":
					query = field+`='`+logData[field]+`',`
				case "postgresql":
					query = field+`=decode(`+logData[field]+`,'HEX'),`
				case "mysql":
					query = field+`=UNHEX("`+logData[field]+`"),`
				}
				addSqlUpdate+=query
			} else {
				addSqlUpdate+= field+`='`+logData[field]+`',`
			}
		}
		err = p.ExecSql("UPDATE ? SET "+addSqlUpdate+" log_id = ? "+where, table, logData["prev_log_id"])
		if err != nil {
			return err
		}
		// подчищаем _log
		err = p.ExecSql("DELETE FROM log_"+table+" WHERE log_id = ?", logId)
		if err != nil {
			return err
		}
		p.rollbackAI("log_"+table, 1)
	} else {
		err = p.ExecSql("DELETE FROM ? "+where, table)
		if err != nil {
			return err
		}
		if rollback {
			p.rollbackAI(table, 1)
		}
	}

	return nil
}

/**
 *
Вычисляем, какой получится профит от суммы $amount
$pct_array = array(
	1394308460=>array('user'=>0.05, 'miner'=>0.10),
	1394308470=>array('user'=>0.06, 'miner'=>0.11),
	1394308480=>array('user'=>0.07, 'miner'=>0.12),
	1394308490=>array('user'=>0.08, 'miner'=>0.13)
	);
 * $holidays_array = array ($start, $end)
 * $points_status_array = array(
	1=>'user',
	9=>'miner',
	10=>'user',
	12=>'miner'
 * );
 * $max_promised_amount_array = array(
	1394308460=>7500,
	1394308471=>2500,
	1394308482=>7500,
	1394308490=>5000
	);
 * $repaid_amount, $holidays_array, $points_status_array, $max_promised_amount_array нужны только для обещанных сумм. у погашенных нет $repaid_amount, $holidays_array, $max_promised_amount_array
 * $repaid_amount нужен чтобы узнать, не будет ли превышения макс. допустимой суммы. считаем amount mining+repaid
 * $currency_id - для иднетификации WOC
 * */
type pctAmount struct {
	pct float64
	amount float64
}
type resultArrType struct {
	num_sec int64
	pct float64
	amount float64
}
func (p *Parser) CalcProfit(amount float64, timeStart, timeFinish int64, pctArray []map[int64]map[string]float64, pointsStatusArray []map[int64]string, holidaysArray [][]int64, maxPromisedAmountArray []map[int64]string, currencyId int64, repaidAmount float64) (float64, error) {
	if timeStart >= timeFinish {
		return 0, nil
	}

	// для WOC майнинг останавливается только если майнера забанил админ, каникулы на WOC не действуют
	if currencyId == 1 {
		holidaysArray = nil
	}

	/* $max_promised_amount_array имеет дефолтные значения от времени = 0
	 * $pct_array имеет дефолтные значения 0% для user/miner от времени = 0
	 * в $points_status_array крайний элемент массива всегда будет относиться к текущим 30-и дням т.к. перед calc_profit всегда идет вызов points_update
	 * */

	var lastStatus string = ""
	var findMinArray []map[string]string
	var newArr []map[int64]float64
	var statusPctArray_ map[string]float64
	// нужно получить массив вида time=>pct, совместив $pct_array и $points_status_array

	findTime := func(key int64, arr []map[int64]float64) bool {
		for i:=0; i< len(arr); i++ {
			if arr[i][key]!=0 {
				return true
			}
		}
		return false
	}

	log.Println("pctArray", pctArray)
	for i:=0; i < len(pctArray); i++ {
		for time, statusPctArray := range pctArray[i] {
			log.Println("i=", i, "pctArray[i]=", pctArray[i])
			findMinArray, pointsStatusArray = findMinPointsStatus(time, pointsStatusArray, "status")
			//log.Println("i", i)
			log.Println("time", time)
			log.Println("findMinArray", findMinArray)
			log.Println("pointsStatusArray", pointsStatusArray)
			for j := 0; j < len(findMinArray); j++ {
				if utils.StrToInt64(findMinArray[j]["time"]) <= time {
					findMinPct := findMinPct(utils.StrToInt64(findMinArray[j]["time"]), pctArray, findMinArray[j]["status"]);
					if !findTime(utils.StrToInt64(findMinArray[j]["time"]), newArr) {
						newArr = append(newArr, map[int64]float64{utils.StrToInt64(findMinArray[j]["time"]) : findMinPct})
						log.Println("findMinPct", findMinPct)
					}
					lastStatus = findMinArray[j]["status"];
				}
			}
			if len(findMinArray) == 0 && len(lastStatus) == 0 {
				findMinArray = append(findMinArray, map[string]string{"status": "user"})
			} else if len(findMinArray) == 0 && len(lastStatus) != 0 { // есть проценты, но кончились points_status
				findMinArray = append(findMinArray, map[string]string{"status": lastStatus})
			}
			if !findTime(time, newArr) {
				newArr = append(newArr, map[int64]float64{time : statusPctArray[findMinArray[len(findMinArray)-1]["status"]]})
			}
			statusPctArray_ = statusPctArray;
		}
	}

	// если в points больше чем в pct
	if len(pointsStatusArray)>0 {
		for i:=0; i < len(pointsStatusArray); i++ {
			for time, status := range pointsStatusArray[i] {
				if !findTime(time, newArr) {
					newArr = append(newArr, map[int64]float64{time : statusPctArray_[status]})
				}
			}
		}
	}

	log.Println("newArr", newArr)


	// newArr - массив, где ключи - это время из pct и points_status, а значения - проценты.

	// $max_promised_amount_array + $pct_array
	/*
	 * в $pct_array сейчас
			[1394308000] =>  0,05
			[1394308100] =>  0,1

		после обработки станет

			[1394308000] => Array
				(
					[pct] => 0,05
					[amount] => 1000
				)
			[1394308005] => Array
				(
					[pct] => 0,05
					[amount] => 100
				)
			[1394308100] => Array
				(
					[pct] => 0,1
					[amount] => 100
				)

	 * */

	findTime2 := func(key int64, arr []map[int64]pctAmount) bool {
		for i:=0; i< len(arr); i++ {
			if arr[i][key].pct!=0 {
				return true
			}
		}
		return false
	}

	var newArr2 []map[int64]pctAmount
	var lastAmount float64
	var amount_ float64
	var pct_ float64
	if len(maxPromisedAmountArray)==0{
		lastAmount = amount
	}

	// нужно получить массив вида time=>pct, совместив newArr и $max_promised_amount_array
	for i:=0; i < len(newArr); i++ {
		log.Println("i ", i)
		for time, pct := range newArr[i] {
			findMinArray, maxPromisedAmountArray = findMinPointsStatus(time, maxPromisedAmountArray, "amount")
			for j:=0; j < len(findMinArray); j++ {
				amount_ = getMaxPromisedAmountCalcProfit(amount, repaidAmount, utils.StrToFloat64(findMinArray[j]["amount"]), currencyId)
				if utils.StrToInt64(findMinArray[j]["time"]) <= time {
					minPct := findMinPct1(utils.StrToInt64(findMinArray[j]["time"]), newArr);
					if !findTime2(utils.StrToInt64(findMinArray[j]["time"]), newArr2) {
						newArr2 = append(newArr2, map[int64]pctAmount{utils.StrToInt64(findMinArray[j]["time"]):{pct:minPct, amount:amount_}})
					}
					lastAmount = amount_;

				}
			}
			if !findTime2(time, newArr2) {
				newArr2 = append(newArr2, map[int64]pctAmount{time:{pct:pct, amount:lastAmount}})
			}
			pct_ = pct;
		}
	}

	// если в max_promised_amount больше чем в pct
	if len(maxPromisedAmountArray) > 0 {
		for i:=0; i<len(maxPromisedAmountArray); i++ {
			for time, maxPromisedAmount := range maxPromisedAmountArray[i] {
				MaxPromisedAmountCalcProfit := getMaxPromisedAmountCalcProfit(amount, repaidAmount, utils.StrToFloat64(maxPromisedAmount), currencyId);
				amount_ = MaxPromisedAmountCalcProfit
				if !findTime2(time, newArr2) {
					newArr2 = append(newArr2, map[int64]pctAmount{time:{pct:pct_, amount:MaxPromisedAmountCalcProfit}})
				}
			}
		}
	}

	log.Println("newArr2", newArr2)

	maxTimeInNewArr2 := func(newArr2 []map[int64]pctAmount) int64 {
		var max int64
		for i:=0; i < len(newArr2); i++ {
			for time, _ := range newArr2[i] {
				if time > max {
					max = time
				}
			}
		}
		return max
	}

	if maxTimeInNewArr2(newArr2) < timeFinish {
		// добавим сразу время окончания
		//newArr2[timeFinish] = pct;
		if !findTime2(timeFinish, newArr2) {
			newArr2 = append(newArr2, map[int64]pctAmount{timeFinish:{pct:pct_, amount:0}})
		}
	}

	var workTime, oldTime int64
	var resultArr []resultArrType
	var oldPctAndAmount pctAmount
	var startHolidays bool
	var finishHolidaysElement int64
	START:
	for i:=0; i < len(newArr2); i++ {

		for time, pctAndAmount := range newArr2[i] {

			log.Println("pctAndAmount", pctAndAmount)

			if (time > timeFinish) {
				continue START
			}
			if (time > timeStart) {
				workTime = time
				for j := 0; j < len(holidaysArray); j++ {

					if holidaysArray[j][1] <= oldTime {
						continue
					}

					log.Println("holidaysArray[j]", holidaysArray[j])

					// полные каникулы в промежутке между time и old_time
					if holidaysArray[j][0]!=-1 && oldTime <= holidaysArray[j][0] && holidaysArray[j][1]!=-1 && workTime >= holidaysArray[j][1] {
						time = holidaysArray[j][0];
						holidaysArray[j][0] = -1
						resultArr = append(resultArr, resultArrType{num_sec : (time-oldTime), pct : oldPctAndAmount.pct, amount : oldPctAndAmount.amount})
						log.Println("resultArr append")
						oldTime = holidaysArray[j][1];
						holidaysArray[j][1] = -1
					}
					if holidaysArray[j][0]!=-1 && workTime >= holidaysArray[j][0] {
						startHolidays = true; // есть начало каникул, но есть ли конец?
						finishHolidaysElement = holidaysArray[j][1]; // для записи в лог
						time = holidaysArray[j][0];
						if time < timeStart {
							time = timeStart
						}
						holidaysArray[j][0] = -1
					} else if holidaysArray[j][1]!=-1 && workTime < holidaysArray[j][1] && holidaysArray[j][0]==-1 {
						// конец каникул заканчивается после $work_time
						time = oldTime
						continue
					} else if holidaysArray[j][1]!=-1 && workTime >= holidaysArray[j][1] {
						oldTime = holidaysArray[j][1]
						holidaysArray[j][1] = -1
						startHolidays = false; // конец каникул есть
					} else if j == len(holidaysArray)-1 && !startHolidays {
						// если это последний полный внутрений холидей, то time должен быть равен текущему workTime
						time = workTime
					}
				}
				if (time > timeFinish) {
					time = timeFinish
				}
				resultArr = append(resultArr, resultArrType{num_sec : (time-oldTime), pct : oldPctAndAmount.pct, amount : oldPctAndAmount.amount})
				log.Println("new", (time-oldTime))
				oldTime = time
			} else {
				oldTime = timeStart
			}
			oldPctAndAmount = pctAndAmount
		}
	}

	log.Println("resultArr", resultArr)

	if (startHolidays && finishHolidaysElement>0) {
		log.Println("finishHolidaysElement:", finishHolidaysElement)
	}

	// время в процентах меньше, чем нужное нам конечное время
	if (oldTime < timeFinish && !startHolidays) {
		// просто берем последний процент и добиваем его до нужного $time_finish
		sec := timeFinish - oldTime;
		resultArr = append(resultArr, resultArrType{num_sec : sec, pct : oldPctAndAmount.pct, amount : oldPctAndAmount.amount})
	}


	var profit, amountAndProfit float64
	for i:=0; i < len(resultArr); i++ {
		pct := 1+resultArr[i].pct
		num := resultArr[i].num_sec
		amountAndProfit = profit +resultArr[i].amount
		//$profit = ( floor( round( $amount_and_profit*pow($pct, $num), 3)*100 ) / 100 ) - $new[$i]['amount'];
		// из-за того, что в front был подсчет без обновления points, а в рабочем методе уже с обновлением points, выходило, что в рабочем методе было больше мелких временных промежуток, и получалось profit <0.01, из-за этого было расхождение в front и попадание минуса в БД
		profit = amountAndProfit*math.Pow(pct, float64(num)) - resultArr[i].amount
	}

	return profit, nil
}


func  getMaxPromisedAmountCalcProfit(amount, repaidAmount, maxPromisedAmount float64, currencyId int64) float64 {
	// для WOC $repaid_amount всегда = 0, т.к. cash_request на WOC послать невозможно
	// если наша сумма больше, чем максимально допустимая ($find_min_array[$i]['amount'])
	var result float64
	if (amount+repaidAmount > maxPromisedAmount) {
		result = maxPromisedAmount-repaidAmount;
	} else if (amount < maxPromisedAmount && currencyId==1) { // для WOC разрешено брать maxPromisedAmount вместо promisedAmount, если promisedAmount < maxPromisedAmount
		result = maxPromisedAmount
	} else {
		result = amount;
	}
	return result
}


// ищем ближайшее время в $points_status_array или $max_promised_amount_array
// $type - status для $points_status_array / amount - для $max_promised_amount_array
func  findMinPointsStatus(needTime int64, pointsStatusArray []map[int64]string, pType string) ([]map[string]string, []map[int64]string) {
	var findTime []int64
	newPointsStatusArray := pointsStatusArray
	var timeStatusArr []map[string]string
	BR:
	for i:=0; i<len(pointsStatusArray); i++ {
		for time, _ := range pointsStatusArray[i] {
			if time > needTime {
				break BR
			}
			findTime = append(findTime, time)
			start:=i+1
			if i+1 > len(pointsStatusArray) {
				start=len(pointsStatusArray)
			}
			newPointsStatusArray = pointsStatusArray[start:]
		}
	}
	if len(findTime) > 0 {
		for i:=0; i<len(findTime); i++ {
			for _, status := range pointsStatusArray[i] {
				timeStatusArr = append(timeStatusArr, map[string]string{"time" : utils.Int64ToStr(findTime[i]), pType : status})
			}
		}
	}
	return timeStatusArr, newPointsStatusArray
}

func findMinPct (needTime int64, pctArray []map[int64]map[string]float64, status string) float64 {
	var findTime int64 = -1
	var pct float64 = 0
	BR:
	for i:=0; i<len(pctArray); i++ {
		for time, _ := range pctArray[i] {
			if time > needTime {
				break BR
			}
			findTime = int64(i)
		}
	}
	if findTime >=0 {
		for _, arr := range pctArray[findTime] {
			pct = arr[status]
		}
	}
	return pct
}


func findMinPct1 (needTime int64, pctArray []map[int64]float64) float64 {
	var findTime int64 = -1
	var pct float64 = 0
	BR:
	for i:=0; i<len(pctArray); i++ {
		for time, _ := range pctArray[i] {
			if time > needTime {
				break BR
			}
			findTime = int64(i)
		}
	}
	if findTime >=0 {
		for _, pct0 := range pctArray[findTime] {
			pct = pct0
		}
	}
	return pct
}

// только для блоков до 24946
func findMinPct1_24946(needTime int64, pctArray []map[int64]float64) float64 {
	var findTime int64 = 0
	var pct float64 = 0
	BR:
	for i:=0; i<len(pctArray); i++ {
		for time, _ := range pctArray[i] {
			if time > needTime {
				break BR
			}
			findTime = int64(i)
		}
	}
	if findTime >=0 {
		for _, pct0 := range pctArray[findTime] {
			pct = pct0
		}
	}
	return pct
}


// только для блоков до 24946
func findMinPct_24946 (needTime int64, pctArray []map[int64]map[string]float64, status string) float64 {
	var findTime int64 = 0
	var pct float64 = 0
	log.Println("pctArray findMinPct_24946", pctArray)
	BR:
	for i:=0; i<len(pctArray); i++ {
		for time, _ := range pctArray[i] {
			log.Println(time, ">", needTime, "?")
			if time > needTime {
				log.Println("break")
				break BR
			}
			findTime = int64(i)
		}
	}
	log.Println("findTime", findTime)
	if findTime >0 {
		for _, arr := range pctArray[findTime] {
			pct = arr[status]
		}
	}
	return pct
}

func (p *Parser) CalcProfit_24946(amount float64, timeStart, timeFinish int64, pctArray []map[int64]map[string]float64, pointsStatusArray []map[int64]string, holidaysArray [][]int64, maxPromisedAmountArray []map[int64]string, currencyId int64, repaidAmount float64) (float64, error) {


	var lastStatus string = ""
	var findMinArray []map[string]string
	var newArr []map[int64]float64
	var statusPctArray_ map[string]float64
	// нужно получить массив вида time=>pct, совместив $pct_array и $points_status_array

	findTime := func(key int64, arr []map[int64]float64) bool {
		for i:=0; i< len(arr); i++ {
			if arr[i][key]!=0 {
				return true
			}
		}
		return false
	}

	log.Println("pctArray", pctArray)
	for i:=0; i < len(pctArray); i++ {
		for time, statusPctArray := range pctArray[i] {
			log.Println("i=", i, "pctArray[i]=", pctArray[i])
			findMinArray, pointsStatusArray = findMinPointsStatus(time, pointsStatusArray, "status")
			//log.Println("i", i)
			log.Println("time", time)
			log.Println("findMinArray", findMinArray)
			log.Println("pointsStatusArray", pointsStatusArray)
			for j := 0; j < len(findMinArray); j++ {
				if utils.StrToInt64(findMinArray[j]["time"]) < time {
					findMinPct := findMinPct_24946(utils.StrToInt64(findMinArray[j]["time"]), pctArray, findMinArray[j]["status"]);
					if !findTime(utils.StrToInt64(findMinArray[j]["time"]), newArr) {
						newArr = append(newArr, map[int64]float64{utils.StrToInt64(findMinArray[j]["time"]) : findMinPct})
						log.Println("findMinPct", findMinPct)
					}
					lastStatus = findMinArray[j]["status"];
				}
			}
			if len(findMinArray) == 0 && len(lastStatus) == 0 {
				findMinArray = append(findMinArray, map[string]string{"status": "user"})
			} else if len(findMinArray) == 0 && len(lastStatus) != 0 { // есть проценты, но кончились points_status
				findMinArray = append(findMinArray, map[string]string{"status": "miner"})
			}
			if !findTime(time, newArr) {
				newArr = append(newArr, map[int64]float64{time : statusPctArray[findMinArray[len(findMinArray)-1]["status"]]})
			}
			statusPctArray_ = statusPctArray;
		}
	}

	// если в points больше чем в pct
	if len(pointsStatusArray)>0 {
		for i:=0; i < len(pointsStatusArray); i++ {
			for time, status := range pointsStatusArray[i] {
				if !findTime(time, newArr) {
					newArr = append(newArr, map[int64]float64{time : statusPctArray_[status]})
				}
			}
		}
	}

	log.Println("newArr", newArr)


	// newArr - массив, где ключи - это время из pct и points_status, а значения - проценты.

	// $max_promised_amount_array + $pct_array
	/*
	 * в $pct_array сейчас
			[1394308000] =>  0,05
			[1394308100] =>  0,1

		после обработки станет

			[1394308000] => Array
				(
					[pct] => 0,05
					[amount] => 1000
				)
			[1394308005] => Array
				(
					[pct] => 0,05
					[amount] => 100
				)
			[1394308100] => Array
				(
					[pct] => 0,1
					[amount] => 100
				)

	 * */

	findTime2 := func(key int64, arr []map[int64]pctAmount) bool {
		for i:=0; i< len(arr); i++ {
			if arr[i][key].pct!=0 {
				return true
			}
		}
		return false
	}

	var newArr2 []map[int64]pctAmount
	var lastAmount float64
	var amount_ float64
	var pct_ float64
	if len(maxPromisedAmountArray)==0{
		lastAmount = amount
	}

	// нужно получить массив вида time=>pct, совместив newArr и $max_promised_amount_array
	for i:=0; i < len(newArr); i++ {
		log.Println("i ", i)
		for time, pct := range newArr[i] {
			findMinArray, maxPromisedAmountArray = findMinPointsStatus(time, maxPromisedAmountArray, "amount")
			for j:=0; j < len(findMinArray); j++ {
				if amount+repaidAmount > utils.StrToFloat64(findMinArray[j]["amount"]) {
					amount_ = utils.StrToFloat64(findMinArray[j]["amount"]) - repaidAmount
				} else if amount < utils.StrToFloat64(findMinArray[j]["amount"]) && currencyId==1 {
					amount_ = utils.StrToFloat64(findMinArray[j]["amount"])
				} else {
					amount_ = amount
				}
				if utils.StrToInt64(findMinArray[j]["time"]) <= time {
					minPct := findMinPct1_24946(utils.StrToInt64(findMinArray[j]["time"]), newArr);
					if !findTime2(utils.StrToInt64(findMinArray[j]["time"]), newArr2) {
						newArr2 = append(newArr2, map[int64]pctAmount{utils.StrToInt64(findMinArray[j]["time"]):{pct:minPct, amount:amount_}})
					}
					lastAmount = amount_;

				}
			}
			if !findTime2(time, newArr2) {
				newArr2 = append(newArr2, map[int64]pctAmount{time:{pct:pct, amount:lastAmount}})
			}
			pct_ = pct
		}
	}


	log.Println("newArr2", newArr2)

	if !findTime2(timeFinish, newArr2) {
		newArr2 = append(newArr2, map[int64]pctAmount{timeFinish:{pct:pct_, amount:0}})
	}

	var workTime, oldTime int64
	var resultArr []resultArrType
	var oldPctAndAmount pctAmount
	var startHolidays bool
	var finishHolidaysElement int64
//START:
	for i:=0; i < len(newArr2); i++ {

		for time, pctAndAmount := range newArr2[i] {

			log.Println("pctAndAmount", pctAndAmount)

			if (time > timeStart) {
				workTime = time
				for j := 0; j < len(holidaysArray); j++ {

					if holidaysArray[j][1] <= oldTime {
						continue
					}

					log.Println("holidaysArray[j]", holidaysArray[j])

					// полные каникулы в промежутке между time и old_time
					if holidaysArray[j][0]!=-1 && workTime >= holidaysArray[j][0] && holidaysArray[j][1]!=-1 && workTime >= holidaysArray[j][1] {
						time = holidaysArray[j][0];
						holidaysArray[j][0] = -1
						resultArr = append(resultArr, resultArrType{num_sec : (time-oldTime), pct : oldPctAndAmount.pct, amount : oldPctAndAmount.amount})
						log.Println("resultArr append")
						oldTime = holidaysArray[j][1];
						holidaysArray[j][1] = -1
					}
					if holidaysArray[j][0]!=-1 && workTime >= holidaysArray[j][0] {
						startHolidays = true; // есть начало каникул, но есть ли конец?
						finishHolidaysElement = holidaysArray[j][1]; // для записи в лог
						time = holidaysArray[j][0];
						if time < timeStart {
							time = timeStart
						}
						holidaysArray[j][0] = -1
					} else if holidaysArray[j][1]!=-1 && workTime < holidaysArray[j][1] && holidaysArray[j][0]==-1 {
						// конец каникул заканчивается после $work_time
						time = oldTime
						continue
					} else if holidaysArray[j][1]!=-1 && workTime >= holidaysArray[j][1] {
						oldTime = holidaysArray[j][1]
						holidaysArray[j][1] = -1
						startHolidays = false; // конец каникул есть
					} else if j == len(holidaysArray)-1 && !startHolidays {
						// если это последний полный внутрений холидей, то time должен быть равен текущему workTime
						time = workTime
					}
				}
				if (time > timeFinish) {
					time = timeFinish
				}
				resultArr = append(resultArr, resultArrType{num_sec : (time-oldTime), pct : oldPctAndAmount.pct, amount : oldPctAndAmount.amount})
				log.Println("new", (time-oldTime))
				oldTime = time
			} else {
				oldTime = timeStart
			}
			oldPctAndAmount = pctAndAmount
		}
	}

	log.Println("resultArr", resultArr)

	if (startHolidays && finishHolidaysElement>0) {
		log.Println("finishHolidaysElement:", finishHolidaysElement)
	}

	// время в процентах меньше, чем нужное нам конечное время
	if (oldTime < timeFinish && !startHolidays) {
		// просто берем последний процент и добиваем его до нужного $time_finish
		sec := timeFinish - oldTime;
		resultArr = append(resultArr, resultArrType{num_sec : sec, pct : oldPctAndAmount.pct, amount : oldPctAndAmount.amount})
	}


	var profit, amountAndProfit float64
	for i:=0; i < len(resultArr); i++ {
		pct := 1+resultArr[i].pct
		num := resultArr[i].num_sec
		amountAndProfit = profit +resultArr[i].amount
		//$profit = ( floor( round( $amount_and_profit*pow($pct, $num), 3)*100 ) / 100 ) - $new[$i]['amount'];
		// из-за того, что в front был подсчет без обновления points, а в рабочем методе уже с обновлением points, выходило, что в рабочем методе было больше мелких временных промежуток, и получалось profit <0.01, из-за этого было расхождение в front и попадание минуса в БД
		profit = amountAndProfit*math.Pow(pct, float64(num)) - resultArr[i].amount
	}
	log.Println("profit", profit)

	return profit, nil
}


func (p *Parser) calcNodeCommission(amount float64, nodeCommission [3]float64) float64 {
	pct := nodeCommission[0]
	minCommission := nodeCommission[1];
	maxCommission := nodeCommission[2];
	nodeCommissionResult := utils.Round ( (amount / 100) * pct , 2 )
	log.Println("nodeCommissionResult", nodeCommissionResult, "amount", amount, "pct", pct)
	if (nodeCommissionResult < minCommission) {
		nodeCommissionResult = minCommission
	} else if (nodeCommissionResult > maxCommission) {
		nodeCommissionResult = maxCommission
	}
	return nodeCommissionResult
}

func (p *Parser) getMyNodeCommission(currencyId, userId int64, amount float64) (float64, error) {
	var nodeCommission float64
	if currencyId >= 1000 {
		currencyId = 1000
	}
	// если это тр-ия без блока, то комиссию нода берем у себя
	if p.BlockData==nil {
		_, _, _, myUserIds , err:= p.GetMyUserId(userId)
		if err != nil {
			return 0, p.ErrInfo(err)
		}

		var commissionJson []byte
		// один элемент в  my_user_ids - это сингл мод
		if len(myUserIds) == 1 {
			commissionJson, err = p.Single("SELECT commission FROM commission WHERE user_id  =  ?", myUserIds[0]).Bytes()
			if err != nil {
				return 0, p.ErrInfo(err)
			}
		} else {
			// если работаем в режиме пула, тогда комиссию берем из config, т.к. майнеры в пуле, у кого комиссиия больше не смогут генерить блоки
			commissionJson, err = p.Single("SELECT commission FROM config").Bytes()
			if err != nil {
				return 0, p.ErrInfo(err)
			}
		}
		commissionMap := make(map[string][3]float64)
		err = json.Unmarshal(commissionJson, &commissionMap)
		if err != nil {
			return 0, p.ErrInfo(err)
		}
		var tmpNodeCommission float64
		currencyIdStr:=utils.Int64ToStr(currencyId)
		if len(commissionMap[currencyIdStr]) > 0 {
			if len(commissionMap[currencyIdStr]) !=3 {
				return 0, p.ErrInfo(err)
			}
			tmpNodeCommission = p.calcNodeCommission(amount, commissionMap[currencyIdStr])
		} else {
			tmpNodeCommission = 0
		}
		if tmpNodeCommission > nodeCommission {
			nodeCommission = tmpNodeCommission
		}
	} else {	// если же тр-ия уже в блоке, то берем комиссию у юзера, который сгенерил этот блок
		commissionJson, err := p.Single("SELECT commission FROM commission WHERE user_id  =  ?", p.BlockData.UserId).Bytes()
		if err != nil {
			return 0, p.ErrInfo(err)
		}
		if len(commissionJson) == 0 {
			nodeCommission = 0
		} else {
			commissionMap := make(map[string][3]float64)
			err = json.Unmarshal(commissionJson, &commissionMap)
			if err != nil {
				return 0, p.ErrInfo(err)
			}
			currencyIdStr:=utils.Int64ToStr(currencyId)
			if len(commissionMap[currencyIdStr]) > 0 {
				log.Println("commissionMap[currencyIdStr]", commissionMap[currencyIdStr])
				nodeCommission = p.calcNodeCommission(amount, commissionMap[currencyIdStr])
				log.Println("nodeCommission", nodeCommission)
			} else {
				nodeCommission = 0
			}
		}
	}
	return nodeCommission, nil

}

func (p *Parser) limitRequestsMoneyOrders(limit int64) (error) {
	num, err := p.Single("SELECT count(tx_hash) FROM log_time_money_orders WHERE user_id  =  ? AND del_block_id  =  0", p.TxUserID).Int64()
	if err != nil {
		return p.ErrInfo(err)
	}
	if num >= limit {
		return p.ErrInfo("[limit_requests] log_time_money_orders")
	} else {
		err = p.ExecSql("INSERT INTO log_time_money_orders ( tx_hash, user_id ) VALUES ( [hex], ? )", p.TxHash, p.TxUserID)
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	return nil
}

func (p *Parser) getWalletsBufferAmount(currencyId int64) (float64, error) {
	return p.Single("SELECT sum(amount) FROM wallets_buffer WHERE user_id = ? AND currency_id = ? AND del_block_id = 0", p.TxUserID, currencyId).Float64()
}


func (p *Parser) getTotalAmount(currencyId int64) (float64, error) {
	var amount float64
	var last_update int64
	err := p.QueryRow("SELECT amount, last_update FROM wallets WHERE user_id = ? AND currency_id = ?", p.TxUserID, currencyId).Scan(&amount, &last_update)
	if err != nil && err!=sql.ErrNoRows {
		return 0, p.ErrInfo(err)
	}
	log.Println("getTotalAmount amount", amount, "p.TxUserID=", p.TxUserID, "currencyId=", currencyId)
	pointsStatus := []map[int64]string {{0:"user"}}
	// getTotalAmount используется только на front, значит используем время из тр-ии - $this->tx_data['time']
	if currencyId >= 1000 { // >=1000 - это CF-валюты, которые не растут
		return amount, nil
	} else {
		pct, err := p.GetPct()
		if err != nil {
			return 0, p.ErrInfo(err)
		}
		profit, err := p.calcProfit_(amount, last_update, p.TxTime, pct[currencyId], pointsStatus, [][]int64{}, []map[int64]string{}, 0, 0)
		if err != nil {
			return 0, p.ErrInfo(err)
		}
		return (amount + profit), nil
	}
	return 0, nil
}

func (p *Parser) updPromisedAmountsRollback(userId int64, cashRequestOutTime bool) error {

	sqlNameCashRequestOutTime := ""
	sqlUpdCashRequestOutTime := ""
	if (cashRequestOutTime) {
		sqlNameCashRequestOutTime = "cash_request_out_time, "
	}

	// идем в обратном порядке (DESC)
	rows, err := p.Query("SELECT log_id FROM promised_amount WHERE status IN ('mining', 'repaid') AND user_id = ? AND currency_id > 1 AND del_block_id = 0 AND del_mining_block_id = 0 ORDER BY id DESC", userId)
	if err != nil {
		return p.ErrInfo(err)
	}
	defer rows.Close()
	for rows.Next() {
		var log_id int64
		err = rows.Scan(&log_id)
		if err != nil {
			return p.ErrInfo(err)
		}
		// данные, которые восстановим
		logData, err := p.OneRow("SELECT tdc_amount, tdc_amount_update, "+sqlNameCashRequestOutTime+" prev_log_id FROM log_promised_amount WHERE log_id  =  ?", log_id).String()
		if err != nil {
			return p.ErrInfo(err)
		}
		if cashRequestOutTime {
			sqlUpdCashRequestOutTime = "cashRequestOutTime = "+logData["cash_request_out_time"]+", ";
		}
		err = p.ExecSql("UPDATE promised_amount SET tdc_amount = ?, tdc_amount_update = ?, "+sqlUpdCashRequestOutTime+" log_id = ? WHERE log_id = ?", logData["tdc_amount"], logData["tdc_amount_update"], logData["prev_log_id"], log_id)
		if err != nil {
			return p.ErrInfo(err)
		}
		// подчищаем _log
		err = p.ExecSql("DELETE FROM log_promised_amount WHERE log_id = ?", log_id)
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.rollbackAI("log_promised_amount", 1)
	}
	return nil
}

func (p *Parser) updPromisedAmounts(userId int64, getTdc, cashRequestOutTimeBool bool, cashRequestOutTime int64) error {
	sqlNameCashRequestOutTime := ""
	sqlValueCashRequestOutTime := ""
	sqlUdpCashRequestOutTime := ""
	if (cashRequestOutTimeBool) {
		sqlNameCashRequestOutTime = "cash_request_out_time, "
		sqlUdpCashRequestOutTime = "cash_request_out_time = "+utils.Int64ToStr(cashRequestOutTime)+", "
	}
	rows, err := p.Query(`
				SELECT  id,
							 currency_id,
							 amount,
							 tdc_amount,
							 tdc_amount_update,
							 `+sqlNameCashRequestOutTime+`
							 log_id
				FROM promised_amount
				WHERE status IN ('mining', 'repaid') AND
							 user_id = ? AND
							 currency_id > 1 AND
							 del_block_id = 0 AND
							 del_mining_block_id = 0
				ORDER BY id ASC
	`, userId)
	if err != nil {
		return p.ErrInfo(err)
	}
	defer rows.Close()
	for rows.Next() {
		var currencyId, tdcAmountUpdate, cashRequestOutTime, amount, log_Id string
		var tdcAmount float64
		var id int64
		if cashRequestOutTimeBool {
			err = rows.Scan(&id, &currencyId, &amount, &tdcAmount, &tdcAmountUpdate, &cashRequestOutTime, &log_Id)
		} else {
			err = rows.Scan(&id, &currencyId, &amount, &tdcAmount, &tdcAmountUpdate, log_Id)
		}
		if err != nil {
			return p.ErrInfo(err)
		}
		if cashRequestOutTimeBool {
			sqlValueCashRequestOutTime = cashRequestOutTime+", "
		}
		logId, err := p.ExecSqlGetLastInsertId(`
					INSERT INTO log_promised_amount (
							tdc_amount,
							tdc_amount_update,
							`+sqlNameCashRequestOutTime+`
							block_id,
							prev_log_id
					)
					VALUES (
							?,
							?,
							`+sqlValueCashRequestOutTime+`
							?,
							?
					)
		`, tdcAmount, tdcAmountUpdate, p.BlockData.BlockId, log_Id)
		if err != nil {
			return p.ErrInfo(err)
		}
		// новая сумма TDC
		var newTdc float64
		if getTdc {
			newTdc, err = p.getTdc(id, userId);
			if err != nil {
				return p.ErrInfo(err)
			}
		} else {
			newTdc = tdcAmount
		}
		err = p.ExecSql("UPDATE promised_amount SET tdc_amount = ?, tdc_amount_update = ?, "+sqlUdpCashRequestOutTime+" log_id = ? WHERE id = ?", newTdc, p.BlockData.Time, logId, id)
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	return nil
}
func (p *Parser) updPromisedAmountsCashRequestOutTime(userId int64) error {
	rows, err := p.Query(`
				SELECT id,
							 cash_request_out_time,
							 log_id
				FROM promised_amount
				WHERE status IN ('mining', 'repaid') AND
							 user_id = ? AND
							 currency_id > 1 AND
							 del_block_id = 0 AND
							 del_mining_block_id = 0 AND
							 cash_request_out_time = 0
				ORDER BY id ASC
	`, userId)
	if err != nil {
		return p.ErrInfo(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, cash_request_out_time, log_id string
		err = rows.Scan(&id, &cash_request_out_time, &log_id)
		if err != nil {
			return p.ErrInfo(err)
		}
		logId, err := p.ExecSqlGetLastInsertId("INSERT INTO log_promised_amount ( cash_request_out_time, block_id, prev_log_id ) VALUES ( ?, ?, ? )", cash_request_out_time, p.BlockData.BlockId, log_id)
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.ExecSql("UPDATE promised_amount SET cash_request_out_time = ?, log_id = ? WHERE id = ?", p.BlockData.Time, logId, id)
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	return nil
}

func (p *Parser) updPromisedAmountsCashRequestOutTimeRollback(userId int64) error {

	// идем в обратном порядке (DESC)
	rows, err := p.Query("SELECT log_id FROM promised_amount WHERE status IN ('mining', 'repaid') AND user_id = ? AND currency_id > 1 AND del_block_id = 0 AND del_mining_block_id = 0 AND cash_request_out_time = ? ORDER BY id DESC", userId, p.BlockData.Time)
	if err != nil {
		return p.ErrInfo(err)
	}
	defer rows.Close()
	for rows.Next() {
		var log_id int64
		err = rows.Scan(&log_id)
		if err != nil {
			return p.ErrInfo(err)
		}
		// данные, которые восстановим
		logData, err := p.OneRow("SELECT cash_request_out_time, prev_log_id FROM log_promised_amount WHERE log_id  =  ?", log_id).Int64()
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.ExecSql("UPDATE promised_amount SET cash_request_out_time = ?, log_id = ? WHERE log_id = ?", logData["cash_request_out_time"], logData["prev_log_id"], log_id)
		if err != nil {
			return p.ErrInfo(err)
		}
		// подчищаем _log
		err = p.ExecSql("DELETE FROM log_promised_amount WHERE log_id = ?", log_id)
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.rollbackAI("log_promised_amount", 1)
		if err != nil {
			return p.ErrInfo(err)
		}
	}
	return nil
}

func (p *Parser) checkSenderMoney(currencyId, fromUserId int64, amount, commission, arbitrator0_commission, arbitrator1_commission, arbitrator2_commission, arbitrator3_commission, arbitrator4_commission float64) (float64, error) {

	// получим все списания (табла wallets_buffer), которые еще не попали в блок и стоят в очереди
	walletsBufferAmount, err := p.getWalletsBufferAmount(currencyId)
	if err != nil {
		return 0, p.ErrInfo(err)
	}
	// получим сумму на кошельке юзера + %
	totalAmount, err := p.getTotalAmount(currencyId)
	if err != nil {
		return 0, p.ErrInfo(err)
	}
	var txTime int64
	if p.BlockData!=nil {// тр-ия пришла в блоке
		txTime = p.BlockData.Time
	} else {
		txTime = time.Now().Unix() - 30  // просто на всякий случай небольшой запас
	}

	// учтем все свежие cash_requests, которые висят со статусом pending
	cashRequestsAmount, err := p.Single("SELECT sum(amount) FROM cash_requests WHERE from_user_id  =  ? AND currency_id  =  ? AND status  =  'pending' AND time > ?", fromUserId, currencyId, (txTime-p.Variables.Int64["cash_request_time"])).Float64()
	if err != nil {
		return 0, p.ErrInfo(err)
	}

	// учитываются все fx-ордеры
	forexOrdersAmount, err := p.Single("SELECT sum(amount) FROM forex_orders WHERE user_id  =  ? AND sell_currency_id  =  ? AND del_block_id  =  0", fromUserId, currencyId).Float64()
	if err != nil {
		return 0, p.ErrInfo(err)
	}

	// учитываем все текущие суммы холдбека
	holdBackAmount, err := p.Single(`
		SELECT sum(sum1) FROM (
			SELECT
			CASE
				WHEN (hold_back_amount - refund - voluntary_refund) < 0 THEN 0
			ELSE (hold_back_amount - refund - voluntary_refund)
			END as sum1
			from orders
			WHERE seller  =  ? AND currency_id  =  ? AND end_time > ?
		) as t1`,
		fromUserId, currencyId, txTime).Float64()
	if err != nil {
		return 0, p.ErrInfo(err)
	}

	amountAndCommission := amount + commission + arbitrator0_commission + arbitrator1_commission + arbitrator2_commission + arbitrator3_commission + arbitrator4_commission
	all := totalAmount - walletsBufferAmount - cashRequestsAmount - forexOrdersAmount - holdBackAmount;
	if all < amountAndCommission {
		return 0, p.ErrInfo(fmt.Sprintf("%f < %f (%f - %f - %f - %f - %f) <  (%f + %f + %f + %f + %f + %f + %f)", all, amountAndCommission, totalAmount, walletsBufferAmount, cashRequestsAmount, forexOrdersAmount, holdBackAmount, amount, commission, arbitrator0_commission, arbitrator1_commission, arbitrator2_commission, arbitrator3_commission, arbitrator4_commission))
	}
	return amountAndCommission, nil
}

func (p *Parser) updateWalletsBuffer(amount float64, currencyId int64) (error) {
	// добавим нашу сумму в буфер кошельков, чтобы юзер не смог послать запрос на вывод всех DC с кошелька.
	hash, err := p.Single("SELECT hash FROM wallets_buffer WHERE hash = [hex]", p.TxHash).String()
	if len(hash) > 0 {
		err = p.ExecSql("UPDATE wallets_buffer SET user_id = ?, currency_id = ?, amount = ? WHERE hash = [hex]", p.TxUserID, currencyId, amount, p.TxHash)
	} else {
		err = p.ExecSql("INSERT INTO wallets_buffer ( hash, user_id, currency_id, amount ) VALUES ( [hex], ?, ?, ? )", p.TxHash, p.TxUserID, currencyId, amount)
	}
	if err != nil {
		return p.ErrInfo(err)
	}
	return nil
}

// нельзя отправить более 10-и ордеров от 1 юзера в 1 блоке с суммой менее эквивалента 0.1$ по текущему курсу этой валюты.
func (p *Parser) checkSpamMoney(currencyId int64, amount float64) (error) {
	if currencyId == consts.USD_CURRENCY_ID {
		if p.TxMaps.Float64["amount"] < 0.1 {
			err := p.limitRequestsMoneyOrders(10)
			if err != nil {
				return p.ErrInfo(err)
			}
		}
	} else {
		// если валюта не доллары, то нужно получить эквивалент на бирже
		dollarEqRate, err := p.Single("SELECT sell_rate FROM forex_orders WHERE sell_currency_id  =  ? AND buy_currency_id  =  ?", currencyId, consts.USD_CURRENCY_ID).Float64()
		if err != nil {
			return p.ErrInfo(err)
		}
		// эквивалент 0.1$
		if dollarEqRate > 0 {
			minAmount := 0.1/dollarEqRate
			if amount < minAmount {
				err = p.limitRequestsMoneyOrders(10)
				if err != nil {
					return p.ErrInfo(err)
				}
			}
		}
	}
	return nil
}
