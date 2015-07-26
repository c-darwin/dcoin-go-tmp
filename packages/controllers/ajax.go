package controllers
import (
	"net/http"
	"regexp"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"encoding/json"
)

func Ajax(w http.ResponseWriter, r *http.Request) {

	log.Debug("Ajax")
	w.Header().Set("Content-type", "text/html")

	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	sessUserId := GetSessUserId(sess)
	sessRestricted := GetSessRestricted(sess)
	sessPublicKey := GetSessPublicKey(sess)
	log.Debug("sessUserId", sessUserId)
	log.Debug("sessRestricted", sessRestricted)
	log.Debug("sessPublicKey", sessPublicKey)
	log.Debug("user_id", sess.Get("user_id"))

	c := new(Controller)
	c.r = r
	c.w = w
	c.sess = sess
	c.SessRestricted = sessRestricted
	dbInit := false;
	if len(configIni["db_user"]) > 0 || configIni["db_type"]=="sqlite" {
		dbInit = true
	}

	c.SessUserId = sessUserId
	if dbInit {
		var err error
		c.DCDB, err = utils.NewDbConnect(configIni)
		if err != nil {
			log.Error("%v", err)
			dbInit = false
		} else {
			defer c.DCDB.Close()
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
				c.MyPrefix = utils.Int64ToStr(sessUserId)+"_";
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
	parameters_ := make(map[string]interface {})
	err := json.Unmarshal([]byte(c.r.PostFormValue("parameters")), &parameters_)
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

	lang:=GetLang(w, r, parameters)
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
	}

	r.ParseForm()
	controllerName := r.FormValue("controllerName")
	log.Debug("controllerName=",controllerName)

	html := ""

	if ok, _ := regexp.MatchString(`^(?i)AvailableKeys|DcoinKey|SynchronizationBlockchain|PoolAddUsers|SaveQueue|AlertMessage|Menu|SaveHost|GetMinerDataMap|SignUpInPool|Check_sign|PoolDataBaseDump|GetSellerData|GenerateNewPrimaryKey|GenerateNewNodeKey|CheckNode|SignLogin|SaveNotifications|ProgressBar|MinersMap|GetMinerData|EncryptComment|Logout|SaveVideo|SaveShopData|SaveRaceCountry|MyNoticeData|HolidaysList|ClearVideo|CheckCfCurrency|WalletsListCfProject|SendTestEmail|SendSms|SaveUserCoords|SaveGeolocation|SaveEmailSms|Profile|DeleteVideo|CropPhoto$`, controllerName); !ok {
		html = "Access denied"
	} else {
		if ok, _ := regexp.MatchString(`^(?i)CfCatalog|CfPagePreview|CfStart|Check_sign|CheckNode|GetBlock|GetMinerData|GetMinerDataMap|GetSellerData|Index|IndexCf|InstallStep0|InstallStep1|InstallStep2|Login|SignLogin|SynchronizationBlockchain|UpdatingBlockchain|Menu$`, controllerName); !ok && c.SessUserId <= 0 {
			html = "Access denied"
		} else {
			// вызываем контроллер в зависимости от шаблона
			html, err = CallController(c, controllerName)
			if err != nil {
				log.Error("%v", err)
			}
		}
	}
	w.Write([]byte(html))

}