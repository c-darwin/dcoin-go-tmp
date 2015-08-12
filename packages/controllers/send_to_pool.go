package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"errors"
	"fmt"
	"net"
	"io/ioutil"
	"time"
	"os"
)

func (c *Controller) SendToPool() (string, error) {

	if c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	filesSign := c.r.FormValue("filesSign")

	conn, err := net.DialTimeout("tcp", c.r.FormValue("pool"), 5 * time.Second)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(240 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(240 * time.Second))

	var data []byte
	data = append(data, utils.DecToBin(c.SessUserId, 5)...)
	data = append(data, utils.EncodeLengthPlusData(filesSign)...)

	if _, err := os.Stat(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_face.jpg"); err == nil {
		file, err := ioutil.ReadFile(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_face.jpg")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		data = append(data, utils.EncodeLengthPlusData(append(utils.DecToBin(0, 1), file...))...)
	}
	if _, err := os.Stat(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_profile.jpg"); err == nil {
		file, err := ioutil.ReadFile(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_profile.jpg")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		data = append(data, utils.EncodeLengthPlusData(append(utils.DecToBin(1, 1), file...))...)
	}
	if _, err := os.Stat(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_video.mp4"); err == nil {
		file, err := ioutil.ReadFile(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_video.mp4")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		data = append(data, utils.EncodeLengthPlusData(append(utils.DecToBin(2, 1), file...))...)
	}
	if _, err := os.Stat(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_video.webm"); err == nil {
		file, err := ioutil.ReadFile(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_video.webm")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		data = append(data, utils.EncodeLengthPlusData(append(utils.DecToBin(3, 1), file...))...)
	}
	if _, err := os.Stat(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_video.ogv"); err == nil {
		file, err := ioutil.ReadFile(*utils.Dir+"/public/"+utils.Int64ToStr(c.SessUserId)+"_user_video.ogv")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		data = append(data, utils.EncodeLengthPlusData(append(utils.DecToBin(4, 1), file...))...)
	}

	// в 4-х байтах пишем размер данных, которые пошлем далее
	size := utils.DecToBin(len(data), 4)
	_, err = conn.Write(size)
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// далее шлем сами данные
	_, err = conn.Write([]byte(data))
	if err != nil {
		return "", utils.ErrInfo(err)
	}

	// в ответ получаем статус
	buf := make([]byte, 1)
	_, err = conn.Read(buf)
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	status := utils.BinToDec(buf)
	fmt.Println(status)

	return "", nil
}
