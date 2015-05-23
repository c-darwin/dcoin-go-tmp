package main

import (
	"fmt"
	"database/sql"
	"utils"
	_ "github.com/lib/pq"
	//"encoding/binary"
	//"bytes"
	//"encoding/hex"
	//"crypto/rand"
	//"crypto/rsa"
	//"crypto/sha1"
	//"daemons"
	"strconv"
	//"errors"
	"dcparser"
	"log"
	"os"
	"github.com/alyu/configparser"
	//"strings"
	//"regexp"
	//"reflect"
)
type Config struct {
	Section struct {
		Name string
		Flag bool
	}
}
type Data struct {
	id int32
	name [16]byte
}

func TypeArray (txType string) int32 {
	x := make([]string, 67)
	// новый юзер
	x[1] = "new_user"
	x[48] = "cf_send_dc"
	for k, v := range x {
		if v == txType {
			return int32(k)
		}
	}
	return 0
	/*
	// новый майнер
	2 => 'new_miner',
	// Добавление новой обещанной суммы
	3 => 'new_promised_amount',
	4 => 'change_promised_amount',
	// голос за претендента на майнера
	5 => 'votes_miner',
	6 => 'new_forex_order',
	7 => 'del_forex_order',
	//  новый набор max_other_currencies от нода-генератора блока
	8 => 'new_max_other_currencies',
	// geolocation. Майнер изменил свои координаты
	9 => 'change_geolocation',
	// votes_promised_amount.
	10 => 'votes_promised_amount',
	// del_promised_amount. Удаление обещанной суммы
	11 => 'del_promised_amount',
	// send_dc
	12 => 'send_dc',
	13 => 'cash_request_out',
	14 => 'cash_request_in',
	// набор голосов по разным валютам
	15 => 'votes_complex',
	16 => 'change_primary_key',
	17 => 'change_node_key',
	18 => 'for_repaid_fix',
	// занесение в БД данных из первого блока
	19 => 'admin_1block',
	// админ разжаловал майнеров в юзеры
	20 => 'admin_ban_miners',
	// админ изменил variables
	21 => 'admin_variables',
	// админ обновил набор точек для проверки лиц
	22 => 'admin_spots',
	// юзер создал кредит
	23 => 'new_credit',
	// админ вернул майнерам звание "майнер"
	24 => 'admin_unban_miners',
	// админ отправил alert message
	25 => 'admin_message',
	// майнер хочет, чтобы указаные им майнеры были разжалованы в юзеры
	26 => 'abuses',
	// майнер хочет, чтобы в указанные дни ему не приходили запросы на обмен DC
	27 => 'new_holidays',
	28 => 'actualization_promised_amounts',
	29 => 'mining',
	// Голосование нода за фото нового майнера
	30 => 'votes_node_new_miner',
	// Юзер исправил проблему с отдачей фото и шлет повторный запрос на получение статуса "майнер"
	31=>'new_miner_update',
	//  новый набор max_promised_amount от нода-генератора блока
	32=>'new_max_promised_amounts',
	//  новый набор % от нода-генератора блока
	33=>'new_pct',
	// добавление новой валюты
	34=>'admin_add_currency',
	35=>'new_cf_project',
	// новая версия, которая кладется каждому в диру public
	36=>'admin_new_version',
	// после того, как новая версия протестируется, выдаем сообщение, что необходимо обновиться
	37=>'admin_new_version_alert',
	// баг репорты
	38=>'message_to_admin',
	// админ может ответить юзеру
	39=>'admin_answer',
	40=>'cf_project_data',
	// блог админа
	41=>'admin_blog',
	// майнер меняет свой хост
	42=>'change_host',
	// майнер меняет комиссию, которую он хочет получать с тр-ий
	43=>'change_commission',
	44=>'del_cf_funding',
	// запуск урезания на основе голосования. генерит нод-генератор блока
	45=>'new_reduction',
	46=>'del_cf_project',
	47=>'cf_comment',
	48=>'cf_send_dc',
	49=>'user_avatar',
	50=>'cf_project_change_category',
	51 =>'change_creditor',
	52 =>'del_credit',
	53 =>'repayment_credit',
	54 =>'change_credit_part',
	55 =>'new_admin',
	// по истечении 30 дней после поступления запроса о восстановлении утерянного ключа, админ может изменить ключ юзера
	56=>'admin_change_primary_key',
	// юзер разрешает или отменяет разрешение на смену своего ключа админом
	57=>'change_key_active',
	// юзер отменяет запрос на смену ключа
	58=>'change_key_close',
	// юзер отправляет с другого акка запрос на получение доступа к акку, ключ к которому он потерял
	59=>'change_key_request',
	// юзер решил стать арбитром или же действующий арбитр меняет комиссии
	60=>'change_arbitrator_conditions',
	// продавец меняет % и кол-во дней для новых сделок.
	61=>'change_seller_hold_back',
	// покупатель или продавец указал список арбитров, кому доверяет
	62=>'change_arbitrator_list',
	// покупатель хочет манибэк
	63=>'money_back_request',
	// магазин добровольно делает манибэк или же арбитр делает манибек
	64=>'money_back',
	// арбитр увеличивает время манибэка, чтобы успеть разобраться в ситуации
	65=>'change_money_back_time',
	// юзер меняет url центров сертификации, где хранятся его приватные ключи
	66=>'change_ca'
	);*/
}

