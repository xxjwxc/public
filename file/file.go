package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"public/tools"
	"time"
)

//上传单个文件
func UploadFile(r *http.Request, w http.ResponseWriter, field, file_type string) (result bool, file_name string) {
	//接受post请求
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		file_name = GetFileName(file_type)
		//开始存储文件
		{
			file, _, err := r.FormFile(field) //文件name
			defer file.Close()
			if err != nil {
				result = false
			}

			if !tools.CheckFileIsExist(tools.GetModelPath() + "/file/" + field + "/") {
				err1 := os.Mkdir(tools.GetModelPath()+"/file/"+field+"/", os.ModePerm) //创建文件夹
				if err1 != nil {
					result = false
				}
			}

			f, err := os.OpenFile(tools.GetModelPath()+"/file/"+field+"/"+file_name, os.O_WRONLY|os.O_CREATE, 0666)
			defer f.Close()
			if err != nil {
				result = false
			}
			io.Copy(f, file)
		}

	} else {
		result = false
	}
	return
}

func GetFileName(exp string) string {
	return fmt.Sprintf("%d%s.%s", tools.GetUtcTime(time.Now()), tools.GetRandomString(4), exp)
}
