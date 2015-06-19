package controllers
import (
	"utils"
	"errors"
	"encoding/json"
	"schema"
)

func (c *Controller) SignUpInPool() (string, error) {
	if c.SessUserId<=0 {
		return "", errors.New("c.SessUserId<=0")
	}
	c.r.ParseForm()
	e:=c.r.FormValue("e")
	n:=c.r.FormValue("n")
	email:=c.r.FormValue("email")
	if len(e) == 0 || len(n) == 0 {
		result, _ := json.Marshal(map[string]string{"error": c.Lang["pool_error"]})
		return "", errors.New(string(result))
	}
	if !utils.ValidateEmail(email) {
		result, _ := json.Marshal(map[string]string{"error": "Incorrect email"})
		return "", errors.New(string(result))
	}

	publicKey := utils.MakeAsn1([]byte(n), []byte(e))

	// если мест в пуле нет, то просто запишем юзера в очередь
	pool_max_users, err := c.Single("SELECT pool_max_users FROM config").Int()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	if len(c.CommunityUsers) >= pool_max_users {
		err = c.ExecSql("INSERT INTO pool_waiting_list ( email, time, user_id ) VALUES ( ?, ?, ? )", email, utils.Time(), c.SessUserId)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		result, _ := json.Marshal(map[string]string{"error": c.Lang["pool_is_full"]})
		return "", errors.New(string(result))
	}

	// регистрируем юзера в пуле
	// вначале убедитмся, что такой user_id у нас уже не зареган
	community, err := c.Single("SELECT user_id FROM community WHERE user_id  =  ?", c.SessUserId).Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	if community != 0 {
		result, _ := json.Marshal(map[string]string{"error": c.Lang["pool_user_id_is_busy"]})
		return "", errors.New(string(result))
	}
	err = c.ExecSql("INSERT IGNORE INTO community ( user_id ) VALUES ( ? )", c.SessUserId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	schema_ := &schema.SchemaStruct{}
	schema_.DCDB = c.DCDB
	schema_.DbType = c.ConfigIni["db_type"]
	schema_.PrefixUserId = int(c.SessUserId)
	schema_.GetSchema()

	prefix := utils.Int64ToStr(c.SessUserId)+"_"
	err = c.ExecSql("INSERT IGNORE INTO "+prefix+"my_table ( user_id, email ) VALUES ( ?, ? )", c.SessUserId, email)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	err = c.ExecSql("INSERT IGNORE INTO "+prefix+"my_keys ( public_key, status ) VALUES ( [hex], 'approved' )", publicKey)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	c.sess.Delete("restricted")

	result, _ := json.Marshal(map[string]string{"success": c.Lang["pool_sign_up_success"]})
	return string(result), nil
}
