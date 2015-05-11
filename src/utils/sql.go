package utils
import (
	 "fmt"
	 _ "github.com/lib/pq"
     _ "github.com/go-sql-driver/mysql"
	 "database/sql"
	"strings"
	"regexp"
	//"errors"
	"log"
	"time"
	"strconv"
	"encoding/json"
	//"github.com/astaxie/beego/config"
	"consts"
//	"os"
)

type DCDB struct {
	 *sql.DB
	configIni map[string]string
}

func NewDbConnect(configIni map[string]string) (*DCDB, error) {
	var db *sql.DB
	var err error
	switch configIni["db_type"] {
	case "sqlite":

		db, err = sql.Open("sqlite3", "./litedb.db")
		if err!=nil {
			return &DCDB{}, err
		}
		ddl := `
				PRAGMA automatic_index = ON;
				PRAGMA cache_size = 32768;
				PRAGMA cache_spill = OFF;
				PRAGMA foreign_keys = ON;
				PRAGMA journal_size_limit = 67110000;
				PRAGMA locking_mode = NORMAL;
				PRAGMA page_size = 4096;
				PRAGMA recursive_triggers = ON;
				PRAGMA secure_delete = ON;
				PRAGMA synchronous = NORMAL;
				PRAGMA temp_store = MEMORY;
				PRAGMA journal_mode = WAL;
				PRAGMA wal_autocheckpoint = 16384;
				PRAGMA encoding = "UTF-8";
				`
		_, err = db.Exec(ddl);
		if err != nil {
			db.Close()
			return &DCDB{}, err
		}
	case "postgresql":
		db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s", configIni["db_user"], configIni["db_password"], configIni["db_name"]))
		//fmt.Println(db)
		//fmt.Println(err)
		if err != nil {
			return &DCDB{}, err
		}
	case "mysql":
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", configIni["db_user"], configIni["db_password"], configIni["db_name"]))
		//fmt.Println("db",db)
		//fmt.Println(err)
		if err != nil {
			return &DCDB{}, err
		}
	}

	return &DCDB{db, configIni}, err
}
/*
func (db *DCDB) DbConnect() {
	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbInfo)
	CheckErr(err)
	err = db.Ping()
	//if err != nil {
	//	panic(err.Error()) // proper error handling instead of panic in your app
	//}
	//return db
}*/

func (db *DCDB) GetMainLockName() (string, error) {
	name, err := db.Single("SELECT script_name FROM main_lock")
	if err != nil {
		return "", err
	}
	return name, nil
}

func (db *DCDB) GetAllTables() ([]string, error) {
	var result []string
	var sql string
	switch db.configIni["db_type"] {
	case "sqlite" :
		sql = "SELECT name FROM sqlite_master WHERE type IN ('table','view') AND name NOT LIKE 'sqlite_%'"
	case "postgresql" :
		sql = "SELECT table_name FROM information_schema.tables WHERE table_type = 'BASE TABLE' AND    table_schema NOT IN ('pg_catalog', 'information_schema')"
	case "mysql" :
		sql = "SHOW TABLES"
	}
	result, err := db.GetList(sql)
	if err != nil {
		return result, err
	}
	return result, nil
}

func (db *DCDB) Single(query string, args ...interface{}) (string, error) {
	var result []byte
	err := db.QueryRow(query, args...).Scan(&result)
	switch {
	case err == sql.ErrNoRows:
		return "", nil
	case err != nil:
		return "", fmt.Errorf("%s in query %s %s", err, query, args)
	}
	if db.configIni["log"]=="1" {
		log.Printf("SQL: %s / %v", query, args)
	}
	return string(result), nil
}

func (db *DCDB) GetList(query string, args ...interface{}) ([]string, error) {
	var result []string
	all, err := db.GetAll(query, -1, args...)
	if err != nil {
		return result, err
	}
	for _, v := range all {
		for _, v2 := range v {
			result = append(result, v2)
		}
	}
	return result, nil
}

