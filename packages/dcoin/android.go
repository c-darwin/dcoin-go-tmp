// +build android

package dcoin

import  (
	"net/http"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"net"
)
/*
func httpListener(ListenHttpHost, BrowserHttpHost string) {
	go func() {
		http.ListenAndServe(ListenHttpHost, nil)
	}()
}
*/

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
		err = http.Serve(NewBoundListener(50, l), http.DefaultServeMux)
		if err != nil {
			log.Error("Error listening: %v (%v)", err, ListenHttpHost)
			panic(err)
			//os.Exit(1)
		}
	}()
}
func tcpListener(db *utils.DCDB) {

}

func signals(countDaemons int) {

}