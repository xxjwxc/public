package myfile

import (
	"fmt"
	"path"
	"time"

	"github.com/xxjwxc/public/tools"
)

// GetExp 获取字符串后缀
func GetExp(exp string) string {
	return path.Ext(exp) //获取文件后缀
}

func getFileName(exp string) string {
	return fmt.Sprintf("%d%s.%s", tools.GetUtcTime(time.Now()), tools.GetRandomString(4), exp)
}
