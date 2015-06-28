package utils

import (
	"image/draw"
	"image"
	"image/color"
	"image/png"
	"log"
	"time"
	"os"
	"encoding/base64"
	"code.google.com/p/freetype-go/freetype"
	"static"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
	"runtime"
	"path/filepath"
	"strconv"
	"errors"
	"crypto"
	"consts"
	"math"
	"crypto/sha256"
	"encoding/hex"
	"crypto/rsa"
	"reflect"
	"regexp"
	"net/http"
	"io"
	"math/big"
	"crypto/x509"
	"math/rand"
	"encoding/pem"
	"sort"
)

type BlockData struct {
	BlockId int64
	Time int64
	UserId int64
	Level int64
	Sign []byte
	Hash []byte
	HeadHash []byte
}

type prevBlockType struct {
	Hash string
	HeadHash string
	BlockId int64
	Time int64
	Level int64
}

//var db *sql.DB
//var err error

func Sleep(sec float64) {
	time.Sleep(time.Duration(sec) * 1000 * time.Millisecond)
}
type SortCfCatalog []map[string]string

func (s SortCfCatalog) Len() int {
	return len(s)
}
func (s SortCfCatalog) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s SortCfCatalog) Less(i, j int) bool {
	return s[i]["name"] < s[j]["name"]
}
func MakeCfCategories(lang map[string]string) []map[string]string {
	var cfCategory []map[string]string
	for i:=0; i < 18; i++ {
		cfCategory = append(cfCategory, map[string]string{"id": IntToStr(i), "name": lang["cf_category_"+IntToStr(i)]})
	}
	sort.Sort(SortCfCatalog(cfCategory))
	return cfCategory
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer file.Close()
	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
	}
	return image.Width, image.Height
}

type ParamType struct {
	X, Y, Width, Height int64
	Bg_path string
}

func KeyToImg(key, resultPath string, userId int64, timeFormat string, param ParamType) error {

	keyBin, _ := base64.StdEncoding.DecodeString(key)
	keyHex := append(BinToHex(keyBin), []byte("00000000")...)
	keyBin = HexToBin(keyHex)

	w, h := getImageDimension(param.Bg_path)
	fSrc, err := os.Open(param.Bg_path)
	if err != nil {
		return ErrInfo(err)
	}
	defer fSrc.Close()
	src, err := png.Decode(fSrc)
	if err != nil {
		return ErrInfo(err)
	}

	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	white := image.NewUniform(color.White)
	black := image.NewUniform(color.Black)
	draw.Draw(dst, dst.Bounds(), src, image.Point{0,0}, draw.Src)

	x :=param.X
	y := param.Y
	var color *image.Uniform
	for i:=0; i < len(keyBin); i++ {
		b:= fmt.Sprintf("%08b ", keyBin[i])
		for j:=0; j < 8; j++{
			if b[j:j+1] == "0" {
				color = black
			} else {
				color = white
			}
			dst.Set(int(x), int(y), color)
			x++
			if (x + 1 - param.X) % param.Width == 0 {
				x = param.X
				y++
			}
		}
	}
	param.Height = y - param.Y + 1

	x = 0
	// теперь пишем инфу, где искать квадрат
	info := fmt.Sprintf("%016s", strconv.FormatInt(param.X, 2)) + fmt.Sprintf("%016s", strconv.FormatInt(param.Y, 2)) + fmt.Sprintf("%016s", strconv.FormatInt(param.Width, 2)) + fmt.Sprintf("%016s", strconv.FormatInt(param.Height, 2))
	for i:=0; i < len(info); i++ {
		if info[i:i+1] == "0" {
			color = black
		} else {
			color = white
		}
		dst.Set(int(x), 0, color)
		x++
	}

	fontBytes, err := static.Asset("static/fonts/luxisr.ttf")
	if err != nil {
		return ErrInfo(err)
	}
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return ErrInfo(err)
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(font)
	c.SetFontSize(15)
	c.SetClip(dst.Bounds())
	c.SetDst(dst)
	c.SetSrc(image.Black)

	// Draw the text.
	pt := freetype.Pt(13, 300)
	_, err = c.DrawString("User ID: "+Int64ToStr(userId), pt)
	if err != nil {
		return ErrInfo(err)
	}

	t := time.Unix(time.Now().Unix(), 0)
	txTime := t.Format(timeFormat)
	pt = freetype.Pt(300, 300)
	_, err = c.DrawString(txTime, pt)
	if err != nil {
		return ErrInfo(err)
	}

	fDst, err := os.Create(resultPath)
	if err != nil {
		return ErrInfo(err)
	}
	defer fDst.Close()
	err = png.Encode(fDst, dst)
	if err != nil {
		return ErrInfo(err)
	}
	return nil
}

func  ParseBlockHeader(binaryBlock *[]byte) (*BlockData) {
	result := new(BlockData)
	// распарсим заголовок блока
	/*
	Заголовок (от 143 до 527 байт )
	TYPE (0-блок, 1-тр-я)        1
	BLOCK_ID   				       4
	TIME       					       4
	USER_ID                         5
	LEVEL                              1
	SIGN                               от 128 до 512 байт. Подпись от TYPE, BLOCK_ID, PREV_BLOCK_HASH, TIME, USER_ID, LEVEL, MRKL_ROOT
	Далее - тело блока (Тр-ии)
	*/
	result.BlockId = BinToDecBytesShift(binaryBlock, 4)
	result.Time = BinToDecBytesShift(binaryBlock, 4)
	result.UserId = BinToDecBytesShift(binaryBlock, 5)
	result.Level = BinToDecBytesShift(binaryBlock, 1)
	signSize := DecodeLength(binaryBlock)
	result.Sign = BytesShift(binaryBlock, signSize)
	return result
}

/*
func Round(f float64, places int) (float64) {
	if places==0 {
		return math.Floor(f + .5)
	} else {
		shift := math.Pow(10, float64(places))
		return math.Floor((f * shift)+.5) / shift;
	}
}
*/


// ищем ближайшее время в $points_status_array или $max_promised_amount_array
// $type - status для $points_status_array / amount - для $max_promised_amount_array
func  findMinPointsStatus(needTime int64, pointsStatusArray []map[int64]string, pType string) ([]map[string]string, []map[int64]string) {
	var findTime []int64
	newPointsStatusArray := pointsStatusArray
	var timeStatusArr []map[string]string
BR:
	for i:=0; i<len(pointsStatusArray); i++ {
		for time, _ := range pointsStatusArray[i] {
			if time > needTime {
				break BR
			}
			findTime = append(findTime, time)
			start:=i+1
			if i+1 > len(pointsStatusArray) {
				start=len(pointsStatusArray)
			}
			newPointsStatusArray = pointsStatusArray[start:]
		}
	}
	if len(findTime) > 0 {
		for i:=0; i<len(findTime); i++ {
			for _, status := range pointsStatusArray[i] {
				timeStatusArr = append(timeStatusArr, map[string]string{"time" : Int64ToStr(findTime[i]), pType : status})
			}
		}
	}
	return timeStatusArr, newPointsStatusArray
}

func findMinPct (needTime int64, pctArray []map[int64]map[string]float64, status string) float64 {
	var findTime int64 = -1
	var pct float64 = 0
BR:
	for i:=0; i<len(pctArray); i++ {
		for time, _ := range pctArray[i] {
			if time > needTime {
				break BR
			}
			findTime = int64(i)
		}
	}
	if findTime >=0 {
		for _, arr := range pctArray[findTime] {
			pct = arr[status]
		}
	}
	return pct
}


func findMinPct1 (needTime int64, pctArray []map[int64]float64) float64 {
	var findTime int64 = -1
	var pct float64 = 0
BR:
	for i:=0; i<len(pctArray); i++ {
		for time, _ := range pctArray[i] {
			if time > needTime {
				break BR
			}
			findTime = int64(i)
		}
	}
	if findTime >=0 {
		for _, pct0 := range pctArray[findTime] {
			pct = pct0
		}
	}
	return pct
}

func  getMaxPromisedAmountCalcProfit(amount, repaidAmount, maxPromisedAmount float64, currencyId int64) float64 {
	// для WOC $repaid_amount всегда = 0, т.к. cash_request на WOC послать невозможно
	// если наша сумма больше, чем максимально допустимая ($find_min_array[$i]['amount'])
	var result float64
	if (amount+repaidAmount > maxPromisedAmount) {
		result = maxPromisedAmount-repaidAmount;
	} else if (amount < maxPromisedAmount && currencyId==1) { // для WOC разрешено брать maxPromisedAmount вместо promisedAmount, если promisedAmount < maxPromisedAmount
		result = maxPromisedAmount
	} else {
		result = amount;
	}
	return result
}

type resultArrType struct {
	num_sec int64
	pct float64
	amount float64
}

type pctAmount struct {
	pct float64
	amount float64
}


