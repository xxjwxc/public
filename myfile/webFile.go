package myfile

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/xxjwxc/public/tools"
)

type myFile struct {
	path       string
	isRelative bool
}

// NewWebFile 新建文件
// @parm  path:目录 isRelative:是否相对路径
func NewWebFile(path string, isRelative bool) *myFile {
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}

	return &myFile{path: path, isRelative: isRelative}
}

// UploadOneFile 单文件上传
func (o *myFile) UploadOneFile(r *http.Request, field string) (string, error) {
	//接受post请求
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		if r.MultipartForm == nil {
			return "", fmt.Errorf("empty")
		}
		_, fh, err := r.FormFile(field) //文件name
		if err != nil {
			return "", err
		}
		return o.SaveOne(fh)
	}

	return "", fmt.Errorf("method not support")
}

// UploadMoreFile 多文件上传
func (o *myFile) UploadMoreFile(r *http.Request) ([]string, error) {
	//接受post请求
	if r.Method == "POST" {
		r.ParseMultipartForm(32 << 20)
		if r.MultipartForm == nil {
			return []string{}, fmt.Errorf("empty")
		}

		var files []*multipart.FileHeader
		for _, v := range r.MultipartForm.File {
			files = append(files, v...)
		}

		result := make([]string, 0, len(files))
		for _, file := range files {
			src, err := o.SaveOne(file)
			if err != nil {
				return []string{}, err
			}
			result = append(result, src)
		}

		return result, nil
	}

	return []string{}, fmt.Errorf("method not support")
}

// SaveOne 保存一个
func (o *myFile) SaveOne(file *multipart.FileHeader) (string, error) {
	filename := getFileName(GetExp(file.Filename))
	path := o.path
	if o.isRelative {
		path = tools.GetCurrentDirectory() + "/" + path
	}

	if !tools.CheckFileIsExist(path) {
		if err := tools.BuildDir(path); err != nil { //创建文件夹
			return "", err
		}
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.Create(path + filename)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return o.path + filename, err
}

// SaveOrigin 原始保存一个
func (o *myFile) SaveOrigin(file *multipart.FileHeader, dir string) (string, error) {
	filename := file.Filename
	_path := path.Join(o.path, dir)
	if o.isRelative {
		_path = path.Join(tools.GetCurrentDirectory(), _path)
	}
	if !strings.HasSuffix(_path, "/") {
		_path += "/"
	}

	if !tools.CheckFileIsExist(_path) {
		if err := tools.BuildDir(_path); err != nil { //创建文件夹
			return "", err
		}
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	out, err := os.OpenFile(path.Join(_path, filename), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return path.Join(o.path, dir, filename), err
}

// GetFileMD5 获取文件的MD5值
// filePath: 文件路径
// 返回: MD5字符串, 错误信息
func GetFileMD5(filePath string) (string, error) {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close() // 函数结束时关闭文件

	// 创建MD5哈希对象
	hash := md5.New()

	// 分块读取文件并写入哈希器（高效处理大文件）
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", err
	}

	// 计算哈希值并转为十六进制字符串
	md5Bytes := hash.Sum(nil)
	md5Str := hex.EncodeToString(md5Bytes)

	return md5Str, nil
}
