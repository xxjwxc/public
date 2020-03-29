package myhttp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/xxjwxc/public/dev"

	"time"

	"github.com/xxjwxc/public/mylog"
	"github.com/xxjwxc/public/tools"
)

//UploadMoreFile 多文件上传,dir:空则使用文件后缀做dir
func UploadMoreFile(r *http.Request, dir string) (result bool, optionDirs []string) {
	//接受post请求
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		if r.MultipartForm == nil {
			result = false
		} else {
			for _, files := range r.MultipartForm.File {
				for _, v := range files {
					var _dir = dir
					file, _ := v.Open()
					defer file.Close()
					ext := getFileType(v.Filename)
					if len(ext) == 0 {
						continue
					}
					if len(_dir) == 0 {
						_dir = ext
					}

					absDir := tools.GetCurrentDirectory() + "/" + dev.GetFileHost() + "/" + _dir + "/"
					fileName := getFileName(ext)
					if !tools.CheckFileIsExist(absDir) {
						tools.BuildDir(absDir)
					}

					//存在则覆盖
					f, err := os.OpenFile(absDir+fileName,
						os.O_WRONLY|os.O_CREATE, 0666)
					defer f.Close()
					if err != nil {
						mylog.Error(err)
						result = false
						return
					}

					io.Copy(f, file)
					optionDirs = append(optionDirs, "/"+dev.GetService()+"/"+dev.GetFileHost()+"/"+_dir+"/"+fileName)
					result = true
				}
			}
		}
	} else {
		result = false
	}
	return
}

func getFileName(exp string) string {
	return fmt.Sprintf("%d%s.%s", tools.GetUtcTime(time.Now()), tools.GetRandomString(4), exp)
}

//获取文件后缀
func getFileType(exp string) string {
	fileSuffix := path.Ext(exp) //获取文件后缀
	if len(fileSuffix) > 1 {
		return fileSuffix[1:]
	}
	return ""
}

//PostFile 模拟客戶端文件上传
//fieldname注意与服务器端保持一致
func PostFile(filename, fieldname, targetURL string) (result string, e error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile(fieldname, filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		e = err
		return
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		e = err
		return
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		e = err
		return
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetURL, contentType, bodyBuf)
	if err != nil {
		e = err
		return
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e = err
		return
	}
	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
	result = string(respBody)
	return
}
