package tools

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// CheckParam 检测参数
func CheckParam(params ...string) bool {
	for _, value := range params {
		if len(value) == 0 {
			return false
		}
	}
	return true
}

// IsPhone 判断是否是手机号
func IsPhone(mobileNum string) bool {
	tmp := `^(13[0-9]|14[579]|15[0-3,5-9]|16[6]|17[0135678]|18[0-9]|19[89])\d{8}$`
	reg := regexp.MustCompile(tmp)
	return reg.MatchString(mobileNum)
}

// IsMail 判断用户是否是邮件用户
func IsMail(username string) (isMail bool) {
	isMail = false
	if strings.Contains(username, "@") {
		isMail = true //是邮箱
	}
	return
}

// IsRunTesting 判断是否在测试环境下使用
func IsRunTesting() bool {
	if len(os.Args) > 1 {
		fmt.Println(os.Args[1])
		return strings.HasPrefix(os.Args[1], "-test")
	}
	return false
}

// IsIDCard 判断是否是18或15位身份证
func IsIDCard(cardNo string) bool {
	//18位身份证 ^(\d{17})([0-9]|X)$
	if m, _ := regexp.MatchString(`(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)`, cardNo); !m {
		return false
	}
	return true
}

var internalType = []string{"string", "bool", "int", "uint", "byte", "rune",
	"int8", "int16", "int32", "int64", "uint8", "uint16", "uint32", "uint64", "uintptr",
	"float32", "float64", "map", "Time"}

// IsInternalType 是否是内部类型
func IsInternalType(t string) bool {
	for _, v := range internalType {
		if strings.EqualFold(t, v) {
			return true
		}
	}
	return false
}

var keywords = []string{"var", "const", "package", "import", "func", "return",
	"defer", "go", "select", "interface", "struct", "break", "case", "continue", "for",
	"fallthrough", "else", "if", "switch", "goto", "default", "chan", "type", "map", "range"}

// IsKeywords 是否是关键字
func IsKeywords(t string) bool {
	for _, v := range keywords {
		if t == v {
			return true
		}
	}
	return false
}
