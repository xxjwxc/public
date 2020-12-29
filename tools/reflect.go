package tools

import (
	"reflect"
	"runtime"
)

// GetFuncName 获取函数名
func GetFuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
