package controllers
import (
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/consts"
	"github.com/astaxie/beego/config"
	"github.com/c-darwin/dcoin-go-tmp/packages/schema"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

type installStep2Struct struct {
	Lang map[string]string
	MyModalIdName string
}

// Шаг 1 - выбор либо стандартных настроек (sqlite и блокчейн с сервера) либо расширенных - pg/mysql и загрузка с нодов
func (c *Controller) InstallStep2() (string, error) {

	c.r.ParseForm()
	installType := c.r.FormValue("type")
	url := c.r.FormValue("url")
	setupPassword := c.r.FormValue("setup_password")
	userId := c.r.FormValue("user_id")
	firstLoad := c.r.FormValue("first_load")
	dbType := c.r.FormValue("db_type")
	dbHost := c.r.FormValue("host")
	dbPort := c.r.FormValue("port")
	dbName := c.r.FormValue("db_name")
	dbUsername := c.r.FormValue("username")
	dbPassword := c.r.FormValue("password")

	if installType=="standard" {
		dbType = "sqlite"
		if len(configIni["blockchain_url"]) > 0 {
			url = configIni["blockchain_url"]
		} else {
			url = consts.BLOCKCHAIN_URL
		}
	}

	confIni, err := config.NewConfig("ini", "config.ini")
	confIni.Set("error_log", "1")
	confIni.Set("log", "0")
	confIni.Set("log_block_id_begin", "0")
	confIni.Set("log_block_id_end", "0")
	confIni.Set("bad_tx_log", "1")
	confIni.Set("nodes_ban_exit", "0")
	confIni.Set("log_tables", "")
	confIni.Set("log_fns", "")
	confIni.Set("sign_hash", "ip")

	fmt.Println(dbType)
	if dbType=="sqlite" {
		confIni.Set("db_user", "")
		confIni.Set("db_host", "")
		confIni.Set("db_port", "")
		confIni.Set("db_password", "")
		confIni.Set("db_name", "")
	} else if dbType=="postgresql" || dbType=="mysql" {
		fmt.Println(dbUsername, dbHost, dbPort)
		confIni.Set("db_type", dbType)
		confIni.Set("db_user", dbUsername)
		confIni.Set("db_host", dbHost)
		confIni.Set("db_port", dbPort)
		confIni.Set("db_password", dbPassword)
		confIni.Set("db_name", dbName)
	}
	err = confIni.SaveConfigFile("config.ini")
	if err != nil {
		return "", err
	}
	fmt.Println(confIni)

	configIni, err = confIni.GetSection("default")
	c.DCDB, err = utils.NewDbConnect(configIni)
	if err != nil {
		return "", err
	} else {
		defer c.DCDB.Close()
	}
	//fmt.Println("c.DCDB", c.DCDB)


	schema_ := &schema.SchemaStruct{}
	schema_.DCDB = c.DCDB
	schema_.DbType = dbType
	schema_.PrefixUserId = 0
	schema_.GetSchema()

	if len(userId)>0 {
		err = c.DCDB.ExecSql("INSERT INTO my_table (user_id) VALUES (?)", userId)
		if err != nil {
			return "", err
		}
	}
	err = c.DCDB.ExecSql("INSERT INTO config (first_load_blockchain, first_load_blockchain_url, setup_password) VALUES (?,?,?)", firstLoad, url, setupPassword)
	if err != nil {
		return "", err
	}
	err = c.DCDB.ExecSql("INSERT INTO payment_systems (name)VALUES ('Adyen'),('Alipay'),('Amazon Payments'),('AsiaPay'),('Atos'),('Authorize.Net'),('BIPS'),('BPAY'),('Braintree'),('CentUp'),('Chargify'),('Citibank'),('ClickandBuy'),('Creditcall'),('CyberSource'),('DataCash'),('DigiCash'),('Digital River'),('Dwolla'),('ecoPayz'),('Edy'),('Elavon'),('Euronet Worldwide'),('eWAY'),('Flooz'),('Fortumo'),('Google'),('GoCardless'),('Heartland Payment Systems'),('HSBC'),('iKobo'),('iZettle'),('IP Payments'),('Klarna'),('Live Gamer'),('Mobilpenge'),('ModusLink'),('MPP Global Solutions'),('Neteller'),('Nochex'),('Ogone'),('Paymate'),('PayPal'),('Payoneer'),('PayPoint'),('Paysafecard'),('PayXpert'),('Payza'),('Peppercoin'),('Playspan'),('Popmoney'),('Realex Payments'),('Recurly'),('RBK Money'),('Sage Group'),('Serve'),('Skrill (Moneybookers)'),('Stripe'),('Square, Inc.'),('TFI Markets'),('TIMWE'),('Use My Services (UMS)'),('Ukash'),('V.me by Visa'),('VeriFone'),('Vindicia'),('WebMoney'),('WePay'),('Wirecard'),('Western Union'),('WorldPay'),('Yandex money'),('Qiwi'),('OK Pay'),('Bitcoin'),('Perfect Money')")
	if err != nil {
		return "", err
	}
	err = c.DCDB.ExecSql(`INSERT INTO cf_lang (id, name) VALUES
	(1, 'English (US)'),
	(2, 'Afrikaans'),
	(3, 'Kiswahili'),
	(4, 'Türkçe'),
	(5, '‏עברית‏'),
	(6, '‏العربية‏'),
	(7, 'Español'),
	(8, 'Français (Canada)'),
	(9, 'Guarani'),
	(10, 'Português (Brasil)'),
	(11, 'Azərbaycan dili'),
	(12, 'Bahasa Indonesia'),
	(13, 'Bahasa Melayu'),
	(14, 'Basa Jawa'),
	(15, 'Bisaya'),
	(16, 'Filipino'),
	(17, 'Tiếng Việt'),
	(18, 'Հայերեն'),
	(19, '‏اردو‏'),
	(20, 'हिन्दी'),
	(21, 'বাংলা'),
	(22, 'ਪੰਜਾਬੀ'),
	(23, 'தமிழ்'),
	(24, 'తెలుగు'),
	(25, 'ಕನ್ನಡ'),
	(26, 'മലയാളം'),
	(27, 'සිංහල'),
	(28, 'ภาษาไทย'),
	(29, '한국어'),
	(30, '中文(台灣)'),
	(31, '中文(简体)'),
	(32, '中文(香港)'),
	(33, '日本語'),
	(35, 'Čeština'),
	(36, 'Magyar'),
	(37, 'Polski'),
	(38, 'Română'),
	(39, 'Slovenčina'),
	(40, 'Slovenščina'),
	(41, 'Български'),
	(42, 'Русский'),
	(43, 'Українська'),
	(45, 'Bosanski'),
	(46, 'Català'),
	(47, 'Cymraeg'),
	(48, 'Dansk'),
	(49, 'Deutsch'),
	(50, 'Eesti'),
	(51, 'English (UK)'),
	(52, 'Español (España)'),
	(53, 'Euskara'),
	(54, 'Français (France)'),
	(55, 'Galego'),
	(56, 'Hrvatski'),
	(57, 'Italiano'),
	(58, 'Latviešu'),
	(59, 'Lietuvių'),
	(60, 'Nederlands'),
	(61, 'Norsk (bokmål)'),
	(62, 'Português (Portugal)'),
	(63, 'Shqip'),
	(64, 'Suomi'),
	(65, 'Svenska'),
	(66, 'Ελληνικά'),
	(67, 'Македонски'),
	(68, 'Српски');`)
	if err != nil {
		return "", err
	}

	err = c.DCDB.ExecSql(`INSERT INTO admin (user_id) VALUES (1)`)
	if err != nil {
		return "", err
	}
	err = c.DCDB.ExecSql(`INSERT INTO install (progress) VALUES ('complete')`)
	if err != nil {
		return "", err
	}


	TemplateStr, err := makeTemplate("install_step_2", "installStep2", &installStep0Struct{
		Lang: c.Lang})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}
