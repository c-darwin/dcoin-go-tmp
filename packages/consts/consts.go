package consts

import (
//"fmt"
)

// чтобы не выдавать одно и тоже голосование
const ASSIGN_TIME = 86400

const DAY = 3600 * 24
const DAY2 = 3600 * 24 * 2

// используется в confirmations
const COUNT_CONFIRMED_NODES = 5
const WAIT_CONFIRMED_NODES = 10

// на сколько % автоматически урезаем денежную массу
const AUTO_REDUCTION_PCT = 10

// У скольких нодов должен быть такой же блок как и у нас, чтобы считать, что блок у большей части DC-сети. для get_confirmed_block_id()
const MIN_CONFIRMED_NODES = 1

// текущая версия
const VERSION = "1.0.1a3"

// примерный текущий крайний блок
const LAST_BLOCK = 270000

// примерный размер блокчейна
const BLOCKCHAIN_SIZE = 65000000

// где лежит блокчейн. для тех, кто не хочет собирать его с нодов
const BLOCKCHAIN_URL = "http://dcoin.club/blockchain"

// на сколько может бежать время в тр-ии
const MAX_TX_FORW = 0

// тр-ия может блуждать по сети сутки и потом попасть в блок
const MAX_TX_BACK = DAY

const USD_CURRENCY_ID = 71

const ARBITRATION_BLOCK_START = 189300

// через какое время админ имеет право изменить ключ юзера, если тот дал на это свое согласие. Это время дается юзеру на то, чтобы отменить запрос.
const CHANGE_KEY_PERIOD_170770 = 3600
const CHANGE_KEY_PERIOD = 3600 * 24 * 30

//  есть ли хотябы X юзеров, у которых на кошелках есть от 0.01 данной валюты
const AUTO_REDUCTION_PROMISED_AMOUNT_MIN = 10

// сколько должно быть процентов PROMISED_AMOUNT от кол-ва DC на кошельках, чтобы запустилось урезание
const AUTO_REDUCTION_PROMISED_AMOUNT_PCT = 1 // X*100%

const LIMIT_NEW_CF_PROJECT = 1
const LIMIT_NEW_CF_PROJECT_PERIOD = 3600 * 24 * 7
const LIMIT_CF_PROJECT_DATA = 10
const LIMIT_CF_PROJECT_DATA_PERIOD = 3600 * 24
const LIMIT_CF_SEND_DC = 10
const LIMIT_CF_SEND_DC_PERIOD = 3600 * 24
const LIMIT_CF_COMMENTS = 10
const LIMIT_CF_COMMENTS_PERIOD = 3600 * 24

// сколько можно делать комментов за сутки за 1 проект
const LIMIT_TIME_COMMENTS_CF_PROJECT = 3600 * 24

const LIMIT_USER_AVATAR = 5
const LIMIT_USER_AVATAR_PERIOD = 3600 * 24

const LIMIT_NEW_CREDIT = 10
const NEW_CREDIT_PERIOD = 3600 * 24
const LIMIT_CHANGE_CREDITOR = 10
const CHANGE_CREDITOR_PERIOD = 3600 * 24
const LIMIT_REPAYMENT_CREDIT = 5
const REPAYMENT_CREDIT_PERIOD = 3600 * 24
const LIMIT_CHANGE_CREDIT_PART = 10
const LIMIT_CHANGE_CREDIT_PART_PERIOD = 3600 * 24
const LIMIT_CHANGE_KEY_ACTIVE = 3
const LIMIT_CHANGE_KEY_ACTIVE_PERIOD = 3600 * 24 * 7
const LIMIT_CHANGE_KEY_REQUEST = 1
const LIMIT_CHANGE_KEY_REQUEST_PERIOD = 3600 * 24 * 7
const LIMIT_CHANGE_ARBITRATION_TRUST_LIST = 3
const LIMIT_CHANGE_ARBITRATION_TRUST_LIST_PERIOD = 3600 * 24
const LIMIT_CHANGE_ARBITRATOR_CONDITIONS = 3
const LIMIT_CHANGE_ARBITRATOR_CONDITIONS_PERIOD = 3600 * 24
const LIMIT_MONEY_BACK_REQUEST = 3
const LIMIT_MONEY_BACK_REQUEST_PERIOD = 3600 * 24
const LIMIT_CHANGE_SELLER_HOLD_BACK = 3
const LIMIT_CHANGE_SELLER_HOLD_BACK_PERIOD = 3600 * 24
const LIMIT_CHANGE_CA = 3
const LIMIT_CHANGE_CA_PERIOD = 3600 * 24