func (db *DCDB) GetAll(query string, countRows int, args ...interface{}) (map[int]map[string]string, error) {

	result := make(map[int]map[string]string)
	// Execute the query
	//fmt.Println("query", query)
	rows, err := db.Query(query, args...)
	if err != nil {
		return result, fmt.Errorf("%s in query %s %s", err, query, args)
	}

	if db.configIni["log"]=="1" {
		log.Printf("SQL: %s / %v", query, args)
	}
	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return result, fmt.Errorf("%s in query %s %s", err, query, args)
	}
	//fmt.Println("columns", columns)

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	r := 0
	// Fetch rows
	for rows.Next() {
		result[r] = make(map[string]string)

		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			return result, fmt.Errorf("%s in query %s %s", err, query, args)
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//fmt.Println(columns[i], ": ", value)
			result[r][columns[i]] = value
		}
		r++
		if countRows!=-1 && r >= countRows {
			break
		}
	}
	if err = rows.Err(); err != nil {
		return result, fmt.Errorf("%s in query %s %s", err, query, args)
	}
	//fmt.Println(result)
	return result, nil
}

func (db *DCDB) OneRow(query string, args ...interface{}) (map[string]string, error) {
	result := make(map[string]string)
	all, err := db.GetAll(query, 1, args ...)
	if err != nil {
		return result, fmt.Errorf("%s in query %s %s", err, query, args)
	}
	return all[0], nil
}

func (db *DCDB) ExecSql(query string, args ...interface{}) (int64, error) {
	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("%s in query %s %s", err, query, args)
	}
	affect, err := res.RowsAffected()
	lastId, err := res.LastInsertId()
	if db.configIni["log"]=="1" {
		log.Printf("SQL: %s / RowsAffected=%d / LastInsertId=%d / %s", query, affect, lastId, args)
	}
	return affect, nil
}


// для юнит-тестов. снимок всех данных в БД
func (db *DCDB) HashTableData(table, where, orderBy string) (string, error) {
	/*var columns string;
	rows, err := db.Query("select column_name from information_schema.columns where table_name= $1", table)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return "", err
		}
		columns+=name+"+"
	}
	columns = columns[:len(columns)-1]

	if len(columns) > 0 {
		if len(orderBy) > 0 {
			orderBy = " ORDER BY "+orderBy;
		}
	}*/
	if len(orderBy) > 0 {
		orderBy = " ORDER BY "+orderBy;
	}
	// это у всех разное, а значит и хэши будут разные, а это будет вызывать путаницу
	q:="SELECT md5(CAST((array_agg(t.* "+orderBy+")) AS text)) FROM \""+table+"\" t "+where
	/*if strings.Count(table, "my_table")>0 {
		columns = strings.Replace(columns,",notification","",-1)
		columns = strings.Replace(columns,"notification,","",-1)
		q="SELECT md5(CAST((array_agg("+columns+" "+orderBy+")) AS text)) FROM \""+table+"\" "+where
	}*/
	/*if strings.Count(columns, "cron_checked_time")>0 {
		columns = strings.Replace(columns, ",cron_checked_time", "", -1)
		columns = strings.Replace(columns, "cron_checked_time,", "", -1)
		q="SELECT md5(CAST((array_agg("+columns+" "+orderBy+")) AS text)) FROM \""+table+"\" "+where
	}*/
	hash, err := db.Single(q)
	if err != nil {
		return "", ErrInfo(err, q)
	}
	return hash, nil
}