type Parser struct {
	db *sql.DB
	txSlice []string
	txMap map[string]string
	blockData map[string]string
}

var configIni *configparser.Section

func main() {

	f, err := os.OpenFile("dclog.txt", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0777)
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	config, err := configparser.Read("config.ini")
	if err != nil {
		log.Fatal(err)
	}
	configIni, err := config.Section("main")



	txType := "NewHolidays";
	txTime := "1426283713";
	blockData := make(map[string]string)

	var txSlice []string
	// hash
	txSlice = append(txSlice, "22cb812e53e22ee539af4a1d39b4596d")
	// type
	txSlice = append(txSlice,  strconv.Itoa(int(TypeArray(txType))))
	// time
	txSlice = append(txSlice, txTime)
	// user_id
	txSlice = append(txSlice, strconv.FormatInt(1, 10));
	//start
	txSlice = append(txSlice, strconv.FormatInt(100000, 10));
	//end
	txSlice = append(txSlice, strconv.FormatInt(4545, 10));
	// sign
	txSlice = append(txSlice, "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111")

	blockData["block_id"] = strconv.FormatInt(185510, 10);
	blockData["time"] = txTime;
	blockData["user_id"] = strconv.FormatInt(1, 10);

	//fmt.Println(txSlice)

	parser := new(dcparser.Parser)
	parser.DCDB = utils.NewDbConnect(configIni)
	parser.TxSlice = txSlice;
	parser.BlockData = blockData;

	/*for i:=0; i<10000; i++ {

		x := func() {
			stmt, err := parser.DCDB.Prepare(`INSERT INTO main_lock(lock_time,script_name) VALUES($1,$2)`)
			defer stmt.Close()
			if err!=nil {
				fmt.Println(err)
			}
			_, err = stmt.Exec(11111, "testblock_generator")
			if err!=nil {
				fmt.Println(err)
			}
		}
		x()
		//stmt, _ := parser.DCDB.Prepare(`INSERT INTO main_lock(lock_time,script_name) VALUES($1,$2)`)
		//fmt.Println(err)
		//defer stmt.Close()
		//_, _ = stmt.Exec(11111, "testblock_generator")
		//fmt.Println(err)
		//_, _ = parser.DCDB.Query("INSERT INTO main_lock(lock_time,script_name) VALUES($1,$2)", 11111, "testblock_generator")
		x2 := func() {
			row, err := parser.DCDB.Query("DELETE FROM main_lock WHERE script_name='testblock_generator'")
			defer row.Close()
			if err!=nil {
				fmt.Println(err)
			}
		}
		x2()
		//parser.DCDB.DbLock()
		//parser.DCDB.DbUnlock()
		//fmt.Println(i)
	}*/
	fmt.Println()


	err = dcparser.MakeTest(parser, txType, hashesStart);
	if err != nil {
		fmt.Println("err", err)
	}
	//go daemons.Testblock_is_ready()

	//parser.Db.HashTableData("holidays", "", "")
	//HashTableData(parser.Db.DB,"holidays", "", "")
	//hashes, err := parser.Db.AllHashes()
	utils.CheckErr(err);
	//fmt.Println(hashes)
	fmt.Println()
/*
	var ptr reflect.Value
	var value reflect.Value
	//var finalMethod reflect.Value

	i := Test{Start: "start"}

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
	fmt.Println(value)
/*
	// check for method on value
	method := value.MethodByName("Finish")
	fmt.Println(method)
	// check for method on pointer
	method = ptr.MethodByName("Finish")
	fmt.Println(method)*/

}
