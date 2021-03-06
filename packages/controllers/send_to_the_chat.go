package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

func (c *Controller) SendToTheChat() (string, error) {

	c.r.ParseForm()
	message := c.r.FormValue("message")
	decryptMessage := c.r.FormValue("decrypt_message")
	sender := utils.StrToInt64(c.r.FormValue("sender"))
	receiver := utils.StrToInt64(c.r.FormValue("receiver"))
	lang := utils.StrToInt64(c.r.FormValue("lang"))
	room := utils.StrToInt64(c.r.FormValue("room"))
	status := utils.StrToInt64(c.r.FormValue("status"))
	signTime := utils.StrToInt64(c.r.FormValue("sign_time"))
	signature := []byte(c.r.FormValue("signature"))

	data := utils.Int64ToByte(lang)
	data = append(data, utils.Int64ToByte(room)...)
	data = append(data, utils.Int64ToByte(receiver)...)
	data = append(data, utils.Int64ToByte(sender)...)
	data = append(data, utils.Int64ToByte(status)...)
	data = append(data, []byte(message)...)
	data = append(data, utils.Int64ToByte(signTime)...)
	data = append(data, []byte(signature)...)

	hash := utils.Md5(data)

	err := c.CheckChatMessage(message, sender, receiver, lang, room, status, signTime, signature)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// на пуле сообщение сразу отобразится у всех
	if status == 1 {
		err = c.ExecSql(`INSERT INTO chat (hash, time, lang, room, receiver, sender, status, enc_message, message, sign_time, signature) VALUES ([hex], ?, ?, ?, ?, ?, ?, ?, ?, ?, [hex])`, hash, utils.Time(), lang, room, receiver, sender, 2, message, decryptMessage, signTime, signature)

	} else {
		err = c.ExecSql(`INSERT INTO chat (hash, time, lang, room, receiver, sender, status, message, sign_time, signature) VALUES ([hex], ?, ?, ?, ?, ?, ?, ?, ?, [hex])`, hash, utils.Time(), lang, room, receiver, sender, status, message, signTime, signature)
	}
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// даем команду рассыльщику, чтобы отправил всем хэш тр-ии сообщения
	utils.ChatNewTx <- true

	return utils.JsonAnswer("success", "success").String(), nil
}
