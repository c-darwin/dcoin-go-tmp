package availablekey

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"github.com/op/go-logging"
	"errors"
	"regexp"
	"math/rand"
	"github.com/c-darwin/dcoin-go-tmp/packages/schema"
)

var log = logging.MustGetLogger("availablekey")

type AvailablekeyStruct struct {
	*utils.DCDB
}

func (a *AvailablekeyStruct) checkAvailableKey(key string) (int64, string, error) {
	publicKeyAsn, err := utils.GetPublicFromPrivate(key)
	if err != nil {
		log.Debug("%v", err)
		return 0, "", utils.ErrInfo(err)
	}
	log.Debug("publicKeyAsn: %s", publicKeyAsn)
	userId, err := a.Single("SELECT user_id FROM users WHERE hex(public_key_0) = ?", publicKeyAsn).Int64()
	if err != nil {
		return 0, "", utils.ErrInfo(err)
	}
	log.Debug("userId: %s", userId)
	if userId == 0 {
		return 0, "", errors.New("null userId")
	}
	allTables, err := a.GetAllTables()
	if err != nil {
		return 0, "", utils.ErrInfo(err)
	}
	// может другой юзер уже начал смену ключа. актуально для пула
	if utils.InSliceString(utils.Int64ToStr(userId)+"_my_table", allTables) {
		return 0, "", errors.New("exists _my_table")
	}
	return userId, string(publicKeyAsn), nil
}

func (a *AvailablekeyStruct) GetAvailableKey() (int64, string, error) {
		keysStr, err := utils.GetHttpTextAnswer("http://dcoin.club/keys")
		if err != nil {
			return 0, "", utils.ErrInfo(err)
		}
		//keysStr = strings.Replace(keysStr, "\n", "", -1)
		r, _ := regexp.Compile("(?s)-----BEGIN RSA PRIVATE KEY-----(.*?)-----END RSA PRIVATE KEY-----")
		keys := r.FindAllString(keysStr, -1)
		for i := range keys {
			j := rand.Intn(i + 1)
			keys[i], keys[j] = keys[j], keys[i]
		}
		community, err := a.GetCommunityUsers()
		if err != nil {
			return 0, "", utils.ErrInfo(err)
		}
		for _, key := range keys {
			userId, pubKey, err := a.checkAvailableKey(key)
			if err != nil {
				log.Error("%s", utils.ErrInfo(err)) // тут ошибка - это нормально
			}
			log.Debug("checkAvailableKey userId: %v", userId)
			if userId > 0 {
				// запишем приватный ключ в БД, чтобы можно было подписать тр-ию на смену ключа
				myPref := ""
				if len(community) > 0 {
					schema_ := &schema.SchemaStruct{}
					schema_.DCDB = a.DCDB
					schema_.DbType = a.ConfigIni["db_type"]
					schema_.PrefixUserId = int(userId)
					schema_.GetSchema()
					myPref = utils.Int64ToStr(userId)+"_"
					err = a.ExecSql("INSERT INTO "+myPref+"my_table (user_id, status) VALUES (?, ?)", userId, "waiting_set_new_key")
					if err != nil {
						return 0, "", utils.ErrInfo(err)
					}
				} else {
					err = a.ExecSql("UPDATE my_table SET user_id = ?, status = ?", userId, "waiting_set_new_key")
					if err != nil {
						return 0, "", utils.ErrInfo(err)
					}
				}

				// пишем приватный в my_keys т.к. им будем подписывать тр-ию на смену ключа
				err = a.ExecSql("INSERT INTO "+myPref+"my_keys (private_key, public_key, status, block_id) VALUES (?, [hex], ?, ?)", key, pubKey, "approved", 1)
				if err != nil {
					return 0, "", utils.ErrInfo(err)
				}
				newPrivKey, newPubKey := utils.GenKeys()
				// сразу генерируем новый ключ и пишем приватный временно в my_keys, чтобы можно было выдавать юзеру для скачивания
				err = a.ExecSql("INSERT INTO "+myPref+"my_keys (private_key, public_key, status) VALUES (?, ?, ?)", newPrivKey, utils.HexToBin([]byte(newPubKey)), "my_pending")
				if err != nil {
					return 0, "", utils.ErrInfo(err)
				}
				return userId, pubKey, nil
			}
		}
		return 0, "", nil
}

