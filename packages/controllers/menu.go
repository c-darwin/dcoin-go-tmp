package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"github.com/c-darwin/dcoin-go-tmp/packages/static"
	"html/template"
	"bytes"
	"strings"

)

type menuPage struct {
	MyModalIdName string
	SetupPassword bool
	Lang map[string]string
	LangInt int64
	PoolAdmin bool
	Community bool
	MinerId int64
	Name string
	UserId int64
	DaemonsStatus string
	MyNotice map[string]string
	BlockId int64
	Avatar string
	NoAvatar string
	FaceUrls string
	Restricted int64
}

func (c *Controller) Menu() (string, error) {

	if !c.dbInit || c.SessUserId==0 {
		return "", nil
	}

	var name, avatar string
	if c.SessUserId > 0 {
		data, err := c.OneRow("SELECT name, avatar FROM users WHERE user_id =  ?", c.SessUserId).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		name, avatar = data["name"], data["avatar"]
	}

	if len(name) == 0 {
		miner, err := c.Single("SELECT miner_id FROM miners_data WHERE user_id  =  ?", c.SessUserId).Int64()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if miner > 0 {
			name = "ID "+utils.Int64ToStr(c.SessUserId)+" (miner)"
		} else {
			name = "ID "+utils.Int64ToStr(c.SessUserId)
		}
	}

	var face_urls []string
	if len(avatar) == 0 {
		data, err := c.OneRow("SELECT photo_block_id, photo_max_miner_id, miners_keepers FROM miners_data WHERE user_id  =  ?", c.SessUserId).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if len(data) > 0 {
			// получим ID майнеров, у которых лежат фото нужного нам юзера
			minersIds := utils.GetMinersKeepers(data["photo_block_id"], data["photo_max_miner_id"], data["miners_keepers"], true);
			if len(minersIds) > 0 {
				hosts, err := c.GetList("SELECT http_host as host FROM miners_data WHERE miner_id IN ("+utils.JoinInts(minersIds, ",")+")").String()
				if err != nil {
					return "", utils.ErrInfo(err)
				}
				for i:=0; i < len(hosts); i++ {
					face_urls = append(face_urls, hosts[i]+"public/face_"+utils.Int64ToStr(c.SessUserId)+".jpg")
				}
			}
		}
	}

	noAvatar := "static/img/noavatar.png"
	minerId, err := c.GetMyMinerId(c.SessUserId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// ID блока вверху
	blockId, err := c.GetBlockId()

	daemonsStatus := ""
	if !c.Community {
		scriptName, err := c.Single("SELECT script_name FROM main_lock").String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if (scriptName == "my_lock") {
			daemonsStatus = `<li title="`+c.Lang["daemons_status_off"]+`"><a href="#" id="start_daemons" style="color:#C90600"><i class="fa fa-power-off" style="font-size: 20px"></i></a></li>`
		} else {
			daemonsStatus = `<li title="`+c.Lang["daemons_status_on"]+`"><a href="#" id="stop_daemons" style="color:#009804"><i class="fa fa-power-off" style="font-size: 20px"></i></a></li>`
		}
	}

	data, err := static.Asset("static/templates/menu.html")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	log.Debug("menu ok")
	modal, err := static.Asset("static/templates/modal.html")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	log.Debug("modal ok")

	t := template.Must(template.New("template").Parse(string(data)))
	t = template.Must(t.Parse(string(modal)))
	b := new(bytes.Buffer)
	err = t.ExecuteTemplate(b, "menu", &menuPage{SetupPassword: false, MyModalIdName: "myModal", Lang: c.Lang, PoolAdmin: c.PoolAdmin, Community: c.Community, MinerId: minerId, Name: name, LangInt: c.LangInt, UserId: c.SessUserId, Restricted: c.SessRestricted, DaemonsStatus: daemonsStatus, MyNotice: c.MyNotice, BlockId: blockId, Avatar: avatar, NoAvatar: noAvatar, FaceUrls: strings.Join(face_urls, ",")})
	log.Debug("ExecuteTemplate")
	if err != nil {
		log.Debug("%s", utils.ErrInfo(err))
		return "", utils.ErrInfo(err)
	}
	log.Debug("b.String():\n %s", b.String())
	return b.String(), nil
}
