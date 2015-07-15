package daemons

import (
    "fmt"
    _ "github.com/lib/pq"
    "time"
	"strconv"
    "crypto/x509"
    "encoding/pem"
    "crypto"
    "crypto/rand"
    "crypto/rsa"
    "crypto/md5"
	"bufio"
	"os"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
    "github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
)

var err error

func TestblockGenerator() {

    const GoroutineName = "TestblockGenerator"
    db := utils.DbConnect(configIni)
    db.GoroutineName = GoroutineName
    db.CheckInstall()

BEGIN:
	for {
        log.Debug("START")

        db.DbLock()

        blockId, err := db.GetBlockId()
		if err != nil {
            db.DbUnlock()
			log.Error("%v", err)
            utils.Sleep(1)
            continue BEGIN
        }
        newBlockId := blockId + 1;
        log.Debug("newBlockId: %v", newBlockId)
        testBlockId, err := db.GetTestBlockId()
        if err != nil {
            db.DbUnlock()
            log.Error("%v", err)
            utils.Sleep(1)
            continue BEGIN
        }

        log.Debug("testBlockId %v", testBlockId)

        if x, err := db.GetMyLocalGateIp(); x!="" {
            if err != nil {
                log.Error("%v", err)
			}
            log.Info("%v", "continue")
            db.DbUnlock()
            utils.Sleep(1)
            continue BEGIN
        }

        if testBlockId==newBlockId {
            db.DbUnlock()
            log.Error("%v", err)
            utils.Sleep(1)
            continue
        }

        prevBlock, myUserId, myMinerId, currentUserId, level, levelsRange, err := db.TestBlock()
        if err != nil {
            db.DbUnlock()
            log.Error("%v", err)
            utils.Sleep(1)
            continue BEGIN
        }
        log.Debug("%v %v %v %v %v %v", prevBlock, myUserId, myMinerId, currentUserId, level, levelsRange)

		if myMinerId==0 {
            db.DbUnlock()
            utils.Sleep(1)
			continue
		}

		sleep, err := db.GetGenSleep(prevBlock, level)
        if err!=nil {
            log.Error("%v", err)
            db.DbUnlock()
            utils.Sleep(1)
            continue
        }
        log.Debug("sleep %v", sleep)

        blockId = prevBlock.BlockId;
        log.Debug("blockId %v", blockId)
        prevHeadHash := prevBlock.HeadHash;
        log.Debug("prevHeadHash %v", prevHeadHash)

        // сколько прошло сек с момента генерации прошлого блока
        diff := time.Now().Unix() - prevBlock.Time;
        log.Debug("diff %v", diff)

        // вычитаем уже прошедшее время
        utils.SleepDiff(&sleep, diff)

        // Если случится откат или придет новый блок, то генератор блоков нужно запускать с начала, т.к. изменится max_miner_id.
        log.Debug("sleep %v", sleep)
        startSleep := time.Now().Unix();
        log.Debug("startSleep %v", startSleep)

        db.DbUnlock()

        for i := 0; i < int(sleep); i++ {
            db.DbLock();
            log.Debug("i %v", i)
            log.Debug("sleep %v", sleep)
			var newHeadHash string
            err := db.QueryRow("SELECT hex(head_hash) FROM info_block").Scan(&newHeadHash)
            utils.CheckErr(err)
            log.Debug("newHeadHash %v", newHeadHash)
            db.DbUnlock();
            if (newHeadHash != prevHeadHash) {
                log.Debug("newHeadHash!=prevHeadHash  %v  %v", newHeadHash, prevHeadHash)
                continue BEGIN
            }
            // из-за задержек с main_lock время уже прошло и выходим раньше, чем закончится цикл
            if time.Now().Unix() - startSleep > sleep {
                log.Debug("break")
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
			log.Error("%v", err)
            db.DbUnlock()
            utils.Sleep(1)
            continue
		}
        log.Debug("%v %v %v %v %v %v", prevBlock, myUserId, myMinerId, currentUserId, level, levelsRange)
        // сколько прошло сек с момента генерации прошлого блока
        diff = time.Now().Unix() - prevBlock.Time;
        log.Debug("diff %v", diff)
        // вычитаем уже прошедшее время
        utils.SleepDiff(&sleep, diff)
        log.Debug("sleep %v", sleep)
        // если нужно доспать, то просто вернемся в начало и доспим нужное время. И на всякий случай убедимся, что блок не изменился
        if sleep > 0 || prevBlock.HeadHash != prevHeadHash {
            log.Debug("continue")
            db.DbUnlock()
            time.Sleep(1000 * time.Millisecond)
            continue
        }
        log.Debug("blockgeneration begin")
        blockId = prevBlock.BlockId;
        if blockId < 1 {
            log.Debug("continue")
            db.DbUnlock()
            time.Sleep(1000 * time.Millisecond)
            continue
        }

        newBlockId = blockId + 1;
		var myPrefix string
        CommunityUser, err := db.GetCommunityUsers()
        if err != nil {
            log.Error("%v", err)
            db.DbUnlock()
            utils.Sleep(1)
            continue
        }
        if len(CommunityUser)>0 {
            myPrefix = strconv.FormatInt(myUserId, 10)+"_"
		} else {
            myPrefix = ""
	    }
        nodePrivateKey, err := db.GetNodePrivateKey(myPrefix)
		if len(nodePrivateKey) < 1 {
            log.Debug("continue")
            db.DbUnlock()
            time.Sleep(1000 * time.Millisecond)
            continue
        }
        prevHeadHash = prevBlock.HeadHash

        log.Debug("%v %v", nodePrivateKey, prevHeadHash)

        //#####################################
        //##		 Формируем блок
        //#####################################
        log.Debug("%v %v", newBlockId, currentUserId)
        if currentUserId < 1 {
            log.Debug("continue")
            db.DbUnlock()
            time.Sleep(1000 * time.Millisecond)
            continue
        }
        if prevBlock.BlockId >= newBlockId {
            log.Debug("continue")
            db.DbUnlock()
            time.Sleep(1000 * time.Millisecond)
            continue
        }
        // откатим transactions_testblock
		p := new(dcparser.Parser)
        p.DCDB = db
        p.RollbackTransactionsTestblock(true)

       Time := time.Now().Unix()

        // переведем тр-ии в `verified` = 1
        err = p.AllTxParser()
        if err != nil {
            db.PrintSleep(err, 1)
            continue
        }

        var mrklArray  [][]byte
		var usedTransactions string
		var mrklRoot []byte
        // берем все данные из очереди. Они уже были проверены ранее, и можно их не проверять, а просто брать
        rows, err := db.Query(db.FormatQuery("SELECT data, hex(hash), type, user_id, third_var FROM transactions WHERE used=0 AND verified = 1"))
        utils.CheckErr(err)
        for rows.Next() {
            var data []byte
            var hash string
            var txType string
            var txUserId string
            var thirdVar string
            err = rows.Scan(&data, &hash, &txType, &txUserId, &thirdVar)
            utils.CheckErr(err)
            log.Debug("data %v", data)
            log.Debug("hash %v", hash)
            transactionType := data[1:2];
            log.Debug("%v", transactionType)
            fmt.Printf("%x", transactionType)
            mrklArray = append(mrklArray, utils.DSha256(data));
            log.Debug("mrklArray %v", mrklArray)

            hash2_ := md5.New()
            hash2_.Write(data)
            hashMd5:=fmt.Sprintf("%x", hash2_.Sum(nil))
            log.Debug(hashMd5)

            dataHex := fmt.Sprintf("%x", data)
            log.Debug("dataHex %v", dataHex)

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
            log.Debug("rows changed %v", affect)
            _, err = db.Exec(`COPY transactions (hash, data, type, user_id, third_var)
			                  FROM '/home/z/psql.sql' with (FORMAT csv, DELIMITER '|')`)
            utils.CheckErr(err)
            affect, err = res.RowsAffected()
            utils.CheckErr(err)
            log.Debug("rows changed %v", affect)

            usedTransactions+="\\x"+hash+",";
        }

        if len(mrklArray) == 0 {
            mrklArray = append(mrklArray,  []byte("0"))
        }
        mrklRoot = utils.MerkleTreeRoot(mrklArray);
        log.Debug("mrklRoot: %s", mrklRoot)
        log.Debug("mrklRoot: %x", mrklRoot)

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
            log.Error("bad key data %v ", utils.GetParent())
            utils.Sleep(1)
            continue BEGIN
        }
        if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
            log.Error("unknown key type %v, want %v / %v ", got, want, utils.GetParent())
            utils.Sleep(1)
            continue BEGIN
        }
        privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
        if err != nil {
            log.Error("err %v %v", err, utils.GetParent())
            utils.Sleep(1)
            continue BEGIN
        }
        var forSign string
        forSign = fmt.Sprintf("0,%v,%v,%v,%v,%v,%v", newBlockId, prevBlock.Hash, Time, myUserId, level, string(mrklRoot))
        log.Debug("forSign: %v", forSign)
        bytes, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, utils.HashSha1(forSign))
        if err != nil {
            log.Error("err %v %v", err, utils.GetParent())
            utils.Sleep(1)
            continue BEGIN
        }
        signatureHex := fmt.Sprintf("%x", bytes)

        // хэш шапки блока. нужен для сравнивания с другими и у кого будет меньше - у того блок круче
        headerHash := utils.DSha256([]byte(fmt.Sprintf("%s,%s,%s", myUserId, newBlockId, prevHeadHash)));

        /*data := fmt.Sprintf("%d|%d|%d|%d|\\x%s|\\x%s|\\x%s", newBlockId, Time, level, myUserId, headerHash, signatureHex, mrklRootHex)
		name := os.TempDir()+"/Dcoin."+strconv.Itoa(math_rand.Intn(999999999))
        log.Debug(name)
        file, _ := os.Create(name)
        defer file.Close()
        writer := bufio.NewWriter(file)
        writer.WriteString(data)
        writer.Flush()
        defer os.Remove(name)*/

        err = db.ExecSql("DELETE FROM testblock WHERE block_id = ?", newBlockId)
        if err != nil {
            db.PrintSleep(err, 1)
            continue BEGIN
        }
        err = db.ExecSql(`INSERT INTO testblock (block_id, time, level, user_id, header_hash, signature, mrkl_root) VALUES (?, ?, ?, ?, [hex], [hex], [hex])`,
            newBlockId, Time, level, myUserId, headerHash, signatureHex, string(mrklRoot))
        if err != nil {
            db.PrintSleep(err, 1)
            continue BEGIN
        }

        /// #######################################
        // Отмечаем транзакции, которые попали в transactions_testblock
        // Пока для эксперимента
        // если не отмечать, то получается, что и в transactions_testblock и в transactions будут провернные тр-ии, которые откатятся дважды
        if len(usedTransactions)>0 {
            usedTransactions := usedTransactions[:len(usedTransactions)-1]
            log.Debug("usedTransactions %v", usedTransactions)
            _, err = db.Exec("UPDATE transactions SET used=1 WHERE hash IN ($1)", usedTransactions)
            if err != nil {
                db.PrintSleep(err, 1)
                continue BEGIN
            }
            // для теста удаляем, т.к. она уже есть в transactions_testblock
            /*  $db->query( __FILE__, __LINE__,  __FUNCTION__,  __CLASS__, __METHOD__, "
				  DELETE FROM `".DB_PREFIX."transactions`
				  WHERE `hash` IN ({$used_transactions})
				  ");*/
        }
        // ############################################

        db.DbUnlock();

		// в sqllite данные в db-файл пишутся только после закрытия всех соединений с БД.
        db.Close()
        db = utils.DbConnect(configIni)

        log.Debug("END")

        time.Sleep(3000 * time.Millisecond)
		break
    }
}

