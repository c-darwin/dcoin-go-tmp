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
	fmt.Println("TxMap", p.TxMap)
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) Admin1BlockFront() (error) {
	// public_key админа еще нет, он в этом блоке
	return nil
}


type firstBlock struct {
	Publickey string `json:"public_key"`
	NodePublicKey string `json:"node_public_key"`
	Host string `json:"host"`
	Currency [][]interface{} `json:"currency"`
	Variables map[string]interface{} `json:"variables"`
	SpotsCompatibility map[string]string `json:"spots_compatibility"`
}
func (p *Parser) Admin1Block() (error) {
	var firstBlock firstBlock
	err := json.Unmarshal(p.TxMap["data"], &firstBlock)
	if err != nil {
		return err
	}
	for _, currencyData := range firstBlock.Currency {
		fmt.Println(currencyData[0], currencyData[1], currencyData[2])
		currencyId, err := p.ExecSqlGetLastInsertId("INSERT INTO currency (name, full_name, max_other_currencies) VALUES (?,?,?)",
			currencyData[0], currencyData[1], currencyData[3])
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.ExecSql("INSERT INTO pct (time, currency_id, miner, user, block_id) VALUES (0,?,0,0,1)",
			currencyId)
		if err != nil {
			return p.ErrInfo(err)
		}
		err = p.ExecSql("INSERT INTO max_promised_amounts (time, currency_id, amount, block_id) VALUES (0,?,?,1)",
			currencyId, currencyData[2])
		if err != nil {
			return p.ErrInfo(err)
		}
	}

	for name, value := range firstBlock.Variables {
		err := p.ExecSql("INSERT INTO variables (name, value) VALUES (?,?)",
			name, value)
		if err != nil {
			return p.ErrInfo(err)
		}
	}

	err = p.ExecSql(`INSERT INTO miners_data (user_id, miner_id, status, node_public_key, host, photo_block_id, photo_max_miner_id, miners_keepers)
		VALUES (1,1,'miner',[hex],?,1,1,1)`,
		firstBlock.NodePublicKey, firstBlock.Host)
	if err != nil {
		return p.ErrInfo(err)
	}

	err = p.ExecSql(`INSERT INTO users (public_key_0) VALUES ([hex])`, firstBlock.Publickey)
	if err != nil {
		return p.ErrInfo(err)
	}

	err = p.ExecSql(`INSERT INTO miners (miner_id, active) VALUES (1,1)`)
	if err != nil {
		return p.ErrInfo(err)
	}

	err = p.ExecSql(`INSERT INTO spots_compatibility (version, example_spots, compatibility, segments, tolerances) VALUES (?,?,?,?,?)`,
		firstBlock.SpotsCompatibility["version"], firstBlock.SpotsCompatibility["example_spots"], firstBlock.SpotsCompatibility["compatibility"], firstBlock.SpotsCompatibility["segments"], firstBlock.SpotsCompatibility["tolerances"])
	if err != nil {
		return p.ErrInfo(err)
	}

	return nil
}

func (p *Parser) Admin1BlockRollback() (error) {
	return nil
}
