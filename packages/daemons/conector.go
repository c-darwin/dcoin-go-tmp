package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"github.com/c-darwin/dcoin-go-tmp/packages/consts"
	"time"
	"net"
	"io/ioutil"
	"strings"
)


func Connector() {

	GoroutineName := "Connector"

	db := DbConnect()
	if db == nil {
		return
	}
	db.GoroutineName = GoroutineName
	if !db.CheckInstall(DaemonCh, AnswerDaemonCh) {
		return
	}

	BEGIN:
	for {
		log.Info(GoroutineName)
		// проверим, не нужно ли нам выйти из цикла
		if CheckDaemonsRestart() {
			break BEGIN
		}

		nodeConfig, err := db.GetNodeConfig()
		if len(nodeConfig["local_gate_ip"]) > 0 {
			utils.Sleep(5)
			continue
		}

		var delMiners []string
		var hosts []map[string]string
		var nodesInc string
		var nodeCount int64
		idArray := make(map[int]int64)

		// ровно стольким нодам мы будем слать хэши блоков и тр-ий
		var maxHosts = consts.OUT_CONNECTIONS
		if utils.StrToInt64(nodeConfig["out_connections"]) > 0 {
			maxHosts = utils.StrToInt(nodeConfig["out_connections"])
		}
		log.Info("%v", maxHosts)
		collective, err := db.GetCommunityUsers()
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
		log.Info("%v", myMinersIds)
		nodesBan, err := db.GetMap(`
				SELECT tcp_host, ban_start
				FROM nodes_ban
				LEFT JOIN miners_data ON miners_data.user_id = nodes_ban.user_id
				`, "tcp_host", "ban_start")
		log.Info("%v", nodesBan)
		nodesConnections, err := db.GetAll(`
				SELECT nodes_connection.host,
							 nodes_connection.user_id,
							 ban_start,
							 miner_id
				FROM nodes_connection
				LEFT JOIN nodes_ban ON nodes_ban.user_id = nodes_connection.user_id
				LEFT JOIN miners_data ON miners_data.user_id = nodes_connection.user_id
				`, -1)
		for _, data := range nodesConnections {

			// проверим, не нужно нам выйти, т.к. обновилась версия софта
			if db.CheckDaemonRestart() {
				utils.Sleep(1)
				break BEGIN
			}

			// проверим соотвествие хоста и user_id
			ok, err := db.Single("SELECT user_id FROM miners_data WHERE user_id  = ? AND tcp_host  =  ?", data["user_id"], data["host"]).Int64()
			if err != nil {
				utils.Sleep(1)
				continue BEGIN
			}
			if ok == 0 {
				err = db.ExecSql("DELETE FROM nodes_connection WHERE host = ? OR user_id = ?", data["host"], data["user_id"])
				if err != nil {
					utils.Sleep(1)
					continue BEGIN
				}
			}

			// если нода забанена недавно
			if utils.StrToInt64(data["ban_start"]) > utils.Time() - consts.NODE_BAN_TIME {
				delMiners = append(delMiners, data["miner_id"])
				err = db.ExecSql("DELETE FROM nodes_connection WHERE host = ? OR user_id = ?", data["host"], data["user_id"])
				if err != nil {
					utils.Sleep(1)
					continue BEGIN
				}
				continue
			}

			hosts = append(hosts, map[string]string{"host": data["host"], "user_id": data["user_id"]})
			nodesInc += data["host"]+";"+data["user_id"]+"\n"
			nodeCount++
		}

		log.Debug("hosts: %v", hosts)
		ch := make(chan *answerType)
		for _, host := range hosts {
			userId := utils.StrToInt64(host["user_id"])
			go func(userId int64, host string) {
				ch_ := make(chan *answerType, 1)
				go func() {
					log.Debug("host: %v / userId: %v", host, userId)
					ch_ <- check(host, userId)
				}()
				select {
					case reachable := <-ch_:
						ch <- reachable
					case <-time.After(consts.WAIT_CONFIRMED_NODES*time.Second):
						ch <-  &answerType{userId: userId, answer: 0}
				}
			}(userId, host["host"])
		}

		// если нода не отвечает, то удалем её из таблы nodes_connection
		for i := 0; i < len(hosts); i++ {
			result := <-ch
			if result.answer == 0 {
				log.Info("delete %v", result.userId)
				err = db.ExecSql("DELETE FROM nodes_connection WHERE user_id = ?", result.userId)
				if err != nil {
					db.PrintSleep(err, 1)
				}
			}
			log.Info("answer: %v", result)
		}

		// добьем недостающие хосты до $max_hosts
		if len(hosts) < maxHosts {
			need := maxHosts - len(hosts)
			max, err := db.Single("SELECT max(miner_id) FROM miners").Int()
			if err != nil {
				db.PrintSleep(err, 1)
				continue BEGIN
			}
			i0:=0
			for {
				rand := 1
				if max > 1 {
					rand = utils.RandInt(1, max+1)
				}
				idArray[rand] = 1
				i0++
				if i0>30 || len(idArray)>=need || len(idArray)>=max {
					break
				}
			}
			log.Info("%v", "idArray", idArray)
			// удалим себя
			for _, id := range myMinersIds {
				delete(idArray, int(id))
			}
			// Удалим забаннные хосты
			for _, id := range delMiners {
				delete(idArray, utils.StrToInt(id))
			}
			log.Info("%v", "idArray", idArray)
			ids := ""
			if len(idArray) > 0 {
				for id, _ := range idArray {
					ids+=utils.IntToStr(id)+","
				}
				ids = ids[:len(ids)-1]
				minersHosts, err := db.GetMap(`
						SELECT tcp_host, user_id
						FROM miners_data
						WHERE miner_id IN (`+ids+`)`, "tcp_host", "user_id")
				for host, userId := range minersHosts {
					if len(nodesBan[host]) > 0 {
						if utils.StrToInt64(nodesBan[host]) > utils.Time() - consts.NODE_BAN_TIME {
							continue
						}
					}
					hosts = append(hosts, map[string]string{"host": host, "user_id": userId})
					err = db.ExecSql("DELETE FROM nodes_connection WHERE host = ?", host)
					if err != nil {
						db.PrintSleep(err, 1)
						continue BEGIN
					}
					log.Debug(host)
					err = db.ExecSql("INSERT INTO nodes_connection ( host, user_id ) VALUES ( ?, ? )", host, userId)
					if err != nil {
						db.PrintSleep(err, 1)
						continue BEGIN
					}
				}
			}
		}

		log.Debug("%v", "hosts", hosts)
		// если хосты не набрались из miner_data, то берем из файла
		if len(hosts) < 10 {
			hostsData_, err := ioutil.ReadFile("nodes.inc")
			if err != nil {
				db.PrintSleep(err, 1)
				continue BEGIN
			}
			hostsData := strings.Split(string(hostsData_), "\n")
			log.Debug("%v", "hostsData_", hostsData_)
			log.Debug("%v", "hostsData", hostsData)
			max := 0
			log.Debug("maxHosts: %v", maxHosts)
			if len(hosts) > maxHosts-1 {
				max = maxHosts
			} else {
				max = len(hosts)
			}
			log.Debug("max: %v", max)
			for i:=0; i < max; i++ {
				r := utils.RandInt(0, max)
				if len(hostsData) <= r {
					continue
				}
				hostUserId := strings.Split(hostsData[r], ";")
				if len(hostUserId) == 1 {
					continue
				}
				host, userId := hostUserId[0], hostUserId[1]
				if utils.InSliceInt64(utils.StrToInt64(userId), collective) {
					continue
				}
				if len(nodesBan[host]) > 0 {
					if utils.StrToInt64(nodesBan[host]) > utils.Time() - consts.NODE_BAN_TIME {
						continue
					}
				}

				err = db.ExecSql("DELETE FROM nodes_connection WHERE host = ?", host)
				if err != nil {
					db.PrintSleep(err, 1)
					continue BEGIN
				}
				log.Debug(host)
				err = db.ExecSql("INSERT INTO nodes_connection ( host, user_id ) VALUES ( ?, ? )", host, userId)
				if err != nil {
					db.PrintSleep(err, 1)
					continue BEGIN
				}
			}
		}

		if nodeCount > 5 {
			nodesInc = nodesInc[:len(nodesInc)-1]
			err := ioutil.WriteFile("nodes.inc", []byte(nodesInc), 0644)
			if err != nil {
				db.PrintSleep(err, 1)
				continue BEGIN
			}
		}

		for i:=0; i < 10; i++ {
			if db.CheckDaemonRestart() {
				utils.Sleep(1)
				break BEGIN
			}
			utils.Sleep(1)
		}
	}
}

