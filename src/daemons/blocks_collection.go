package daemons

import (
//    "fmt"
    _ "github.com/lib/pq"
    //"time"
//    "database/sql"
	//"strconv"
   // "crypto/x509"
    //"encoding/pem"
    //"crypto"
    //"crypto/rand"
   // "crypto/rsa"
   // math_rand "math/rand"
    //"crypto/md5"
	//"bufio"
	//"os"
	//"errors"
	"utils"
	"log"
	"consts"
    //"io/ioutil"
 //   "log"
  //  "net/http"
 //   "io"
	////"os"
   // "github.com/alyu/configparser"
	//"io/ioutil"
    //"github.com/astaxie/beego/config"
	"fmt"
	"os"
	"dcparser"
//	"bufio"
//    "io/ioutil"
	"static"
)


func BlocksCollection(configIni map[string]string) {

    const GoroutineName = "blocks_collection"

    db := utils.DbConnect(configIni)

    // Возможна ситуация, когда инсталяция еще не завершена. База данных может быть создана, а таблицы еще не занесены
    INSTALL:
    progress, err := db.Single("SELECT progress FROM install").String()
    if err != nil || progress != "complete" {
        utils.Sleep(1)
        goto INSTALL
    }

	var cur bool
    BEGIN:
	for {

		fmt.Println("BlocksCollection")
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

        myPrefix, err:= db.GetMyPrefix()
        if err != nil {
            log.Print(utils.ErrInfo(err))
            utils.Sleep(1)
            continue BEGIN
        }
		fmt.Println("myPrefix",myPrefix)

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

        fmt.Println("config", config)
        fmt.Println("currentBlockId", currentBlockId)

		// на время тесто
		if !cur {
            currentBlockId = 0
            cur = true
        }
        if currentBlockId==0 {

			if config["first_load_blockchain"]=="file" {

                fmt.Println("first_load_blockchain=file")
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

                fmt.Println("GO!")
                for {
                    b1 := make([]byte, 5)
                    file.Read(b1)
                    dataSize := utils.BinToDec(b1)
                    fmt.Println("dataSize", dataSize)
                    if dataSize > 0 {

                        data := make([]byte, dataSize)
                        file.Read(data)
                        fmt.Printf("data %x\n", data)
                        blockId := utils.BinToDec(data[0:5])
                        if blockId == 138500 {
                           break BEGIN
                        }
                        fmt.Println("blockId", blockId)
                        data2:=data[5:]
                        length := utils.DecodeLength(&data2)
                        fmt.Println("length", length)
                        fmt.Printf("data2 %x\n", data2)
                        blockBin := utils.BytesShift(&data2, length)
                        fmt.Printf("blockBin %x\n", blockBin)

                        if blockId > 134999 {

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
		}

		db.DbUnlock(GoroutineName);

        // в sqllite данные в db-файл пишутся только после закрытия всех соединений с БД.
        db.Close()
        db = utils.DbConnect(configIni)

        utils.Sleep(3)
		//break
    }
}

