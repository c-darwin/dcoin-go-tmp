package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"log"
)

type notificationsPage struct {
	SignData string
	ShowSignData bool
	Alert string
	Lang map[string]string
	CountSignArr []int
	MyNotifications map[string]map[string]string
	LangInt int64
	NodeAdmin bool
	Data map[string]string
}

func (c *Controller) Notifications() (string, error) {

	var err error
	data, err := c.OneRow(`
			SELECT email,
						 sms_http_get_request,
						 use_smtp,
						 smtp_server,
						 smtp_port,
						 smtp_ssl,
						 smtp_auth,
						 smtp_username,
						 smtp_password
			FROM `+c.MyPrefix+`my_table
			`).String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	myNotifications := make(map[string]map[string]string)
	myNotifications_, err := c.GetAll("SELECT * FROM "+c.MyPrefix+"my_notifications ORDER BY sort ASC", -1)
	for _, data := range myNotifications_ {
		myNotifications[data["name"]] = map[string]string {"email": data["email"], "sms": data["sms"], "important": data["important"]}
	}
	log.Println("myNotifications", myNotifications)

	TemplateStr, err := makeTemplate("notifications", "notifications", &notificationsPage{
		Alert: c.Alert,
		Lang: c.Lang,
		CountSignArr: c.CountSignArr,
		ShowSignData: c.ShowSignData,
		SignData: "",
		MyNotifications: myNotifications,
		NodeAdmin: c.NodeAdmin,
		LangInt: c.LangInt,
		Data: data})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}


