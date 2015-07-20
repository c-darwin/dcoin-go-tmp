package daemons

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/dcparser"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"crypto/rand"
	"regexp"
	"net/http"
	"net/url"
	"bytes"
)


/*
 * Важно! отключать демона при обнулении данных в БД
*/

func Shop() {

	const GoroutineName = "Shop"

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

		myBlockId, err := db.GetMyBlockId()
		blockId, err := db.GetBlockId()
		if myBlockId > blockId {
			db.PrintSleep(utils.ErrInfo(err), 1)
			continue
		}
		currencyList, err := db.GetCurrencyList(false)
		if err != nil {
			db.PrintSleep(utils.ErrInfo(err), 1)
			continue
		}
		// нужно знать текущий блок, который есть у большинства нодов
		blockId, err = db.GetConfirmedBlockId()
		if err != nil {
			db.PrintSleep(utils.ErrInfo(err), 1)
			continue
		}

		// сколько должно быть подтверждений, т.е. кол-во блоков сверху
		confirmations := int64(5)

		// берем всех юзеров по порядку
		community, err := db.GetCommunityUsers()
		if err != nil {
			db.PrintSleep(utils.ErrInfo(err), 1)
			continue
		}
		for _, userId := range community {
			privateKey := ""
			myPrefix := utils.Int64ToStr(userId)+"_"
			allTables, err := db.GetAllTables()
			if err != nil {
				db.PrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			if !utils.InSliceString(myPrefix+"my_keys", allTables) {
				continue
			}
			// проверим, майнер ли
			minerId, err := db.GetMyMinerId(userId)
			if err != nil {
				db.PrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			if minerId > 0 {
				// наш приватный ключ нода, которым будем расшифровывать комменты
				privateKey, err = db.GetNodePrivateKey(myPrefix)
				if err != nil {
					db.PrintSleep(utils.ErrInfo(err), 1)
					continue BEGIN
				}
			}
			// возможно, что комменты будут зашифрованы юзерским ключем
			if len(privateKey) == 0 {
				privateKey, err = db.GetMyPrivateKey(myPrefix)
				if err != nil {
					db.PrintSleep(utils.ErrInfo(err), 1)
					continue BEGIN
				}
			}
			// если это еще не майнер и админ ноды не указал его приватный ключ в табле my_keys, то $private_key будет пуст
			if len(privateKey) == 0 {
				continue
			}
			myData, err := db.OneRow("SELECT shop_secret_key, shop_callback_url FROM "+myPrefix+"my_table").String()
			if err != nil {
				db.PrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}

			// Получаем инфу о входящих переводах и начисляем их на счета юзеров
			rows, err := db.Query(`
					SELECT id, block_id, type_id, currency_id, amount, to_user_id, comment_status, comment
					FROM `+myPrefix+`my_dc_transactions
					WHERE type = 'from_user' AND
								 block_id < ? AND
								 merchant_checked = 0 AND
								 status = 'approved'
					ORDER BY id DESC
					`, blockId - confirmations)
			if err != nil {
				db.PrintSleep(utils.ErrInfo(err), 1)
				continue BEGIN
			}
			defer rows.Close()
			for rows.Next() {
				var id, block_id, type_id, currency_id, to_user_id int64
				var comment_status, comment string
				var amount float64
				err = rows.Scan(&id, &block_id, &type_id, &currency_id, &amount, &to_user_id, &comment_status, &comment)
				if err != nil {
					db.PrintSleep(utils.ErrInfo(err), 1)
					continue BEGIN
				}
				if len(myData["shop_callback_url"]) == 0 {
					// отметим merchant_checked=1, чтобы больше не брать эту тр-ию
					err = db.ExecSql("UPDATE "+myPrefix+"my_dc_transactions SET merchant_checked = 1 WHERE id = ?", id)
					if err != nil {
						db.PrintSleep(utils.ErrInfo(err), 1)
						continue BEGIN
					}
					continue
				}

				// вначале нужно проверить, точно ли есть такой перевод в блоке
				binaryData, err := db.Single("SELECT data FROM block_chain WHERE id  =  ?", blockId).Bytes()
				if err != nil {
					db.PrintSleep(utils.ErrInfo(err), 1)
					continue BEGIN
				}
				p := new(dcparser.Parser)
				p.BinaryData = binaryData
				p.ParseDataLite()
				for _, txMap := range p.TxMapsArr {

					// пропускаем все ненужные тр-ии
					if txMap.Int64["type"] != utils.TypeInt("SendDc") {
						continue
					}

					// сравнение данных из таблы my_dc_transactions с тем, что в блоке
					if txMap.Int64["user_id"] == userId && txMap.Int64["currency_id"] == currency_id && txMap.Money["amount"] == amount && txMap.Int64["to_user_id"] == to_user_id {
						decryptedComment := ""
						// расшифруем коммент
						if comment_status == "encrypted" {

							block, _ := pem.Decode([]byte(privateKey));
							if block == nil || block.Type != "RSA PRIVATE KEY" {
								db.PrintSleep(utils.ErrInfo(err), 1)
								continue BEGIN
							}
							private_key, err := x509.ParsePKCS1PrivateKey(block.Bytes);
							if err != nil {
								db.PrintSleep(utils.ErrInfo(err), 1)
								continue BEGIN
							}
							decryptedComment_, err := rsa.DecryptPKCS1v15(rand.Reader, private_key, []byte(comment))
							if err != nil {
								db.PrintSleep(utils.ErrInfo(err), 1)
								continue BEGIN
							}
							decryptedComment = string(decryptedComment_)
							// запишем расшифрованный коммент, чтобы потом можно было найти перевод в ручном режиме
							err = db.ExecSql("UPDATE "+myPrefix+"my_dc_transactions SET comment = ?, comment_status = 'decrypted' WHERE id = ?", decryptedComment, id)
							if err != nil {
								db.PrintSleep(utils.ErrInfo(err), 1)
								continue BEGIN
							}
						} else {
							decryptedComment = comment
						}

						// возможно, что чуть раньше было reduction, а это значит, что все тр-ии,
						// которые мы ещё не обработали и которые были До блока с reduction нужно принимать с учетом reduction
						// т.к. средства на нашем счете уже урезались, а  вот те, что после reduction - остались в том виде, в котором пришли
						lastReduction, err := db.OneRow("SELECT block_id, pct FROM reduction WHERE currency_id  = ? ORDER BY block_id", currency_id).Int64()
						if err != nil {
							db.PrintSleep(utils.ErrInfo(err), 1)
							continue BEGIN
						}
						if blockId <= lastReduction["block_id"] {
							// сумму с учетом reduction
							k0 := (100 - lastReduction["pct"]) / 100
							amount = amount * float64(k0)
						}

						// делаем запрос к callback скрипту
						r, _ := regexp.Compile(`(?i)\s*#\s*([0-9]+)\s*`)
						order := r.FindStringSubmatch(decryptedComment)
						orderId := 0
						if len(order) > 0 {
							orderId = utils.StrToInt(order[1])
						}
						txId := id
						sign := fmt.Sprintf("%v:%v:%v:%v:%v:%v:%v:%v", amount, currencyList[currency_id], orderId, decryptedComment, txMap.Int64["user_id"], blockId, txId, myData["shop_secret_key"])
						data := url.Values{}
						data.Add("amount", utils.Float64ToStrPct(amount))
						data.Add("currency", currencyList[currency_id])
						data.Add("order_id", utils.IntToStr(orderId))
						data.Add("message", decryptedComment)
						data.Add("user_id", utils.Int64ToStr(txMap.Int64["user_id"]))
						data.Add("block_id", utils.Int64ToStr(txMap.Int64["block_id"]))
						data.Add("tx_id", utils.Int64ToStr(txId))
						data.Add("sign", sign)

						client := &http.Client{}
						req, err := http.NewRequest("POST", myData["shop_callback_url"], bytes.NewBufferString(data.Encode()))
						if err != nil {
							db.PrintSleep(utils.ErrInfo(err), 1)
							continue BEGIN
						}
						req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
						req.Header.Add("Content-Length", utils.IntToStr(len(data.Encode())))

						resp, err := client.Do(req)
						if err != nil {
							db.PrintSleep(utils.ErrInfo(err), 1)
							continue BEGIN
						}
						//contents, _ := ioutil.ReadAll(resp.Body)
						if resp.StatusCode == 200 {
							// отметим merchant_checked=1, чтобы больше не брать эту тр-ию
							err = db.ExecSql("UPDATE "+myPrefix+"my_dc_transactions SET merchant_checked = 1 WHERE id = ?", id)
							if err != nil {
								db.PrintSleep(utils.ErrInfo(err), 1)
								continue BEGIN
							}
						}
					}
				}
			}
		}

		for i:=0; i < 60; i++ {
			utils.Sleep(1)
			// проверим, не нужно ли нам выйти из цикла
			if CheckDaemonsRestart() {
				break BEGIN
			}
		}
	}
}
