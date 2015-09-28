package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"

	"strings"
	"io/ioutil"
	"fmt"
	"os"
)

type upgrade7Page struct {
	SignData string
	ShowSignData bool
	TxType string
	TxTypeId int64
	TimeNow int64
	UserId int64
	Alert string
	Lang map[string]string
	CountSignArr []int
	Latitude string
	Longitude string
	NodePublicKey string
	SaveAndGotoStep string
	UpgradeMenu string
	ProfileHash string
	FaceHash string
	MyTable map[string]string
	NoExistsMp4 bool
	Data map[string]string
	Mobile bool
}

func (c *Controller) Upgrade7() (string, error) {

	log.Debug("Upgrade7")

	txType := "NewMiner"
	txTypeId := utils.TypeInt(txType)
	timeNow := utils.Time()

	// Формируем контент для подписи
	myTable, err := c.OneRow("SELECT user_id, race, country, geolocation, http_host, tcp_host, face_coords, profile_coords, video_url_id, video_type FROM "+c.MyPrefix+"my_table").String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	if len(myTable["video_url_id"]) == 0 {
		myTable["video_url_id"] = "null"
	}
	if len(myTable["video_type"]) == 0 {
		myTable["video_type"] = "null"
	}
	var profileHash, faceHash string

	if _, err := os.Stat(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_face.jpg"); err == nil {
		file, err := ioutil.ReadFile(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_face.jpg")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		faceHash = string(utils.DSha256(file))
	}
	if _, err := os.Stat(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_profile.jpg"); err == nil {
		file, err := ioutil.ReadFile(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_profile.jpg")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		profileHash = string(utils.DSha256(file))
	}

	latitude := "0"
	longitude := "0"
	if len(myTable["geolocation"]) > 0 {
		x := strings.Split(myTable["geolocation"], ", ")
		latitude = x[0]
		longitude = x[1]
	}

	// проверим, есть ли не обработанные ключи в локальной табле
	nodePublicKey, err := c.Single("SELECT public_key FROM "+c.MyPrefix+"my_node_keys WHERE block_id  =  0").String()
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	if len(nodePublicKey) == 0 {
		//  сгенерим ключ для нода
		priv, pub := utils.GenKeys()
		err = c.ExecSql("INSERT INTO "+c.MyPrefix+"my_node_keys ( public_key, private_key ) VALUES ( [hex], ? )", pub, priv)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
	} else {
		nodePublicKey = string(utils.BinToHex([]byte(nodePublicKey)))
	}

	saveAndGotoStep := strings.Replace(c.Lang["save_and_goto_step"], "[num]", "7", -1)
	upgradeMenu := utils.MakeUpgradeMenu(7)

	var noExistsMp4 bool
	if _, err := os.Stat(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_video.mp4"); os.IsNotExist(err) {
		noExistsMp4 = true
	}

	TemplateStr, err := makeTemplate("upgrade_7", "upgrade7", &upgrade7Page {
		Alert: c.Alert,
		Lang: c.Lang,
		CountSignArr: c.CountSignArr,
		ShowSignData: c.ShowSignData,
		UserId: c.SessUserId,
		TimeNow: timeNow,
		TxType: txType,
		TxTypeId: txTypeId,
		NoExistsMp4: noExistsMp4,
		SignData: fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v", txTypeId, timeNow, c.SessUserId, myTable["race"],myTable["country"],latitude,longitude,myTable["http_host"],myTable["tcp_host"],faceHash,profileHash,myTable["face_coords"],myTable["profile_coords"],myTable["video_type"],myTable["video_url_id"],nodePublicKey),
		SaveAndGotoStep: saveAndGotoStep,
		UpgradeMenu: upgradeMenu,
		Latitude: latitude,
		Longitude: longitude,
		NodePublicKey: nodePublicKey,
		ProfileHash: profileHash,
		FaceHash: faceHash,
		Data: myTable,
		MyTable: myTable,
		Mobile: utils.Mobile()})
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	return TemplateStr, nil
}