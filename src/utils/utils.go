package utils

import (
	"fmt"
	_ "github.com/lib/pq"
	//"time"
	"strings"
	//"database/sql"
	"runtime"
	"path/filepath"
	"strconv"
	"errors"
	//"encoding/json"
	//"crypto/x509"
	//"encoding/pem"
	"crypto"
	"consts"
	//"crypto/rand"
	//"crypto/rsa"
	//"strings"
	"math"
	//math_rand "math/rand"
	"crypto/sha256"
	//"crypto/md5"
	"encoding/hex"
	"crypto/rsa"
	//"bufio"
	//"os"
	//"errors"
	//"io/ioutil"*/
	"reflect"
	"regexp"
	//"fmt"
	"os"
	"net/http"
	"io"
	"time"
	"log"
	"math/big"
	"crypto/x509"
	"encoding/base64"
	"math/rand"
	"encoding/pem"
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

func round(num float64) int {
	//num += ROUND_FIX
	//	return int(StrToFloat64(Float64ToStr(num)) + math.Copysign(0.5, num))
	return int(num + math.Copysign(0.5, num))
}

func Round(num float64, precision int) float64 {
	num += consts.ROUND_FIX
	//num = StrToFloat64(Float64ToStr(num))
	output := math.Pow(10, float64(precision))
	return float64(round(num * output)) / output
}



func CheckDaemonRestart(name string) bool {
	return false
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
	switch dataType {
	case "arbitration_trust_list":
		if ok, _ := regexp.MatchString(`^\[[0-9]{1,10}(,[0-9]{1,10}){0,100}\]$`, data); ok{
			return true
		}
	case "votes_comment", "cf_comment":
		if ok, _ := regexp.MatchString(`^[\pL0-9\,\s\.\-\:\=\;\?\!\%\)\(\@\/\n\r]{1,140}$`, data); ok{
			return true
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
			return []byte(fmt.Sprintf("0%x", len0))
		}
		return []byte(fmt.Sprintf("%x", len0))
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

func InterfaceToStr(i []interface {}) []string {
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
	n_ := []byte(HexToBin(hex_n))
	n_ = append([]byte("02"), EncodeLength(int64(len(HexToBin(hex_n))))...)
	n_ = append(n_, hex_n...)
	e_ := append([]byte("02"), EncodeLength(int64(len(HexToBin(hex_e))))...)
	e_ = append(e_, hex_e...)
	length := EncodeLength(int64(len(HexToBin(append(n_,e_...)))))
	rez := append([]byte("30"), length...)
	rez = append(rez, n_...)
	rez = append(rez, e_...)
	return rez
	//b64:=base64.StdEncoding.EncodeToString([]byte(utils.HexToBin("30"+length+bin_enc)))
	//fmt.Println(b64)
}

func CheckSign(publicKeys [][]byte, forSign string, signs []byte, nodeKeyOrLogin bool ) (bool, error) {
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
		//fmt.Println("publicKeys[i]", publicKeys[i])
		key := base64.StdEncoding.EncodeToString(publicKeys[i])
		key = "-----BEGIN PUBLIC KEY-----\n"+key+"\n-----END PUBLIC KEY-----"
		//fmt.Printf("%x\n", publicKeys[i])
		//fmt.Println("key", key)
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
		}
		err = rsa.VerifyPKCS1v15(pub, crypto.SHA1,  HashSha1(forSign), signsSlice[i])
		if err != nil {
			//fmt.Println("pub", pub)
			//fmt.Println("crypto.SHA1", crypto.SHA1)
			//fmt.Println("HashSha1(forSign)", HashSha1(forSign))
			//fmt.Println("HashSha1(forSign)", HashSha1(forSign))
			//fmt.Println("forSign", forSign)
			//fmt.Printf("sign: %x\n", signsSlice[i])
			return false, ErrInfoFmt("incorrect sign: key = %s ; hash = %x; forSign = %v", key, HashSha1(forSign), forSign)
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
func GetMaxVote(array  []map[int64]int64, min, max, step int64) int64 {

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

func TypeArray (txType string) int64 {
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
