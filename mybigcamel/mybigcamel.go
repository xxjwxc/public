package mybigcamel

import (
	"bytes"
	"strings"
)

/*
大驼峰命名规则转换
基本满足大驼峰命名法则 首字母大写 “_” 忽略后大写
device_id 对应 DeviceID create_time 对应 CreateTime location 对应 Location
*/

/*
	转换为大驼峰命名法则
	首字母大写，“_” 忽略后大写
*/
func Marshal(name string) string {
	if name == "" {
		return ""
	}

	temp := strings.Split(name, "_")
	var s string
	for _, v := range temp {
		vv := []rune(v)
		if len(vv) > 0 {
			if bool(vv[0] >= 'a' && vv[0] <= 'z') { //首字母大写
				vv[0] -= 32
			}
			s += string(vv)
		}
	}

	s = uncommonInitialismsReplacer.Replace(s)
	//smap.Set(name, s)
	return s
}

// UnMarshal  回退网络模式
func UnMarshal(name string) string {
	const (
		lower = false
		upper = true
	)

	if name == "" {
		return ""
	}

	var (
		value                                    = commonInitialismsReplacer.Replace(name)
		buf                                      = bytes.NewBufferString("")
		lastCase, currCase, nextCase, nextNumber bool
	)

	for i, v := range value[:len(value)-1] {
		nextCase = bool(value[i+1] >= 'A' && value[i+1] <= 'Z')
		nextNumber = bool(value[i+1] >= '0' && value[i+1] <= '9')

		if i > 0 {
			if currCase == upper {
				if lastCase == upper && (nextCase == upper || nextNumber == upper) {
					buf.WriteRune(v)
				} else {
					if value[i-1] != '_' && value[i+1] != '_' {
						buf.WriteRune('_')
					}
					buf.WriteRune(v)
				}
			} else {
				buf.WriteRune(v)
				if i == len(value)-2 && (nextCase == upper && nextNumber == lower) {
					buf.WriteRune('_')
				}
			}
		} else {
			currCase = upper
			buf.WriteRune(v)
		}
		lastCase = currCase
		currCase = nextCase
	}

	buf.WriteByte(value[len(value)-1])

	s := strings.ToLower(buf.String())
	return s
}

// UnSmallMarshal 小驼峰模式
func UnSmallMarshal(name string) string {
	if name == "" {
		return ""
	}

	var (
		value = commonInitialismsReplacer.Replace(name)
	)

	strArry := []rune(value)

	if bool(strArry[0] >= 'A' && strArry[0] <= 'Z') {
		strArry[0] = strArry[0] + 32
	}

	return string(strArry)
}
