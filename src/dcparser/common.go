package dcparser
import (
	//"database/sql"
	_ "github.com/lib/pq"
	"fmt"
	"utils"
	"os"
)

type Parser struct {
	*utils.DCDB
	TxSlice []string
	TxMap map[string]string
	BlockData map[string]string
}

func (p *Parser) GetTxMap(fields []string) (map[string]string, error) {
	//fmt.Println("p.TxSlice", p.TxSlice)
	//fmt.Println("fields", fields)
	if len(p.TxSlice) != len(fields)+4 {
		return nil, fmt.Errorf("bad transaction_array %d != %d (type=%d)",  len(p.TxSlice),  len(fields)+4, p.TxSlice[0])
	}
	TxMap := make(map[string]string)
	TxMap["hash"] = p.TxSlice[0]
	TxMap["type"] = p.TxSlice[1]
	TxMap["time"] = p.TxSlice[2]
	TxMap["user_id"] = p.TxSlice[3]
	for i, field := range fields {
		TxMap[field] = p.TxSlice[i+4]
	}
	//fmt.Println("TxMap", TxMap)
	//fmt.Println("TxMap[hash]", TxMap["hash"])
	//fmt.Println("p.TxSlice[0]", p.TxSlice[0])
	return TxMap, nil
}

func (p *Parser) GetMyUserId(userId int64) (int64, int64, string, []int64) {
	var myUserId int64
	var myPrefix string
	var myUserIds []int64
	var myBlockId int64
	collective := p.GetCommunityUsers()
	if len(collective) > 0 {// если работаем в пуле
		myUserIds = collective
		// есть ли юзер, который задействован среди юзеров нашего пула
		if utils.In_array(userId, collective) {
			myPrefix = fmt.Sprintf("%d_", userId)
			// чтобы не было проблем с change_primary_key нужно получить user_id только тогда, когда он был реально выдан
			// в будущем можно будет переделать, чтобы user_id можно было указывать всем и всегда заранее.
			// тогда при сбросе будут собираться более полные таблы my_, а не только те, что заполнятся в change_primary_key
			err := p.QueryRow("SELECT user_id FROM "+myPrefix+"my_table").Scan(&myUserId)
			utils.CheckErr(err)
		}
	} else {
		err := p.QueryRow("SELECT user_id FROM my_table").Scan(&myUserId)
		utils.CheckErr(err)
		myUserIds = append(myUserIds, myUserId)
	}
	err := p.QueryRow("SELECT my_block_id FROM config").Scan(&myBlockId)
	utils.CheckErr(err)
	return myUserId, myBlockId, myPrefix, myUserIds
}

func MakeTest(parser *Parser, txType string, hashesStart map[string]string) error {
	//fmt.Println("dcparser."+txType+"Init")
	err := utils.CallMethod(parser, txType+"Init")
	//fmt.Println(err)

	if i, ok := err.(error); ok {
		fmt.Println(err.(error), i)
		return err.(error)
	}

	if len(os.Args)==1 {
		err = utils.CallMethod(parser, txType)
		if i, ok := err.(error); ok {
			fmt.Println(err.(error), i)
			return err.(error)
		}

		//fmt.Println("-------------------")
		// узнаем, какие таблицы были затронуты в результате выполнения основного метода
		hashesMiddle, err := parser.AllHashes()
		if err != nil {
			return utils.ErrInfo(err)
		}
		var tables []string
		for table, hash := range hashesMiddle {
			if hash!=hashesStart[table] {
				tables = append(tables, table)
			}
		}
		fmt.Println(tables)

		// rollback
		err0 := utils.CallMethod(parser, txType+"Rollback")
		if i, ok := err0.(error); ok {
			fmt.Println(err0.(error), i)
			return err0.(error)
		}

		// сраниим хэши, которые были до начала и те, что получились после роллбэка
		hashesEnd, err := parser.AllHashes()
		if err != nil {
			return utils.ErrInfo(err)
		}
		for table, hash := range hashesEnd {
			if hash!=hashesStart[table] {
				fmt.Println("ERROR in table ", table)
			}
		}

	} else if os.Args[1] == "w" {
		err = utils.CallMethod(parser, txType)
		if i, ok := err.(error); ok {
			fmt.Println(err.(error), i)
			return err.(error)
		}
	} else if os.Args[1] == "r" {
		err = utils.CallMethod(parser, txType+"Rollback")
		if i, ok := err.(error); ok {
			fmt.Println(err.(error), i)
			return err.(error)
		}
	}
	return nil
}

// откатываем ID на кол-во затронутых строк
func (p *Parser) rollbackAI(table string, num int64) (error) {
	//fmt.Println("table", table)
	current, err := p.Single("SELECT id FROM "+table+" ORDER BY id DESC LIMIT 1", )
	if err != nil {
		return utils.ErrInfo(err)
	}

	pg_get_serial_sequence, err := p.Single("SELECT pg_get_serial_sequence('"+table+"', 'id')")
	if err != nil {
		return utils.ErrInfo(err)
	}

	_, err = p.ExecSql("ALTER SEQUENCE "+pg_get_serial_sequence+" RESTART WITH "+utils.Int64ToStr(utils.StrToInt64(current)+num))
	if err != nil {
		return utils.ErrInfo(err)
	}
	return err
}
