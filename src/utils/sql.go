package utils
import (
	 "fmt"
	 _ "github.com/lib/pq"
     _ "github.com/go-sql-driver/mysql"
	 "database/sql"
	"strings"
	//"regexp"
	//"errors"
	"log"
	"time"
	"strconv"
	"encoding/json"
	//"github.com/astaxie/beego/config"
	"consts"
	"sync"
//	"sync/atomic"
//	"os"
	//"dcparser"
)

type DCDB struct {
	 *sql.DB
	ConfigIni map[string]string
}

func replQ(q string) string {
	q1:=strings.Split(q, "?")
	fmt.Println(q1)
	result:=""
	for i:=0; i < len(q1); i++ {
		if i != len(q1)-1 {
			result+=q1[i]+"$"+IntToStr(i+1)
		} else {
			result+=q1[i]
		}
	}
	return result
}

func NewDbConnect(ConfigIni map[string]string) (*DCDB, error) {
	var db *sql.DB
	var err error
	switch ConfigIni["db_type"] {
	case "sqlite":

		fmt.Println("sqlite")
		db, err = sql.Open("sqlite3", "litedb.db")
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
		db, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s", ConfigIni["db_user"], ConfigIni["db_password"], ConfigIni["db_name"]))
		//fmt.Println(db)
		//fmt.Println(err)
		if err != nil {
			return &DCDB{}, err
		}
	case "mysql":
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", ConfigIni["db_user"], ConfigIni["db_password"], ConfigIni["db_name"]))
		//fmt.Println("db",db)
		//fmt.Println(err)
		if err != nil {
			return &DCDB{}, err
		}
	}

	return &DCDB{db, ConfigIni}, err
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

func (db *DCDB) GetConfigIni(name string) string {
	return db.ConfigIni[name]
}

func (db *DCDB) GetMainLockName() (string, error) {
	name, err := db.Single("SELECT script_name FROM main_lock").String()
	if err != nil {
		return "", err
	}
	return name, nil
}

func (db *DCDB) GetAllTables() ([]string, error) {
	var result []string
	var sql string
	switch db.ConfigIni["db_type"] {
	case "sqlite" :
		sql = "SELECT name FROM sqlite_master WHERE type IN ('table','view') AND name NOT LIKE 'sqlite_%'"
	case "postgresql" :
		sql = "SELECT table_name FROM information_schema.tables WHERE table_type = 'BASE TABLE' AND    table_schema NOT IN ('pg_catalog', 'information_schema')"
	case "mysql" :
		sql = "SHOW TABLES"
	}
	result, err := db.GetList(sql).String()
	if err != nil {
		return result, err
	}
	return result, nil
}

type Variables struct {
	Int64 map[string]int64
	String map[string]string
	Float64 map[string]float64
}
func (db *DCDB) GetAllVariables() (*Variables, error) {
	result := new(Variables)
	result.Int64 = make(map[string]int64)
	result.String = make(map[string]string)
	result.Float64 = make(map[string]float64)
	all, err := db.GetAll("SELECT * FROM variables", -1)
	fmt.Println(all)
	if err != nil {
		return result, err
	}
	for _, v := range all {
		switch v["name"] {
		case "alert_error_time", "error_time", "promised_amount_points", "promised_amount_votes_0", "promised_amount_votes_1", "promised_amount_votes_period", "holidays_max", "limit_abuses", "limit_abuses_period", "limit_promised_amount", "limit_promised_amount_period", "limit_cash_requests_out", "limit_cash_requests_out_period", "limit_change_geolocation", "limit_change_geolocation_period", "limit_holidays", "limit_holidays_period", "limit_message_to_admin", "limit_message_to_admin_period", "limit_mining", "limit_mining_period", "limit_node_key", "limit_node_key_period", "limit_primary_key", "limit_primary_key_period", "limit_votes_miners", "limit_votes_miners_period", "limit_votes_complex", "limit_votes_complex_period", "limit_commission", "limit_commission_period", "limit_new_miner", "limit_new_miner_period", "limit_new_user", "limit_new_user_period", "max_block_size", "max_block_user_transactions", "max_day_votes", "max_tx_count", "max_tx_size", "max_user_transactions", "miners_keepers", "miner_points", "miner_votes_0", "miner_votes_1", "miner_votes_attempt", "miner_votes_period", "mining_votes_0", "mining_votes_1", "mining_votes_period", "min_miners_keepers", "node_voting", "node_voting_period", "rollback_blocks_1", "rollback_blocks_2", "limit_change_host", "limit_change_host_period", "min_miners_of_voting", "min_hold_time_promise_amount", "min_promised_amount", "points_update_time", "reduction_period", "new_pct_period", "new_max_promised_amount", "new_max_other_currencies", "cash_request_time", "limit_for_repaid_fix", "limit_for_repaid_fix_period", "miner_newbie_time", "system_commission":
			result.Int64[v["name"]] = StrToInt64(v["value"])
		case "points_factor" :
			result.Float64[v["name"]] = StrToFloat64(v["value"])
		case "sleep" :
			result.String[v["name"]] = v["value"]
		}
	}
	return result, err
}
/*
func (db *DCDB) SingleInt64(query string, args ...interface{}) (int64, error) {
	result, err := db.Single(query, args...)
	if err != nil {
		return 0, err
	}
	return StrToInt64(result), nil
}
*/

type singleResult struct {
	result []byte
	err error
}

type listResult struct {
	result []string
	err error
}

type oneRow struct {
	result map[string]string
	err error
}

func (r *listResult) Int64() ([]int64, error) {
	var result []int64
	if r.err != nil {
		return result, r.err
	}
	for _, v := range r.result {
		result = append(result, StrToInt64(v))
	}
	return result, nil
}

func (r *listResult) MapInt() (map[int]int, error) {
	result := make(map[int]int)
	if r.err != nil {
		return result, r.err
	}
	i := 0
	for _, v := range r.result {
		result[i] = StrToInt(v)
		i++
	}
	return result, nil
}

func (r *listResult) String() ([]string, error) {
	if r.err != nil {
		return r.result, r.err
	}
	return r.result, nil
}


func (r *oneRow) String() (map[string]string, error) {
	if r.err != nil {
		return r.result, r.err
	}
	return r.result, nil
}

func (r *oneRow) Int64() (map[string]int64, error) {
	result := make(map[string]int64)
	if r.err != nil {
		return result, r.err
	}
	for k, v := range r.result {
		result[k] = StrToInt64(v)
	}
	return result, nil
}

func (r *oneRow) Int() (map[string]int, error) {
	result := make(map[string]int)
	if r.err != nil {
		return result, r.err
	}
	for k, v := range r.result {
		result[k] = StrToInt(v)
	}
	return result, nil
}

func (r *singleResult) Int64() (int64, error) {
	if r.err!=nil {
		return 0, r.err
	}
	return BytesToInt64(r.result), nil
}
func (r *singleResult) Int() (int, error) {
	if r.err!=nil {
		return 0, r.err
	}
	return BytesToInt(r.result), nil
}

func (r *singleResult) Float64() (float64, error) {
	if r.err!=nil {
		return 0, r.err
	}
	return StrToFloat64(string(r.result)), nil
}

func (r *singleResult) String() (string, error) {
	if r.err!=nil {
		return "", r.err
	}
	return string(r.result), nil
}

func (r *singleResult) Bytes() ([]byte, error) {
	if r.err!=nil {
		return []byte(""), r.err
	}
	return r.result, nil
}


func (db *DCDB) Single(query string, args ...interface{}) *singleResult {
	switch db.ConfigIni["db_type"] {
	case "sqlite":
		query = strings.Replace(query, "[hex]", "?", -1)
	case "postgresql":
		query = strings.Replace(query, "[hex]", "decode(?,'HEX')", -1)
		query = replQ(query)
	case "mysql":
		query = strings.Replace(query, "[hex]", "UNHEX(?)", -1)
	}
	var result []byte
	err := db.QueryRow(query, args...).Scan(&result)
	if err!=nil {
		fmt.Println("[error] single", fmt.Sprintf("%s in query %s %s\n", err, query, args))
	}
	switch {
	case err == sql.ErrNoRows:
		return &singleResult{[]byte(""), nil}
	case err != nil:
		return  &singleResult{[]byte(""), fmt.Errorf("%s in query %s %s", err, query, args)}
	}
	if db.ConfigIni["log"]=="1" {
		log.Printf("SQL: %s / %v", query, args)
	}
	return &singleResult{result, nil}
}

func (db *DCDB) GetMap(query string, name, value string, args ...interface{}) (map[string]string, error) {
	result := make(map[string]string)
	all, err := db.GetAll(query, -1, args ...)
	fmt.Println(all)
	if err != nil {
		return result, err
	}
	for _, v := range all {
		result[v[name]] = v[value]
	}
	return result, err
}

func (db *DCDB) GetList(query string, args ...interface{}) *listResult {
	var result []string
	all, err := db.GetAll(query, -1, args...)
	if err != nil {
		return &listResult{result, err}
	}
	for _, v := range all {
		for _, v2 := range v {
			result = append(result, v2)
		}
	}
	return &listResult{result, nil}
}

func (db *DCDB) GetAll(query string, countRows int, args ...interface{}) ([]map[string]string, error) {
	if db.ConfigIni["db_type"] == "postgresql" {
		query = replQ(query)
	}
	var result []map[string]string
	// Execute the query
	//fmt.Println("query", query)
	rows, err := db.Query(query, args...)
	if err != nil {
		return result, fmt.Errorf("%s in query %s %s", err, query, args)
	}
	defer rows.Close()

	if db.ConfigIni["log"]=="1" {
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
		//result[r] = make(map[string]string)

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
			result = append(result, map[string]string{columns[i]:value})
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

func (db *DCDB) OneRow(query string, args ...interface{}) *oneRow {
	result := make(map[string]string)
	all, err := db.GetAll(query, 1, args ...)
	if err != nil {
		return &oneRow{result, fmt.Errorf("%s in query %s %s", err, query, args)}
	}
	if len(all) == 0 {
		return &oneRow{result, nil}
	}
	return &oneRow{all[0], nil}
}

func (db *DCDB) InsertInLogTx(binaryTx []byte, time int64) error {
	txMD5 := Md5(binaryTx)
	err := db.ExecSql("INSERT INTO log_transactions (hash, time) VALUES ([hex], ?)", txMD5, time)
	if err != nil {
		return ErrInfo(err)
	}
	return nil
}

func (db *DCDB) DelLogTx(binaryTx []byte) error {
	txMD5 := Md5(binaryTx)
	err := db.ExecSql("DELETE FROM log_transactions WHERE hash=[hex]", txMD5)
	if err != nil {
		return ErrInfo(err)
	}
	return nil
}

func formatQuery(query, dbType string) string {
	switch dbType {
	case "sqlite":
		query = strings.Replace(query, "[hex]", "?", -1)
	case "postgresql":
		query = strings.Replace(query, "[hex]", "decode(?,'HEX')", -1)
		query = replQ(query)
	case "mysql":
		query = strings.Replace(query, "[hex]", "UNHEX(?)", -1)
	}
	return query
}

func (db *DCDB) ExecSqlGetLastInsertId(query string, args ...interface{}) (int64, error) {
	query = formatQuery(query, db.ConfigIni["db_type"])
	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("%s in query %s %s", err, query, args)
	}
	affect, err := res.RowsAffected()
	lastId, err := res.LastInsertId()
	if db.ConfigIni["log"]=="1" {
		log.Printf("SQL: %s / RowsAffected=%d / LastInsertId=%d / %s", query, affect, lastId, args)
	}
	return lastId, nil
}

func (db *DCDB) ExecSql(query string, args ...interface{}) (error) {
	query = formatQuery(query, db.ConfigIni["db_type"])
	res, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("%s in query %s %s", err, query, args)
	}
	affect, err := res.RowsAffected()
	lastId, err := res.LastInsertId()
	if db.ConfigIni["log"]=="1" {
		log.Printf("SQL: %s / RowsAffected=%d / LastInsertId=%d / %s", query, affect, lastId, args)
	}
	return nil
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
	hash, err := db.Single(q).String()
	if err != nil {
		return "", ErrInfo(err, q)
	}
	return hash, nil
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
	lastBlockBin, err := db.Single("SELECT data FROM block_chain WHERE id =?", confirmedBlockId).Bytes()
	if err != nil || len(lastBlockBin)==0 {
		return result, ErrInfo(err)
	}
	// ID блока
	result["blockId"] = int64(BinToDec(lastBlockBin[1:5]))
	// Время последнего блока
	result["lastBlockTime"] = int64(BinToDec(lastBlockBin[5:9]))
	return result, nil
}

func (db *DCDB) GetMyNoticeData(sessRestricted int, sessUserId int64, myPrefix string, lang map[string]string) (map[string]string, error) {
	result := make(map[string]string)
	if sessRestricted == 0 {
		my_table, err := db.OneRow("SELECT user_id, miner_id, status FROM "+myPrefix+"my_table").String()
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
		status, err := db.Single("SELECT status FROM miners_data WHERE user_id = $1", sessUserId).String()
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

	result["connections"], err = db.Single("SELECT count(*) FROM nodes_connection").String()
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
	result, err := db.Single("SELECT pool_admin_user_id FROM config").Int64()
	if err != nil {
		return 0, ErrInfo(err)
	}
	return result, nil
}

func (db *DCDB) GetMyPublicKey(myPrefix string) ([]byte, error) {
	result, err := db.Single("SELECT public_key FROM "+myPrefix+"my_keys WHERE block_id = (SELECT max(block_id) FROM "+myPrefix+"my_keys)").Bytes()
	if err != nil {
		return []byte(""), ErrInfo(err)
	}
	return result, nil
}

func (db *DCDB) GetDataAuthorization(hash []byte) (string, error) {
	// получим данные для подписи
	var sql string
	switch db.ConfigIni["db_type"] {
	case "sqlite":
		sql = `SELECT data FROM authorization WHERE hash = $1`
	case "postgresql":
		sql = `select data from "authorization" where "hash" = '\x$1'`
	case "mysql":
		sql = `SELECT data FROM authorization WHERE hash = 0x$1`
	}
	data, err := db.Single(sql, hash).String()
	if err != nil {
		return "", ErrInfo(err)
	}
	return data, nil
}

func (db *DCDB) GetAdminUserId() (int64, error) {
	result, err := db.Single("SELECT user_id FROM admin").Int64()
	if err != nil {
		return 0, ErrInfo(err)
	}
	return result, nil
}

func (db *DCDB) GetUserPublicKey(userId int64) (string, error) {
	result, err := db.Single("SELECT public_key_0 FROM users WHERE user_id = $1", userId).String()
	if err != nil {
		return "", ErrInfo(err)
	}
	return result, nil
}

func (db *DCDB) GetNodePrivateKey(myPrefix string) (string, error) {
	var key string
	key, err := db.Single("SELECT private_key FROM "+myPrefix+"my_node_keys WHERE block_id = (SELECT max(block_id) FROM "+myPrefix+"my_node_keys)").String()
	if err != nil {
		return "", ErrInfo(err)
	}
	return key, nil
}

func (db *DCDB) GetMaxPromisedAmount(currencyId int64) (float64, error) {
	result, err := db.Single("SELECT amount FROM max_promised_amounts WHERE currencyId = ? ORDER BY time DESC", currencyId).Float64()
	if err != nil {
		return 0, err
	}
	return result, nil
}


func (db *DCDB) GetMaxPromisedAmounts() (map[int64][]map[int64]string, error) {
	var result map[int64][]map[int64]string
	rows, err := db.Query("SELECT currency_id, time, amount  FROM max_promised_amounts ORDER BY time ASC")
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var currency_id, time int64
		var amount string
		err = rows.Scan(&currency_id, &time, &amount)
		if err!= nil {
			return result, err
		}
		result[currency_id] = append(result[currency_id], map[int64]string{time:amount})
	}
	return result, nil
}




func (db *DCDB) GetPrivateKey(myPrefix string) (string, error) {
	var key string
	key, err := db.Single("SELECT private_key FROM "+myPrefix+"my_keys WHERE block_id = (SELECT max(block_id) FROM "+myPrefix+"my_keys)").String()
	if err != nil {
		return "", ErrInfo(err)
	}
	return key, nil
}

func (db *DCDB) GetNodeConfig() (map[string]string, error) {
	var result map[string]string
	result, err := db.OneRow("SELECT * FROM config").String()
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
	maxMinerId, err := db.Single("SELECT max(miner_id) FROM miners").Int64()
	if err != nil {
		return prevBlock, userId, minerId, currentUserId, level, levelsRange, err
	}

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
		currentUserId, err = db.Single("SELECT user_id  FROM miners_data  WHERE miner_id = " + strconv.FormatInt(currentMinerId, 10)).Int64()
		if err != nil {
			return prevBlock, userId, minerId, currentUserId, level, levelsRange, err
		}
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

func  (db *DCDB) GetSleepData() (map[string][]int64, error) {
	sleepDataMap := make(map[string][]int64)
	var sleepDataJson []byte
	sleepDataJson, err := db.Single("SELECT value FROM variables WHERE name = 'sleep'").Bytes()
	if err != nil {
		return sleepDataMap, ErrInfo(err)
	}
	if len(sleepDataJson) > 0 {
		err = json.Unmarshal(sleepDataJson, &sleepDataMap)
		if err != nil {
			return sleepDataMap, ErrInfo(err)
		}
	}
	return sleepDataMap, nil;
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
	localGateIp, err := db.Single("SELECT local_gate_ip FROM config").String()
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
		result, err := db.Single("SELECT max(block_id) FROM confirmations WHERE good >= ?", consts.MIN_CONFIRMED_NODES).Int64()
		if err != nil {
			return 0, err
		}
		//log.Print("result int64",StrToInt64(result))
		return result, nil
	}
}


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

func (db *DCDB) GetMyUserId(myPrefix string) (int64, error) {
	userId, err := db.Single("SELECT user_id FROM "+myPrefix+"my_table").Int64()
	if err != nil {
		return 0, err
	}
	return userId, nil
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
	blockId, err := db.Single("SELECT block_id FROM info_block").Int64()
	if err != nil {
		return 0, err
	}
	return blockId, nil
}
// наличие cash_requests с pending означает, что у юзера все обещанные суммы в for_repaid. Возможно, временно, если это свежий запрос и юзер еще не успел послать cash_requests_in
func (db *DCDB) CheckCashRequests(userId int64) (error) {
	cashRequestStatus, err := db.Single("SELECT status FROM cash_requests WHERE to_user_id  =  ? AND del_block_id  =  0 AND for_repaid_del_block_id  =  0 AND status  =  'pending'", userId).String()
	if err != nil {
		return err
	}
	if len(cashRequestStatus) > 0 {
		return fmt.Errorf("cashRequestStatus")
	}
	return nil
}

func(db *DCDB) GetPct() (map[int64][]map[int64]map[string]float64, error) {
	result := make(map[int64][]map[int64]map[string]float64)
	rows, err := db.Query("SELECT currency_id, time, user, miner FROM pct GROUP BY time ASC")
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var currency_id, time int64
		var user, miner float64
		err = rows.Scan(&currency_id, &time, &user, &miner)
		if err!= nil {
			return result, err
		}
		result[currency_id] = append(result[currency_id], map[int64]map[string]float64{time:{"miner":miner, "user":user}})
	}
	return result, nil
}


