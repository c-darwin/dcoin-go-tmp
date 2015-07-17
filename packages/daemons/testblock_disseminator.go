package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	//"log"
	"strings"
)

/**
 * Демон, который мониторит таблу testblock и если видит status=active,
 * то шлет блок строго тем, кто находятся на одном с нами уровне. Если пошлет
 * тем, кто не на одном уровне, то блок просто проигнорируется
 *
 */
func TestblockDisseminator() {

	GoroutineName := "TestblockDisseminator"
	db := DbConnect()
	db.GoroutineName = GoroutineName
	db.CheckInstall()
	for {

		// проверим, не нужно нам выйти, т.к. обновилась версия софта
		if db.CheckDaemonRestart() {
			utils.Sleep(1)
			break
		}

		nodeConfig, err := db.GetNodeConfig()
		if len(nodeConfig["local_gate_ip"]) != 0 {
			db.PrintSleep("local_gate_ip", 60)
			continue
		}

		_, _, _, _, level, levelsRange, err := db.TestBlock()
		if err != nil {
			db.PrintSleep(err, 1)
			continue
		}
		log.Debug("level: %v", level)
		log.Debug("levelsRange: %v", levelsRange)
		// получим id майнеров, которые на нашем уровне
		nodesIds := utils.GetOurLevelNodes(level, levelsRange)
		if len(nodesIds) == 0 {
			db.PrintSleep("len(nodesIds) == 0", 1)
			continue
		}
		log.Debug("nodesIds: %v", nodesIds)

		// получим хосты майнеров, которые на нашем уровне
		hosts, err := db.GetList("SELECT tcp_host FROM miners_data WHERE miner_id IN ("+strings.Join(utils.SliceInt64ToString(nodesIds), `,`)+")").String()
		if err != nil {
			db.PrintSleep(err, 1)
			continue
		}
		log.Debug("hosts: %v", hosts)

		// шлем block_id, user_id, mrkl_root, signature
		data, err := db.OneRow("SELECT block_id, time, user_id, mrkl_root, signature FROM testblock WHERE status  =  'active'").String()
		if err != nil {
			db.PrintSleep(err, 1)
			continue
		}
		if len(data) > 0 {

			dataToBeSent := utils.DecToBin(utils.StrToInt64(data["block_id"]), 4)
			dataToBeSent = append(dataToBeSent, utils.DecToBin(data["time"], 4)...)
			dataToBeSent = append(dataToBeSent, utils.DecToBin(data["user_id"], 4)...)
			dataToBeSent = append(dataToBeSent, []byte(data["mrkl_root"])...)
			dataToBeSent = append(dataToBeSent, utils.EncodeLengthPlusData(data["signature"])...)

			for _, host := range hosts {
				go func(host string) {
					log.Debug("host: %v", host)
					conn, err := utils.TcpConn(host)
					if err != nil {
						log.Info("%v", utils.ErrInfo(err))
						return
					}
					defer conn.Close()

					// вначале шлем тип данных
					_, err = conn.Write(utils.DecToBin(6, 1))
					if err != nil {
						log.Info("%v", utils.ErrInfo(err))
						return
					}

					// в 4-х байтах пишем размер данных, которые пошлем далее
					_, err = conn.Write(utils.DecToBin(len(dataToBeSent), 4))
					if err != nil {
						log.Info("%v", utils.ErrInfo(err))
						return
					}
					// далее шлем сами данные
					log.Debug("dataToBeSent: %x", dataToBeSent)
					_, err = conn.Write(dataToBeSent)
					if err != nil {
						log.Info("%v", utils.ErrInfo(err))
						return
					}

					/*
					 * Получаем тр-ии, которые есть у юзера, в ответ выдаем те, что недостают и
					 * их порядок следования, чтобы получить валидный блок
					 */
					buf := make([]byte, 4)
					_, err =conn.Read(buf)
					if err != nil {
						log.Info("%v", utils.ErrInfo(err))
						return
					}
					dataSize := utils.BinToDec(buf)
					// и если данных менее 10мб, то получаем их
					if dataSize < 10485760 {

						data, err := db.OneRow("SELECT * FROM testblock").String()
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}

						responseBinaryData := utils.DecToBin(utils.StrToInt64(data["block_id"]), 4)
						responseBinaryData = append(responseBinaryData, utils.DecToBin(utils.StrToInt64(data["time"]), 4)...)
						responseBinaryData = append(responseBinaryData, utils.DecToBin(utils.StrToInt64(data["user_id"]), 5)...)
						responseBinaryData = append(responseBinaryData, utils.EncodeLengthPlusData(data["signature"])...)

						addSql := ""
						if dataSize > 0 {
							binaryData := make([]byte, dataSize)
							_, err := conn.Read(binaryData)
							if err != nil {
								log.Info("%v", utils.ErrInfo(err))
								return
							}

							// разбираем присланные данные
							// получим хэши тр-ий, которые надо исключить
							for {
								if len(binaryData) < 16 {
									break
								}
								txHex := utils.BinToHex(utils.BytesShift(&binaryData, 16))
								// проверим
								addSql+=string(txHex)+","
								if len(binaryData) == 0 {
									break
								}
							}
							addSql = addSql[:len(addSql)-1]
							addSql = "WHERE id NOT IN ("+addSql+")"
						}
						// сами тр-ии
						var transactions []byte
						transactions_testblock, err := db.GetList(`SELECT data FROM transactions_testblock `+addSql).String()
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}
						for _, txData := range transactions_testblock {
							transactions = append(transactions, utils.EncodeLengthPlusData(txData)...)
						}

						responseBinaryData = append(responseBinaryData, utils.EncodeLengthPlusData(transactions)...)

						// порядок тр-ий
						transactions_testblock, err = db.GetList(`SELECT hash FROM transactions_testblock ORDER BY id ASC`).String()
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}
						for _, txHash := range transactions_testblock {
							responseBinaryData = append(responseBinaryData, []byte(txHash)...)
						}

						// в 4-х байтах пишем размер данных, которые пошлем далее
						_, err = conn.Write(utils.DecToBin(len(responseBinaryData), 4))
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}

						// далее шлем сами данные
						_, err = conn.Write(responseBinaryData)
						if err != nil {
							log.Info("%v", utils.ErrInfo(err))
							return
						}
					}

				}(host)
			}
		}

		utils.Sleep(1)

		log.Info("%v", "Happy end")
	}


}


