package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/consts"
	"github.com/c-darwin/dcoin-go-tmp/packages/static"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"
	"regexp"
	"fmt"
)

var myUserIdForChat int64

func (d *daemon) chatConnector() {
	log.Debug("start chatConnector")
	maxMinerId, err := d.Single("SELECT max(miner_id) FROM miners_data").Int64()
	if err != nil {
		log.Error("%v", err)
	}
	if maxMinerId == 0 {
		maxMinerId = 1
	}
	q := ""
	if d.ConfigIni["db_type"] == "postgresql" {
		q = "SELECT DISTINCT ON (tcp_host) tcp_host, user_id FROM miners_data WHERE miner_id IN (" + strings.Join(utils.RandSlice(1, maxMinerId, consts.COUNT_CHAT_NODES), ",") + ")"
	} else {
		q = "SELECT tcp_host, user_id FROM miners_data WHERE miner_id IN  (" + strings.Join(utils.RandSlice(1, maxMinerId, consts.COUNT_CHAT_NODES), ",") + ") GROUP BY tcp_host"
	}
	hosts, err := d.GetAll(q, consts.COUNT_CHAT_NODES)
	if err != nil {
		log.Error("%v", err)
	}
	// исключим себя
	myTcpHost, err := d.Single(`SELECT tcp_host FROM miners_data WHERE user_id = ?`, myUserIdForChat).String()
	if err != nil {
		log.Error("%v", err)
	}
	fmt.Println("myTcpHost:", myTcpHost)

	// исключим хосты, к которым уже подключены
	var uids string
	for userId, _ := range utils.ChatOutConnections {
		uids+=utils.Int64ToStr(userId)+","
	}
	uids = uids[:len(uids)-1]
	existsTcpHost, err := d.GetList(`SELECT tcp_host FROM miners_data WHERE user_id IN (`+uids+`)`).String()
	if err != nil {
		log.Error("%v", err)
	}

	log.Debug("hosts: %v", hosts)
	for _, data := range hosts {

		host := data["tcp_host"]
		userId := utils.StrToInt64(data["user_id"])

		if host == myTcpHost || utils.InSliceString(host, existsTcpHost) {
			continue
		}
		go func(host string, userId int64) {

			log.Debug("host: %v", host)
			log.Debug("userId: %d", userId)
			re := regexp.MustCompile(`(.*?):[0-9]+$`)
			match := re.FindStringSubmatch(host)
			log.Debug("match: %v", match)

			if len(match) != 0 {

				log.Debug("myUserIdForChat %v", myUserIdForChat)
				log.Debug("chat host: %v", match[1]+consts.CHAT_PORT)
				chatHost := match[1]+":"+consts.CHAT_PORT
				//chatHost := "192.168.150.30:8087"

				// проверим, нет ли уже созданных каналов для такого хоста
				if _, ok := utils.ChatOutConnections[userId]; !ok {

					// канал для приема тр-ий чата
					conn, err := net.DialTimeout("tcp", chatHost, 5*time.Second)
					if err != nil {
						log.Error("%v", utils.ErrInfo(err))
						return
					} else {
						log.Debug(conn.RemoteAddr().String(), conn)
						myUid := utils.DecToBin(myUserIdForChat, 4)
						log.Debug("myUid %x", myUid)
						n, err := conn.Write(myUid)
						log.Debug("n: %d", n)
						if err != nil {
							log.Error("%v", utils.ErrInfo(err))
							return
						}
						n, err = conn.Write(utils.DecToBin(1, 1))
						log.Debug("n: %d", n)
						if err != nil {
							log.Error("%v", utils.ErrInfo(err))
							return
						}
						fmt.Println("connector ChatInput", conn.RemoteAddr(), utils.Time())
						utils.ChatMutex.Lock()
						utils.ChatInConnections[userId] = 1
						utils.ChatMutex.Unlock()
						go utils.ChatInput(conn, userId)
					}

					// канал для отправки тр-ий чата
					conn2, err := net.DialTimeout("tcp", chatHost, 5*time.Second)
					if err != nil {
						log.Error("%v", utils.ErrInfo(err))
						return
					} else {
						log.Debug(conn2.RemoteAddr().String(), conn2)
						n, err := conn2.Write(utils.DecToBin(myUserIdForChat, 4))
						log.Debug("n: %d", n)
						if err != nil {
							log.Error("%v", utils.ErrInfo(err))
							return
						}
						n, err = conn2.Write(utils.DecToBin(0, 1))
						log.Debug("n: %d", n)
						if err != nil {
							log.Error("%v", utils.ErrInfo(err))
							return
						}

						fmt.Println("connector ADD", userId, conn2.RemoteAddr(), utils.Time())
						utils.ChatMutex.Lock()
						utils.ChatOutConnections[userId] = 1
						utils.ChatMutex.Unlock()
						fmt.Println("ChatOutConnections", utils.ChatOutConnections)
						utils.ChatTxDisseminator(conn2, userId)
					}
				}
			}
		}(host, userId)
	}
}