func(db *DCDB) GetHolidays(userId int64) ([][]int64, error) {
	var result [][]int64
	rows, err := db.Query("SELECT start_time, end_time FROM holidays WHERE user_id = ? AND delete = 0", userId)
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		var start_time, end_time int64
		err = rows.Scan(&start_time, &end_time)
		if err!= nil {
			return result, err
		}
		result = append(result, []int64{start_time, end_time})
	}
	return result, nil
}


func (db *DCDB) GetRepaidAmount(currencyId, userId int64) (float64, error) {
	amount, err := db.Single("SELECT amount FROM promised_amount WHERE status = 'repaid' AND currency_id = ? AND user_id = ? AND del_block_id = 0 AND del_mining_block_id = 0", currencyId, userId).Float64()
	if err != nil {
		return 0, err
	}
	return amount, nil
}
/*	function get_repaid_amount($currency_id, $user_id)
	{
		return $this->db->query( __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__, "
				SELECT `amount`
				FROM `".DB_PREFIX."promised_amount`
				WHERE `status` = 'repaid' AND
							 `currency_id` = {$currency_id} AND
							 `user_id` = {$user_id} AND
							 `del_block_id` = 0 AND
							 `del_mining_block_id` = 0
				", 'fetch_one');
	}*/
func (db *DCDB) CheckCurrency(currency_id int64) (bool, error) {
	id, err := db.Single("SELECT id FROM currency WHERE id = ?", currency_id).Int()
	if err != nil {
		return false, err
	}
	if id == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (db *DCDB) GetUserIdByPublicKey(publicKey []byte) (string, error) {
	var sql string
	switch db.ConfigIni["db_type"] {
	case "sqlite":
		sql = `SELECT user_id FROM users WHERE public_key_0 = $1`
	case "postgresql":
		sql = `SELECT user_id FROM users WHERE public_key_0 = '\x$1'`
	case "mysql":
		sql = `SELECT user_id FROM users WHERE public_key_0 = 0x$1`
	}
	userId, err := db.Single(sql, publicKey).String()
	if err != nil{
		return "", ErrInfo(err)
	}
	return userId, nil
}

func (db *DCDB) InsertIntoMyKey(userId string, publicKey []byte, curBlockId string) error {
	var sql string
	switch db.ConfigIni["db_type"] {
	case "sqlite":
		sql = `INSERT INTO `+userId+`_my_keys (public_key, status, block_id) VALUES ($1,'approved', $2)`
	case "postgresql":
		sql = `INSERT INTO `+userId+`_my_keys (public_key, status, block_id) VALUES ('\x$1','approved', $2)`
	case "mysql":
		sql = `INSERT INTO `+userId+`_my_keys (public_key, status, block_id) VALUES (0x$1,'approved', $2)`
	}
	err := db.ExecSql(sql, publicKey, curBlockId )
	if err != nil {
		return err
	}
	return nil
}

func (db *DCDB) GetInfoBlock() (map[string]string, error) {
	var result map[string]string
	result, err := db.OneRow("SELECT * FROM info_block").String()
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
	result, err := db.Single("SELECT local_gate_ip FROM config").String()
	if err != nil {
		return "", err
	}
	return result, nil
}
func (db *DCDB) GetNodePublicKey(userId int64) ([]byte, error) {
	result, err := db.Single("SELECT node_public_key FROM miners_data WHERE user_id = ?", userId).Bytes()
	if err != nil {
		return []byte(""), err
	}
	return result, nil
}
func (db *DCDB) GetCountCurrencies() (int64, error) {
	result, err := db.Single("SELECT count(id) FROM currency").Int64()
	if err != nil {
		return 0, err
	}
	return result, nil
}

func (db *DCDB) UpdMainLock() error {
	err := db.ExecSql("UPDATE main_lock SET lock_time = ?", time.Now().Unix())
	return err
}

func (db *DCDB) DbLock(name string) error {
	var mutex = &sync.Mutex{}
	var ok bool
	for {
		mutex.Lock()
		exists, err := db.OneRow("SELECT lock_time, script_name FROM main_lock").String()
		if err != nil {
			return ErrInfo(err)
		}
		if exists["script_name"] == name {
			err = db.ExecSql("UPDATE main_lock SET lock_time = ?", time.Now().Unix())
			if err != nil {
				return ErrInfo(err)
			}
			ok = true
		} else if len(exists["script_name"])==0 {
			err = db.ExecSql(`INSERT INTO main_lock(lock_time, script_name) VALUES(?, ?)`, time.Now().Unix(), name)
			if err != nil {
				return ErrInfo(err)
			}
			ok = true
		}
		mutex.Unlock()
		if !ok {
			time.Sleep(time.Duration(RandInt(100, 200)) * time.Millisecond)
		} else {
			break
		}
	}
	fmt.Println("OK")
	return nil
}

func (db *DCDB) DbUnlock(name string) error {
	err := db.ExecSql("DELETE FROM main_lock WHERE script_name=?", name)
	if err != nil {
		return err
	}
	return nil
}

func (db *DCDB) GetIsReadySleep(level int64, data []int64) int64 {
	//SleepData := db.GetSleepData();
	return GetIsReadySleep0(level, data)
}

func (db *DCDB) GetGenSleep(prevBlock *prevBlockType, level int64) (int64, error) {

	sleepData, err := db.GetSleepData()
	if err != nil {
		return 0, err
	}

	// узнаем время, которые было затрачено в ожидании is_ready предыдущим блоком
	isReadySleep := db.GetIsReadySleep(prevBlock.Level , sleepData["is_ready"])
	//fmt.Println("isReadySleep", isReadySleep)

	// сколько сек должен ждать нод, перед тем, как начать генерить блок, если нашел себя в одном из уровней.
	generatorSleep := GetGeneratorSleep(level, sleepData["generator"])
	//fmt.Println("generatorSleep", generatorSleep)

	// сумма is_ready всех предыдущих уровней, которые не успели сгенерить блок
	isReadySleep2 := GetIsReadySleepSum(level , sleepData["is_ready"])
	//fmt.Println("isReadySleep2", isReadySleep2)

	// узнаем, сколько нам нужно спать
	sleep := isReadySleep + generatorSleep + isReadySleep2;
	return sleep, nil;
}

func (db *DCDB) UpdDaemonTime(name string) {

}

func (db *DCDB) GetBlockDataFromBlockChain(blockId int64) (*BlockData, error) {
	BlockData := new(BlockData)
	data, err := db.OneRow("SELECT * FROM block_chain WHERE id = ?", blockId).String()
	if err!=nil {
		return BlockData, err
	}
	if len(data["data"]) > 0 {
		binaryData := []byte(data["data"])
		BytesShift(&binaryData, 1)  // не нужно. 0 - блок, >0 - тр-ии
		BlockData = ParseBlockHeader(&binaryData)
		BlockData.Hash = BinToHex([]byte(data["hash"]))
		BlockData.HeadHash = BinToHex([]byte(data["head_hash"]))
	}
	return BlockData, nil
}
/*	function get_prev_block($block_id)
	{
		if (!$this->prev_block) {
			$data = $this->db->query( __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__, "
						SELECT `hash`, `head_hash`, `data`
						FROM `".DB_PREFIX."block_chain`
						WHERE `id` = {$block_id}
						", 'fetch_array' );
			$binary_data = $data['data'];
			string_shift ($binary_data, 1 ); // 0 - блок, >0 - тр-ии
			$this->block_info = parse_block_header($binary_data);
			$this->block_info['hash'] = bin2hex($data['hash']);
			$this->block_info['head_hash'] = bin2hex($data['head_hash']);
		}
		$this->prev_block = $this->block_info;

	}*/
