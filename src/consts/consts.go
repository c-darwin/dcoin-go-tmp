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
const BLOCKCHAIN_URL = "http://dcoin.me/blockchain"

var LangMap = map[string]int{"en":1, "ru":42}

var MyTables = []string {"my_admin_messages","my_cash_requests","my_comments","my_commission","my_complex_votes","my_dc_transactions","my_holidays","my_keys","my_new_users","my_node_keys","my_notifications","my_promised_amount","my_table","my_tasks","my_cf_funding"}

var ReductionDC = []int {0,10,25,50,90}

func init() {
}
