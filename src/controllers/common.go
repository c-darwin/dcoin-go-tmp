package controllers
import (
	//"database/sql"
	//_ "github.com/lib/pq"
	"reflect"
	"utils"
	"net/http"
	"fmt"
	//"bufio"
	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/session"
	"consts"
	"strconv"
	"time"
	"log"
	"regexp"
	"unicode"
	"os"
	"io/ioutil"
	"static"
	"bytes"
	"encoding/json"
	"html/template"
	"strings"
)

type Controller struct {
	dbInit bool
	*utils.DCDB
	r *http.Request
	w http.ResponseWriter
	sess session.SessionStore
	Lang map[string]string
	LangInt int64
	Navigate string
	Periods map[int64]string
	Community bool
	CommunityUsers []int64
	ShowSignData bool
	MyPrefix string
	Alert string
	UserId int64
	Admin bool
	SessRestricted int64
	SessUserId int64
	MyNotice map[string]string
	Parameters map[string]string
	Variables *utils.Variables
	CountSign int
	CountSignArr []int
	TimeFormat string
	PoolAdmin bool
	PoolAdminUserId int64
	NodeAdmin bool
	NodeConfig map[string]string
	ConfigCommission map[int64][]float64
	CurrencyList map[int64]string
	CurrencyListCf map[int64]string
	PaymentSystems map[string]string
	ConfirmedBlockId int64
	MinerId int64
	Races map[int64]string
}

var configIni map[string]string
var globalSessions *session.Manager

// в гоурутинах используется только для чтения
var globalLangReadOnly map[int]map[string]string

func init() {
	globalSessions, _ = session.NewManager("file",`{"cookieName":"gosessionid","gclifetime":864000,"ProviderConfig":"./tmp"}`)
	go globalSessions.GC()

	if _, err := os.Stat("config.ini"); os.IsNotExist(err) {
		fmt.Println("NO")
		d1 := []byte(`
error_log=1
log=1
log_block_id_begin=0
log_block_id_end=0
bad_tx_log=1
nodes_ban_exit=0
log_tables=
log_fns=
sign_hash=ip
db_type=sqlite
DB_USER=
DB_PASSWORD=
DB_NAME=`)
		ioutil.WriteFile("config.ini", d1, 0644)
	} else {
		fmt.Println("YES")
	}

	configIni_, err := config.NewConfig("ini", "config.ini")
	if err != nil {
		log.Fatal(err)
	}
	configIni, err = configIni_.GetSection("default")

	globalLangReadOnly = make(map[int]map[string]string)
	for _, v := range consts.LangMap{
		data, err := static.Asset(fmt.Sprintf("static/lang/%d.ini", v))
		if err != nil {
			fmt.Println(err)
		}
		iniconf_, err := config.NewConfigData("ini", []byte(data))
		if err != nil {
			fmt.Println(err)
		}
		//fmt.Println(iniconf_)
		iniconf, err := iniconf_.GetSection("default")
		globalLangReadOnly[v] = make(map[string]string)
		globalLangReadOnly[v] = iniconf
	}

}

func CallController(c *Controller, name string)  (string, error) {
	// имя экспортируемого метода должно начинаться с заглавной буквы
	a := []rune(name)
	a[0] = unicode.ToUpper(a[0])
	name = string(a)
	fmt.Println("Controller", name)
	html, err := CallMethod(c, name)
	if err != nil {
		log.Print(err)
		html = fmt.Sprintf(`{"error":%q}`, err)
		log.Println(html)
	}
	return html, err
}


func CallMethod(i interface{}, methodName string) (string, error) {
	var ptr reflect.Value
	var value reflect.Value
	var finalMethod reflect.Value

	value = reflect.ValueOf(i)

	// if we start with a pointer, we need to get value pointed to
	// if we start with a value, we need to get a pointer to that value
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(i))
		temp := ptr.Elem()
		temp.Set(value)
	}

	// check for method on value
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	// check for method on pointer
	method = ptr.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}

	if (finalMethod.IsValid()) {
		x:=finalMethod.Call([]reflect.Value{})
		err_, found := x[1].Interface().(error)
		var err error
		if found {
			err = err_
		} else {
			err = nil
		}
		return x[0].Interface().(string), err
	}

	// return or panic, method not found of either type
	return "", fmt.Errorf("not found")
}

