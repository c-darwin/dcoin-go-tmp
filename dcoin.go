package main
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/dcoin"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"fmt"
	"github.com/c-darwin/dcoin-go-tmp/packages/daemons"
)

func init() {

}

func main() {
	dcoin.Start("")

	// мониторинг демонов
	var daemonsTable map[string]string
	go func() {
		fmt.Println("daemonsTable")
		for {
			fmt.Println("wait daemonNameAndTime")
			daemonNameAndTime := <-daemons.MonitorDaemonCh
			fmt.Printf("daemonNameAndTime: %v\n", daemonNameAndTime)
			daemonsTable[daemonNameAndTime[0]] = daemonNameAndTime[1]
			if utils.Time()%100 == 0 {
				fmt.Printf("daemonsTable: %v\n", daemonsTable)
			}
		}
		fmt.Println("end daemonsTable")

	} ()
}
