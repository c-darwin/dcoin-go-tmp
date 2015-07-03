package controllers
import (
	"utils"
	"log"
	"errors"
	"encoding/json"
)

func (c *Controller) GenerateNewNodeKey() (string, error) {

	if c.SessRestricted!=0 {
		return "", errors.New("Permission denied")
	}

	priv, pub := utils.GenKeys()
	json, err := json.Marshal(map[string]string{"private_key": priv, "public_key": pub})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	log.Println(json)
	return string(json), nil
}