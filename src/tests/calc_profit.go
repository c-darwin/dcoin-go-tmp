package main

import (
	"fmt"
	"dcparser"
)

type calcProfitTest struct {
	amount float64
	repaid_amount float64
	time_start int64
	time_finish int64
	pct_array []map[int64]map[string]float64
	points_status_array []map[int64]string
	holidays_array [][]int64
	max_promised_amount_array []map[int64]int64
	currency_id int64
}

func main() {
	var test_data [21]*calcProfitTest
	test_data[20] = new(calcProfitTest)
	test_data[20].amount = 1500
	test_data[20].repaid_amount = 50
	test_data[20].time_start = 50
	test_data[20].time_finish = 140
	test_data[20].pct_array = []map[int64]map[string]float64{
		{0:{"user":0.0059, "miner":0.05}},
		{36:{"user":0.0088, "miner":0.08}},
		{36:{"user":0.0088, "miner":0.08}},
		{164:{"user":0.0049, "miner":0.04}},
		{223:{"user":0.0029, "miner":0.02}},
	}
	test_data[20].points_status_array = []map[int64]string{
		{0:"miner"},
		{98:"miner"},
		{101:"user"},
		{101:"user"},
		{295:"miner"},
	}
	test_data[20].holidays_array = [][]int64 {
		{0, 10},
		{10, 20},
		{30, 30},
		{40, 50},
		{66, 99},
		{233, 1999},
	}
	test_data[20].max_promised_amount_array = []map[int64]int64 {
		{0:1000},
		{63:1525},
		{64:1550},
		{139:500},
	}
	test_data[20].currency_id = 10

	p:=new(dcparser.Parser)
	profit, err:=p.CalcProfit(test_data[20].amount, test_data[20].time_start, test_data[20].time_finish, test_data[20].pct_array,  test_data[20].points_status_array,  test_data[20].holidays_array,  test_data[20].max_promised_amount_array,  test_data[20].currency_id, test_data[20].repaid_amount)
	fmt.Println(profit)
	fmt.Println(err)
}
