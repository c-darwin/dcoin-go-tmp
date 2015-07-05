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
        // отметимся в БД, что мы живы.
        db.UpdDaemonTime(GoroutineName)
        // проверим, не нужно нам выйти из цикла
        if utils.CheckDaemonRestart("BlocksCollection") {
			break
		}

        config, err := db.GetNodeConfig()
        if err != nil {
            log.Print(utils.ErrInfo(err))
            utils.Sleep(1)
            continue BEGIN
        }

        /*myPrefix, err:= db.GetMyPrefix()
        if err != nil {
            log.Print(utils.ErrInfo(err))
            utils.Sleep(1)
            continue BEGIN
        }
		log.Println("myPrefix",myPrefix)*/

       err = db.DbLock(GoroutineName);
        if err != nil {
            log.Print(utils.ErrInfo(err))
            utils.Sleep(1)
            continue BEGIN
        }

        // если это первый запуск во время инсталяции
        currentBlockId, err := db.GetBlockId()
        if err != nil {
            log.Print(utils.ErrInfo(err))
            utils.Sleep(1)
            db.DbUnlock(GoroutineName);
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
                /*
                На время тестов не какчаем
                blockchainSize, err := utils.DownloadToFile(consts.BLOCKCHAIN_URL, "public/blockchain")
                if err != nil || blockchainSize < consts.BLOCKCHAIN_SIZE {
                    if err != nil {
                        log.Print(utils.ErrInfo(err))
                    } else {
                        log.Print(fmt.Sprintf("%v < %v", blockchainSize, consts.BLOCKCHAIN_SIZE))
					}
                    utils.Sleep(1)
                    db.DbUnlock(GoroutineName);
                    continue BEGIN
                }*/

                first := true
                // блокчейн мог быть загружен ранее. проверим его размер
                file, err := os.Open("public/blockchain")
                if err != nil {
                    log.Print(utils.ErrInfo(err))
                    utils.Sleep(1)
                    db.DbUnlock(GoroutineName);
                    continue BEGIN
                }

                stat, err := file.Stat()
                if err != nil {
                    log.Print(utils.ErrInfo(err))
                    utils.Sleep(1)
                    db.DbUnlock(GoroutineName);
                    file.Close()
                    continue BEGIN
                }
                if stat.Size() < consts.BLOCKCHAIN_SIZE {
                    log.Print(fmt.Sprintf("%v < %v", stat.Size(), consts.BLOCKCHAIN_SIZE))
                    utils.Sleep(1)
                    db.DbUnlock(GoroutineName);
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
                                log.Print(utils.ErrInfo(err))
                                utils.Sleep(1)
                                db.DbUnlock(GoroutineName);
                                file.Close()
                                break
                            }
                            parser.InsertIntoBlockchain()

                            // отметимся в БД, что мы живы.
                            parser.DCDB.UpdDaemonTime(GoroutineName)
                            // отметимся, чтобы не спровоцировать очистку таблиц
                            err = parser.DCDB.UpdMainLock()
                            if err != nil {
                                log.Print(utils.ErrInfo(err))
                                utils.Sleep(1)
                                db.DbUnlock(GoroutineName);
                                file.Close()
                                break
                            }
                            if utils.CheckDaemonRestart(GoroutineName) {
                                log.Print(utils.ErrInfo(err))
                                utils.Sleep(1)
                                db.DbUnlock(GoroutineName);
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
                parser := new(dcparser.Parser)
                parser.DCDB = db
                parser.GoroutineName = GoroutineName
                parser.BinaryData = newBlock
                if err != nil {
                    log.Print(utils.ErrInfo(err))
                    utils.Sleep(1)
                    db.DbUnlock(GoroutineName);
                    break
                }
                parser.CurrentVersion = consts.VERSION

                err = parser.ParseDataFull()
                if err != nil {
                    log.Print(utils.ErrInfo(err))
                    utils.Sleep(1)
                    db.DbUnlock(GoroutineName);
                    break
                }
                parser.InsertIntoBlockchain()
			}

			utils.Sleep(1)
			db.DbUnlock(GoroutineName);
            continue BEGIN
		}
		db.DbUnlock(GoroutineName);

        myConfig, err := db.OneRow("SELECT local_gate_ip, static_node_user_id FROM config").String()
        if err != nil {
			db.UnlockPrintSleep(utils.ErrInfo(err), 1)
            continue
        }
		var hosts []map[string]string
		var getMaxBlockScriptName, getBlockScriptName, addNodeHost string
		if len(myConfig["local_gate_ip"]) > 0 {
            hosts = append(hosts, map[string]string{"host": myConfig["local_gate_ip"], "user_id": myConfig["static_node_user_id"]})
            nodeHost, err := db.Single("SELECT host FROM miners_data WHERE user_id  =  ?", myConfig["static_node_user_id"]).String()
            if err != nil {
                db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                continue
			}
            getMaxBlockScriptName = "ajax?controllerName=protectedGetMaxBlock&nodeHost=".nodeHost;
			getBlockScriptName = "ajax?controllerName=protectedGetBlock";
			addNodeHost = "&nodeHost=".nodeHost;
        } else {
            // получим список нодов, с кем установлено рукопожатие
            hosts, err := db.GetAll("SELECT * FROM nodes_connection", -1)
            if err != nil {
                db.UnlockPrintSleep(utils.ErrInfo(err), 1)
                continue
			}
            getMaxBlockScriptName = "ajax?controllerName=getMaxBlock";
            getBlockScriptName = "ajax?controllerName=getBlock";
            addNodeHost = "";
        }

		log.Println(hosts)

		if len(hosts) == 0 {
            db.UnlockPrintSleep(utils.ErrInfo(errors.New("len(hosts) == 0")), 1)
            continue
		}

        maxBlockId := 1
        // получим максимальный номер блока
        for _, host := range hosts {
            url := host["host"]+"/"+getMaxBlockScriptName
            resp, err := http.Get(url)
            if err != nil {
                log.Print(utils.ErrInfo(err))
                continue
            }
            defer resp.Body.Close()
            html, err := ioutil.ReadAll(resp.Body)
			id = utils.StrToInt()
		}


        // в sqllite данные в db-файл пишутся только после закрытия всех соединений с БД.
        db.Close()
        db = utils.DbConnect(configIni)

        utils.Sleep(3)
		//break

		//p := new(dcparser.Parser)
        //p.GetBlocks()

    }
}

