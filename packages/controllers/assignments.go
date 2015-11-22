package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/consts"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"strings"
)

type AssignmentsPage struct {
	Alert              string
	SignData           string
	ShowSignData       bool
	TxType             string
	TxTypeId           int64
	TimeNow            int64
	UserId             int64
	Lang               map[string]string
	CountSignArr       []int
	CurrencyList       map[int64]string
	CashRequestsStatus map[string]string
	MyCashRequests     []map[string]string
	ActualData         map[string]string
	MainQuestion       string
	NewPromiseAmount   string
	MyRace             string
	MyCountry          string
	ExamplePoints      map[string]string
	VideoHost          string
	PhotoHosts         []string
	PromisedAmountData map[string]string
	UserInfo           map[string]string
	CloneHosts         map[int64][]string
}

func (c *Controller) Assignments() (string, error) {

	var randArr []int64
	// Нельзя завершить голосование юзеров раньше чем через сутки, даже если набрано нужное кол-во голосов.
	// В голосовании нодов ждать сутки не требуется, т.к. там нельзя поставить поддельных нодов

	// Модерация новых майнеров
	// берем тех, кто прошел проверку нодов (type='node_voting')
	num, err := c.Single("SELECT count(id) FROM votes_miners WHERE votes_end  =  0 AND type  =  'user_voting'").Int64()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	if num > 0 {
		randArr = append(randArr, 1)
	}

	// Модерация promised_amount
	// вначале получим ID валют, которые мы можем проверять.
	currency, err := c.GetList("SELECT currency_id FROM promised_amount WHERE status IN ('mining', 'repaid') AND user_id = ?", c.SessUserId).String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	addSql := ""
	currencyIds := strings.Join(currency, ",")
	if len(currencyIds) > 0 || c.SessUserId == 1 {
		if c.SessUserId != 1 {
			addSql = "AND currency_id IN (" + currencyIds + ")"
		}
		num, err := c.Single("SELECT count(id) FROM promised_amount WHERE status  =  'pending' AND del_block_id  =  0 " + addSql + "").Int64()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if num > 0 {
			randArr = append(randArr, 2)
		}
	}

	var AssignType int64
	if len(randArr) > 0 {
		AssignType = randArr[utils.RandInt(0, len(randArr))]
	}

	cloneHosts := make(map[int64][]string)
	var photoHosts []string
	examplePoints := make(map[string]string)
	tplName := "assignments"
	tplTitle := "assignments"

	var txType string
	var txTypeId int64
	var timeNow int64
	var myRace, myCountry, mainQuestion, newPromiseAmount, videoHost string
	var promisedAmountData, userInfo map[string]string

	switch AssignType {
	case 1:

		// ***********************************
		// задания по модерации новых майнеров
		// ***********************************
		txType = "VotesMiner"
		txTypeId = utils.TypeInt(txType)
		timeNow = utils.Time()

		userInfo, err = c.OneRow(`
				SELECT miners_data.user_id,
							 id as vote_id,
							 face_coords,
							 profile_coords,
							 video_type,
							 video_url_id,
							 photo_block_id,
							 photo_max_miner_id,
							 miners_keepers,
							 http_host
				FROM votes_miners
				LEFT JOIN miners_data
						 ON miners_data.user_id = votes_miners.user_id
				WHERE votes_end = 0 AND
							 type = 'user_voting'
				`).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		// проверим, не голосовали ли мы за это в последние 30 минут
		repeated, err := c.Single("SELECT id FROM "+c.MyPrefix+"my_tasks WHERE type  =  'miner' AND id  =  ? AND time > ?", userInfo["vote_id"], utils.Time()-consts.ASSIGN_TIME).Int64()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if repeated > 0 {
			tplName = "assignments"
			break
		}
		examplePoints, err = c.GetPoints(c.Lang)
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		// получим ID майнеров, у которых лежат фото нужного нам юзера
		minersIds := utils.GetMinersKeepers(userInfo["photo_block_id"], userInfo["photo_max_miner_id"], userInfo["miners_keepers"], true)
		if len(minersIds) > 0 {
			photoHosts, err = c.GetList("SELECT http_host FROM miners_data WHERE miner_id  IN (" + utils.JoinInts(minersIds, ",") + ")").String()
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}

		// отрезки майнера, которого проверяем
		relations, err := c.OneRow("SELECT * FROM faces WHERE user_id  =  ?", userInfo["user_id"]).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		// получим допустимые расхождения между точками и совместимость версий
		data_, err := c.OneRow("SELECT tolerances, compatibility FROM spots_compatibility").String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		tolerances := make(map[string]map[string]string)
		if err := json.Unmarshal([]byte(data_["tolerances"]), &tolerances); err != nil {
			return "", utils.ErrInfo(err)
		}
		var compatibility []int
		if err := json.Unmarshal([]byte(data_["compatibility"]), &compatibility); err != nil {
			return "", utils.ErrInfo(err)
		}

		// формируем кусок SQL-запроса для соотношений отрезков
		addSqlTolerances := ""
		typesArr := []string{"face", "profile"}
		for i := 0; i < len(typesArr); i++ {
			for j := 1; j <= len(tolerances[typesArr[i]]); j++ {
				currentRelations := utils.StrToFloat64(relations[typesArr[i][:1]+utils.IntToStr(j)])
				diff := utils.StrToFloat64(tolerances[typesArr[i]][utils.IntToStr(j)]) * currentRelations
				if diff == 0 {
					continue
				}
				min := currentRelations - diff
				max := currentRelations + diff
				addSqlTolerances += typesArr[i][:1] + utils.IntToStr(j) + ">" + utils.Float64ToStr(min) + " AND " + typesArr[i][:1] + utils.IntToStr(j) + " < " + utils.Float64ToStr(max) + " AND "
			}
		}
		addSqlTolerances = addSqlTolerances[:len(addSqlTolerances)-4]

		// формируем кусок SQL-запроса для совместимости версий
		addSqlCompatibility := ""
		for i := 0; i < len(compatibility); i++ {
			addSqlCompatibility += fmt.Sprintf(`%d,`, compatibility[i])
		}
		addSqlCompatibility = addSqlCompatibility[:len(addSqlCompatibility)-1]

		// получаем из БД похожие фото
		rows, err := c.Query(c.FormatQuery(`
				SELECT miners_data.user_id,
							 photo_block_id,
							 photo_max_miner_id,
							 miners_keepers
				FROM faces
				LEFT JOIN miners_data ON
						miners_data.user_id = faces.user_id
				WHERE `+addSqlTolerances+` AND
							version IN (`+addSqlCompatibility+`) AND
				             faces.status = 'used' AND
				             miners_data.user_id != ?
				LIMIT 0, 100
				`), userInfo["user_id"])
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		defer rows.Close()
		for rows.Next() {
			var photo_block_id, photo_max_miner_id, miners_keepers string
			var user_id int64
			err = rows.Scan(&user_id, &photo_block_id, &photo_max_miner_id, &miners_keepers)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			// майнеры, у которых можно получить фото нужного нам юзера
			minersIds := utils.GetMinersKeepers(photo_block_id, photo_max_miner_id, miners_keepers, true)
			if len(minersIds) > 0 {
				photoHosts, err = c.GetList("SELECT http_host FROM miners_data WHERE miner_id  IN (" + utils.JoinInts(minersIds, ",") + ")").String()
				if err != nil {
					return "", utils.ErrInfo(err)
				}
			}
			cloneHosts[user_id] = photoHosts
		}

		data, err := c.OneRow("SELECT race, country FROM " + c.MyPrefix + "my_table").Int64()
		myRace = c.Races[data["race"]]
		myCountry = consts.Countries[int(data["country"])]

		tplName = "assignments_new_miner"
		tplTitle = "assignmentsNewMiner"

	case 2:
		promisedAmountData, err = c.OneRow(`
				SELECT id,
							 currency_id,
							 amount,
							 user_id,
							 video_type,
							 video_url_id
				FROM promised_amount
				WHERE status =  'pending' AND
							 del_block_id = 0
				` + addSql + `
		`).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		promisedAmountData["currency_name"] = c.CurrencyList[utils.StrToInt64(promisedAmountData["currency_id"])]

		// проверим, не голосовали ли мы за это в последние 30 минут
		repeated, err := c.Single("SELECT id FROM "+c.MyPrefix+"my_tasks WHERE type  =  'promised_amount' AND id  =  ? AND time > ?", promisedAmountData["id"], utils.Time()-consts.ASSIGN_TIME).Int64()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		if repeated > 0 {
			tplName = "assignments"
			tplTitle = "assignments"
			break
		}

		// если нету видео на ютубе, то получаем host юзера, где брать видео
		if promisedAmountData["video_url_id"] == "null" {
			videoHost, err = c.Single("SELECT http_host FROM miners_data WHERE user_id  =  ?", promisedAmountData["user_id"]).String()
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}

		// каждый раз обязательно проверяем, где находится юзер
		userInfo, err = c.OneRow(`
				SELECT latitude,
							 user_id,
							 longitude,
							 photo_block_id,
							 photo_max_miner_id,
							 miners_keepers,
							 http_host
				FROM miners_data
				WHERE user_id = ?
				`, promisedAmountData["user_id"]).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		// получим ID майнеров, у которых лежат фото нужного нам юзера
		minersIds := utils.GetMinersKeepers(userInfo["photo_block_id"], userInfo["photo_max_miner_id"], userInfo["miners_keepers"], true)
		if len(minersIds) > 0 {
			photoHosts, err = c.GetList("SELECT http_host FROM miners_data WHERE miner_id  IN (" + utils.JoinInts(minersIds, ",") + ")").String()
			if err != nil {
				return "", utils.ErrInfo(err)
			}
		}

		txType = "VotesPromisedAmount"
		txTypeId = utils.TypeInt(txType)
		timeNow = utils.Time()

		newPromiseAmount = strings.Replace(c.Lang["new_promise_amount"], "[amount]", promisedAmountData["amount"], -1)
		newPromiseAmount = strings.Replace(newPromiseAmount, "[currency]", promisedAmountData["currency_name"], -1)

		mainQuestion = strings.Replace(c.Lang["main_question"], "[amount]", promisedAmountData["amount"], -1)
		mainQuestion = strings.Replace(mainQuestion, "[currency]", promisedAmountData["currency_name"], -1)

		tplName = "assignments_promised_amount"
		tplTitle = "assignmentsPromisedAmount"

	default:
		tplName = "assignments"
		tplTitle = "assignments"
	}

	TemplateStr, err := makeTemplate(tplName, tplTitle, &AssignmentsPage{
		Alert:              c.Alert,
		Lang:               c.Lang,
		CountSignArr:       c.CountSignArr,
		ShowSignData:       c.ShowSignData,
		UserId:             c.SessUserId,
		TimeNow:            timeNow,
		TxType:             txType,
		TxTypeId:           txTypeId,
		SignData:           "",
		CurrencyList:       c.CurrencyList,
		MainQuestion:       mainQuestion,
		NewPromiseAmount:   newPromiseAmount,
		MyRace:             myRace,
		MyCountry:          myCountry,
		ExamplePoints:      examplePoints,
		VideoHost:          videoHost,
		PhotoHosts:         photoHosts,
		PromisedAmountData: promisedAmountData,
		UserInfo:           userInfo,
		CloneHosts:         cloneHosts})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