// для юнит-тестов. снимок всех данных в БД
func (db *DCDB) AllHashes() (map[string]string, error) {
	//var orderBy string
	result:=make(map[string]string)
	//var columns string;
	rows, err := db.Query(`
		SELECT table_name
		FROM
		information_schema.tables
		WHERE
		table_type = 'BASE TABLE'
		AND
		table_schema NOT IN ('pg_catalog', 'information_schema');`)
	if err != nil {
		//fmt.Println(err)
		return result, err
	}
	for rows.Next() {
		var table string
		err = rows.Scan(&table)
		if err != nil {
			return result, err
		}
		//fmt.Println(table)

		orderByFns := func(table string) string {
			// ошибки не проверяются т.к. некритичны
			match, _ := regexp.MatchString("^(log_forex_orders|log_forex_orders_main|cf_comments|cf_currency|cf_funding|cf_lang|cf_projects|cf_projects_data)$", table)
			if match {
				return "id"
			}
			match, _ = regexp.MatchString("^log_time_(.*)$", table)
			if match && table!="log_time_money_orders" {
				return "user_id, time"
			}
			match, _ = regexp.MatchString("^log_transactions$", table)
			if match {
				return "time"
			}
			match, _ = regexp.MatchString("^log_votes$", table)
			if match {
				return "user_id, voting_id"
			}
			match, _ = regexp.MatchString("^log_(.*)$", table)
			if match && table!="log_time_money_orders" && table!="log_minute" {
				return "log_id"
			}
			match, _ = regexp.MatchString("^wallets$", table)
			if match {
				return "last_update"
			}
			return ""
		}
		orderBy := orderByFns(table)
		hash, err := db.HashTableData(table, "", orderBy)
		if err != nil {
			return result, ErrInfo(err)
		}
		result[table] = hash
	}
	return result, nil
}

func (db *DCDB) GetLastBlockData() (map[string]int64, error) {
	result := make(map[string]int64)
	confirmedBlockId, err := db.GetConfirmedBlockId()
	if err != nil {
		return result, ErrInfo(err)
	}
	if confirmedBlockId == 0 {
		confirmedBlockId = 1
	}
	log.Print("confirmedBlockId", confirmedBlockId)
	// получим время из последнего подвержденного блока
	lastBlockBin, err := db.Single("SELECT data FROM block_chain WHERE id =?", confirmedBlockId)
	if err != nil || len(lastBlockBin)==0 {
		return result, ErrInfo(err)
	}
	// ID блока
	result["blockId"] = int64(BinToDec([]byte(lastBlockBin[1:5])))
	// Время последнего блока
	result["lastBlockTime"] = int64(BinToDec([]byte(lastBlockBin[5:9])))
	return result, nil
}

func (db *DCDB) GetMyNoticeData(sessRestricted int, sessUserId int64, myPrefix string, lang map[string]string) (map[string]string, error) {
	result := make(map[string]string)
	if sessRestricted == 0 {
		my_table, err := db.OneRow("SELECT user_id, miner_id, status FROM "+myPrefix+"my_table")
		if err != nil {
			return result, ErrInfo(err)
		}
		if my_table["user_id"] == "0" {
			result["account_status"] = "searching"
		} else if my_table["status"] == "bad_key" {
			result["account_status"] = "bad_key"
		} else if my_table["miner_id"] != "0" {
			result["account_status"] = "miner"
		} else if my_table["user_id"] != "0" {
			result["account_status"] = "user"
		}
	} else {
		// user_id уже есть, т.к. мы смогли зайти в урезанном режиме по паблик-кею
		// проверим, может есть что-то в miners_data
		status, err := db.Single("SELECT status FROM miners_data WHERE user_id = $1", sessUserId)
		if err != nil {
			return result, ErrInfo(err)
		}
		if len(status) > 0 {
			result["account_status"] = status
		} else {
			result["account_status"] = "user"
		}
	}
	result["account_status"] = lang["status_"+result["account_status"]]

	// Инфа о последнем блоке
	blockData, err := db.GetLastBlockData()
	if err != nil {
		return result, ErrInfo(err)
	}
	result["cur_block_id"] = Int64ToStr(blockData["blockId"])
	t := time.Unix(blockData["lastBlockTime"], 0)
	result["time_last_block"] = t.Format("2006-01-02 15:04:05")
	result["time_last_block_int"] = Int64ToStr(blockData["lastBlockTime"])

	result["connections"], err = db.Single("SELECT count(*) FROM nodes_connection")
	if err != nil {
		return result, ErrInfo(err)
	}

	if time.Now().Unix() - blockData["lastBlockTime"] > 1800 {
		result["main_status"] = lang["downloading_blocks"]
		result["main_status_complete"] = "0"
	} else {
		result["main_status"] = lang["downloading_complete"]
		result["main_status_complete"] = "1"
	}

	return result, nil
}

func (db *DCDB) GetPoolAdminUserId() (int64, error)  {
	result, err := db.Single("SELECT pool_admin_user_id FROM config")
	if err != nil {
		return 0, ErrInfo(err)
	}
	return StrToInt64(result), nil
}

