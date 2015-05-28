package main

import (
	"fmt"
	"utils"
	"tests_utils"
)

func main() {

	tests_utils.InitLog()

	txType := "ChangeHost";
	txTime := "1399278817";
	userId := []byte("2")
	var blockId int64 = 1415
	host:="http://fdfdfd.ru/"

	var txSlice [][]byte
	// hash
	txSlice = append(txSlice, []byte("22cb812e53e22ee539af4a1d39b4596d"))
	// type
	txSlice = append(txSlice,  utils.Int64ToByte(utils.TypeArray(txType)))
	// time
	txSlice = append(txSlice, []byte(txTime))
	// user_id
	txSlice = append(txSlice, userId)
	// promised_amount_id
	txSlice = append(txSlice, []byte(host))

	dataForSign := fmt.Sprintf("%v,%v,%s,%s", utils.TypeArray(txType), txTime, userId, host)

	err := tests_utils.MakeFrontTest(txSlice, utils.StrToInt64(txTime), dataForSign, txType, utils.BytesToInt64(userId), "", blockId)
	if err != nil {
		fmt.Println(err)
	}
}
