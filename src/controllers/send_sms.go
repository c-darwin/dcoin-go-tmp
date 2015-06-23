package controllers
import (
//	"utils"
	"errors"
	"io/ioutil"
	"net/http"
	"fmt"
	"encoding/json"
)
func (c *Controller) SendSms() (string, error) {

	if c.SessRestricted != 0 || !c.NodeAdmin {
		return "", errors.New("Permission denied")
	}

	c.r.ParseForm()
	text := c.r.FormValue("text")

	sms_http_get_request, err := c.Single("SELECT sms_http_get_request FROM "+c.MyPrefix+"my_table").String()
	if err != nil {
		result, _ := json.Marshal(map[string]string{"error": fmt.Sprintf(`%s`, err)})
		return string(result), nil
	}
	resp, err := http.Get(sms_http_get_request+text)
	if err != nil {
		result, _ := json.Marshal(map[string]string{"error": fmt.Sprintf(`%s`, err)})
		return string(result), nil
	}
	defer resp.Body.Close()
	htmlData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		result, _ := json.Marshal(map[string]string{"error": fmt.Sprintf(`%s`, err)})
		return string(result), nil
	}
	result, _ := json.Marshal(map[string]string{"success": string(htmlData)})
	return string(result), nil

}

