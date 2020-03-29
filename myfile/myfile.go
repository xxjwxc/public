package myfile

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/xxjwxc/public/tools"
)

//上传单个文件
func UploadOneFile(r *http.Request, field, file_type, dir string) (result bool, file_name string) {
	//接受post请求
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file_name = getFileName(file_type)
		//开始存储文件
		{
			file, _, err := r.FormFile(field) //文件name
			defer file.Close()
			if err != nil {
				result = false
			}

			if !tools.CheckFileIsExist(tools.GetModelPath() + "/file/" + dir + "/") {
				err1 := os.Mkdir(tools.GetModelPath()+"/file/"+dir+"/", os.ModePerm) //创建文件夹
				if err1 != nil {
					result = false
				}
			}

			f, err := os.OpenFile(tools.GetModelPath()+"/file/"+dir+"/"+file_name, os.O_WRONLY|os.O_CREATE, 0666)
			defer f.Close()
			if err != nil {
				result = false
			}
			io.Copy(f, file)
			result = true
		}

	} else {
		result = false
	}
	return
}

//多文件上传
func UploadMoreFile(r *http.Request, field, file_type, dir string) (result bool, optionDirs []string) {
	//接受post请求
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		files := r.MultipartForm.File[field]
		l := len(files)
		//		optionDirs := make([]string, l)
		for i := 0; i < l; i++ {
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				result = false
			}
			file_name := getFileName(file_type)

			if !tools.CheckFileIsExist(tools.GetModelPath() + "/file/" + dir + "/") {
				err1 := os.Mkdir(tools.GetModelPath()+"/file/"+dir+"/", os.ModePerm) //创建文件夹
				if err1 != nil {
					result = false
				}
			}
			f, err := os.OpenFile(tools.GetModelPath()+"/file/"+dir+"/"+file_name, os.O_WRONLY|os.O_CREATE, 0666)
			defer f.Close()
			if err != nil {
				result = false
			}
			io.Copy(f, file)
			optionDirs = append(optionDirs, file_name)
			result = true
		}

	} else {
		result = false
	}
	return
}
func getFileName(exp string) string {
	return fmt.Sprintf("%d%s.%s", tools.GetUtcTime(time.Now()), tools.GetRandomString(4), exp)
}