func Connector() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("daemon Recovered", r)
			panic(r)
		}
	}()

	if _, err := os.Stat(*utils.Dir + "/nodes.inc"); os.IsNotExist(err) {
		data, err := static.Asset("static/nodes.inc")
		if err != nil {
			log.Error("%v", err)
		}
		err = ioutil.WriteFile(*utils.Dir+"/nodes.inc", []byte(data), 0644)
		if err != nil {
			log.Error("%v", err)
		}
	}


	GoroutineName := "Connector"
	d := new(daemon)
	d.DCDB = DbConnect()
	if d.DCDB == nil {
		return
	}
	d.goRoutineName = GoroutineName
	if utils.Mobile() {
		d.sleepTime = 600
	} else {
		d.sleepTime = 30
	}
	if !d.CheckInstall(DaemonCh, AnswerDaemonCh) {
		return
	}
	d.DCDB = DbConnect()
	if d.DCDB == nil {
		return
	}


	// соединения для чата иногда отваливаются, поэтому в цикле мониторим состояние
	go func() {
		for {
			if myUserIdForChat == 0 {
				utils.Sleep(1)
				continue
			}
			if len(utils.ChatOutConnections) < 1 || len(utils.ChatInConnections) < 1 {
				go d.chatConnector()
			}
			utils.Sleep(30)
		}
	}()

	BEGIN:
	for {
		log.Info(GoroutineName)
		MonitorDaemonCh <- []string{GoroutineName, utils.Int64ToStr(utils.Time())}

		// проверим, не нужно ли нам выйти из цикла
		if CheckDaemonsRestart() {
			break BEGIN
		}

		nodeConfig, err := d.GetNodeConfig()
		if len(nodeConfig["local_gate_ip"]) > 0 {
			utils.Sleep(2)
			continue
		}

		var delMiners []string
		var hosts []map[string]string
		var nodeCount int64
		idArray := make(map[int]int64)
		nodesInc := make(map[string]string)

		// ровно стольким нодам мы будем слать хэши блоков и тр-ий
		var maxHosts = consts.OUT_CONNECTIONS
		if utils.StrToInt64(nodeConfig["out_connections"]) > 0 {
			maxHosts = utils.StrToInt(nodeConfig["out_connections"])
		}
		log.Info("%v", maxHosts)

		collective, err := d.GetCommunityUsers()
		if err != nil {
			log.Error("%v", err)
			return
		}
		if len(collective) == 0 {
			myUserId, err := d.GetMyUserId("")
			if err != nil {
				log.Error("%v", err)
				return
			}
			collective = append(collective, myUserId)
			myUserIdForChat = myUserId
		} else {
			myUserIdForChat, err = d.Single(`SELECT pool_admin_user_id FROM config`).Int64()
			if err != nil {
				log.Error("%v", err)
				return
			}
		}

		// в сингл-моде будет только $my_miners_ids[0]
		myMinersIds, err := d.GetMyMinersIds(collective)
		if err != nil {
			if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
			continue
		}
		log.Info("%v", myMinersIds)
		nodesBan, err := d.GetMap(`
				SELECT tcp_host, ban_start
				FROM nodes_ban
				LEFT JOIN miners_data ON miners_data.user_id = nodes_ban.user_id
				`, "tcp_host", "ban_start")
		log.Info("%v", nodesBan)
		nodesConnections, err := d.GetAll(`
				SELECT nodes_connection.host,
							 nodes_connection.user_id,
							 ban_start,
							 miner_id
				FROM nodes_connection
				LEFT JOIN nodes_ban ON nodes_ban.user_id = nodes_connection.user_id
				LEFT JOIN miners_data ON miners_data.user_id = nodes_connection.user_id
				`, -1)
		//fmt.Println("nodesConnections", nodesConnections)
		log.Debug("nodesConnections: %v", nodesConnections)
		for _, data := range nodesConnections {

			// проверим, не нужно нам выйти, т.к. обновилась версия софта
			if CheckDaemonsRestart() {
				break BEGIN
			}

			/*// проверим соотвествие хоста и user_id
			ok, err := d.Single("SELECT user_id FROM miners_data WHERE user_id  = ? AND tcp_host  =  ?", data["user_id"], data["host"]).Int64()
			if err != nil {
				utils.Sleep(1)
				continue BEGIN
			}
			if ok == 0 {
				err = d.ExecSql("DELETE FROM nodes_connection WHERE host = ? OR user_id = ?", data["host"], data["user_id"])
				if err != nil {
					utils.Sleep(1)
					continue BEGIN
				}
			}*/

			// если нода забанена недавно
			if utils.StrToInt64(data["ban_start"]) > utils.Time()-consts.NODE_BAN_TIME {
				delMiners = append(delMiners, data["miner_id"])
				err = d.ExecSql("DELETE FROM nodes_connection WHERE host = ? OR user_id = ?", data["host"], data["user_id"])
				if err != nil {
					if d.dPrintSleep(utils.ErrInfo(err), d.sleepTime) {	break BEGIN }
					continue BEGIN
				}
				continue
			}

			hosts = append(hosts, map[string]string{"host": data["host"], "user_id": data["user_id"]})
			nodesInc[data["host"]] = data["user_id"]
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
				case <-time.After(consts.WAIT_CONFIRMED_NODES * time.Second):
					ch <- &answerType{userId: userId, answer: 0}
				}
			}(userId, host["host"])
		}

		log.Debug("%v", "hosts", hosts)
		var newHosts []map[string]string
		var countOk int
		// если нода не отвечает, то удалем её из таблы nodes_connection
		for i := 0; i < len(hosts); i++ {
			result := <-ch
			if result.answer == 0 {
				log.Info("delete %v", result.userId)
				err = d.ExecSql("DELETE FROM nodes_connection WHERE user_id = ?", result.userId)
				if err != nil {
					if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
				}
				for _, data := range hosts {
					if utils.StrToInt64(data["user_id"]) != result.userId {
						newHosts = append(newHosts, data)
					}
				}
			} else {
				countOk++
			}
			log.Info("answer: %v", result)
		}
		hosts = newHosts
		log.Debug("%v", "hosts", hosts)

		// добьем недостающие хосты до $max_hosts
		if len(hosts) < maxHosts {
			need := maxHosts - len(hosts)
			max, err := d.Single("SELECT max(miner_id) FROM miners").Int()
			if err != nil {
				if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
			i0 := 0
			for {
				rand := 1
				if max > 1 {
					rand = utils.RandInt(1, max+1)
				}
				idArray[rand] = 1
				i0++
				if i0 > 30 || len(idArray) >= need || len(idArray) >= max {
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
					ids += utils.IntToStr(id) + ","
				}
				ids = ids[:len(ids)-1]
				minersHosts, err := d.GetMap(`
						SELECT tcp_host, user_id
						FROM miners_data
						WHERE miner_id IN (`+ids+`)`, "tcp_host", "user_id")
				for host, userId := range minersHosts {
					if len(nodesBan[host]) > 0 {
						if utils.StrToInt64(nodesBan[host]) > utils.Time()-consts.NODE_BAN_TIME {
							continue
						}
					}
					hosts = append(hosts, map[string]string{"host": host, "user_id": userId})
					err = d.ExecSql("DELETE FROM nodes_connection WHERE host = ?", host)
					if err != nil {
						if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
						continue BEGIN
					}
					log.Debug(host)
					err = d.ExecSql("INSERT INTO nodes_connection ( host, user_id ) VALUES ( ?, ? )", host, userId)
					if err != nil {
						if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
						continue BEGIN
					}
				}
			}
		}

		log.Debug("%v", "hosts", hosts)
		// если хосты не набрались из miner_data, то берем из файла
		if len(hosts) < 10 {
			hostsData_, err := ioutil.ReadFile(*utils.Dir + "/nodes.inc")
			if err != nil {
				if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
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
				max = len(hostsData)
			}
			log.Debug("max: %v", max)
			for i := 0; i < max; i++ {
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
					if utils.StrToInt64(nodesBan[host]) > utils.Time()-consts.NODE_BAN_TIME {
						continue
					}
				}

				err = d.ExecSql("DELETE FROM nodes_connection WHERE host = ?", host)
				if err != nil {
					if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
					continue BEGIN
				}
				log.Debug(host)
				err = d.ExecSql("INSERT INTO nodes_connection ( host, user_id ) VALUES ( ?, ? )", host, userId)
				if err != nil {
					if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
					continue BEGIN
				}
				nodesInc[host] = userId

			}
		}

		if nodeCount > 5 {
			nodesFile := ""
			for k, v := range nodesInc {
				nodesFile += k + ";" + v + "\n"
			}
			nodesFile = nodesFile[:len(nodesFile)-1]
			err := ioutil.WriteFile(*utils.Dir+"/nodes.inc", []byte(nodesFile), 0644)
			if err != nil {
				if d.dPrintSleep(err, d.sleepTime) {	break BEGIN }
				continue BEGIN
			}
		}

		var sleepTime int
		if countOk < 2 {
			sleepTime = 5
		} else {
			sleepTime = d.sleepTime
		}

		if d.dSleep(sleepTime) {
			break BEGIN
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
		log.Error("%v", utils.ErrInfo(err))
		return &answerType{userId: userId, answer: 0}
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)*/
	conn, err := net.DialTimeout("tcp", host, 5*time.Second)

	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		return &answerType{userId: userId, answer: 0}
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(consts.READ_TIMEOUT * time.Second))
	conn.SetWriteDeadline(time.Now().Add(consts.WRITE_TIMEOUT * time.Second))

	// вначале шлем тип данных, чтобы принимающая сторона могла понять, как именно надо обрабатывать присланные данные
	_, err = conn.Write(utils.DecToBin(5, 1))
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		return &answerType{userId: userId, answer: 0}
	}

	// в 5-и байтах пишем userID, чтобы проверить, верный ли у него нодовский ключ, т.к. иначе ему нельзя слать зашифрованные данные
	_, err = conn.Write(utils.DecToBin(userId, 5))
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		return &answerType{userId: userId, answer: 0}
	}

	// ответ всегда 1 байт. 0 или 1
	answer := make([]byte, 1)

	_, err = conn.Read(answer)
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		return &answerType{userId: userId, answer: 0}
	}

	// создадим канал для чата
	if utils.BinToDec(answer) == 1 {

	}
	log.Debug("host: %v / answer: %v", host, answer)
	return &answerType{userId: userId, answer: utils.BinToDec(answer)}
}
