package dev

import (
	"os"
	"strings"
)

var _isDev = true

//是否开发模式
func OnSetDev(isDev bool) {
	_isDev = isDev
}

//OnIsDev ... 是否是开发版本
func OnIsDev() bool {
	return _isDev
}

//判断是否在测试环境下使用
func IsRunTesting() bool {
	if len(os.Args) > 1 {
		return strings.HasPrefix(os.Args[1], "-test")
	}
	return false
}
