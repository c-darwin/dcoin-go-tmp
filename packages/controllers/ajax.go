package controllers

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"net/http"
	"regexp"
)

func Ajax(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Error("ajax Recovered", r)
			panic(r)
		}
	}()
	log.Debug("Ajax")
	w.Header().Set("Content-type", "text/html")

	sess, err := globalSessions.SessionStart(w, r)
	if err != nil {
		log.Error("%v", err)
		return
	}
	defer sess.SessionRelease(w)
	sessUserId := GetSessUserId(sess)
	sessRestricted := GetSessRestricted(sess)
	sessPublicKey := GetSessPublicKey(sess)
	log.Debug("sessUserId", sessUserId)
	log.Debug("sessRestricted", sessRestricted)
	log.Debug("sessPublicKey", sessPublicKey)
	log.Debug("user_id: %v", sess.Get("user_id"))

	c := new(Controller)
	c.r = r
	c.w = w
	c.sess = sess
	c.SessRestricted = sessRestricted
	dbInit := false
	if len(configIni["db_user"]) > 0 || configIni["db_type"] == "sqlite" {
		dbInit = true
	}

	c.SessUserId = sessUserId
	if dbInit {
		var err error

		//c.DCDB, err = utils.NewDbConnect(configIni)

		c.DCDB = utils.DB

		if utils.DB == nil || utils.DB.DB == nil {
			log.Error("utils.DB == nil")
			dbInit = false
		}
		if dbInit {
			c.Variables, err = c.GetAllVariables()
			var communityUsers []int64
			communityUsers, err = c.GetCommunityUsers()
			if err != nil {
				log.Error("%v", err)
			}
			c.CommunityUsers = communityUsers
			if len(communityUsers) > 0 {
				c.Community = true
			}
			if c.Community {
				poolAdminUserId, err := c.GetPoolAdminUserId()
				if err != nil {
					log.Error("%v", err)
				}
				c.PoolAdminUserId = poolAdminUserId
				if c.SessUserId == poolAdminUserId {
					c.PoolAdmin = true
				}
				c.MyPrefix = utils.Int64ToStr(sessUserId) + "_"
			} else {
				c.PoolAdmin = true
			}
			c.NodeAdmin, err = c.NodeAdminAccess(c.SessUserId, c.SessRestricted)
			if err != nil {
				log.Error("%v", err)
			}
		}
	}
	c.dbInit = dbInit

	c.Parameters, err = c.GetParameters()
	log.Debug("parameters=", c.Parameters)

	lang := GetLang(w, r, c.Parameters)
	log.Debug("lang", lang)
	c.Lang = globalLangReadOnly[lang]
	c.LangInt = int64(lang)
	if lang == 42 {
		c.TimeFormat = "2006-01-02 15:04:05"
	} else {
		c.TimeFormat = "2006-02-01 15:04:05"
	}

	if dbInit {
		myNotice, err := c.GetMyNoticeData(sessRestricted, sessUserId, c.MyPrefix, globalLangReadOnly[lang])
		if err != nil {
			log.Error("%v", err)
		}
		c.MyNotice = myNotice
		config, err := c.GetNodeConfig()
		if err != nil {
			log.Error("%v", err)
		}
		c.NodeConfig = config
	}

	r.ParseForm()
	controllerName := r.FormValue("controllerName")
	log.Debug("controllerName=", controllerName)

	html := ""

	//w.Header().Set("Access-Control-Allow-Origin", "*")

	if ok, _ := regexp.MatchString(`^(?i)rewritePrimaryKey|EImportData|EDataBaseDump|AlertFromAdmin|Update|UpdateDcoin|exchangeAdmin|newPhoto|NodeConfigControl|SaveDecryptComment|EncryptChatMessage|GetChatMessages|SendToTheChat|SaveToken|SendToPool|CheckSetupPassword|getBlock|AcceptNewKeyStatus|availableKeys|ClearDbLite|ClearDb|UploadVideo|DcoinKey|SynchronizationBlockchain|PoolAddUsers|SaveQueue|AlertMessage|Menu|SaveHost|GetMinerDataMap|SignUpInPool|Check_sign|PoolDataBaseDump|GetSellerData|GenerateNewPrimaryKey|GenerateNewNodeKey|CheckNode|SignLogin|SaveNotifications|ProgressBar|MinersMap|GetMinerData|EncryptComment|Logout|SaveVideo|SaveShopData|SaveRaceCountry|MyNoticeData|HolidaysList|ClearVideo|CheckCfCurrency|WalletsListCfProject|SendTestEmail|SendSms|SaveUserCoords|SaveGeolocation|SaveEmailSms|Profile|DeleteVideo|CropPhoto$`, controllerName); !ok {
		html = "Access denied 0"
	} else {
		pages := "CheckSetupPassword|AcceptNewKeyStatus|availableKeys|CfCatalog|CfPagePreview|CfStart|Check_sign|CheckNode|GetBlock|GetMinerData|GetMinerDataMap|GetSellerData|Index|IndexCf|InstallStep0|InstallStep1|InstallStep2|Login|SynchronizationBlockchain|UpdatingBlockchain|Menu|SignUpInPool|SignLogin"
		if utils.IOS() { // На IOS можно сгенерить ключ без сессии
			pages += "|DcoinKey"
		}
		if ok, _ := regexp.MatchString(`^(?i)`+pages+`$`, controllerName); !ok && c.SessUserId <= 0 {
			html = "Access denied 1"
		} else {
			// без БД будет выдавать панику
			if ok, _ := regexp.MatchString(`^(?i)GetChatMessages$`, controllerName); ok && !dbInit {
				html = "Please wait. nill dbInit"
			} else {
				// вызываем контроллер в зависимости от шаблона
				html, err = CallController(c, controllerName)
				if err != nil {
					log.Error("ajax error: %v", err)
				}
			}
		}
	}
	w.Write([]byte(html))

}
