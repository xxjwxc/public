package mylog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"sync"
	"time"
)

type errDeal struct {
	once  sync.Once
	_path string
}

func (s *errDeal) initPath() {
	s._path = getCurrentDirectory() + "/err"
	os.MkdirAll(s._path, os.ModePerm) //生成多级目录
}

// SaveError 保存错误信息
func (s *errDeal) SaveError(errstring, flag string) {
	s.once.Do(s.initPath)

	now := time.Now()                                            //获取当前时间
	timeStr := now.Format("2006-01-02_15")                       //设定时间格式
	fname := fmt.Sprintf("%s/%s_%s.log", s._path, flag, timeStr) //保存错误信息文件名:程序名-进程ID-当前时间（年月日时分秒）

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

// Panic return trace of error
func (s *errDeal) Panic(a ...interface{}) {
	s.once.Do(s.initPath)

	now := time.Now()                                                 //获取当前时间
	pid := os.Getpid()                                                //获取进程ID
	timeStr := now.Format("2006-01-02")                               //设定时间格式
	fname := fmt.Sprintf("%s/panic_%s-%x.log", s._path, timeStr, pid) //保存错误信息文件名:程序名-进程ID-当前时间（年月日时分秒）
	fmt.Println("panic to file ", fname)

	f, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return
	}
	defer f.Close()

	f.WriteString("=========================" + now.Format("2006-01-02 15:04:05 ========================= \r\n"))
	f.WriteString(getStr(a...)) //输出堆栈信息
	f.WriteString("=========================end=========================")
}

// GetCurrentDirectory 获取exe所在目录
func getCurrentDirectory() string {
	dir, _ := os.Executable()
	exPath := filepath.Dir(dir)
	// fmt.Println(exPath)

	return exPath
}

func getStr(a ...interface{}) string {
	if len(a) == 1 {
		switch v := a[0].(type) {
		case string, []byte, int8, int, int32, int64, float32, float64, time.Time, bool, error: // 系统变量
			return fmt.Sprintf("%v", v)
		default:
			return fmt.Sprintf("%#v", v)
		}
	}

	var rep string
	for _, v := range a {
		switch v.(type) {
		case string, []byte, int8, int, int32, int64, float32, float64, time.Time, bool, error: // 系统变量
			rep += "[%v]"
		default:
			rep += "[%#v]"
		}
	}
	return fmt.Sprintf(rep, a...)
}
