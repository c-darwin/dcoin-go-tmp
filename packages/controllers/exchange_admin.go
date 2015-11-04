package controllers

import (
	"errors"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

type ExchangeAdminPage struct {
	Alert        string
	UserId       int64
	Lang         map[string]string
	Withdraw []map[string]string
}

func (c *Controller) ExchangeAdmin() (string, error) {

	if !c.NodeAdmin || c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	log.Debug("c.Parameters", c.Parameters)

	withdraw, err := c.GetAll(`SELECT e_withdraw.id, open_time, close_time, user_id, currency_id, account, amount,  wd_amount, method, email
    		FROM e_withdraw
    		LEFT JOIN e_users on e_users.id = e_withdraw.user_id
   			ORDER BY open_time DESC`, 100)
	/*rows, err := c.Query(c.FormatQuery(`
			SELECT e_withdraw.id, open_time, close_time, user_id, currency_id, account, amount,  wd_amount, method, email
    		FROM e_withdraw
    		LEFT JOIN e_users on e_users.id = e_withdraw.user_id
   			ORDER BY open_time DESC
			`))
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, open_time, close_time, user_id, currency_id int64
		var account, method, email string
		var amount, wd_amount float64
		err = rows.Scan(&id, &open_time, &close_time, &user_id, &currency_id, &account, &amount,  &wd_amount, &method, &email)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	}*/


	TemplateStr, err := makeTemplate("exchange_admin", "exchangeAdmin", &ExchangeAdminPage{
		Alert:        c.Alert,
		Lang:         c.Lang,
		Withdraw: withdraw,
		UserId:       c.SessUserId})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