func CalcProfit_24946(amount float64, timeStart, timeFinish int64, pctArray []map[int64]map[string]float64, pointsStatusArray []map[int64]string, holidaysArray [][]int64, maxPromisedAmountArray []map[int64]string, currencyId int64, repaidAmount float64) (float64, error) {

	log.Println("amount", amount)
	log.Println("timeStart", timeStart)
	log.Println("timeFinish", timeFinish)
	log.Println(pctArray)
	log.Println(pointsStatusArray)
	log.Println(holidaysArray)
	log.Println(maxPromisedAmountArray)
	log.Println("currencyId", currencyId)
	log.Println("repaidAmount", repaidAmount)


	var lastStatus string = ""
	var findMinArray []map[string]string
	var newArr []map[int64]float64
	var statusPctArray_ map[string]float64
	// нужно получить массив вида time=>pct, совместив $pct_array и $points_status_array

	findTime := func(key int64, arr []map[int64]float64) bool {
		for i:=0; i< len(arr); i++ {
			if arr[i][key]!=0 {
				return true
			}
		}
		return false
	}

	log.Println("pctArray", pctArray)
	for i:=0; i < len(pctArray); i++ {
		for time, statusPctArray := range pctArray[i] {
			log.Println("i=", i, "pctArray[i]=", pctArray[i])
			findMinArray, pointsStatusArray = findMinPointsStatus(time, pointsStatusArray, "status")
			//log.Println("i", i)
			log.Println("time", time)
			log.Println("findMinArray", findMinArray)
			log.Println("pointsStatusArray", pointsStatusArray)
			for j := 0; j < len(findMinArray); j++ {
				if StrToInt64(findMinArray[j]["time"]) < time {
					findMinPct := findMinPct_24946(StrToInt64(findMinArray[j]["time"]), pctArray, findMinArray[j]["status"]);
					if !findTime(StrToInt64(findMinArray[j]["time"]), newArr) {
						newArr = append(newArr, map[int64]float64{StrToInt64(findMinArray[j]["time"]) : findMinPct})
						log.Println("findMinPct", findMinPct)
					}
					lastStatus = findMinArray[j]["status"];
				}
			}
			if len(findMinArray) == 0 && len(lastStatus) == 0 {
				findMinArray = append(findMinArray, map[string]string{"status": "user"})
			} else if len(findMinArray) == 0 && len(lastStatus) != 0 { // есть проценты, но кончились points_status
				findMinArray = append(findMinArray, map[string]string{"status": "miner"})
			}
			if !findTime(time, newArr) {
				newArr = append(newArr, map[int64]float64{time : statusPctArray[findMinArray[len(findMinArray)-1]["status"]]})
			}
			statusPctArray_ = statusPctArray;
		}
	}

	// если в points больше чем в pct
	if len(pointsStatusArray)>0 {
		for i:=0; i < len(pointsStatusArray); i++ {
			for time, status := range pointsStatusArray[i] {
				if !findTime(time, newArr) {
					newArr = append(newArr, map[int64]float64{time : statusPctArray_[status]})
				}
			}
		}
	}

	log.Println("newArr", newArr)


	// newArr - массив, где ключи - это время из pct и points_status, а значения - проценты.

	// $max_promised_amount_array + $pct_array
	/*
	 * в $pct_array сейчас
			[1394308000] =>  0,05
			[1394308100] =>  0,1

		после обработки станет

			[1394308000] => Array
				(
					[pct] => 0,05
					[amount] => 1000
				)
			[1394308005] => Array
				(
					[pct] => 0,05
					[amount] => 100
				)
			[1394308100] => Array
				(
					[pct] => 0,1
					[amount] => 100
				)

	 * */

	findTime2 := func(key int64, arr []map[int64]pctAmount) bool {
		for i:=0; i< len(arr); i++ {
			if arr[i][key].pct!=0 {
				return true
			}
		}
		return false
	}

	var newArr2 []map[int64]pctAmount
	var lastAmount float64
	var amount_ float64
	var pct_ float64
	if len(maxPromisedAmountArray)==0{
		lastAmount = amount
	}

	// нужно получить массив вида time=>pct, совместив newArr и $max_promised_amount_array
	for i:=0; i < len(newArr); i++ {
		log.Println("i ", i)
		for time, pct := range newArr[i] {
			findMinArray, maxPromisedAmountArray = findMinPointsStatus(time, maxPromisedAmountArray, "amount")
			for j:=0; j < len(findMinArray); j++ {
				if amount+repaidAmount > StrToFloat64(findMinArray[j]["amount"]) {
					amount_ = StrToFloat64(findMinArray[j]["amount"]) - repaidAmount
				} else if amount < StrToFloat64(findMinArray[j]["amount"]) && currencyId==1 {
					amount_ = StrToFloat64(findMinArray[j]["amount"])
				} else {
					amount_ = amount
				}
				if StrToInt64(findMinArray[j]["time"]) <= time {
					minPct := findMinPct1_24946(StrToInt64(findMinArray[j]["time"]), newArr);
					if !findTime2(StrToInt64(findMinArray[j]["time"]), newArr2) {
						newArr2 = append(newArr2, map[int64]pctAmount{StrToInt64(findMinArray[j]["time"]):{pct:minPct, amount:amount_}})
					}
					lastAmount = amount_;

				}
			}
			if !findTime2(time, newArr2) {
				newArr2 = append(newArr2, map[int64]pctAmount{time:{pct:pct, amount:lastAmount}})
			}
			pct_ = pct
		}
	}


	log.Println("newArr2", newArr2)

	if !findTime2(timeFinish, newArr2) {
		newArr2 = append(newArr2, map[int64]pctAmount{timeFinish:{pct:pct_, amount:0}})
	}

	var workTime, oldTime int64
	var resultArr []resultArrType
	var oldPctAndAmount pctAmount
	var startHolidays bool
	var finishHolidaysElement int64
	//START:
	for i:=0; i < len(newArr2); i++ {

		for time, pctAndAmount := range newArr2[i] {

			log.Println("pctAndAmount", pctAndAmount)

			if (time > timeStart) {
				workTime = time
				for j := 0; j < len(holidaysArray); j++ {

					if holidaysArray[j][1] <= oldTime {
						continue
					}

					log.Println("holidaysArray[j]", holidaysArray[j])

					// полные каникулы в промежутке между time и old_time
					if holidaysArray[j][0]!=-1 && workTime >= holidaysArray[j][0] && holidaysArray[j][1]!=-1 && workTime >= holidaysArray[j][1] {
						time = holidaysArray[j][0];
						holidaysArray[j][0] = -1
						resultArr = append(resultArr, resultArrType{num_sec : (time-oldTime), pct : oldPctAndAmount.pct, amount : oldPctAndAmount.amount})
						log.Println("resultArr append")
						oldTime = holidaysArray[j][1];
						holidaysArray[j][1] = -1
					}
					if holidaysArray[j][0]!=-1 && workTime >= holidaysArray[j][0] {
						startHolidays = true; // есть начало каникул, но есть ли конец?
						finishHolidaysElement = holidaysArray[j][1]; // для записи в лог
						time = holidaysArray[j][0];
						if time < timeStart {
							time = timeStart
						}
						holidaysArray[j][0] = -1
					} else if holidaysArray[j][1]!=-1 && workTime < holidaysArray[j][1] && holidaysArray[j][0]==-1 {
						// конец каникул заканчивается после $work_time
						time = oldTime
						continue
					} else if holidaysArray[j][1]!=-1 && workTime >= holidaysArray[j][1] {
						oldTime = holidaysArray[j][1]
						holidaysArray[j][1] = -1
						startHolidays = false; // конец каникул есть
					} else if j == len(holidaysArray)-1 && !startHolidays {
						// если это последний полный внутрений холидей, то time должен быть равен текущему workTime
						time = workTime
					}
				}
				if (time > timeFinish) {
					time = timeFinish
				}
				resultArr = append(resultArr, resultArrType{num_sec : (time-oldTime), pct : oldPctAndAmount.pct, amount : oldPctAndAmount.amount})
				log.Println("new", (time-oldTime))
				oldTime = time
			} else {
				oldTime = timeStart
			}
			oldPctAndAmount = pctAndAmount
		}
	}

	log.Println("resultArr", resultArr)

	if (startHolidays && finishHolidaysElement>0) {
		log.Println("finishHolidaysElement:", finishHolidaysElement)
	}

	// время в процентах меньше, чем нужное нам конечное время
	if (oldTime < timeFinish && !startHolidays) {
		// просто берем последний процент и добиваем его до нужного $time_finish
		sec := timeFinish - oldTime;
		resultArr = append(resultArr, resultArrType{num_sec : sec, pct : oldPctAndAmount.pct, amount : oldPctAndAmount.amount})
	}


	var profit, amountAndProfit float64
	for i:=0; i < len(resultArr); i++ {
		pct := 1+resultArr[i].pct
		num := resultArr[i].num_sec
		amountAndProfit = profit +resultArr[i].amount
		//$profit = ( floor( round( $amount_and_profit*pow($pct, $num), 3)*100 ) / 100 ) - $new[$i]['amount'];
		// из-за того, что в front был подсчет без обновления points, а в рабочем методе уже с обновлением points, выходило, что в рабочем методе было больше мелких временных промежуток, и получалось profit <0.01, из-за этого было расхождение в front и попадание минуса в БД
		profit = amountAndProfit*math.Pow(pct, float64(num)) - resultArr[i].amount
	}
	log.Println("profit", profit)

	return profit, nil
}


