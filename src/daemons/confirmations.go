package daemons

import (
	"utils"
	"consts"
	"time"
	"os"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)


func check(u string) string {
	resp, err := http.Get(u)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	return string(htmlData)
}

func IsReachable(url string, ch0 chan string) {
	log.Println("IsReachable", url)
	ch := make(chan string, 1)
	go func() {
		ch <- check(url)
	}()
	select {
	case reachable := <-ch:
	ch0 <- reachable
	case <-time.After(consts.WAIT_CONFIRMED_NODES*time.Second):
	ch0 <-  "0"
	}
}

func Confirmations(configIni map[string]string) {

	const GoroutineName = "confirmations"

	db := utils.DbConnect(configIni)
	for {
		blockId, err := db.GetBlockId()
		hash, err := db.Single("SELECT hash FROM block_chain WHERE id =  ?", blockId).String()
		if err != nil {
			log.Println(err)
		}
		log.Println(hash)

		var hosts []map[string]string
		if db.ConfigIni["test_mode"] == "1" {
			hosts = []map[string]string {{"host":"http://localhost:8089/", "user_id":"1"}}
		} else {
			maxMinerId, err := db.Single("SELECT max(miner_id) FROM miners_data").Int64()
			if err != nil {
				log.Println(err)
			}
			q := ""
			if db.ConfigIni["db_type"] == "postgresql" {
				q = "SELECT DISTINCT ON (host) host, user_id FROM miners_data WHERE miner_id IN ("+strings.Join(utils.RandSlice(1, maxMinerId, consts.COUNT_CONFIRMED_NODES), ",")+")"
			} else {
				q = "SELECT host, user_id FROM miners_data WHERE miner_id IN  ("+strings.Join(utils.RandSlice(1, maxMinerId, consts.COUNT_CONFIRMED_NODES), ",")+") GROUP BY host"
			}
			hosts, err = db.GetAll(q, consts.COUNT_CONFIRMED_NODES)
			if err != nil {
				log.Println(err)
			}
		}

		ch := make(chan string)
		for i := 0; i < len(hosts); i++ {
			log.Println("hosts[i]", hosts[i])
			host:=hosts[i]["host"];
			log.Println("host", host)
			go func() {
				IsReachable(host+"tools?controllerName=checkNode&block_id="+utils.Int64ToStr(blockId), ch)
			}()
		}
		var answer string
		var st0, st1 int64
		for i := 0; i < len(hosts); i++ {
			answer = <-ch
			log.Println("answer == hash", answer, hash)
			if answer == hash {
				st1 ++
			} else {
				st0 ++
			}
			log.Println("CHanswer", answer)
		}
		exists, err := db.Single("SELECT block_id FROM confirmations WHERE block_id= ?", blockId).Int64()
		if exists > 0 {
			err = db.ExecSql("UPDATE confirmations SET good = ?, bad = ?, time = ? WHERE block_id = ?", st1, st0, time.Now().Unix(), blockId)
			if err != nil {
				log.Println(err)
			}
		} else {
			err = db.ExecSql("INSERT INTO confirmations ( block_id, good, bad, time ) VALUES ( ?, ?, ?, ? )", blockId, st1, st0, time.Now().Unix())
			if err != nil {
				log.Println(err)
			}
		}
		utils.Sleep(60)
	}
}
