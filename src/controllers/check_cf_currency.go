package controllers
import (
	"errors"
	"utils"
)

func (c *Controller) CheckCfCurrency() (string, error) {

	if c.SessUserId == 0 || c.SessRestricted != 0 {
		return "", errors.New("Permission denied")
	}

	c.r.ParseForm()
	projectCurrencyName := c.r.FormValue("project_currency_name")
	if !utils.CheckInputData(projectCurrencyName, "cf_currency_name")  {
		return "", errors.New("incorrect project_currency_name")
	}

	// проверим, не занято ли имя валюты
	currency, err := c.Single("SELECT id FROM cf_projects WHERE project_currency_name = ? AND close_block_id = 0 AND del_block_id = 0", projectCurrencyName).Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	if currency > 0 {
		return `{"error":"`+c.Lang["currency_name_busy"]+`"}`, nil
	}

	// проверим, не занято ли имя валюты
	currency, err = c.Single("SELECT id FROM cf_currency WHERE name = ?", projectCurrencyName).Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	if currency > 0 {
		return `{"error":"`+c.Lang["currency_name_busy"]+`"}`, nil
	}

	return `{"success":"`+c.Lang["name_is_not_occupied"]+`"}`, nil
}
