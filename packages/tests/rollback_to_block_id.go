package main

import (
	"fmt"
//	"github.com/c-darwin/dcoin-tmp/packages/utils"
	"tests_utils"
	"github.com/c-darwin/dcoin-tmp/packages/dcparser"
)

func main() {

	f:=tests_utils.InitLog()
	defer f.Close()

	db := tests_utils.DbConn()
	parser := new(dcparser.Parser)
	parser.DCDB = db
	err := parser.RollbackToBlockId(136711)
	if err!=nil {
		fmt.Println(err)
	}

}
