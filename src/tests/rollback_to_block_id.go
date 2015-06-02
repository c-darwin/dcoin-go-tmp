package main

import (
	"fmt"
//	"utils"
	"tests_utils"
	"dcparser"
)

func main() {

	f:=tests_utils.InitLog()
	defer f.Close()

	db := tests_utils.DbConn()
	parser := new(dcparser.Parser)
	parser.DCDB = db
	err := parser.RollbackToBlockId(108209)
	if err!=nil {
		fmt.Println(err)
	}

}