// только для блоков до 24946
func findMinPct1_24946(needTime int64, pctArray []map[int64]float64) float64 {
	var findTime int64 = 0
	var pct float64 = 0
BR:
	for i:=0; i<len(pctArray); i++ {
		for time, _ := range pctArray[i] {
			if time > needTime {
				break BR
			}
			findTime = int64(i)
		}
	}
	if findTime >=0 {
		for _, pct0 := range pctArray[findTime] {
			pct = pct0
		}
	}
	return pct
}


// только для блоков до 24946
func findMinPct_24946 (needTime int64, pctArray []map[int64]map[string]float64, status string) float64 {
	var findTime int64 = 0
	var pct float64 = 0
	log.Println("pctArray findMinPct_24946", pctArray)
BR:
	for i:=0; i<len(pctArray); i++ {
		for time, _ := range pctArray[i] {
			log.Println(time, ">", needTime, "?")
			if time > needTime {
				log.Println("break")
				break BR
			}
			findTime = int64(i)
		}
	}
	log.Println("findTime", findTime)
	if findTime >0 {
		for _, arr := range pctArray[findTime] {
			pct = arr[status]
		}
	}
	return pct
}
func CalcProfit(amount float64, timeStart, timeFinish int64, pctArray []map[int64]map[string]float64, pointsStatusArray []map[int64]string, holidaysArray [][]int64, maxPromisedAmountArray []map[int64]string, currencyId int64, repaidAmount float64) (float64, error) {

	log.Println("CalcProfit")
	log.Println("amount", amount)
	log.Println("timeStart", timeStart)
	log.Println("timeFinish", timeFinish)
	log.Println(pctArray)
	log.Println(pointsStatusArray)
	log.Println(holidaysArray)
	log.Println(maxPromisedAmountArray)
	log.Println("currencyId", currencyId)
	log.Println("repaidAmount", repaidAmount)

	if timeStart >= timeFinish {
		return 0, nil
	}

	// для WOC майнинг останавливается только если майнера забанил админ, каникулы на WOC не действуют
	if currencyId == 1 {
		holidaysArray = nil
	}

	/* $max_promised_amount_array имеет дефолтные значения от времени = 0
	 * $pct_array имеет дефолтные значения 0% для user/miner от времени = 0
	 * в $points_status_array крайний элемент массива всегда будет относиться к текущим 30-и дням т.к. перед calc_profit всегда идет вызов points_update
	 * */

	var lastStatus string = ""
	var findMinArray []map[string]string
	var newArr []map[int64]float64
	var statusPctArray_ map[string]float64
	// нужно получить массив вида time=>pct, совместив $pct_array и $points_status_array

	findTime := func(key int64, arr []map[int64]float64) bool {
		for i:=0; i< len(arr); i++ {
			if arr[i][key]!=0 {
				return true
			}
		}
		return false
	}

	log.Println("pctArray", pctArray)
	for i:=0; i < len(pctArray); i++ {
		for time, statusPctArray := range pctArray[i] {
			log.Println("i=", i, "pctArray[i]=", pctArray[i])
			findMinArray, pointsStatusArray = findMinPointsStatus(time, pointsStatusArray, "status")
			//log.Println("i", i)
			log.Println("time", time)
			log.Println("findMinArray", findMinArray)
			log.Println("pointsStatusArray", pointsStatusArray)
			for j := 0; j < len(findMinArray); j++ {
				if StrToInt64(findMinArray[j]["time"]) <= time {
					findMinPct := findMinPct(StrToInt64(findMinArray[j]["time"]), pctArray, findMinArray[j]["status"]);
					if !findTime(StrToInt64(findMinArray[j]["time"]), newArr) {
						newArr = append(newArr, map[int64]float64{StrToInt64(findMinArray[j]["time"]) : findMinPct})
						log.Println("findMinPct", findMinPct)
					}
					lastStatus = findMinArray[j]["status"];
				}
			}
			if len(findMinArray) == 0 && len(lastStatus) == 0 {
				findMinArray = append(findMinArray, map[string]string{"status": "user"})
			} else if len(findMinArray) == 0 && len(lastStatus) != 0 { // есть проценты, но кончились points_status
				findMinArray = append(findMinArray, map[string]string{"status": lastStatus})
			}
			if !findTime(time, newArr) {
				newArr = append(newArr, map[int64]float64{time : statusPctArray[findMinArray[len(findMinArray)-1]["status"]]})
			}
			statusPctArray_ = statusPctArray;
		}
	}

	// если в points больше чем в pct
	if len(pointsStatusArray)>0 {
		for i:=0; i < len(pointsStatusArray); i++ {
			for time, status := range pointsStatusArray[i] {
				if !findTime(time, newArr) {
					newArr = append(newArr, map[int64]float64{time : statusPctArray_[status]})
				}
			}
		}
	}

	log.Println("newArr", newArr)


	// newArr - массив, где ключи - это время из pct и points_status, а значения - проценты.

	// $max_promised_amount_array + $pct_array
	/*
	 * в $pct_array сейчас
			[1394308000] =>  0,05
			[1394308100] =>  0,1

		после обработки станет

			[1394308000] => Array
				(
					[pct] => 0,05
					[amount] => 1000
				)
			[1394308005] => Array
				(
					[pct] => 0,05
					[amount] => 100
				)
			[1394308100] => Array
				(
					[pct] => 0,1
					[amount] => 100
				)

	 * */

	findTime2 := func(key int64, arr []map[int64]pctAmount) bool {
		for i:=0; i< len(arr); i++ {
			if arr[i][key].pct!=0 {
				return true
			}
		}
		return false
	}

	var newArr2 []map[int64]pctAmount
	var lastAmount float64
	var amount_ float64
	var pct_ float64
	if len(maxPromisedAmountArray)==0{
		lastAmount = amount
	}

	log.Println("newArr201", newArr)

	// нужно получить массив вида time=>pct, совместив newArr и $max_promised_amount_array
	for i:=0; i < len(newArr); i++ {
		log.Println("i ", i)
		for time, pct := range newArr[i] {
			findMinArray, maxPromisedAmountArray = findMinPointsStatus(time, maxPromisedAmountArray, "amount")
			for j:=0; j < len(findMinArray); j++ {
				amount_ = getMaxPromisedAmountCalcProfit(amount, repaidAmount, StrToFloat64(findMinArray[j]["amount"]), currencyId)
				if StrToInt64(findMinArray[j]["time"]) <= time {
					minPct := findMinPct1(StrToInt64(findMinArray[j]["time"]), newArr);
					if !findTime2(StrToInt64(findMinArray[j]["time"]), newArr2) {
						newArr2 = append(newArr2, map[int64]pctAmount{StrToInt64(findMinArray[j]["time"]):{pct:minPct, amount:amount_}})
					}
					lastAmount = amount_;

				}
			}
			if !findTime2(time, newArr2) {
				newArr2 = append(newArr2, map[int64]pctAmount{time:{pct:pct, amount:lastAmount}})
				log.Println("findTime2", time, pct)
			}
			pct_ = pct;
		}
	}

	log.Println("newArr21", newArr2)

	// если в max_promised_amount больше чем в pct
	if len(maxPromisedAmountArray) > 0 {
		log.Println("maxPromisedAmountArray", maxPromisedAmountArray)

		for i:=0; i<len(maxPromisedAmountArray); i++ {
			for time, maxPromisedAmount := range maxPromisedAmountArray[i] {
				MaxPromisedAmountCalcProfit := getMaxPromisedAmountCalcProfit(amount, repaidAmount, StrToFloat64(maxPromisedAmount), currencyId);
				amount_ = MaxPromisedAmountCalcProfit
				if !findTime2(time, newArr2) {
					newArr2 = append(newArr2, map[int64]pctAmount{time:{pct:pct_, amount:MaxPromisedAmountCalcProfit}})
				}
			}
		}
	}


	maxTimeInNewArr2 := func(newArr2 []map[int64]pctAmount) int64 {
		var max int64
		for i:=0; i < len(newArr2); i++ {
			for time, _ := range newArr2[i] {
				if time > max {
					max = time
				}
			}
		}
		return max
	}

	if maxTimeInNewArr2(newArr2) < timeFinish {
		// добавим сразу время окончания
		//newArr2[timeFinish] = pct;
		if !findTime2(timeFinish, newArr2) {
			newArr2 = append(newArr2, map[int64]pctAmount{timeFinish:{pct:pct_, amount:0}})
		}
	}

	var workTime, oldTime int64
	var resultArr []resultArrType
	var oldPctAndAmount pctAmount
	var startHolidays bool
	var finishHolidaysElement int64
	log.Println("newArr2", newArr2)
START:
	for i:=0; i < len(newArr2); i++ {

		for time, pctAndAmount := range newArr2[i] {

			log.Println(time, timeFinish)
			log.Println("pctAndAmount", pctAndAmount)
			if (time > timeFinish) {
				log.Println("continue START", time, timeFinish)
				continue START
			}
			if (time > timeStart) {
				workTime = time
				for j := 0; j < len(holidaysArray); j++ {

					if holidaysArray[j][1] <= oldTime {
						continue
					}

					log.Println("holidaysArray[j]", holidaysArray[j])

					// полные каникулы в промежутке между time и old_time
					if holidaysArray[j][0]!=-1 && oldTime <= holidaysArray[j][0] && holidaysArray[j][1]!=-1 && workTime >= holidaysArray[j][1] {
						time = holidaysArray[j][0];
						holidaysArray[j][0] = -1
						resultArr = append(resultArr, resultArrType{num_sec : (time-oldTime), pct : oldPctAndAmount.pct, amount : oldPctAndAmount.amount})
						log.Println("resultArr append")
						oldTime = holidaysArray[j][1];
						holidaysArray[j][1] = -1
					}
					if holidaysArray[j][0]!=-1 && workTime >= holidaysArray[j][0] {
						startHolidays = true; // есть начало каникул, но есть ли конец?
						finishHolidaysElement = holidaysArray[j][1]; // для записи в лог
						time = holidaysArray[j][0];
						if time < timeStart {
							time = timeStart
						}
						holidaysArray[j][0] = -1
					} else if holidaysArray[j][1]!=-1 && workTime < holidaysArray[j][1] && holidaysArray[j][0]==-1 {
						// конец каникул заканчивается после $work_time
						time = oldTime
						continue
					} else if holidaysArray[j][1]!=-1 && workTime >= holidaysArray[j][1] {
						oldTime = holidaysArray[j][1]
						holidaysArray[j][1] = -1
						startHolidays = false; // конец каникул есть
					} else if j == len(holidaysArray)-1 && !startHolidays {
						// если это последний полный внутрений холидей, то time должен быть равен текущему workTime
						time = workTime
					}
				}
				if (time > timeFinish) {
					time = timeFinish
				}
				resultArr = append(resultArr, resultArrType{num_sec : (time-oldTime), pct : oldPctAndAmount.pct, amount : oldPctAndAmount.amount})
				log.Println("new", (time-oldTime))
				oldTime = time
			} else {
				oldTime = timeStart
			}
			oldPctAndAmount = pctAndAmount
			log.Println("oldPctAndAmount", oldPctAndAmount)
		}
	}
	log.Println("oldTime", oldTime)
	log.Println("timeFinish", timeFinish)

	log.Println("resultArr", resultArr)

	if (startHolidays && finishHolidaysElement>0) {
		log.Println("finishHolidaysElement:", finishHolidaysElement)
	}

	// время в процентах меньше, чем нужное нам конечное время
	if (oldTime < timeFinish && !startHolidays) {
		log.Println("oldTime < timeFinish")
		// просто берем последний процент и добиваем его до нужного $time_finish
		sec := timeFinish - oldTime;
		resultArr = append(resultArr, resultArrType{num_sec : sec, pct : oldPctAndAmount.pct, amount : oldPctAndAmount.amount})
	}

	log.Println("resultArr", resultArr)

	var profit, amountAndProfit float64
	for i:=0; i < len(resultArr); i++ {
		pct := 1+resultArr[i].pct
		num := resultArr[i].num_sec
		amountAndProfit = profit +resultArr[i].amount
		//$profit = ( floor( round( $amount_and_profit*pow($pct, $num), 3)*100 ) / 100 ) - $new[$i]['amount'];
		// из-за того, что в front был подсчет без обновления points, а в рабочем методе уже с обновлением points, выходило, что в рабочем методе было больше мелких временных промежуток, и получалось profit <0.01, из-за этого было расхождение в front и попадание минуса в БД
		profit = amountAndProfit*math.Pow(pct, float64(num)) - resultArr[i].amount
	}

	log.Println("total profit w/o amount:", profit)

	return profit, nil
}