func (db *DCDB) GetMyPublicKey(myPrefix string) (string, error) {
	result, err := db.Single("SELECT public_key FROM "+myPrefix+"my_keys WHERE block_id = (SELECT max(block_id) FROM "+myPrefix+"my_keys)")
	if err != nil {
		return "", ErrInfo(err)
	}
	return result, nil
}

func (db *DCDB) GetDataAuthorization(hash string) (string, error) {
	// получим данные для подписи
	var sql string
	switch db.configIni["db_type"] {
	case "sqlite":
		sql = `SELECT data FROM authorization WHERE hash = $1`
	case "postgresql":
		sql = `select data from "authorization" where "hash" = '\x$1'`
	case "mysql":
		sql = `SELECT data FROM authorization WHERE hash = 0x$1`
	}
	data, err := db.Single(sql, hash)
	if err != nil {
		return "", ErrInfo(err)
	}
	return data, nil
}

func (db *DCDB) GetAdminUserId() (int64, error) {
	result, err := db.Single("SELECT user_id FROM admin")
	if err != nil {
		return 0, ErrInfo(err)
	}
	return StrToInt64(result), nil
}

func (db *DCDB) GetUserPublicKey(userId int64) (string, error) {
	result, err := db.Single("SELECT public_key_0 FROM users WHERE user_id = $1", userId)
	if err != nil {
		return "", ErrInfo(err)
	}
	return result, nil
}

func (db *DCDB) GetNodePrivateKey(myPrefix string) string {
	var key string;
	rows, err := db.Query("SELECT private_key FROM "+myPrefix+"my_node_keys WHERE block_id = (SELECT max(block_id) FROM "+myPrefix+"my_node_keys)")
	CheckErr(err)
	if  ok := rows.Next(); ok {
		err = rows.Scan(&key)
		CheckErr(err)
	}
	return key
}

func (db *DCDB) GetNodeConfig() (map[string]string, error) {
	var result map[string]string
	result, err := db.OneRow("SELECT * FROM config")
	if err != nil{
		return result, ErrInfo(err)
	}
	return result, nil
}

