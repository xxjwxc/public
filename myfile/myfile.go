package myfile

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/xxjwxc/public/tools"
)

// GetExp 获取字符串后缀
func GetExp(ext string) string {
	return strings.TrimLeft(path.Ext(ext), ".")
}

func getFileName(exp string) string {
	return fmt.Sprintf("%d%s.%s", tools.GetUtcTime(time.Now()), tools.GetRandomString(4), exp)
}
