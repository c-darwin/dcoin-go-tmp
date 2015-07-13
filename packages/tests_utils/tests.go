package tests_utils

import (
	"regexp"
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"fmt"
	"crypto/rsa"
	"crypto/rand"
	"crypto"
	"os"
	"github.com/astaxie/beego/config"
	"log"
	"io"
	"strings"
//	"crypto/rand"
//	"crypto/rsa"
	"encoding/pem"
	"crypto/x509"
	"encoding/base64"
)

func genKeys() (string, string) {
	privatekey, _ := rsa.GenerateKey(rand.Reader, 1024)
	var pemkey = &pem.Block{Type : "RSA PRIVATE KEY", Bytes : x509.MarshalPKCS1PrivateKey(privatekey)}
	PrivBytes0 := pem.EncodeToMemory(&pem.Block{Type:  "RSA PRIVATE KEY", Bytes: pemkey.Bytes})

	PubASN1, _ := x509.MarshalPKIXPublicKey(&privatekey.PublicKey)
	pubBytes := pem.EncodeToMemory(&pem.Block{Type:  "RSA PUBLIC KEY", Bytes: PubASN1})
	s := strings.Replace(string(pubBytes),"-----BEGIN RSA PUBLIC KEY-----","",-1)
	s = strings.Replace(s,"-----END RSA PUBLIC KEY-----","",-1)
	sDec, _ := base64.StdEncoding.DecodeString(s)

	return string(PrivBytes0), fmt.Sprintf("%x", sDec)
}

// для юнит-тестов. снимок всех данных в БД
func AllHashes(db *utils.DCDB) (map[string]string, error) {
	//var orderBy string
	result:=make(map[string]string)
	//var columns string;
	tables, err:=db.GetAllTables()
	if err != nil {
		return result, err
	}
	/*rows, err := db.Query(`
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
*/
	for _, table :=range tables {
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
			return result, utils.ErrInfo(err)
		}
		result[table] = hash
	}
	return result, nil
}

func DbConn() *utils.DCDB {
	configIni_, err := config.NewConfig("ini", "config.ini")
	if err != nil {
		fmt.Println(err)
	}
	configIni, err := configIni_.GetSection("default")
	db := utils.DbConnect(configIni)
	return db
}


func InitLog() *os.File {
	f, err := os.OpenFile("dclog.txt", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0777)
	if err!=nil{
		fmt.Println(err)
	}
	//log.SetOutput(f)
	log.SetOutput(io.MultiWriter(f, os.Stdout))
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	return f
}

func MakeFrontTest(transactionArray [][]byte, time int64, dataForSign string, txType string, userId int64, MY_PREFIX string, blockId int64) error {

	db := DbConn()

	priv, pub  := genKeys()

	nodeArr := []string{"new_admin", "votes_node_new_miner", "NewPct"}
	var binSign []byte
	if utils.InSliceString(txType, nodeArr) {

		err:=db.ExecSql("UPDATE my_node_keys SET private_key = ?", priv)
		if err != nil {
			return utils.ErrInfo(err)
		}
		err=db.ExecSql("UPDATE miners_data SET node_public_key = [hex] WHERE user_id = ?", pub, userId)
		if err != nil {
			return utils.ErrInfo(err)
		}


		k, err := db.GetNodePrivateKey(MY_PREFIX)
		if err != nil {
			return utils.ErrInfo(err)
		}
		fmt.Println("k", k)
		privateKey, err := utils.MakePrivateKey(k)
		if err != nil {
			return utils.ErrInfo(err)
		}
		//fmt.Println("privateKey.PublicKey", privateKey.PublicKey)
		//fmt.Println("privateKey.D", privateKey.D)
		//fmt.Printf("privateKey.N %x\n", privateKey.N)
		//fmt.Println("privateKey.Public", privateKey.Public())
		binSign, err = rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, utils.HashSha1(dataForSign))
		//nodePublicKey, err := db.GetNodePublicKey(userId)
		//fmt.Println("nodePublicKey", nodePublicKey)
		//if err != nil {
		//	return utils.ErrInfo(err)
		//}
		//CheckSignResult, err := utils.CheckSign([][]byte{nodePublicKey}, dataForSign, binSign, true);
		//fmt.Printf("binSign: %x\n", binSign)
		//fmt.Println("err", err)
		//fmt.Println("CheckSignResult", CheckSignResult)

	} else {

		err:=db.ExecSql("UPDATE my_keys SET private_key = ?", priv)
		if err != nil {
			return utils.ErrInfo(err)
		}
		err=db.ExecSql("UPDATE users SET public_key_0 = [hex]", pub)
		if err != nil {
			return utils.ErrInfo(err)
		}

		k, err := db.GetPrivateKey(MY_PREFIX)
		privateKey, err := utils.MakePrivateKey(k)
		if err != nil {
			return utils.ErrInfo(err)
		}
		binSign, err = rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, utils.HashSha1(dataForSign))
		binSign = utils.EncodeLengthPlusData(binSign)
	}

	//fmt.Println("HashSha1", utils.HashSha1(dataForSign))
	//fmt.Printf("binSign %x\n", binSign)
	//fmt.Println("dataForSign", dataForSign)
	transactionArray = append(transactionArray, binSign)

	parser := new(dcparser.Parser)
	parser.DCDB = db
	parser.GoroutineName = "test"
	parser.TxSlice = transactionArray
	parser.BlockData = &utils.BlockData{BlockId: blockId, Time: time, UserId: userId}
	parser.TxHash = []byte("111111111111111")
	parser.Variables, _ = parser.DCDB.GetAllVariables()

	err0 := utils.CallMethod(parser, txType+"Init")
	if i, ok := err0.(error); ok {
	fmt.Println(err0.(error), i)
	return err0.(error)
	}
	err0 = utils.CallMethod(parser, txType+"Front")
	if i, ok := err0.(error); ok {
	fmt.Println(err0.(error), i)
	return err0.(error)
	}
	err0 = utils.CallMethod(parser, txType+"RollbackFront")
	if i, ok := err0.(error); ok {
	fmt.Println(err0.(error), i)
	return err0.(error)
	}
	return nil
}


