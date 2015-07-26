package controllers
import (
	"net/http"
	"regexp"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"encoding/json"
)

func Content(w http.ResponseWriter, r *http.Request) {

	var err error

	w.Header().Set("Content-type", "text/html")

	sess, _ := globalSessions.SessionStart(w, r)

	defer sess.SessionRelease(w)
	sessUserId := GetSessUserId(sess)
	sessRestricted := GetSessRestricted(sess)
	sessPublicKey := GetSessPublicKey(sess)
	sessAdmin := GetSessAdmin(sess)
	log.Debug("sessUserId", sessUserId)
	log.Debug("sessRestricted", sessRestricted)
	log.Debug("sessPublicKey", sessPublicKey)

	c := new(Controller)
	c.r = r
	c.sess = sess
	c.SessRestricted = sessRestricted
	c.SessUserId = sessUserId
	if sessAdmin == 1 {
		c.Admin = true
	}
	c.ContentInc = true

	var installProgress, configExists string
	var lastBlockTime int64

	dbInit := false;
	if len(configIni["db_user"]) > 0 || (configIni["db_type"]=="sqlite") {
		dbInit = true
	}

	if dbInit {
		var err error
		c.DCDB, err = utils.NewDbConnect(configIni)
		if err != nil {
			log.Error("%v", err)
			dbInit = false
		} else {
			defer utils.DbClose(c.DCDB)
			//defer c.DCDB.Close()
		}
		if dbInit {
			// отсутвие таблы выдаст ошибку, значит процесс инсталяции еще не пройден и надо выдать 0-й шаг
			_, err = c.DCDB.Single("SELECT progress FROM install").String()
			if err != nil {
				log.Error("%v", err)
				dbInit = false
			}
		}
	}

	c.dbInit = dbInit;

	if dbInit {
		var err error
		installProgress, err = c.DCDB.Single("SELECT progress FROM install").String()
		if err != nil {
			log.Error("%v", err)
		}
		configExists, err = c.DCDB.Single("SELECT first_load_blockchain_url FROM config").String()
		if err != nil {
			log.Error("%v", err)
		}

		c.Variables, err = c.GetAllVariables()

		// Инфа о последнем блоке
		blockData, err := c.DCDB.GetLastBlockData()
		if err != nil {
			log.Error("%v", err)
		}
		//время последнего блока
		lastBlockTime = blockData["lastBlockTime"]
		log.Debug("installProgress", installProgress, "configExists", configExists,  "lastBlockTime", lastBlockTime)

		// валюты
		currencyListCf, err := c.GetCurrencyList(true)
		if err != nil {
			log.Error("%v", err)
		}
		c.CurrencyListCf = currencyListCf
		currencyList, err := c.GetCurrencyList(false)
		if err != nil {
			log.Error("%v", err)
		}
		c.CurrencyList = currencyList

		confirmedBlockId, err := c.GetConfirmedBlockId()
		if err != nil {
			log.Error("%v", err)
		}
		c.ConfirmedBlockId = confirmedBlockId

		c.MinerId, err = c.GetMyMinerId(c.SessUserId)
		if err != nil {
			log.Error("%v", err)
		}

		paymentSystems, err := c.GetPaymentSystems()
		if err != nil {
			log.Error("%v", err)
		}
		c.PaymentSystems = paymentSystems
	}
	r.ParseForm()
	tplName := r.FormValue("tpl_name")
	parameters_ := make(map[string]interface {})
	err = json.Unmarshal([]byte(c.r.PostFormValue("parameters")), &parameters_)
	if err != nil {
		log.Error("%v", err)
	}
	log.Debug("parameters_=",parameters_)
	parameters := make(map[string]string)
	for k, v := range parameters_ {
		parameters[k] = utils.InterfaceToStr(v)
	}
	c.Parameters = parameters
	log.Debug("parameters=",parameters)
	log.Debug("tpl_name=",tplName)

	// если в параметрах пришел язык, то установим его
	newLang := utils.StrToInt(parameters["lang"])
	if newLang > 0 {
		log.Debug("newLang", newLang)
		SetLang(w, r, newLang)
	}
	// уведомления
	//if utils.CheckInputData(parameters["alert"], "alert") {
	c.Alert = parameters["alert"]
	//}

	lang:=GetLang(w, r, parameters)
	log.Debug("lang", lang)

	c.Lang = globalLangReadOnly[lang]
	c.LangInt = int64(lang)
	if lang == 42 {
		c.TimeFormat = "2006-01-02 15:04:05"
	} else {
		c.TimeFormat = "2006-02-01 15:04:05"
	}

	c.Periods = map[int64]string{86400 : "1"+c.Lang["day"], 604800 : "1"+c.Lang["week"], 31536000 : "1"+c.Lang["year"], 2592000 : "1"+c.Lang["month"], 1209600 : "2"+c.Lang["weeks"], }

	c.Races = map[int64]string{1: c.Lang["race_1"], 2: c.Lang["race_2"], 3: c.Lang["race_3"]}

	match, _ := regexp.MatchString("^installStep[0-9_]+$", tplName)
	// CheckInputData - гарантирует, что tplName чист
	if tplName!="" && utils.CheckInputData(tplName, "tpl_name") && (sessUserId > 0 || match) {
		tplName = tplName
	} else if dbInit && installProgress=="complete" && len(configExists)==0  {
		// первый запуск, еще не загружен блокчейн
		tplName = "after_install"
	} else if dbInit && installProgress=="complete" {
		tplName = "login"
	} else {
		tplName = "installStep0" // самый первый запуск
	}
	log.Debug("dbInit", dbInit, "installProgress", installProgress,  "configExists", configExists)
	log.Debug("tplName>>>>>>>>>>>>>>>>>>>>>>", tplName)

	var communityUsers []int64
	if dbInit {
		communityUsers, err = c.DCDB.GetCommunityUsers()
		if err != nil {
			log.Error("%v", err)
		}
		c.CommunityUsers = communityUsers
		if len(communityUsers) == 0 {
			c.MyPrefix = "";
		} else {
			c.MyPrefix = utils.Int64ToStr(sessUserId)+"_";
			c.Community = true
		}
		// нужна мин. комиссия на пуле для перевода монет
		config, err := c.GetNodeConfig()
		if err != nil {
			log.Error("%v", err)
		}
		configCommission_ := make(map[string][]float64)
		if len(config["commission"]) > 0 {
			err = json.Unmarshal([]byte(config["commission"]), &configCommission_)
			if err != nil {
				log.Error("%v", err)
			}
		}
		configCommission := make(map[int64][]float64)
		for k, v := range configCommission_{
			configCommission[utils.StrToInt64(k)] = v
		}
		c.NodeConfig = config
		c.ConfigCommission = configCommission

		c.NodeAdmin, err = c.NodeAdminAccess(c.SessUserId, c.SessRestricted)
		if err != nil {
			log.Error("%v", err)
		}
	}




	log.Debug("dbInit", dbInit)
	// идет загрузка блокчейна
	wTime := int64(2)
	if c.ConfigIni["test_mode"] == "1" {
		wTime = 2*365*86400
		log.Debug("%v", wTime)
		log.Debug("%v", lastBlockTime)
	}
	if dbInit && tplName!="installStep0" && (utils.Time()-lastBlockTime > 3600*wTime) && len(configExists)>0 {
		if len(communityUsers) > 0 {
			// исключение - админ пула
			poolAdminUserId, err := c.DCDB.Single("SELECT pool_admin_user_id FROM config").String()
			if err != nil {
				log.Error("%v", err)
			}
			if sessUserId != utils.StrToInt64(poolAdminUserId) {
				tplName = "updatingBlockchain"
			}
		} else {
			tplName = "updatingBlockchain"
		}
	}

	log.Debug("tplName2=",tplName)

	// кол-во ключей=подписей у юзера
	var countSign int
	var userId int64
	//	var myUserId int64
	if sessUserId > 0 && dbInit && installProgress == "complete" {
		userId = sessUserId
		//	myUserId = sessUserId
		countSign = 1
		pk, err := c.DCDB.OneRow("SELECT public_key_1, public_key_2 FROM users WHERE user_id=?", userId).String()
		if err != nil {
			log.Error("%v", err)
		}
		if len(pk["public_key_1"]) > 0 {
			countSign = 2
		}
		if len(pk["public_key_2"]) > 0 {
			countSign = 3
		}
	} else {
		userId = 0
		//myUserId = 0
	}
	c.UserId = userId
	var CountSignArr []int
	for i:=0; i < countSign; i++ {
		CountSignArr = append(CountSignArr, i)
	}
	c.CountSign = countSign
	c.CountSignArr = CountSignArr

	log.Debug("tplName::", tplName, sessUserId, installProgress)

	if ok, _ := regexp.MatchString(`^(?i)CfPagePreview|CfCatalog|AddCfProjectData|CfProjectChangeCategory|NewCfProject|MyCfProjects|DelCfProject|DelCfFunding|CfStart|PoolAdminControl|Credits|Home|WalletsList|Information|Notifications|Interface|MiningMenu|Upgrade5|NodeConfigControl|Upgrade7|Upgrade6|Upgrade5|Upgrade4|Upgrade3|Upgrade2|Upgrade1|Upgrade0|StatisticVoting|ProgressBar|MiningPromisedAmount|CurrencyExchangeDelete|CurrencyExchange|ChangeCreditor|ChangeCommission|CashRequestOut|ArbitrationSeller|ArbitrationBuyer|ArbitrationArbitrator|Arbitration|InstallStep2|InstallStep1|InstallStep0|DbInfo|ChangeHost|Assignments|NewUser|NewPhoto|Voting|VoteForMe|RepaymentCredit|PromisedAmountList|PromisedAmountActualization|NewPromisedAmount|Login|ForRepaidFix|DelPromisedAmount|DelCredit|ChangePromisedAmount|ChangePrimaryKey|ChangeNodeKey|ChangeAvatar|BugReporting|Abuse|UpgradeResend|UpdatingBlockchain|Statistic|RewritePrimaryKey|RestoringAccess|PoolTechWorks|Points|NewHolidays|NewCredit|MoneyBackRequest|MoneyBack|ChangeMoneyBack|ChangeKeyRequest|ChangeKeyClose|ChangeGeolocation|ChangeCountryRace|ChangeArbitratorConditions|CashRequestIn|BlockExplorer$`, tplName); !ok {
		w.Write([]byte("Access denied"))
	} else if len(tplName) > 0 && sessUserId > 0 && installProgress == "complete" {
		// если ключ юзера изменился, то выбрасываем его
		userPublicKey, err := c.DCDB.GetUserPublicKey(userId);
		if err != nil {
			log.Error("%v", err)
		}
		if userPublicKey != sessPublicKey {
			sess.Delete("user_id")
			sess.Delete("private_key")
			sess.Delete("public_key")
			log.Debug("sess.Delete user_id private_key public_key")
			w.Write([]byte("<script language=\"javascript\">window.location.href = \"/\"</script>If you are not redirected automatically, follow the <a href=\"/\">/</a>"))
			return;
		}
		if tplName == "login" {
			tplName = "home"
		}

		c.TplName = tplName

		log.Debug("communityUsers:", communityUsers)
		if dbInit && len(communityUsers) > 0 {
			poolAdminUserId, err := c.GetPoolAdminUserId()
			if err != nil {
				log.Error("%v", err)
			}
			c.PoolAdminUserId = poolAdminUserId
			if c.SessUserId == poolAdminUserId {
				c.PoolAdmin = true
			}
		} else {
			c.PoolAdmin = true
		}

		if dbInit {
			// проверим, не идут ли тех. работы на пуле
			config, err := c.DCDB.OneRow("SELECT pool_admin_user_id, pool_tech_works FROM config").String()
			if err != nil {
				log.Error("%v", err)
			}
			if len(config["pool_admin_user_id"]) > 0 && utils.StrToInt64(config["pool_admin_user_id"]) != sessUserId && config["pool_tech_works"] == "1" {
				tplName = "pool_tech_works"
			}
			// Если у юзера только 1 праймари ключ, то выдавать форму, где показываются данные для подписи и форма ввода подписи не нужно.
			// Только если он сам не захочет, указав это в my_table
			showSignData := false
			if sessRestricted == 0 { // у незареганных в пуле юзеров нет MyPrefix, поэтому сохранять значение show_sign_data им негде
				showSignData_, err := c.DCDB.Single("SELECT show_sign_data FROM "+c.MyPrefix+"my_table").String()
				if err != nil {
					log.Error("%v", err)
				}
				if showSignData_ == "1" {
					showSignData = true
				} else {
					showSignData = false
				}
			}
			if showSignData || countSign > 1 {
				c.ShowSignData = true
			} else {
				c.ShowSignData = false
			}
		}

		if dbInit && tplName !="updatingBlockchain" {
			html, err :=  CallController(c, "AlertMessage")
			if err != nil {
				log.Error("%v", err)
			}
			w.Write([]byte(html))
		}
		w.Write([]byte("<input type='hidden' id='tpl_name' value='"+tplName+"'>"))

		myNotice, err := c.DCDB.GetMyNoticeData(sessRestricted, sessUserId, c.MyPrefix, globalLangReadOnly[lang])
		if err != nil {
			log.Error("%v", err)
		}
		c.MyNotice = myNotice

		log.Debug("tplName==", tplName)

		// подсвечиваем красным номер блока, если идет процесс обновления
		var blockJs string
		blockId, err := c.GetBlockId()
		if err != nil {
			log.Error("%v", err)
		}
		if myNotice["main_status_complete"] != "1" {
			blockJs = "$('#block_id').html("+utils.Int64ToStr(blockId)+");$('#block_id').css('color', '#ff0000');";
		} else {
			blockJs = "$('#block_id').html("+utils.Int64ToStr(blockId)+");$('#block_id').css('color', '#428BCA');";
		}
		w.Write([]byte(`<script>
								$( document ).ready(function() {
								$('.lng_1').attr('href', '#`+tplName+`/lang=1');
								$('.lng_42').attr('href', '#`+tplName+`/lang=42');
								`+blockJs+`
								});
								</script>`))
		skipRestrictedUsers := []string{"cash_requests_in", "cash_requests_out", "upgrade", "notifications"}
		// тем, кто не зареган на пуле не выдаем некоторые страницы
		if ( sessRestricted == 0 || !utils.InSliceString(tplName, skipRestrictedUsers) ) {
			// вызываем контроллер в зависимости от шаблона
			html, err :=  CallController(c, tplName)
			if err != nil {
				log.Error("%v", err)
			}
			w.Write([]byte(html))
		}
	} else if len(tplName) > 0 {
		log.Debug("tplName",tplName)
		html := ""
		if ok, _ := regexp.MatchString(`^(?i)CfCatalog|CfPagePreview|CfStart|Check_sign|CheckNode|GetBlock|GetMinerData|GetMinerDataMap|GetSellerData|Index|IndexCf|InstallStep0|InstallStep1|InstallStep2|Login|SignLogin|SynchronizationBlockchain|UpdatingBlockchain|Menu$`, tplName); !ok && c.SessUserId <= 0 {
			html = "Access denied"
		} else {
			// вызываем контроллер в зависимости от шаблона
			html, err = CallController(c, tplName)
			if err != nil {
				log.Error("%v", err)
			}
		}
		w.Write([]byte(html))
	} else {
		html, err :=  CallController(c, "login")
		if err != nil {
			log.Error("%v", err)
		}
		w.Write([]byte(html))
	}
	//sess.Set("username", 11111)

}