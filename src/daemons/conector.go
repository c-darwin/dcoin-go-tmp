package daemons

import (
	"utils"
	"consts"
	"time"
	"log"
	"strings"
	"net"
)

func check(host string, blockId int64) string {

	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Println(utils.ErrInfo(err))
		return "0"
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Println(utils.ErrInfo(err))
		return "0"
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(consts.READ_TIMEOUT * time.Second))
	conn.SetWriteDeadline(time.Now().Add(consts.WRITE_TIMEOUT * time.Second))

	// вначале шлем тип данных, чтобы принимающая сторона могла понять, как именно надо обрабатывать присланные данные
	_, err = conn.Write(utils.DecToBin(4, 1))
	if err != nil {
		log.Println(utils.ErrInfo(err))
		return "0"
	}

	// в 4-х байтах пишем ID блока, хэш которого хотим получить
	size := utils.DecToBin(blockId, 4)
	_, err = conn.Write(size)
	if err != nil {
		log.Println(utils.ErrInfo(err))
		return "0"
	}

	// ответ всегда 16 байт
	hash := make([]byte, 16)
	_, err = conn.Read(hash)
	if err != nil {
		log.Println(utils.ErrInfo(err))
		return "0"
	}
	return string(utils.BinToHex(hash))
}

func isReachable(host string, blockId int64, ch0 chan string) {
	log.Println("IsReachable", host)
	ch := make(chan string, 1)
	go func() {
		ch <- check(host, blockId)
	}()
	select {
		case reachable := <-ch:
		ch0 <- reachable
		case <-time.After(consts.WAIT_CONFIRMED_NODES*time.Second):
		ch0 <-  "0"
	}
}

func Connector(configIni map[string]string) {

	GoroutineName := "connector"
	db := utils.DbConnect(configIni)
	db.GoroutineName = GoroutineName
	BEGIN:
	for {

		// проверим, не нужно нам выйти, т.к. обновилась версия софта
		if db.CheckDaemonRestart() {
			utils.Sleep(1)
			break
		}

		var hosts []map[string]string
		nodeConfig, err := db.GetNodeConfig()
		if len(nodeConfig["local_gate_ip"]) > 0 {
			utils.Sleep(5)
			continue
		}

		// ровно стольким нодам мы будем слать хэши блоков и тр-ий
		var maxHosts = consts.OUT_CONNECTIONS
		if utils.StrToInt64(nodeConfig["out_connections"]) > 0 {
			maxHosts = utils.StrToInt64(nodeConfig["out_connections"])
		}
		collective, er := db.GetCommunityUsers()
		if err != nil {
			db.PrintSleep(err, 1)
			continue
		}
		if len(collective) == 0 {
			myUserId, err := db.GetMyUserId("")
			if err != nil {
				db.PrintSleep(err, 1)
				continue
			}
			collective = append(collective, myUserId)
		}
		// в сингл-моде будет только $my_miners_ids[0]
		myMinersIds, err := db.GetMyMinersIds(collective);
		if err != nil {
			db.PrintSleep(err, 1)
			continue
		}
		nodesBan, err := db.GetList(`
				SELECT host, ban_start
				FROM nodes_ban
				LEFT JOIN miners_data ON miners_data.user_id = nodes_ban.user_id
				`, "host", "ban_start").String()

		nodesConnections, err := db.GetAll(`
				SELECT nodes_connection.host,
							 nodes_connection.user_id,
							 ban_start,
							 miner_id
				FROM nodes_connection
				LEFT JOIN nodes_ban ON nodes_ban.user_id = nodes_connection.user_id
				LEFT JOIN miners_data ON miners_data.user_id = nodes_connection.user_id
				`, -1)
		for _, data := rage nodesConnections {

			// проверим, не нужно нам выйти, т.к. обновилась версия софта
			if db.CheckDaemonRestart() {
				utils.Sleep(1)
				break
			}
		}
	}

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
				IsReachable(host, blockId, ch)
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
