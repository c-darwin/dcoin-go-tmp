package main
import (
	"fmt"
	"schema"
	//"encoding/binary"
	//"bytes"
	//"encoding/hex"
	//"crypto/rand"
	//"crypto/rsa"
	//"crypto/sha1"
	"daemons"
	//"dcparser"
	"log"
	"os"
	//"github.com/alyu/configparser"
	"net/http"
	//"html/template"
	//"image/jpeg"
	_ "image/png"
	//"image"
	"controllers"
	//"bufio"
	//"utils"
	//"strconv"
	//"bytes"
	//"encoding/binary"
	//"regexp"
	//"time"
	//"github.com/beego/i18n"
	//"github.com/nicksnyder/go-i18n/i18n"
	//"consts"
	//"strconv"
	//"reflect"
	"github.com/astaxie/beego/config"
    "github.com/elazarl/go-bindata-assetfs"
	"bindatastatic"
	_ "github.com/mattn/go-sqlite3"
	//"runtime"
//	"database/sql"
	//"os/exec"
	"io/ioutil"
)

var configIni map[string]string
func main() {
/*
	t := time.Unix(time.Now().Unix(), 0)
	fmt.Println(t.Format("2006-01-02 15:04:05"))
*/

	/*enc := utils.Encode_length(443343);
	fmt.Println(enc);
	bin_enc := utils.HexToBin(enc)+utils.HexToBin("FFFFFF")
	fmt.Println("bin_enc", utils.BinToHex(bin_enc));
	DecodeLength(&bin_enc)
	fmt.Println("bin_enc", utils.BinToHex(bin_enc));
	/*
	xxx:="9876"
	Shift := StringShift(&xxx, 1);
	fmt.Println("Shift", Shift);
	fmt.Println("xxx", xxx);*/

/*
	langIni_, err := configparser.Read("lang/1.ini")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(langIni_.AllSections())
	//_, err := langIni_.Section("")
	//fmt.Println(langIni)
	//v:= langIni.ValueOf("picture_description")
	//fmt.Println(langIni.String())
*/
	//fmt.Println(xx)
	schema.GetSchema("sqlite")
	// читаем config.ini
	if _, err := os.Stat("config.ini"); os.IsNotExist(err) {
		fmt.Println("NO")
		d1 := []byte(`
error_log=1
log=1
log_block_id_begin=0
log_block_id_end=0
bad_tx_log=1
nodes_ban_exit=0
log_tables=
log_fns=
sign_hash=ip
db_type=sqlite
DB_USER=
DB_PASSWORD=
DB_NAME=`)
		ioutil.WriteFile("config.ini", d1, 0644)
	} else {
		fmt.Println("YES")
	}
	configIni_, err := config.NewConfig("ini", "config.ini")
	if err != nil {
		//log.Fatal(err)
	}
	configIni, err = configIni_.GetSection("default")
	fmt.Println("configIni[log]", configIni["log"])

	f, err := os.OpenFile("dclog.txt", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0777)
	defer f.Close()
	log.SetOutput(f)
	log.SetFlags(log.LstdFlags | log.Lshortfile)


	if len(os.Args)>1 {
		// запускаем всех демонов
		go daemons.Testblock_generator(configIni)
		go daemons.Testblock_is_ready()
	}
	// включаем листинг TCP-сервером и обработку входящих запросов

	// включаем листинг веб-сервером для клиентской части
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/content", controllers.Content)
	http.HandleFunc("/ajax", controllers.Ajax)

	http.Handle("/static/", http.FileServer(&assetfs.AssetFS{Asset: bindatastatic.Asset, AssetDir: bindatastatic.AssetDir, Prefix: ""}))
/*
fmt.Println(runtime.GOOS)
switch runtime.GOOS {
case "linux":
err = exec.Command("xdg-open", "http://localhost:8089/").Start()
case "windows", "darwin":
err = exec.Command("open", "http://localhost:8089/").Start()
default:
err = fmt.Errorf("unsupported platform")
}*/


	http.ListenAndServe(":8089", nil)


	fmt.Scanln()

}
