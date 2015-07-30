package main
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/daemons"
	"os"
	"flag"
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
	"os/signal"
	"syscall"
	"github.com/c-darwin/dcoin-go-tmp/packages/tcpserver"
)
/*
#include <stdio.h>
#include <signal.h>

extern void go_callback_int();
static inline void SigBreak_Handler(int n_signal){
      printf("closed\n");
	go_callback_int();
}
static inline void waitSig() {
    #if (WIN32 || WIN64)
    signal(SIGBREAK, &SigBreak_Handler);
    signal(SIGINT, &SigBreak_Handler);
    #endif
}
*/
import "C"

//export go_callback_int
func go_callback_int(){
	SigChan <- syscall.Signal(1)
}

var SigChan chan os.Signal

var log = logging.MustGetLogger("example")
//var format = logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{shortfile} %{shortfunc} [%{level:.4s}] %{color:reset} %{message}")
var format = logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{shortfile} %{shortfunc} [%{level:.4s}] %{color:reset} %{message}"+string(byte(0)))

var configIni map[string]string

var console = flag.Int64("console", 0, "Start from console")

/*func init() {
	flag.StringVar(console, "console", 0, "Description")
}
*/
func main() {
	dir, err := utils.GetCurrentDir()
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	// читаем config.ini
	if _, err := os.Stat(dir+"config.ini"); os.IsNotExist(err) {
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
		ioutil.WriteFile(dir+"/config.ini", d1, 0644)
	}
	configIni_, err := config.NewConfig("ini", dir+"/config.ini")
	if err != nil {
		panic(err)
		os.Exit(1)
	}
	configIni, err = configIni_.GetSection("default")

	f, err := os.OpenFile(dir+"/dclog.txt", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0777)
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
	logLevel, err := logging.LogLevel(configIni["log_level"])
	if err!= nil {
		logLevel = logging.DEBUG
	}

	log.Debug("logLevel: %v", logLevel)
	backendLeveled.SetLevel(logLevel, "")
	logging.SetBackend(backendLeveled)

	rand.Seed( time.Now().UTC().UnixNano())

	if _, err := os.Stat(dir+"/public"); os.IsNotExist(err) {
		err = os.Mkdir(dir+"/public", 0755)
		if err != nil {
			log.Error("%v", err)
			panic(err)
			os.Exit(1)
		}
	}

	daemons.DaemonCh = make(chan bool, 1)
	daemons.AnswerDaemonCh = make(chan bool, 1)
	log.Debug("daemonsStart")
	//TestblockIsReady,TestblockGenerator,TestblockDisseminator,Shop,ReductionGenerator,QueueParserTx,QueueParserTestblock,QueueParserBlocks,PctGenerator,Notifications,NodeVoting,MaxPromisedAmountGenerator,MaxOtherCurrenciesGenerator,ElectionsAdmin,Disseminator,Confirmations,Connector,Clear,CleaningDb,CfProjects,BlocksCollection
	daemonsStart := map[string]func(){"TestblockIsReady":daemons.TestblockIsReady,"TestblockGenerator":daemons.TestblockGenerator,"TestblockDisseminator":daemons.TestblockDisseminator,"Shop":daemons.Shop,"ReductionGenerator":daemons.ReductionGenerator,"QueueParserTx":daemons.QueueParserTx,"QueueParserTestblock":daemons.QueueParserTestblock,"QueueParserBlocks":daemons.QueueParserBlocks,"PctGenerator":daemons.PctGenerator,"Notifications":daemons.Notifications,"NodeVoting":daemons.NodeVoting,"MaxPromisedAmountGenerator":daemons.MaxPromisedAmountGenerator,"MaxOtherCurrenciesGenerator":daemons.MaxOtherCurrenciesGenerator,"ElectionsAdmin":daemons.ElectionsAdmin,"Disseminator":daemons.Disseminator,"Confirmations":daemons.Confirmations,"Connector":daemons.Connector,"Clear":daemons.Clear,"CleaningDb":daemons.CleaningDb,"CfProjects":daemons.CfProjects,"BlocksCollection":daemons.BlocksCollection}

	countDaemons := 0
	if len(configIni["daemons"]) > 0 && configIni["daemons"] != "null" {
		daemonsConf := strings.Split(configIni["daemons"], ",")
		for _, fns := range daemonsConf {
			go daemonsStart[fns]()
			countDaemons++
		}
	} else if configIni["daemons"] != "null" {
		for _, fns := range daemonsStart {
			go fns()
			countDaemons++
		}
	}

	SigChan = make(chan os.Signal, 1)
	C.waitSig()
	var Term os.Signal = syscall.SIGTERM
	go func() {
		signal.Notify(SigChan, os.Interrupt, os.Kill, Term)
		<-SigChan
		log.Debug("countDaemons %v", countDaemons)
		for i:=0; i < countDaemons; i++ {
			daemons.DaemonCh <- true
			log.Debug("daemons.DaemonCh <- true")
			answer := <-daemons.AnswerDaemonCh
			log.Debug("answer: %v", answer)
		}
		os.Exit(1)
	}()


	db := utils.DbConnect(configIni)
	BrowserHttpHost, HandleHttpHost, ListenHttpHost := db.GetHttpHost()
	log.Error("BrowserHttpHost: %v, HandleHttpHost: %v, ListenHttpHost: %v", BrowserHttpHost, HandleHttpHost, ListenHttpHost)
	// включаем листинг веб-сервером для клиентской части
	http.HandleFunc(HandleHttpHost+"/", controllers.Index)
	http.HandleFunc(HandleHttpHost+"/content", controllers.Content)
	http.HandleFunc(HandleHttpHost+"/ajax", controllers.Ajax)
	http.HandleFunc(HandleHttpHost+"/tools", controllers.Tools)
	http.HandleFunc(HandleHttpHost+"/cf/", controllers.IndexCf)
	http.HandleFunc(HandleHttpHost+"/cf/content", controllers.ContentCf)
	http.Handle(HandleHttpHost+"/public/", noDirListing(http.FileServer(http.Dir("./"))))
	http.Handle(HandleHttpHost+"/static/", http.FileServer(&assetfs.AssetFS{Asset: static.Asset, AssetDir: static.AssetDir, Prefix: ""}))

	log.Debug("tcp")
	go func() {
		tcpHost := db.GetTcpHost()
		log.Debug("tcpHost: %v", tcpHost)
		// включаем листинг TCP-сервером и обработку входящих запросов
		l, err := net.Listen("tcp", tcpHost)
		if err != nil {
			log.Error("Error listening: %v", err)
			panic(err)
			//os.Exit(1)
		}
		//defer l.Close()
		go func() {
			for {
				conn, err := l.Accept()
				if err != nil {
					log.Error("Error accepting: %v", err)
					utils.Sleep(1)
					//panic(err)
					//os.Exit(1)
				} else {
					go func(conn net.Conn) {
						t := new(tcpserver.TcpServer)
						t.DCDB = db
						t.Conn = conn
						t.HandleTcpRequest()
					}(conn)
				}
			}
		}()
	}()

	log.Debug("ListenHttpHost", ListenHttpHost)
	//err = http.ListenAndServe(ListenHttpHost, nil)
	l, err := net.Listen("tcp", ListenHttpHost)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		err = http.Serve(NewBoundListener(20, l), http.DefaultServeMux)
		if err != nil {
			log.Error("Error listening: %v (%v)", err, ListenHttpHost)
			panic(err)
			//os.Exit(1)
		}
	}()

	utils.Sleep(3)

	flag.Parse()
	if *console == 0 {
		log.Debug("runtime.GOOS: %v", runtime.GOOS)
		err = nil
		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", BrowserHttpHost).Start()
		case "windows", "darwin":
			err = exec.Command("open", BrowserHttpHost).Start()
			if err!=nil {
				exec.Command("cmd", "/c", "start", BrowserHttpHost).Start()
			}
		default:
			err = fmt.Errorf("unsupported platform")
		}
		if err != nil {
			log.Error("%v", err)
			//panic(err)
			//os.Exit(1)
		}
	}

	log.Debug("ALL RIGHT")
	utils.Sleep(3600*24*90)
	log.Debug("EXIT")
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


func NewBoundListener(maxActive int, l net.Listener) net.Listener {
	return &boundListener{l, make(chan bool, maxActive)}
}

type boundListener struct {
	net.Listener
	active chan bool
}

type boundConn struct {
	net.Conn
	active chan bool
}

func (l *boundListener) Accept() (net.Conn, error) {
	l.active <- true
	c, err := l.Listener.Accept()
	if err != nil {
		<-l.active
		return nil, err
	}
	return &boundConn{c, l.active}, err
}

func (l *boundConn) Close() error {
	err := l.Conn.Close()
	<-l.active
	return err
}

