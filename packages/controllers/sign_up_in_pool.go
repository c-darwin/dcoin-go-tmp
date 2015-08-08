package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"errors"
	"encoding/json"
	"github.com/c-darwin/dcoin-go-tmp/packages/schema"
	"fmt"
)

func jsonErr(err interface{}) error {
	var error_ string
	switch err.(type) {
	case string:
		error_ = err.(string)
	case error:
		error_ =  fmt.Sprintf("%v", err)
	}
	result, _ := json.Marshal(map[string]string{"error": fmt.Sprintf("%v", error_)})
	return errors.New(string(result))
}

func (c *Controller) SignUpInPool() (string, error) {

	log.Debug("1")
	c.w.Header().Set("Access-Control-Allow-Origin", "*")

	successResult, _ := json.Marshal(map[string]string{"success": c.Lang["pool_sign_up_success"]})

	if !c.Community {
		return "", jsonErr("Not pool")
	}

	c.r.ParseForm()
	log.Debug("2")

	var userId int64
	var codeSign string
	if c.SessUserId<=0 {
		// запрос пришел с десктопного кошелька юзера
		codeSign = c.r.FormValue("code_sign")
		if !utils.CheckInputData(codeSign, "hex_sign") {
			return "", jsonErr("Incorrect code_sign")
		}
		userId = utils.StrToInt64(c.r.FormValue("user_id"))
		if !utils.CheckInputData(userId, "int64") {
			return "", jsonErr("Incorrect userId")
		}
		// получим данные для подписи
		var hash []byte
		if configIni["sign_hash"] == "ip" {
			hash = utils.Md5(c.r.RemoteAddr);
		} else {
			hash = utils.Md5(c.r.Header.Get("User-Agent")+c.r.RemoteAddr);
		}
		log.Debug("hash %s", hash)
		forSign, err := c.GetDataAuthorization(hash)
		publicKey, err := c.GetUserPublicKey(userId)
		if err != nil {
			return "", jsonErr(utils.ErrInfo(err))
		}
		// проверим подпись
		resultCheckSign, err := utils.CheckSign([][]byte{[]byte(publicKey)}, forSign, utils.HexToBin([]byte(codeSign)), true);
		if err != nil {
			return "", jsonErr(utils.ErrInfo(err))
		}
		if !resultCheckSign {
			return "", jsonErr("Incorrect codeSign")
		}
	} else {
		// запрос внутри пула
		userId = c.SessUserId
	}
	/*e:=c.r.FormValue("e")
	n:=c.r.FormValue("n")
	if len(e) == 0 || len(n) == 0 {
		result, _ := json.Marshal(map[string]string{"error": c.Lang["pool_error"]})
		return "", errors.New(string(result))
	}*/
	email:=c.r.FormValue("email")
	if !utils.ValidateEmail(email) {
		return "", jsonErr("Incorrect email")
	}
	nodePrivateKey:=c.r.FormValue("node_private_key")
	if !utils.CheckInputData(nodePrivateKey, "private_key") {
		return "", jsonErr("Incorrect private_key")
	}
	//publicKey := utils.MakeAsn1([]byte(n), []byte(e))
	log.Debug("3")

	// если мест в пуле нет, то просто запишем юзера в очередь
	pool_max_users, err := c.Single("SELECT pool_max_users FROM config").Int()
	if err != nil {
		return "", jsonErr(utils.ErrInfo(err))
	}
	if len(c.CommunityUsers) >= pool_max_users {
		err = c.ExecSql("INSERT INTO pool_waiting_list ( email, time, user_id ) VALUES ( ?, ?, ? )", email, utils.Time(), userId)
		if err != nil {
			return "", jsonErr(utils.ErrInfo(err))
		}
		return "", jsonErr(c.Lang["pool_is_full"])
	}

	// регистрируем юзера в пуле
	// вначале убедитмся, что такой user_id у нас уже не зареган
	community, err := c.Single("SELECT user_id FROM community WHERE user_id  =  ?", userId).Int64()
	if err != nil {
		return string(successResult), nil
	}
	if community != 0 {
		return "", jsonErr(c.Lang["pool_user_id_is_busy"])
	}
	err = c.ExecSql("INSERT INTO community ( user_id ) VALUES ( ? )", userId)
	if err != nil {
		return "", jsonErr(utils.ErrInfo(err))
	}

	schema_ := &schema.SchemaStruct{}
	schema_.DCDB = c.DCDB
	schema_.DbType = c.ConfigIni["db_type"]
	schema_.PrefixUserId = int(userId)
	schema_.GetSchema()

	prefix := utils.Int64ToStr(userId)+"_"
	err = c.ExecSql("INSERT INTO "+prefix+"my_table ( user_id, email ) VALUES ( ?, ? )", userId, email)
	if err != nil {
		return "", jsonErr(utils.ErrInfo(err))
	}
	publicKey, err := c.GetUserPublicKey(userId)
	err = c.ExecSql("INSERT INTO "+prefix+"my_keys ( public_key, status ) VALUES ( [hex], 'approved' )", utils.BinToHex(publicKey))
	if err != nil {
		return "", jsonErr(utils.ErrInfo(err))
	}
	nodePublicKey, err := utils.GetPublicFromPrivate(nodePrivateKey)
	if err != nil {
		return "", jsonErr(utils.ErrInfo(err))
	}
	err = c.ExecSql("INSERT INTO "+prefix+"my_node_keys ( private_key, public_key, status ) VALUES ( ?, [hex] 'pending' )", nodePrivateKey, nodePublicKey)
	if err != nil {
		return "", jsonErr(utils.ErrInfo(err))
	}

	c.sess.Delete("restricted")
	return string(successResult), nil
}
