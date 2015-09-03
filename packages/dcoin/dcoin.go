package dcoin
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/daemons"
	"os"
	"net/http"
	_ "image/png"
	"github.com/c-darwin/dcoin-go-tmp/packages/controllers"
	"github.com/astaxie/beego/config"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/c-darwin/dcoin-go-tmp/packages/static"
	"io"
	"io/ioutil"
	"math/rand"
	"time"
	"strings"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"github.com/op/go-logging"
	"runtime"
	"os/exec"
	"fmt"
	"github.com/astaxie/beego/session"
	"github.com/c-darwin/dcoin-go-tmp/packages/consts"
)
/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation
#import <Foundation/Foundation.h>

void
logNS(char* text) {
    NSLog(@"golog: %s", text);
}

*/
import "C"

var (
	log = logging.MustGetLogger("dcoin")
	format = logging.MustStringFormatter("%{color}%{time:15:04:05.000} %{shortfile} %{shortfunc} [%{level:.4s}] %{color:reset} %{message}["+consts.VERSION+"]"+string(byte(0)))
	configIni map[string]string
	globalSessions *session.Manager
)

func iosLog(text string) {
	if utils.IOS() {
		C.logNS(C.CString(text))
	}
}


func Stop() {
	var err error
	utils.DB, err = utils.NewDbConnect(configIni)
	log.Debug("%v", utils.DB)
	iosLog("utils.DB:"+fmt.Sprintf("%v", utils.DB))
	if err != nil {
		iosLog("err:"+fmt.Sprintf("%s", utils.ErrInfo(err)))
		log.Error("%v", utils.ErrInfo(err))
		panic(err)
		os.Exit(1)
	}
	err = utils.DB.ExecSql(`INSERT INTO stop_daemons(stop_time) VALUES (?)`, utils.Time())
	if err != nil {
		iosLog("err:"+fmt.Sprintf("%s", utils.ErrInfo(err)))
		log.Error("%v", utils.ErrInfo(err))
	}
}

