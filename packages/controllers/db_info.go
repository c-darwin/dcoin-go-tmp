package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

type DbInfoPage struct {
	TimeNow               string
	NodesBan              []map[string]string
	NodesConnection       []map[string]string
	MainLock              []map[string]string
	Variables             map[string]string
	QueueTx               int64
	TransactionsTestblock int64
	Transactions          int64
	Lang                  map[string]string
}

func (c *Controller) DbInfo() (string, error) {

	var err error

	timeNow := utils.TimeF(c.TimeFormat)

	nodesBan, err := c.GetAll(`
			SELECT nodes_ban.ban_start,
						  nodes_ban.user_id,
						  miners_data.tcp_host,
						  nodes_ban.info
			FROM nodes_ban
			LEFT JOIN miners_data ON miners_data.user_id = nodes_ban.user_id
			ORDER BY ban_start
			`, -1)

	nodesConnection, err := c.GetAll(`
			SELECT *
			FROM nodes_connection
			`, -1)

	mainLock, err := c.GetAll(`
			SELECT *
			FROM main_lock
			`, -1)

	queueTx, err := c.Single("SELECT count(*) FROM queue_tx").Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	transactionsTestblock, err := c.Single("SELECT count(*) FROM transactions_testblock").Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	transactions, err := c.Single("SELECT count(*) FROM transactions").Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	variables, err := c.GetMap("SELECT name, value	FROM variables", "name", "value")
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	TemplateStr, err := makeTemplate("db_info", "dbInfo", &DbInfoPage{
		Lang:                  c.Lang,
		TimeNow:               timeNow,
		NodesBan:              nodesBan,
		NodesConnection:       nodesConnection,
		MainLock:              mainLock,
		Variables:             variables,
		QueueTx:               queueTx,
		TransactionsTestblock: transactionsTestblock,
		Transactions:          transactions})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
