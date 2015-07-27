package tcpserver

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"net"
	"runtime"
	"fmt"
	"github.com/op/go-logging"
	"sync"
)

var log = logging.MustGetLogger("tcpserver")
var counter int64

type TcpServer struct {
	*utils.DCDB
	Conn net.Conn
	variables *utils.Variables
}

func (t *TcpServer) deferClose() {
	t.Conn.Close()
	var mutex = &sync.Mutex{}
	mutex.Lock()
	counter--
	fmt.Println(counter)
	mutex.Unlock()
}

func (t *TcpServer) HandleTcpRequest() {

	var mutex = &sync.Mutex{}

	fmt.Println("NumCPU:", runtime.NumCPU(),
		" NumCgoCall:", runtime.NumCgoCall(),
		" NumGoRoutine:", runtime.NumGoroutine(),
		" t.counter:", counter)

	var err error

	log.Debug("HandleTcpRequest from %v", t.Conn.RemoteAddr())
	defer t.deferClose()

	mutex.Lock()
	if counter > 20 {
		return
	} else {
		counter++
		fmt.Println(counter)
	}
	mutex.Unlock()

	t.variables, err = t.GetAllVariables()
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		return
	}

	// тип данных
	buf := make([]byte, 1)
	_, err = t.Conn.Read(buf)
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		return
	}
	dataType := utils.BinToDec(buf)
	log.Debug("dataType %v", dataType)
	switch dataType {
	case 1:
		t.Type1()
	case 2:
		t.Type2()
	case 3:
		t.Type3()
	case 4 :
		t.Type4()
	case 5:
		t.Type5()
	case 6:
		t.Type6()
	case 7:
		t.Type7()
	case 8:
		t.Type8()
	case 9:
		t.Type9()
	case 10:
		t.Type10()
	}
	log.Debug("END")
}