func (db *DCDB) TestBlock () (*prevBlockType, int64, int64, int64, int64, [][][]int64, error) {

	var minerId, userId, level, i, currentMinerId, currentUserId int64;
	prevBlock := new(prevBlockType)
	var levelsRange [][][]int64
	// последний успешно записанный блок
	rows, err := db.Query(`
            SELECT LOWER(encode(hash, 'hex')),
            LOWER(encode(head_hash, 'hex')),
            block_id,
            time,
            level
            FROM info_block
            `)

	defer rows.Close()

	if  ok := rows.Next(); ok {
		err = rows.Scan(&prevBlock.Hash, &prevBlock.HeadHash, &prevBlock.BlockId, &prevBlock.Time, &prevBlock.Level)
		if err!= nil {
			return prevBlock, userId, minerId, currentUserId, level, levelsRange, err
		}
	}
	//fmt.Println("prevBlock", prevBlock)

	// общее кол-во майнеров
	row, err := db.Single("SELECT max(miner_id) FROM miners")
	if err != nil {
		return prevBlock, userId, minerId, currentUserId, level, levelsRange, err
	}
	maxMinerId := StrToInt64(row)

	for currentUserId == 0 {
		// если майнера заморозили то у него исчезает miner_id, чтобы не попасть на такой пустой miner_id
		// нужно пербирать энтропию, пока не дойдем до существующего miner_id
		var entropy int64
		if (i == 0) {
			entropy = GetEntropy(prevBlock.HeadHash);
		} else {
			time.Sleep(1000 * time.Millisecond)

			blockId := prevBlock.BlockId - i;
			if (blockId < 1) {
				break;
			}

			rows, err = db.Query("SELECT LOWER(encode(head_hash, 'hex'))   FROM block_chain  WHERE id = " + strconv.FormatInt(blockId, 10))
			defer rows.Close()
			CheckErr(err)
			var newHeadHash string
			if  ok := rows.Next(); ok {
				err = rows.Scan(&newHeadHash)
				CheckErr(err)
			}
			//fmt.Println("newHeadHash", newHeadHash)
			entropy = GetEntropy(newHeadHash);
		}
		currentMinerId = GetBlockGeneratorMinerId(maxMinerId, entropy);

		// получим ID юзера по его miner_id
		row, err = db.Single("SELECT user_id  FROM miners_data  WHERE miner_id = " + strconv.FormatInt(currentMinerId, 10))
		if err != nil {
			return prevBlock, userId, minerId, currentUserId, level, levelsRange, err
		}
		currentUserId = StrToInt64(row)
		i++;
	}

	collective, err := db.GetMyUsersIds(true)
	if err != nil {
		return prevBlock, userId, minerId, currentUserId, level, levelsRange, err
	}

	// в сингл-моде будет только $my_miners_ids[0]
	myMinersIds, err := db.GetMyMinersIds(collective);
	if err != nil {
		return prevBlock, userId, minerId, currentUserId, level, levelsRange, err
	}

	// есть ли кто-то из нашего пула (или сингл-мода), кто находится на 0-м уровне
	if InSliceInt64(currentMinerId, myMinersIds) {
		level = 0;
		levelsRange = append(levelsRange, [][]int64 {{1,1}});
		minerId = currentMinerId;
	} else {
		levelsRange = GetBlockGeneratorMinerIdRange (currentMinerId, maxMinerId);
		//fmt.Println("levelsRange", levelsRange)
		if len(myMinersIds)>0 {
			minerId, level = FindMinerIdLevel(myMinersIds,levelsRange);
		} else {
			level = -1; // у нас нет уровня, т.к. пуст $my_miners_ids, т.е. на сервере нет майнеров
			minerId = 0;
		}
	}
	err = db.QueryRow("SELECT user_id FROM miners_data WHERE miner_id = $1", 1).Scan(&userId)
	if err != nil {
		return prevBlock, userId, minerId, currentUserId, level, levelsRange, err
	}
	return prevBlock, userId, minerId, currentUserId, level, levelsRange, nil
}

func  (db *DCDB) GetSleepData() map[string][]int64 {
	var sleepDataJson []byte
	err := db.QueryRow("SELECT value FROM variables WHERE name = 'sleep'").Scan(&sleepDataJson)
	CheckErr(err)
	var sleepDataMap map[string][]int64
	err = json.Unmarshal(sleepDataJson, &sleepDataMap)
	CheckErr(err)
	return sleepDataMap;
}


func  (db *DCDB) GetMyMinersIds(collective []int64) ([]int64, error) {
	var miners []int64
	rows, err := db.Query("SELECT miner_id FROM miners_data WHERE user_id IN ($1) AND miner_id > 0", strings.Join(SliceInt64ToString(collective), ","))
	if err != nil {
		return miners, err
	}
	defer rows.Close()
	for rows.Next() {
		var minerId int64
		err = rows.Scan(&minerId)
		if err != nil {
			return miners, err
		}
		miners = append(miners, minerId);
	}
	return miners, nil;
}

func (db *DCDB) GetConfirmedBlockId() (int64, error) {
	localGateIp, err := db.Single("SELECT local_gate_ip FROM config")
	if err != nil {
		return 0, err
	}
	if localGateIp!="" {
		blockId, err := db.GetBlockId()
		if err != nil {
			return 0, err
		}
		return blockId, nil
	} else {
		result, err := db.Single("SELECT max(block_id) FROM confirmations WHERE good >= ?", consts.MIN_CONFIRMED_NODES)
		if err != nil {
			return 0, err
		}
		//log.Print("result int64",StrToInt64(result))
		return StrToInt64(result), nil
	}
}
/*function get_confirmed_block_id($db)
{
	// в защищенном режиме нет прямого выхода в интернет, поэтому просто берем get_block_id
	$config['local_gate_ip'] = $db->query( __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__, "
			SELECT `local_gate_ip`
			FROM `".DB_PREFIX."config`
			", 'fetch_one');
	if ($config['local_gate_ip'])
		return get_block_id($db);
	else
		return $db->query( __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__, "
				SELECT max(`block_id`)
				FROM `".DB_PREFIX."confirmations`
				WHERE `good` >= ".MIN_CONFIRMED_NODES."
				", 'fetch_one');
}*/

