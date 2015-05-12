package dcparser
import (
	//"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"utils"
	"os"
)

type BlockData struct {
	BlockId int64
	Time int64
	UserId int64
	Level int64
	Sign []byte
}

type Parser struct {
	*utils.DCDB
	TxSlice []string
	TxMap map[string]string
	BlockData BlockData
	BinaryData []byte
	blockHashHex string
	dataType int64
	blockHex string
	variables map[string]string
	CurrentBlockId int64
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

	p.BlockData.BlockId = utils.BinToDec(utils.BytesShift(&p.BinaryData, 4))
	p.BlockData.Time = utils.BinToDec(utils.BytesShift(&p.BinaryData, 4))
	p.BlockData.UserId = utils.BinToDec(utils.BytesShift(&p.BinaryData, 5))
	p.BlockData.Level = utils.BinToDec(utils.BytesShift(&p.BinaryData, 1))
	signSize := utils.DecodeLength(&p.BinaryData)
	p.BlockData.Sign = utils.BytesShift(&p.BinaryData, signSize)
	p.CurrentBlockId = p.BlockData.BlockId
	fmt.Println(p.BlockData)
	fmt.Println(signSize)

	return nil
}


func (p *Parser) CheckBlockHeader() error {
	// инфа о предыдущем блоке (т.е. последнем занесенном)

/*// инфа о предыдущем блоке (т.е. последнем занесенном)
		if (!$this->prev_block) // инфа может быть передана прямо в массиве
			$this->get_prev_block($this->block_data['block_id']-1);
		//$this->get_info_block(); убрано, т.к. CheckBlockHeader используется и при сборе новых блоков при вилке

		// для локальных тестов
		if ($this->prev_block['block_id']==1) {
			$ini_array = parse_ini_file(ABSPATH . "config.ini", true);
			if (isset($ini_array['local']['start_block_id'])) {
				$this->prev_block['block_id'] = $ini_array['local']['start_block_id'];
			}
		}

		debug_print("this->prev_block: ".print_r_hex($this->prev_block), __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__);
		debug_print("this->block_data: ".print_r_hex($this->block_data), __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__);

		// меркель рут нужен для проверки подписи блока, а также проверки лимитов MAX_TX_SIZE и MAX_TX_COUNT
		if ($this->block_data['block_id']==1) {
			$this->global_variables['max_tx_size'] = 1024*1024;
			$first = true;
		}
		else
			$first = false;

		$this->mrkl_root = self::getMrklroot($this->binary_data, $this->global_variables, $first);

		// проверим время
		if ( !check_input_data ($this->block_data['time'], 'int') )
			return 'error time';

		// проверим уровень
		if ( !check_input_data ($this->block_data['level'], 'level') )
			return 'error level';

		// получим значения для сна
		$sleep_data = $this->db->query( __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__, "
					SELECT `value`
					FROM `".DB_PREFIX."variables`
					WHERE `name` = 'sleep'
					", 'fetch_one' );
		$sleep_data = json_decode($sleep_data, true);
		debug_print("sleep_data:", __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__);
		//print_R($sleep_data);

		// узнаем время, которые было затрачено в ожидании is_ready предыдущим блоком
		$is_ready_sleep = testblock::get_is_ready_sleep($this->prev_block['level'], $sleep_data['is_ready']);
		// сколько сек должен ждать нод, перед тем, как начать генерить блок, если нашел себя в одном из уровней.
		$generator_sleep = testblock::get_generator_sleep($this->block_data['level'] , $sleep_data['generator']);
		// сумма is_ready всех предыдущих уровней, которые не успели сгенерить блок
		$is_ready_sleep2 = testblock::get_is_ready_sleep_sum($this->block_data['level'] , $sleep_data['is_ready']);

		debug_print("is_ready_sleep={$is_ready_sleep}\ngenerator_sleep={$generator_sleep}\nis_ready_sleep2={$is_ready_sleep2}", __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__);

		debug_print('prev_block:'.print_r_hex($this->prev_block), __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__);
		debug_print('block_data:'.print_r_hex($this->block_data), __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__);

		// не слишком ли рано прислан этот блок. допустима погрешность = error_time
		if (!$first)
		if ( $this->prev_block['time'] + $is_ready_sleep + $generator_sleep + $is_ready_sleep2 - $this->block_data['time'] > $this->global_variables['error_time'] )
			return "error block time {$this->prev_block['time']} + {$is_ready_sleep} + {$generator_sleep} + {$is_ready_sleep2} - {$this->block_data['time']} > {$this->global_variables['error_time']}\n";
		// исключим тех, кто сгенерил блок с бегущими часами
		if ( $this->block_data['time'] > time() )
			return "error block time";

		// проверим ID блока
		if ( !check_input_data ($this->block_data['block_id'], 'int') )
			return 'block_id';

		// проверим, верный ли ID блока
		if (!$first)
			if ( $this->block_data['block_id'] != $this->prev_block['block_id']+1 )
				return "error block_id ({$this->block_data['block_id'] }!=".($this->prev_block['block_id']+1).")";

		// проверим, есть ли такой майнер и заодно получим public_key
		// ================================== my_table заменить ===============================================
		$this->node_public_key = $this->db->query( __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__, "
					SELECT `node_public_key`
					FROM `".DB_PREFIX."miners_data`
					WHERE `user_id` = {$this->block_data['user_id']}
					LIMIT 1
					", 'fetch_one' );

		if (!$first)
			if  ( !$this->node_public_key )
				return 'user_id';

		// SIGN от 128 байта до 512 байт. Подпись от TYPE, BLOCK_ID, PREV_BLOCK_HASH, TIME, USER_ID, LEVEL, MRKL_ROOT
		$for_sign = "0,{$this->block_data['block_id']},{$this->prev_block['hash']},{$this->block_data['time']},{$this->block_data['user_id']},{$this->block_data['level']},{$this->mrkl_root}";
		debug_print("checkSign", __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__);
		// проверяем подпись
		if (!$first) {
			$error = self::checkSign ($this->node_public_key, $for_sign, $this->block_data['sign'], true);
			if ($error)
				return $error;
		}
*/
}

func (p *Parser) ParseDataFull() error {
	p.dataPre()
	if p.dataType != 0  { // парсим только блоки
		return fmt.Errorf("incorrect dataType")
	}
	var err error
	p.variables, err = p.DCDB.GetAllVariables()
	if err != nil {
		return err
	}
	err = p.ParseBlock()
	if err != nil {
		return err
	}

	// проверим данные, указанные в заголовке блока
	err = p.CheckBlockHeader()
	if err != nil {
		return err
	}


	return nil
}

func (p *Parser) GetTxMap(fields []string) (map[string]string, error) {
	//fmt.Println("p.TxSlice", p.TxSlice)
	//fmt.Println("fields", fields)
	if len(p.TxSlice) != len(fields)+4 {
		return nil, fmt.Errorf("bad transaction_array %d != %d (type=%d)",  len(p.TxSlice),  len(fields)+4, p.TxSlice[0])
	}
	TxMap := make(map[string]string)
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
