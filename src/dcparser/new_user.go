package dcparser

import (
	"fmt"
	"utils"
	"encoding/json"
	"regexp"
)


func (p *Parser) NewUser() (error) {
	fields := []string {"public_key", "sign"}
	TxMap := make(map[string][]byte)
	TxMap, err := p.GetTxMap(fields);
	p.TxMap = TxMap;
	fmt.Println("TxMap", p.TxMap)
	if err != nil {
		return err
	}
	return nil
}

func (p *Parser) NewUserFront() (error) {
	err := p.generalCheck()
	if err != nil {
		return err
	}

	// является ли данный юзер майнером
	err = p.checkMiner(p.TxMap["user_id"])
	if err != nil {
		return err
	}

	// прошло ли 30 дней с момента регистрации майнера
	err = p.checkMinerNewbie()
	if err != nil {
		return err
	}

	// чтобы не записали слишком мелкий или слишком крупный ключ
	if !utils.CheckInputData(p.TxMap["public_key_hex"], "public_key") {
		return utils.ErrInfoFmt("incorrect public_key")
	}

	// публичный ключ должен быть без паролей
	if ok, _ := regexp.MatchString("DEK-Info: (.+),(.+)", string(p.TxMap["public_key"])); ok{
		return utils.ErrInfoFmt("incorrect public_key")
	}

	forSign := fmt.Sprintf("%v,%v,%v,%v", p.TxMap["type"], p.TxMap["time"], p.TxMap["user_id"], p.TxMap["public_key_hex"])
	err = utils.CheckSign(p.PublicKeys, forSign, p.TxMap["sign"], false);
	if err != nil {
		return err
	}

	// один ключ не может быть у двух юзеров
	num_, err = p.DCDB.Single("SELECT count(user_id) FROM users WHERE public_key_0 = [hex] OR public_key_1 = [hex] OR public_key_2 = [hex]",
		p.TxMap["public_key_hex"], p.TxMap["public_key_hex"], p.TxMap["public_key_hex"])
	num := utils.StrToInt64(num_)
	if num > 0 {
		return utils.ErrInfoFmt("exists public_ke")
	}
	err = p.getAdminUserId()
	if err != nil {
		return err
	}
	if utils.StrToInt64(p.TxMap["user_id"]) == p.AdminUserId {
		err =
	}
	/*
		if ($this->tx_data['user_id'] == $this->admin_user_id)
			$error = $this -> limit_requests( 1000, 'new_user', 86400 );
		else
			$error = $this -> limit_requests( $this->variables['limit_new_user'], 'new_user', $this->variables['limit_new_user_period'] );
		if ($error)
			return $error;
*/
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
func (p *Parser) NewUser() (error) {
	var firstBlock firstBlock
	err := json.Unmarshal(p.TxMap["data"], &firstBlock)
	if err != nil {
		return err
	}
	for _, currencyData := range firstBlock.Currency {
		fmt.Println(currencyData[0], currencyData[1], currencyData[2])
		currencyId, err := p.DCDB.ExecSqlGetLastInsertId("INSERT INTO currency (name, full_name, max_other_currencies) VALUES (?,?,?)",
			currencyData[0], currencyData[1], currencyData[3])
		if err != nil {
			return utils.ErrInfo(err)
		}
		_, err = p.DCDB.ExecSql("INSERT INTO pct (time, currency_id, miner, user, block_id) VALUES (0,?,0,0,1)",
			currencyId)
		if err != nil {
			return utils.ErrInfo(err)
		}
		_, err = p.DCDB.ExecSql("INSERT INTO max_promised_amounts (time, currency_id, amount, block_id) VALUES (0,?,?,1)",
			currencyId, currencyData[2])
		if err != nil {
			return utils.ErrInfo(err)
		}
	}

	for name, value := range firstBlock.Variables {
		_, err := p.DCDB.ExecSql("INSERT INTO variables (name, value) VALUES (?,?)",
			name, value)
		if err != nil {
			return utils.ErrInfo(err)
		}
	}

	_, err = p.DCDB.ExecSql(`INSERT INTO miners_data (user_id, miner_id, status, node_public_key, host, photo_block_id, photo_max_miner_id, miners_keepers)
		VALUES (1,1,'miner',[hex],?,1,1,1)`,
		firstBlock.NodePublicKey, firstBlock.Host)
	if err != nil {
		return utils.ErrInfo(err)
	}

	_, err = p.DCDB.ExecSql(`INSERT INTO users (public_key_0) VALUES ([hex])`,
		firstBlock.Publickey)
	if err != nil {
		return utils.ErrInfo(err)
	}

	_, err = p.DCDB.ExecSql(`INSERT INTO miners (miner_id, active) VALUES (1,1)`)
	if err != nil {
		return utils.ErrInfo(err)
	}

	_, err = p.DCDB.ExecSql(`INSERT INTO spots_compatibility (version, example_spots, compatibility, segments, tolerances) VALUES (?,?,?,?,?)`,
		firstBlock.SpotsCompatibility["version"], firstBlock.SpotsCompatibility["example_spots"], firstBlock.SpotsCompatibility["compatibility"], firstBlock.SpotsCompatibility["segments"], firstBlock.SpotsCompatibility["tolerances"])
	if err != nil {
		return utils.ErrInfo(err)
	}

	return nil
}

func (p *Parser) NewUserRollback() (error) {
	return nil
}
