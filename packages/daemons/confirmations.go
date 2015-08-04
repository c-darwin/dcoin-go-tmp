package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"github.com/c-darwin/dcoin-go-tmp/packages/consts"
	"time"
	//"log"
	"strings"
	"net"
)

/*
Получаем кол-во нодов, у которых такой же хэш последнего блока как и у нас
Нужно чтобы следить за вилками
*/

func Confirmations() {

	const GoroutineName = "Confirmations"
	d := new(daemon)
	d.DCDB = DbConnect()
	if d.DCDB == nil {
		return
	}
	d.goRoutineName = GoroutineName
	if !d.CheckInstall(DaemonCh, AnswerDaemonCh) {
		return
	}

	BEGIN:
	for {
		log.Info(GoroutineName)
		// проверим, не нужно ли нам выйти из цикла
		if CheckDaemonsRestart() {
			break BEGIN
		}
		blockId, err := d.GetBlockId()
		hash, err := d.Single("SELECT hash FROM block_chain WHERE id =  ?", blockId).String()
		if err != nil {
			log.Info("%v", err)
		}
		log.Info("%v", hash)

		var hosts []map[string]string
		if d.ConfigIni["test_mode"] == "1" {
			hosts = []map[string]string {{"tcp_host":"localhost:8088", "user_id":"1"}}
		} else {
			maxMinerId, err := d.Single("SELECT max(miner_id) FROM miners_data").Int64()
			if err != nil {
				log.Info("%v", err)
			}
			if maxMinerId == 0 {
				maxMinerId = 1
			}
			q := ""
			if d.ConfigIni["db_type"] == "postgresql" {
				q = "SELECT DISTINCT ON (tcp_host) tcp_host, user_id FROM miners_data WHERE miner_id IN ("+strings.Join(utils.RandSlice(1, maxMinerId, consts.COUNT_CONFIRMED_NODES), ",")+")"
			} else {
				q = "SELECT tcp_host, user_id FROM miners_data WHERE miner_id IN  ("+strings.Join(utils.RandSlice(1, maxMinerId, consts.COUNT_CONFIRMED_NODES), ",")+") GROUP BY tcp_host"
			}
			hosts, err = d.GetAll(q, consts.COUNT_CONFIRMED_NODES)
			if err != nil {
				log.Info("%v", err)
			}
		}

		ch := make(chan string)
		for i := 0; i < len(hosts); i++ {
			log.Info("hosts[i] %v", hosts[i])
			host:=hosts[i]["tcp_host"];
			log.Info("host %v", host)
			go func() {
				IsReachable(host, blockId, ch)
			}()
		}
		var answer string
		var st0, st1 int64
		for i := 0; i < len(hosts); i++ {
			answer = <-ch
			log.Info("answer == hash (%x = %x)", answer, hash)
			log.Info("answer == hash (%s = %s)", answer, hash)
			if answer == hash {
				st1 ++
			} else {
				st0 ++
			}
			log.Info("%v", "CHanswer", answer)
		}
		exists, err := d.Single("SELECT block_id FROM confirmations WHERE block_id= ?", blockId).Int64()
		if exists > 0 {
			log.Debug("UPDATE confirmations SET good = %v, bad = %v, time = %v WHERE block_id = %v", st1, st0, time.Now().Unix(), blockId)
			err = d.ExecSql("UPDATE confirmations SET good = ?, bad = ?, time = ? WHERE block_id = ?", st1, st0, time.Now().Unix(), blockId)
			if err != nil {
				log.Info("%v", err)
			}
		} else {
			log.Debug("INSERT INTO confirmations ( block_id, good, bad, time) VALUES ( %v, %v, %v, %v )", blockId, st1, st0, time.Now().Unix())
			err = d.ExecSql("INSERT INTO confirmations ( block_id, good, bad, time ) VALUES ( ?, ?, ?, ? )", blockId, st1, st0, time.Now().Unix())
			if err != nil {
				log.Info("%v", err)
			}
		}
		for i:=0; i < 60; i++ {
			utils.Sleep(1)
			// проверим, не нужно ли нам выйти из цикла
			if CheckDaemonsRestart() {
				break BEGIN
			}
		}
	}
}

func checkConf(host string, blockId int64) string {

	log.Debug("host: %v", host)
	/*tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		log.Info("%v", utils.ErrInfo(err))
		return "0"
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)*/
	conn, err := net.DialTimeout("tcp", host, 5 * time.Second)
	if err != nil {
		log.Info("%v", utils.ErrInfo(err))
		return "0"
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(consts.READ_TIMEOUT * time.Second))
	conn.SetWriteDeadline(time.Now().Add(consts.WRITE_TIMEOUT * time.Second))

	// вначале шлем тип данных, чтобы принимающая сторона могла понять, как именно надо обрабатывать присланные данные
	_, err = conn.Write(utils.DecToBin(4, 1))
	if err != nil {
		log.Info("%v", utils.ErrInfo(err))
		return "0"
	}

	// в 4-х байтах пишем ID блока, хэш которого хотим получить
	size := utils.DecToBin(blockId, 4)
	_, err = conn.Write(size)
	if err != nil {
		log.Info("%v", utils.ErrInfo(err))
		return "0"
	}

	// ответ всегда 32 байта
	hash := make([]byte, 32)
	_, err = conn.Read(hash)
	if err != nil {
		log.Info("%v", utils.ErrInfo(err))
		return "0"
	}
	return string(hash)
}

func IsReachable(host string, blockId int64, ch0 chan string) {
	log.Info("IsReachable %v", host)
	ch := make(chan string, 1)
	go func() {
		ch <- checkConf(host, blockId)
	}()
	select {
	case reachable := <-ch:
	ch0 <- reachable
	case <-time.After(consts.WAIT_CONFIRMED_NODES*time.Second):
	ch0 <-  "0"
	}
}
