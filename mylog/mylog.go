package mylog

var _log myLog

func init() {
	_log = GetDefaultStd()
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
}

// SetLog set log
func SetLog(log myLog) {
	_log = log
}

// Info info level debug
func Info(a ...interface{}) {
	_log.Info(a...)
}

// Infof info level debug
func Infof(msg string, a ...interface{}) {
	_log.Infof(msg, a...)
}

// Error 记录错误信息
func Error(a ...interface{}) {
	_log.Error(a...)
}

// Errorf 记录错误信息
func Errorf(msg string, a ...interface{}) {
	_log.Errorf(msg, a...)
}

// Debug level debug
func Debug(a ...interface{}) {
	_log.Debug(a...)
}

// Debugf level debug
func Debugf(msg string, a ...interface{}) {
	_log.Debugf(msg, a...)
}

// Fatal level fatal
func Fatal(a ...interface{}) {
	_log.Fatal(a...)
}

// Fatalf level fatal
func Fatalf(msg string, a ...interface{}) {
	_log.Fatal(a...)
}

// JSON json info
func JSON(a ...interface{}) {
	_log.JSON(a...)
}

// TraceError return trace of error
func TraceError(err error) error {
	return _log.TraceError(err)
}

// SaveError return trace of error
func SaveError(errstring, flag string) {
	_log.SaveError(errstring, flag)
}

// Panic return trace of error
func Panic(a ...interface{}) {
	_log.Panic(a...)
}

// ErrorString 打印错误信息
func ErrorString(a ...interface{}) {
	_log.ErrorString(a...)
}
