package main

import (
	"fmt"
	"dcoin/packages/utils"
	"tests_utils"
)

func main() {

	f:=tests_utils.InitLog()
	defer f.Close()

	txType := "SendDc";
	txTime := "1409288580";
	userId := []byte("2")
	var blockId int64 = 10000

	var txSlice [][]byte
	// hash
	txSlice = append(txSlice, []byte("22cb812e53e22ee539af4a1d39b4596d"))
	// type
	txSlice = append(txSlice,  utils.Int64ToByte(utils.TypeInt(txType)))
	// time
	txSlice = append(txSlice, []byte(txTime))
	// user_id
	txSlice = append(txSlice, userId)
	// to_user_id
	txSlice = append(txSlice, []byte("2"))
	// currency_id
	txSlice = append(txSlice, []byte("72"))
	// amount
	txSlice = append(txSlice, []byte("8"))
	// commission
	txSlice = append(txSlice, []byte("0.1"))
/*	for i:=0; i<5; i++ {
		txSlice = append(txSlice, []byte("0"))
	}
	for i:=0; i<5; i++ {
		txSlice = append(txSlice, []byte("0"))
	}*/
	// comment
	txSlice = append(txSlice, []byte("1111111111111111111111111111111111"))
	// sign
	txSlice = append(txSlice, []byte("11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111"))

	blockData := new(utils.BlockData)
	blockData.BlockId = blockId
	blockData.Time = utils.StrToInt64(txTime)
	blockData.UserId = utils.BytesToInt64(userId)

	err := tests_utils.MakeTest(txSlice, blockData, txType, "work_and_rollback");
	if err != nil {
		fmt.Println(err)
	}

}
