package controllers
import (
	"html/template"
	"bytes"
	"github.com/c-darwin/dcoin-go-tmp/packages/static"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"runtime"
	"github.com/c-darwin/dcoin-go-tmp/packages/consts"
)

type updatingBlockchainStruct struct {
	Lang map[string]string
	WaitText string
	BlockTime int64
	BlockId int64
	StartDaemons string
	BlockMeter int64
	CheckTime string
}


func (c *Controller) UpdatingBlockchain() (string, error) {

	funcMap := template.FuncMap{
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
	}
	data, err := static.Asset("static/templates/updating_blockchain.html")
	t := template.New("template").Funcs(funcMap)
	t, err = t.Parse(string(data))
	if err != nil {
		return "", err
	}

	var blockTime, blockId int64
	var waitText, startDaemons, checkTime string
	ConfirmedBlockId, err := c.DCDB.GetConfirmedBlockId()
	if err != nil {
		return "", err
	}
	if ConfirmedBlockId == 0 {
		firstLoadBlockchain, err := c.DCDB.Single("SELECT first_load_blockchain FROM config").String()
		if err != nil {
			return "", err
		}
		if firstLoadBlockchain == "file" {
			waitText = c.Lang["loading_blockchain_please_wait"]
		} else {
			waitText = c.Lang["is_synchronized_with_the_dc_network"]
		}
	} else {
		LastBlockData, err := c.DCDB.GetLastBlockData()
		if err != nil {
			return "", err
		}
		blockTime = LastBlockData["lastBlockTime"]
		blockId = LastBlockData["blockId"]
	}

	// для сингл-мода, кнопка включения и выключения демонов
	if !c.Community {
		lockName, err := c.DCDB.GetMainLockName()
		if err != nil {
			return "", err
		}
		if lockName == "main_lock" {
			startDaemons =`<a href="#" id="start_daemons" style="color:#C90600">Start daemons</a>`
		}
		// инфа о синхронизации часов
		switch runtime.GOOS {
		case "linux":
			checkTime = c.Lang["check_time_nix"]
		case "windows":
			checkTime = c.Lang["check_time_win"]
		case "darwin":
			checkTime = c.Lang["check_time_mac"]
		default:
			checkTime = c.Lang["check_time_nix"]
		}
		checkTime = c.Lang["check_time"] + checkTime
	}

	blockMeter := int64(utils.Round(float64((blockId/consts.LAST_BLOCK)*100), 0))
	b := new(bytes.Buffer)
	t.Execute(b, &updatingBlockchainStruct{Lang: c.Lang, WaitText: waitText, BlockId: blockId, BlockTime: blockTime, StartDaemons: startDaemons, BlockMeter:blockMeter, CheckTime: checkTime})
	return b.String(), nil
}
