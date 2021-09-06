package tools

import (
	"crypto/md5"
	crand "crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"reflect"
	"sort"
	"strings"
	"unicode"

	"github.com/xxjwxc/public/errors"
)

// Md5Encoder md5 加密
func Md5Encoder(src string) string {
	h := md5.New()
	h.Write([]byte(src)) // 需要加密的字符串
	//    fmt.Printf("%x\n", h.Sum(nil)) // 输出加密结
	ret := fmt.Sprintf("%x", h.Sum(nil))
	return strings.ToUpper(ret)
}

// Copy 合并数组
func Copy(dest []interface{}, src []interface{}) (result []interface{}) {
	result = make([]interface{}, len(dest)+len(src))
	copy(result, dest)
	copy(result[len(dest):], src)
	return
}

// DeleteArray 删除数组
func DeleteArray(src []interface{}, index int) (result []interface{}) {
	result = append(src[:index], src[(index+1):]...)
	return
}

// GetMd5String 生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// GetTotalPageNum 获取总页数
func GetTotalPageNum(pageSize, totalCount int) int {
	return (totalCount + pageSize - 1) / pageSize
}

// UniqueID 生成32位guid
func UniqueID() string {
	b := make([]byte, 48)
	if _, err := io.ReadFull(crand.Reader, b); err != nil {
		return ""
	}
	return GetMd5String(base64.URLEncoding.EncodeToString(b))
}

// DeleteSlice 删除切片index
func DeleteSlice(slice interface{}, index int) (interface{}, error) {
	sliceValue := reflect.ValueOf(slice)
	length := sliceValue.Len()
	if slice == nil || length == 0 || (length-1) < index {
		return nil, errors.New("error")
	}
	if length-1 == index {
		return sliceValue.Slice(0, index).Interface(), nil
	}

	return reflect.AppendSlice(sliceValue.Slice(0, index), sliceValue.Slice(index+1, length)).Interface(), nil
}

// MinimumInt 查找int最小值
func MinimumInt(rest []int) int {
	minimum := rest[0]
	for _, v := range rest {
		if v < minimum {
			minimum = v
		}
	}
	return minimum
}

// func LoadTemplate(list ...string) *template.Template {
// 	var tmp []string
// 	for _, v := range list {
// 		if CheckFileIsExist(GetModelPath() + config.Static_host[0] + v) {
// 			tmp = append(tmp, GetModelPath()+config.Static_host[0]+v)
// 		} else {
// 			mylog.Debug("file does not exist:" + GetModelPath() + config.Static_host[0] + v)
// 			panic(GetModelPath() + config.Static_host[0] + v)
// 		}
// 	}
// 	return template.Must(template.ParseFiles(tmp...))
// }

/*
	执行模版渲染，
	name:模版名字，""则无名字
	data:传参列表
	list:模版列表
*/
// func ExecuteTemplate(w rest.ResponseWriter, name string, data interface{}, list ...string) error {
// 	t := LoadTemplate(list...)
// 	w.(http.ResponseWriter).Header().Set("Content-Type", "text/html; charset=utf-8")

// 	if len(name) == 0 {
// 		return t.Execute(w.(http.ResponseWriter), data)
// 	} else {
// 		return t.ExecuteTemplate(w.(http.ResponseWriter), name, data)
// 	}
// }

// DictSort 按字典顺序排序
func DictSort(res []string) (str string) {
	sort.Strings(res)
	if len(res) > 0 {
		for _, v := range res {
			str += v
		}
	}
	return
}

// Sha1Encrypt SHA1加密
func Sha1Encrypt(str string) string {
	h := sha1.New()
	h.Write([]byte(str))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}

// GetUtf8Str 中文字符切割时有问题。采用此方式不会有问题
func GetUtf8Str(str string) []rune {
	return []rune(str)
}

// GetUtf8Len 获取中文字符的长度
func GetUtf8Len(str string) int {
	return len([]rune(str))
}

// IsHan 判断是否有中文
func IsHan(str string) bool {
	for _, r := range str {
		if unicode.Is(unicode.Han, r) {
			return true
		}
	}
	return false
}