const CRON_CHECKED_TIME_SEC = 86400 * 3

const ROUND_FIX = 0.00000000001

// таймауты для TCP
const READ_TIMEOUT = 20
const WRITE_TIMEOUT = 20

// дефолтное знаение, со скольким нодами устанавляиваем связь
const OUT_CONNECTIONS = 10

// на какое время баним нода, давшего нам плохие данные
const NODE_BAN_TIME = 3600

// через сколько можно делать следующее урезание.
// важно учитывать то, что не должно быть роллбеков дальше чем на 1 урезание
// т.к. при урезании используется backup в этой же табле вместо отдельной таблы log_
const AUTO_REDUCTION_PERIOD = 3600 * 24 * 2

const LIMIT_ACTUALIZATION = 1
const LIMIT_ACTUALIZATION_PERIOD = 3600 * 24 * 14

// на сколько арбитр может продлить время рассмотрения манибека
const MAX_MONEY_BACK_TIME = 180

const CHAT_PORT = "8150"

const COUNT_CHAT_NODES = 10

const ALERT_KEY = `30820122300d06092a864886f70d01010105000382010f003082010a0282010100d4a48242d0fb2c7c295bc9c87b1aa0c6d23b5f8cab2ec20c2dfde35513ef6066b92ee3935f9a38100493717b60bb7832411daee02012f44a9f58ac91056b2603661544116bfbc55181e5a693bace5ec9325ba0232b9c9c0a29096569d217243e5bf891cc7fc4bcd2e7d6518acc6f982aaa43a9ed737e3ea2845d6432a823ee5b40d1548f802d0c108bf6e5cb5a4daa7edb48764dcbfa6b7a961208833996cfee265ca2ce2655d444cf3c177b3841b1cc4f3102f89cb2bdb1e5a68eac270506147dd8391b7b3af40a50be13c3970077faffaf98ccc5b8c011146be9c2eb9dfd3454f67a68daaf385d334366d132308bffede27656a515ff69a260bbe2452bd2c30203010001`

var LangMap = map[string]int{"en": 1, "ru": 42}

var MyTables = []string{"my_admin_messages", "my_cash_requests", "my_comments", "my_commission", "my_complex_votes", "my_dc_transactions", "my_holidays", "my_keys", "my_new_users", "my_node_keys", "my_notifications", "my_promised_amount", "my_table", "my_tasks", "my_cf_funding"}

