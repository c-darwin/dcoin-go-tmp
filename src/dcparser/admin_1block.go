package dcparser

import (
	"fmt"
	//"utils"
	"encoding/json"
)

func (p *Parser) Admin1BlockInit() (error) {
	fields := []string {"data", "sign"}
	TxMap := make(map[string][]byte)
	TxMap, err := p.GetTxMap(fields);
	p.TxMap = TxMap;
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) Admin1BlockFront() (error) {
	// public_key админа еще нет, он в этом блоке
	return nil
}

func (p *Parser) Admin1Block() (error) {
	var data []byte
	err := json.Unmarshal(p.TxMap["data"], &data)
	if err!=nil {
		return err
	}
	fmt.Println("data", data)
	return nil
}

func (p *Parser) Admin1BlockRollback() (error) {
	return nil
}
