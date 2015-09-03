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

import  (
	"net/http"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)


func iosLog(text string) {
	if utils.IOS() {
		C.logNS(C.CString(text))
	}
}

func httpListener(ListenHttpHost, BrowserHttpHost string) {
	go func() {
		http.ListenAndServe(ListenHttpHost, nil)
	}()
}

func tcpListener(db *utils.DCDB) {

}

func signals(countDaemons int) {

}