func Content1(w http.ResponseWriter, r *http.Request) {

}

func GetSessUserId(sess session.SessionStore) int64 {
	sessUserId := sess.Get("user_id")
	switch sessUserId.(type) {
	case int64:
		return sessUserId.(int64)
	default:
		return 0
	}
	return 0
}

func GetSessAdmin(sess session.SessionStore) int64 {
	admin := sess.Get("admin")
	switch admin.(type) {
	case int64:
		return admin.(int64)
	default:
		return 0
	}
	return 0
}

func DelSessResctricted(sess session.SessionStore) {
	sess.Delete("restricted")
}

func GetSessInt64(sessName string, sess session.SessionStore) int64 {
	sess_ := sess.Get(sessName)
	switch sess_.(type) {
	default:
		return 0
	case int64:
		return sess_.(int64)
	}
	return 0
}

func GetSessRestricted(sess session.SessionStore) int64 {
	sessRestricted := sess.Get("restricted")
	switch sessRestricted.(type) {
	default:
		return 0
	case int64:
		return sessRestricted.(int64)
	}
	return 0
}

func GetSessPublicKey(sess session.SessionStore) string {
	sessPublicKey := sess.Get("public_key")
	switch sessPublicKey.(type) {
	default:
		return ""
	case string:
		return sessPublicKey.(string)
	}
	return ""
}

// ключ в сессии хранится до того момента, пока юзер его не сменит
// т.е. ключ там лежит до тех пор, пока юзер играется и еще не начал пользоваться Dcoin-ом
func GetSessPrivateKey(w http.ResponseWriter, r *http.Request) string {
	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	sessPrivateKey := sess.Get("private_key")
	switch sessPrivateKey.(type) {
	default:
		return ""
	case string:
		return sessPrivateKey.(string)
	}
	return ""
}

func SetLang(w http.ResponseWriter, r *http.Request, lang int) {
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "lang", Value: strconv.Itoa(lang), Expires: expiration}
	http.SetCookie(w, &cookie)
}


// если в lang прислали какую-то гадость
func CheckLang(lang int) bool {
	for _, v := range consts.LangMap{
		if lang == v {
			return true
		}
	}
	return false
}


func GetLang(w http.ResponseWriter, r *http.Request, parameters map[string]string) int {
	var lang int = 1
	lang = utils.StrToInt(parameters["lang"])
	if !CheckLang(lang) {
		if langCookie, err := r.Cookie("lang"); err==nil {
			lang, _ = strconv.Atoi(langCookie.Value)
		}
	}
	if !CheckLang(lang) {
		al := r.Header.Get("Accept-Language")  // en-US,en;q=0.5
		if len(al) >= 2 {
			if _, ok := consts.LangMap[al[:2]]; ok {
				lang = consts.LangMap[al[:2]]
			}
		}
	}
	return lang
}


