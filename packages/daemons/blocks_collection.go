package daemons

import (
    _ "github.com/lib/pq"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"github.com/c-darwin/dcoin-go-tmp/packages/consts"
	"fmt"
	"os"
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
	"github.com/c-darwin/dcoin-go-tmp/packages/static"
	"errors"
	"encoding/json"
    "flag"
)
var startBlockId = flag.Int64("startBlockId", 0, "Start block for blockCollection daemon")
var endBlockId = flag.Int64("endBlockId", 0, "End block for blockCollection daemon")

func BlocksCollection() {

    const GoroutineName = "BlocksCollection"
    d := new(daemon)    
    d.DCDB = DbConnect()
    if d.DCDB == nil {
        return
    }
    d.goRoutineName = GoroutineName
    if !d.CheckInstall(DaemonCh, AnswerDaemonCh) {
        return
    }
    d.DCDB = DbConnect()
    if d.DCDB == nil {
        return
    }
	//var cur bool
    BEGIN:
    for {
        log.Info(GoroutineName)
        log.Info("00000000")
        MonitorDaemonCh <- []string{GoroutineName, utils.Int64ToStr(utils.Time())}
        log.Info("11111111")

        // проверим, не нужно ли нам выйти из цикла
        if CheckDaemonsRestart() {
            break BEGIN
        }

        config, err := d.GetNodeConfig()
        if err != nil {
            d.PrintSleep(err, 1)
            continue BEGIN
        }

        err, restart := d.dbLock()
        if restart {
            break BEGIN
        }
        if err != nil {
            d.PrintSleep(err, 1)
            continue BEGIN
        }

        // если это первый запуск во время инсталяции
        currentBlockId, err := d.GetBlockId()
        if err != nil {
            d.unlockPrintSleep(err, 1)
            continue BEGIN
        }

        log.Info("config", config)
        log.Info("currentBlockId", currentBlockId)

		// на время тестов
		/*if !cur {
            currentBlockId = 0
            cur = true
        }*/
        parser := new(dcparser.Parser)
        parser.DCDB = d.DCDB
        parser.GoroutineName = GoroutineName
        if currentBlockId==0 || *startBlockId > 0 {
            /*
            IsNotExistBlockChain := false
            if _, err := os.Stat(*utils.Dir+"/public/blockchain"); os.IsNotExist(err) {
                IsNotExistBlockChain = true
            }*/
			if config["first_load_blockchain"]=="file"/* && IsNotExistBlockChain*/ {

                log.Info("first_load_blockchain=file")
                nodeConfig, err:= d.GetNodeConfig()
                blockchain_url:=nodeConfig["first_load_blockchain_url"]
                if len(blockchain_url) == 0 {
                    blockchain_url = consts.BLOCKCHAIN_URL
                }
                log.Debug("blockchain_url: %s", blockchain_url)
                // возможно сервер отдаст блокчейн не с первой попытки
                var blockchainSize int64
                for i:=0;i<10;i++ {
                    blockchainSize, err = utils.DownloadToFile(blockchain_url, *utils.Dir+"/public/blockchain", 3600, DaemonCh, AnswerDaemonCh)
                    if blockchainSize > consts.BLOCKCHAIN_SIZE {
                        break
                    }
                }
                if err != nil || blockchainSize < consts.BLOCKCHAIN_SIZE {
                    if err != nil {
                        log.Error("%v", utils.ErrInfo(err))
                    } else {
                        log.Info(fmt.Sprintf("%v < %v", blockchainSize, consts.BLOCKCHAIN_SIZE))
					}
                    d.unlockPrintSleep(err, 1)
                    continue BEGIN
                }

                first := true
                /*// блокчейн мог быть загружен ранее. проверим его размер


                stat, err := file.Stat()
                if err != nil {
                    d.unlockPrintSleep(err, 1)
                    file.Close()
                    continue BEGIN
                }
                if stat.Size() < consts.BLOCKCHAIN_SIZE {
                    d.unlockPrintSleep(fmt.Errorf("%v < %v", stat.Size(), consts.BLOCKCHAIN_SIZE), 1)
                    file.Close()
                    continue BEGIN
                }*/

                log.Debug("GO!")
                file, err := os.Open(*utils.Dir+"/public/blockchain")
                if err != nil {
                    d.unlockPrintSleep(err, 1)
                    continue BEGIN
                }
                for {
                    // проверим, не нужно ли нам выйти из цикла
                    if CheckDaemonsRestart() {
                        d.unlockPrintSleep(fmt.Errorf("DaemonsRestart"), 0)
                        break BEGIN
                    }
                    b1 := make([]byte, 5)
                    file.Read(b1)
                    dataSize := utils.BinToDec(b1)
                    log.Debug("dataSize", dataSize)
                    if dataSize > 0 {

                        data := make([]byte, dataSize)
                        file.Read(data)
                        log.Debug("data %x\n", data)
                        blockId := utils.BinToDec(data[0:5])
                        if *endBlockId > 0 && blockId == *endBlockId {
                            d.PrintSleep(err, 1)
                            file.Close()
                            continue BEGIN
                        }
                        log.Info("blockId", blockId)
                        data2:=data[5:]
                        length := utils.DecodeLength(&data2)
                        log.Debug("length", length)
                        log.Debug("data2 %x\n", data2)
                        blockBin := utils.BytesShift(&data2, length)
                        log.Debug("blockBin %x\n", blockBin)

                        if *startBlockId == 0 || (*startBlockId > 0 && blockId > *startBlockId) {

                            // парсинг блока
                            parser.BinaryData = blockBin;

                            if first {
                                parser.CurrentVersion = consts.VERSION
                                first = false
                            }
                            err = parser.ParseDataFull()
                            if err != nil {
                                d.PrintSleep(err, 1)
                                file.Close()
                                continue BEGIN
                            }
                            parser.InsertIntoBlockchain()

                            // отметимся, чтобы не спровоцировать очистку таблиц
                            err = parser.UpdMainLock()
                            if err != nil {
                                d.PrintSleep(err, 1)
                                file.Close()
                                continue BEGIN
                            }
                            if CheckDaemonsRestart() {
                                d.PrintSleep(err, 1)
                                file.Close()
                                continue BEGIN
                            }
                        }
                        // ненужный тут размер в конце блока данных
                        data = make([]byte, 5)
                        file.Read(data)
                    } else {
                        d.unlockPrintSleep(err, 1)
                        continue BEGIN
					}
                   // utils.Sleep(1)
                }
                file.Close()
	        } else {
                newBlock, err := static.Asset("static/1block.bin")
                if err != nil {
                    d.PrintSleep(err, 1)
                    continue BEGIN
                }
                parser.BinaryData = newBlock
                parser.CurrentVersion = consts.VERSION

                err = parser.ParseDataFull()
                if err != nil {
                    d.PrintSleep(err, 1)
                    continue BEGIN
                }
                err = parser.InsertIntoBlockchain()

                if err != nil {
                    d.PrintSleep(err, 1)
                    continue BEGIN
                }
			}

			utils.Sleep(1)
			d.dbUnlock();
            continue BEGIN
		}
		d.dbUnlock();

        myConfig, err := d.OneRow("SELECT local_gate_ip, static_node_user_id FROM config").String()
        if err != nil {
            d.PrintSleep(err, 1)
            continue
        }
		var hosts []map[string]string
		var nodeHost string
        var dataTypeMaxBlockId, dataTypeBlockBody int64
		if len(myConfig["local_gate_ip"]) > 0 {
            hosts = append(hosts, map[string]string{"host": myConfig["local_gate_ip"], "user_id": myConfig["static_node_user_id"]})
            nodeHost, err = d.Single("SELECT tcp_host FROM miners_data WHERE user_id  =  ?", myConfig["static_node_user_id"]).String()
            if err != nil {
                d.PrintSleep(err, 1)
                continue
            }
            dataTypeMaxBlockId = 9
            dataTypeBlockBody = 8
			//getBlockScriptName = "ajax?controllerName=protectedGetBlock";
			//addNodeHost = "&nodeHost="+nodeHost;
        } else {
            // получим список нодов, с кем установлено рукопожатие
            hosts, err = d.GetAll("SELECT * FROM nodes_connection", -1)
            if err != nil {
                d.PrintSleep(err, 1)
                continue
            }
            dataTypeMaxBlockId = 10
            dataTypeBlockBody = 7
            //getBlockScriptName = "ajax?controllerName=getBlock";
            //addNodeHost = "";
        }

		log.Info("%v", hosts)

		if len(hosts) == 0 {
            d.PrintSleep("len hosts = 0", 1)
            continue
		}

        maxBlockId := int64(1)
        maxBlockIdHost := ""
        var maxBlockIdUserId int64
        // получим максимальный номер блока
        for i:=0; i < len(hosts); i++ {
            if CheckDaemonsRestart() {
                break BEGIN
            }
            conn, err := utils.TcpConn(hosts[i]["host"])
            if err != nil {
                d.PrintSleep(err, 1)
                continue
            }
            // шлем тип данных
            _, err = conn.Write(utils.DecToBin(dataTypeMaxBlockId, 1))
            if err != nil {
                conn.Close()
                d.PrintSleep(err, 1)
                continue
            }
            if len(nodeHost) > 0 { // защищенный режим
                err = utils.WriteSizeAndData([]byte(nodeHost), conn)
                if err != nil {
                    conn.Close()
                    d.PrintSleep(err, 1)
                    continue
                }
            }
            // в ответ получаем номер блока
            blockIdBin := make([]byte, 4)
            _, err = conn.Read(blockIdBin)
            if err != nil {
                conn.Close()
                d.PrintSleep(err, 1)
                continue
            }
            conn.Close()
			id := utils.BinToDec(blockIdBin)
            if id > maxBlockId || i == 0 {
                maxBlockId = id
                maxBlockIdHost = hosts[i]["host"]
                maxBlockIdUserId = utils.StrToInt64(hosts[i]["user_id"])
			}
            if CheckDaemonsRestart() {
                utils.Sleep(1)
                break BEGIN
            }
		}

        // получим наш текущий имеющийся номер блока
        // ждем, пока разлочится и лочим сами, чтобы не попасть в тот момент, когда данные из блока уже занесены в БД, а info_block еще не успел обновиться
        err, restart = d.dbLock()
        if restart {
            break BEGIN
        }
        if err != nil {
            d.PrintSleep(err, 1)
            continue BEGIN
        }

        currentBlockId, err = d.Single("SELECT block_id FROM info_block").Int64()
        if err != nil {
            d.unlockPrintSleep(utils.ErrInfo(err), 1)
            continue
        }
		log.Info("currentBlockId", currentBlockId, "maxBlockId", maxBlockId)
        if maxBlockId <= currentBlockId {
            d.unlockPrintSleep(utils.ErrInfo(errors.New("maxBlockId <= currentBlockId")), 1)
            continue
		}

        // в цикле собираем блоки, пока не дойдем до максимального
        for blockId := currentBlockId+1; blockId < maxBlockId+1; blockId++ {
			d.UpdMainLock()
            if CheckDaemonsRestart() {
                d.unlockPrintSleep(utils.ErrInfo(err), 1)
                break BEGIN
            }
			variables, err := d.GetAllVariables()
            if err != nil {
                d.unlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
            }
            // качаем тело блока с хоста maxBlockIdHost
            binaryBlock, err := utils.GetBlockBody(maxBlockIdHost, blockId, dataTypeBlockBody, nodeHost)

            if len(binaryBlock) == 0 {
                // баним на 1 час хост, который дал нам пустой блок, хотя должен был дать все до максимального
                // для тестов убрал, потом вставить.
                //nodes_ban ($db, $max_block_id_user_id, substr($binary_block, 0, 512)."\n".__FILE__.', '.__LINE__.', '. __FUNCTION__.', '.__CLASS__.', '. __METHOD__);
                //p.NodesBan(maxBlockIdUserId, "len(binaryBlock) == 0")
                d.unlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
			}
            binaryBlockFull := binaryBlock
			utils.BytesShift(&binaryBlock, 1) // уберем 1-й байт - тип (блок/тр-я)
            // распарсим заголовок блока
            blockData := utils.ParseBlockHeader(&binaryBlock)
			log.Info("blockData: %v, blockId: %v", blockData, blockId)

            // если существуют глючная цепочка, тот тут мы её проигнорируем
            badBlocks_, err := d.Single("SELECT bad_blocks FROM config").Bytes()
            if err != nil {
                d.unlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
            }
            badBlocks := make(map[int64]string)
            if len(badBlocks_) > 0 {
                err = json.Unmarshal(badBlocks_, &badBlocks)
                if err != nil {
                    d.unlockPrintSleep(utils.ErrInfo(err), 1)
                    continue BEGIN
                }
            }
            if badBlocks[blockData.BlockId] == string(utils.BinToHex(blockData.Sign)) {
                d.NodesBan(maxBlockIdUserId, fmt.Sprintf("bad_block = %v => %v", blockData.BlockId, badBlocks[blockData.BlockId]))
                d.unlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
            }

            // размер блока не может быть более чем max_block_size
            if currentBlockId > 1 {
				if int64(len(binaryBlock)) > variables.Int64["max_block_size"] {
                    d.NodesBan(maxBlockIdUserId, fmt.Sprintf(`len(binaryBlock) > variables.Int64["max_block_size"]  %v > %v`, len(binaryBlock), variables.Int64["max_block_size"]))
                    d.unlockPrintSleep(utils.ErrInfo(err), 1)
                    continue BEGIN
				}
			}

			if blockData.BlockId != blockId {
                d.NodesBan(maxBlockIdUserId, fmt.Sprintf(`blockData.BlockId != blockId  %v > %v`, blockData.BlockId, blockId))
                d.unlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
			}

            // нам нужен хэш предыдущего блока, чтобы проверить подпись
            prevBlockHash := ""
            if blockId > 1 {
                prevBlockHash, err = d.Single("SELECT hash FROM block_chain WHERE id = ?", blockId-1).String()
                if err != nil {
                    d.unlockPrintSleep(utils.ErrInfo(err), 1)
                    continue BEGIN
				}
                prevBlockHash = string(utils.BinToHex([]byte(prevBlockHash)))
			} else {
                prevBlockHash = "0"
			}
			first := false
			if blockId == 1 {
				first = true
			}
            // нам нужен меркель-рут текущего блока
            mrklRoot, err := utils.GetMrklroot(binaryBlock, variables, first)
			if err != nil {
                d.NodesBan(maxBlockIdUserId, fmt.Sprintf(`%v`, err))
                d.unlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
			}

            // публичный ключ того, кто этот блок сгенерил
            nodePublicKey, err := d.GetNodePublicKey(blockData.UserId)
            if err != nil {
                d.unlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
            }

            // SIGN от 128 байта до 512 байт. Подпись от TYPE, BLOCK_ID, PREV_BLOCK_HASH, TIME, USER_ID, LEVEL, MRKL_ROOT
            forSign := fmt.Sprintf("0,%v,%v,%v,%v,%v,%s", blockData.BlockId, prevBlockHash, blockData.Time, blockData.UserId, blockData.Level, mrklRoot)

            // проверяем подпись
			if !first {
                _, err = utils.CheckSign([][]byte{nodePublicKey}, forSign, blockData.Sign, true);
			}

            // качаем предыдущие блоки до тех пор, пока отличается хэш предыдущего.
            // другими словами, пока подпись с prevBlockHash будет неверной, т.е. пока что-то есть в $error
			if err != nil {
                log.Info("%v", utils.ErrInfo(err))
				if blockId < 1 {
                    d.unlockPrintSleep(utils.ErrInfo(err), 1)
                    continue BEGIN
				}
                // нужно привести данные в нашей БД в соответствие с данными у того, у кого качаем более свежий блок
				//func (p *Parser) GetOldBlocks (userId, blockId int64, host string, hostUserId int64, goroutineName, getBlockScriptName, addNodeHost string) error {
                err := parser.GetOldBlocks(blockData.UserId, blockId-1, maxBlockIdHost, maxBlockIdUserId, GoroutineName, dataTypeBlockBody, nodeHost)
				log.Info("%v", err)
				if err != nil {
                    d.NodesBan(maxBlockIdUserId, fmt.Sprintf(`blockId: %v / %v`, blockId, err))
                    d.unlockPrintSleep(utils.ErrInfo(err), 1)
                    continue BEGIN
				}

            } else {

                log.Info("plug found blockId=%v\n", blockId)

                // получим наши транзакции в 1 бинарнике, просто для удобства
                var transactions []byte
                rows, err := d.Query("SELECT data FROM transactions WHERE verified = 1 AND used = 0")
                if err != nil {
                    d.unlockPrintSleep(utils.ErrInfo(err), 1)
                    continue BEGIN
				}
                defer rows.Close()
                for rows.Next() {
                    var data []byte
                    err = rows.Scan(&data)
                    if err != nil {
                        d.unlockPrintSleep(utils.ErrInfo(err), 1)
                        continue BEGIN
					}
                    transactions = append(transactions, utils.EncodeLengthPlusData(data)...)
                }
				if len(transactions) > 0 {
                    // отмечаем, что эти тр-ии теперь нужно проверять по новой
                    err = d.ExecSql("UPDATE transactions SET verified = 0 WHERE verified = 1 AND used = 0")
                    if err != nil {
                        d.unlockPrintSleep(utils.ErrInfo(err), 1)
                        continue BEGIN
					}
                    // откатываем по фронту все свежие тр-ии
                    parser.BinaryData = transactions
                    err = parser.ParseDataRollbackFront(false)
                    if err != nil {
                        utils.Sleep(1)
                        continue BEGIN
					}
                }

                err = parser.RollbackTransactionsTestblock(true)
                if err != nil {
                    d.unlockPrintSleep(utils.ErrInfo(err), 1)
                    continue BEGIN
                }
                err = d.ExecSql("DELETE FROM testblock")
                if err != nil {
                    d.unlockPrintSleep(utils.ErrInfo(err), 1)
                    continue BEGIN
				}
            }

            // теперь у нас в таблицах всё тоже самое, что у нода, у которого качаем блок
            // и можем этот блок проверить и занести в нашу БД
            parser.BinaryData = binaryBlockFull
            err = parser.ParseDataFull()
			if err == nil {
				parser.InsertIntoBlockchain()
			}
            // начинаем всё с начала уже с другими нодами. Но у нас уже могут быть новые блоки до $block_id, взятые от нода, которого с в итоге мы баним
            if err != nil {
                d.NodesBan(maxBlockIdUserId, fmt.Sprintf(`blockId: %v / %v`, blockId, err))
                d.unlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
			}
        }

		d.dbUnlock()

        // в sqllite данные в db-файл пишутся только после закрытия всех соединений с БД.
        for i:=0; i < 60; i++ {
            utils.Sleep(1)
            // проверим, не нужно ли нам выйти из цикла
            if CheckDaemonsRestart() {
                break BEGIN
            }
        }
    }
}

