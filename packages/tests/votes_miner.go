package main

import (
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"tests_utils"
)

func main() {

	f:=tests_utils.InitLog()
	defer f.Close()

	txType := "VotesMiner";
	txTime := "1427383713";
	userId := []byte("2")
	var blockId int64 = 128008

	var txSlice [][]byte
	// hash
	txSlice = append(txSlice, []byte("1111111111"))
	// type
	txSlice = append(txSlice,  utils.Int64ToByte(utils.TypeInt(txType)))
	// time
	txSlice = append(txSlice, []byte(txTime))
	// user_id
	txSlice = append(txSlice, []byte("50"))
	// vote_id
	txSlice = append(txSlice, []byte("14"))
	// result
	txSlice = append(txSlice, []byte("1"))
	// comment
	txSlice = append(txSlice, []byte("222222222222222"))
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
