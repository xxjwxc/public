package mylog

import (
	"data/config"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"public/myelastic"
	"public/myqueue"
	"time"
)

var ptr_que *myqueue.MyQueue = nil
var elastic myelastic.MyElastic
var isSaveFile bool = true        //默认存文件
var isSaveToEs bool = false       //默认不保存
var local_Log_file string = "log" //默认存放文件的目录
var exe_path string

func init() {
	if config.IsRunTesting() { //测试时候不创建
		return
	}

	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	exe_path = filepath.Dir(path)

	BuildDir(local_Log_file)
	ptr_que = myqueue.NewSyncQueue()
	es_path := config.GetEsAddrUrl()
	if len(es_path) > 0 {
		elastic = myelastic.OnInitES(es_path)
		elastic.CreateIndex(Http_log_index, mapping)
	}
	go onConsumerLog()
}

/*
	发送日志请求
*/
func OnLog(es_index, es_type, es_id string, data LogInfo) {
	//	local, _ := time.LoadLocation("Local")
	//	data.Creat_time = time.Now().In(local)
	b, err := json.Marshal(data.Data)
	if err != nil {
		log.Println("OnLog error:", err)
	}
	data.Data = string(b)
	var info EsLogInfo
	info.Es_index = es_index
	info.Es_type = es_type
	info.Es_id = es_id
	info.Info = data

	ptr_que.Push(info) //加入日志队列
}

/*
	更新本地存储文件地址
	isSave:是否本地存储
	LogFile:本地存储相对程序位置(log/  ==> 当前可执行文件的 log/目录)
*/
func InitLogFileInfo(isSave, isSaveEs bool, LogFile string) {
	isSaveFile = isSave
	isSaveToEs = isSaveEs
	local_Log_file = LogFile
	if isSave {
		BuildDir(local_Log_file)
	}
}

func BuildDir(logfile string) {
	os.MkdirAll(exe_path+"/"+logfile, os.ModePerm) //生成多级目录
}

/*
	消费者 消费日志
*/
func onConsumerLog() {
	for {
		var info EsLogInfo
		info = ptr_que.Pop().(EsLogInfo)

		if isSaveToEs && elastic.Client != nil {
			if !elastic.Add(info.Es_index, info.Es_type, info.Es_id, info.Info) {
				log.Println("elastic add error ")
			}
		}

		if isSaveFile {
			saveLogTofile(info)
		}
		Debug(info)
	}
}

var _f *os.File
var _err error
var saveFaile string

func saveLogTofile(info EsLogInfo) {

	time_str := time.Now().Format("2006-01-02-15") //设定时间格式
	fname := fmt.Sprintf("%s/%s/%s.log", exe_path, local_Log_file, time_str)
	if saveFaile != fname {
		if _f != nil {
			_f.Close()
		}
		_f, _err = os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if _err != nil {
			log.Println(_err)
			return
		}
	}

	b, _ := json.Marshal(info)
	_f.WriteString(string(b) + "\r\n") //输出堆栈信息
}
