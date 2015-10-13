package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"strings"
	"crypto/rsa"
	"crypto/rand"
	"text/template"
)

var chatIds = make(map[int64][]int)

func (c *Controller) GetChatMessages() (string, error) {

	c.r.ParseForm()
	first := c.r.FormValue("first")
	room := utils.StrToInt64(c.r.FormValue("room"))
	lang := utils.StrToInt64(c.r.FormValue("lang"))

	if first == "1" {
		chatIds[c.SessUserId] = []int{}
	}
	ids := ""
	if len(chatIds[c.SessUserId]) > 0 {
		ids = `AND id NOT IN(`+strings.Join(utils.IntSliceToStr(chatIds[c.SessUserId]), ",")+`)`
	}
	var result string
	chatData, err := c.GetAll(`SELECT * FROM chat WHERE room = ? AND lang = ?  `+ids+` ORDER BY id ASC LIMIT 100`, 100, room, lang)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	for _, data := range chatData {
		status := data["status"]
		message := data["message"]
		receiver := utils.StrToInt64(data["receiver"])
		sender := utils.StrToInt64(data["sender"])
		if status == "1" {
			// Если юзер хранит приватый ключ в БД, то сможем расшифровать прямо тут
			if receiver == c.SessUserId {
				privateKey, err := c.GetMyPrivateKey(c.MyPrefix)
				if err != nil {
					log.Error("%v", utils.ErrInfo(err))
					continue
				}
				if len(privateKey) > 0 {
					rsaPrivateKey, err := utils.MakePrivateKey(privateKey)
					if err != nil {
						log.Error("%v", utils.ErrInfo(err))
						continue
					}
					decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, rsaPrivateKey, utils.HexToBin([]byte(data["message"])))
					if err != nil {
						log.Error("%v", utils.ErrInfo(err))
						continue
					}
					if len(decrypted) > 0 {
						err = c.ExecSql(`UPDATE chat SET enc_message = message, message = ?, status = ? WHERE id = ?`, decrypted, 2, data["id"])
						if err != nil {
							log.Error("%v", utils.ErrInfo(err))
							continue
						}
						message = string(decrypted)
						status = "2"
					}
				}
			}
		}

		name := data["sender"]
		ava := "/static/img/noavatar.png"
		// возможно у отпарвителя есть ник
		nameAva, err := c.OneRow(`SELECT name, avatar FROM users WHERE user_id = ?`, sender).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if len(nameAva["name"]) > 0 {
			name = nameAva["name"]
		}

		minerStatus, err := c.Single(`SELECT status FROM miners_data WHERE user_id = ?`, sender).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if minerStatus == "miner" && len(nameAva["avatar"]) > 0 {
			ava = nameAva["avatar"]
		}

		row := ""
		message = template.HTMLEscapeString(message)
		name = `<a class="chatNick" onclick='setReceiver("`+name+`", "`+data["sender"]+`")'>`+name+`</a>`
		if status == "2" { // успешно расшифровали
			row = `<tr><td><img src="`+ava+`"></strong>`+name+`</strong>: <i class="fa fa-lock"></i> `+message+`</td></tr>`
		} else if status == "1" && receiver == c.SessUserId { // либо нет ключа, либо какая-то ошибка
			row = `<tr><td><img src="`+ava+`"></strong>`+name+`</strong>: <div id="comment_`+data["id"]+`" style="display: inline-block;"><input type="hidden" value="`+message+`" id="encrypt_comment_`+data["id"]+`"><a class="btn btn-default btn-lg" onclick="decrypt_comment(`+data["id"]+`, 'chat')"> <i class="fa fa-lock"></i> Decrypt</a></div></td></tr>`
		} else if status == "0" {
			row = `<tr><td><img src="`+ava+`"></strong>`+name+`</strong>: `+message+`</td></tr>`
		}
		result += row
		chatIds[c.SessUserId] = append(chatIds[c.SessUserId], utils.StrToInt(data["id"]))
	}

	return utils.JsonAnswer(result, "messages").String(), nil
}