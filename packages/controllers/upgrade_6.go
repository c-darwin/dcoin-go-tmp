package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"strings"
	"net"
	"time"
)

type upgrade6Page struct {
	Alert string
	UserId int64
	Lang map[string]string
	HttpHost string
	TcpHost string
	SaveAndGotoStep string
	UpgradeMenu string
	Community bool
	HostType string
}

func (c *Controller) Upgrade6() (string, error) {

	log.Debug("Upgrade6")

	var hostType string

	hostData, err := c.OneRow("SELECT http_host, tcp_host FROM "+c.MyPrefix+"my_table").String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// в режиме пула выдаем только хост ноды
	log.Debug("c.Community: %v", c.Community)
	log.Debug("c.PoolAdminUserId: %v", c.PoolAdminUserId)
	if c.Community /*&& len(data["http_host"]) == 0 && len(data["tcp_host"]) == 0*/ {
		hostType = "pool"
		hostData, err = c.OneRow("SELECT http_host, tcp_host FROM miners_data WHERE user_id  =  ?", c.PoolAdminUserId).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	} else {
		// если смогли подключиться из вне
		ip, err := utils.GetHttpTextAnswer("https://api.ipify.org");
		/*httpHost, err := c.Single("SELECT http_host FROM "+c.MyPrefix+"my_table").String()
		if err!=nil {
			return "", utils.ErrInfo(err)
		}
		port := "8089"
		if len(httpHost) > 0 {
			re := regexp.MustCompile(`https?:\/\/(?:[0-9a-z\_\.\-]+):([0-9]+)`)
			match := re.FindStringSubmatch(httpHost)
			if len(match) != 0 {
				port = match[1];
			}
		}*/
		conn, err := net.DialTimeout("tcp", ip+":8089", 5 * time.Second)
		if err != nil {
			// если не смогли подключиться, то в JS будем искать рабочий пул и региться на нем. и дадим юзеру указать другие хост:ip
			hostType = "findPool"


		} else {
			hostType = "normal"
			defer conn.Close()
			hostData["http_host"] = ip+":8089"
			hostData["tcp_host"] = ip+":8088"
		}
	}


	saveAndGotoStep := strings.Replace(c.Lang["save_and_goto_step"], "[num]", "7", -1)
	upgradeMenu := utils.MakeUpgradeMenu(6)

	TemplateStr, err := makeTemplate("upgrade_6", "upgrade6", &upgrade6Page {
		Alert: c.Alert,
		Lang: c.Lang,
		SaveAndGotoStep: saveAndGotoStep,
		UpgradeMenu: upgradeMenu,
		HttpHost: hostData["http_host"],
		TcpHost: hostData["tcp_host"],
		Community: c.Community,
		HostType: hostType,
		UserId: c.SessUserId})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

