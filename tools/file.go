package tools

import (
	"bufio"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/xxjwxc/public/mylog"
)

// CheckFileIsExist 检查目录是否存在
func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		mylog.Debug(err)
		exist = false
	}
	return exist
}

// BuildDir 创建目录
func BuildDir(absDir string) error {
	return os.MkdirAll(path.Dir(absDir), os.ModePerm) //生成多级目录
}

// DeleteFile 删除文件或文件夹
func DeleteFile(absDir string) error {
	return os.RemoveAll(absDir)
}

// GetPathDirs 获取目录所有文件夹
func GetPathDirs(absDir string) (re []string) {
	if CheckFileIsExist(absDir) {
		files, _ := ioutil.ReadDir(absDir)
		for _, f := range files {
			if f.IsDir() {
				re = append(re, f.Name())
			}
		}
	}
	return
}

// GetPathFiles 获取目录所有文件
func GetPathFiles(absDir string) (re []string) {
	if CheckFileIsExist(absDir) {
		files, _ := ioutil.ReadDir(absDir)
		for _, f := range files {
			if !f.IsDir() {
				re = append(re, f.Name())
			}
		}
	}
	return
}

// GetModelPath 获取程序运行目录
func GetModelPath() string {
	dir, _ := os.Getwd()
	return strings.Replace(dir, "\\", "/", -1)
}

// GetCurrentDirectory 获取exe所在目录
func GetCurrentDirectory() string {
	dir, _ := os.Executable()
	exPath := filepath.Dir(dir)
	// fmt.Println(exPath)

	return strings.Replace(exPath, "\\", "/", -1)
}

// SaveToFile 写入文件
func SaveToFile(fname string, src []string, isClear bool) bool {
	return WriteFile(fname, src, isClear)
}

// WriteFile 写入文件
func WriteFile(fname string, src []string, isClear bool) bool {
	BuildDir(fname)
	flag := os.O_CREATE | os.O_WRONLY | os.O_TRUNC
	if !isClear {
		flag = os.O_CREATE | os.O_RDWR | os.O_APPEND
	}
	f, err := os.OpenFile(fname, flag, 0666)
	if err != nil {
		mylog.Error(err)
		return false
	}
	defer f.Close()

	for _, v := range src {
		f.WriteString(v)
		f.WriteString("\r\n")
	}

	return true
}

// ReadFile 读取文件
func ReadFile(fname string) (src []string) {
	f, err := os.OpenFile(fname, os.O_RDONLY, 0666)
	if err != nil {
		return []string{}
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	for {
		line, _, err := rd.ReadLine()
		if err != nil || io.EOF == err {
			break
		}
		src = append(src, string(line))
	}

	return src
}

// MoveFile 移动文件或文件夹(/结尾)
func MoveFile(from, to string) error {
	// if !CheckFileIsExist(to) {
	// 	BuildDir(to)
	// }
	return os.Rename(from, to)
}

func CopyFile(src, des string) error {
	if !CheckFileIsExist(des) {
		BuildDir(des)
	}
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	desFile, err := os.Create(des)
	if err != nil {
		return err
	}
	defer desFile.Close()

	_, err = io.Copy(desFile, srcFile)
	return err
}