func round(num float64) int64 {
	log.Println("num", num)
	//num += ROUND_FIX
	//	return int(StrToFloat64(Float64ToStr(num)) + math.Copysign(0.5, num))
	log.Println("num", num)
	return int64(num + math.Copysign(0.5, num))
}

func Round(num float64, precision int) float64 {
	num += consts.ROUND_FIX
	log.Println("num", num)
	//num = StrToFloat64(Float64ToStr(num))
	log.Println("precision", precision)
	log.Println("float64(precision)", float64(precision))
	output := math.Pow(10, float64(precision))
	log.Println("output", output)
	return float64(round(num * output)) / output
}

func CheckDaemonRestart(name string) bool {
	return false
}

func RandSlice(min, max, count int64) []string {
	var result []string
	for i:=0; i<int(count); i++ {
		result = append(result, IntToStr(RandInt(int(min), int(max))))
	}
	return result
}

func RandInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func PpLenght(p1, p2 [2]int) float64 {
	return math.Sqrt(math.Pow(float64(p1[0]-p2[0]), 2) + math.Pow(float64(p1[1]-p2[1]), 2))
}

func CheckInputData(data_ interface{}, dataType string) bool {
	return CheckInputData_(data_, dataType, "")
}
// функция проверки входящих данных
func CheckInputData_(data_ interface{}, dataType string, info string) bool {
	var data string
	switch data_.(type) {
	case int:
		data = IntToStr(data_.(int))
	case int64:
		data = Int64ToStr(data_.(int64))
	case float64:
		data = Float64ToStr(data_.(float64))
	case string:
		data = data_.(string)
	case []byte:
		data = string(data_.([]byte))
	}
	log.Println("CheckInputData_:"+data)
	switch dataType {
	case "arbitration_trust_list":
		if ok, _ := regexp.MatchString(`^\[[0-9]{1,10}(,[0-9]{1,10}){0,100}\]$`, data); ok{
			return true
		}
	case "votes_comment", "cf_comment":
		if ok, _ := regexp.MatchString(`^[\pL0-9\,\s\.\-\:\=\;\?\!\%\)\(\@\/\n\r]{1,140}$`, data); ok{
			return true
		}
	case "type":
		if ok, _ := regexp.MatchString(`^[\w]+$`, data); ok{
			if StrToInt(data) <= 30 {
				return true
			}
		}
	case "referral":
		if ok, _ := regexp.MatchString(`^[0-9]{1,2}$`, data); ok{
			if StrToInt(data) <= 30 {
				return true
			}
		}
	case "currency_id":
		if ok, _ := regexp.MatchString(`^[0-9]{1,3}$`, data); ok{
			if StrToInt(data) <= 255 {
				return true
			}
		}
	case "tinyint":
		if ok, _ := regexp.MatchString(`^[0-9]{1,3}$`, data); ok{
			if StrToInt(data) <= 127 {
				return true
			}
		}
	case "smallint":
		if ok, _ := regexp.MatchString(`^[0-9]{1,5}$`, data); ok{
			if StrToInt(data) <= 65535 {
				return true
			}
		}
	case "reduction_type":
		if ok, _ := regexp.MatchString(`^(manual|promised_amount)$`, data); ok{
			if StrToInt(data) <= 30 {
				return true
			}
		}
	case "img_url":
		regex := `https?\:\/\/`; // SCHEME
		regex += `[\w-.]*\.[a-z]{2,4}`; // Host or IP
		regex += `(\:[0-9]{2,5})?`; // Port
		regex += `(\/[\w_-]+)*\/?`; // Path
		regex += `\.(png|jpg)`; // Img
		if ok, _ := regexp.MatchString(`^`+regex+`$`, data); ok{
			if len(data) < 50 {
				return true
			}
		}
	case "ca_url", "arbitrator_url":
		regex := `https?\:\/\/`; // SCHEME
		regex += `[\w-.]*\.[a-z]{2,4}`; // Host or IP
		regex += `(\:[0-9]{2,5})?`; // Port
		regex += `(\/[\w_-]+)*\/?`; // Path
		if ok, _ := regexp.MatchString(`^`+regex+`$`, data); ok{
			if len(data) <= 30 {
				return true
			}
		}
	case "credit_pct", "pct":
		if ok, _ := regexp.MatchString(`^[0-9]{1,3}(\.[0-9]{2})?$`, data); ok{
			return true
		}
	case "user_name":
		if ok, _ := regexp.MatchString(`^[\w\s]{1,30}$`, data); ok{
			return true
		}
	case "admin_currency_list":
		if ok, _ := regexp.MatchString(`^((\d{1,3}\,){0,9}\d{1,3}|ALL)$`, data); ok{
			return true
		}
	case "cf_currency_name":
		if ok, _ := regexp.MatchString(`^[A-Z0-9]{7}$`, data); ok{
			return true
		}
	case "users_ids":
		if ok, _ := regexp.MatchString(`^([0-9]{1,12},){0,1000}[0-9]{1,12}$`, data); ok{
			return true
		}
	case "version":
		if ok, _ := regexp.MatchString(`^[0-9]{1,2}\.[0-9]{1,2}\.[0-9]{1,2}([a-z]{1,2}[0-9]{1,2})?$`, data); ok{
			return true
		}
	case "soft_type":
		if ok, _ := regexp.MatchString(`^[a-z]{3,10}$`, data); ok{
			return true
		}
	case "currency_name":
		if ok, _ := regexp.MatchString(`^[A-Z]{3}$`, data); ok{
			return true
		}
	case "currency_full_name":
		if ok, _ := regexp.MatchString(`^[a-zA-Z\s]{3,50}$`, data); ok{
			return true
		}
	case "currency_commission":
		if ok, _ := regexp.MatchString(`^[0-9]{1,7}(\.[0-9]{1,2})?$`, data); ok{
			return true
		}
	case "sell_rate":
		if ok, _ := regexp.MatchString(`^[0-9]{0,10}(\.[0-9]{0,10})?$`, data); ok{
			return true
		}
	case "amount":
		if ok, _ := regexp.MatchString(`^[0-9]{0,10}(\.[0-9]{0,2})?$`, data); ok{
			return true
		}
	case "tpl_name":
		if ok, _ := regexp.MatchString("^[\\w]{1,30}$", data); ok{
			return true
		}
	case "example_spots":
		r1 := `"\d{1,2}":\["\d{1,3}","\d{1,3}",(\[("[a-z_]{1,30}",?){0,20}\]|""),"\d{1,2}","\d{1,2}"\]`
		reg := `^\{(\"(face|profile)\":\{(`+r1+`,?){1,20}\},?){2}}$`
		if ok, _ := regexp.MatchString(reg, data); ok{
			return true
		}
	case "segments":
		r1 := `"\d{1,2}":\["\d{1,2}","\d{1,2}"\]`
		face := `"face":\{(`+r1+`\,){1,20}`+r1+`\}`
		profile := `"profile":\{(`+r1+`\,){1,20}`+r1+`\}`
		reg := `^\{`+face+`,`+profile+`\}$`
		if ok, _ := regexp.MatchString(reg, data); ok{
			return true
		}
	case "tolerances":
		r1 := `"\d{1,2}":"0\.\d{1,2}"`
		face := `"face":\{(`+r1+`\,){1,50}`+r1+`\}`
		profile := `"profile":\{(`+r1+`\,){1,50}`+r1+`\}`
		reg := `^\{`+face+`,`+profile+`\}$`
		if ok, _ := regexp.MatchString(reg, data); ok{
			return true
		}
	case "compatibility":
		if ok, _ := regexp.MatchString(`^\[(\d{1,5},)*\d{1,5}\]$`, data); ok{
			return true
		}
	case "race":
		if ok, _ := regexp.MatchString("^[1-3]$", data); ok{
			return true
		}
	case "country":
		if ok, _ := regexp.MatchString("^[0-9]{1,3}$", data); ok{
			return true
		}
	case "vote", "boolean":
		if ok, _ := regexp.MatchString(`^0|1$`, data); ok{
			return true
		}
	case "coordinate":
		if ok, _ := regexp.MatchString(`^\-?[0-9]{1,3}(\.[0-9]{1,5})?$`, data); ok{
			return true
		}
	case "cf_links":
		regex := `\["https?\:\/\/(goo\.gl|bit\.ly|t\.co)\/[\w-]+",[0-9]+,[0-9]+,[0-9]+,[0-9]+\]`
		if ok, _ := regexp.MatchString(`^\[`+regex+`(\,`+regex+`)*\]$`, data); ok{
			if len(data) < 512 {
				return true
			}
		}
	case "host":
		if ok, _ := regexp.MatchString(`^https?:\/\/[0-9a-z\_\.\-\/:]{1,100}[\/]$`, data); ok{
			return true
		}
	case "coords":
		xy := `\[\d{1,3}\,\d{1,3}\]`;
		r := `^\[(`+xy+`\,){`+info+`}`+xy+`\]$`;
		if ok, _ := regexp.MatchString(r, data); ok{
			return true
		}
		fmt.Println(r)
		fmt.Println(data)
	case "lang":
		if ok, _ := regexp.MatchString("^(en|ru)$", data); ok{
			return true
		}
	case "payment_systems_ids":
		if ok, _ := regexp.MatchString("^([0-9]{1,4},){0,4}[0-9]{1,4}$", data); ok{
			return true
		}
	case "video_type":
		if ok, _ := regexp.MatchString("^(youtube|vimeo|youku|null)$", data); ok{
			return true
		}
	case "video_url_id":
		if ok, _ := regexp.MatchString("^(?i)([0-9a-z_-]{5,32}|null)$", data); ok{
			return true
		}
	case "photo_hash", "sha256":
		if ok, _ := regexp.MatchString("^[0-9a-z]{64}$", data); ok{
			return true
		}
	case "alert":
		if ok, _ := regexp.MatchString("^[\\pL0-9\\,\\s\\.\\-\\:\\=\\;\\?\\!\\%\\)\\(\\@\\/]{1,512}$", data); ok{
			return true
		}
	case "int":
		if ok, _ := regexp.MatchString("^[0-9]{1,10}$", data); ok{
			return true
		}
	case "float":
		if ok, _ := regexp.MatchString(`^[0-9]{1,5}(\.[0-9]{1,5})?$`, data); ok{
			return true
		}
	case "sleep_var":
		if ok, _ := regexp.MatchString(`^\{\"is_ready\"\:\[([0-9]{1,5},){1,100}[0-9]{1,5}\],\"generator\"\:\[([0-9]{1,5},){1,100}[0-9]{1,5}\]\}$`, data); ok{
			return true
		}
	case "int64", "bigint", "user_id":
		if ok, _ := regexp.MatchString("^[0-9]{1,15}$", data); ok{
			return true
		}
	case "level":
		if StrToInt(data) >= 0 && StrToInt(data) <= 34 {
			return true
		}
	case "comment":
		if len(data) >= 1 && len(data) <= 512 {
			return true
		}
	case "hex_sign", "hex", "public_key":
		if ok, _ := regexp.MatchString("^[0-9a-z]+$", data); ok{
			if len(data) < 2048 {
				return true
			}
		}
	}

	return false
}