func Start(dir string) {

	iosLog("start")

	defer func() {
		if r := recover(); r != nil {
			log.Error("Recovered", r)
			panic(r)
		}
	}()

	if dir!="" {
		*utils.Dir = dir
	}

	iosLog("dir:"+dir)

	fmt.Println("dcVersion:", consts.VERSION)
	log.Debug("dcVersion: %v", consts.VERSION)
	// читаем config.ini
	if _, err := os.Stat(*utils.Dir+"/config.ini"); os.IsNotExist(err) {
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
		ioutil.WriteFile(*utils.Dir+"/config.ini", d1, 0644)
		iosLog("config ok")
	}
	configIni_, err := config.NewConfig("ini", *utils.Dir+"/config.ini")
	if err != nil {
		iosLog("err:"+fmt.Sprintf("%s", utils.ErrInfo(err)))
		log.Error("%v", utils.ErrInfo(err))
		panic(err)
		os.Exit(1)
	}
	configIni, err = configIni_.GetSection("default")

	controllers.SessInit()
	controllers.ConfigInit()
	daemons.ConfigInit()

	go func() {
		utils.DB, err = utils.NewDbConnect(configIni)
		log.Debug("%v", utils.DB)
		iosLog("utils.DB:"+fmt.Sprintf("%v", utils.DB))
		if err != nil {
			iosLog("err:"+fmt.Sprintf("%s", utils.ErrInfo(err)))
			log.Error("%v", utils.ErrInfo(err))
			panic(err)
			os.Exit(1)
		}
	}()

	f, err := os.OpenFile(*utils.Dir+"/dclog.txt", os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0777)
	if err != nil {
		iosLog("err:"+fmt.Sprintf("%s", utils.ErrInfo(err)))
		log.Error("%v", utils.ErrInfo(err))
		panic(err)
		os.Exit(1)
	}
	defer f.Close()
	iosLog("configIni:"+fmt.Sprintf("%v", configIni))
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

	log.Debug("public")
	iosLog("public")
	if _, err := os.Stat(*utils.Dir+"/public"); os.IsNotExist(err) {
		err = os.Mkdir(*utils.Dir+"/public", 0755)
		if err != nil {
			log.Error("%v", utils.ErrInfo(err))
			panic(err)
			os.Exit(1)
		}
	}

	daemons.DaemonCh = make(chan bool, 1)
	daemons.AnswerDaemonCh = make(chan bool, 1)
	log.Debug("daemonsStart")
	iosLog("daemonsStart")
	//TestblockIsReady,TestblockGenerator,TestblockDisseminator,Shop,ReductionGenerator,QueueParserTx,QueueParserTestblock,QueueParserBlocks,PctGenerator,Notifications,NodeVoting,MaxPromisedAmountGenerator,MaxOtherCurrenciesGenerator,ElectionsAdmin,Disseminator,Confirmations,Connector,Clear,CleaningDb,CfProjects,BlocksCollection
	daemonsStart := map[string]func(){"TestblockIsReady":daemons.TestblockIsReady,"TestblockGenerator":daemons.TestblockGenerator,"TestblockDisseminator":daemons.TestblockDisseminator,"Shop":daemons.Shop,"ReductionGenerator":daemons.ReductionGenerator,"QueueParserTx":daemons.QueueParserTx,"QueueParserTestblock":daemons.QueueParserTestblock,"QueueParserBlocks":daemons.QueueParserBlocks,"PctGenerator":daemons.PctGenerator,"Notifications":daemons.Notifications,"NodeVoting":daemons.NodeVoting,"MaxPromisedAmountGenerator":daemons.MaxPromisedAmountGenerator,"MaxOtherCurrenciesGenerator":daemons.MaxOtherCurrenciesGenerator,"ElectionsAdmin":daemons.ElectionsAdmin,"Disseminator":daemons.Disseminator,"Confirmations":daemons.Confirmations,"Connector":daemons.Connector,"Clear":daemons.Clear,"CleaningDb":daemons.CleaningDb,"CfProjects":daemons.CfProjects,"BlocksCollection":daemons.BlocksCollection}
	if utils.Mobile() {
		daemonsStart = map[string]func(){"QueueParserTx":daemons.QueueParserTx,"Notifications":daemons.Notifications,"Disseminator":daemons.Disseminator,"Confirmations":daemons.Confirmations,"Connector":daemons.Connector,"Clear":daemons.Clear,"CleaningDb":daemons.CleaningDb,"BlocksCollection":daemons.BlocksCollection}

	}

	countDaemons := 0
	if len(configIni["daemons"]) > 0 && configIni["daemons"] != "null" {
		daemonsConf := strings.Split(configIni["daemons"], ",")
		for _, fns := range daemonsConf {
			log.Debug("start daemon %s", fns)
			go daemonsStart[fns]()
			countDaemons++
		}
	} else if configIni["daemons"] != "null" {
		for dName, fns := range daemonsStart {
			log.Debug("start daemon %s", dName)
			go fns()
			countDaemons++
		}
	}


	iosLog("MonitorDaemons")
	// мониторинг демонов
	daemonsTable := make(map[string]string)
	go func() {
		for {
			daemonNameAndTime := <-daemons.MonitorDaemonCh
			daemonsTable[daemonNameAndTime[0]] = daemonNameAndTime[1]
			if utils.Time()%10 == 0 {
				log.Error("daemonsTable: %v\n", daemonsTable)
			}
		}
	} ()



	iosLog("signals")
	// сигналы демонам для выхода
	signals(countDaemons)

	utils.Sleep(1)
	db := utils.DB


	iosLog("stop_daemons")
	// мониторим сигнал из БД о том, что демонам надо завершаться
	go func() {
		for {
			dExtit, err := db.Single(`SELECT stop_time FROM stop_daemons`).Int64()
			if err != nil {
				iosLog("err:"+fmt.Sprintf("%s", utils.ErrInfo(err)))
				log.Error("%v", utils.ErrInfo(err))
			}
			if dExtit > 0 {
				for i := 0; i < countDaemons; i++ {
					daemons.DaemonCh <- true
					log.Debug("daemons.DaemonCh <- true")
					answer := <-daemons.AnswerDaemonCh
					log.Debug("answer: %v", answer)
				}
				err := utils.DB.Close()
				if err != nil {
					log.Error("%v", utils.ErrInfo(err))
					panic(err)
				}
				os.Exit(1)
			}
			utils.Sleep(1)
		}
	} ()



	iosLog("BrowserHttpHost")
	BrowserHttpHost := "http://localhost:8089"
	HandleHttpHost := ""
	ListenHttpHost := ":8089"
	if db != nil {
		BrowserHttpHost, HandleHttpHost, ListenHttpHost = db.GetHttpHost()
	}
	iosLog(fmt.Sprintf("BrowserHttpHost: %v, HandleHttpHost: %v, ListenHttpHost: %v", BrowserHttpHost, HandleHttpHost, ListenHttpHost))
	log.Debug("BrowserHttpHost: %v, HandleHttpHost: %v, ListenHttpHost: %v", BrowserHttpHost, HandleHttpHost, ListenHttpHost)
	// включаем листинг веб-сервером для клиентской части
	http.HandleFunc(HandleHttpHost+"/", controllers.Index)
	http.HandleFunc(HandleHttpHost+"/content", controllers.Content)
	http.HandleFunc(HandleHttpHost+"/ajax", controllers.Ajax)
	http.HandleFunc(HandleHttpHost+"/tools", controllers.Tools)
	http.HandleFunc(HandleHttpHost+"/cf/", controllers.IndexCf)
	http.HandleFunc(HandleHttpHost+"/cf/content", controllers.ContentCf)
	http.Handle(HandleHttpHost+"/public/", noDirListing(http.FileServer(http.Dir(*utils.Dir))))
	http.Handle(HandleHttpHost+"/static/", http.FileServer(&assetfs.AssetFS{Asset: static.Asset, AssetDir: static.AssetDir, Prefix: ""}))

	log.Debug("ListenHttpHost", ListenHttpHost)

	iosLog(fmt.Sprintf("ListenHttpHost: %v", ListenHttpHost))

	httpListener(ListenHttpHost, BrowserHttpHost)

	tcpListener(db)

	utils.Sleep(3)

	iosLog("Sleep")

	if *utils.Console == 0 && !utils.Mobile() {
		openBrowser(BrowserHttpHost)
	}

	log.Debug("ALL RIGHT")
	iosLog("ALL RIGHT")
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


func openBrowser(BrowserHttpHost string) {
	log.Debug("runtime.GOOS: %v", runtime.GOOS)
	var err error
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
	}
}

