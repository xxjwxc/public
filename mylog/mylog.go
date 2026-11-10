package mylog

import "sync/atomic"

// _logPtr 保存当前日志实现。
// 使用 atomic.Pointer 而非普通包级变量,避免 SetLog 与日志写调用之间产生数据竞争。
var _logPtr atomic.Pointer[myLog]

func init() {
	l := myLog(GetDefaultStd())
	_logPtr.Store(&l)
}

type myLog interface {
	Debug(a ...interface{})
	Debugf(msg string, a ...interface{})
	Info(a ...interface{})
	Infof(msg string, a ...interface{})
	Error(a ...interface{})
	Errorf(msg string, a ...interface{})
	ErrorString(a ...interface{})        // ErrorString 打印错误信息(不记录错误文件)
	Fatal(v ...interface{})              // 系统级错误,退出
	Fatalf(msg string, a ...interface{}) // 系统级错误,退出
	SaveError(errstring, flag string)
	Panic(a ...interface{})

	JSON(a ...interface{})

	TraceError(err error) error
	Close()
	Printf(string, ...interface{})
}

// GetLog 获取当前日志实现
func GetLog() myLog {
	if p := _logPtr.Load(); p != nil {
		return *p
	}
	return GetDefaultStd()
}

// SetLog set log
func SetLog(log myLog) {
	_logPtr.Store(&log)
}

// Info info level debug
func Info(a ...interface{}) {
	GetLog().Info(a...)
}

// Infof info level debug
func Infof(msg string, a ...interface{}) {
	GetLog().Infof(msg, a...)
}

// Error 记录错误信息
func Error(a ...interface{}) {
	GetLog().Error(a...)
}

// Errorf 记录错误信息
func Errorf(msg string, a ...interface{}) {
	GetLog().Errorf(msg, a...)
}

// Debug level debug
func Debug(a ...interface{}) {
	GetLog().Debug(a...)
}

// Debugf level debug
func Debugf(msg string, a ...interface{}) {
	GetLog().Debugf(msg, a...)
}

// Fatal level fatal
func Fatal(a ...interface{}) {
	GetLog().Fatal(a...)
}

// Fatalf level fatal
func Fatalf(msg string, a ...interface{}) {
	// 原实现误传 a...,丢失了格式串 msg
	GetLog().Fatalf(msg, a...)
}

// JSON json info
func JSON(a ...interface{}) {
	GetLog().JSON(a...)
}

// TraceError return trace of error
func TraceError(err error) error {
	return GetLog().TraceError(err)
}

// SaveError return trace of error
func SaveError(errstring, flag string) {
	GetLog().SaveError(errstring, flag)
}

// Panic return trace of error
func Panic(a ...interface{}) {
	GetLog().Panic(a...)
}

// ErrorString 打印错误信息
func ErrorString(a ...interface{}) {
	GetLog().ErrorString(a...)
}
