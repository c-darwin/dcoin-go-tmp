package controllers
import (
	"fmt"
	"dcoin/packages/utils"
	"time"
	"log"
)


func (c *Controller) Check_sign() (string, error) {

	var checkError bool

	c.r.ParseForm()
	n := []byte(c.r.FormValue("n"))
	e := []byte(c.r.FormValue("e"))
	sign := []byte(c.r.FormValue("sign"))
	setupPassword := c.r.FormValue("setup_password")
	if !utils.CheckInputData(n, "hex") {
		return `{"result":"incorrect n"}`, nil
	}
	if !utils.CheckInputData(e, "hex") {
		return `{"result":"incorrect e"}`, nil
	}
	if !utils.CheckInputData(string(sign), "hex_sign") {
		return `{"result":"incorrect sign"}`, nil
	}

	allTables, err := c.DCDB.GetAllTables()
	if err != nil {
		return "{\"result\":0}", err
	}

	var hash []byte
	if configIni["sign_hash"] == "ip" {
		hash = utils.Md5(c.r.RemoteAddr);
	} else {
		hash = utils.Md5(c.r.Header.Get("User-Agent")+c.r.RemoteAddr);
	}
	community, err := c.DCDB.GetCommunityUsers()
	if err != nil {
		return "{\"result\":0}", err
	}
	if len(community) > 0 {

		// в цикле проверяем, кому подойдет присланная подпись
		for _, userId := range community {

			myPrefix := utils.Int64ToStr(userId)+"_"
			if !utils.InSliceString(myPrefix+"my_keys", allTables) {
				continue
			}

			// получим открытый ключ юзера
			publicKey, err := c.DCDB.GetMyPublicKey(myPrefix)
			if err != nil {
				return "{\"result\":0}", err
			}

			// получим данные для подписи
			forSign, err := c.DCDB.GetDataAuthorization(hash)

			// проверим подпись
			resultCheckSign, err := utils.CheckSign([][]byte{publicKey}, forSign, utils.HexToBin(sign), true);
			if err != nil {
				return "{\"result\":0}", err
			}
			// если подпись верная, значит мы нашли юзера, который эту подпись смог сделать
			if resultCheckSign {
				myUserId := userId
				// убираем ограниченный режим
				c.sess.Delete("restricted")
				c.sess.Set("user_id", myUserId)
				public_key, err := c.DCDB.GetUserPublicKey(myUserId)
				if err != nil {
					return "{\"result\":0}", err
				}
				// паблик кей в сессии нужен чтобы выбрасывать юзера, если ключ изменился
				c.sess.Set("public_key", public_key)

				adminUSerID, err := c.DCDB.GetAdminUserId();
				if err != nil {
					return "{\"result\":0}", err
				}
				if adminUSerID == myUserId {
					c.sess.Set("admin", 1)
				}
				return "{\"result\":1}", nil
			}
		}

		// если дошли досюда, значит ни один ключ не подошел и даем возможность войти в ограниченном режиме
		publicKey := utils.MakeAsn1(n, e)
		userId, err := c.DCDB.GetUserIdByPublicKey(publicKey)
		if err != nil {
			return "{\"result\":0}", err
		}

		// юзер с таким ключем есть в БД
		if len(userId) > 0 {

			// получим данные для подписи
			forSign, err := c.DCDB.GetDataAuthorization(hash)

			// проверим подпись
			resultCheckSign, err := utils.CheckSign([][]byte{publicKey}, forSign, utils.HexToBin(sign), true);
			if err != nil {
				return "{\"result\":0}", err
			}
			if resultCheckSign {

				// если юзер смог подписать наш хэш, значит у него актуальный праймари ключ
				c.sess.Set("user_id", userId)

				// возможно в табле my_keys старые данные, но если эта табла есть, то нужно добавить туда ключ
				if utils.InSliceString(userId+"_my_keys", allTables) {
					curBlockId, err := c.DCDB.GetBlockId()
					if err != nil {
						return "{\"result\":0}", err
					}
					err = c.DCDB.InsertIntoMyKey(userId+"_", publicKey, utils.Int64ToStr(curBlockId))
					if err != nil {
						return "{\"result\":0}", err
					}
					c.sess.Delete("restricted")
				} else {
					c.sess.Set("restricted", 1)
				}
				return "{\"result\":1}", nil
			} else {
				return "{\"result\":0}", nil
			}
		}
	} else {
		// получим открытый ключ юзера
		publicKey, err := c.DCDB.GetMyPublicKey("")
		if err != nil {
			return "{\"result\":0}", err
		}

		// Если ключ еще не успели установить
		if len(publicKey) == 0 {

			fmt.Println("len(publicKey) == 0")
			// пока не собрана цепочка блоков не даем ввести ключ
			infoBlock, err := c.DCDB.GetInfoBlock()
			if err != nil {
				return "{\"result\":0}", err
			}
			fmt.Println("infoBlock", infoBlock)
			// если последний блок не старше 2-х часов
			wTime := int64(2)
			if c.ConfigIni["test_mode"] == "1" {
				wTime = 365*24
			}
			log.Println(time.Now().Unix(), utils.StrToInt64(infoBlock["time"]), wTime)
			if (time.Now().Unix() - utils.StrToInt64(infoBlock["time"])) < 3600*wTime  {

				// проверим, верный ли установочный пароль, если он, конечно, есть
				setupPassword_, err := c.Single("SELECT setup_password FROM config").String()
				if err != nil {
					return "{\"result\":0}", err
				}
				if len(setupPassword_) > 0 && setupPassword_ != fmt.Sprintf("%x", utils.HashSha1(setupPassword)) {
					log.Println(setupPassword_, fmt.Sprintf("%x", utils.HashSha1(setupPassword)), setupPassword)
					return "{\"result\":0}", nil
				}

				publicKey := utils.MakeAsn1(n, e)
				log.Println("new key", string(publicKey))
				userId, err := c.GetUserIdByPublicKey(publicKey)
				if err != nil {
					return "{\"result\":0}", err
				}

				// получим данные для подписи
				forSign, err := c.DCDB.GetDataAuthorization(hash)
				// проверим подпись
				resultCheckSign, err := utils.CheckSign([][]byte{utils.HexToBin(publicKey)}, forSign, utils.HexToBin(sign), true);
				if err != nil {
					return "{\"result\":0}", err
				}
				if !resultCheckSign {
					return "{\"result\":0}", nil
				}

				if len(userId) > 0 {
					err := c.InsertIntoMyKey("", publicKey, "0")
					if err != nil {
						return "{\"result\":0}", err
					}
					myUserId, err := c.GetMyUserId("")
					if myUserId > 0 {
						c.ExecSql("UPDATE my_table SET user_id=?, status = 'user'", userId)
					} else {
						c.ExecSql("INSERT INTO my_table (user_id, status) VALUES (?, 'user')", userId)
					}
				} else {
					checkError = true
				}
			} else {
				checkError = true
			}
		} else {

			if configIni["sign_hash"] == "ip" {
				hash = utils.Md5(c.r.RemoteAddr);
			} else {
				hash = utils.Md5(c.r.Header.Get("User-Agent")+c.r.RemoteAddr);
			}
			log.Println("hash", hash)

			// получим данные для подписи
			forSign, err := c.DCDB.GetDataAuthorization(hash)
			log.Println("forSign", forSign)
			log.Printf("publicKey %x\n", string(publicKey))
			log.Println("sign_",string(sign))
			// проверим подпись
			resultCheckSign, err := utils.CheckSign([][]byte{publicKey}, forSign, utils.HexToBin(sign), true);
			if err != nil {
				return "{\"result\":0}", err
			}
			if !resultCheckSign {
				return "{\"result\":0}", nil
			}
		}

		if checkError {
			return "{\"result\":0}", nil
		} else {
			myUserId, err := c.DCDB.GetMyUserId("")
			if myUserId==0 {
				myUserId = -1
			}
			if err != nil {
				return "{\"result\":0}", err
			}
			c.sess.Delete("restricted")
			c.sess.Set("user_id", myUserId)

			// если уже пришел блок, в котором зареган ключ юзера
			if (myUserId!=-1) {

				public_key, err := c.DCDB.GetUserPublicKey(myUserId)
				if err != nil {
					return "{\"result\":0}", err
				}
				// паблик кей в сессии нужен чтобы выбрасывать юзера, если ключ изменился
				c.sess.Set("public_key", public_key)

				AdminUserId, err := c.DCDB.GetAdminUserId()
				if err != nil {
					return "{\"result\":0}", err
				}
				if AdminUserId == myUserId {
					c.sess.Set("admin", int64(1))
				}
				return "{\"result\":1}", nil
			}
		}
	}
	return "{\"result\":0}", nil
}
