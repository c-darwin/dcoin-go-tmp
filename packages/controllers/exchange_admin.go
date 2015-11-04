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
	Lock int64
	Users []map[string]string
}

func (c *Controller) ExchangeAdmin() (string, error) {

	if !c.NodeAdmin || c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	log.Debug("c.Parameters", c.Parameters)
	withdrawId := utils.StrToInt64(c.Parameters["withdraw_id"])
	if withdrawId > 0 {
		err := c.ExecSql(`UPDATE e_withdraw SET close_time = ? WHERE id = ?`, utils.Time(), withdrawId)
		if err!=nil {
			return "", utils.ErrInfo(err)
		}
	}

	lock, err := c.Single(`SELECT time FROM e_reduction_lock`).Int64()
	if len(c.Parameters["change_reduction_lock"]) > 0 {
		if lock > 0 {
			err := c.ExecSql(`DELETE FROM e_reduction_lock`)
			if err!=nil {
				return "", utils.ErrInfo(err)
			}
		} else {
			err := c.ExecSql(`INSERT INTO e_reduction_lock (time) VALUES (?)`, utils.Time())
			if err!=nil {
				return "", utils.ErrInfo(err)
			}
			lock =  utils.Time()
		}

	}

	withdraw, err := c.GetAll(`SELECT e_withdraw.id, open_time, close_time, e_users.user_id, currency_id, account, amount,  wd_amount, method, email
    		FROM e_withdraw
    		LEFT JOIN e_users on e_users.id = e_withdraw.user_id
   			ORDER BY id DESC`, 100)
	if err!=nil {
		return "", utils.ErrInfo(err)
	}

	users, err := c.GetAll(`SELECT id, email, ip, lock, user_id
    		FROM e_users
   			ORDER BY id DESC`, 100)
	if err!=nil {
		return "", utils.ErrInfo(err)
	}

	TemplateStr, err := makeTemplate("exchange_admin", "exchangeAdmin", &ExchangeAdminPage{
		Alert:        c.Alert,
		Lang:         c.Lang,
		Withdraw: withdraw,
		Lock: lock,
		Users: users,
		UserId:       c.SessUserId})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
