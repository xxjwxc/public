package dev

import (
	"os"
	"strings"
)

var _isDev = true
var service = "service"
var fileHost = "file"

//是否开发模式
func OnSetDev(isDev bool) {
	_isDev = isDev
}

//IsDev ... 是否是开发版本
func IsDev() bool {
	return _isDev
}

//判断是否在测试环境下使用
func IsRunTesting() bool {
	if len(os.Args) > 1 {
		return strings.HasPrefix(os.Args[1], "-test")
	}
	return false
}

//设置服务名
func SetService(s string) {
	service = s
}

//获取服务名
func GetService() string {
	return service
}

//设置服务名
func SetFileHost(s string) {
	fileHost = s
}

//获取文件host
func GetFileHost() string {
	return fileHost
}
