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
	Log_Error   = iota //打印 Error 及以上级别
	Log_warning        //打印 warning 及以上级别
	Log_Info           //默认的返回值，为0，自增 //打印 Info 及以上级别
)

//
func Print(log_level int, describ string) {
	log.Println(color.Info.Render(describ))
	return
}

//
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

//
func Info(describ string) {
	log.Println(color.FgGreen.Render(describ))
	return
}

//
func Error(err error) {
	err = errors.Cause(err) //获取原始对象
	log.Println(color.Error.Render(fmt.Sprintf(":Cause:%+v", err)))
	SaveError(fmt.Sprintf("%+v", err), "err")
}

//打印错误信息
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

//保存错误信息
func SaveError(errstring, flag string) {
	_logerTuple.once.Do(initPath)

	now := time.Now()                                                       //获取当前时间
	time_str := now.Format("2006-01-02_15")                                 //设定时间格式
	fname := fmt.Sprintf("%s/%s_%s.log", _logerTuple._path, flag, time_str) //保存错误信息文件名:程序名-进程ID-当前时间（年月日时分秒）

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

//
func Debug(describ ...interface{}) {
	if dev.OnIsDev() {
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
