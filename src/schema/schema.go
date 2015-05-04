package schema

import (
	"fmt"
	"strings"
	"regexp"
	"utils"
)

type dbSet struct {
	mysql string
	postgresql string
	sqlite string
	comment string
}
type recmap map[string]interface{}
type recmapi map[int]interface{}
type recmap2 map[string]string

func GetSchema(dbType string) string {
	/*
	schema:=make(map[string]map[string]map[string]dbSet)
	schema["log_arbitrator_conditions"] = make(map[string]map[string]dbSet)
	schema["log_arbitrator_conditions"]["fileds"] = make(map[string]dbSet)
	schema["log_arbitrator_conditions"]["fileds"]["log_id"].mysql = "bigint(20) unsigned NOT NULL AUTO_INCREMENT"
	schema["log_arbitrator_conditions"]["fileds"]["log_id"].postgresql = ""
	schema["log_arbitrator_conditions"]["fileds"]["log_id"].sqlite = "bigint(20) unsigned NOT NULL AUTO_INCREMENT"
	schema["log_arbitrator_conditions"]["fileds"]["log_id"].comment = ""

	schema["log_arbitrator_conditions"]["fileds"]["conditions"].mysql = "text NOT NULL"
	schema["log_arbitrator_conditions"]["fileds"]["conditions"].postgresql = ""
	schema["log_arbitrator_conditions"]["fileds"]["conditions"].sqlite = "text NOT NULL"
	schema["log_arbitrator_conditions"]["fileds"]["conditions"].comment = ""

	//schema["log_arbitrator_conditions"]["table"] = "111111111"
	fmt.Println(schema)
	*/
	s:=make(recmap)
	s1:=make(recmap)
	s2:=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint  NOT NULL  default nextval('log_arbitrator_conditions_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"conditions", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"prev_log_id", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_arbitrator_conditions"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"conditions", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id"}
	s1["comment"] = ""
	s["arbitrator_conditions"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_change_ca"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_change_seller_hold_back"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_change_arbitrator_conditions"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_change_arbitration_trust_list"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_money_back_request"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"arbitrator_user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = "Список арбитров, кому доверяют юзеры"
	s["arbitration_trust_list"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('log_arbitration_trust_list_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"arbitration_trust_list", "mysql":"varchar(512) NOT NULL", "sqlite":"varchar(512) NOT NULL","postgresql":"varchar(512) NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"prev_log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_arbitration_trust_list"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('orders_id_seq')", "comment": "Этот ID указывается в тр-ии при запросе манибека"}
	s2[1] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Время блока, в котором запечатана данная сделка"}
	s2[2] = map[string]string{"name":"buyer", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "user_id покупателя"}
	s2[3] = map[string]string{"name":"seller", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "user_id продавца"}
	s2[4] = map[string]string{"name":"arbitrator0", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "user_id арбитра 0"}
	s2[5] = map[string]string{"name":"arbitrator1", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "user_id арбитра 1"}
	s2[6] = map[string]string{"name":"arbitrator2", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "user_id арбитра 2"}
	s2[7] = map[string]string{"name":"arbitrator3", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "user_id арбитра 3"}
	s2[8] = map[string]string{"name":"arbitrator4", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "user_id арбитра 4"}
	s2[9] = map[string]string{"name":"amount", "mysql":"decimal(15,2) NOT NULL", "sqlite":"decimal(15,2) NOT NULL","postgresql":"decimal(15,2) NOT NULL", "comment": "Сумма сделки"}
	s2[10] = map[string]string{"name":"hold_back_amount", "mysql":"decimal(15,2) NOT NULL", "sqlite":"decimal(15,2) NOT NULL","postgresql":"decimal(15,2) NOT NULL", "comment": "Сумма, которая замораживается на счету продавца. % для новых сделок задается в users.hold_back_pct"}
	s2[11] = map[string]string{"name":"currency_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[12] = map[string]string{"name":"end_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Время окончения возможности сделать манибек через арбитра. Может быть однократно увеличино арбитром. Используется для подсчета, сколько на данный момент времени есть активных сделок, чтобы посчитать от них 10% и не дать их списать со счета продавца"}
	s2[13] = map[string]string{"name":"end_time_changed", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": "Если арбитр изменил время окончания, то тут будет 1, чтобы нельзя было изменить повторно"}
	s2[14] = map[string]string{"name":"status", "mysql":"enum('normal','refund') NOT NULL DEFAULT 'normal'", "sqlite":"varchar(100)  NOT NULL DEFAULT 'normal'","postgresql":"enum('normal','refund') NOT NULL DEFAULT 'normal'", "comment": "Чтобы арбитр мог понять, что покупатель сделал запрос манибека. Когда юзер шлет тр-ию с запросом манибека, тут меняется статус на refund"}
	s2[15] = map[string]string{"name":"refund", "mysql":"decimal(15,2) NOT NULL", "sqlite":"decimal(15,2) NOT NULL","postgresql":"decimal(15,2) NOT NULL", "comment": "Сумма к возврату, которую определил арбитр. Она не может быть больше, чем сумма сделки.  Повторно отправить транзакцию с  манибеком  не даем, дабы не захламлять тр-ми сеть"}
	s2[16] = map[string]string{"name":"refund_arbitrator_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Для статы. ID арбитра, который сделал манибек юзеру"}
	s2[17] = map[string]string{"name":"arbitrator_refund_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Для статы. Время, когда арбитр сделал манибек"}
	s2[18] = map[string]string{"name":"voluntary_refund", "mysql":"decimal(15,2) NOT NULL", "sqlite":"decimal(15,2) NOT NULL","postgresql":"decimal(15,2) NOT NULL", "comment": "Сумма, которую продавец добровольно вернул покупателю.  Повторно отправить транзакцию с добровольным манибеком  не даем, дабы не захламлять тр-ми сеть. Если сумма не вся, то арбитр может довести процесс до конца, если посчитает нужным"}
	s2[19] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[20] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = ""
	s["orders"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('log_orders_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"end_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"status", "mysql":"enum('normal','refund') NOT NULL DEFAULT 'normal'", "sqlite":"varchar(100)  NOT NULL DEFAULT 'normal'","postgresql":"enum('normal','refund') NOT NULL DEFAULT 'normal'", "comment": ""}
	s2[3] = map[string]string{"name":"end_time_changed", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"refund", "mysql":"decimal(15,2) NOT NULL", "sqlite":"decimal(15,2) NOT NULL","postgresql":"decimal(15,2) NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"arbitrator_refund_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"refund_arbitrator_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[7] = map[string]string{"name":"voluntary_refund", "mysql":"decimal(15,2) NOT NULL", "sqlite":"decimal(15,2) NOT NULL","postgresql":"decimal(15,2) NOT NULL", "comment": ""}
	s2[8] = map[string]string{"name":"remaining_refund", "mysql":"decimal(15,2) NOT NULL", "sqlite":"decimal(15,2) NOT NULL","postgresql":"decimal(15,2) NOT NULL", "comment": ""}
	s2[9] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[10] = map[string]string{"name":"prev_log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_orders"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"referral", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"amount", "mysql":"decimal(15,2) NOT NULL", "sqlite":"decimal(15,2) NOT NULL","postgresql":"decimal(15,2) NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"currency_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = "Для вывода статы по рефам"
	s["referral_stats"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"hash", "mysql":"binary(16) NOT NULL", "sqlite":"binary(16) NOT NULL","postgresql":"bytea  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"type", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"error", "mysql":"varchar(255) NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"hash"}
	s1["comment"] = "Для удобства незарегенных юзеров на пуле. Показываем им статус их тр-ий"
	s["transactions_status"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"block_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"good", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"bad", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"block_id"}
	s1["comment"] = "Результаты сверки имеющегося у нас блока с блоками у случайных нодов"
	s["confirmations"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_change_key_request"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_change_key_active"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["admin"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s1["fileds"] = s2
	s1["comment"] = ""
	s["admin"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('log_admin_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"prev_log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_admin"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"admin_user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id"}
	s1["comment"] = ""
	s["votes_admin"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('log_votes_admin_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"admin_user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"prev_log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_votes_admin"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_new_credit"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_change_creditor"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_repayment_credit"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_change_credit_part"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('credits_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"del_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"amount", "mysql":"decimal(10,2) NOT NULL", "sqlite":"decimal(10,2) NOT NULL","postgresql":"decimal(10,2) NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"from_user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"to_user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"currency_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[7] = map[string]string{"name":"pct", "mysql":"decimal(5,2) NOT NULL", "sqlite":"decimal(5,2) NOT NULL","postgresql":"decimal(5,2) NOT NULL", "comment": ""}
	s2[8] = map[string]string{"name":"monthly_payment", "mysql":"decimal(10,2) NOT NULL", "sqlite":"decimal(10,2) NOT NULL","postgresql":"decimal(10,2) NOT NULL", "comment": "Ежемесячный платеж по кредиту. Пока не используется"}
	s2[9] = map[string]string{"name":"last_payment", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Время последнего платежа по кредиту. Пока не используется"}
	s2[10] = map[string]string{"name":"surety_1", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Поручитель 1. Пока не используется"}
	s2[11] = map[string]string{"name":"surety_2", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Поручитель 2. Пока не используется"}
	s2[12] = map[string]string{"name":"surety_3", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Поручитель 3. Пока не используется"}
	s2[13] = map[string]string{"name":"surety_4", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Поручитель 4. Пока не используется"}
	s2[14] = map[string]string{"name":"surety_5", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Поручитель 5. Пока не используется"}
	s2[15] = map[string]string{"name":"tx_hash", "mysql":"binary(16) NOT NULL", "sqlite":"binary(16) NOT NULL","postgresql":"bytea  NOT NULL", "comment": ""}
	s2[16] = map[string]string{"name":"tx_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[17] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = ""
	s["credits"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('log_credits_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"amount", "mysql":"decimal(10,2) NOT NULL", "sqlite":"decimal(10,2) NOT NULL","postgresql":"decimal(10,2) NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"to_user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"last_payment", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"tx_hash", "mysql":"binary(16) NOT NULL", "sqlite":"binary(16) NOT NULL","postgresql":"bytea  NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"tx_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[7] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_credits"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"project_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"ps1", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"ps2", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"ps3", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"ps4", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"ps5", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"ps6", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[7] = map[string]string{"name":"ps7", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[8] = map[string]string{"name":"ps8", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"project_id"}
	s1["comment"] = "Каждому CF-проекту вручную указывается платежные системы"
	s["cf_projects_ps"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"project_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"project_id"}
	s1["comment"] = "Какие проекты не выводим в CF-каталоге"
	s["cf_blacklist"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"email", "mysql":"varchar(200) CHARACTER SET utf8 NOT NULL", "sqlite":"varchar(200) NOT NULL","postgresql":"varchar(200) NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"email"}
	s1["comment"] = ""
	s["pool_waiting_list"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('cf_lang_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"name", "mysql":"varchar(200) CHARACTER SET utf8 NOT NULL", "sqlite":"varchar(200) NOT NULL","postgresql":"varchar(200) NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = ""
	s["cf_lang"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s1["fileds"] = s2
	s1["comment"] = ""
	s["cf_lang"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_user_avatar"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_cf_comments"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_new_cf_project"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_cf_project_data"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_cf_send_dc"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('cf_comments_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"project_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"lang_id", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"comment", "mysql":"varchar(255) CHARACTER SET utf8 NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Для того, чтобы можно было отсчитать время до размещения следующего коммента"}
	s2[6] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Для откатов"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = ""
	s["cf_comments"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('cf_funding_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"project_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"amount", "mysql":"decimal(15,2) NOT NULL", "sqlite":"decimal(15,2) NOT NULL","postgresql":"decimal(15,2) NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"amount_backup", "mysql":"decimal(15,2) NOT NULL", "sqlite":"decimal(15,2) NOT NULL","postgresql":"decimal(15,2) NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"currency_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "DC растут с юзерским %"}
	s2[7] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Для откатов"}
	s2[8] = map[string]string{"name":"del_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Фундер передумал и до завершения проекта вернул деньги"}
	s2[9] = map[string]string{"name":"checked", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": "Для определения по крону cf_project.funding и cf_project.funding"}
	s2[10] = map[string]string{"name":"del_checked", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": "Для определения по крону cf_project.funding и cf_project.funding"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = ""
	s["cf_funding"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('cf_currency_id_seq')", "comment": "ID идет от 1000, чтобы id CF-валют не пересекались с DC-валютами"}
	s2[1] = map[string]string{"name":"name", "mysql":"char(7) NOT NULL", "sqlite":"char(7) NOT NULL","postgresql":"char(7) NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"project_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["AI_START"] = "1000"
	s1["comment"] = ""
	s["cf_currency"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('cf_projects_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"currency_id", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"amount", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"funding", "mysql":"decimal(15,2) NOT NULL", "sqlite":"decimal(15,2) NOT NULL","postgresql":"decimal(15,2) NOT NULL", "comment": "Получаем в кроне. Сколько собрано средств. Нужно для вывода проектов в каталоге, чтобы не дергать cf_funding"}
	s2[5] = map[string]string{"name":"funders", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Получаем в кроне. Кол-во инвесторов. Нужно для вывода проектов в каталоге, чтобы не дергать cf_funding"}
	s2[6] = map[string]string{"name":"project_currency_name", "mysql":"char(7) NOT NULL", "sqlite":"char(7) NOT NULL","postgresql":"char(7) NOT NULL", "comment": ""}
	s2[7] = map[string]string{"name":"start_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[8] = map[string]string{"name":"end_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[9] = map[string]string{"name":"latitude", "mysql":"decimal(8,5) NOT NULL", "sqlite":"decimal(8,5) NOT NULL","postgresql":"decimal(8,5) NOT NULL", "comment": ""}
	s2[10] = map[string]string{"name":"longitude", "mysql":"decimal(8,5) NOT NULL", "sqlite":"decimal(8,5) NOT NULL","postgresql":"decimal(8,5) NOT NULL", "comment": ""}
	s2[11] = map[string]string{"name":"country", "mysql":"varchar(100) NOT NULL", "sqlite":"varchar(100) NOT NULL","postgresql":"varchar(100) NOT NULL", "comment": ""}
	s2[12] = map[string]string{"name":"city", "mysql":"varchar(100) NOT NULL", "sqlite":"varchar(100) NOT NULL","postgresql":"varchar(100) NOT NULL", "comment": ""}
	s2[13] = map[string]string{"name":"category_id", "mysql":"smallint(6) NOT NULL", "sqlite":"smallint(6) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[14] = map[string]string{"name":"close_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Чтобы знать, когда проект завершился и можно было бы удалить старые данные из cf_funding. Также используется для определения статус проекта - открыт/закрыт"}
	s2[15] = map[string]string{"name":"del_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Проект был закрыт автором, а средства возвращены инвесторам"}
	s2[16] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Для откатов"}
	s2[17] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[18] = map[string]string{"name":"geo_checked", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": "По крону превращаем координаты в названия страны и города и отмечаем тут"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = ""
	s["cf_projects"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('cf_projects_data_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"hide", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"project_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"lang_id", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"blurb_img", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"head_img", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"description_img", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[7] = map[string]string{"name":"picture", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": "Если нет видео, то выводится эта картинка"}
	s2[8] = map[string]string{"name":"video_type", "mysql":"varchar(10) NOT NULL", "sqlite":"varchar(10) NOT NULL","postgresql":"varchar(10) NOT NULL", "comment": ""}
	s2[9] = map[string]string{"name":"video_url_id", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[10] = map[string]string{"name":"news_img", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[11] = map[string]string{"name":"links", "mysql":"varchar(512) NOT NULL", "sqlite":"varchar(512) NOT NULL","postgresql":"varchar(512) NOT NULL", "comment": ""}
	s2[12] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["UNIQ"] = []string{"project_id","lang_id"}
	s1["AI"] = "id"
	s1["comment"] = ""
	s["cf_projects_data"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('log_cf_projects_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"category_id", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_cf_projects"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('log_cf_projects_data_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"hide", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"lang_id", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"blurb_img", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"head_img", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"description_img", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"picture", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[7] = map[string]string{"name":"video_type", "mysql":"varchar(10) NOT NULL", "sqlite":"varchar(10) NOT NULL","postgresql":"varchar(10) NOT NULL", "comment": ""}
	s2[8] = map[string]string{"name":"video_url_id", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[9] = map[string]string{"name":"news_img", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[10] = map[string]string{"name":"links", "mysql":"varchar(512) NOT NULL", "sqlite":"varchar(512) NOT NULL","postgresql":"varchar(512) NOT NULL", "comment": ""}
	s2[11] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[12] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_cf_projects_data"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"from_user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"comment", "mysql":"varchar(255) CHARACTER SET utf8 NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"time", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = "Абузы на майнеров от майнеров"
	s["abuses"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('admin_blog_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"lng", "mysql":"varchar(5) NOT NULL", "sqlite":"varchar(5) NOT NULL","postgresql":"varchar(5) NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"title", "mysql":"varchar(255) CHARACTER SET utf8 NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"message", "mysql":"text CHARACTER SET utf8 NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = "Блог админа"
	s["admin_blog"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('alert_messages_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"notification", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"close", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": "Юзер может закрыть сообщение и оно больше не появится"}
	s2[3] = map[string]string{"name":"message", "mysql":"text CHARACTER SET utf8 NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": "json. Каждому языку свое сообщение и gen - для тех, на кого языков не хватило"}
	s2[4] = map[string]string{"name":"currency_list", "mysql":"varchar(1024) NOT NULL", "sqlite":"varchar(1024) NOT NULL","postgresql":"varchar(1024) NOT NULL", "comment": "Для каких валют выводим сообщение. ALL - всем"}
	s2[5] = map[string]string{"name":"block_id", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Для откатов"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = "Сообщения от админа, которые выводятся в интерфейсе софта"
	s["alert_messages"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"hash", "mysql":"binary(32) NOT NULL", "sqlite":"binary(32) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "Хэш от полного заголовка блока (new_block_id,prev_block_hash,merkle_root,time,user_id,level). Используется как PREV_BLOCK_HASH"}
	s2[2] = map[string]string{"name":"head_hash", "mysql":"binary(32) NOT NULL", "sqlite":"binary(32) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "Хэш от заголовка блока (user_id,block_id,prev_head_hash). Используется для обновления head_hash в info_block при восстановлении после вилки в upd_block_info()"}
	s2[3] = map[string]string{"name":"data", "mysql":"longblob NOT NULL", "sqlite":"longblob NOT NULL","postgresql":"bytea NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["comment"] = "Главная таблица. Хранит цепочку блоков"
	s["block_chain"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"bigint(20) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint  NOT NULL  default nextval('cash_requests_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Время создания запроса. От него отсчитываем 48 часов"}
	s2[2] = map[string]string{"name":"from_user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"to_user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"notification", "mysql":"tinyint(1) unsigned NOT NULL", "sqlite":"tinyint(1)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"amount", "mysql":"decimal(13,2)  NOT NULL", "sqlite":"decimal(13,2)  NOT NULL","postgresql":"decimal(13,2)  NOT NULL", "comment": "На эту сумму должны быть выданы наличные"}
	s2[7] = map[string]string{"name":"hash_code", "mysql":"binary(32) NOT NULL", "sqlite":"binary(32) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "Хэш от кода, а сам код пердается при личной встрече. "}
	s2[8] = map[string]string{"name":"status", "mysql":"enum('approved','pending') NOT NULL DEFAULT 'pending'", "sqlite":"varchar(100)  NOT NULL DEFAULT 'pending'","postgresql":"enum('approved','pending') NOT NULL DEFAULT 'pending'", "comment": "Если в блоке указан верный код для хэша, то тут будет approved. Rejected нет, т.к. можно и без него понять, что запрос невыполнен, просто посмотрев время"}
	s2[9] = map[string]string{"name":"for_repaid_del_block_id", "mysql":"int(11)  unsigned NOT NULL", "sqlite":"int(11)   NOT NULL","postgresql":"int   NOT NULL", "comment": "если больше нет for_repaid ни по одной валюте у данного юзера, то нужно проверить, нет ли у него просроченных cash_requests, которым нужно отметить for_repaid_del_block_id, чтобы cash_request_out не переводил более обещанные суммы данного юзера в for_repaid из-за просроченных cash_requests"}
	s2[10] = map[string]string{"name":"del_block_id", "mysql":"int(11)  unsigned NOT NULL", "sqlite":"int(11)   NOT NULL","postgresql":"int   NOT NULL", "comment": "Во время reduction все текущие cash_requests, т.е. по которым не прошло 2 суток удаляются"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = "Запросы на обмен DC на наличные"
	s["cash_requests"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"tinyint(3) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"smallint  NOT NULL  default nextval('currency_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"name", "mysql":"char(3) NOT NULL", "sqlite":"char(3) NOT NULL","postgresql":"char(3) NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"full_name", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"max_other_currencies", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": "Со сколькими валютами данная валюта может майниться"}
	s2[4] = map[string]string{"name":"tmp_curs", "mysql":"double NOT NULL", "sqlite":"double NOT NULL","postgresql":"money NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = ""
	s["currency"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"int(11) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int  NOT NULL  default nextval('log_currency_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"max_other_currencies", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[3] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_currency"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"name", "mysql":"char(15) NOT NULL", "sqlite":"char(15) NOT NULL","postgresql":"char(15) NOT NULL", "comment": "Кодовое обозначение демона"}
	s2[1] = map[string]string{"name":"script", "mysql":"char(40) NOT NULL", "sqlite":"char(40) NOT NULL","postgresql":"char(40) NOT NULL", "comment": "Название скрипта"}
	s2[2] = map[string]string{"name":"param", "mysql":"char(5) NOT NULL", "sqlite":"char(5) NOT NULL","postgresql":"char(5) NOT NULL", "comment": "Параметры для запуска"}
	s2[3] = map[string]string{"name":"pid", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Pid демона для детекта дублей"}
	s2[4] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Время последней активности демона"}
	s2[5] = map[string]string{"name":"first", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"memory", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[7] = map[string]string{"name":"restart", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": "Команда демону, что нужно выйти"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"script"}
	s1["comment"] = "Демоны"
	s["daemons"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"race", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": "Раса. От 1 до 3"}
	s2[2] = map[string]string{"name":"country", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"version", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Версия набора точек"}
	s2[4] = map[string]string{"name":"status", "mysql":"enum('pending','used') NOT NULL DEFAULT 'pending'", "sqlite":"varchar(100)  NOT NULL DEFAULT 'pending'","postgresql":"enum('pending','used') NOT NULL DEFAULT 'pending'", "comment": "При new_miner ставим pending, при отрицательном завершении юзерского голосования - pending. used ставится только если юзерское голосование завершилось положительно"}
	s2[5] = map[string]string{"name":"f1", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": "Отрезок 1"}
	s2[6] = map[string]string{"name":"f2", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[7] = map[string]string{"name":"f3", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[8] = map[string]string{"name":"f4", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[9] = map[string]string{"name":"f5", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[10] = map[string]string{"name":"f6", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[11] = map[string]string{"name":"f7", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[12] = map[string]string{"name":"f8", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[13] = map[string]string{"name":"f9", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[14] = map[string]string{"name":"f10", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[15] = map[string]string{"name":"f11", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[16] = map[string]string{"name":"f12", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[17] = map[string]string{"name":"f13", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[18] = map[string]string{"name":"f14", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[19] = map[string]string{"name":"f15", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[20] = map[string]string{"name":"f16", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[21] = map[string]string{"name":"f17", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[22] = map[string]string{"name":"f18", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[23] = map[string]string{"name":"f19", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[24] = map[string]string{"name":"f20", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[25] = map[string]string{"name":"p1", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[26] = map[string]string{"name":"p2", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[27] = map[string]string{"name":"p3", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[28] = map[string]string{"name":"p4", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[29] = map[string]string{"name":"p5", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[30] = map[string]string{"name":"p6", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[31] = map[string]string{"name":"p7", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[32] = map[string]string{"name":"p8", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[33] = map[string]string{"name":"p9", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[34] = map[string]string{"name":"p10", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[35] = map[string]string{"name":"p11", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[36] = map[string]string{"name":"p12", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[37] = map[string]string{"name":"p13", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[38] = map[string]string{"name":"p14", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[39] = map[string]string{"name":"p15", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[40] = map[string]string{"name":"p16", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[41] = map[string]string{"name":"p17", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[42] = map[string]string{"name":"p18", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[43] = map[string]string{"name":"p19", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[44] = map[string]string{"name":"p20", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[45] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id"}
	s1["comment"] = "Точки по каждому юзеру"
	s["faces"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('holidays_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"delete", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": "1-удалено. нужно для отката"}
	s2[3] = map[string]string{"name":"start_time", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"end_time", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = "Время, в которое майнер не получает %, т.к. отдыхает"
	s["holidays"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"hash", "mysql":"binary(32) NOT NULL", "sqlite":"binary(32) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "Хэш от полного заголовка блока (new_block_id,prev_block_hash,merkle_root,time,user_id,level). Используется как prev_hash"}
	s2[1] = map[string]string{"name":"head_hash", "mysql":"binary(32) NOT NULL", "sqlite":"binary(32) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "Хэш от заголовка блока (user_id,block_id,prev_head_hash)"}
	s2[2] = map[string]string{"name":"block_id", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Время создания блока"}
	s2[4] = map[string]string{"name":"level", "mysql":"tinyint(4) unsigned NOT NULL", "sqlite":"tinyint(4)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": "На каком уровне был сгенерирован блок"}
	s2[5] = map[string]string{"name":"current_version", "mysql":"varchar(50) NOT NULL DEFAULT '0.0.1'", "sqlite":"varchar(50) NOT NULL DEFAULT '0.0.1'","postgresql":"varchar(50) NOT NULL DEFAULT '0.0.1'", "comment": ""}
	s2[6] = map[string]string{"name":"sent", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": "Был ли блок отправлен нодам, указанным в nodes_connections"}
	s1["fileds"] = s2
	s1["comment"] = "Текущий блок, данные из которого мы уже занесли к себе"
	s["info_block"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('promised_amount_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"del_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"user_id", "mysql":"bigint(16) NOT NULL", "sqlite":"bigint(16) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"amount", "mysql":"decimal(13,2) NOT NULL", "sqlite":"decimal(13,2) NOT NULL","postgresql":"decimal(13,2) NOT NULL", "comment": "Обещанная сумма. На неё влияет reduction и она будет урезаться при обновлении max_promised_amount (очень важно на случай деноминации фиата). Если же статус = repaid, то тут храниться кол-во денег, которые майнер отдал. Нужно хранить только чтобы знать общую сумму и не превысить max_promised_amount. Для WOC  amount не нужен, т.к. WOC полностью зависит от max_promised_amount"}
	s2[4] = map[string]string{"name":"amount_backup", "mysql":"decimal(13,2) NOT NULL", "sqlite":"decimal(13,2) NOT NULL","postgresql":"decimal(13,2) NOT NULL", "comment": "Нужно для откатов при reduction"}
	s2[5] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"ps1", "mysql":"smallint(5) unsigned NOT NULL", "sqlite":"smallint(5)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": "ID платежной системы, в валюте которой он готов сделать перевод в случае входящего запроса"}
	s2[7] = map[string]string{"name":"ps2", "mysql":"smallint(5) unsigned NOT NULL", "sqlite":"smallint(5)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[8] = map[string]string{"name":"ps3", "mysql":"smallint(5) unsigned NOT NULL", "sqlite":"smallint(5)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[9] = map[string]string{"name":"ps4", "mysql":"smallint(5) unsigned NOT NULL", "sqlite":"smallint(5)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[10] = map[string]string{"name":"ps5", "mysql":"smallint(5) unsigned NOT NULL", "sqlite":"smallint(5)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[11] = map[string]string{"name":"start_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Используется, когда нужно узнать, кто имеет право голосовать за данную валюту, т.е. прошло ли 60 дней с момента получения статуса miner или repaid(учитывая время со статусом miner). Изменяется при каждой смене статуса. Сущетвует только со статусом mining и repaid. Это защита от атаки клонов, когда каким-то образом 100500 майнеров прошли проверку, добавили какую-то валюту и проголосовали за reduction 90%. 90 дней - это время админу, чтобы заметить и среагировать на такую атаку"}
	s2[12] = map[string]string{"name":"status", "mysql":"enum('pending','mining','rejected','repaid','change_geo','suspended') NOT NULL DEFAULT 'pending'", "sqlite":"varchar(100)  NOT NULL DEFAULT 'pending'","postgresql":"enum('pending','mining','rejected','repaid','change_geo','suspended') NOT NULL DEFAULT 'pending'", "comment": "pending - при первом добавлении или при повтороном запросе.  change_geo ставится когда идет смена местоположения, suspended - когда админ разжаловал майнера в юзеры. TDC набегают только когда статус mining, repaid с майнерским или же юзерским % (если статус майнера = passive_miner)"}
	s2[13] = map[string]string{"name":"status_backup", "mysql":"enum('pending','mining','rejected','repaid','change_geo','suspended','null') NOT NULL DEFAULT 'null'", "sqlite":"varchar(100)  NOT NULL DEFAULT 'null'","postgresql":"enum('pending','mining','rejected','repaid','change_geo','suspended','null') NOT NULL DEFAULT 'null'", "comment": "Когда админ банит майнера, то в status пишется suspended, а сюда - статус из  status"}
	s2[14] = map[string]string{"name":"tdc_amount", "mysql":"decimal(13,2) NOT NULL", "sqlite":"decimal(13,2) NOT NULL","postgresql":"decimal(13,2) NOT NULL", "comment": "Набежавшая сумма за счет % роста. Пересчитывается при переводе TDC на кошелек"}
	s2[15] = map[string]string{"name":"tdc_amount_backup", "mysql":"decimal(13,2) NOT NULL", "sqlite":"decimal(13,2) NOT NULL","postgresql":"decimal(13,2) NOT NULL", "comment": "Нужно для откатов при reduction"}
	s2[16] = map[string]string{"name":"tdc_amount_update", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Время обновления tdc_amount"}
	s2[17] = map[string]string{"name":"video_type", "mysql":"varchar(100) NOT NULL", "sqlite":"varchar(100) NOT NULL","postgresql":"varchar(100) NOT NULL", "comment": ""}
	s2[18] = map[string]string{"name":"video_url_id", "mysql":"varchar(255) NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": "Если пусто, то видео берем по ID юзера.flv. На видео майнер говорит, что хочет майнить выбранную валюту"}
	s2[19] = map[string]string{"name":"votes_start_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "При каждой смене местоположения начинается новое голосование. Менять местоположение можно не чаще раза в сутки"}
	s2[20] = map[string]string{"name":"votes_0", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[21] = map[string]string{"name":"votes_1", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[22] = map[string]string{"name":"woc_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Нужно для отката добавления woc"}
	s2[23] = map[string]string{"name":"cash_request_out_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Любой cash_request_out приводит к появлению данной записи у получателя запроса. Убирается она только после того, как у юзера не остается непогашенных cash_request-ов. Нужно для reduction_generator, чтобы учитывать только те обещанные суммы, которые еще не заморожены невыполенными cash_request-ами"}
	s2[24] = map[string]string{"name":"cash_request_out_time_backup", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Используется в new_reduction()"}
	s2[25] = map[string]string{"name":"cash_request_in_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Нужно для отката cash_request_in"}
	s2[26] = map[string]string{"name":"del_mining_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Нужно для отката del_promised_amount"}
	s2[27] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = ""
	s["promised_amount"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint  NOT NULL  default nextval('log_promised_amount_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"del_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"amount", "mysql":"decimal(13,2) NOT NULL", "sqlite":"decimal(13,2) NOT NULL","postgresql":"decimal(13,2) NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"amount_backup", "mysql":"decimal(13,2) NOT NULL", "sqlite":"decimal(13,2) NOT NULL","postgresql":"decimal(13,2) NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"start_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"status", "mysql":"enum('pending','mining','rejected','repaid','change_geo','suspended')  NOT NULL", "sqlite":"varchar(100)   NOT NULL","postgresql":"enum('pending','mining','rejected','repaid','change_geo','suspended')  NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"status_backup", "mysql":"enum('pending','mining','rejected','repaid','change_geo','suspended','null')  NOT NULL DEFAULT 'null'", "sqlite":"varchar(100)   NOT NULL DEFAULT 'null'","postgresql":"enum('pending','mining','rejected','repaid','change_geo','suspended','null')  NOT NULL DEFAULT 'null'", "comment": ""}
	s2[7] = map[string]string{"name":"tdc_amount", "mysql":"decimal(13,2)  NOT NULL", "sqlite":"decimal(13,2)  NOT NULL","postgresql":"decimal(13,2)  NOT NULL", "comment": ""}
	s2[8] = map[string]string{"name":"tdc_amount_update", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[9] = map[string]string{"name":"video_type", "mysql":"varchar(100) NOT NULL", "sqlite":"varchar(100) NOT NULL","postgresql":"varchar(100) NOT NULL", "comment": ""}
	s2[10] = map[string]string{"name":"video_url_id", "mysql":"varchar(255) NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": "Если пусто, то видео берем по ID юзера.flv. На видео майнер говорит, что хочет майнить выбранную валюту"}
	s2[11] = map[string]string{"name":"votes_start_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "При каждой смене местоположения начинается новое голосование"}
	s2[12] = map[string]string{"name":"votes_0", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[13] = map[string]string{"name":"votes_1", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[14] = map[string]string{"name":"cash_request_out_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[15] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[16] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_promised_amount"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('log_faces_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"race", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"country", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"version", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Версия набора точек"}
	s2[5] = map[string]string{"name":"status", "mysql":"enum('approved','rejected','pending') NOT NULL", "sqlite":"varchar(100)  NOT NULL","postgresql":"enum('approved','rejected','pending') NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"f1", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": "Отрезок 1"}
	s2[7] = map[string]string{"name":"f2", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[8] = map[string]string{"name":"f3", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[9] = map[string]string{"name":"f4", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[10] = map[string]string{"name":"f5", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[11] = map[string]string{"name":"f6", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[12] = map[string]string{"name":"f7", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[13] = map[string]string{"name":"f8", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[14] = map[string]string{"name":"f9", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[15] = map[string]string{"name":"f10", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[16] = map[string]string{"name":"f11", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[17] = map[string]string{"name":"f12", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[18] = map[string]string{"name":"f13", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[19] = map[string]string{"name":"f14", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[20] = map[string]string{"name":"f15", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[21] = map[string]string{"name":"f16", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[22] = map[string]string{"name":"f17", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[23] = map[string]string{"name":"f18", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[24] = map[string]string{"name":"f19", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[25] = map[string]string{"name":"f20", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[26] = map[string]string{"name":"p1", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[27] = map[string]string{"name":"p2", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[28] = map[string]string{"name":"p3", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[29] = map[string]string{"name":"p4", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[30] = map[string]string{"name":"p5", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[31] = map[string]string{"name":"p6", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[32] = map[string]string{"name":"p7", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[33] = map[string]string{"name":"p8", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[34] = map[string]string{"name":"p9", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[35] = map[string]string{"name":"p10", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[36] = map[string]string{"name":"p11", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[37] = map[string]string{"name":"p12", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[38] = map[string]string{"name":"p13", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[39] = map[string]string{"name":"p14", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[40] = map[string]string{"name":"p15", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[41] = map[string]string{"name":"p16", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[42] = map[string]string{"name":"p17", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[43] = map[string]string{"name":"p18", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[44] = map[string]string{"name":"p19", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[45] = map[string]string{"name":"p20", "mysql":"float NOT NULL", "sqlite":"float NOT NULL","postgresql":"float NOT NULL", "comment": ""}
	s2[46] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[47] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = "Точки по каждому юзеру"
	s["log_faces"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('log_miners_data_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"miner_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"reg_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"status", "mysql":"enum('miner','user','passive_miner','suspended_miner') NOT NULL", "sqlite":"varchar(100)  NOT NULL","postgresql":"enum('miner','user','passive_miner','suspended_miner') NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"node_public_key", "mysql":"varbinary(512) NOT NULL", "sqlite":"varbinary(512) NOT NULL","postgresql":"bytea  NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"face_hash", "mysql":"varchar(128) NOT NULL", "sqlite":"varchar(128) NOT NULL","postgresql":"varchar(128) NOT NULL", "comment": ""}
	s2[7] = map[string]string{"name":"profile_hash", "mysql":"varchar(128) NOT NULL", "sqlite":"varchar(128) NOT NULL","postgresql":"varchar(128) NOT NULL", "comment": ""}
	s2[8] = map[string]string{"name":"photo_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[9] = map[string]string{"name":"photo_max_miner_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[10] = map[string]string{"name":"miners_keepers", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[11] = map[string]string{"name":"face_coords", "mysql":"varchar(1024) NOT NULL", "sqlite":"varchar(1024) NOT NULL","postgresql":"varchar(1024) NOT NULL", "comment": ""}
	s2[12] = map[string]string{"name":"profile_coords", "mysql":"varchar(1024) NOT NULL", "sqlite":"varchar(1024) NOT NULL","postgresql":"varchar(1024) NOT NULL", "comment": ""}
	s2[13] = map[string]string{"name":"video_type", "mysql":"varchar(100) NOT NULL", "sqlite":"varchar(100) NOT NULL","postgresql":"varchar(100) NOT NULL", "comment": ""}
	s2[14] = map[string]string{"name":"video_url_id", "mysql":"varchar(255) NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": ""}
	s2[15] = map[string]string{"name":"host", "mysql":"varchar(255) CHARACTER SET utf8 NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": ""}
	s2[16] = map[string]string{"name":"latitude", "mysql":"decimal(8,5) NOT NULL", "sqlite":"decimal(8,5) NOT NULL","postgresql":"decimal(8,5) NOT NULL", "comment": ""}
	s2[17] = map[string]string{"name":"longitude", "mysql":"decimal(8,5) NOT NULL", "sqlite":"decimal(8,5) NOT NULL","postgresql":"decimal(8,5) NOT NULL", "comment": ""}
	s2[18] = map[string]string{"name":"country", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[19] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[20] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_miners_data"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"count", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Сколько новых транзакций сделал юзер за минуту"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id"}
	s1["comment"] = ""
	s["log_minute"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('log_recycle_bin_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"profile_file_name", "mysql":"varchar(64) NOT NULL", "sqlite":"varchar(64) NOT NULL","postgresql":"varchar(64) NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"face_file_name", "mysql":"varchar(64) NOT NULL", "sqlite":"varchar(64) NOT NULL","postgresql":"varchar(64) NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[5] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_recycle_bin"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('log_spots_compatibility_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"version", "mysql":"double NOT NULL", "sqlite":"double NOT NULL","postgresql":"money NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"example_spots", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"compatibility", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"segments", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"tolerances", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[7] = map[string]string{"name":"prev_log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_spots_compatibility"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_actualization"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = "можно создавать только 1 тр-ю с абузами за 24h"
	s["log_time_abuses"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_for_repaid_fix"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_commission"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = "Для учета кол-ва запр. на доб. / удал. / изменение promised_amount. Чистим кроном"
	s["log_time_promised_amount"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_cash_requests"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_change_geolocation"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_holidays"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_message_to_admin"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_mining"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_change_host"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_new_miner"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_new_user"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_node_key"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_primary_key"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = "Храним данные за 1 сутки"
	s["log_time_votes"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = "Лимиты для повторых запросов, за которые голосуют ноды"
	s["log_time_votes_miners"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = "Голоса от нодов"
	s["log_time_votes_nodes"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["log_time_votes_complex"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"hash", "mysql":"binary(16) NOT NULL", "sqlite":"binary(16) NOT NULL","postgresql":"bytea  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"hash"}
	s1["comment"] = "Храним данные за сутки, чтобы избежать дублей."
	s["log_transactions"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('log_users_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"name", "mysql":"varchar(30) CHARACTER SET utf8 NOT NULL", "sqlite":"varchar(30) NOT NULL","postgresql":"varchar(30) NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"avatar", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"public_key_0", "mysql":"varbinary(512) NOT NULL", "sqlite":"varbinary(512) NOT NULL","postgresql":"bytea  NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"public_key_1", "mysql":"varbinary(512) NOT NULL", "sqlite":"varbinary(512) NOT NULL","postgresql":"bytea  NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"public_key_2", "mysql":"varbinary(512) NOT NULL", "sqlite":"varbinary(512) NOT NULL","postgresql":"bytea  NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"ca1", "mysql":"varchar(30) NOT NULL", "sqlite":"varchar(30) NOT NULL","postgresql":"varchar(30) NOT NULL", "comment": ""}
	s2[7] = map[string]string{"name":"ca2", "mysql":"varchar(30) NOT NULL", "sqlite":"varchar(30) NOT NULL","postgresql":"varchar(30) NOT NULL", "comment": ""}
	s2[8] = map[string]string{"name":"ca3", "mysql":"varchar(30) NOT NULL", "sqlite":"varchar(30) NOT NULL","postgresql":"varchar(30) NOT NULL", "comment": ""}
	s2[9] = map[string]string{"name":"referral", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[10] = map[string]string{"name":"credit_part", "mysql":"decimal(5,2) NOT NULL", "sqlite":"decimal(5,2) NOT NULL","postgresql":"decimal(5,2) NOT NULL", "comment": ""}
	s2[11] = map[string]string{"name":"change_key", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[12] = map[string]string{"name":"change_key_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[13] = map[string]string{"name":"change_key_close", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[14] = map[string]string{"name":"seller_hold_back_pct", "mysql":"decimal(5,2) NOT NULL", "sqlite":"decimal(5,2) NOT NULL","postgresql":"decimal(5,2) NOT NULL", "comment": ""}
	s2[15] = map[string]string{"name":"arbitration_days_refund", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[16] = map[string]string{"name":"url", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[17] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[18] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_users"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('log_variables_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"data", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_variables"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "Кто голосует"}
	s2[1] = map[string]string{"name":"voting_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "За что голосует. тут может быть id geolocation и пр"}
	s2[2] = map[string]string{"name":"type", "mysql":"enum('votes_miners','promised_amount') NOT NULL", "sqlite":"varchar(100)  NOT NULL","postgresql":"enum('votes_miners','promised_amount') NOT NULL", "comment": "Нужно для voting_id"}
	s2[3] = map[string]string{"name":"del_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было удаление. Нужно для чистки по крону старых данных и для откатов."}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id","type"}
	s1["comment"] = "Чтобы 1 юзер не смог проголосовать 2 раза за одно и тоже"
	s["log_votes"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('log_wallets_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"amount", "mysql":"decimal(15,2) UNSIGNED NOT NULL", "sqlite":"decimal(15,2)  NOT NULL","postgresql":"decimal(15,2)  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"amount_backup", "mysql":"decimal(15,2)  NOT NULL", "sqlite":"decimal(15,2)  NOT NULL","postgresql":"decimal(15,2)  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"last_update", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[5] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "Id предыдщуего log_id, который запишем в wallet"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = "Таблица, где будет браться инфа при откате блока"
	s["log_wallets"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"lock_time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"script_name", "mysql":"varchar(100) NOT NULL", "sqlite":"varchar(100) NOT NULL","postgresql":"varchar(100) NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"uniq", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["UNIQ"] = []string{"uniq"}
	s1["comment"] = "Полная блокировка на поступление новых блоков/тр-ий"
	s["main_lock"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"miner_id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('miners_miner_id_seq')", "comment": "Если есть забаненные, то на их место становятся новички, т.о. все miner_id будут заняты без пробелов"}
	s2[1] = map[string]string{"name":"active", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": "1 - активен, 0 - забанен"}
	s2[2] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Без log_id нельзя определить, был ли апдейт в табле miners или же инсерт, т.к. по AUTO_INCREMENT не понять, т.к. обновление может быть в самой последней строке"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"miner_id"}
	s1["AI"] = "miner_id"
	s1["comment"] = ""
	s["miners"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('log_miners_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[2] = map[string]string{"name":"prev_log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_miners"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"miner_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Из таблицы miners"}
	s2[2] = map[string]string{"name":"reg_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Время, когда майнер получил miner_id по итогам голосования. Определеяется один раз и не меняется. Нужно, чтобы не давать новым майнерам генерить тр-ии регистрации новых юзеров и исходящих запросов"}
	s2[3] = map[string]string{"name":"ban_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке майнер был разжалован в suspended_miner. Нужно для исключения пересечения тр-ий разжалованного майнера и самой тр-ии разжалования"}
	s2[4] = map[string]string{"name":"status", "mysql":"enum('miner','user','passive_miner','suspended_miner') NOT NULL DEFAULT 'user'", "sqlite":"varchar(100)  NOT NULL DEFAULT 'user'","postgresql":"enum('miner','user','passive_miner','suspended_miner') NOT NULL DEFAULT 'user'", "comment": "Измнеения вызывают персчет TDC в promised_amount"}
	s2[5] = map[string]string{"name":"node_public_key", "mysql":"varbinary(512) NOT NULL", "sqlite":"varbinary(512) NOT NULL","postgresql":"bytea  NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"face_hash", "mysql":"varchar(128) NOT NULL", "sqlite":"varchar(128) NOT NULL","postgresql":"varchar(128) NOT NULL", "comment": "Хэш фото юзера"}
	s2[7] = map[string]string{"name":"profile_hash", "mysql":"varchar(128) NOT NULL", "sqlite":"varchar(128) NOT NULL","postgresql":"varchar(128) NOT NULL", "comment": ""}
	s2[8] = map[string]string{"name":"photo_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Блок, в котором было добавлено фото"}
	s2[9] = map[string]string{"name":"photo_max_miner_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "Макс. майнер id в момент добавления фото. Это и photo_block_id нужны для определения 10-и нодов, где лежат фото"}
	s2[10] = map[string]string{"name":"miners_keepers", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": "Скольким майнерам копируем фото юзера. По дефолту = 10"}
	s2[11] = map[string]string{"name":"face_coords", "mysql":"varchar(1024) NOT NULL", "sqlite":"varchar(1024) NOT NULL","postgresql":"varchar(1024) NOT NULL", "comment": ""}
	s2[12] = map[string]string{"name":"profile_coords", "mysql":"varchar(1024) NOT NULL", "sqlite":"varchar(1024) NOT NULL","postgresql":"varchar(1024) NOT NULL", "comment": ""}
	s2[13] = map[string]string{"name":"video_type", "mysql":"varchar(100) NOT NULL", "sqlite":"varchar(100) NOT NULL","postgresql":"varchar(100) NOT NULL", "comment": ""}
	s2[14] = map[string]string{"name":"video_url_id", "mysql":"varchar(255) NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": "Если пусто, то видео берем по ID юзера.flv"}
	s2[15] = map[string]string{"name":"host", "mysql":"varchar(255) CHARACTER SET utf8 NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": "адрес хоста или IP, где находится нод майнера"}
	s2[16] = map[string]string{"name":"latitude", "mysql":"decimal(8,5) NOT NULL", "sqlite":"decimal(8,5) NOT NULL","postgresql":"decimal(8,5) NOT NULL", "comment": "Местоположение можно сменить без проблем, но это одновременно ведет запуск голосования у promised_amount по всем валютам, где статус mining или hold"}
	s2[17] = map[string]string{"name":"longitude", "mysql":"decimal(8,5) NOT NULL", "sqlite":"decimal(8,5) NOT NULL","postgresql":"decimal(8,5) NOT NULL", "comment": ""}
	s2[18] = map[string]string{"name":"country", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[19] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id"}
	s1["comment"] = ""
	s["miners_data"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"version", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"alert", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"notification", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"version"}
	s1["comment"] = "Сюда пишется новая версия, которая загружена в public"
	s["new_version"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"ban_start", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"info", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id"}
	s1["comment"] = "Баним на 1 час тех, кто дает нам данные с ошибками"
	s["nodes_ban"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"host", "mysql":"varchar(100) NOT NULL", "sqlite":"varchar(100) NOT NULL","postgresql":"varchar(100) NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Чтобы получать открытый ключ, которым шифруем блоки и тр-ии"}
	s2[2] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "ID блока, который есть у данного нода. Чтобы слать ему только >="}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"host"}
	s1["comment"] = "Ноды, которым шлем данные и от которых принимаем данные"
	s["nodes_connection"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('pct_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Время блока, в котором были новые %"}
	s2[2] = map[string]string{"name":"notification", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"miner", "mysql":"decimal(13,13) NOT NULL", "sqlite":"decimal(13,13) NOT NULL","postgresql":"decimal(13,13) NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"user", "mysql":"decimal(13,13) NOT NULL", "sqlite":"decimal(13,13) NOT NULL","postgresql":"decimal(13,13) NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"block_id", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Нужно для откатов"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = "% майнера, юзера. На основе  pct_votes"
	s["pct"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('max_promised_amounts_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Время блока, в котором были новые max_promised_amount"}
	s2[2] = map[string]string{"name":"notification", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"amount", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"block_id", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Нужно для откатов"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = "На основе votes_max_promised_amount"
	s["max_promised_amounts"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"time"}
	s1["comment"] = "Время последнего обновления max_other_currencies_time в currency "
	s["max_other_currencies_time"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('reduction_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Время блока, в котором было произведено уполовинивание"}
	s2[2] = map[string]string{"name":"notification", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"type", "mysql":"enum('manual','auto') NOT NULL DEFAULT 'auto'", "sqlite":"varchar(100)  NOT NULL DEFAULT 'auto'","postgresql":"enum('manual','auto') NOT NULL DEFAULT 'auto'", "comment": ""}
	s2[5] = map[string]string{"name":"pct", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"block_id", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Нужно для откатов"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = "Когда была последняя процедура урезания для конкретной валюты. Чтобы отсчитывать 2 недели до следующей"
	s["reduction"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Нужно только для того, чтобы определять, голосовал ли юзер или нет. От этого зависит, будет он получать майнерский или юзерский %"}
	s2[2] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"pct", "mysql":"decimal(13,13) NOT NULL", "sqlite":"decimal(13,13) NOT NULL","postgresql":"decimal(13,13) NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id","currency_id"}
	s1["comment"] = "Голосвание за %. Каждые 14 дней пересчет"
	s["votes_miner_pct"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint  NOT NULL  default nextval('log_votes_miner_pct_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"time", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"pct", "mysql":"decimal(13,13) NOT NULL", "sqlite":"decimal(13,13) NOT NULL","postgresql":"decimal(13,13) NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[5] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_votes_miner_pct"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"pct", "mysql":"decimal(13,13) NOT NULL", "sqlite":"decimal(13,13) NOT NULL","postgresql":"decimal(13,13) NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id","currency_id"}
	s1["comment"] = "Голосвание за %. Каждые 14 дней пересчет"
	s["votes_user_pct"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint  NOT NULL  default nextval('log_votes_user_pct_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"pct", "mysql":"decimal(13,13) NOT NULL", "sqlite":"decimal(13,13) NOT NULL","postgresql":"decimal(13,13) NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[4] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_votes_user_pct"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Учитываются только свежие голоса, т.е. один голос только за одно урезание"}
	s2[2] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"pct", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id","currency_id"}
	s1["comment"] = "Голосвание за уполовинивание денежной массы. Каждые 14 дней пересчет"
	s["votes_reduction"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint  NOT NULL  default nextval('log_votes_reduction_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"time", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"pct", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[5] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_votes_reduction"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"amount", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Возможные варианты задаются в скрипте, иначе будут проблемы с поиском варианта-победителя"}
	s2[3] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id","currency_id"}
	s1["comment"] = ""
	s["votes_max_promised_amount"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint  NOT NULL  default nextval('log_votes_max_promised_amount_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"amount", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[4] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_votes_max_promised_amount"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"count", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Возможные варианты задаются в скрипте, иначе будут проблемы с поиском варианта-победителя"}
	s2[3] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id","currency_id"}
	s1["comment"] = ""
	s["votes_max_other_currencies"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint  NOT NULL  default nextval('log_votes_max_other_currencies_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"currency_id", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"count", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[4] = map[string]string{"name":"prev_log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_votes_max_other_currencies"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time_start", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": "От какого времени отсчитывается 1 месяц"}
	s2[2] = map[string]string{"name":"points", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Баллы, полученные майнером за голосования"}
	s2[3] = map[string]string{"name":"log_id", "mysql":"bigint(20)  NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id"}
	s1["comment"] = "Баллы майнеров, по которым решается - получат они майнерские % или юзерские"
	s["points"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('log_points_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"time_start", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"points", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[4] = map[string]string{"name":"prev_log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_points"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"time_start", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Время начала действия статуса. До какого времени действует данный статус определяем простым добавлением в массив времени, которое будет через 30 дней"}
	s2[2] = map[string]string{"name":"status", "mysql":"enum('user','miner') NOT NULL DEFAULT 'user'", "sqlite":"varchar(100)  NOT NULL DEFAULT 'user'","postgresql":"enum('user','miner') NOT NULL DEFAULT 'user'", "comment": ""}
	s2[3] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Нужно для удобного отката"}
	s1["fileds"] = s2
	s1["comment"] = "Статусы юзеров на основе подсчета points"
	s["points_status"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"head_hash", "mysql":"binary(32) NOT NULL", "sqlite":"binary(32) NOT NULL","postgresql":"bytea  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"hash", "mysql":"binary(32) NOT NULL", "sqlite":"binary(32) NOT NULL","postgresql":"bytea  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"user_id", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"head_hash","hash"}
	s1["comment"] = "Блоки, которые мы должны забрать у указанных нодов"
	s["queue_blocks"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"head_hash", "mysql":"binary(32) NOT NULL", "sqlite":"binary(32) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "Хэш от заголовка блока (user_id,block_id,prev_head_hash)"}
	s2[1] = map[string]string{"name":"data", "mysql":"longblob NOT NULL", "sqlite":"longblob NOT NULL","postgresql":"bytea NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"head_hash"}
	s1["comment"] = "Очередь на фронтальную проверку соревнующихся блоков"
	s["queue_testblock"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"hash", "mysql":"binary(16) NOT NULL", "sqlite":"binary(16) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "md5 от тр-ии"}
	s2[1] = map[string]string{"name":"high_rate", "mysql":"tinyint(1) NOT NULL DEFAULT '0'", "sqlite":"tinyint(1) NOT NULL DEFAULT '0'","postgresql":"smallint NOT NULL DEFAULT '0'", "comment": "Если 1, значит это админская тр-ия"}
	s2[2] = map[string]string{"name":"data", "mysql":"longblob NOT NULL", "sqlite":"longblob NOT NULL","postgresql":"bytea NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"_tmp_node_user_id", "mysql":"VARCHAR(255)", "sqlite":"VARCHAR(255)","postgresql":"VARCHAR(255)", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"hash"}
	s1["comment"] = "Тр-ии, которые мы должны проверить"
	s["queue_tx"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"profile_file_name", "mysql":"varchar(64) NOT NULL", "sqlite":"varchar(64) NOT NULL","postgresql":"varchar(64) NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"face_file_name", "mysql":"varchar(64) NOT NULL", "sqlite":"varchar(64) NOT NULL","postgresql":"varchar(64) NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id"}
	s1["comment"] = ""
	s["recycle_bin"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"version", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"example_spots", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": "Точки, которые наносим на 2 фото-примера (анфас и профиль)"}
	s2[2] = map[string]string{"name":"compatibility", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": "С какими версиями совместимо"}
	s2[3] = map[string]string{"name":"segments", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": "Нужно для составления отрезков в new_miner()"}
	s2[4] = map[string]string{"name":"tolerances", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": "Допустимые расхождения между точками при поиске фото-дублей"}
	s2[5] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"version"}
	s1["comment"] = "Совместимость текущей версии точек с предыдущими"
	s["spots_compatibility"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "ID тестируемого блока"}
	s2[1] = map[string]string{"name":"time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": "Время, когда блок попал сюда"}
	s2[2] = map[string]string{"name":"level", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": "Пишем сюда для использования при формировании заголовка"}
	s2[3] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "По id вычисляем хэш шапки"}
	s2[4] = map[string]string{"name":"header_hash", "mysql":"binary(32) NOT NULL", "sqlite":"binary(32) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "Хэш шапки, им меряемся, у кого меньше - тот круче. Хэш генерим у себя, при получении данных блока"}
	s2[5] = map[string]string{"name":"signature", "mysql":"blob NOT NULL", "sqlite":"blob NOT NULL","postgresql":"bytea NOT NULL", "comment": "Подпись блока юзером, чей минимальный хэш шапки мы приняли"}
	s2[6] = map[string]string{"name":"mrkl_root", "mysql":"binary(32) NOT NULL", "sqlite":"binary(32) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "Хэш тр-ий. Чтобы каждый раз не проверять теже самые данные, просто сравниваем хэши"}
	s2[7] = map[string]string{"name":"status", "mysql":"enum('active','pending') NOT NULL DEFAULT 'active'", "sqlite":"varchar(100)  NOT NULL DEFAULT 'active'","postgresql":"enum('active','pending') NOT NULL DEFAULT 'active'", "comment": "Указание скрипту testblock_disseminator.php"}
	s2[8] = map[string]string{"name":"uniq", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"block_id"}
	s1["UNIQ"] = []string{"uniq"}
	s1["comment"] = "Нужно на этапе соревнования, у кого меньше хэш"
	s["testblock"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"lock_time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"script_name", "mysql":"varchar(30) NOT NULL", "sqlite":"varchar(30) NOT NULL","postgresql":"varchar(30) NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"uniq", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["UNIQ"] = []string{"uniq"}
	s1["comment"] = ""
	s["testblock_lock"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"hash", "mysql":"binary(16) NOT NULL", "sqlite":"binary(16) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "Все хэши из этой таблы шлем тому, у кого хотим получить блок (т.е. недостающие тр-ии для составления блока)"}
	s2[1] = map[string]string{"name":"data", "mysql":"longblob NOT NULL", "sqlite":"longblob NOT NULL","postgresql":"bytea NOT NULL", "comment": "Само тело тр-ии"}
	s2[2] = map[string]string{"name":"verified", "mysql":"tinyint(1) NOT NULL DEFAULT '1'", "sqlite":"tinyint(1) NOT NULL DEFAULT '1'","postgresql":"smallint NOT NULL DEFAULT '1'", "comment": "Оставшиеся после прихода нового блока тр-ии отмечаются как \"непроверенные\" и их нужно проверять по новой"}
	s2[3] = map[string]string{"name":"used", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": "После того как попадют в блок, ставим 1, а те, у которых уже стояло 1 - удаляем"}
	s2[4] = map[string]string{"name":"high_rate", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": "1 - админские, 0 - другие"}
	s2[5] = map[string]string{"name":"for_self_use", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": "для new_pct(pct_generator.php), т.к. эта тр-ия валидна только вместе с блоком, который сгенерил тот, кто сгенерил эту тр-ию"}
	s2[6] = map[string]string{"name":"type", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": "Тип тр-ии. Нужно для недопущения попадения в блок 2-х тр-ий одного типа от одного юзера"}
	s2[7] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Нужно для недопущения попадения в блок 2-х тр-ий одного типа от одного юзера"}
	s2[8] = map[string]string{"name":"third_var", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Для исключения пересения в одном блоке удаления обещанной суммы и запроса на её обмен на DC. И для исключения голосования за один и тот же объект одним и тем же юзеров и одном блоке"}
	s2[9] = map[string]string{"name":"counter", "mysql":"tinyint(3) NOT NULL", "sqlite":"tinyint(3) NOT NULL","postgresql":"smallint NOT NULL", "comment": "Чтобы избежать зацикливания при проверке тр-ии: verified=1, новый блок, verified=0. При достижении 10-и - удаляем тр-ию "}
	s2[10] = map[string]string{"name":"sent", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": "Была отправлена нодам, указанным в nodes_connections"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"hash"}
	s1["comment"] = "Все незанесенные в блок тр-ии, которые у нас есть"
	s["transactions"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"int(11) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"int NOT NULL  default nextval('transactions_testblock_id_seq')", "comment": "Порядок следования очень важен"}
	s2[1] = map[string]string{"name":"hash", "mysql":"binary(16) NOT NULL", "sqlite":"binary(16) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "md5 для обмена только недостающими тр-ми"}
	s2[2] = map[string]string{"name":"data", "mysql":"longblob NOT NULL", "sqlite":"longblob NOT NULL","postgresql":"bytea NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"type", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": "Тип тр-ии. Нужно для недопущения попадения в блок 2-х тр-ий одного типа от одного юзера"}
	s2[4] = map[string]string{"name":"user_id", "mysql":"tinyint(4) NOT NULL", "sqlite":"tinyint(4) NOT NULL","postgresql":"smallint NOT NULL", "comment": "Нужно для недопущения попадения в блок 2-х тр-ий одного типа от одного юзера"}
	s2[5] = map[string]string{"name":"third_var", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Для исключения пересения в одном блоке удаления обещанной суммы и запроса на её обмен на DC. И для исключения голосования за один и тот же объект одним и тем же юзеров и одном блоке"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["UNIQ"] = []string{"hash"}
	s1["AI"] = "id"
	s1["comment"] = "Тр-ии, которые используются в текущем testblock"
	s["transactions_testblock"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint  NOT NULL  default nextval('log_commission_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"commission", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "В каком блоке было занесено. Нужно для удаления старых данных"}
	s2[3] = map[string]string{"name":"prev_log_id", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = "Каждый майнер определяет, какая комиссия с тр-ий будет доставаться ему, если он будет генерить блок"
	s["log_commission"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"commission", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": "Комиссии по всем валютам в json. Если какой-то валюты нет в списке, то комиссия будет равна нулю. currency_id, %, мин., макс."}
	s2[2] = map[string]string{"name":"log_id", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id"}
	s1["comment"] = "Каждый майнер определяет, какая комиссия с тр-ий будет доставаться ему, если он будет генерить блок"
	s["commission"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint  NOT NULL  default nextval('users_user_id_seq')", "comment": "На него будут слаться деньги"}
	s2[1] = map[string]string{"name":"name", "mysql":"varchar(30) CHARACTER SET utf8 NOT NULL", "sqlite":"varchar(30) NOT NULL","postgresql":"varchar(30) NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"avatar", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"public_key_0", "mysql":"varbinary(512) NOT NULL", "sqlite":"varbinary(512) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "Открытый ключ которым проверяются все транзакции от юзера"}
	s2[4] = map[string]string{"name":"public_key_1", "mysql":"varbinary(512) NOT NULL", "sqlite":"varbinary(512) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "2-й ключ, если есть"}
	s2[5] = map[string]string{"name":"public_key_2", "mysql":"varbinary(512) NOT NULL", "sqlite":"varbinary(512) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "3-й ключ, если есть"}
	s2[6] = map[string]string{"name":"ca1", "mysql":"varchar(30) NOT NULL", "sqlite":"varchar(30) NOT NULL","postgresql":"varchar(30) NOT NULL", "comment": ""}
	s2[7] = map[string]string{"name":"ca2", "mysql":"varchar(30) NOT NULL", "sqlite":"varchar(30) NOT NULL","postgresql":"varchar(30) NOT NULL", "comment": ""}
	s2[8] = map[string]string{"name":"ca3", "mysql":"varchar(30) NOT NULL", "sqlite":"varchar(30) NOT NULL","postgresql":"varchar(30) NOT NULL", "comment": ""}
	s2[9] = map[string]string{"name":"referral", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "Тот, кто зарегал данного юзера и теперь получает с него рефские"}
	s2[10] = map[string]string{"name":"credit_part", "mysql":"decimal(5,2) NOT NULL", "sqlite":"decimal(5,2) NOT NULL","postgresql":"decimal(5,2) NOT NULL", "comment": "% от поступлений, которые юзер осталяет себе. Если есть активные кредиты, то можно только уменьшать"}
	s2[11] = map[string]string{"name":"change_key", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[12] = map[string]string{"name":"change_key_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[13] = map[string]string{"name":"change_key_close", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[14] = map[string]string{"name":"seller_hold_back_pct", "mysql":"decimal(5,2) NOT NULL", "sqlite":"decimal(5,2) NOT NULL","postgresql":"decimal(5,2) NOT NULL", "comment": "% холдбека для новых сделок"}
	s2[15] = map[string]string{"name":"arbitration_days_refund", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Продавец тут указывает кол-во дней для новых сделок, в течение которых он готов сделать манибек. Если стоит 0, значит продавец больше не работает с манибеком"}
	s2[16] = map[string]string{"name":"url", "mysql":"varchar(50) NOT NULL", "sqlite":"varchar(50) NOT NULL","postgresql":"varchar(50) NOT NULL", "comment": ""}
	s2[17] = map[string]string{"name":"log_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id"}
	s1["AI"] = "user_id"
	s1["comment"] = ""
	s["users"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"name", "mysql":"varchar(35) NOT NULL", "sqlite":"varchar(35) NOT NULL","postgresql":"varchar(35) NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"value", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"comment", "mysql":"varchar(255) CHARACTER SET utf8 NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"name"}
	s1["comment"] = ""
	s["variables"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"bigint(20) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint  NOT NULL  default nextval('votes_miners_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"type", "mysql":"enum('node_voting','user_voting') NOT NULL", "sqlite":"varchar(100)  NOT NULL","postgresql":"enum('node_voting','user_voting') NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": "За кого голосуем"}
	s2[3] = map[string]string{"name":"votes_start_time", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"votes_0", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"votes_1", "mysql":"int(11) unsigned NOT NULL", "sqlite":"int(11)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"votes_end", "mysql":"tinyint(1) unsigned NOT NULL", "sqlite":"tinyint(1)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[7] = map[string]string{"name":"end_block_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": "В каком блоке мы выставили принудительное end для node"}
	s2[8] = map[string]string{"name":"cron_checked_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "По крону проверили, не нужно ли нам скачать фотки юзера к себе на сервер"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = "Отдел. от miners_data, чтобы гол. шли точно за свежие данные"
	s["votes_miners"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"first", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"second", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"third", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id"}
	s1["comment"] = "Голосвание за рефские %. Каждые 14 дней пересчет"
	s["votes_referral"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint  NOT NULL  default nextval('log_votes_referral_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"first", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"second", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"third", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"prev_log_id", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_votes_referral"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"first", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"second", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"third", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"log_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["referral"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"log_id", "mysql":"bigint(20) unsigned NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint  NOT NULL  default nextval('log_referral_log_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"first", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"second", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"third", "mysql":"tinyint(2) unsigned NOT NULL", "sqlite":"tinyint(2)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"prev_log_id", "mysql":"int(10) unsigned NOT NULL", "sqlite":"int(10)  NOT NULL","postgresql":"int  NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"log_id"}
	s1["AI"] = "log_id"
	s1["comment"] = ""
	s["log_referral"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"progress", "mysql":"varchar(10) NOT NULL", "sqlite":"varchar(10) NOT NULL","postgresql":"varchar(10) NOT NULL", "comment": "На каком шаге остановились"}
	s1["fileds"] = s2
	s1["comment"] = "Используется только в момент установки"
	s["install"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('wallets_user_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"currency_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"amount", "mysql":"decimal(15,2) unsigned NOT NULL", "sqlite":"decimal(15,2)  NOT NULL","postgresql":"decimal(15,2)  NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"amount_backup", "mysql":"decimal(15,2) unsigned NOT NULL", "sqlite":"decimal(15,2)  NOT NULL","postgresql":"decimal(15,2)  NOT NULL", "comment": "Может неравномерно обнуляться из-за обработки, а затем - отката new_reduction()"}
	s2[4] = map[string]string{"name":"last_update", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Время последнего пересчета суммы с учетом % из miner_pct"}
	s2[5] = map[string]string{"name":"log_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "ID log_wallets, откуда будет брать данные при откате на 1 блок. 0 - значит при откате нужно удалить строку"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id","currency_id"}
	s1["AI"] = "user_id"
	s1["comment"] = "У кого сколько какой валюты"
	s["wallets"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"hash", "mysql":"binary(32) NOT NULL", "sqlite":"binary(32) NOT NULL","postgresql":"bytea  NOT NULL", "comment": "Хэш транзакции. Нужно для удаления данных из буфера, после того, как транзакция была обработана в блоке, либо анулирована из-за ошибок при повторной проверке"}
	s2[1] = map[string]string{"name":"del_block_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "Т.к. удалять нельзя из-за возможного отката блока, приходится делать delete=1, а через сутки - чистить"}
	s2[2] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[3] = map[string]string{"name":"currency_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": ""}
	s2[4] = map[string]string{"name":"amount", "mysql":"decimal(15,2) unsigned NOT NULL", "sqlite":"decimal(15,2)  NOT NULL","postgresql":"decimal(15,2)  NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"block_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "Может быть = 0. Номер блока, в котором была занесена запись. Если блок в процессе фронт. проверки окажется невалдиным, то просто удалим все данные по block_id"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"hash"}
	s1["comment"] = "Суммируем все списания, которые еще не в блоке"
	s["wallets_buffer"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('forex_orders_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "Чей ордер"}
	s2[2] = map[string]string{"name":"sell_currency_id", "mysql":"int(10) NOT NULL", "sqlite":"int(10) NOT NULL","postgresql":"int NOT NULL", "comment": "Что продается"}
	s2[3] = map[string]string{"name":"sell_rate", "mysql":"decimal(20,10) NOT NULL", "sqlite":"decimal(20,10) NOT NULL","postgresql":"decimal(20,10) NOT NULL", "comment": "По какому курсу к buy_currency_id"}
	s2[4] = map[string]string{"name":"amount", "mysql":"decimal(15,2)  NOT NULL", "sqlite":"decimal(15,2)  NOT NULL","postgresql":"decimal(15,2)  NOT NULL", "comment": "Сколько осталось на данном ордере"}
	s2[5] = map[string]string{"name":"amount_backup", "mysql":"decimal(15,2)  NOT NULL", "sqlite":"decimal(15,2)  NOT NULL","postgresql":"decimal(15,2)  NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"buy_currency_id", "mysql":"int(10) NOT NULL", "sqlite":"int(10) NOT NULL","postgresql":"int NOT NULL", "comment": "Какая валюта нужна"}
	s2[7] = map[string]string{"name":"commission", "mysql":"decimal(15,2)  NOT NULL", "sqlite":"decimal(15,2)  NOT NULL","postgresql":"decimal(15,2)  NOT NULL", "comment": "Какую отдали комиссию ноду-генератору"}
	s2[8] = map[string]string{"name":"empty_block_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "Если ордер опустошили, то тут будет номер блока. Чтобы потом удалить старые записи"}
	s2[9] = map[string]string{"name":"del_block_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "Если юзер решил удалить ордер, то тут будет номер блока"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = ""
	s["forex_orders"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('log_forex_orders_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"main_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "ID из log_forex_orders_main. Для откатов"}
	s2[2] = map[string]string{"name":"order_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": "Какой ордер был задействован. Для откатов"}
	s2[3] = map[string]string{"name":"amount", "mysql":"decimal(15,2) unsigned NOT NULL", "sqlite":"decimal(15,2)  NOT NULL","postgresql":"decimal(15,2)  NOT NULL", "comment": "Какая сумма была вычтена из ордера"}
	s2[4] = map[string]string{"name":"to_user_id", "mysql":"bigint(20) unsigned NOT NULL", "sqlite":"bigint(20)  NOT NULL","postgresql":"bigint  NOT NULL", "comment": "Какому юзеру была начислено amount "}
	s2[5] = map[string]string{"name":"new", "mysql":"tinyint(3) unsigned NOT NULL", "sqlite":"tinyint(3)  NOT NULL","postgresql":"smallint  NOT NULL", "comment": "Если 1, то был создан новый  ордер. при 1 amount не указывается, т.к. при откате будет просто удалена запись из forex_orders"}
	s2[6] = map[string]string{"name":"block_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "Для откатов"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = "Все ордеры, который были затронуты в результате тр-ии"
	s["log_forex_orders"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('log_forex_orders_main_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"block_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "Чтобы можно было понять, какие данные можно смело удалять из-за их давности"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = "Каждый ордер пишется сюда. При откате любого ордера просто берем последнюю строку отсюда"
	s["log_forex_orders_main"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"tx_hash", "mysql":"binary(16)", "sqlite":"binary(16)","postgresql":"bytea ", "comment": "По этому хэшу отмечается, что данная тр-ия попала в блок и ставится del_block_id"}
	s2[1] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[2] = map[string]string{"name":"del_block_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": "block_id сюда пишется в тот момент, когда тр-ия попала в блок и уже не используется для фронтальной проверки. Нужно чтобы можно было понять, какие данные можно смело удалять из-за их давности"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"tx_hash"}
	s1["comment"] = "В один блок не должно попасть более чем 10 тр-ий перевода средств или создания forex-ордеров на суммы менее эквивалента 0.05-0.1$ по текущему курсу"
	s["log_time_money_orders"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('_my_admin_messages_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"add_time", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "для удаления старых my_pending"}
	s2[2] = map[string]string{"name":"user_int_message_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "ID сообщения, который присылает юзер"}
	s2[3] = map[string]string{"name":"parent_user_int_message_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Parent_id, который присылает юзер"}
	s2[4] = map[string]string{"name":"user_id", "mysql":"bigint(20) NOT NULL", "sqlite":"bigint(20) NOT NULL","postgresql":"bigint NOT NULL", "comment": ""}
	s2[5] = map[string]string{"name":"type", "mysql":"enum('from_user','to_user') NOT NULL", "sqlite":"varchar(100)  NOT NULL","postgresql":"enum('from_user','to_user') NOT NULL", "comment": ""}
	s2[6] = map[string]string{"name":"subject", "mysql":"varchar(255) CHARACTER SET utf8 NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": ""}
	s2[7] = map[string]string{"name":"encrypted", "mysql":"blob NOT NULL", "sqlite":"blob NOT NULL","postgresql":"bytea NOT NULL", "comment": ""}
	s2[8] = map[string]string{"name":"decrypted", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[9] = map[string]string{"name":"message", "mysql":"text CHARACTER SET utf8 NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": ""}
	s2[10] = map[string]string{"name":"message_type", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[11] = map[string]string{"name":"message_subtype", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[12] = map[string]string{"name":"status", "mysql":"enum('my_pending','approved') NOT NULL DEFAULT 'my_pending'", "sqlite":"varchar(100)  NOT NULL DEFAULT 'my_pending'","postgresql":"enum('my_pending','approved') NOT NULL DEFAULT 'my_pending'", "comment": ""}
	s2[13] = map[string]string{"name":"close", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": "Воспрос закрыли, чтобы больше не маячил"}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = "Эта табла видна только админу"
	s["_my_admin_messages"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"hash", "mysql":"binary(16) NOT NULL", "sqlite":"binary(16) NOT NULL","postgresql":"bytea  NOT NULL", "comment": ""}
	s2[1] = map[string]string{"name":"data", "mysql":"varchar(20) NOT NULL", "sqlite":"varchar(20) NOT NULL","postgresql":"varchar(20) NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"hash"}
	s1["comment"] = ""
	s["authorization"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"user_id"}
	s1["comment"] = "Если не пусто, то работаем в режиме пула"
	s["community"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"uniq", "mysql":"enum('1') NOT NULL DEFAULT '1'", "sqlite":"varchar(100)  NOT NULL DEFAULT '1'","postgresql":"enum('1') NOT NULL DEFAULT '1'", "comment": ""}
	s2[1] = map[string]string{"name":"data", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"uniq"}
	s1["comment"] = ""
	s["backup_community"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"id", "mysql":"bigint(20) NOT NULL AUTO_INCREMENT", "sqlite":"INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL","postgresql":"bigint NOT NULL  default nextval('payment_systems_id_seq')", "comment": ""}
	s2[1] = map[string]string{"name":"name", "mysql":"varchar(255) NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": ""}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"id"}
	s1["AI"] = "id"
	s1["comment"] = "Для тех, кто не хочет встречаться для обмена кода на наличные"
	s["payment_systems"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s1["fileds"] = s2
	s1["comment"] = "Для тех, кто не хочет встречаться для обмена кода на наличные"
	s["payment_systems"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"ip", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Раз в минуту удаляется"}
	s2[1] = map[string]string{"name":"req", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Кол-во запросов от ip. "}
	s1["fileds"] = s2
	s1["PRIMARY"] = []string{"ip"}
	s1["comment"] = "Защита от случайного ддоса"
	s["ddos_protection"] = s1
	printSchema(s, dbType)

	s=make(recmap)
	s1=make(recmap)
	s2=make(recmapi)
	s2[0] = map[string]string{"name":"php_path", "mysql":"varchar(255) NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": "Нужно для запуска демонов"}
	s2[1] = map[string]string{"name":"my_block_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Параллельно с info_block пишем и сюда. Нужно при обнулении рабочих таблиц, чтобы знать до какого блока не трогаем таблы my_"}
	s2[2] = map[string]string{"name":"local_gate_ip", "mysql":"varchar(255) NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": "Если тут не пусто, то connector.php будет не активным, а ip для disseminator.php будет браться тут. Нужно для защищенного режима"}
	s2[3] = map[string]string{"name":"static_node_user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Все исходящие тр-ии будут подписаны публичным ключом этой ноды. Нужно для защищенного режима"}
	s2[4] = map[string]string{"name":"in_connections_ip_limit", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Кол-во запросов от 1 ip за минуту"}
	s2[5] = map[string]string{"name":"in_connections", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Кол-во нодов и просто юзеров, от кого принимаем данные. Считаем кол-во ip за 1 минуту"}
	s2[6] = map[string]string{"name":"out_connections", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Кол-во нодов, кому шлем данные"}
	s2[7] = map[string]string{"name":"bad_blocks", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": "Номера и sign плохих блоков. Нужно, чтобы не подцепить более длинную, но глючную цепочку блоков"}
	s2[8] = map[string]string{"name":"pool_max_users", "mysql":"int(11) NOT NULL DEFAULT '100'", "sqlite":"int(11) NOT NULL DEFAULT '100'","postgresql":"int NOT NULL DEFAULT '100'", "comment": ""}
	s2[9] = map[string]string{"name":"pool_admin_user_id", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": ""}
	s2[10] = map[string]string{"name":"pool_tech_works", "mysql":"tinyint(1) NOT NULL", "sqlite":"tinyint(1) NOT NULL","postgresql":"smallint NOT NULL", "comment": ""}
	s2[11] = map[string]string{"name":"exchange_api_url", "mysql":"varchar(255) NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": "На home далается ajax-запрос к api биржи и выдается инфа о курсе и пр."}
	s2[12] = map[string]string{"name":"cf_url", "mysql":"varchar(255) NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": "URL, который отображается в соц. кнопках и с которого подгружаются css/js/img/fonts при прямом заходе в CF-каталог"}
	s2[13] = map[string]string{"name":"pool_url", "mysql":"varchar(255) NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": "URL, на который ссылается кнопка Contribute now из внешнего CF-каталога "}
	s2[14] = map[string]string{"name":"pool_email", "mysql":"varchar(255) NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": "В режиме пула используется как адрес отправителя при рассылке уведомлений"}
	s2[15] = map[string]string{"name":"cf_available_coins_url", "mysql":"varchar(255) NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": "URL биржи, где можно узнать, сколько там осталось монет в продаже по курсу 1"}
	s2[16] = map[string]string{"name":"cf_exchange_url", "mysql":"varchar(255) NOT NULL", "sqlite":"varchar(255) NOT NULL","postgresql":"varchar(255) NOT NULL", "comment": "URL биржи. Просто, чтобы дать на неё ссылку в сообщении, где говорится, что монеты на бирже кончились"}
	s2[17] = map[string]string{"name":"cf_top_html", "mysql":"text CHARACTER SET utf8 NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": "html-код с платежными системами для страницы cf_page_preview"}
	s2[18] = map[string]string{"name":"cf_bottom_html", "mysql":"text CHARACTER SET utf8 NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": "html-код с платежными системами для страницы cf_page_preview"}
	s2[19] = map[string]string{"name":"cf_ps", "mysql":"text CHARACTER SET utf8 NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": "Массива с платежными системами, которые будут выводиться на cf_page_preview"}
	s2[20] = map[string]string{"name":"auto_reload", "mysql":"int(11) NOT NULL", "sqlite":"int(11) NOT NULL","postgresql":"int NOT NULL", "comment": "Если произойдет сбой и в main_lock будет висеть запись более auto_reload секунд, тогда будет запущен сбор блоков с чистого листа"}
	s2[21] = map[string]string{"name":"commission", "mysql":"text NOT NULL", "sqlite":"text NOT NULL","postgresql":"text NOT NULL", "comment": "Максимальная комиссия, которую могут поставить ноды на данном пуле"}
	s2[22] = map[string]string{"name":"first_load_blockchain", "mysql":"enum('nodes','file')", "sqlite":"varchar(100) ","postgresql":"enum('nodes','file')", "comment": ""}
	s1["fileds"] = s2
	s1["comment"] = ""
	s["config"] = s1
	printSchema(s, dbType)



	/*
		m3["mysql"] = "text NOT NULL"
		m3["postgresql"] = "text NOT NULL"
		m2["conditions"]=m3
		m1["fileds"]=m2
		m["log_111111111"]=m1*/

	return ""
}

func typeMysql(s recmap) {
	for table_name, v := range s {
		fmt.Printf("DROP TABLE IF EXISTS %[1]s; CREATE TABLE IF NOT EXISTS %[1]s (\n", table_name)
		var tableComment string
		primaryKey := ""
		uniqKey := ""
		var tableSlice []string
		for k, v1 := range v.(recmap) {
			if k=="comment" {
				tableComment = v1.(string)
				//fmt.Println(k, v1.(string), v1)
			} else if k=="fileds" {
				//fmt.Println(k, v1)
				//i:=0
				//end:=""
				for i:=0; i<len(v1.(recmapi)); i++ {
					/*if i == len(v1.(recmap)) - 1 {
						end = ""
					} else {
						end = ","
					}*/
					tableSlice = append(tableSlice, fmt.Sprintf("`%s` %s COMMENT '%s'",  v1.(recmapi)[i].(map[string]string)["name"],  v1.(recmapi)[i].(map[string]string)["mysql"], v1.(recmapi)[i].(map[string]string)["comment"]))
					//fmt.Println(i)
					//i++
				}
			} else if k=="PRIMARY" {
				primaryKey = fmt.Sprintf("PRIMARY KEY (`%s`)", strings.Join(v1.([]string), "`,`"))
			} else if k=="UNIQ" {
				uniqKey = fmt.Sprintf("UNIQUE KEY (`%v`)", strings.Join(v1.([]string), "`,`"))
			}
		}
		if len(uniqKey) > 0 {
			tableSlice = append(tableSlice, uniqKey)
			//fmt.Printf("%s,\n", uniqKey)
		}
		if len(primaryKey) > 0 {
			tableSlice = append(tableSlice, primaryKey)
			//fmt.Printf("%s\n", primaryKey)
		}
		//fmt.Println(tableSlice)
		for i, line:= range tableSlice {
			if i == len(tableSlice) - 1{
				fmt.Printf("%s\n", line)
			} else {
				fmt.Printf("%s,\n", line)
			}
		}
		fmt.Printf(") ENGINE=MyISAM  DEFAULT CHARSET=latin1 COMMENT='%s';\n\n", tableComment)
	}
}

func typePostgresql(s recmap) {
	for table_name, v := range s {
		//var tableComment string
		primaryKey := ""
		uniqKey := ""
		AI := ""
		AI_START := "1"
		var tableSlice []string
		for k, v1 := range v.(recmap) {
			if k=="fileds" {
				for i:=0; i < len(v1.(recmapi)); i++ {
					var enumSlice []string
					dType := v1.(recmapi)[i].(map[string]string)["postgresql"]
					if ok, _ := regexp.MatchString(`enum`, dType); ok {
						//enum('normal','refund') NOT NULL DEFAULT 'normal'
						r, _ := regexp.Compile(`'([\w]+)'`)
						//fmt.Println(dType)
						for _, match := range r.FindAllStringSubmatch(dType, -1) {
							//fmt.Printf("==>%s\n",  match[1])
							if ok, _ := regexp.MatchString(`^([\w]+)$`, match[1]); ok{
								if !utils.InSliceString(match[1], enumSlice) {
									enumSlice = append(enumSlice, match[1])
								}
								//fmt.Println(match)
								//fmt.Println(enumSlice)
							}
						}
						name := v1.(recmapi)[i].(map[string]string)["name"]
						fmt.Printf("DROP TYPE IF EXISTS %s_enum_%s CASCADE;\n", table_name, name)
						fmt.Printf("CREATE TYPE %s_enum_%s AS ENUM ('%s');\n", table_name, name, strings.Join(enumSlice, "','"))
					}
				}

				for i:=0; i < len(v1.(recmapi)); i++ {
					dType := v1.(recmapi)[i].(map[string]string)["postgresql"]
					if ok, _ := regexp.MatchString(`enum`, dType); ok {
						//NOT NULL DEFAULT 'user'
						r, _ := regexp.Compile(`^enum\(.*?\)(.*)$`)
						rest := r.FindStringSubmatch(dType)
						dType = fmt.Sprintf("%s_enum_%s %s", table_name, v1.(recmapi)[i].(map[string]string)["name"], rest[1])
					}
					tableSlice = append(tableSlice, fmt.Sprintf("\"%s\" %s",  v1.(recmapi)[i].(map[string]string)["name"], dType))
				}
			} else if k == "PRIMARY" {
				primaryKey = fmt.Sprintf("ALTER TABLE ONLY \"%[1]s\" ADD CONSTRAINT %[1]s_pkey PRIMARY KEY (%[2]s);", table_name, strings.Join(v1.([]string), ","))
			} else if k == "UNIQ" {
				uniqKey = fmt.Sprintf("CREATE UNIQUE INDEX %[1]s_%[2]s ON \"%[1]s\" USING btree (%[3]s);", table_name, v1.([]string)[0], strings.Join(v1.([]string), ","))
			} else if k == "AI" {
				AI = v1.(string)
			} else if k == "AI_START" {
				AI_START = v1.(string)
			}
		}

		if len(AI) > 0 {
			fmt.Printf("DROP SEQUENCE IF EXISTS %[3]s_%[1]s_seq CASCADE;\nCREATE SEQUENCE %[3]s_%[1]s_seq START WITH %[2]s;\n", AI, AI_START, table_name)
		}

		fmt.Printf("DROP TABLE IF EXISTS \"%[1]s\"; CREATE TABLE \"%[1]s\" (\n", table_name)
		//fmt.Println(tableSlice)
		for i, line := range tableSlice {
			if i == len(tableSlice) - 1{
				fmt.Printf("%s\n", line)
			} else {
				fmt.Printf("%s,\n", line)
			}
		}
		fmt.Println(");")

		if len(uniqKey) > 0 {
			fmt.Println(uniqKey)
		}

		if len(AI) > 0 {
			fmt.Printf("ALTER SEQUENCE %[2]s_%[1]s_seq owned by %[2]s.%[1]s;\n", AI, table_name)
		}

		if len(primaryKey) > 0 {
			fmt.Println(primaryKey)
		}

		fmt.Println("\n\n")
	}
}

func typeSqlite(s recmap) {
	for table_name, v := range s {
		fmt.Printf("DROP TABLE IF EXISTS \"%[1]s\"; CREATE TABLE \"%[1]s\" (\n", table_name)
		//var tableComment string
		primaryKey := ""
		uniqKey := ""
		AI := ""
		var tableSlice []string
		for k, v1 := range v.(recmap) {
			/*if k=="comment" {
				tableComment = v1.(string)
				//fmt.Println(k, v1.(string), v1)
			} else*/ if k=="fileds" {
				//fmt.Println(k, v1)
				//i:=0
				//end:=""
				for i:=0; i<len(v1.(recmapi)); i++ {
					/*if i == len(v1.(recmap)) - 1 {
						end = ""
					} else {
						end = ","
					}*/
					tableSlice = append(tableSlice, fmt.Sprintf("\"%s\" %s",  v1.(recmapi)[i].(map[string]string)["name"],  v1.(recmapi)[i].(map[string]string)["sqlite"]))
					//fmt.Println(i)
					//i++
				}
			} else if k=="PRIMARY" {
				primaryKey = fmt.Sprintf("PRIMARY KEY (`%s`)", strings.Join(v1.([]string), "`,`"))
			} else if k=="UNIQ" {
				uniqKey = fmt.Sprintf("UNIQUE (`%v`)", strings.Join(v1.([]string), "`,`"))
			} else if k=="AI" {
				AI = v1.(string)
			}
		}
		if len(uniqKey) > 0 {
			tableSlice = append(tableSlice, uniqKey)
			//fmt.Printf("%s,\n", uniqKey)
		}
		if len(primaryKey) > 0 && len(AI) == 0 {
			tableSlice = append(tableSlice, primaryKey)
			//fmt.Printf("%s\n", primaryKey)
		}
		//fmt.Println(tableSlice)
		for i, line:= range tableSlice {
			if i == len(tableSlice) - 1{
				fmt.Printf("%s\n", line)
			} else {
				fmt.Printf("%s,\n", line)
			}
		}
		fmt.Println(");\n\n")
	}
}

func printSchema(s recmap, dbType string) {
	switch dbType {
	case "mysql":
		typeMysql(s)
	case "sqlite":
		typeSqlite(s)
	case "postgresql":
		typePostgresql(s)

	}
}
