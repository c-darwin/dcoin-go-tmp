// +build darwin
// +build arm arm64

package dcoin

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

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"github.com/hydrogen18/stoppableListener"
	"net"
	"net/http"
)

var stop = make(chan bool)

func IosLog(text string) {
	if utils.IOS() {
		C.logNS(C.CString(text))
	}
}

func StartHTTPServer(ListenHttpHost string) {
	originalListener, err := net.Listen("tcp", ListenHttpHost)
	if err != nil {
		panic(err)
	}
	sl, err := stoppableListener.New(originalListener)
	if err != nil {
		panic(err)
	}
	server := http.Server{}
	go func() {
		server.Serve(sl)
	}()
	<-stop
	sl.Stop()
}

func StopHTTPServer() {
	log.Debug("StopHTTPServer()")
	IosLog("StopHTTPServer 0")
	go func() { stop <- true }()
	utils.Sleep(1)
	IosLog("StopHTTPServer 1")
}

func httpListener(ListenHttpHost, BrowserHttpHost string) {
	go StartHTTPServer(ListenHttpHost)
}

func tcpListener(db *utils.DCDB) {

}

func signals(countDaemons int) {

}
