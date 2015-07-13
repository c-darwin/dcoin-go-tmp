package main
import (
	"fmt"
	"dcoin/packages/daemons"
	"log"
	"os"
	"net/http"
	_ "image/png"
	"dcoin/packages/controllers"
	"github.com/astaxie/beego/config"
    "github.com/elazarl/go-bindata-assetfs"
	"dcoin/packages/static"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"math/rand"
	"time"
	"strings"
	"net"
	"dcoin/packages/utils"
)

var configIni map[string]string
func main() {

	rand.Seed( time.Now().UTC().UnixNano())

	if _, err := os.Stat("public"); os.IsNotExist(err) {
		os.Mkdir("public", 0755)
	}

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
db_user=
db_host=
db_port=
db_password=
db_name=`)
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
	//log.SetOutput(f)
	log.SetOutput(io.MultiWriter(f, os.Stdout))
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// запускаем всех демонов
	go daemons.TestblockIsReady(configIni)
	go daemons.TestblockGenerator(configIni)
	go daemons.TestblockDisseminator(configIni)
	go daemons.Shop(configIni)
	go daemons.ReductionGenerator(configIni)
	go daemons.QueueParserTx(configIni)
	go daemons.QueueParserTestblock(configIni)
	go daemons.QueueParserBlocks(configIni)
	go daemons.PctGenerator(configIni)
	go daemons.Notifications(configIni)
	go daemons.NodeVoting(configIni)
	go daemons.MaxPromisedAmountGenerator(configIni)
	go daemons.MaxOtherCurrenciesGenerator(configIni)
	go daemons.ElectionsAdmin(configIni)
	go daemons.Disseminator(configIni)
	go daemons.Confirmations(configIni)
	go daemons.Connector(configIni)
	go daemons.Clear(configIni)
	go daemons.CleaningDb(configIni)
	go daemons.CfProjects(configIni)
	go daemons.BlocksCollection(configIni)


	// включаем листинг веб-сервером для клиентской части
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/content", controllers.Content)
	http.HandleFunc("/ajax", controllers.Ajax)
	http.HandleFunc("/tools", controllers.Tools)
	http.HandleFunc("/cf/", controllers.IndexCf)
	http.HandleFunc("/cf/content", controllers.ContentCf)
	http.Handle("/public/", noDirListing(http.FileServer(http.Dir("./"))))
	http.Handle("/static/", http.FileServer(&assetfs.AssetFS{Asset: static.Asset, AssetDir: static.AssetDir, Prefix: ""}))


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


	// включаем листинг TCP-сервером и обработку входящих запросов
	l, err := net.Listen("tcp", "localhost:8088")
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	defer l.Close()
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				log.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}
			go utils.HandleTcpRequest(conn, configIni)
		}
	}()
	fmt.Scanln()

}

// http://grokbase.com/t/gg/golang-nuts/12a9yhgr64/go-nuts-disable-directory-listing-with-http-fileserver#201210093cnylxyosmdfuf3wh5xqnwiut4
func noDirListing(h http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}
