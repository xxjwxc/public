package tools

import (
	"strconv"
	"strings"

	"github.com/axgle/mahonia"
)

// FormatInt int数组转字符串
func FormatInt(list []int, seg string) string {
	s := make([]string, len(list))
	for i := range list {
		s[i] = strconv.Itoa(list[i])
	}

	return strings.Join(s, seg)
}

//字符串去除特殊字符
//func DealStr(str, replace string) string {
//	var s string = strings.TrimSpace(str) //去除尾部空格
//	strings.Replace(s, replace, "", -1)
//	return s
//}

// Substr 截取字符串 不包括str
func Substr(str string, start, length int) string {
	rs := []rune(str)
	rl := len(rs)
	end := 0

	if start < 0 {
		start = rl - 1 + start
	}
	end = start + length

	if start > end {
		start, end = end, start
	}

	if start < 0 {
		start = 0
	}
	if start > rl {
		start = rl
	}
	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	return string(rs[start:end])
}

// MinimumString 查找字符串最小值
func MinimumString(rest []string) string {
	minimum := rest[0]
	for _, v := range rest {

		if v := v; v < minimum {
			minimum = v
		}

	}
	return minimum
}

// ConvertToString 字符集转换
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	return string(cdata)
}

// ConvertString 系统转其他
func ConvertString(src string, tagCode string) string {
	enc := mahonia.NewEncoder(tagCode)
	return enc.ConvertString(src)
}

// GetGBK 获取gbk
func GetGBK(src string) string {
	return string(ConvertString(src, "gbK"))
}

// Reverse 反转字符串
func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
