package mylog

import (
	"data/config"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime/debug"
	"time"
)

func init() {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	path = filepath.Dir(path)
	os.MkdirAll(path+"/err", os.ModePerm) //生成多级目录
}

const ( //
	Log_Error   = iota //打印 Error 及以上级别
	Log_warning        //打印 warning 及以上级别
	Log_Info           //默认的返回值，为0，自增 //打印 Info 及以上级别
)

//
func Print(log_level int, describ string) {
	log.Println(describ)
	return
}

//
func Println(describ ...interface{}) {
	for _, e := range describ {
		switch v := e.(type) {
		case string:
			log.Println(v)
		case []byte:
			log.Println(string(v))
		default:
			log.Println(v)
		}
	}
	return
}

//
func Info(describ string) {
	log.Println(describ)
	return
}

//
func Error(err error) {
	log.Println(err)
	SaveError(err.Error(), "err")
}

//保存错误信息
func SaveError(errstring, flag string) {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	path = filepath.Dir(path)

	now := time.Now()                                              //获取当前时间
	time_str := now.Format("2006-01-02_15")                        //设定时间格式
	fname := fmt.Sprintf("%s/err/%s_%s.log", path, flag, time_str) //保存错误信息文件名:程序名-进程ID-当前时间（年月日时分秒）

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
	if config.OnIsDev() {
		for _, e := range describ {
			switch v := e.(type) {
			case string:
				log.Println(v)
			case []byte:
				log.Println(string(v))
			default:
				log.Println(v)
			}
		}
	}
}

//刷新
func Flush() {

}
