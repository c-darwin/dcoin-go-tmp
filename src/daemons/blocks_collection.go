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
)


func BlocksCollection(configIni map[string]string) {

    const GoroutineName = "blocks_collection"

    db := utils.DbConnect(configIni)

    // Возможна ситуация, когда инсталяция еще не завершена. База данных может быть создана, а таблицы еще не занесены
    INSTALL:
    progress, err := db.Single("SELECT progress FROM install")
    if err != nil || progress != "complete" {
        utils.Sleep(1)
        goto INSTALL
    }

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
		fmt.Println(myPrefix)

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
        if currentBlockId==0 {
			if config["first_load_blockchain"]=="file" {
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
                }
			}

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
                    fmt.Println("blockId", blockId)
					data2:=data[5:]
                    length := utils.DecodeLength(&data2)
                    fmt.Println("length", length)
                    fmt.Printf("data2 %x\n", data2)
					blockBin := utils.BytesShift(&data2, length)
                    fmt.Printf("blockDin %x\n", blockBin)

					// парсинг блока
                    parser := new(dcparser.Parser)
                    parser.DCDB = db
                    parser.BinaryData = blockBin;
					parser.GoroutineName = GoroutineName

                    err = parser.ParseDataFull()
                    if err != nil {
                        log.Print(utils.ErrInfo(err))
                        utils.Sleep(1)
                        db.DbUnlock(GoroutineName);
                        file.Close()
                        continue BEGIN
					}
                    // ненужный тут размер в конце блока данных
                    data = make([]byte, 5)
                    file.Read(data)
                }
				/*
                // читаем побайтно
                buf := bufio.NewReader(file)
                b5, err := buf.Peek(5)
                if err != nil {
                    log.Print(utils.ErrInfo(err))
                    utils.Sleep(1)
                    db.DbUnlock(GoroutineName);
                    continue BEGIN
                }
                //file.Close()
                dataSize := utils.BinToDec(b5)
                fmt.Println(dataSize)
                if dataSize > 0 {
                    //buf := bufio.NewReader(file)
                    data, err := buf.Peek(int(dataSize))
                    if err != nil {
                        log.Print(utils.ErrInfo(err))
                        utils.Sleep(1)
                        db.DbUnlock(GoroutineName);
                        continue BEGIN
                    }
                    file.Close()
					fmt.Printf("%x", data)
				} else {
                    file.Close()
					break
				}*/
                utils.Sleep(1)
            }
            file.Close()

	    }

		db.DbUnlock(GoroutineName);

        // в sqllite данные в db-файл пишутся только после закрытия всех соединений с БД.
        db.Close()
        db = utils.DbConnect(configIni)

        utils.Sleep(3)
		//break
    }
}

