package controllers

import (
	"errors"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"archive/zip"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)


func (c *Controller) UpdateDcoin() (string, error) {

	if c.SessRestricted != 0 {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	_, url, err := c.getUpdVerAndUrl()
	if err!= nil {
		return "", utils.ErrInfo(err)
	}

	if len(url) > 0 {
		_, err := utils.DownloadToFile(url, "tmp_dc.zip", 3600, nil, nil)
		if err!= nil {
			return "", utils.ErrInfo(err)
		}
		zipfile := "tmp_dc.zip"
		reader, err := zip.OpenReader(zipfile)
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		defer reader.Close()

		f_ := reader.Reader.File
		f := f_[0]
		zipped, err := f.Open()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		defer zipped.Close()

		writer, err := os.OpenFile("tmp_dc", os.O_WRONLY|os.O_CREATE, f.Mode())
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		defer writer.Close()

		if _, err = io.Copy(writer, zipped); err != nil {
			return "", utils.ErrInfo(err)
		}

		exec.Command("tmp_dc", "-oldFileName", filepath.Base(os.Args[0]))
		return utils.JsonAnswer("success", "success").String(), nil
	}
	return "", nil
}