func Ajax(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Ajax")
	w.Header().Set("Content-type", "text/html")

	sess, _ := globalSessions.SessionStart(w, r)
	defer sess.SessionRelease(w)
	sessUserId := GetSessUserId(sess)
	sessRestricted := GetSessRestricted(sess)
	sessPublicKey := GetSessPublicKey(sess)
	log.Println("sessUserId", sessUserId)
	log.Println("sessRestricted", sessRestricted)
	log.Println("sessPublicKey", sessPublicKey)
	log.Println("user_id", sess.Get("user_id"))

	c := new(Controller)
	c.r = r
	c.w = w
	c.sess = sess
	dbInit := false;
	if len(configIni["db_user"]) > 0 || configIni["db_type"]=="sqlite" {
		dbInit = true
	}

	c.SessUserId = sessUserId
	if dbInit {
		var err error
		c.DCDB, err = utils.NewDbConnect(configIni)
		if err != nil {
			log.Print(err)
			dbInit = false
		} else {
			defer c.DCDB.Close()
		}
		c.Variables, err = c.GetAllVariables()
		var communityUsers []int64
		communityUsers, err = c.GetCommunityUsers()
		if err != nil {
			log.Print(err)
		}
		c.CommunityUsers = communityUsers
		if len(communityUsers) > 0 {
			c.Community = true
		}
		if c.Community {
			poolAdminUserId, err := c.GetPoolAdminUserId()
			if err != nil {
				log.Print(err)
			}
			if c.SessUserId == poolAdminUserId {
				c.PoolAdmin = true
				c.PoolAdminUserId = poolAdminUserId
			}
			c.MyPrefix = utils.Int64ToStr(sessUserId)+"_";
		} else {
			c.PoolAdmin = true
		}
		c.NodeAdmin, err = c.NodeAdminAccess(c.SessUserId, c.SessRestricted)
		if err != nil {
			log.Print(err)
		}
	}
	c.dbInit = dbInit
	parameters_ := make(map[string]interface {})
	err := json.Unmarshal([]byte(c.r.PostFormValue("parameters")), &parameters_)
	if err != nil {
		log.Print(err)
	}
	fmt.Println("parameters_=",parameters_)
	parameters := make(map[string]string)
	for k, v := range parameters_ {
		parameters[k] = utils.InterfaceToStr(v)
	}
	c.Parameters = parameters
	fmt.Println("parameters=",parameters)

	lang:=GetLang(w, r, parameters)
	fmt.Println("lang", lang)
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
			log.Print(err)
		}
		c.MyNotice = myNotice
	}

	r.ParseForm()
	controllerName := r.FormValue("controllerName")
	fmt.Println("controllerName=",controllerName)
	// вызываем контроллер в зависимости от шаблона
	html, err :=  CallController(c, controllerName)
	if err != nil {
		log.Print(err)
	}
	w.Write([]byte(html))
}


func Tools(w http.ResponseWriter, r *http.Request) {

	log.Println("Tools")
	w.Header().Set("Content-type", "text/html")

	c := new(Controller)
	c.r = r
	dbInit := false;
	if len(configIni["db_user"]) > 0 || configIni["db_type"]=="sqlite" {
		dbInit = true
	}

	if dbInit {
		var err error
		c.DCDB, err = utils.NewDbConnect(configIni)
		if err != nil {
			log.Print(err)
			dbInit = false
		} else {
			defer c.DCDB.Close()
		}
		c.Variables, err = c.GetAllVariables()

	}

	r.ParseForm()
	controllerName := r.FormValue("controllerName")
	fmt.Println("controllerName=",controllerName)
	// вызываем контроллер в зависимости от шаблона
	html, err :=  CallController(c, controllerName)
	if err != nil {
		log.Print(err)
	}
	w.Write([]byte(html))
}