func Time() int64 {
	return time.Now().Unix()
}

func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

// без проверки на ошибки т.к. тут ошибки не могут навредить
func StrToInt64(s string) int64 {
	int64, _ := strconv.ParseInt(s, 10, 64)
	return int64
}
func BytesToInt64(s []byte) int64 {
	int64, _ := strconv.ParseInt(string(s), 10, 64)
	return int64
}
func StrToUint64(s string) uint64 {
	int64, _ := strconv.ParseInt(s, 10, 64)
	return uint64(int64)
}
func StrToInt(s string) int {
	int_, _ := strconv.Atoi(s)
	return int_
}
func Float64ToStr(f float64) string {
	return strconv.FormatFloat(f,'f', 13, 64)
}
func Float64ToBytes(f float64) []byte {
	return []byte(strconv.FormatFloat(f,'f', 13, 64))
}
func Float64ToStrPct(f float64) string {
	return strconv.FormatFloat(f,'f', 2, 64)
}
func StrToFloat64(s string) float64 {
	Float64, _ := strconv.ParseFloat(s, 64)
	return Float64
}
func BytesToFloat64(s []byte) float64 {
	Float64, _ := strconv.ParseFloat(string(s), 64)
	return Float64
}
func BytesToInt(s []byte) int {
	int_, _ := strconv.Atoi(string(s))
	return int_
}
func StrToMoney(str string) float64 {
	ind:=strings.Index(str, ".")
	new:=""
	if ind!=-1 {
		end := 2
		if len(str[ind + 1 : ]) > 1 {
			end = 3
		}
		new = str[ : ind] + "." + str[ind + 1 : ind + end]
	} else {
		new = str
	}
	return StrToFloat64(new)
}
func DownloadToFile(url, file string) (int64, error) {
	out, err := os.Create(file)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	countBytes, err := io.Copy(out, resp.Body)
	if err != nil {
		return countBytes, err
	}
	return countBytes, nil
}

func CheckErr(err error) {
	if err != nil {
		panic(fmt.Sprintf("%s", err))
	}
}

func ErrInfoFmt(err string, a ...interface{}) error {
	err_ := fmt.Sprintf(err, a...)
	return fmt.Errorf("%s (%s)", err_, Caller(1))
}

func ErrInfo(err error, additionally...string) error {
	if err != nil {
		if len(additionally) > 0 {
			return fmt.Errorf("%s # %s (%s)", err, additionally, Caller(1))
		} else {
			return fmt.Errorf("%s (%s)", err, Caller(1))
		}
	}
	return err
}

func CallMethod(i interface{}, methodName string) interface{} {
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
		return finalMethod.Call([]reflect.Value{})[0].Interface()
	}

	// return or panic, method not found of either type
	return fmt.Errorf("method %s not found", methodName)
}

func Caller(steps int) string {
	name := "?"
	if pc, _, num, ok := runtime.Caller(steps + 1); ok {
		//fmt.Println(num)
		name = fmt.Sprintf("%s :  %d", filepath.Base(runtime.FuncForPC(pc).Name()), num)
	}
	return name
}

func GetEntropy(hash string) int64 {
	if len(hash)>=6 {
		result, err := strconv.ParseInt(hash[0:6], 16, 0)
		CheckErr(err)
		return result;
	} else {
		return 0;
	}
}

/**
 * Определяем, какой юзер должен генерить блок
 *
 * @param int $max_miner_id Общее число майнеров
 * @param string $ctx Энтропия
 * @return int ID майнера
 */
