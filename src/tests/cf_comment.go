package main

import (
	"fmt"
	"utils"
	"tests_utils"
)

func main() {

	f:=tests_utils.InitLog()
	defer f.Close()

	txType := "CfComment";
	txTime := "1427383713";
	userId := []byte("2")
	var blockId int64 = 128008

	var txSlice [][]byte
	// hash
	txSlice = append(txSlice, []byte("1111111111"))
	// type
	txSlice = append(txSlice,  utils.Int64ToByte(utils.TypeArray(txType)))
	// time
	txSlice = append(txSlice, []byte(txTime))
	// user_id
	txSlice = append(txSlice, []byte("1"))
	//project_id
	txSlice = append(txSlice, []byte("1"))
	//lang_id
	txSlice = append(txSlice, []byte("1"))
	//comment
	txSlice = append(txSlice, []byte(`К"))о"))м"))м"))е"))н"))т)) http://google.com/ 54fds56f4sd`))
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
