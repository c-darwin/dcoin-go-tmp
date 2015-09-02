// +build android ios

package dcoin

import  (
	"net/http"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
)

func httpListener(ListenHttpHost, BrowserHttpHost string) {
	go func() {
		http.ListenAndServe(ListenHttpHost, nil)
	}()
}

func tcpListener(db *utils.DCDB) {

}

func signals(countDaemons int) {

}