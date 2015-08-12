package tcpserver

import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"io/ioutil"
	"os"
)

func (t *TcpServer) Type11() {

	/* Получаем данные от send_to_pool */

	// размер данных
	buf := make([]byte, 4)
	_, err := t.Conn.Read(buf)
	if err != nil {
		log.Error("%v", utils.ErrInfo(err))
		return
	}
	size := utils.BinToDec(buf)
	if size < 32<<20 {
		// сами данные
		binaryData := make([]byte, size)
		binaryData, err = ioutil.ReadAll(t.Conn)
		if err != nil {
			log.Error("%v", utils.ErrInfo(err))
			return
		}
		userId := utils.BinToDec(utils.BytesShift(&binaryData, 5))
		// проверим, есть ли такой юзер на пуле
		inPool, err := t.Single(`SELECT user_id FROM community WHERE user_id=?`, userId).Int64()
		if inPool<=0 {
			log.Error("%v", utils.ErrInfo("inPool<=0"))
			_, err = t.Conn.Write(utils.DecToBin(0, 1))
			return
		}
		filesSign := utils.BytesShift(&binaryData, utils.DecodeLength(&binaryData))
		forSign := ""
		var files []string
		for i:=0; i < 5; i++ {
			size := utils.DecodeLength(&binaryData)
			data := utils.BytesShift(&binaryData, size)
			if len(binaryData)==0 {
				break
			}
			fileType := utils.BinToDec(utils.BytesShift(&data, 1))
			var name string
			switch fileType {
			case 0:
				name = utils.Int64ToStr(userId)+"_user_face.jpg"
			case 1:
				name = utils.Int64ToStr(userId)+"_user_profile.jpg"
			case 2:
				name = utils.Int64ToStr(userId)+"_user_video.mp4"
			case 3:
				name = utils.Int64ToStr(userId)+"_user_video.webm"
			case 4:
				name = utils.Int64ToStr(userId)+"_user_video.ogv"
			}
			forSign = forSign+","+string(data)
			err = ioutil.WriteFile(os.TempDir()+"/"+name, data, 0644)
			if err != nil {
				log.Error("%v", utils.ErrInfo(err))
				return
			}
			files = append(files, name)
		}

		forSign = forSign[:len(forSign)-1]

		// проверим подпись
		publicKey, err := t.GetUserPublicKey(userId)
		resultCheckSign, err := utils.CheckSign([][]byte{[]byte(publicKey)}, forSign, utils.HexToBin(filesSign), true);
		if err != nil {
			log.Error("%v", utils.ErrInfo(err))
			_, err = t.Conn.Write(utils.DecToBin(0, 1))
			return
		}
		if resultCheckSign {
			for i:=0; i < len(files); i++ {
				utils.CopyFileContents(os.TempDir()+"/"+files[i], *utils.Dir+"/public/"+files[i])
			}
		}

		// и возвращаем статус
		_, err = t.Conn.Write(utils.DecToBin(1, 1))
		if err != nil {
			log.Error("%v", utils.ErrInfo(err))
			return
		}
	}
}