var ReductionDC = []int64{0, 10, 25, 50, 90}
var Countries = []string{"Afghanistan", "Albania", "Algeria", "American Samoa", "Andorra", "Angola", "Anguilla", "Antarctica", "Antigua and Barbuda", "Argentina", "Armenia", "Aruba", "Australia", "Austria", "Azerbaijan", "Bahamas", "Bahrain", "Bangladesh", "Barbados", "Belarus", "Belgium", "Belize", "Benin", "Bermuda", "Bhutan", "Bolivia", "Bosnia and Herzegovina", "Botswana", "Bouvet Island", "Brazil", "British Indian Ocean Territory", "British Virgin Islands", "Brunei", "Bulgaria", "Burkina Faso", "Burundi", "Cambodia", "Cameroon", "Canada", "Cape Verde", "Cayman Islands", "Central African Republic", "Chad", "Chile", "China", "Christmas Island", "Cocos [Keeling] Islands", "Colombia", "Comoros", "Congo [DRC]", "Congo [Republic]", "Cook Islands", "Costa Rica", "Croatia", "Cuba", "Cyprus", "Czech Republic", "Côte d\"Ivoire", "Denmark", "Djibouti", "Dominica", "Dominican Republic", "Ecuador", "Egypt", "El Salvador", "Equatorial Guinea", "Eritrea", "Estonia", "Ethiopia", "Falkland Islands [Islas Malvinas]", "Faroe Islands", "Fiji", "Finland", "France", "French Guiana", "French Polynesia", "French Southern Territories", "Gabon", "Gambia", "Gaza Strip", "Georgia", "Germany", "Ghana", "Gibraltar", "Greece", "Greenland", "Grenada", "Guadeloupe", "Guam", "Guatemala", "Guernsey", "Guinea", "Guinea-Bissau", "Guyana", "Haiti", "Heard Island and McDonald Islands", "Honduras", "Hong Kong", "Hungary", "Iceland", "India", "Indonesia", "Iran", "Iraq", "Ireland", "Isle of Man", "Israel", "Italy", "Jamaica", "Japan", "Jersey", "Jordan", "Kazakhstan", "Kenya", "Kiribati", "Kosovo", "Kuwait", "Kyrgyzstan", "Laos", "Latvia", "Lebanon", "Lesotho", "Liberia", "Libya", "Liechtenstein", "Lithuania", "Luxembourg", "Macau", "Macedonia [FYROM]", "Madagascar", "Malawi", "Malaysia", "Maldives", "Mali", "Malta", "Marshall Islands", "Martinique", "Mauritania", "Mauritius", "Mayotte", "Mexico", "Micronesia", "Moldova", "Monaco", "Mongolia", "Montenegro", "Montserrat", "Morocco", "Mozambique", "Myanmar [Burma]", "Namibia", "Nauru", "Nepal", "Netherlands", "Netherlands Antilles", "New Caledonia", "New Zealand", "Nicaragua", "Niger", "Nigeria", "Niue", "Norfolk Island", "North Korea", "Northern Mariana Islands", "Norway", "Oman", "Pakistan", "Palau", "Palestinian Territories", "Panama", "Papua New Guinea", "Paraguay", "Peru", "Philippines", "Pitcairn Islands", "Poland", "Portugal", "Puerto Rico", "Qatar", "Romania", "Russia", "Rwanda", "Réunion", "Saint Helena", "Saint Kitts and Nevis", "Saint Lucia", "Saint Pierre and Miquelon", "Saint Vincent and the Grenadines", "Samoa", "San Marino", "Saudi Arabia", "Senegal", "Serbia", "Seychelles", "Sierra Leone", "Singapore", "Slovakia", "Slovenia", "Solomon Islands", "Somalia", "South Africa", "South Georgia and the South Sandwich Islands", "South Korea", "Spain", "Sri Lanka", "Sudan", "Suriname", "Svalbard and Jan Mayen", "Swaziland", "Sweden", "Switzerland", "Syria", "São Tomé and Príncipe", "Taiwan", "Tajikistan", "Tanzania", "Thailand", "Timor-Leste", "Togo", "Tokelau", "Tonga", "Trinidad and Tobago", "Tunisia", "Turkey", "Turkmenistan", "Turks and Caicos Islands", "Tuvalu", "U.S. Minor Outlying Islands", "U.S. Virgin Islands", "Uganda", "Ukraine", "United Arab Emirates", "United Kingdom", "United States", "Uruguay", "Uzbekistan", "Vanuatu", "Vatican City", "Venezuela", "Vietnam", "Wallis and Futuna", "Western Sahara", "Yemen", "Zambia", "Zimbabwe"}

