package controllers

import (
	"errors"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"encoding/json"
	"github.com/c-darwin/dcoin-go-tmp/packages/consts"
	"crypto/rsa"
	"crypto"
	"github.com/mcuadros/go-version"
	"runtime"
	"strings"
	"fmt"
)


type updateType struct {
	Message map[string]string
	Signature string
}

func (c *Controller) Update() (string, error) {

	if c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	ver, _, err := c.getUpdVerAndUrl()
	if err!= nil {
		return "", utils.ErrInfo(err)
	}
	return utils.JsonAnswer(ver, "success").String(), nil
}

func (c *Controller) getUpdVerAndUrl() (string, string, error) {

	update, err := utils.GetHttpTextAnswer("http://dcoin.club/update.json")
	if len(update) > 0 {

		fmt.Println(update)

		updateData := new(updateType)
		err = json.Unmarshal([]byte(update), &updateData)
		if err != nil {
			return "", "", utils.ErrInfo(err)
		}

		fmt.Println(updateData)

		dataJson, err := json.Marshal(updateData.Message)
		if err != nil {
			return "", "", utils.ErrInfo(err)
		}

		pub, err := utils.BinToRsaPubKey(utils.HexToBin(consts.ALERT_KEY))
		if err != nil {
			return "", "", utils.ErrInfo(err)
		}
		fmt.Println(updateData.Signature)
		fmt.Println(string(dataJson))
		err = rsa.VerifyPKCS1v15(pub, crypto.SHA1, utils.HashSha1(string(dataJson)), []byte(utils.HexToBin(updateData.Signature)))
		if err != nil {
			return "", "", utils.ErrInfo(err)
		}

		fmt.Println(runtime.GOOS+"_"+runtime.GOARCH)
		fmt.Println(updateData.Message["dcoin_"+runtime.GOOS+"_"+runtime.GOARCH])
		if len(updateData.Message["dcoin_"+runtime.GOOS+"_"+runtime.GOARCH]) > 0 && version.Compare(updateData.Message["version"], consts.VERSION, ">") {
			newVersion := strings.Replace(c.Lang["new_version"], "[ver]", updateData.Message["version"], -1)
			return newVersion, updateData.Message[runtime.GOOS+"_"+runtime.GOARCH], nil
		}
	}
	return "", "", nil
}