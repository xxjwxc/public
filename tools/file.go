package tools

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/xie1xiao1jun/public/mylog"
)

//检查目录是否存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		mylog.Debug(filename + " not exist")
		exist = false
	}
	return exist
}

//创建目录
func BuildDir(abs_dir string) error {
	return os.MkdirAll(abs_dir, os.ModePerm) //生成多级目录
}

//删除文件或文件夹
func DeleteFile(abs_dir string) error {
	return os.RemoveAll(abs_dir)
}

//获取目录所有文件夹
func GetPathDirs(abs_dir string) (re []string) {
	if CheckFileIsExist(abs_dir) {
		files, _ := ioutil.ReadDir(abs_dir)
		for _, f := range files {
			if f.IsDir() {
				re = append(re, f.Name())
			}
		}
	}
	return
}

//获取目录所有文件夹
func GetPathFiles(abs_dir string) (re []string) {
	if CheckFileIsExist(abs_dir) {
		files, _ := ioutil.ReadDir(abs_dir)
		for _, f := range files {
			if !f.IsDir() {
				re = append(re, f.Name())
			}
		}
	}
	return
}

//获取目录地址
func GetModelPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	// if len(path) > 0 {
	// 	path += "/"
	// }
	path = filepath.Dir(path)
	return path
}
