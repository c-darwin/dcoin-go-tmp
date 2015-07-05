package controllers
import (
	"utils"
	"log"
	"errors"
	"io/ioutil"
	"consts"
	"encoding/json"
)

type nodeConfigPage struct {
	Alert string
	SignData string
	ShowSignData bool
	CountSignArr []int
	Config map[string]string
	WaitingList []map[string]string
	MyStatus string
	MyMode string
	ConfigIni string
	UserId int64
	Lang map[string]string
	Users []map[int64]map[string]string
}

func (c *Controller) NodeConfigControl() (string, error) {

	if !c.NodeAdmin || c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	log.Println("c.Parameters", c.Parameters)
	if _, ok := c.Parameters["save_config"]; ok {
		err := c.ExecSql("UPDATE config SET in_connections_ip_limit = ?, in_connections = ?, out_connections = ?, cf_url = ?, pool_url = ?, pool_admin_user_id = ?, exchange_api_url = ?, auto_reload = ?", c.Parameters["in_connections_ip_limit"], c.Parameters["in_connections"], c.Parameters["out_connections"] , c.Parameters["cf_url"], c.Parameters["pool_url"], c.Parameters["pool_admin_user_id"], c.Parameters["exchange_api_url"],  c.Parameters["auto_reload"])
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	}

	if _, ok := c.Parameters["switch_pool_mode"]; ok {
		dq := `"`;
		if c.ConfigIni["db_type"] == "mysql" {
			dq = ``
		}
		if !c.Community{ // сингл-мод

			// переключаемся в пул-мод
			myUserId, err := c.GetMyUserId("")
			for _, table := range consts.MyTables {

				err = c.ExecSql("ALTER TABLE "+dq+table+dq+" RENAME TO "+dq+utils.Int64ToStr(myUserId)+"_"+table+dq)
				if err != nil {
					return "", utils.ErrInfo(err)
				}
			}
			err = c.ExecSql("INSERT INTO community (user_id) VALUES (?)", myUserId)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			commission, err := c.Single("SELECT commission FROM commission WHERE user_id = ?", myUserId).String()
			if err != nil {
				return "", utils.ErrInfo(err)
			}

			err = c.ExecSql("UPDATE config SET pool_admin_user_id = ?, pool_max_users = 100, commission = ?", myUserId, commission)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		} else {
			communityUsers := c.CommunityUsers
			jsonData, _ := json.Marshal(communityUsers)
			backup_community, err := c.Single("SELECT data FROM backup_community").String()
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			if len(backup_community) > 0 {
				err := c.ExecSql("UPDATE backup_community SET data = ?", jsonData)
				if err != nil {
					return "", utils.ErrInfo(err)
				}
			} else {
				err = c.ExecSql("INSERT INTO backup_community (data) VALUES (?)", jsonData)
				if err != nil {
					return "", utils.ErrInfo(err)
				}
			}
			myUserId, err := c.GetPoolAdminUserId()
			for _, table := range consts.MyTables {
				err = c.ExecSql("ALTER TABLE "+dq+utils.Int64ToStr(myUserId)+"_"+table+dq+" RENAME TO "+dq+table+dq)
				if err != nil {
					return "", utils.ErrInfo(err)
				}
			}
			err = c.ExecSql("DELETE FROM community")
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}
	}
	scriptName, err := c.Single("SELECT script_name FROM main_lock").String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	myStatus := "ON"
	if scriptName == "my_lock" {
		myStatus = "OFF"
	}
	myMode := "Single";
	if c.Community {
		myMode = "Pool"
	}
	configIni, err := ioutil.ReadFile("config.ini")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	config, err := c.GetNodeConfig()
	TemplateStr, err := makeTemplate("node_config", "nodeConfig", &nodeConfigPage {
		Alert: c.Alert,
		Lang: c.Lang,
		ShowSignData: c.ShowSignData,
		SignData: "",
		Config: config,
		UserId: c.SessUserId,
		MyStatus: myStatus,
		MyMode: myMode,
		ConfigIni: string(configIni),
		CountSignArr: c.CountSignArr})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