func (db *DCDB)  GetCommunityUsers() ([]int64, error) {
	var users []int64
	rows, err := db.Query("SELECT user_id FROM community")
	if err != nil {
		return users, ErrInfo(err)
	}
	defer rows.Close()
	for rows.Next() {
		var userId int64
		err = rows.Scan(&userId)
		if err != nil{
			return users, ErrInfo(err)
		}
		users = append(users, userId);
	}
	return users, err;
}

func (db *DCDB) GetMyUsersId(myPrefix string) (int64, error) {
	userId, err := db.Single("SELECT user_id FROM "+myPrefix+"my_table")
	if err != nil {
		return 0, err
	}
	return StrToInt64(userId), nil
}

func (db *DCDB) GetMyUsersIds(checkCommission bool) ([]int64, error) {
	var usersIds []int64
	usersIds, err := db.GetCommunityUsers();
	if err != nil {
		return usersIds, err
	}
	if len(usersIds) == 0 { // сингл-мод
		rows, err := db.Query("SELECT user_id FROM my_table")
		if err != nil {
			return usersIds, err
		}
		defer rows.Close()
		if  ok := rows.Next(); ok {
			var x int64;
			err = rows.Scan(&x)
			if err != nil {
				return usersIds, err
			}
			usersIds = append(usersIds, x)
		}
	} else{
		// нельзя допустить, чтобы блок подписал майнер, у которого комиссия больше той, что разрешана в пуле,
		// т.к. это приведет к попаднию в блок некорректной тр-ии, что приведет к сбою пула
		if checkCommission {
			// комиссия на пуле
			rows, err := db.Query("SELECT commission FROM config")
			if err != nil {
				return usersIds, err
			}
			defer rows.Close()
			if  ok := rows.Next(); ok {
				var commissionJson []byte;
				err = rows.Scan(&commissionJson)
				if err != nil {
					return usersIds, err
				}

				var commissionPoolMap map[string][]float64
				err := json.Unmarshal(commissionJson, &commissionPoolMap)
				if err != nil {
					return usersIds, err
				}

				rows, err := db.Query("SELECT user_id, commission FROM commission WHERE user_id IN ("+strings.Join(SliceInt64ToString(usersIds), ",")+")")
				if err != nil {
					return usersIds, err
				}
				defer rows.Close()
				if  ok := rows.Next(); ok {
					var uid int64;
					var commJson []byte;
					err = rows.Scan(&uid, &commJson)
					if err != nil {
						return usersIds, err
					}

					var commissionUserMap map[string][]float64
					err := json.Unmarshal(commJson, &commissionUserMap)
					if err != nil {
						return usersIds, err
					}

					for currencyId, Commissions := range commissionUserMap {

						if Commissions[0] > commissionPoolMap[currencyId][0] || Commissions[1] > commissionPoolMap[currencyId][1] {
							//fmt.Println("del_user_id_from_array")
							DelUserIdFromArray(&usersIds, uid);
						}
					}
				}

			}
		}
	}
	return usersIds, nil;
}

func (db *DCDB) GetBlockId() (int64, error) {
	blockId, err := db.Single("SELECT block_id FROM info_block")
	if err != nil {
		return 0, err
	}
	return StrToInt64(blockId), nil
}

func (db *DCDB) GetUserIdByPublicKey(publicKey string) (string, error) {
	var sql string
	switch db.configIni["db_type"] {
	case "sqlite":
		sql = `SELECT user_id FROM users WHERE public_key_0 = $1`
	case "postgresql":
		sql = `SELECT user_id FROM users WHERE public_key_0 = '\x$1'`
	case "mysql":
		sql = `SELECT user_id FROM users WHERE public_key_0 = 0x$1`
	}
	userId, err := db.Single(sql, publicKey)
	if err != nil{
		return "", ErrInfo(err)
	}
	return userId, nil
}

