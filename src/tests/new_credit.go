package main

import (
	"fmt"
	"utils"
	"tests_utils"
)

func main() {

	f:=tests_utils.InitLog()
	defer f.Close()

	txType := "NewCredit";
	txTime := "1427383713";
	userId := []byte("2")
	var blockId int64 = 128008

	var txSlice [][]byte
	// hash
	txSlice = append(txSlice, []byte("22cb812e53e22ee539af4a1d39b4596d"))
	// type
	txSlice = append(txSlice,  utils.Int64ToByte(utils.TypeInt(txType)))
	// time
	txSlice = append(txSlice, []byte(txTime))
	// user_id
	txSlice = append(txSlice, []byte("1"))
	//to_user_id
	txSlice = append(txSlice, []byte("2"))
	//amount
	txSlice = append(txSlice, []byte("1000"))
	//currency_id
	txSlice = append(txSlice, []byte("72"))
	//pct
	txSlice = append(txSlice, []byte("10.00"))
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
