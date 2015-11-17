package controllers

import (
	"errors"
	"github.com/c-darwin/dcoin-go-tmp/packages/utils"
	"os"
	"os/exec"
	"fmt"
	"path/filepath"
	"github.com/kardianos/osext"
	"archive/zip"
	"io"
)


func (c *Controller) UpdateDcoin() (string, error) {

	if c.SessRestricted != 0 || !c.NodeAdmin {
		return "", utils.ErrInfo(errors.New("Permission denied"))
	}

	_, url, err := c.getUpdVerAndUrl()
	if err!= nil {
		return "", utils.ErrInfo(err)
	}

	fmt.Println(url)
	if len(url) > 0 {
		_, err := utils.DownloadToFile(url, *utils.Dir+"/dc.zip", 3600, nil, nil)
		if err!= nil {
			return "", utils.ErrInfo(err)
		}
		zipfile := *utils.Dir+"/dc.zip"
		fmt.Println(zipfile)
		reader, err := zip.OpenReader(zipfile)
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		f_ := reader.Reader.File
		f := f_[0]
		zipped, err := f.Open()
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		writer, err := os.OpenFile(*utils.Dir+"/dc.tmp", os.O_WRONLY|os.O_CREATE, f.Mode())
		if err != nil {
			return "", utils.ErrInfo(err)
		}

		if _, err = io.Copy(writer, zipped); err != nil {
			return "", utils.ErrInfo(err)
		}
		reader.Close()
		zipped.Close()
		writer.Close()

		pwd, err := os.Getwd()
		if err != nil {
			return "", utils.ErrInfo(err)
		}
		fmt.Print(pwd)

		folderPath, err := osext.ExecutableFolder()
		if err != nil {
			log.Fatal(err)
		}

		old := ""
		fmt.Println(0, os.Args[0])
		fmt.Println(1, folderPath+"/"+filepath.Base(os.Args[0]))
		fmt.Println(2, *utils.Dir+"/"+filepath.Base(os.Args[0]))
		if _, err := os.Stat(os.Args[0]); err == nil {
			old = os.Args[0]
		} else if _, err := os.Stat(folderPath+"/"+filepath.Base(os.Args[0])); err == nil {
			old = folderPath+"/"+filepath.Base(os.Args[0])
		} else {
			old = *utils.Dir+"/"+filepath.Base(os.Args[0])
		}

		/*file, _ := os.Open(*utils.Dir+"/dc")
		defer file.Close()
		stat, _ := file.Stat()
		fmt.Println("File size is " , stat.Size())*/

		fmt.Println(*utils.Dir+"/dc.tmp", "-oldFileName", old, "-dir", *utils.Dir, "-logLevel", "DEBUG")
		/*var cmdOut []byte
		if cmdOut, err = exec.Command(*utils.Dir+"/dc", "-oldFileName", pwd+"/"+filepath.Base(os.Args[0]), "-dir", *utils.Dir, "-logLevel", "DEBUG").Output(); err != nil {
			fmt.Println(os.Stderr)
			return "", utils.ErrInfo(err)
		}
		fmt.Println(cmdOut)*/
		err = exec.Command(*utils.Dir+"/dc.tmp", "-oldFileName", old, "-dir", *utils.Dir, "-logLevel", "DEBUG").Start()
		if err != nil {
			fmt.Println(os.Stderr)
			return "", utils.ErrInfo(err)
		}
/*
		fmt.Println(os.Args[0])
		fmt.Println(*utils.Dir+"/dc", "-oldFileName="+filepath.Base(os.Args[0]))
		cmd := exec.Command(*utils.Dir+"/dc", "-oldFileName", os.Args[0], "-dir", *utils.Dir)
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		fmt.Printf("%s\n", out.String())
		if err != nil {
			fmt.Println(err)
			return "", utils.ErrInfo(err)
		}*/
		return utils.JsonAnswer("success", "success").String(), nil
	}
	return "", nil
}