func GetBlockGeneratorMinerId(maxMinerId, ctx int64) int64 {
	var x, hi, lo float64
	hi = float64(ctx) / 127773;
	lo = float64(ctx % 127773);
	x = 16807 * lo - 2836 * hi;
	if (x <= 0) {
		x += 0x7fffffff;
	}
	rez := int64(x) % (maxMinerId + 1);
	if rez == 0 {
		rez = 1;
	}
	return rez;
}
/*
new
static function get_block_generator_miner_id ($max_miner_id, $ctx)
{
	$n = ceil( log($max_miner_id) / log(16) );
	$hash = $ctx;
	do {
	  $hash = hash('sha256', $hash);
	  $c = substr($hash, 0, $n);
	  $level_0_miner_id = hexdec($c);
	} while( $level_0_miner_id > $max_miner_id || !$level_0_miner_id );

	return $level_0_miner_id;
}
*/



func SliceInt64ToString(int64 []int64) []string  {
	result := make([]string, len(int64))
	for i, v := range int64 {
		result[i] = strconv.FormatInt(v, 10)
	}
	return result
}

func RemoveInt64Slice(slice []int64, pos int) {
	slice = append(slice[:pos], slice[pos+1:]...)
}

func DelUserIdFromArray(array *[]int64, userId int64) {
	for i, v := range *array {
		if v == userId {
			RemoveInt64Slice(*array, i);
		}
	}
}



func InSliceInt64(search int64, slice []int64) bool {
	for _, v := range slice {
		if v == search {
			return true
		}
	}
	return false
}

func InSliceString(search string, slice []string) bool {
	for _, v := range slice {
		if v == search {
			return true
		}
	}
	return false
}

func GetBlockGeneratorMinerIdRange(currentMinerId, maxMinerId int64) [][][]int64 {
	var end int64
	var begin int64
	minus1Ok := 0
	minusStop := false
	//var result [][][2]int64{1, 2, 3}
	result:= [][][]int64{{{currentMinerId, currentMinerId}}}
	// на верхнем уровне тот, кто генерит блок первым
	var i float64 = 1
	for {
		needUsers := math.Pow(2, i)
		//fmt.Println("needUsers", needUsers)
		if begin == 0 {
			var x int64
			if i > 1 {
				x = int64(math.Pow(2, i-1))
			}
			begin = currentMinerId + 1 + x
		} else {
			begin = end + 1
		}
		//fmt.Println("begin", begin)
		if begin == maxMinerId+1 && !minusStop && currentMinerId > 1 && begin != 2 {
			begin = 1;
			//fmt.Println("begin ", begin)
			minusStop = true
		} else {
			if begin == currentMinerId || end == currentMinerId || begin > maxMinerId {
				break;
			}
		}

		end = begin + int64(needUsers) - 1
		if end > currentMinerId && minus1Ok > 0 {
			//fmt.Println("$end > $cur_miner_id && $minus_1_ok ")
			end = currentMinerId - 1;
		}


		end_p := end;
		if end_p > maxMinerId  {
			//fmt.Println("$end_p > $max_miner_id")
			end_p = maxMinerId
		}

		if end_p == maxMinerId && end_p == currentMinerId {
			//fmt.Println("$end_p == $max_miner_id && $end_p == $cur_miner_id")
			end_p =currentMinerId - 1;
		}

		result = append(result, [][]int64{{begin, end_p}});
		var minus int64 = 0
		if end > maxMinerId && !minusStop {
			//fmt.Println("$end > $max_miner_id && !$minus_stop")
			minus = maxMinerId  - end;
			if int64(math.Abs(float64(minus)))>=currentMinerId {
				minus = - (currentMinerId - 1)
			}
			end = int64(math.Abs(float64(minus)));
			minus1Ok = 1
		}

		if minus!=0 {
			result[int(i)] = append(result[int(i)], []int64{1, int64(math.Abs(float64(minus)))});
		}
		i++;
	}
	return result
}

func FindMinerIdLevel(minersIds []int64, levelsRange [][][]int64) (int64, int64) {
	for _, minerId := range minersIds {
		for level, ranges := range levelsRange {
			if minerId >= ranges[0][0] && minerId <= ranges[0][1] {
				return minerId, int64(level);
			}
			if len(ranges) == 2 {
				if minerId >= ranges[1][0] && minerId <= ranges[1][1] {
					return minerId, int64(level)
				}
			}
		}
	}
	return 0, -1
}

func GetIsReadySleep0(level int64, data []int64) int64 {
	if int64(len(data)) > level {
		return data[level]
	} else {
		return 0
	}
}

func GetOurLevelNodes(level int64, levelsRange [][][]int64) []int64 {
	var result []int64
	if level != -1 {
		for i := levelsRange[level][0][0]; i <= levelsRange[level][0][1]; i++ {
			result = append(result, i);
		}
		if len(levelsRange[level]) == 2 {
			for i := levelsRange[level][1][0]; i <= levelsRange[level][1][1]; i++ {
				result = append(result, i);
			}
		}
	}
	return result
}
// на 0-м уровне всегда большее значение, чтобы успели набраться тр-ии
// на остальных уровнях - это время, за которое нужно успеть получить новый блок и занести его в БД
func GetGeneratorSleep(level int64, data []int64) int64 {
	var sleep int64;
	if int64(len(data)) > level {
		// суммируем время со всех уровней, которые не успели сгенерить блок до нас
		for i := 0; i <= int(level); i++ {
			sleep+=data[i];
		}
	}
	return sleep;
}

// сумма is_ready всех предыдущих уровней, которые не успели сгенерить блок
func GetIsReadySleepSum(level int64, data []int64) int64 {
	var sum int64;
	for i := 0; i < int(level); i++ {
		sum += data[i];
	}
	return sum;
}

func EncodeLengthPlusData(data []byte) []byte  {
	//fmt.Println("len(data)", len(data))
	//fmt.Printf("EncodeLength(int64(len(data)) %s\n", EncodeLength(int64(len(data))))
	return append(EncodeLength(int64(len(data))) , data...)
}

func EncodeLength(len0 int64) []byte  {
	if len0<=127 {
		if len0%2 > 0 {
			return HexToBin([]byte(fmt.Sprintf("0%x", len0)))
		}
		return HexToBin([]byte(fmt.Sprintf("%x", len0)))
	}
	temphex:= fmt.Sprintf("%x", len0)
	if len(temphex)%2 > 0 {
		temphex = "0"+temphex;
	}
	str, _ := hex.DecodeString(temphex)
	temp := string(str)
	t1 := (0x80 | len(temp))
	t1hex:= fmt.Sprintf("%x", t1)
	len_and_t1 := t1hex+temphex
	len_and_t1_bin, _ := hex.DecodeString(len_and_t1)
	//fmt.Println("len_and_t1_bin", len_and_t1_bin)
	//fmt.Printf("len_and_t1_bin %x\n", len_and_t1_bin)
	return len_and_t1_bin
}

func DecToHex(dec int64) string {
	return strconv.FormatInt(dec, 16)
}

func HexToDec(h string) int64 {
	int64, _ := strconv.ParseInt(h, 16, 0)
	return int64
}

func HexToDecBig(hex string) string {
	i := new(big.Int)
	i.SetString(hex, 16)
	return fmt.Sprintf("%d", i)
}

func DecToHexBig(hex string) string {
	i := new(big.Int)
	i.SetString(hex, 10)
	hex = fmt.Sprintf("%x", i)
	if len(hex)%2 >0{
		hex = "0"+hex
	}
	return hex
}

func Int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}
func Int64ToByte(num int64) []byte {
	return []byte(strconv.FormatInt(num, 10))
}

func IntToStr(num int) string {
	return strconv.Itoa(num)
}

func DecToBin(dec, sizeBytes int64) []byte {
	Hex := fmt.Sprintf("%0"+Int64ToStr(sizeBytes*2)+"x", dec)
	//fmt.Println("Hex", Hex)
	return HexToBin([]byte(Hex))
}
func BinToHex(bin []byte) []byte {
	return []byte(fmt.Sprintf("%x", bin))
}

func HexToBin(hex_ []byte) []byte {
	// без проверки на ошибки т.к. эта функция используется только там где валидность была проверена ранее
	var str []byte
	str, err := hex.DecodeString(string(hex_))
	if err!=nil {
		fmt.Println(err)
	}
	return str
}




func BinToDec(bin []byte) int64 {
	var a uint64
	l := len(bin)
	for i, b := range bin {
		shift := uint64((l-i-1) * 8)
		a |= uint64(b) << shift
	}
	return int64(a)
}

func BinToDecBytesShift(bin *[]byte, num int64) int64 {
	return BinToDec(BytesShift(bin, num))
}

func BytesShift(str *[]byte, index int64) []byte {
	var substr []byte
	var str_ []byte
	substr = *str
	substr = substr[0:index]
	str_ = *str
	str_ = str_[index:]
	*str = str_
	return substr
}

func InterfaceToStr(v interface {}) string {
	var str string
	switch v.(type) {
		case int:
			str = IntToStr(v.(int))
		case float64:
			str = Float64ToStr(v.(float64))
		case int64:
			str = Int64ToStr(v.(int64))
		case string:
			str = v.(string)
		case []byte:
			str = string(v.([]byte))
	}
	return str
}
func InterfaceSliceToStr(i []interface {}) []string {
	var str []string
	for _, v := range i {
		switch v.(type) {
		case int:
			str = append(str, IntToStr(v.(int)))
		case float64:
			str = append(str, Float64ToStr(v.(float64)))
		case int64:
			str = append(str, Int64ToStr(v.(int64)))
		case string:
			str = append(str, v.(string))
		case []byte:
			str = append(str, string(v.([]byte)))
		}
	}
	return str
}

