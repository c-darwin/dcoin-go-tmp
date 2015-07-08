package daemons

import (
	"utils"
	"log"
	"strings"
	"net"
)

/*
 * просто шлем всем, кто есть в nodes_connection хэши блока и тр-ий
 * если мы не майнер, то шлем всю тр-ию целиком, блоки слать не можем
 * если майнер - то шлем только хэши, т.к. у нас есть хост, откуда всё можно скачать
 * */
func Disseminator(configIni map[string]string) string {

	GoroutineName := "disseminator"
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
		if len(nodeConfig["local_gate_ip"]) == 0 {
			// обычный режим
			hosts, err = db.GetAll(`
					SELECT miners_data.user_id, miners_data.host, node_public_key
					FROM nodes_connection
					LEFT JOIN miners_data ON nodes_connection.user_id = miners_data.user_id
					`, -1)
			if err != nil {
				db.PrintSleep(err, 1)
				continue
			}
			if len(hosts) == 0 {
				utils.Sleep(1)
				continue
			}
		} else {
			// защищенный режим
			nodeData, err := db.OneRow("SELECT node_public_key, host FROM miners_data WHERE user_id  =  ?", nodeConfig["static_node_user_id"]).String()
			if err != nil {
				db.PrintSleep(err, 1)
				continue
			}
			hosts = append(hosts, map[string]string{"host":nodeConfig["local_gate_ip"], "node_public_key":nodeData["node_public_key"], "user_id":nodeConfig["static_node_user_id"]})
		}

		myUsersIds, err := db.GetMyUsersIds(false)
		myMinersIds, err := db.GetMyMinersIds(myUsersIds)

		// если среди тр-ий есть смена нодовского ключа, то слать через отправку хэшей с последющей отдачей данных может не получиться
		// т.к. при некорректном нодовском ключе придет зашифрованый запрос на отдачу данных, а мы его не сможем расшифровать т.к. ключ у нас неверный
		changeNodeKey, err := db.Single(`
				SELECT count(*)
				FROM transactions
				WHERE type = ? AND
							 user_id IN (`+strings.Join(utils.SliceInt64ToString(myUsersIds), ",")+`)
				`, utils.TypeInt("ChangeNodeKey")).Int64()
		if err != nil {
			db.PrintSleep(err, 1)
			continue BEGIN
		}

		// если я майнер и работаю в обычном режиме, то должен слать хэши
		if len(myMinersIds) > 0 && len(nodeConfig["local_gate_ip"]) == 0 && changeNodeKey == 0 {

			// определим, от кого будем слать
			r := utils.RandInt(0, len(myMinersIds))
			myMinerId := myMinersIds[r]
			myUserId, err := db.Single("SELECT user_id FROM miners_data WHERE miner_id  =  ?", myMinerId).Int64()
			if err != nil {
				db.PrintSleep(err, 1)
				continue BEGIN
			}

			// возьмем хэш текущего блока и номер блока
			// для теста ролбеков отключим на время
			data, err := db.OneRow("SELECT block_id, hash, head_hash FROM info_block WHERE sent  =  0").Bytes()
			if err != nil {
				db.PrintSleep(err, 1)
				continue BEGIN
			}
			err = db.ExecSql("UPDATE info_block SET sent = 1")
			if err != nil {
				db.PrintSleep(err, 1)
				continue BEGIN
			}

			/*
			 * Составляем данные на отправку
			 * */
			// 5 байт = наш user_id. Но они будут не первые, т.к. m_curl допишет вперед user_id получателя (нужно для пулов)
			toBeSent := utils.DecToBin(myUserId, 5);
			if len(data) > 0 {  // блок
				// если 5-й байт = 0, то на приемнике будем читать блок, если = 1 , то сразу хэши тр-ий
				toBeSent = append(toBeSent, utils.DecToBin(0, 1)...)
				toBeSent = append(toBeSent, utils.DecToBin(utils.BytesToInt64(data["block_id"]), 3)...)
				toBeSent = append(toBeSent, data["hash"]...)
				toBeSent = append(toBeSent, data["head_hash"]...)
				err = db.ExecSql("UPDATE info_block SET sent = 1")
				if err != nil {
					db.PrintSleep(err, 1)
					continue BEGIN
				}
			} else { // тр-ии без блока
				toBeSent = append(toBeSent, utils.DecToBin(1, 1)...)
			}

			// возьмем хэши тр-ий
			transactions, err := db.GetAll("SELECT hash, high_rate FROM transactions WHERE sent = 0 AND for_self_use = 0", -1)
			if err != nil {
				db.PrintSleep(err, 1)
				continue BEGIN
			}
			for _, data := range transactions {
				hexHash := utils.BinToHex([]byte(data["hash"]))
				toBeSent = append(toBeSent, utils.DecToBin(utils.StrToInt64(data["high_rate"]), 1)...)
				toBeSent = append(toBeSent, []byte(data["hash"])...)
				err = db.ExecSql("UPDATE transactions SET sent = 1 WHERE hash = [hex]", hexHash)
				if err != nil {
					db.PrintSleep(err, 1)
					continue BEGIN
				}
			}

			// отправляем блок и хэши тр-ий, если есть что отправлять
			if len(toBeSent) > 0 {
				for _, data := range hosts {
					go func() {
						SendHashes(toBeSent, data["host"], []byte(data["node_public_key"]), utils.StrToInt64(data["user_id"]))
					}()
				}
			}


		}

		db.DbUnlock()

		utils.Sleep(1)

		log.Println("Happy end")
	}

	return ""
}

/*
 * $remote_node_user_id - это когда идет пересылка уже зашифрованной тр-ии внутри сети. Чтобы принимающая сторона могла понять,
 * какому ноду слать эту тр-ию, пишем в первые 5 байт user_id
 * */
func SendHashes(data []byte, host string, nodePublicKey []byte, userId int64, remoteNodeHost string) error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err!=nil {
		return err
	} else {

		randTestblockHash, err := db.Single("SELECT head_hash FROM queue_testblock ORDER BY RAND()").String()
		if err != nil {
			return err
		}
		// получаем IV + ключ + зашифрованный текст
		encryptedData, err := utils.EncryptData(data, nodePublicKey, randTestblockHash)
		// т.к. на приеме может быть пул, то нужно дописать в начало user_id, чьим нодовским ключем шифруем
		encryptedData = append(utils.DecToBin(userId, 5), encryptedData...)
		if len(remoteNodeHost) > 0 {
			encryptedData = append(utils.EncodeLengthPlusData(remoteNodeHost), encryptedData...)
		}
		// в первых 4-х байтах пишем размер данных, которые пошлем далее
		size := utils.DecToBin(len(encryptedData), 4)
		conn.Write(size)
		// далее шлем сами данные
		conn.Write(encryptedData)
		// в ответ получаем размер данных, которые нам хочет передать сервер
		buf := make([]byte, 4)
		_, err :=conn.Read(buf)
		if err != nil {
			return err
		}
		dataSize := utils.BinToDec(buf)
		// и если данных менее 1мб, то получаем их
		if dataSize < 1048576 {
			buf := make([]byte, dataSize)
			_, err := conn.Read(buf)
			if err != nil {
				return err
			}
			// разбираем полученные данные
			// ключ нужен чтобы зашифровать данные, которые пошлем
			key, data, err := utils.DecryptData(&buf)
			if err != nil {
				return err
			}
			// Разбираем список транзакций



		}
		conn.Close()
	}
}
