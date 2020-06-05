package mylog

import (
	"encoding/json"
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

// GetDefaultStd get default std logger
func GetDefaultStd() *stdLog {
	return &stdLog{}
}

// Info level info msg
func (s *stdLog) Info(a ...interface{}) {
	log.Println(color.FgGreen.Render(getStr(a...)))
}

// Info level info msg
func (s *stdLog) Infof(msg string, a ...interface{}) {
	log.Println(color.FgGreen.Render(fmt.Sprintf(msg, a...)))
}

// Error 记录错误信息
func (s *stdLog) Error(a ...interface{}) {
	// err = errors.Cause(err) //获取原始对象
	log.Println(color.Error.Render(getStr(a...)))
	s.SaveError(getStr(a...), "err")
}

// Errorf 记录错误信息
func (s *stdLog) Errorf(msg string, a ...interface{}) {
	log.Println(color.Error.Render(fmt.Sprintf(msg, a...)))
	s.SaveError(fmt.Sprintf(msg, a...), "err")
}

// Debug level info msg
func (s *stdLog) Debug(a ...interface{}) {
	if dev.IsDev() {
		log.Println(color.Debug.Render(getStr(a...)))
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
	log.Output(2, color.Error.Render(getStr(a...)))
	os.Exit(1)
}

//Fatalf 系统级错误
func (s *stdLog) Fatalf(msg string, a ...interface{}) {
	log.Output(2, color.Error.Render(fmt.Sprintf(msg, a...)))
	os.Exit(1)
}

//JSON json输出
func (s *stdLog) JSON(a ...interface{}) {
	for _, v := range a {
		b, _ := json.MarshalIndent(v, "", "     ")
		log.Println(color.FgGreen.Render(string(b)))
	}
}

// TraceError return trace of error
func (s *stdLog) TraceError(err error) error {
	e := errors.Cause(err) //获取原始对象
	log.Println(color.Error.Render(getStr(e)))
	s.SaveError(getStr(err), "err")
	return errors.WithStack(err)
}

// ErrorString 打印错误信息
func (s *stdLog) ErrorString(a ...interface{}) {
	log.Output(2, color.Error.Render(getStr(a...)))
}

func (s *stdLog) Close() {

}
