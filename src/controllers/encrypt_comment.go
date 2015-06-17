package controllers
import (
	//"database/sql"
	//_ "github.com/lib/pq"
	//"reflect"
	//"fmt"
	//"html/template"
	//"bufio"
	//"bytes"
	//"time"
	//"regexp"
    "encoding/json"
	"utils"
	//"time"
	//"log"
//	"math"
	//"time"
	"errors"
	"crypto/rsa"
	"crypto/rand"
)

func (c *Controller) EncryptComment() (string, error) {

	var err error

	c.r.ParseForm()

	txType := c.r.FormValue("type")
	var toId int64
	var toIds []int64
	toIds_ := c.r.FormValue("toIds")
	if len (toIds_) == 0 {
		toId = utils.StrToInt64(c.r.FormValue("toId"))
	} else {
		err = json.Unmarshal([]byte(toIds_), &toIds)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	}

	comment := c.r.FormValue("comment")
	if len(comment) > 1024 {
		return `incorrect comment`, errors.New("incorrect comment")
	}

	var toUserId int64
	if txType == "project" {
		toUserId, err = c.Single("SELECT user_id FROM cf_projects WHERE id  =  ?", toId).Int64()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	} else {
		toUserId = toId
	}

	if len (toIds) == 0 {
		toIds = []int64{toUserId}
	}

	enc := make(map[int]string)
	for i:=0; i < len(toIds); i++ {
		if toIds[i] == 0 {
			enc[i] = "0"
			continue
		}
		// если получатель майнер, тогда шифруем нодовским ключем
		minersData, err := c.OneRow("SELECT miner_id, node_public_key FROM miners_data WHERE user_id  =  ?", toIds[i]).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		var publicKey string
		if utils.StrToInt(minersData["miner_id"]) > 0 && txType!="cash_request" && txType!="bug_reporting" && txType!="project" && txType!="money_back" {
			publicKey = minersData["node_public_key"]
		} else {
			publicKey, err = c.Single("SELECT public_key_0 FROM users WHERE user_id  =  ?", toIds[i]).String()
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}

		enc[i], err = rsa.EncryptPKCS1v15(rand.Reader, publicKey, comment)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	}
	if txType != "arbitration_arbitrators" {
		return enc[0], nil
	} else {
		result, err := json.Marshal(enc)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		return string(result), nil
	}
}
