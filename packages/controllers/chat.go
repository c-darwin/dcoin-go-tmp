package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

type chatPage struct {
	Lang                  map[string]string
	CountSign             int
	CountSignArr          []int
	SignData              string
	ShowSignData          bool
	IOS                   bool
	Mobile				  bool
	MyChatName			  string
	UserId                int64
}

func (c *Controller) Chat() (string, error) {

	myChatName := utils.Int64ToStr(c.SessUserId)
	// возможно у отпарвителя есть ник
	name, err := c.Single(`SELECT name FROM users WHERE user_id = ?`, c.SessUserId).String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	if len(name) > 0 {
		myChatName = name
	}

	TemplateStr, err := makeTemplate("chat", "chat", &chatPage{
		CountSignArr:          c.CountSignArr,
		CountSign:             c.CountSign,
		Lang:                  c.Lang,
		ShowSignData:          c.ShowSignData,
		SignData:              "",
		MyChatName:			   myChatName,
		UserId:                c.SessUserId,
		IOS:                   utils.IOS(),
		Mobile:				   utils.Mobile()})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
