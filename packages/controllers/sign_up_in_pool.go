package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"errors"
	"encoding/json"
	"github.com/c-darwin/dcoin-go-tmp/packages/schema"
	"fmt"
)

func jsonAnswer(err interface{}, answType string) error {
	var error_ string
	switch err.(type) {
	case string:
		error_ = err.(string)
	case error:
		error_ =  fmt.Sprintf("%v", err)
	}
	result, _ := json.Marshal(map[string]string{answType: fmt.Sprintf("%v", error_)})
	return errors.New(string(result))
}

func (c *Controller) SignUpInPool() (string, error) {

	log.Debug("1")
	c.w.Header().Set("Access-Control-Allow-Origin", "*")

	if !c.Community {
		return "", jsonAnswer("Not pool", "error")
	}

	c.r.ParseForm()
	log.Debug("2")

	var userId int64
	var codeSign string
	if c.SessUserId<=0 {
		// запрос пришел с десктопного кошелька юзера
		codeSign = c.r.FormValue("code_sign")
		if !utils.CheckInputData(codeSign, "hex_sign") {
			return "", jsonAnswer("Incorrect code_sign", "error")
		}
		userId = utils.StrToInt64(c.r.FormValue("user_id"))
		if !utils.CheckInputData(userId, "int64") {
			return "", jsonAnswer("Incorrect userId", "error")
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
		log.Debug("forSign: %v", forSign)
		publicKey, err := c.GetUserPublicKey(userId)
		log.Debug("publicKey: %x", publicKey)
		if err != nil {
			return "", jsonAnswer(utils.ErrInfo(err), "error")
		}
		// проверим подпись
		resultCheckSign, err := utils.CheckSign([][]byte{[]byte(publicKey)}, forSign, utils.HexToBin([]byte(codeSign)), true);
		if err != nil {
			return "", jsonAnswer(utils.ErrInfo(err), "error")
		}
		if !resultCheckSign {
			return "", jsonAnswer("Incorrect codeSign", "error")
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
		return "", jsonAnswer("Incorrect email", "error")
	}
	nodePrivateKey:=c.r.FormValue("node_private_key")
	if !utils.CheckInputData(nodePrivateKey, "private_key") {
		return "", jsonAnswer("Incorrect private_key", "error")
	}
	//publicKey := utils.MakeAsn1([]byte(n), []byte(e))
	log.Debug("3")

	// если мест в пуле нет, то просто запишем юзера в очередь
	pool_max_users, err := c.Single("SELECT pool_max_users FROM config").Int()
	if err != nil {
		return "", jsonAnswer(utils.ErrInfo(err), "error")
	}
	if len(c.CommunityUsers) >= pool_max_users {
		err = c.ExecSql("INSERT INTO pool_waiting_list ( email, time, user_id ) VALUES ( ?, ?, ? )", email, utils.Time(), userId)
		if err != nil {
			return "", jsonAnswer(utils.ErrInfo(err), "error")
		}
		return "", jsonAnswer(c.Lang["pool_is_full"], "error")
	}

	// регистрируем юзера в пуле
	// вначале убедитмся, что такой user_id у нас уже не зареган
	community, err := c.Single("SELECT user_id FROM community WHERE user_id  =  ?", userId).Int64()
	if err != nil {
		return "", jsonAnswer(utils.ErrInfo(err), "error")
	}
	if community != 0 {
		return fmt.Sprintf("%s", jsonAnswer(c.Lang["pool_user_id_is_busy"], "success")), nil
	}
	err = c.ExecSql("INSERT INTO community ( user_id ) VALUES ( ? )", userId)
	if err != nil {
		return "", jsonAnswer(utils.ErrInfo(err), "error")
	}

	schema_ := &schema.SchemaStruct{}
	schema_.DCDB = c.DCDB
	schema_.DbType = c.ConfigIni["db_type"]
	schema_.PrefixUserId = int(userId)
	schema_.GetSchema()

	prefix := utils.Int64ToStr(userId)+"_"
	err = c.ExecSql("INSERT INTO "+prefix+"my_table ( user_id, email ) VALUES ( ?, ? )", userId, email)
	if err != nil {
		return "", jsonAnswer(utils.ErrInfo(err), "error")
	}
	publicKey, err := c.GetUserPublicKey(userId)
	err = c.ExecSql("INSERT INTO "+prefix+"my_keys ( public_key, status ) VALUES ( [hex], 'approved' )", utils.BinToHex(publicKey))
	if err != nil {
		return "", jsonAnswer(utils.ErrInfo(err), "error")
	}
	nodePublicKey, err := utils.GetPublicFromPrivate(nodePrivateKey)
	if err != nil {
		return "", jsonAnswer(utils.ErrInfo(err), "error")
	}
	err = c.ExecSql("INSERT INTO "+prefix+"my_node_keys ( private_key, public_key ) VALUES ( ?, [hex] )", nodePrivateKey, nodePublicKey)
	if err != nil {
		return "", jsonAnswer(utils.ErrInfo(err), "error")
	}

	c.sess.Delete("restricted")
	return fmt.Sprintf("%s", jsonAnswer(c.Lang["pool_sign_up_success"], "success")), nil
}