func InterfaceToFloat64(i interface {}) float64 {
	var result float64
		switch i.(type) {
		case int:
			result = float64(i.(int))
		case float64:
			result = i.(float64)
		case int64:
			result = float64(i.(int64))
		case string:
			result = StrToFloat64(i.(string))
		case []byte:
			result = BytesToFloat64(i.([]byte))
		}
	return result
}

func BytesShiftReverse(str *[]byte, index_ interface{}) []byte {
	var index int64
	switch index_.(type) {
	case int:
		index = int64(index_.(int))
	case int64:
		index = index_.(int64)
	}

	var substr []byte
	var str_ []byte
	substr = *str
	substr = substr[int64(len(substr))-index:]
	//fmt.Println(substr)
	str_ = *str
	str_ = str_[0:int64(len(str_))-index]
	*str = str_
	//fmt.Println(utils.BinToHex(str_))
	return substr
}


func DecodeLength(str *[]byte) int64 {
	var str_ []byte
	str_ = *str
	if len(str_)== 0 {
		return 0
	}
	length_ := []byte(BytesShift(&str_, 1))
	*str = str_
	length := int64(length_[0])
	//fmt.Println(length)
	t1 := (length & 0x80)
	//fmt.Printf("length&0x80 %x", t1)
	if t1>0 {
		//fmt.Println("1")
		length &= 0x7F;
		//fmt.Printf("length %x\n", length)
		temp := BytesShift(&str_, length)
		*str = str_
		//fmt.Printf("temp %x\n", temp)
		temp2 := fmt.Sprintf("%08x", temp)
		//fmt.Println("temp2", temp2)
		temp3 := HexToDec(temp2)
		//fmt.Println("temp3", temp3)
		return temp3
	}
	return length
}

func SleepDiff(sleep *int64, diff int64) {
	// вычитаем уже прошедшее время
	if *sleep > diff {
		*sleep = *sleep - diff;
	} else {
		*sleep = 0;
	}
}


func RollbackTransactionsTestBlock(truncate bool){

}

func AllTxParser() {

}
func CopyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func RandSeq(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func MakeAsn1(hex_n, hex_e []byte) []byte {
	//hex_n = append([]byte("00"), hex_n...)
	n_ := []byte(HexToBin(hex_n))
	n_ = append([]byte("02"), BinToHex(EncodeLength(int64(len(HexToBin(hex_n)))))...)
	//log.Println("n_length", string(n_))
	n_ = append(n_, hex_n...)
	//log.Println("n_", string(n_))
	e_ := append([]byte("02"), BinToHex(EncodeLength(int64(len(HexToBin(hex_e)))))...)
	e_ = append(e_, hex_e...)
	//log.Println("e_", string(e_))
	length := BinToHex(EncodeLength(int64(len(HexToBin(append(n_,e_...))))))
	//log.Println("length", string(length))
	rez := append([]byte("30"), length...)
	rez = append(rez, n_...)
	rez = append(rez, e_...)
	rez = append([]byte("00"), rez...)
	//log.Println(string(rez))
	//log.Println(len(string(rez)))
	//log.Println(len(HexToBin(rez)))
	rez = append(BinToHex(EncodeLength(int64(len(HexToBin(rez))))), rez...)
	rez = append([]byte("03"), rez...)
	//log.Println(string(rez))
	rez = append([]byte("300d06092a864886f70d0101010500"), rez...)
	//log.Println(string(rez))
	rez = append(BinToHex(EncodeLength(int64(len(HexToBin(rez))))), rez...)
	//log.Println(string(rez))
	rez = append([]byte("30"), rez...)
	//log.Println(string(rez))

	return rez
	//b64:=base64.StdEncoding.EncodeToString([]byte(utils.HexToBin("30"+length+bin_enc)))
	//fmt.Println(b64)
}



func BinToRsaPubKey(publicKey []byte) (*rsa.PublicKey, error) {
	key := base64.StdEncoding.EncodeToString(publicKey)
	key = "-----BEGIN PUBLIC KEY-----\n"+key+"\n-----END PUBLIC KEY-----"
	//fmt.Printf("%x\n", publicKeys[i])
	log.Println("key", key)
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, ErrInfo(fmt.Errorf("incorrect key"))
	}
	re, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, ErrInfo(err)
	}
	pub := re.(*rsa.PublicKey)
	if err != nil {
		return nil, ErrInfo(err)
	}
	return pub, nil
}

func CheckSign(publicKeys [][]byte, forSign string, signs []byte, nodeKeyOrLogin bool ) (bool, error) {

	log.Println("forSign", forSign)
	//fmt.Println("publicKeys", publicKeys)
	var signsSlice [][]byte
	// у нода всегда 1 подпись
	if nodeKeyOrLogin {
		signsSlice = append(signsSlice, signs)
	} else {
		// в 1 signs может быть от 1 до 3-х подписей
		for {
			if len(signs)==0 {
				break
			}
			length := DecodeLength(&signs)
			//fmt.Println("length", length)
			//fmt.Printf("signs %x", signs)
			signsSlice = append(signsSlice, BytesShift(&signs, length))
		}
		if len(publicKeys) != len(signsSlice) {
			log.Println("signsSlice", signsSlice)
			log.Println("publicKeys", publicKeys)
			return false, fmt.Errorf("sign error %d!=%d", len(publicKeys), len(signsSlice) )
		}
	}

	for i:=0; i<len(publicKeys); i++ {
		/*log.Println("publicKeys[i]", string(publicKeys[i]))
		key := base64.StdEncoding.EncodeToString(publicKeys[i])
		key = "-----BEGIN PUBLIC KEY-----\n"+key+"\n-----END PUBLIC KEY-----"
		//fmt.Printf("%x\n", publicKeys[i])
		log.Println("key", key)
		block, _ := pem.Decode([]byte(key))
		if block == nil {
			return false, ErrInfo(fmt.Errorf("incorrect key"))
		}
		re, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return false, ErrInfo(err)
		}
		pub := re.(*rsa.PublicKey)
		if err != nil {
			return false,  ErrInfo(err)
		}*/
		pub, err := BinToRsaPubKey(publicKeys[i])
		if err != nil {
			return false,  ErrInfo(err)
		}
		err = rsa.VerifyPKCS1v15(pub, crypto.SHA1,  HashSha1(forSign), signsSlice[i])
		if err != nil {
			log.Println("pub", pub)
			log.Println("crypto.SHA1", crypto.SHA1)
			log.Println("HashSha1(forSign)", HashSha1(forSign))
			log.Println("HashSha1(forSign)", string(HashSha1(forSign)))
			log.Println("forSign", forSign)
			log.Printf("sign: %x\n", signsSlice[i])
			return false, ErrInfoFmt("incorrect sign:  hash = %x; forSign = %v",  HashSha1(forSign), forSign)
		}
	}
	return true, nil
}

func HashSha1(msg string) []byte {
	sh := crypto.SHA1.New()
	sh.Write([]byte(msg))
	hash := sh.Sum(nil)
	return hash
}

func Md5(msg_ interface {}) []byte {
	var msg []byte
	switch msg_.(type) {
	case string:
		msg = []byte(msg_.(string))
	case []byte:
		msg = msg_.([]byte)
	}
	sh := crypto.MD5.New()
	sh.Write(msg)
	hash := sh.Sum(nil)
	return BinToHex(hash)
}

func DSha256(data_ interface{}) []byte {
	var data []byte
	switch data_.(type) {
	case string:
		data = []byte(data_.(string))
	case []byte:
		data = data_.([]byte)
	}
	sha256_ := sha256.New()
	sha256_.Write(data)
	hashSha256:=fmt.Sprintf("%x", sha256_.Sum(nil))
	sha256_ = sha256.New()
	sha256_.Write([]byte(hashSha256))
	return []byte(fmt.Sprintf("%x", sha256_.Sum(nil)))
}
func Sha256(data_ interface{}) []byte {
	var data []byte
	switch data_.(type) {
	case string:
		data = []byte(data_.(string))
	case []byte:
		data = data_.([]byte)
	}
	sha256_ := sha256.New()
	sha256_.Write(data)
	return []byte(fmt.Sprintf("%x", sha256_.Sum(nil)))
}

func GetMrklroot(binaryData []byte, variables *Variables, first bool) ([]byte, error) {
	fmt.Println(variables)
	var mrklSlice [][]byte
	var txSize int64
	// [error] парсим после вызова функции
	if len(binaryData) > 0 {
		for {
			// чтобы исключить атаку на переполнение памяти
			if !first {
				if txSize > variables.Int64["max_tx_size"] {
					return nil, ErrInfoFmt("[error] MAX_TX_SIZE")
				}
			}
			txSize = DecodeLength(&binaryData)

			// отчекрыжим одну транзакцию от списка транзакций
			if txSize > 0 {
				transactionBinaryData := BytesShift(&binaryData, txSize)
				mrklSlice = append(mrklSlice, DSha256(transactionBinaryData))
			}

			// чтобы исключить атаку на переполнение памяти
			if !first {
				if len(mrklSlice) > int(variables.Int64["max_tx_count"]) {
					return nil, ErrInfo(fmt.Errorf("[error] MAX_TX_COUNT (%v > %v)", len(mrklSlice), variables.Int64["max_tx_count"]))
				}
			}
			if len(binaryData) == 0 {
				break
			}
		}
	} else {
		mrklSlice = append(mrklSlice, []byte("0"))
	}
	return MerkleTreeRoot(mrklSlice), nil
}

