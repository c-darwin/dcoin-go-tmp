// +build android

package dcoin

import  (
	"net/http"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

func iosLog(text string) {
}

func httpListener(ListenHttpHost, BrowserHttpHost string) {
	go func() {
		http.ListenAndServe(ListenHttpHost, nil)
		http.C
	}()
}

func tcpListener(db *utils.DCDB) {

}

func signals(countDaemons int) {

}