package controllers
import (
	"utils"
	"log"
	"strings"
)

type VotingPage struct {
	Alert string
	SignData string
	ShowSignData bool
	UserId int64
	Lang map[string]string
	CountSignArr []int
}

func (c *Controller) Voting() (string, error) {

	log.Println("Voting")

	txType := "votes_complex";
	txTypeId := utils.TypeInt(txType)
	timeNow := utils.Time()

	waitVoting := make(map[int64]string)

	// голосовать майнер может только после того, как пройдет  miner_newbie_time сек
	regTime, err := c.Single("SELECT reg_time FROM miners_data WHERE user_id  =  ?", c.SessUserId).Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	minerNewbie := ""
	if regTime > utils.Time() - c.Variables.Int64["miner_newbie_time"] && c.SessUserId !=1 {
		minerNewbie = strings.Replace(c.Lang["hold_time_wait2"], "[sec]", utils.TimeLeft(c.Variables.Int64["miner_newbie_time"] - (utils.Time() - regTime), c.Lang), -1)
	} else {
		// валюты
		rows, err := c.Query(c.FormatQuery(`
				SELECT currency_id,
							  name,
							  full_name,
							  start_time
				FROM promised_amount
					LEFT JOIN currency ON currency.id = promised_amount.currency_id
				WHERE user_id = ? AND
							 status IN ('mining', 'repaid') AND
							 start_time > 0 AND
							 del_block_id = 0
				GROUP BY currency_id
				`), c.SessUserId)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		defer rows.Close()
		for rows.Next() {
			var currency_id, start_time int64
			var name, full_name string
			err = rows.Scan(&currency_id, &name, &full_name, &start_time)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			// после добавления обещанной суммы должно пройти не менее min_hold_time_promise_amount сек, чтобы за неё можно было голосовать
			if start_time > utils.Time() - c.Variables["min_hold_time_promise_amount"] {
				waitVoting[currency_id] = strings.Replace(c.Lang["hold_time_wait"], "[sec]", utils.TimeLeft(c.Variables.Int64["min_hold_time_promise_amount"] - (utils.Time() - start_time), c.Lang), -1)
				continue
			}
			// если по данной валюте еще не набралось >1000 майнеров, то за неё голосовать нельзя.
			countMiners, err := c.Single(`
					SELECT count(user_id)
					FROM promised_amount
					WHERE start_time < ? AND
								 del_block_id = 0 AND
								 status IN ('mining', 'repaid') AND
								 currency_id = ? AND
								 del_block_id = 0
					GROUP BY  user_id
					`, currency_id).Int64()
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}
	}


	TemplateStr, err := makeTemplate("voting", "Voting", &VotingPage{
		Alert: c.Alert,
		Lang: c.Lang,
		CountSignArr: c.CountSignArr,
		ShowSignData: c.ShowSignData,
		UserId: c.SessUserId,
		SignData: ""})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

