package consts
import (
	//"fmt"
)

// У скольких нодов должен быть такой же блок как и у нас, чтобы считать, что блок у большей части DC-сети. для get_confirmed_block_id()
const MIN_CONFIRMED_NODES = 3

// текущая версия
const VERSION = "1.0.1b1"

// примерный текущий крайний блок
const LAST_BLOCK = 250000

// примерный размер блокчейна
const BLOCKCHAIN_SIZE = 65000000

// где лежит блокчейн. для тех, кто не хочет собирать его с нодов
const BLOCKCHAIN_URL = "http://localhost/blockchain"

// на сколько может бежать время в тр-ии
const MAX_TX_FORW = 0

// тр-ия может блуждать по сети сутки и потом попасть в блок
const MAX_TX_BACK = 3600*24


var LangMap = map[string]int{"en":1, "ru":42}

var MyTables = []string {"my_admin_messages","my_cash_requests","my_comments","my_commission","my_complex_votes","my_dc_transactions","my_holidays","my_keys","my_new_users","my_node_keys","my_notifications","my_promised_amount","my_table","my_tasks","my_cf_funding"}

var ReductionDC = []int {0,10,25,50,90}

var TxTypes = map[int]string {
	// новый юзер
	1 : "NewUser",
	// новый майнер
	2 : "NewMiner",
	// Добавление новой обещанной суммы
	3 : "NewPromisedAmount",
	4 : "ChangePromisedAmount",
	// голос за претендента на майнера
	5 : "VotesMiner",
	6 : "new_forex_order",
	7 : "del_forex_order",
	//  новый набор max_other_currencies от нода-генератора блока
	8 : "new_max_other_currencies",
	// geolocation. Майнер изменил свои координаты
	9 : "change_geolocation",
	// votes_promised_amount.
	10 : "votes_promised_amount",
	// del_promised_amount. Удаление обещанной суммы
	11 : "del_promised_amount",
	// send_dc
	12 : "send_dc",
	13 : "cash_request_out",
	14 : "cash_request_in",
	// набор голосов по разным валютам
	15 : "votes_complex",
	16 : "change_primary_key",
	17 : "change_node_key",
	18 : "for_repaid_fix",
	// занесение в БД данных из первого блока
	19 : "Admin1Block",
	// админ разжаловал майнеров в юзеры
	20 : "admin_ban_miners",
	// админ изменил variables
	21 : "AdminVariables",
	// админ обновил набор точек для проверки лиц
	22 : "AdminSpots",
	// юзер создал кредит
	23 : "new_credit",
	// админ вернул майнерам звание "майнер"
	24 : "admin_unban_miners",
	// админ отправил alert message
	25 : "admin_message",
	// майнер хочет, чтобы указаные им майнеры были разжалованы в юзеры
	26 : "abuses",
	// майнер хочет, чтобы в указанные дни ему не приходили запросы на обмен DC
	27 : "new_holidays",
	28 : "actualization_promised_amounts",
	29 : "mining",
	// Голосование нода за фото нового майнера
	30 : "VotesNodeNewMiner",
	// Юзер исправил проблему с отдачей фото и шлет повторный запрос на получение статуса "майнер"
	31 : "new_miner_update",
	//  новый набор max_promised_amount от нода-генератора блока
	32 : "new_max_promised_amounts",
	//  новый набор % от нода-генератора блока
	33 : "new_pct",
	// добавление новой валюты
	34 : "admin_add_currency",
	35 : "new_cf_project",
	// новая версия, которая кладется каждому в диру public
	36 : "admin_new_version",
	// после того, как новая версия протестируется, выдаем сообщение, что необходимо обновиться
	37 : "admin_new_version_alert",
	// баг репорты
	38 : "message_to_admin",
	// админ может ответить юзеру
	39 : "admin_answer",
	40 : "cf_project_data",
	// блог админа
	41 : "admin_blog",
	// майнер меняет свой хост
	42 : "change_host",
	// майнер меняет комиссию, которую он хочет получать с тр-ий
	43 : "change_commission",
	44 : "del_cf_funding",
	// запуск урезания на основе голосования. генерит нод-генератор блока
	45 : "new_reduction",
	46 : "del_cf_project",
	47 : "cf_comment",
	48 : "cf_send_dc",
	49 : "user_avatar",
	50 : "cf_project_change_category",
	51 : "change_creditor",
	52 : "del_credit",
	53 : "repayment_credit",
	54 : "change_credit_part",
	55 : "new_admin",
	// по истечении 30 дней после поступления запроса о восстановлении утерянного ключа, админ может изменить ключ юзера
	56 : "admin_change_primary_key",
	// юзер разрешает или отменяет разрешение на смену своего ключа админом
	57 : "change_key_active",
	// юзер отменяет запрос на смену ключа
	58 : "change_key_close",
	// юзер отправляет с другого акка запрос на получение доступа к акку, ключ к которому он потерял
	59 : "change_key_request",
	// юзер решил стать арбитром или же действующий арбитр меняет комиссии
	60 : "change_arbitrator_conditions",
	// продавец меняет % и кол-во дней для новых сделок.
	61 : "change_seller_hold_back",
	// покупатель или продавец указал список арбитров, кому доверяет
	62 : "change_arbitrator_list",
	// покупатель хочет манибэк
	63 : "money_back_request",
	// магазин добровольно делает манибэк или же арбитр делает манибек
	64 : "money_back",
	// арбитр увеличивает время манибэка, чтобы успеть разобраться в ситуации
	65 : "change_money_back_time",
	// юзер меняет url центров сертификации, где хранятся его приватные ключи
	66 : "change_ca",
}

func init() {
}
