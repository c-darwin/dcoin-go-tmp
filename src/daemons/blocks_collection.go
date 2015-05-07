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
	"bufio"
)

func BlocksCollection(configIni map[string]string) {

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
        utils.UpdDaemonTime("BlocksCollection")
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
        userPublicKey, err:= db.GetMyPublicKey(myPrefix)
        if err != nil {
            log.Print(utils.ErrInfo(err))
            utils.Sleep(1)
            continue BEGIN
        }

        // если это первый запуск во время инсталяции, то нужно дождаться, пока юзер загрузит свой ключ
        // т.к. возможно этот ключ уже есть в блоках и нужно обновить внутренние таблицы
        if (len(userPublicKey)==0 && config["first_load_blockchain"]!="file" && config["first_load_blockchain"]!="nodes") {
            log.Print("continue")
            utils.Sleep(1)
            continue BEGIN
        }

        db.DbLock();

        // если это первый запуск во время инсталяции
        currentBlockId, err := db.GetBlockId()
        if err != nil {
            log.Print(utils.ErrInfo(err))
            utils.Sleep(1)
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
                    continue BEGIN
                }
			}

			// блокчейн мог быть загружен ранее. проверим его размер
            file, err := os.Open("public/blockchain")
            if err != nil {
                log.Print(utils.ErrInfo(err))
                utils.Sleep(1)
                continue BEGIN
            }

            stat, err := file.Stat()
            if err != nil {
                log.Print(utils.ErrInfo(err))
                utils.Sleep(1)
                continue BEGIN
            }
			if stat.Size() < consts.BLOCKCHAIN_SIZE {
                log.Print(fmt.Sprintf("%v < %v", stat.Size(), consts.BLOCKCHAIN_SIZE))
                utils.Sleep(1)
                continue BEGIN
			}

			// читаем побайтно
            buf := bufio.NewReader(file)
            b4, err := buf.Peek(5)
            if err != nil {
                log.Print(utils.ErrInfo(err))
                utils.Sleep(1)
                continue BEGIN
            }
            fmt.Printf("5 bytes: %s\n", string(b4))

            file.Close()


	    }

		db.DbUnlock();

        // в sqllite данные в db-файл пишутся только после закрытия всех соединений с БД.
        db.Close()
        db = utils.DbConnect(configIni)

        utils.Sleep(3)
		break
    }
}

