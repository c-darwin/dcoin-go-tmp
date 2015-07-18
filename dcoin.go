package main
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/daemons"
	"os"
	"net/http"
	_ "image/png"
	"github.com/c-darwin/dcoin-go-tmp/packages/controllers"
	"github.com/astaxie/beego/config"
    "github.com/elazarl/go-bindata-assetfs"
	"github.com/c-darwin/dcoin-go-tmp/packages/static"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"io/ioutil"
	"math/rand"
	"time"
	"strings"
	"net"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"github.com/op/go-logging"
	"runtime"
	"os/exec"
	"fmt"
)

var log = logging.MustGetLogger("example")
//var format = logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{shortfile} %{shortfunc} [%{level:.4s}] %{color:reset} %{message}")
var format = logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{shortfile} %{shortfunc} [%{level:.4s}] %{color:reset} %{message}"+string(byte(0)))

var configIni map[string]string

func main() {
	// читаем config.ini
	if _, err := os.Stat("config.ini"); os.IsNotExist(err) {
		d1 := []byte(`
error_log=1
log=1
sql_log=0
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
log_level=DEBUG
log_output=file
db_name=`)
		ioutil.WriteFile("config.ini", d1, 0644)
	}
	configIni_, err := config.NewConfig("ini", "config.ini")
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	configIni, err = configIni_.GetSection("default")

	f, err := os.OpenFile("dclog.txt", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0777)
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	defer f.Close()
	var backend *logging.LogBackend
	switch configIni["log_output"] {
	case "file":
		backend = logging.NewLogBackend(f, "", 0)
	case "console":
		backend = logging.NewLogBackend(os.Stderr, "", 0)
	case "file_console":
		backend = logging.NewLogBackend(io.MultiWriter(f, os.Stderr), "", 0)
	default:
		backend = logging.NewLogBackend(f, "", 0)
	}
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(logging.DEBUG, "")
	logging.SetBackend(backendLeveled)

	rand.Seed( time.Now().UTC().UnixNano())

	if _, err := os.Stat("public"); os.IsNotExist(err) {
		err = os.Mkdir("public", 0755)
		if err != nil {
			log.Error("%v", err)
			panic(err)
			os.Exit(1)
		}
	}

	daemonsStart := map[string]func(){"TestblockIsReady":daemons.TestblockIsReady,"TestblockGenerator":daemons.TestblockGenerator,"TestblockDisseminator":daemons.TestblockDisseminator,"Shop":daemons.Shop,"ReductionGenerator":daemons.ReductionGenerator,"QueueParserTx":daemons.QueueParserTx,"QueueParserTestblock":daemons.QueueParserTestblock,"QueueParserBlocks":daemons.QueueParserBlocks,"PctGenerator":daemons.PctGenerator,"Notifications":daemons.Notifications,"NodeVoting":daemons.NodeVoting,"MaxPromisedAmountGenerator":daemons.MaxPromisedAmountGenerator,"MaxOtherCurrenciesGenerator":daemons.MaxOtherCurrenciesGenerator,"ElectionsAdmin":daemons.ElectionsAdmin,"Disseminator":daemons.Disseminator,"Confirmations":daemons.Confirmations,"Connector":daemons.Connector,"Clear":daemons.Clear,"CleaningDb":daemons.CleaningDb,"CfProjects":daemons.CfProjects,"BlocksCollection":daemons.BlocksCollection}

	if len(configIni["daemons"]) > 0 && configIni["daemons"] != "null" {
		daemonsConf := strings.Split(configIni["daemons"], ",")
		for _, fns := range daemonsConf {
			go daemonsStart[fns]()
		}
	} else if configIni["daemons"] != "null" {
		for _, fns := range daemonsStart {
			go fns()
		}
	}

	// включаем листинг веб-сервером для клиентской части
	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/content", controllers.Content)
	http.HandleFunc("/ajax", controllers.Ajax)
	http.HandleFunc("/tools", controllers.Tools)
	http.HandleFunc("/cf/", controllers.IndexCf)
	http.HandleFunc("/cf/content", controllers.ContentCf)
	http.Handle("/public/", noDirListing(http.FileServer(http.Dir("./"))))
	http.Handle("/static/", http.FileServer(&assetfs.AssetFS{Asset: static.Asset, AssetDir: static.AssetDir, Prefix: ""}))

	db := utils.DbConnect(configIni)
	go func() {
		tcpPort := db.GetTcpPort()
		// включаем листинг TCP-сервером и обработку входящих запросов
		l, err := net.Listen("tcp", ":"+tcpPort)
		log.Debug("tcpPort: %v", tcpPort)
		if err != nil {
			log.Error("Error listening: %v", err)
			panic(err)
			os.Exit(1)
		}
		defer l.Close()
		go func() {
			for {
				conn, err := l.Accept()
				if err != nil {
					log.Error("Error accepting: %v", err)
					panic(err)
					os.Exit(1)
				}

				go utils.HandleTcpRequest(conn, configIni)
			}
		}()
	}()

	HttpPort := db.GetHttpPort()
	err = http.ListenAndServe(":"+HttpPort, nil)
	if err != nil {
		log.Error("Error listening: %v", err)
		panic(err)
		os.Exit(1)
	}


	log.Debug("runtime.GOOS: %v", runtime.GOOS)
	err = nil
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", "http://localhost:"+HttpPort+"/").Start()
	case "windows", "darwin":
		err = exec.Command("open", "http://localhost:"+HttpPort+"/").Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Error("%v", err)
		panic(err)
		os.Exit(1)
	}

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