func Content(w http.ResponseWriter, r *http.Request) {

	var err error
	log.Println("content")
	w.Header().Set("Content-type", "text/html")

	sess, _ := globalSessions.SessionStart(w, r)

	defer sess.SessionRelease(w)
	sessUserId := GetSessUserId(sess)
	sessRestricted := GetSessRestricted(sess)
	sessPublicKey := GetSessPublicKey(sess)
	sessAdmin := GetSessAdmin(sess)
	log.Println("sessUserId", sessUserId)
	log.Println("sessRestricted", sessRestricted)
	log.Println("sessPublicKey", sessPublicKey)

	c := new(Controller)
	c.r = r
	c.sess = sess
	c.SessRestricted = sessRestricted
	c.SessUserId = sessUserId
	if sessAdmin==1 {
		c.Admin = true
	}
	if err != nil {
		log.Print(err)
	}

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
			log.Print(err)
			dbInit = false
		} else {
			defer utils.DbClose(c.DCDB)
			//defer c.DCDB.Close()
		}
		// отсутвие таблы выдаст ошибку, значит процесс инсталяции еще не пройден и надо выдать 0-й шаг
		_, err = c.DCDB.Single("SELECT progress FROM install").String()
		if err != nil {
			fmt.Println(err)
			dbInit = false
		}
	}

	if dbInit {
		var err error
		installProgress, err = c.DCDB.Single("SELECT progress FROM install").String()
		if err != nil {
			log.Print(err)
		}
		configExists, err = c.DCDB.Single("SELECT first_load_blockchain_url FROM config").String()
		if err != nil {
			log.Print(err)
		}

		c.Variables, err = c.GetAllVariables()

		// Инфа о последнем блоке
		blockData, err := c.DCDB.GetLastBlockData()
		if err != nil {
			log.Print(err)
		}
		//время последнего блока
		lastBlockTime = blockData["lastBlockTime"]
		fmt.Println("installProgress", installProgress, "configExists", configExists,  "lastBlockTime", lastBlockTime)

		// валюты
		currencyListCf, err := c.GetCurrencyList(true)
		if err != nil {
			log.Print(err)
		}
		c.CurrencyListCf = currencyListCf
		currencyList, err := c.GetCurrencyList(false)
		if err != nil {
			log.Print(err)
		}
		c.CurrencyList = currencyList

		confirmedBlockId, err := c.GetConfirmedBlockId()
		if err != nil {
			log.Print(err)
		}
		c.ConfirmedBlockId = confirmedBlockId

		c.MinerId, err = c.GetMyMinerId(c.SessUserId)
		if err != nil {
			log.Print(err)
		}

		paymentSystems, err := c.GetPaymentSystems()
		if err != nil {
			log.Print(err)
		}
		c.PaymentSystems = paymentSystems
	}
	r.ParseForm()
	tplName := r.FormValue("tpl_name")
	parameters_ := make(map[string]interface {})
	err = json.Unmarshal([]byte(c.r.PostFormValue("parameters")), &parameters_)
	if err != nil {
		log.Print(err)
	}
	fmt.Println("parameters_=",parameters_)
	parameters := make(map[string]string)
	for k, v := range parameters_ {
		parameters[k] = utils.InterfaceToStr(v)
	}
	c.Parameters = parameters
	fmt.Println("parameters=",parameters)
	fmt.Println("tpl_name=",tplName)

	// если в параметрах пришел язык, то установим его
	newLang := utils.StrToInt(parameters["lang"])
	if newLang > 0 {
		fmt.Println("newLang", newLang)
		SetLang(w, r, newLang)
	}
	// уведомления
	//if utils.CheckInputData(parameters["alert"], "alert") {
		c.Alert = parameters["alert"]
	//}

	lang:=GetLang(w, r, parameters)
	fmt.Println("lang", lang)

	c.Lang = globalLangReadOnly[lang]
	c.LangInt = int64(lang)
	if lang == 42 {
		c.TimeFormat = "2006-01-02 15:04:05"
	} else {
		c.TimeFormat = "2006-02-01 15:04:05"
	}

	c.Periods = map[int64]string{86400 : "1"+c.Lang["day"], 604800 : "1"+c.Lang["week"], 31536000 : "1"+c.Lang["year"], 2592000 : "1"+c.Lang["month"], 1209600 : "2"+c.Lang["weeks"], }

	c.Races = map[int64]string{1: c.Lang["race_1"], 2: c.Lang["race_2"], 3: c.Lang["race_3"]}

	match, _ := regexp.MatchString("^install_step_[0-9_]+$", tplName)
	// CheckInputData - гарантирует, что tplName чист
	if tplName!="" && utils.CheckInputData(tplName, "tpl_name") && (sessUserId > 0 || match) {
		tplName = tplName
	} else if dbInit && installProgress=="complete" && len(configExists)==0  {
		// первый запуск, еще не загружен блокчейн
		tplName = "after_install"
	} else if dbInit && installProgress=="complete" {
		tplName = "login"
	} else {
		tplName = "install_step_0" // самый первый запуск
	}
	fmt.Println("dbInit", dbInit, "installProgress", installProgress,  "configExists", configExists)
	fmt.Println("tplName>>>>>>>>>>>>>>>>>>>>>>", tplName)

	var communityUsers []int64
	if dbInit {
		communityUsers, err = c.DCDB.GetCommunityUsers()
		if err != nil {
			log.Print(err)
		}
		// нужна мин. комиссия на пуле для перевода монет
		config, err := c.GetNodeConfig()
		if err != nil {
			log.Print(err)
		}
		configCommission_ := make(map[string][]float64)
		if len(config["commission"]) > 0 {
			err = json.Unmarshal([]byte(config["commission"]), &configCommission_)
			if err != nil {
				log.Print(err)
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
			log.Print(err)
		}
	}




	fmt.Println("dbInit", dbInit)
	// идет загрузка блокчейна
	wTime := int64(2)
	if c.ConfigIni["test_mode"] == "1" {
		wTime = 365*24
		log.Println(wTime)
		log.Println(lastBlockTime)
	}
	if dbInit && tplName!="install_step_0" && (time.Now().Unix()-lastBlockTime > 3600*wTime) && len(configExists)>0 {
		if len(communityUsers) > 0 {
			// исключение - админ пула
			poolAdminUserId, err := c.DCDB.Single("SELECT pool_admin_user_id FROM config").String()
			if err != nil {
				log.Print(err)
			}
			if sessUserId != utils.StrToInt64(poolAdminUserId) {
				tplName = "updating_blockchain"
			}
		} else {
			tplName = "updating_blockchain"
		}
	}

	fmt.Println("tplName2=",tplName)

	// кол-во ключей=подписей у юзера
	var countSign int
	var userId int64
	var myUserId int64
	if sessUserId > 0 && dbInit && installProgress == "complete" {
		userId = sessUserId
		myUserId = sessUserId
		countSign = 1
		pk, err := c.DCDB.OneRow("SELECT public_key_1, public_key_2 FROM users WHERE user_id=?", userId).String()
		if err != nil {
			log.Print(err)
		}
		if len(pk["public_key_1"]) > 0 {
			countSign = 2
		}
		if len(pk["public_key_2"]) > 0 {
			countSign = 3
		}
	} else {
		userId = 0
		myUserId = 0
	}
	c.UserId = userId
	var CountSignArr []int
	for i:=0; i < countSign; i++ {
		CountSignArr = append(CountSignArr, i)
	}
	c.CountSign = countSign
	c.CountSignArr = CountSignArr
	fmt.Println(countSign, userId, myUserId ,countSign)

	if len(tplName) > 0 && sessUserId > 0 && installProgress == "complete" {
		// если ключ юзера изменился, то выбрасываем его
		userPublicKey, err := c.DCDB.GetUserPublicKey(userId);
		if err != nil {
			log.Print(err)
		}
		if userPublicKey != sessPublicKey {
			sess.Delete("user_id")
			sess.Delete("private_key")
			sess.Delete("public_key")
			w.Write([]byte("<script language=\"javascript\">window.location.href = \"index.php\"</script>If you are not redirected automatically, follow the <a href=\"index.php\">index.php</a>"))
			return;
		}
		if tplName == "login" {
			tplName = "home"
		}

		if dbInit && len(communityUsers) > 0 {
			c.Community = true
			poolAdminUserId, err := c.GetPoolAdminUserId()
			if err != nil {
				log.Print(err)
			}
			if c.SessUserId == poolAdminUserId {
				c.PoolAdmin = true
				c.PoolAdminUserId = poolAdminUserId
			}
		} else {
			c.PoolAdmin = true
		}

		if dbInit {
			// проверим, не идут ли тех. работы на пуле
			config, err := c.DCDB.OneRow("SELECT pool_admin_user_id, pool_tech_works FROM config").String()
			if err != nil {
				log.Print(err)
			}
			if len(config["pool_admin_user_id"]) > 0 && utils.StrToInt64(config["pool_admin_user_id"]) != sessUserId && config["pool_tech_works"] == "1" {
				tplName = "pool_tech_works"
			}
			if len(communityUsers) == 0 {
				c.MyPrefix = "";
			} else {
				c.MyPrefix = utils.Int64ToStr(sessUserId)+"_";
			}
			// Если у юзера только 1 праймари ключ, то выдавать форму, где показываются данные для подписи и форма ввода подписи не нужно.
			// Только если он сам не захочет, указав это в my_table
			showSignData := false
			if sessRestricted == 0 { // у незареганных в пуле юзеров нет MyPrefix, поэтому сохранять значение show_sign_data им негде
				showSignData_, err := c.DCDB.Single("SELECT show_sign_data FROM "+c.MyPrefix+"my_table").String()
				if err != nil {
					log.Print(err)
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

		if dbInit && tplName !="updating_blockchain" {
			html, err :=  CallController(c, "AlertMessage")
			if err != nil {
				log.Print(err)
			}
			w.Write([]byte(html))
		}
		w.Write([]byte("<input type='hidden' id='tpl_name' value='"+tplName+"'>"))

		myNotice, err := c.DCDB.GetMyNoticeData(sessRestricted, sessUserId, c.MyPrefix, globalLangReadOnly[lang])
		if err != nil {
			log.Print(err)
		}
		c.MyNotice = myNotice

		fmt.Println("tplName==", tplName)

		// подсвечиваем красным номер блока, если идет процесс обновления
		var blockJs string
		blockId, err := c.GetBlockId()
		if err != nil {
			log.Print(err)
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
				log.Print(err)
			}
			w.Write([]byte(html))
		}
	} else if len(tplName) > 0 {
		fmt.Println("tplName",tplName)
		// вызываем контроллер в зависимости от шаблона
		html, err :=  CallController(c, tplName)
		if err != nil {
			log.Print(err)
		}
		w.Write([]byte(html))
	} else {
		html, err :=  CallController(c, "login")
		if err != nil {
			log.Print(err)
		}
		w.Write([]byte(html))
	}
	//sess.Set("username", 11111)

}


func makeTemplate(html, name string, tData interface {}) (string, error) {
	data, err := static.Asset("static/templates/"+html+".html")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	signatures, err := static.Asset("static/templates/signatures.html")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	alert_success, err := static.Asset("static/templates/alert_success.html")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	funcMap := template.FuncMap{
		"div": func(a, b interface{}) float64 {
			return utils.InterfaceToFloat64(a)/utils.InterfaceToFloat64(b)
		},
		"round": func(a  interface{}, num int) float64 {
			return utils.Round(utils.InterfaceToFloat64(a), num)
		},
		"len": func(s []map[string]string) int {
			return len(s)
		},
		"lenMap": func(s map[string]string) int {
			return len(s)
		},
		"sum": func(a, b interface{}) float64 {
			return utils.InterfaceToFloat64(a)+utils.InterfaceToFloat64(b)
		},
		"minus": func(a, b interface{}) float64 {
			return utils.InterfaceToFloat64(a)-utils.InterfaceToFloat64(b)
		},
		"noescape": func(s string) template.HTML {
			return template.HTML(s)
		},
		"js": func(s string) template.JS {
			return template.JS(s)
		},
		"join": func(s []string, sep string) string {
			return strings.Join(s, sep)
		},
		"strToInt64": func(text string) int64 {
			return utils.StrToInt64(text)
		},
		"bin2hex": func(text string) string {
			return string(utils.BinToHex([]byte(text)))
		},
		"int64ToStr": func(text int64) string {
			return utils.Int64ToStr(text)
		},
		"rand": func() int {
			return utils.RandInt(0, 99999999)
		},
		"append": func(args ...interface{}) string {
			var result string
			for _, value := range args {
				switch value.(type) {
				case int64:
					result+=utils.Int64ToStr(value.(int64))
				case float64:
					result+=utils.Float64ToStr(value.(float64))
				case string:
					result+=value.(string)
				}
			}
			return result
		},
		"replaceCurrencyName": func(text, name string) string {
			return strings.Replace(text, "[currency]", "D"+name, -1)
		},
		"cfCategoryLang": func(lang map[string]string, name string) string {
			return lang["cf_category_"+name]
		},
		"checkProjectPs": func(ProjectPs map[string]string, id string) bool {
			if len(ProjectPs["ps"+id]) > 0 {
				return true
			} else {
				return false
			}
		},
		"cfPageTypeLang": func(lang map[string]string, name string) string {
			return lang["cf_"+name]
		},
		"notificationsLang": func(lang map[string]string, name string) string {
			return lang["notifications_"+name]
		},

	}
	t := template.Must(template.New("template").Funcs(funcMap).Parse(string(data)))
	t = template.Must(t.Parse(string(alert_success)))
	t = template.Must(t.Parse(string(signatures)))
	b := new(bytes.Buffer)
	err = t.ExecuteTemplate(b, name, tData)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return b.String(), nil
}
