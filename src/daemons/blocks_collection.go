package daemons

import (
    _ "github.com/lib/pq"
	"utils"
	"log"
	"consts"
	"fmt"
	"os"
	"dcparser"
	"static"
	"errors"
	"net/http"
	"io/ioutil"
	"encoding/json"
)

func BlocksCollection(configIni map[string]string) {

    const GoroutineName = "blocks_collection"

    db := utils.DbConnect(configIni)
    db.GoroutineName = GoroutineName

    // Возможна ситуация, когда инсталяция еще не завершена. База данных может быть создана, а таблицы еще не занесены
    INSTALL:
    progress, err := db.Single("SELECT progress FROM install").String()
    if err != nil || progress != "complete" {
		log.Println(`progress != "complete"`)
		if err!=nil {
            log.Print(utils.ErrInfo(err))
        }
        utils.Sleep(1)
        goto INSTALL
    }

	var cur bool
    BEGIN:
	for {

		log.Println("BlocksCollection")
        // проверим, не нужно нам выйти из цикла
        if db.CheckDaemonRestart() {
			break
		}

        config, err := db.GetNodeConfig()
        if err != nil {
            db.PrintSleep(err, 1)
            continue BEGIN
        }

        /*myPrefix, err:= db.GetMyPrefix()
        if err != nil {
            log.Print(utils.ErrInfo(err))
            utils.Sleep(1)
            continue BEGIN
        }
		log.Println("myPrefix",myPrefix)*/

       err = db.DbLock();
        if err != nil {
            db.PrintSleep(err, 1)
            continue BEGIN
        }

        // если это первый запуск во время инсталяции
        currentBlockId, err := db.GetBlockId()
        if err != nil {
            db.UnlockPrintSleep(err, 1)
            continue BEGIN
        }

        log.Println("config", config)
        log.Println("currentBlockId", currentBlockId)

		// на время тесто
		if !cur {
            currentBlockId = 0
            cur = true
        }
        if currentBlockId==0 {

			if config["first_load_blockchain"]=="file" {

                log.Println("first_load_blockchain=file")

                blockchainSize, err := utils.DownloadToFile(consts.BLOCKCHAIN_URL, "public/blockchain")
                if err != nil || blockchainSize < consts.BLOCKCHAIN_SIZE {
                    if err != nil {
                        log.Print(utils.ErrInfo(err))
                    } else {
                        log.Print(fmt.Sprintf("%v < %v", blockchainSize, consts.BLOCKCHAIN_SIZE))
					}
                    db.UnlockPrintSleep(err, 1)
                    continue BEGIN
                }

                first := true
                // блокчейн мог быть загружен ранее. проверим его размер
                file, err := os.Open("public/blockchain")
                if err != nil {
                    db.UnlockPrintSleep(err, 1)
                    continue BEGIN
                }

                stat, err := file.Stat()
                if err != nil {
                    db.UnlockPrintSleep(err, 1)
                    file.Close()
                    continue BEGIN
                }
                if stat.Size() < consts.BLOCKCHAIN_SIZE {
                    db.UnlockPrintSleep(fmt.Errorf("%v < %v", stat.Size(), consts.BLOCKCHAIN_SIZE), 1)
                    file.Close()
                    continue BEGIN
                }

                log.Println("GO!")
                for {
                    b1 := make([]byte, 5)
                    file.Read(b1)
                    dataSize := utils.BinToDec(b1)
                    log.Println("dataSize", dataSize)
                    if dataSize > 0 {

                        data := make([]byte, dataSize)
                        file.Read(data)
                        fmt.Printf("data %x\n", data)
                        blockId := utils.BinToDec(data[0:5])
                        if blockId == 244790 {
                            break BEGIN
                        }
                        log.Println("blockId", blockId)
                        data2:=data[5:]
                        length := utils.DecodeLength(&data2)
                        log.Println("length", length)
                        fmt.Printf("data2 %x\n", data2)
                        blockBin := utils.BytesShift(&data2, length)
                        fmt.Printf("blockBin %x\n", blockBin)

                        if blockId > 244790 {

                            // парсинг блока
                            parser := new(dcparser.Parser)
                            parser.DCDB = db
                            parser.BinaryData = blockBin;
                            parser.GoroutineName = GoroutineName

                            if first {
                                parser.CurrentVersion = consts.VERSION
                                first = false
                            }
                            err = parser.ParseDataFull()
                            if err != nil {
                                db.UnlockPrintSleep(err, 1)
                                file.Close()
                                break
                            }
                            parser.InsertIntoBlockchain()

                            // отметимся, чтобы не спровоцировать очистку таблиц
                            err = parser.DCDB.UpdMainLock()
                            if err != nil {
                                db.UnlockPrintSleep(err, 1)
                                file.Close()
                                break
                            }
                            if db.CheckDaemonRestart() {
                                db.UnlockPrintSleep(err, 1)
                                file.Close()
                                break BEGIN
                            }
                        }
                        // ненужный тут размер в конце блока данных
                        data = make([]byte, 5)
                        file.Read(data)
                    } else {
						break
					}
                   // utils.Sleep(1)
                }
                file.Close()
	        } else {
                newBlock, err := static.Asset("static/1block.bin")
                if err != nil {
                    db.UnlockPrintSleep(err, 1)
                    break
                }
                parser := new(dcparser.Parser)
                parser.DCDB = db
                parser.GoroutineName = GoroutineName
                parser.BinaryData = newBlock
                parser.CurrentVersion = consts.VERSION

                err = parser.ParseDataFull()
                if err != nil {
                    db.UnlockPrintSleep(err, 1)
                    break
                }
                err = parser.InsertIntoBlockchain()

                if err != nil {
                    db.UnlockPrintSleep(err, 1)
                    break
                }
			}

			utils.Sleep(1)
			db.DbUnlock();
            continue BEGIN
		}
		db.DbUnlock();

        myConfig, err := db.OneRow("SELECT local_gate_ip, static_node_user_id FROM config").String()
        if err != nil {
            utils.Sleep(1)
            continue
        }
		var hosts []map[string]string
		var getMaxBlockScriptName, getBlockScriptName, addNodeHost string
		if len(myConfig["local_gate_ip"]) > 0 {
            hosts = append(hosts, map[string]string{"host": myConfig["local_gate_ip"], "user_id": myConfig["static_node_user_id"]})
            nodeHost, err := db.Single("SELECT host FROM miners_data WHERE user_id  =  ?", myConfig["static_node_user_id"]).String()
            if err != nil {
                utils.Sleep(1)
                continue
			}
            getMaxBlockScriptName = "ajax?controllerName=protectedGetMaxBlock&nodeHost="+nodeHost;
			getBlockScriptName = "ajax?controllerName=protectedGetBlock";
			addNodeHost = "&nodeHost="+nodeHost;
        } else {
            // получим список нодов, с кем установлено рукопожатие
            hosts, err = db.GetAll("SELECT * FROM nodes_connection", -1)
            if err != nil {
                utils.Sleep(1)
                continue
			}
            getMaxBlockScriptName = "ajax?controllerName=getMaxBlock";
            getBlockScriptName = "ajax?controllerName=getBlock";
            addNodeHost = "";
        }

		log.Println(hosts)

		if len(hosts) == 0 {
            utils.Sleep(1)
            continue
		}

        maxBlockId := int64(1)
        maxBlockIdHost := ""
        var maxBlockIdUserId int64
        // получим максимальный номер блока
        for i:=0; i < len(hosts); i++ {
            url := hosts[i]["host"]+"/"+getMaxBlockScriptName
            resp, err := http.Get(url)
            if err != nil {
                utils.Sleep(1)
                continue
            }
            defer resp.Body.Close()
            html, err := ioutil.ReadAll(resp.Body)
			id := utils.BytesToInt64(html)
            if id > maxBlockId || i == 0 {
                maxBlockId = id
                maxBlockIdHost = hosts[i]["host"]
                maxBlockIdUserId = utils.StrToInt64(hosts[i]["user_id"])
			}
            if db.CheckDaemonRestart() {
                utils.Sleep(1)
                break BEGIN
            }
		}

        // получим наш текущий имеющийся номер блока
        // ждем, пока разлочится и лочим сами, чтобы не попасть в тот момент, когда данные из блока уже занесены в БД, а info_block еще не успел обновиться
		db.DbLock()
        currentBlockId, err = db.Single("SELECT block_id FROM info_block").Int64()
        if err != nil {
            db.UnlockPrintSleep(utils.ErrInfo(err), 1)
            continue
        }
		log.Println("currentBlockId", currentBlockId, "maxBlockId", maxBlockId)
        if maxBlockId <= currentBlockId {
            db.UnlockPrintSleep(utils.ErrInfo(errors.New("maxBlockId <= currentBlockId")), 1)
            continue
		}

        // в цикле собираем блоки, пока не дойдем до максимального
        for blockId := currentBlockId+1; blockId < maxBlockId+1; blockId++ {
			db.UpdMainLock()
            if db.CheckDaemonRestart() {
                db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                break BEGIN
            }
			variables, err := db.GetAllVariables()
            if err != nil {
                db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
            }
            url := maxBlockIdHost+"/"+getBlockScriptName+"?id="+utils.Int64ToStr(blockId)+addNodeHost
            resp, err := http.Get(url)
            if err != nil {
                db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
            }
            defer resp.Body.Close()
            binaryBlock, err := ioutil.ReadAll(resp.Body)
			if len(binaryBlock) == 0 {
                // баним на 1 час хост, который дал нам пустой блок, хотя должен был дать все до максимального
                // для тестов убрал, потом вставить.
                //nodes_ban ($db, $max_block_id_user_id, substr($binary_block, 0, 512)."\n".__FILE__.', '.__LINE__.', '. __FUNCTION__.', '.__CLASS__.', '. __METHOD__);
                //p.NodesBan(maxBlockIdUserId, "len(binaryBlock) == 0")
                db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
			}
            binaryBlockFull := binaryBlock
			utils.BytesShift(&binaryBlock, 1) // уберем 1-й байт - тип (блок/тр-я)
            // распарсим заголовок блока
            blockData := utils.ParseBlockHeader(&binaryBlock)
			log.Println("blockData", blockData, "blockId", blockId)

            // если существуют глючная цепочка, тот тут мы её проигнорируем
            badBlocks_, err := db.Single("SELECT bad_blocks FROM config").Bytes()
            if err != nil {
                db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
            }
            badBlocks := make(map[int64]string)
            err = json.Unmarshal(badBlocks_, &badBlocks)
            if err != nil {
                db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
            }
            if badBlocks[blockData.BlockId] == string(utils.BinToHex(blockData.Sign)) {
                db.NodesBan(maxBlockIdUserId, fmt.Sprintf("bad_block = %v => %v", blockData.BlockId, badBlocks[blockData.BlockId]))
                db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
            }

            // размер блока не может быть более чем max_block_size
            if currentBlockId > 1 {
				if int64(len(binaryBlock)) > variables.Int64["max_block_size"] {
                    db.NodesBan(maxBlockIdUserId, fmt.Sprintf(`len(binaryBlock) > variables.Int64["max_block_size"]  %v > %v`, len(binaryBlock), variables.Int64["max_block_size"]))
                    db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                    continue BEGIN
				}
			}

			if blockData.BlockId != blockId {
                db.NodesBan(maxBlockIdUserId, fmt.Sprintf(`blockData.BlockId != blockId  %v > %v`, blockData.BlockId, blockId))
                db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
			}

            // нам нужен хэш предыдущего блока, чтобы проверить подпись
            prevBlockHash := ""
            if blockId > 1 {
                prevBlockHash, err := db.Single("SELECT hash FROM block_chain WHERE id = ?", blockId-1).String()
                if err != nil {
                    db.UnlockPrintSleep(utils.ErrInfo(err), 1)
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
                db.NodesBan(maxBlockIdUserId, fmt.Sprintf(`%v`, err))
                db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
			}

            // публичный ключ того, кто этот блок сгенерил
            nodePublicKey, err := db.GetNodePublicKey(blockData.UserId)
            if err != nil {
                db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
            }

            // SIGN от 128 байта до 512 байт. Подпись от TYPE, BLOCK_ID, PREV_BLOCK_HASH, TIME, USER_ID, LEVEL, MRKL_ROOT
            forSign := fmt.Sprintf("0,%v,%v,%v,%v,%v,%v", blockData.BlockId, prevBlockHash, blockData.Time, blockData.UserId, blockData.Level, mrklRoot)

            // проверяем подпись
			if !first {
                _, err = utils.CheckSign([][]byte{nodePublicKey}, forSign, blockData.Sign, true);
			}

            // качаем предыдущие блоки до тех пор, пока отличается хэш предыдущего.
            // другими словами, пока подпись с $prev_block_hash будет неверной, т.е. пока что-то есть в $error
			if err != nil {
                log.Println(utils.ErrInfo(err))
				if blockId < 1 {
                    db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                    continue BEGIN
				}
                // нужно привести данные в нашей БД в соответствие с данными у того, у кого качаем более свежий блок
                p := new(dcparser.Parser)
				p.DCDB = db
				//func (p *Parser) GetOldBlocks (userId, blockId int64, host string, hostUserId int64, goroutineName, getBlockScriptName, addNodeHost string) error {
                err := p.GetOldBlocks(blockData.UserId, blockId-1, maxBlockIdHost, maxBlockIdUserId, GoroutineName, getBlockScriptName, addNodeHost)
				log.Println(err)
				if err != nil {
                    db.NodesBan(maxBlockIdUserId, fmt.Sprintf(`blockId: %v / %v`, blockId, err))
                    db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                    continue BEGIN
				}

            } else {

                log.Printf("plug found blockId=%v\n", blockId)

                // получим наши транзакции в 1 бинарнике, просто для удобства
                var transactions []byte
                rows, err := db.Query("SELECT data FROM transactions WHERE verified = 1 AND used = 0")
                if err != nil {
                    db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                    continue BEGIN
				}
                defer rows.Close()
                for rows.Next() {
                    var data []byte
                    err = rows.Scan(&data)
                    if err != nil {
                        db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                        continue BEGIN
					}
                    transactions = append(transactions, utils.EncodeLengthPlusData(data)...)
                }
				if len(transactions) > 0 {
                    // отмечаем, что эти тр-ии теперь нужно проверять по новой
                    err = db.ExecSql("UPDATE transactions SET verified = 0 WHERE verified = 1 AND used = 0")
                    if err != nil {
                        db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                        continue BEGIN
					}
                    // откатываем по фронту все свежие тр-ии
                    parser := new(dcparser.Parser)
                    parser.DCDB = db
                    parser.BinaryData = transactions
                    err = parser.ParseDataRollbackFront(false)
                    if err != nil {
                        utils.Sleep(1)
                        continue BEGIN
					}
                }

                parser := new(dcparser.Parser)
                parser.DCDB = db
                err = parser.RollbackTransactionsTestblock(true)
                if err != nil {
                    utils.Sleep(1)
                    continue BEGIN
                }
                err = db.ExecSql("TRUNCATE TABLE testblock")
                if err != nil {
                    db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                    continue BEGIN
				}
            }

            // теперь у нас в таблицах всё тоже самое, что у нода, у которого качаем блок
            // и можем этот блок проверить и занести в нашу БД
            parser := new(dcparser.Parser)
            parser.DCDB = db
            parser.BinaryData = binaryBlockFull
            err = parser.ParseDataFull()
			if err == nil {
				parser.InsertIntoBlockchain()
			}
            // начинаем всё с начала уже с другими нодами. Но у нас уже могут быть новые блоки до $block_id, взятые от нода, которого с в итоге мы баним
            if err != nil {
                db.NodesBan(maxBlockIdUserId, fmt.Sprintf(`blockId: %v / %v`, blockId, err))
                db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                continue BEGIN
			}
        }

		db.DbUnlock()

        // в sqllite данные в db-файл пишутся только после закрытия всех соединений с БД.
        db.Close()
        db = utils.DbConnect(configIni)

        utils.Sleep(60)
    }
}

