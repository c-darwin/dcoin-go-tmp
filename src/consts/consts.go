package consts
import (
	"fmt"
)

// У скольких нодов должен быть такой же блок как и у нас, чтобы считать, что блок у большей части DC-сети. для get_confirmed_block_id()
const MIN_CONFIRMED_NODES = 3

var LangMap map[string]int

func init() {
	LangMap = map[string]int{"en":1, "ru":42}
	fmt.Println(LangMap)
}
