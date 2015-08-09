package controllers
import (
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"bytes"
	"io"
	"io/ioutil"
)

func (c *Controller) UploadVideo() (string, error) {

	if c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	c.r.ParseMultipartForm(32 << 20)
	file, _, err := c.r.FormFile("file")
	if err != nil {
		return "", utils.ErrInfo(err)
	}
	fmt.Println(c.r.MultipartForm.File["file"][0].Filename)
	fmt.Println(c.r.MultipartForm.File["file"][0].Header.Get("Content-Type"))
	fmt.Println(c.r.MultipartForm.Value["type"][0])

	var contentType, videoType string
	if _, ok := c.r.MultipartForm.File["file"]; ok {
		contentType = c.r.MultipartForm.File["file"][0].Header.Get("Content-Type")
	}
	if _, ok := c.r.MultipartForm.Value["type"]; ok {
		videoType = c.r.MultipartForm.Value["type"][0]
	}
	end := "mp4"
  	switch contentType {
	case "video/mp4", "video/quicktime":
		end = "mp4"
	case "video/ogg":
		end = "ogv"
	case "video/webm":
		end = "webm"
	case "video/3gpp":
		fmt.Println("3gpp")
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("3gp", "3gp.3gp")
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		_, err = io.Copy(part, file)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		err = writer.Close()
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		fmt.Println("http://3gp.dcoin.me")

		req, err := http.NewRequest("POST", "http://3gp.dcoin.me", body)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		log.Debug("%v", req)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		htmlData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		fmt.Println("htmlData",string(htmlData))

	}
	log.Debug(videoType, end)
			/*
				buffer := new(bytes.Buffer)
				_, err = io.Copy(buffer, file)
				if err != nil {
					return "", utils.ErrInfo(err)
				}
				defer file.Close()
				log.Debug(buffer.String())

				var mainMap JsonBackup
				err = json.Unmarshal(buffer.Bytes(), &mainMap)
				if err != nil {
					return "", utils.ErrInfo(err)
				}

				schema_ := &schema.SchemaStruct{}
				schema_.DCDB = c.DCDB
				schema_.DbType = c.ConfigIni["db_type"]
				for i:=0; i <len(mainMap.Community); i++ {
					schema_.PrefixUserId = utils.StrToInt(mainMap.Community[i])
					schema_.GetSchema()
					c.ExecSql(`INSERT INTO community (user_id) VALUES (?)`, mainMap.Community[i])
				}

				allTables, err := c.GetAllTables()

				for table, arr := range mainMap.Data {
					if !utils.InSliceString(table, allTables) {
						continue
					}
					//_ = c.ExecSql(`DROP TABLE `+table)
					//if err != nil {
					//	return "", utils.ErrInfo(err)
					//}
					log.Debug(table)
					for i, data := range arr {
						log.Debug("%v", i)
						colNames := ""
						values := []interface {} {}
						qq := ""
						for name, value := range data {

							if ok, _ := regexp.MatchString("my_table", table); ok{
								if name == "host" {
									name = "http_host"
								}
							}
							if name == "show_progressbar" {
								name = "show_progress_bar"
							}

							colNames += name+","
							values = append(values, value)
							if ok, _ := regexp.MatchString("(hash_code|public_key|encrypted)", name); ok{
								qq+="[hex],"
							} else {
								qq+="?,"
							}
						}
						colNames = colNames[0:len(colNames)-1]
						qq = qq[0:len(qq)-1]
						query := `INSERT INTO `+table+` (`+colNames+`) VALUES (`+qq+`)`
						log.Debug("%v", query)
						log.Debug("%v", values)
						err = c.ExecSql(query, values...)
						if err != nil {
							return "", utils.ErrInfo(err)
						}
					}
				}
			*/

	return "", nil
}
