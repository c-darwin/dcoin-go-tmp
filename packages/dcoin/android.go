// +build android

package dcoin

import (
	"net/http"
)

func IosLog(text string) {
}

func httpListener(ListenHttpHost, BrowserHttpHost string) {
	go func() {
		http.ListenAndServe(ListenHttpHost, nil)
	}()
}

func tcpListener() {

}

func signals(countDaemons int) {

}
