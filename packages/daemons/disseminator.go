package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	//"log"
	"strings"
)

/*
 * просто шлем всем, кто есть в nodes_connection хэши блока и тр-ий
 * если мы не майнер, то шлем всю тр-ию целиком, блоки слать не можем
 * если майнер - то шлем только хэши, т.к. у нас есть хост, откуда всё можно скачать
 * */
func Disseminator() {

	GoroutineName := "Disseminator"

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

		var hosts []map[string]string
		var nodeData map[string]string
		nodeConfig, err := db.GetNodeConfig()
		if len(nodeConfig["local_gate_ip"]) == 0 {
			// обычный режим
			hosts, err = db.GetAll(`
					SELECT miners_data.user_id, miners_data.tcp_host as host, node_public_key
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
			nodeData, err = db.OneRow("SELECT node_public_key, tcp_host FROM miners_data WHERE user_id  =  ?", nodeConfig["static_node_user_id"]).String()
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
		var changeNodeKey int64
		if len(myUsersIds) > 0 {
			changeNodeKey, err = db.Single(`
				SELECT count(*)
				FROM transactions
				WHERE type = ? AND
							 user_id IN (`+strings.Join(utils.SliceInt64ToString(myUsersIds), ",")+`)
				`, utils.TypeInt("ChangeNodeKey")).Int64()
			if err != nil {
				db.PrintSleep(err, 1)
				continue BEGIN
			}
		}

		var dataType int64 // это тип для того, чтобы принимающая сторона могла понять, как именно надо обрабатывать присланные данные

		// если я майнер и работаю в обычном режиме, то должен слать хэши
		if len(myMinersIds) > 0 && len(nodeConfig["local_gate_ip"]) == 0 && changeNodeKey == 0 {

			dataType = 1

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
			if len(transactions) == 0 {
				db.PrintSleep("len(transactions) == 0", 1)
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
				for _, host := range hosts {
					userId := utils.StrToInt64(host["user_id"])
					go func() {

						// шлем данные указанному хосту
						conn, err := utils.TcpConn(host["host"])
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}
						defer conn.Close()

						randTestblockHash, err := db.Single("SELECT head_hash FROM queue_testblock").String()
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}
						// получаем IV + ключ + зашифрованный текст
						dataToBeSent, key, iv, err := utils.EncryptData(toBeSent, []byte(host["node_public_key"]), randTestblockHash)
						log.Debug("key: %s", key)
						log.Debug("iv: %s", iv)
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}

						// вначале шлем тип данных, чтобы принимающая сторона могла понять, как именно надо обрабатывать присланные данные
						_, err = conn.Write(utils.DecToBin(dataType, 1))
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}

						// т.к. на приеме может быть пул, то нужно дописать в начало user_id, чьим нодовским ключем шифруем
						dataToBeSent = append(utils.DecToBin(userId, 5), dataToBeSent...)
						log.Debug("dataToBeSent: %x", dataToBeSent)

						// в 4-х байтах пишем размер данных, которые пошлем далее
						size := utils.DecToBin(len(dataToBeSent), 4)
						_, err = conn.Write(size)
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}
						// далее шлем сами данные
						_, err = conn.Write(dataToBeSent)
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}
						// в ответ получаем размер данных, которые нам хочет передать сервер
						buf := make([]byte, 4)
						_, err =conn.Read(buf)
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}
						dataSize := utils.BinToDec(buf)
						log.Debug("dataSize %d", dataSize)
						// и если данных менее 1мб, то получаем их
						if dataSize < 1048576 {
							encBinaryTxHashes := make([]byte, dataSize)
							_, err := conn.Read(encBinaryTxHashes)
							if err != nil {
								log.Info("%v", utils.ErrInfo(err))
								return
							}
							// разбираем полученные данные
							binaryTxHashes, err := utils.DecryptCFB(iv, encBinaryTxHashes, key)
							if err != nil {
								log.Info("%v", utils.ErrInfo(err))
								return
							}
							log.Debug("binaryTxHashes %x", binaryTxHashes)
							var binaryTx []byte
							for {
								// Разбираем список транзакций
								txHash := make([]byte, 16)
								if len(binaryTxHashes) >= 16 {
									txHash = utils.BytesShift(&binaryTxHashes, 16)
								}
								txHash = utils.BinToHex(txHash)
								tx, err := db.Single("SELECT data FROM transactions WHERE hash  =  [hex]", txHash).Bytes()
								log.Debug("tx %x", tx)
								if err != nil {
									log.Info("%v", utils.ErrInfo(err))
									return
								}
								if len(tx) > 0 {
									binaryTx = append(binaryTx, utils.EncodeLengthPlusData(tx)...)
								}
								if len (binaryTxHashes) == 0 {
									break
								}
							}

							log.Debug("binaryTx %x", binaryTx)
							// шифруем тр-ии. Вначале encData добавляется IV
							encData, _, err := utils.EncryptCFB(binaryTx, key, iv)
							if err != nil {
								log.Info("%v", utils.ErrInfo(err))
								return
							}
							log.Debug("encData %x", encData)

							// шлем серверу
							// в первых 4-х байтах пишем размер данных, которые пошлем далее
							size := utils.DecToBin(len(encData), 4)
							_, err = conn.Write(size)
							if err != nil {
								log.Info("%v", utils.ErrInfo(err))
								return
							}
							// далее шлем сами данные
							_, err = conn.Write(encData)
							if err != nil {
								log.Info("%v", utils.ErrInfo(err))
								return
							}
						}
					}()
				}
			}
		} else {
			var remoteNodeHost string
			// если просто юзер или работаю в защищенном режиме, то шлю тр-ии целиком. слать блоки не имею права.
			if len(nodeConfig["local_gate_ip"]) > 0 {
				dataType = 3
				remoteNodeHost = nodeData["host"]
			} else {
				dataType = 2
				remoteNodeHost = ""
			}

			log.Debug("dataType: %d", dataType)

			var toBeSent []byte // сюда пишем все тр-ии, которые будем слать другим нодам
			// возьмем хэши и сами тр-ии
			rows, err := db.Query("SELECT hash, data FROM transactions WHERE sent  =  0")
			if err != nil {
				db.PrintSleep(err, 1)
				continue BEGIN
			}
			defer rows.Close()
			for rows.Next() {
				var hash, data []byte
				err = rows.Scan(&hash, &data)
				if err != nil {
					db.PrintSleep(err, 1)
					continue BEGIN
				}
				log.Debug("hash %x", hash)
				hashHex := utils.BinToHex(hash)
				err = db.ExecSql("UPDATE transactions SET sent = 1 WHERE hash = [hex]", hashHex)
				if err != nil {
					db.PrintSleep(err, 1)
					continue BEGIN
				}
				toBeSent = append(toBeSent, data...)
			}

			// шлем тр-ии
			if len(toBeSent) > 0 {
				for _, host := range hosts {
					userId := utils.StrToInt64(host["user_id"])
					go func() {

						conn, err := utils.TcpConn(host["host"])
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}
						defer conn.Close()

						randTestblockHash, err := db.Single("SELECT head_hash FROM queue_testblock").String()
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}
						// получаем IV + ключ + зашифрованный текст
						encryptedData, _, _, err := utils.EncryptData(toBeSent, []byte(host["node_public_key"]), randTestblockHash)
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}

						// вначале шлем тип данных, чтобы принимающая сторона могла понять, как именно надо обрабатывать присланные данные
						_, err = conn.Write(utils.DecToBin(dataType, 1))
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}

						// т.к. на приеме может быть пул, то нужно дописать в начало user_id, чьим нодовским ключем шифруем
						/*_, err = conn.Write(utils.DecToBin(userId, 5))
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}*/
						encryptedData = append(utils.DecToBin(userId, 5), encryptedData...)

						// это может быть защищенное локальное соедниение (dataType = 3) и принимающему ноду нужно знать, куда дальше слать данные и чьим они зашифрованы ключем
						if len(remoteNodeHost) > 0 {
							/*
							_, err = conn.Write([]byte(remoteNodeHost))
							if err != nil {
								log.Info("%v", utils.ErrInfo(err))
								return
							}*/
							encryptedData = append([]byte(remoteNodeHost), encryptedData...)
						}

						// в 4-х байтах пишем размер данных, которые пошлем далее
						size := utils.DecToBin(len(encryptedData), 4)
						_, err = conn.Write(size)
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}
						// далее шлем сами данные
						_, err = conn.Write(encryptedData)
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}

					}()
				}
			}
		}

		db.DbUnlock()

		utils.Sleep(1)

		log.Info("%v", "Happy end")
	}


}


