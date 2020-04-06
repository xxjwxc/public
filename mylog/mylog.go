package mylog

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"

	"github.com/gookit/color"
	"github.com/xxjwxc/public/dev"
	"github.com/xxjwxc/public/errors"
)

var _logerTuple logerTuple

type logerTuple struct {
	once  sync.Once
	_path string
}

const ( //
	logError   = iota //打印 Error 及以上级别
	logwarning        //打印 warning 及以上级别
	logInfo           //默认的返回值，为0，自增 //打印 Info 及以上级别
)

// Print 打印信息
func Print(logLevel int, describ string) {
	log.Println(color.Info.Render(describ))
	return
}

// Println 打印信息
func Println(describ ...interface{}) {
	for _, e := range describ {
		switch v := e.(type) {
		// case string:
		// 	log.Println(color.Info.Render(v))
		case []byte:
			log.Println(color.Info.Render(string(v)))
		default:
			log.Println(color.Info.Render(v))
		}
	}
	return
}

// Info ...
func Info(describ string) {
	log.Println(color.FgGreen.Render(describ))
	return
}

// Error 记录错误信息
func Error(err error) {
	err = errors.Cause(err) //获取原始对象
	log.Println(color.Error.Render(fmt.Sprintf(":Cause:%+v", err)))
	SaveError(fmt.Sprintf("%+v", err), "err")
}

// TraceError 追踪器
func TraceError(err error) error {
	Error(err)

	return errors.WithStack(err)
}

// ErrorString 打印错误信息
func ErrorString(v ...interface{}) {
	log.Output(2, color.Error.Render(fmt.Sprint(v...)))
}

//Fatal 系统级错误
func Fatal(v ...interface{}) {
	log.Output(2, color.Error.Render(fmt.Sprint(v...)))
	os.Exit(1)
}

func initPath() {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	path = filepath.Dir(path)
	_logerTuple._path = path + "/err"
	os.MkdirAll(_logerTuple._path, os.ModePerm) //生成多级目录
}

// SaveError 保存错误信息
func SaveError(errstring, flag string) {
	_logerTuple.once.Do(initPath)

	now := time.Now()                                                      //获取当前时间
	timeStr := now.Format("2006-01-02_15")                                 //设定时间格式
	fname := fmt.Sprintf("%s/%s_%s.log", _logerTuple._path, flag, timeStr) //保存错误信息文件名:程序名-进程ID-当前时间（年月日时分秒）

	f, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	f.WriteString("=========================" + now.Format("2006-01-02 15:04:05 ========================= \r\n"))
	f.WriteString(errstring + "\r\n")             //输出堆栈信息
	f.WriteString(string(debug.Stack()) + "\r\n") //输出堆栈信息)
	f.WriteString("=========================end=========================\r\n")
}

// Debug debug
func Debug(describ ...interface{}) {
	if dev.IsDev() {
		for _, e := range describ {
			switch v := e.(type) {
			case string:
				log.Println(color.Note.Render(v))
			case []byte:
				log.Println(color.Note.Render(string(v)))
			default:
				log.Println(color.Note.Render(fmt.Sprintf("%+v", v)))
			}
		}
	}
}
