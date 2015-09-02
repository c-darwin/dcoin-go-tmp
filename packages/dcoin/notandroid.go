// +build !android,!ios

package dcoin

import  (
	_ "github.com/mattn/go-sqlite3"
	"syscall"
	"net"
	"os"
	"net/http"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"github.com/c-darwin/dcoin-go-tmp/packages/tcpserver"
	"github.com/c-darwin/dcoin-go-tmp/packages/daemons"
	"os/signal"
	"time"
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

func waitSig() {
	C.waitSig()
}


func httpListener(ListenHttpHost, BrowserHttpHost string) {
	l, err := net.Listen("tcp", ListenHttpHost)
	if err != nil {
		log.Error("%v", err)
		// Если это повторный запуск и он не из консоли, то открываем окно браузера, т.к. скорее всего юзер тыкнул по иконке
		if *utils.Console == 0 {
			openBrowser(BrowserHttpHost)
		}
		log.Error("%v", utils.ErrInfo(err))
		panic(err)
		os.Exit(1)
	}

	go func() {
		err = http.Serve(NewBoundListener(100, l), http.TimeoutHandler(http.DefaultServeMux, time.Duration(600 * time.Second), "Your request has timed out"))
		if err != nil {
			log.Error("Error listening: %v (%v)", err, ListenHttpHost)
			panic(err)
			//os.Exit(1)
		}
	}()
}

func tcpListener(db *utils.DCDB) {
	log.Debug("tcp")
	go func() {
		if db == nil {
			for {
				db = utils.DB
				if db!=nil {
					break
				}
			}
		}
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
}

func signals(countDaemons int) {
	SigChan = make(chan os.Signal, 1)
	waitSig()
	var Term os.Signal = syscall.SIGTERM
	go func() {
		signal.Notify(SigChan, os.Interrupt, os.Kill, Term)
		<-SigChan
		log.Debug("countDaemons %v", countDaemons)
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
	}()
}