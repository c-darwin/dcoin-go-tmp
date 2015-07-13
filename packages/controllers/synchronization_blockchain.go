package controllers
import (
	"fmt"
	"time"
	"dcoin/packages/utils"
    "encoding/json"
	"log"
)


func (c *Controller) SynchronizationBlockchain() (string, error) {

	blockData, err:=c.DCDB.GetInfoBlock()
	if err != nil {
		return "", err
	}
	blockId := blockData["block_id"]
	blockTime := blockData["time"]
	if len(blockId)==0 {
		blockId = "0"
	}
	if len(blockTime)==0 {
		blockTime = "0"
	}

	wTime := int64(12)
	wTimeReady := int64(1)
	if c.ConfigIni["test_mode"] == "1" {
		wTime = 365*24
		wTimeReady = 365*24
	}
	// если время менее 12 часов от текущего, то выдаем не подвержденные, а просто те, что есть в блокчейне
	if time.Now().Unix() - utils.StrToInt64(blockData["time"]) < 3600*wTime  {
		lastBlockData, err := c.DCDB.GetLastBlockData()
		if err != nil {
			return "", err
		}
		log.Println("lastBlockData", lastBlockData)
		// если уже почти собрали все блоки
		if time.Now().Unix() - lastBlockData["lastBlockTime"] < 3600*wTimeReady {
			blockId = "-1"
			blockTime = "-1"
		}
	}

	result := map[string]string{"block_id": blockId, "block_time": blockTime}
	resultJ, _ := json.Marshal(result)
	fmt.Println(string(resultJ))

	return string(resultJ), nil
}