var TxTypes = map[int]string{
	// новый юзер
	1: "NewUser",
	// новый майнер
	2: "NewMiner",
	// Добавление новой обещанной суммы
	3: "NewPromisedAmount",
	4: "ChangePromisedAmount",
	// голос за претендента на майнера
	5: "VotesMiner",
	6: "NewForexOrder",
	7: "DelForexOrder",
	//  новый набор max_other_currencies от нода-генератора блока
	8: "NewMaxOtherCurrencies",
	// geolocation. Майнер изменил свои координаты
	9: "ChangeGeolocation",
	// votes_promised_amount.
	10: "VotesPromisedAmount",
	// del_promised_amount. Удаление обещанной суммы
	11: "DelPromisedAmount",
	// send_dc
	12: "SendDc",
	13: "CashRequestOut",
	14: "CashRequestIn",
	// набор голосов по разным валютам
	15: "VotesComplex",
	16: "ChangePrimaryKey",
	17: "ChangeNodeKey",
	18: "ForRepaidFix",
	// занесение в БД данных из первого блока
	19: "Admin1Block", //
	// админ разжаловал майнеров в юзеры
	20: "AdminBanMiners",
	// админ изменил variables
	21: "AdminVariables",
	// админ обновил набор точек для проверки лиц
	22: "AdminSpots",
	// юзер создал кредит
	23: "NewCredit",
	// админ вернул майнерам звание "майнер"
	24: "AdminUnbanMiners",
	// админ отправил alert message
	25: "AdminMessage",
	// майнер хочет, чтобы указаные им майнеры были разжалованы в юзеры
	26: "Abuses",
	// майнер хочет, чтобы в указанные дни ему не приходили запросы на обмен DC
	27: "NewHolidays",
	28: "ActualizationPromisedAmounts",
	29: "Mining",
	// Голосование нода за фото нового майнера
	30: "VotesNodeNewMiner",
	// Юзер исправил проблему с отдачей фото и шлет повторный запрос на получение статуса "майнер"
	31: "NewMinerUpdate",
	//  новый набор max_promised_amount от нода-генератора блока
	32: "NewMaxPromisedAmounts",
	//  новый набор % от нода-генератора блока
	33: "NewPct",
	// добавление новой валюты
	34: "AdminAddCurrency",
	35: "NewCfProject",
	// новая версия, которая кладется каждому в диру public
	36: "AdminNewVersion",
	// после того, как новая версия протестируется, выдаем сообщение, что необходимо обновиться
	37: "AdminNewVersionAlert",
	// баг репорты
	38: "MessageToAdmin",
	// админ может ответить юзеру
	39: "AdminAnswer",
	40: "CfProjectData",
	// блог админа
	41: "AdminBlog",
	// майнер меняет свой хост
	42: "ChangeHost",
	// майнер меняет комиссию, которую он хочет получать с тр-ий
	43: "ChangeCommission",
	44: "DelCfFunding",
	// запуск урезания на основе голосования. генерит нод-генератор блока
	45: "NewReduction",
	46: "DelCfProject",
	47: "CfComment",
	48: "CfSendDc",
	49: "UserAvatar",
	50: "CfProjectChangeCategory",
	51: "ChangeCreditor",
	52: "DelCredit",
	53: "RepaymentCredit",
	54: "ChangeCreditPart",
	55: "NewAdmin",
	// по истечении 30 дней после поступления запроса о восстановлении утерянного ключа, админ может изменить ключ юзера
	56: "AdminChangePrimaryKey",
	// юзер разрешает или отменяет разрешение на смену своего ключа админом
	57: "ChangeKeyActive",
	// юзер отменяет запрос на смену ключа
	58: "ChangeKeyClose",
	// юзер отправляет с другого акка запрос на получение доступа к акку, ключ к которому он потерял
	59: "ChangeKeyRequest",
	// юзер решил стать арбитром или же действующий арбитр меняет комиссии
	60: "ChangeArbitratorConditions",
	// продавец меняет % и кол-во дней для новых сделок.
	61: "ChangeSellerHoldBack",
	// покупатель или продавец указал список арбитров, кому доверяет
	62: "ChangeArbitratorList",
	// покупатель хочет манибэк
	63: "MoneyBackRequest",
	// магазин добровольно делает манибэк или же арбитр делает манибек
	64: "MoneyBack",
	// арбитр увеличивает время манибэка, чтобы успеть разобраться в ситуации
	65: "ChangeMoneyBackTime",
	// юзер меняет url центров сертификации, где хранятся его приватные ключи
	66: "ChangeCa",
	// бан юзера в чате
	67: "AdminBanUnbanChat",
	// голосование за биржи
	68: "VotesExchange",
}

func init() {
}
