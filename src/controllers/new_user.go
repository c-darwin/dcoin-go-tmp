package controllers
import (
	"utils"
	"log"
	"strings"
	"os"
)

type newUserPage struct {
	SignData string
	ShowSignData bool
	TxType string
	TxTypeId int64
	TimeNow int64
	UserId int64
	Alert string
	Lang map[string]string
	CountSignArr []int
	MyRefs map[int64]myRefs
	GlobalRefs map[int64]globalRefs
	CurrencyList map[int64]string
	PoolUrl string
}

func (c *Controller) NewUser() (string, error) {

	log.Println("NewUser")

	txType := "new_user";
	txTypeId := utils.TypeInt(txType)
	timeNow := utils.Time()

	param := map[string]string{"x": 176, "y": 100, "width": 100, "bg_path": "static/img/k_bg.png"}

	myRefsKeys := make(map[string]string)
	if c.SessRestricted == 0 {
		rows, err := c.Query(c.FormatQuery(`
				SELECT user_id,	private_key,  log_id
				FROM `+c.MyPrefix+`my_new_users
				LEFT JOIN users ON users.user_id = `+c.MyPrefix+`my_new_users.user_id
				WHERE status = 'approved'
				`))
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		defer rows.Close()
		for rows.Next() {
			var user_id, log_id int64
			var private_key string
			err = rows.Scan(&user_id, &private_key, &log_id)
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			// проверим, не сменил ли уже юзер свой ключ
			if log_id > 0 {
				myRefsKeys[user_id] = map[string]string{"user_id": user_id}
			} else {
				myRefsKeys[user_id] =  map[string]string{"user_id": user_id, "private_key": private_key}
				md5:=string(utils.Md5(private_key))
				kPath := "public/"+md5[0:16]
				kPathPng := kPath+".png"
				kPathTxt := kPath+".txt"
				if _, err := os.Stat(kPathPng); os.IsNotExist(err) {
					privKey := strings.Replace(private_key, "-----BEGIN RSA PRIVATE KEY-----", "", -1)
					privKey = strings.Replace(privKey, "-----END RSA PRIVATE KEY-----", "", -1)
					/*$gd = key_to_img($private_key, $param, $row['user_id']);
				imagepng($gd, $k_path_png);
				file_put_contents($k_path_txt, trim($private_key));*/
				}
			}
		}
	}

	refs := make(map[int64]map[int64]float64)
	// инфа по рефам юзера
	rows, err := c.Query(c.FormatQuery(`
			SELECT referral, sum(amount) as amount, currency_id
			FROM referral_stats
			WHERE user_id = ?
			GROUP BY currency_id,  referral
			`), c.SessUserId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	defer rows.Close()
	for rows.Next() {
		var referral, currency_id int64
		var amount float64
		err = rows.Scan(&referral, &amount, &currency_id)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		refs[referral] = map[int64]float64{currency_id:amount}
	}

	myRefsAmounts := make(map[int64]myRefs)
	for refUserId, refData := range refs {
		data, err := c.OneRow("SELECT * FROM miners_data WHERE user_id  =  ?", refUserId).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		// получим ID майнеров, у которых лежат фото нужного нам юзера
		if len(data) == 0 {
			continue
		}
		minersIds := utils.GetMinersKeepers(data["photo_block_id"], data["photo_max_miner_id"], data["miners_keepers"], true)
		if len(minersIds) > 0 {
			hosts, err := c.GetList("SELECT host FROM miners_data WHERE miner_id  IN ("+utils.JoinInts(minersIds, ",")+")").String()
			if err != nil {
				return "", utils.ErrInfo(err)
			}
			myRefsAmounts[refUserId] = myRefs{amounts: refData, hosts: hosts}
		}
	}
	myRefs := make(map[int64]myRefs)
	for refUserId, refData := range myRefsAmounts {
		myRefs[refUserId] = refData
	}
	for refUserId, refData := range myRefsKeys {
		myRefs[refUserId].key = refData["private_key"]
		md5 := string(utils.Md5(refData["private_key"]))
		myRefs[refUserId].keyUrl = c.NodeConfig["pool_url"] + "public/"+md5[0:16]
	}

	/*
	 * Общая стата по рефам
	 */
	globalRefs := make(map[int64]globalRefs)
	// берем лидеров по USD
	rows, err := c.Query(c.FormatQuery(`
			SELECT user_id, sum(amount) as amount
			FROM referral_stats
			WHERE currency_id = 72
			GROUP BY user_id
			ORDER BY amount DESC
			`), c.SessUserId)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	defer rows.Close()
	for rows.Next() {
		var user_id int64
		var amount float64
		err = rows.Scan(&user_id, &amount)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		// вся прибыль с рефов у данного юзера
		refAmounts, err := c.GetAll(`
				SELECT ROUND(sum(amount)) as amount,  currency_id
				FROM referral_stats
				WHERE user_id = ?
				GROUP BY currency_id
				`, -1, user_id)
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		data, err := c.OneRow("SELECT * FROM miners_data WHERE user_id  =  ?", user_id).String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		// получим ID майнеров, у которых лежат фото нужного нам юзера
		minersIds := utils.GetMinersKeepers(data["photo_block_id"], data["photo_max_miner_id"], data["miners_keepers"], true)
		hosts, err := c.GetList("SELECT host FROM miners_data WHERE miner_id  IN ("+utils.JoinInts(minersIds, ",")+")").String()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		globalRefs[user_id] = globalRefs{amounts: refAmounts, hosts: hosts}
	}

	lastTx, err := c.GetLastTx(c.SessUserId, utils.TypesToIds([]string{"new_user"}), 1, c.TimeFormat)
	lastTxFormatted := ""
	if len(lastTx) > 0 {
		lastTxFormatted, _ = utils.MakeLastTx(lastTx, c.Lang)
	}

	TemplateStr, err := makeTemplate("new_user", "newUser", &newUserPage{
		Alert: c.Alert,
		Lang: c.Lang,
		CountSignArr: c.CountSignArr,
		ShowSignData: c.ShowSignData,
		UserId: c.SessUserId,
		TimeNow: timeNow,
		TxType: txType,
		TxTypeId: txTypeId,
		SignData: "",
		MyRefs: myRefs,
		GlobalRefs: globalRefs,
		CurrencyList: c.CurrencyList,
		PoolUrl: c.NodeConfig["pool_url"]})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}

type myRefs struct {
	amounts map[int64]float64
	hosts []string
	key string
	keyUrl string
}

type globalRefs struct {
	amounts []map[string]string
	hosts []string
}
