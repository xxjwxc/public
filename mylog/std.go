package mylog

import (
	"fmt"
	"log"
	"os"

	"github.com/xxjwxc/public/dev"

	"github.com/gookit/color"
	"github.com/xxjwxc/public/errors"
)

type stdLog struct {
	errDeal
}

const ( //
	logError   = iota //打印 Error 及以上级别
	logwarning        //打印 warning 及以上级别
	logInfo           //默认的返回值，为0，自增 //打印 Info 及以上级别
)

// GetDefaultStd get default std logger
func GetDefaultStd() *stdLog {
	return &stdLog{}
}

// Info level info msg
func (s *stdLog) Info(a ...interface{}) {
	log.Println(color.FgGreen.Render(a...))
}

// Info level info msg
func (s *stdLog) Infof(msg string, a ...interface{}) {
	log.Println(color.FgGreen.Render(fmt.Sprintf(msg, a...)))
}

// Error 记录错误信息
func (s *stdLog) Error(a ...interface{}) {
	// err = errors.Cause(err) //获取原始对象
	log.Println(color.Error.Render(a...))
	s.SaveError(getStr(a...), "err")
}

// Errorf 记录错误信息
func (s *stdLog) Errorf(msg string, a ...interface{}) {
	log.Println(color.Error.Render(a...))
	s.SaveError(fmt.Sprintf(msg, a...), "err")
}

// Debug level info msg
func (s *stdLog) Debug(a ...interface{}) {
	if dev.IsDev() {
		log.Println(color.Debug.Render(a...))
	}
}

// Debug level info msg
func (s *stdLog) Debugf(msg string, a ...interface{}) {
	if dev.IsDev() {
		log.Println(color.Debug.Render(fmt.Sprintf(msg, a...)))
	}
}

//Fatal 系统级错误
func (s *stdLog) Fatal(a ...interface{}) {
	log.Output(2, color.Error.Render(fmt.Sprint(a...)))
	os.Exit(1)
}

//Fatalf 系统级错误
func (s *stdLog) Fatalf(msg string, a ...interface{}) {
	log.Output(2, color.Error.Render(fmt.Sprintf(msg, a...)))
	os.Exit(1)
}

// TraceError return trace of error
func (s *stdLog) TraceError(err error) error {
	e := errors.Cause(err) //获取原始对象
	log.Println(color.Error.Render(fmt.Sprintf(":Cause:%+v", e)))
	s.SaveError(getStr(err), "err")
	return errors.WithStack(err)
}

// ErrorString 打印错误信息
func (s *stdLog) ErrorString(a ...interface{}) {
	log.Output(2, color.Error.Render(fmt.Sprint(a...)))
}

func (s *stdLog) Close() {

}
