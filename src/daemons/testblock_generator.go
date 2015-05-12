package daemons

import (
    "fmt"
    _ "github.com/lib/pq"
    "time"
//    "database/sql"
	"strconv"
    "crypto/x509"
    "encoding/pem"
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    math_rand "math/rand"
    "crypto/md5"
	"bufio"
	"os"
	"errors"
	"utils"
	"log"
   // "github.com/alyu/configparser"
	//"io/ioutil"
    //"github.com/astaxie/beego/config"
)


var err error

func Testblock_generator(configIni map[string]string) {

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

        db.DbLock()

        blockId, err := db.GetBlockId()
		if err != nil {
            db.DbUnlock(mainName)
			log.Print(err)
            utils.Sleep(1)
            continue BEGIN
        }
        newBlockId := blockId + 1;
        fmt.Println(newBlockId, "newBlockId")
        testBlockId, err := db.GetTestBlockId()
        if err != nil {
            db.DbUnlock(mainName)
            log.Print(err)
            utils.Sleep(1)
            continue BEGIN
        }
        fmt.Println(testBlockId, "testBlockId")

        if x, err := db.GetMyLocalGateIp(); x!="" {
            if err != nil {
                log.Print(err)
			}
            log.Print("continue")
            db.DbUnlock(mainName)
            utils.Sleep(1)
            continue BEGIN
        }

        if testBlockId==newBlockId {
            db.DbUnlock(mainName)
            log.Print(err)
            utils.Sleep(1)
            continue
        }

        prevBlock, myUserId, myMinerId, currentUserId, level, levelsRange, err := db.TestBlock()
        if err != nil {
            db.DbUnlock(mainName)
            log.Print(err)
            utils.Sleep(1)
            continue BEGIN
        }
        fmt.Println(prevBlock, myUserId, myMinerId, currentUserId, level, levelsRange)

		if myMinerId==0 {
            db.DbUnlock(mainName)
            utils.Sleep(1)
			continue
		}

		sleep := db.GetGenSleep(prevBlock, level)
        fmt.Println("sleep", sleep)

        blockId = prevBlock.BlockId;
        fmt.Println("blockId", blockId)
        prevHeadHash := prevBlock.HeadHash;
        fmt.Println("prevHeadHash", prevHeadHash)

        // сколько прошло сек с момента генерации прошлого блока
        diff := time.Now().Unix() - prevBlock.Time;
        fmt.Println("diff", diff)

        // вычитаем уже прошедшее время
        utils.SleepDiff(&sleep, diff)

        // Если случится откат или придет новый блок, то генератор блоков нужно запускать с начала, т.к. изменится max_miner_id.
        fmt.Println("sleep", sleep)
        startSleep := time.Now().Unix();
        fmt.Println("startSleep", startSleep)

        db.DbUnlock(mainName)

        for i := 0; i < int(sleep); i++ {
            db.DbLock();
            fmt.Println("i", i)
            fmt.Println("sleep", sleep)
			var newHeadHash string
            err := db.QueryRow("SELECT LOWER(encode(head_hash, 'hex')) FROM info_block").Scan(&newHeadHash)
            utils.CheckErr(err)
            fmt.Println("newHeadHash", newHeadHash)
            db.DbUnlock(mainName);
            if (newHeadHash != prevHeadHash) {
                fmt.Println("newHeadHash!=prevHeadHash", newHeadHash, prevHeadHash)
                continue BEGIN
            }
            // из-за задержек с main_lock время уже прошло и выходим раньше, чем закончится цикл
            if time.Now().Unix() - startSleep > sleep {
                fmt.Println("break")
                break
            }
            time.Sleep(1000 * time.Millisecond) // спим 1 сек. общее время = $sleep
        }

        /*
		 *  Закончили спать, теперь генерим блок
		 * Но, всё, что было до main_unlock может стать недействительным, т.е. надо обновить данные
		 * */
        db.DbLock();

        prevBlock, myUserId, myMinerId, currentUserId, level, levelsRange, err = db.TestBlock();
		if err != nil {
			log.Print(err)
            db.DbUnlock(mainName)
            utils.Sleep(1)
            continue
		}
        fmt.Println(prevBlock, myUserId, myMinerId, currentUserId, level, levelsRange)
        // сколько прошло сек с момента генерации прошлого блока
        diff = time.Now().Unix() - prevBlock.Time;
        fmt.Println("diff", diff)
        // вычитаем уже прошедшее время
        utils.SleepDiff(&sleep, diff)
        fmt.Println("sleep", sleep)
        // если нужно доспать, то просто вернемся в начало и доспим нужное время. И на всякий случай убедимся, что блок не изменился
        if sleep > 0 || prevBlock.HeadHash != prevHeadHash {
            fmt.Println("continue")
            db.DbUnlock(mainName)
            time.Sleep(1000 * time.Millisecond)
            continue
        }
        fmt.Println("Закончили спать, теперь генерим блок")
        blockId = prevBlock.BlockId;
        if blockId < 1 {
            fmt.Println("continue")
            db.DbUnlock(mainName)
            time.Sleep(1000 * time.Millisecond)
            continue
        }

        newBlockId = blockId + 1;
		var myPrefix string
        CommunityUser, err := db.GetCommunityUsers()
        if len(CommunityUser)>0 {
            myPrefix = strconv.FormatInt(myUserId, 10)+"_"
		} else {
            myPrefix = ""
	    }
        nodePrivateKey := db.GetNodePrivateKey(myPrefix)
		if len(nodePrivateKey) < 1 {
            fmt.Println("continue")
            db.DbUnlock(mainName)
            time.Sleep(1000 * time.Millisecond)
            continue
        }
        prevHeadHash = prevBlock.HeadHash

        fmt.Println(nodePrivateKey, prevHeadHash)

        //#####################################
        //##		 Формируем блок
        //#####################################
        fmt.Println(newBlockId, currentUserId)
        if currentUserId < 1 {
            fmt.Println("continue")
            db.DbUnlock(mainName)
            time.Sleep(1000 * time.Millisecond)
            continue
        }
        if prevBlock.BlockId >= newBlockId {
            fmt.Println("continue")
            db.DbUnlock(mainName)
            time.Sleep(1000 * time.Millisecond)
            continue
        }
        // отакатим transactions_testblock
        utils.RollbackTransactionsTestBlock(true)

       Time := time.Now().Unix()

        // переведем тр-ии в `verified` = 1
        utils.AllTxParser();

        var mrklArray  []string
		var usedTransactions string
		var mrklRoot string
        // берем все данные из очереди. Они уже были проверены ранее, и можно их не проверять, а просто брать
        rows, err := db.Query("SELECT data, LOWER(encode(hash, 'hex')),type,user_id,third_var FROM transactions WHERE used=0 AND verified = 1")
        utils.CheckErr(err)
        for rows.Next() {
            var data []byte
            var hash string
            var txType string
            var txUserId string
            var thirdVar string
            err = rows.Scan(&data, &hash, &txType, &txUserId, &thirdVar)
            utils.CheckErr(err)
            fmt.Println("data", data)
            fmt.Println("hash", hash)
            transactionType := data[1:2];
            fmt.Println(transactionType)
            fmt.Printf("%x", transactionType)
            mrklArray = append(mrklArray, utils.DSha256(data));
            fmt.Println("mrklArray", mrklArray)

            hash2_ := md5.New()
            hash2_.Write(data)
            hashMd5:=fmt.Sprintf("%x", hash2_.Sum(nil))
            fmt.Println(hashMd5)

            dataHex := fmt.Sprintf("%x", data)
            fmt.Println("dataHex", dataHex)

            file, _ := os.Create("/home/z/psql.sql")
            writer := bufio.NewWriter(file)
            writer.Write([]byte("\\x"+hashMd5))
            writer.Write([]byte("|"))
            writer.Write([]byte("\\x"+dataHex))
            writer.Write([]byte("|"))
            writer.Write([]byte(txType))
            writer.Write([]byte("|"))
            writer.Write([]byte(txUserId))
            writer.Write([]byte("|"))
            writer.Write([]byte(thirdVar))
            writer.Flush()
            res, err := db.Exec("DELETE FROM transactions WHERE hash = decode($1, 'hex')", hashMd5)
            utils.CheckErr(err)
            affect, err := res.RowsAffected()
            utils.CheckErr(err)
            fmt.Println(affect, "rows changed")
            _, err = db.Exec(`COPY transactions (hash, data, type, user_id, third_var)
			                  FROM '/home/z/psql.sql' with (FORMAT csv, DELIMITER '|')`)
            utils.CheckErr(err)
            affect, err = res.RowsAffected()
            utils.CheckErr(err)
            fmt.Println(affect, "rows changed")

            usedTransactions+="\\x"+hash+",";
            mrklRoot = utils.MerkleTreeRoot(mrklArray);
        }

        /*
		Заголовок
		TYPE (0-блок, 1-тр-я)     FF (256)
		BLOCK_ID   				       FF FF FF FF (4 294 967 295)
		TIME       					       FF FF FF FF (4 294 967 295)
		USER_ID                         FF FF FF FF FF (1 099 511 627 775)
		LEVEL                              FF (256)
		SIGN                               от 128 байта до 512 байт. Подпись от TYPE, BLOCK_ID, PREV_BLOCK_HASH, TIME, USER_ID, LEVEL, MRKL_ROOT
		Далее - тело блока (Тр-ии)
		*/

        // подписываем нашим нод-ключем заголовок блока

        // Extract the PEM-encoded data block
        block, _ := pem.Decode([]byte(nodePrivateKey))
        if block == nil {
            utils.CheckErr(errors.New("bad key data"))
        }
        if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
            utils.CheckErr(errors.New("unknown key type "+got+", want "+want));
        }
        //fmt.Println("block", block, "end")

        // Decode the RSA private key
        privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
        utils.CheckErr(err);
        //fmt.Println(privateKey)

		var forSign string
        forSign = fmt.Sprintf("0,%s,%s,%s,%s,%s,%s", newBlockId, prevBlock.Hash, Time, myUserId, level, mrklRoot)
        bytes, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, utils.HashSha1(forSign))
        utils.CheckErr(err)
        signatureHex := fmt.Sprintf("%x", bytes)

        // хэш шапки блока. нужен для сравнивания с другими и у кого будет меньше - у того блок круче
        headerHash := utils.DSha256([]byte(fmt.Sprintf("%s,%s,%s", myUserId, newBlockId, prevHeadHash)));

        data := fmt.Sprintf("%d|%d|%d|%d|\\x%s|\\x%s|\\x%s", newBlockId, Time, level, myUserId, headerHash, signatureHex, mrklRoot)
        /*file, err := ioutil.TempFile(os.TempDir(), "prefix")
		file.WriteString(data)
        file.Close()
        defer os.Remove(file.Name())*/
		name := os.TempDir()+"/Dcoin."+strconv.Itoa(math_rand.Intn(999999999))
        fmt.Println(name)
        file, _ := os.Create(name)
        defer file.Close()
        writer := bufio.NewWriter(file)
        writer.WriteString(data)
        writer.Flush()
        defer os.Remove(name)

        res, err := db.Exec("DELETE FROM testblock WHERE block_id = $1", newBlockId)
        utils.CheckErr(err)
        affect, err := res.RowsAffected()
        utils.CheckErr(err)
        res, err = db.Exec(`COPY testblock (block_id, time, level, user_id, header_hash, signature, mrkl_root)
			                               FROM '`+name+`' with (FORMAT csv, DELIMITER '|')`)
        utils.CheckErr(err)
        affect, err = res.RowsAffected()
        utils.CheckErr(err)
        fmt.Println(affect, "rows changed")

        /// #######################################
        // Отмечаем транзакции, которые попали в transactions_testblock
        // Пока для эксперимента
        // если не отмечать, то получается, что и в transactions_testblock и в transactions будут провернные тр-ии, которые откатятся дважды
        if len(usedTransactions)>0 {
            usedTransactions := usedTransactions[:len(usedTransactions)-1]
            fmt.Println("usedTransactions", usedTransactions)
            _, err = db.Exec("UPDATE transactions SET used=1 WHERE hash IN ($1)", usedTransactions)
            // для теста удаляем, т.к. она уже есть в transactions_testblock
            /*  $db->query( __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__, "
				  DELETE FROM `".DB_PREFIX."transactions`
				  WHERE `hash` IN ({$used_transactions})
				  ");*/
        }
        // ############################################

        db.DbUnlock(mainName);

		// в sqllite данные в db-файл пишутся только после закрытия всех соединений с БД.
        db.Close()
        db = utils.DbConnect(configIni)

        time.Sleep(3000 * time.Millisecond)
		break
    }
}