type answerType struct {
	userId int64
	answer int64
}

func check(host string, userId int64) *answerType {

	/*tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Info("%v", utils.ErrInfo(err))
		return &answerType{userId: userId, answer: 0}
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)*/
	conn, err := net.DialTimeout("tcp", host, 5 * time.Second)

	if err != nil {
		log.Info("%v", utils.ErrInfo(err))
		return &answerType{userId: userId, answer: 0}
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(consts.READ_TIMEOUT * time.Second))
	conn.SetWriteDeadline(time.Now().Add(consts.WRITE_TIMEOUT * time.Second))

	// вначале шлем тип данных, чтобы принимающая сторона могла понять, как именно надо обрабатывать присланные данные
	_, err = conn.Write(utils.DecToBin(5, 1))
	if err != nil {
		log.Info("%v", utils.ErrInfo(err))
		return &answerType{userId: userId, answer: 0}
	}

	// в 5-и байтах пишем userID, чтобы проверить, верный ли у него нодовский ключ, т.к. иначе ему нельзя слать зашифрованные данные
	_, err = conn.Write(utils.DecToBin(userId, 5))
	if err != nil {
		log.Info("%v", utils.ErrInfo(err))
		return &answerType{userId: userId, answer: 0}
	}

	// ответ всегда 1 байт. 0 или 1
	answer := make([]byte, 1)

	_, err = conn.Read(answer)
	if err != nil {
		log.Info("%v", utils.ErrInfo(err))
		return &answerType{userId: userId, answer: 0}
	}
	log.Debug("host: %v / answer: %v", host, answer)
	return &answerType{userId: userId, answer: utils.BinToDec(answer)}
}
