package file

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"time"

	"github.com/xxjwxc/public/tools"
)

// 上传单个文件
func UploadFile(r *http.Request, w http.ResponseWriter, field, file_type string) (result bool, file_name string) {
	//接受post请求
	if r.Method != "POST" {
		return false, ""
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		return false, ""
	}
	file_name = GetFileName(file_type)

	//开始存储文件
	file, _, err := r.FormFile(field) //文件name
	if err != nil {
		return false, ""
	}
	defer file.Close()

	dir := tools.GetModelPath() + "/file/" + field + "/"
	if !tools.CheckFileIsExist(dir) {
		if err = os.Mkdir(dir, os.ModePerm); err != nil { //创建文件夹
			return false, ""
		}
	}

	// O_TRUNC 保证重复文件名不会残留旧内容;原实现漏了截断,会写出脏数据
	f, err := os.OpenFile(dir+file_name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return false, ""
	}
	defer f.Close()

	if _, err = io.Copy(f, file); err != nil {
		return false, ""
	}
	// 原实现所有分支都只把 result 置 false,从未置 true,导致成功也返回失败
	return true, file_name
}

func GetFileName(exp string) string {
	return fmt.Sprintf("%d%s.%s", tools.GetUtcTime(time.Now()), tools.GetRandomString(4), exp)
}
