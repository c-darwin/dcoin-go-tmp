package controllers
import (
	//"database/sql"
	//_ "github.com/lib/pq"
	//"reflect"
	"fmt"
	//"html/template"
	//"bufio"
	//"bytes"
	"time"
	"utils"
    "encoding/json"
)

type synchronizationBlockchainStruct struct {

}

func (c *Controller) Synchronization_blockchain() (string, error) {
	fmt.Println("Synchronization_blockchain")

	blockData, err:=c.DCDB.GetInfoBlock()
	if err != nil {
		return "", err
	}
	blockId := blockData["bock_id"]
	blockTime := blockData["time"]
	if len(blockId)==0 {
		blockId = "0"
	}
	if len(blockTime)==0 {
		blockTime = "0"
	}

	// если время более 12 часов от текущего, то выдаем не подвержденные, а просто те, что есть в блокчейне
	if time.Now().Unix() - utils.StrToInt64(blockData["time"]) < 3600*12  {
		lastBlockData, err := c.DCDB.GetLastBlockData()
		if err != nil {
			return "", err
		}
		// если уже почти собрали все блоки
		if time.Now().Unix() - lastBlockData["lastBlockTime"] < 3600 {
			blockId = "-1"
			blockTime = "-1"
		}
	}

	result := map[string]string{"block_id": blockId, "block_time": blockTime}
	resultJ, _ := json.Marshal(result)
	fmt.Println(string(resultJ))

	return string(resultJ), nil
}