func (db *DCDB) InsertIntoMyKey(userId, publicKey, curBlockId string) error {
	var sql string
	switch db.configIni["db_type"] {
	case "sqlite":
		sql = `INSERT INTO `+userId+`_my_keys (public_key, status, block_id) VALUES ($1,'approved', $2)`
	case "postgresql":
		sql = `INSERT INTO `+userId+`_my_keys (public_key, status, block_id) VALUES ('\x$1','approved', $2)`
	case "mysql":
		sql = `INSERT INTO `+userId+`_my_keys (public_key, status, block_id) VALUES (0x$1,'approved', $2)`
	}
	_, err := db.ExecSql(sql, publicKey, curBlockId )
	if err != nil {
		return err
	}
	return nil
}

func (db *DCDB) GetInfoBlock() (map[string]string, error) {
	var result map[string]string
	result, err := db.OneRow("SELECT * FROM info_block")
	if err != nil{
		return result, ErrInfo(err)
	}
	if len(result)==0 {
		return result, fmt.Errorf("empty info_block")
	}
	return result, nil
}

func (db *DCDB) GetTestBlockId() (int64, error) {
	rows, err := db.Query("SELECT block_id FROM testblock")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	if  ok := rows.Next(); ok {
		var block_id int64
		err = rows.Scan(&block_id)
		if err != nil {
			return 0, err
		}
		return block_id, nil
	}
	return 0, nil
}

func (db *DCDB) GetMyPrefix() (string, error) {
	collective, err := db.GetCommunityUsers()
	if err != nil {
		return "", ErrInfo(err)
	}
	if len(collective) == 0 {
		return "", nil
	} else {
		myUserId, err := db.GetPoolAdminUserId()
		if err != nil || myUserId == 0  {
			if err != nil {
				return "", ErrInfo(err)
			} else {
				return "", fmt.Errorf("myUserId==0")
			}
		}
		return Int64ToStr(myUserId)+"_", nil
	}
}


func (db *DCDB) GetMyLocalGateIp() (string, error) {
	result, err := db.Single("SELECT local_gate_ip FROM config")
	if err != nil {
		return "", err
	}
	return result, nil
}

func (db *DCDB) DbLock() {
	var affect int64;
	for affect==0 {
		t := time.Now().Unix();
		stmt, err := db.Prepare(`INSERT INTO main_lock(lock_time,script_name)
                                                    VALUES($1,$2)`)
		CheckErr(err)
		defer stmt.Close()

		res, err := stmt.Exec(t, "testblock_generator")
		if fmt.Sprintf("%s", err)=="pq: duplicate key value violates unique constraint \"main_lock_pkey\"" {
			//fmt.Println(err)
		} else {
			CheckErr(err)
			affect, err = res.RowsAffected()
			CheckErr(err)
			//fmt.Println(affect, "rows changed")
		}
		defer stmt.Close()

		if affect == 0 {
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func (db *DCDB) DbUnlock() {
	rows, err := db.Query("DELETE FROM main_lock WHERE script_name='testblock_generator'")
	defer rows.Close()
	CheckErr(err)
	defer rows.Close()
}

func (db *DCDB) GetIsReadySleep(level int64) int64 {
	SleepData := db.GetSleepData();
	return GetIsReadySleep0(level, SleepData["is_ready"])
}

func (db *DCDB) GetGenSleep(prevBlock *prevBlockType, level int64) int64 {

	sleepData := db.GetSleepData()

	// узнаем время, которые было затрачено в ожидании is_ready предыдущим блоком
	isReadySleep := db.GetIsReadySleep(prevBlock.Level)
	//fmt.Println("isReadySleep", isReadySleep)

	// сколько сек должен ждать нод, перед тем, как начать генерить блок, если нашел себя в одном из уровней.
	generatorSleep := GetGeneratorSleep(level, sleepData["generator"])
	//fmt.Println("generatorSleep", generatorSleep)

	// сумма is_ready всех предыдущих уровней, которые не успели сгенерить блок
	isReadySleep2 := GetIsReadySleepSum(level , sleepData["is_ready"])
	//fmt.Println("isReadySleep2", isReadySleep2)

	// узнаем, сколько нам нужно спать
	sleep := isReadySleep + generatorSleep + isReadySleep2;
	return sleep;
}