func MakeTest(txSlice [][]byte, blockData *utils.BlockData, txType string, testType string) error {

	db := DbConn()

	parser := new(dcparser.Parser)
	parser.DCDB = db
	parser.TxSlice = txSlice
	parser.BlockData = blockData
	parser.TxHash = []byte("111111111111111")
	parser.Variables, _ = db.GetAllVariables()


	// делаем снимок БД в виде хэшей до начала тестов
	hashesStart, err := AllHashes(db)
	if err!=nil {
		return err
	}

	//fmt.Println("dcparser."+txType+"Init")
	err0 := utils.CallMethod(parser, txType+"Init")
	if i, ok := err0.(error); ok {
		fmt.Println(err0.(error), i)
		return err0.(error)
	}

	if testType == "work_and_rollback" {

		err0 = utils.CallMethod(parser, txType)
		if i, ok := err0.(error); ok {
			fmt.Println(err0.(error), i)
			return err0.(error)
		}

		//fmt.Println("-------------------")
		// узнаем, какие таблицы были затронуты в результате выполнения основного метода
		hashesMiddle, err := AllHashes(db)
		if err != nil {
			return utils.ErrInfo(err)
		}
		var tables []string
		//fmt.Println("hashesMiddle", hashesMiddle)
		//fmt.Println("hashesStart", hashesStart)
		for table, hash := range hashesMiddle {
			if hash!=hashesStart[table] {
				tables = append(tables, table)
			}
		}
		fmt.Println("tables", tables)

		// rollback
		err0 := utils.CallMethod(parser, txType+"Rollback")
		if i, ok := err0.(error); ok {
			fmt.Println(err0.(error), i)
			return err0.(error)
		}

		// сраниим хэши, которые были до начала и те, что получились после роллбэка
		hashesEnd, err := AllHashes(db)
		if err != nil {
			return utils.ErrInfo(err)
		}
		for table, hash := range hashesEnd {
			if hash!=hashesStart[table] {
				fmt.Println("ERROR in table ", table)
			}
		}

	} else if (len(os.Args)>1 && os.Args[1] == "w") || testType == "work" {
		err0 = utils.CallMethod(parser, txType)
		if i, ok := err0.(error); ok {
			fmt.Println(err0.(error), i)
			return err0.(error)
		}
	} else if (len(os.Args)>1 && os.Args[1] == "r")  || testType == "rollback" {
		err0 = utils.CallMethod(parser, txType+"Rollback")
		if i, ok := err0.(error); ok {
			fmt.Println(err0.(error), i)
			return err0.(error)
		}
	}
	return nil
}
