package controllers
import (
	"utils"
	"log"
	"errors"
	"regexp"
	"encoding/json"
	"strings"
)

func (c *Controller) GenerateNewPrimaryKey() (string, error) {

	if c.SessRestricted!=0 {
		return "", errors.New("Permission denied")
	}

	c.r.ParseForm()
	if ok, _ := regexp.MatchString(`install`, c.r.FormValue("tpl_name")); ok {
		return "", nil
	}

	show := false
	// проверим, есть ли сообщения от админа
	data, err := c.OneRow("SELECT * FROM alert_messages WHERE close  =  0 ORDER BY id DESC").String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	var message map[string]string
	var adminMessage string
	if len(data["message"]) > 0 {
		err = json.Unmarshal([]byte(data["message"]), &message)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if len(message[c.LangInt]) {
			adminMessage = message[c.LangInt]
		} else {
			adminMessage = message["gen"]
		}
		if data["currency_list"] != "ALL" {
			// проверим, есть ли у нас обещанные суммы с такой валютой
			promisedAmount, err := c.Single("SELECT id FROM promised_amount WHERE currency_id IN (?)", data["currency_list"]).Int64()
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			if promisedAmount > 0 {
				show = true
			}
		} else {
			show = true
		}
	}
	result := ""
	if show {
		result += `<script>
			$('#close_alert').bind('click', function () {
				$.post( 'ajax?controllerName=closeAlert', {
					'id' : '`+data["id"]+`'
				} );
			});
			</script>
			 <div class="alert alert-danger alert-dismissable" style='margin-top: 30px'><button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>
			 <h4>Warning!</h4>
			    `+data["message"]+`
			  </div>`
	}

	// сообщение о новой версии движка
	myVer, err := c.Single("SELECT current_version FROM info_block").String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// возможны 2 сценария:
	// 1. информация о новой версии есть в блоках
	newVer, err := c.GetMap("SELECT version FROM new_version WHERE alert  =  1").String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	newMaxVer := "0"
	for i:=0; i < len(newVer); i++ {
		if VersionOrdinal(newVer[i]) > VersionOrdinal(myVer) && newMaxVer == "0" {
			newMaxVer = newVer[i]
		}
		if VersionOrdinal(newVer[i]) > VersionOrdinal(newMaxVer) && newMaxVer != "0" {
			newMaxVer = newVer[i]
		}
	}
	var newVersion string
	if newMaxVer != "0" {
		newVersion = strings.Replace(c.Lang["new_version"], "[ver]", newMaxVer, -1)
	}

	// для пулов и ограниченных юзеров выводим сообщение без кнопок
	if (c.Community || c.SessRestricted!=0) && newMaxVer != "0" {
		newVersion = strings.Replace(c.Lang["new_version_pulls"], "[ver]", newMaxVer, -1)
	}

	if newMaxVer != "0" && len(myVer) > 0 {
		result += `<script>
				$('#btn_install').bind('click', function () {
					$('#new_version_text').text('Please wait');
					$.post( 'ajax?controllerName=installNewVersion', {}, function(data) {
						$('#new_version_text').text(data);
					});
				});
				$('#btn_upgrade').bind('click', function () {
					$('#new_version_text').text('Please wait');
					$.post( 'ajax?controllerName=upgradeToNewVersion', {}, function(data) {
						$('#new_version_text').text(data);
					});
				});
			</script>

			<div class="alert alert-danger alert-dismissable" style='margin-top: 30px'><button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>
			    <h4>Warning!</h4>
			   <div id='new_version_text'>`+c.Lang["new_version"]+`</div>
			 </div>`
	}

	if c.SessRestricted == 0 && !c.Community {
		
	}

		log.Println(json)
	return string(json), nil
}