func SliceReverse(s []int64) []int64 {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func MerkleTreeRoot(dataArray [][]byte) []byte {
	result := make(map[int32][][]byte)
	for _, v := range dataArray {
		result[0] = append(result[0], DSha256(v))
	}
	var j int32
	for len(result[j]) > 1 {
		for i := 0; i < len(result[j]); i = i + 2 {
			if len(result[j]) <= (i+1) {
				if _, ok := result[j+1]; !ok {
					result[j+1] = [][]byte {result[j][i]}
				} else {
					result[j+1] = append(result[j+1], result[j][i])
				}
			} else {
				if _, ok := result[j+1]; !ok {
					result[j+1] = [][]byte {DSha256(append(result[j][i], result[j][i+1]...))}
				} else {
					result[j+1] = append(result[j+1], DSha256([]byte(append(result[j][i], result[j][i+1]...))))
				}
			}
		}
		j++
	}
	result_ := result[int32(len(result)-1)];
	return []byte(result_[0])
}

func DbConnect(configIni map[string]string) *DCDB {
	DB_CONNECT:
		db, err := NewDbConnect(configIni)
		if err != nil {
			Sleep(1)
			goto DB_CONNECT
		}
		return db
}

func DbClose(c *DCDB) {
	err := c.Close()
	if err != nil {
		log.Print(err)
	}
}

func GetAllMaxPromisedAmount() []int64 {
	var arr []int64
	var iEnd int64
	for i := int64(1); i < 1000000000; i=i*10 {
		if i == 10 {
			continue
		}
		if i < 100 {
			iEnd = 100
		} else {
			iEnd = 90
		}
		for j:=int64(0); j<iEnd; j++ {
			if i<100 {
				arr = append(arr, int64(i + j))
			} else {
				arr = append(arr, int64(i + j*i/10))
			}
		}
	}
	return arr
}


func MaxInMap(m map[int64]int64) (int64, int64) {
	var max int64
	var maxK int64
	for k, v := range(m) {
		if max == 0 {
			max = v
			maxK = k
		} else if v > max {
			max = v
			maxK = k
		}
	}
	return max, maxK
}

// время выполнения - 0,01сек на 1000. т.е. на 100 валют майнеру и юзеру уйдет 2 сек.
func GetMaxVote(array []map[int64]int64, min, max, step int64) int64 {

	// если <=10 то тупо берем максимальное кол-во голосов
	var maxPct int64
	if len(array)<=10 {
		if len(array)>0 {
			_, maxPct = MaxInSliceMap(array)
			//fmt.Println("0maxPct", maxPct)
		} else {
			maxPct = 0
			//fmt.Println("0maxPct", maxPct)
		}
		return maxPct
	}

	// делим набор данных от $min до $max на секции кратные $step
	// таких секций будет 100 при $max=1000
	// цель такого деления - найти ту секцию, где больше всего голосов

	var sums []map[int64]int64
	dataBank := make(map[int64][]map[int64]int64)
	for i:=min; i<max; i=i+step/10 {
		//dataBank[i] = make(map[int64]int64)
		min0 := i
		max0 := i + step
		//fmt.Println("min0", min0)
		//fmt.Println("max0", max0)
		// берем из массива те данные, которые попадают в данную секцию
		for i0:=0; i0<len(array); i0++ {
			for number, votes := range array[i0] {
				if number >= min0 && number < max0 {
					dataBank[i] = append(dataBank[i], map[int64]int64{number:votes})
					//fmt.Println("i", i, "number", number, "votes", votes)
				}
			}
		}
		if len(dataBank[i])>0 {
			sums = append(sums, map[int64]int64{i:arraySum(dataBank[i])})
			//fmt.Println("i", i, dataBank[i])
		}
	}

	// ищем id секции, где больше всего голосов
	_, maxI := MaxInSliceMap(sums)
	//fmt.Println(sums)
	//fmt.Println(dataBank)
	//fmt.Println(maxI)

	// если в этой секции <= 10-и элементов, то просто выбираем максимальное кол-во голосов
	if len(dataBank[maxI]) <= 10 {
		_, maxPct := MaxInSliceMap(dataBank[maxI])
		return maxPct
	} else { // без рекурсии, просто один раз вызываем свою функцию
		return GetMaxVote(dataBank[maxI], maxI, maxI+step, step/10)
	}
	return 0
}

func arraySum(m []map[int64]int64) int64 {
	var sum int64
	for i:=0; i<len(m); i++ {
		for _, v :=range m[i] {
			sum+=v
		}
	}
	return sum
}

func MaxInSliceMap(m []map[int64]int64) (int64, int64) {
	var max int64
	var maxK int64
	for i:=0; i<len(m); i++ {
		for k, v := range (m[i]) {
			if max == 0 {
				max = v
				maxK = k
			} else if v > max {
				max = v
				maxK = k
			}
		}
	}
	return max, maxK
}

func TypesToIds(arr []string) []int64 {
	var result []int64
	for _, v := range(arr) {
		result = append(result, TypeInt(v))
	}
	return result
}

func TypeInt (txType string) int64 {
	for k, v := range consts.TxTypes {
		if v == txType {
			return int64(k)
		}
	}
	return 0
}

func MakePrivateKey(key string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return nil, errors.New("bad key data")
	}
	if got, want := block.Type, "RSA PRIVATE KEY"; got != want {
		return nil, errors.New("unknown key type "+got+", want "+want)
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func JoinInts(arr map[int]int, sep string) string {
	var arrStr []string
	for _, v := range arr {
		arrStr = append(arrStr, IntToStr(v))
	}
	return strings.Join(arrStr, sep)
}

func TimeLeft(sec int64, lang map[string]string) string {
	result := ""
	if sec > 0 {
		days := int64(math.Floor(float64(sec / 86400)))
		sec -= days*86400
		result += fmt.Sprintf(`%d %s `, days, lang["time_days"])
	}
	if sec > 0 {
		hours := int64(math.Floor(float64(sec / 3600)))
		sec -= hours*3600
		result += fmt.Sprintf(`%d %s `, hours, lang["time_hours"])
	}
	if sec > 0 {
		minutes := int64(math.Floor(float64(sec / 60)))
		sec -= minutes*3600
		result += fmt.Sprintf(`%d %s `, minutes, lang["time_minutes"])
	}
	return result
}

func GetMinersKeepers(ctx0, maxMinerId0, minersKeepers0 string, arr0 bool) map[int]int {
	ctx:=StrToInt(ctx0)
	maxMinerId:=StrToInt(maxMinerId0)
	minersKeepers:=StrToInt(minersKeepers0)
	result := make(map[int]int)
	newResult := make(map[int]int)
	var ctx_ float64
	ctx_ = float64(ctx)
	for i:=0; i<minersKeepers; i++ {
		//log.Println("ctx", ctx)
		//var hi float34
		hi := ctx_ / float64(127773)
		//log.Println("hi", hi)
		lo := int(ctx_) % 127773
		//log.Println("lo", lo)
		x := (float64(16807) * float64(lo)) - (float64(2836) * hi)
		//log.Println("x", x, float64(16807), float64(lo), float64(2836), hi)
		if x <= 0 {
			x += 0x7fffffff
		}
		ctx_ = x
		rez := int(ctx_) % (maxMinerId+1)
		//log.Println("rez", rez)
		if rez == 0 {
			rez = 1
		}
		result[rez] = 1
	}
	if arr0 {
		i:=0
		for k, _ := range result {
			newResult[i] = k
			i++
		}
	} else {
		newResult = result
	}
	return newResult
}

func MakeLastTx(lastTx []map[string]string, lng map[string]string) (string, map[int64]int64) {
	pendingTx := make(map[int64]int64)
	result := `<h3>`+lng["transactions"]+`</h3><table class="table" style="width:500px;">`
	result+=`<tr><th>`+lng["time"]+`</th><th>`+lng["result"]+`</th></tr>`
	for _, data := range(lastTx) {
		result+="<tr>"
		result+="<td class='unixtime'>"+data["time_int"]+"</td>"
		if len(data["block_id"]) > 0 {
			result+="<td>"+lng["in_the_block"]+" "+data["block_id"]+"</td>"
		} else if len(data["error"]) > 0 {
			result+="<td>Error: "+data["error"]+"</td>"
		} else if (len(data["queue_tx"]) == 0 && len(data["tx"]) == 0) || time.Now().Unix() - StrToInt64(data["time_int"]) > 7200  {
			result+="<td>"+lng["lost"]+"</td>"
		} else {
			result+="<td>"+lng["status_pending"]+"</td>"
			pendingTx[StrToInt64(data["type"])] = 1
		}
		result+="</tr>"
	}
	result+="</table>"
	return result, pendingTx
}
