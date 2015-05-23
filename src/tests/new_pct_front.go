package main

import (
	"fmt"
//	"database/sql"
	//"dcparser"
	"utils"
	"tests_utils"
	//_ "github.com/lib/pq"
	//"encoding/binary"
	//"bytes"
	//"encoding/hex"
	//"crypto/rand"
	//"crypto/rsa"
	//"crypto/sha1"
	//"daemons"
//	"strconv"
	//"errors"
	"log"
	"os"
	//"github.com/alyu/configparser"
	"github.com/astaxie/beego/config"
	//"strings"
	//"regexp"
	//"reflect"
//	"consts"
	"io"
)


func main() {

	f, err := os.OpenFile("dclog.txt", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0777)
	defer f.Close()
	//log.SetOutput(f)
	log.SetOutput(io.MultiWriter(f, os.Stdout))
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	txType := "NewPct";
	txTime := "1427383713";
	userId := []byte("2")
	var blockId int64 = 1288
	data:=`{"1":{"miner_pct":"0.0000000044318","user_pct":"0.0000000027036"},"72":{"miner_pct":"0.0000000047610","user_pct":"0.0000000029646"}}`

	var txSlice [][]byte
	// hash
	txSlice = append(txSlice, []byte("22cb812e53e22ee539af4a1d39b4596d"))
	// type
	txSlice = append(txSlice,  utils.Int64ToByte(utils.TypeArray(txType)))
	// time
	txSlice = append(txSlice, []byte(txTime))
	// user_id
	txSlice = append(txSlice, userId)
	//new_pct
	txSlice = append(txSlice, []byte(data))

	dataForSign := fmt.Sprintf("%v,%v,%s,%v", utils.TypeArray(txType), txTime, userId, data)

	blockData := new(utils.BlockData)
	blockData.BlockId = blockId
	blockData.Time = utils.StrToInt64(txTime)
	blockData.UserId = utils.BytesToInt64(userId)

	configIni_, err := config.NewConfig("ini", "config.ini")
	if err != nil {
		fmt.Println(err)
	}
	configIni, err := configIni_.GetSection("default")

	db := utils.DbConnect(configIni)
	err = tests_utils.MakeFrontTest(txSlice, utils.StrToInt64(txTime), dataForSign, txType, utils.BytesToInt64(userId), "", blockId, db)
	if err != nil {
		fmt.Println(err)
	}

}
