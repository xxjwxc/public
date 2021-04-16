package tools

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/xxjwxc/public/errors"

	"github.com/xxjwxc/public/mylog"
)

// RawBytes ...
type RawBytes []byte

var errNilPtr = errors.New("destination pointer is nil")

func cloneBytes(b []byte) []byte {
	if b == nil {
		return nil
	}

	c := make([]byte, len(b))
	copy(c, b)
	return c
}

// AsString 转成string
func AsString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int:
		return strconv.Itoa(v)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		return GetTimeStr(v)
	case bool:
		return strconv.FormatBool(v)
	default:
		{
			b, _ := json.Marshal(v)
			return string(b)
		}
	}
	// return fmt.Sprintf("%v", src)
}

// EncodeByte 编码二进制
func EncodeByte(data interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(data)
	if err != nil {
		mylog.Error(err)
		return nil, err
	}
	return buf.Bytes(), nil
}

// DecodeByte 解码二进制
func DecodeByte(data []byte, to interface{}) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	return dec.Decode(to)
}

// ByteToHex byte转16进制字符串
func ByteToHex(data []byte) string {

	return hex.EncodeToString(data)

	// buffer := new(bytes.Buffer)
	// for _, b := range data {

	// 	s := strconv.FormatInt(int64(b&0xff), 16)
	// 	if len(s) == 1 {
	// 		buffer.WriteString("0")
	// 	}
	// 	buffer.WriteString(s)
	// }

	// return buffer.String()
}

// HexToBye 16进制字符串转[]byte
func HexToBye(hexStr string) []byte {
	hr, _ := hex.DecodeString(hexStr)
	return hr

	// length := len(hex) / 2
	// slice := make([]byte, length)
	// rs := []rune(hex)

	// for i := 0; i < length; i++ {
	// 	s := string(rs[i*2 : i*2+2])
	// 	value, _ := strconv.ParseInt(s, 16, 10)
	// 	slice[i] = byte(value & 0xFF)
	// }
	// return slice
}

// UnicodeEmojiDecode Emoji表情解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

// UnicodeEmojiCode Emoji表情转换
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u
		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

// DbcToSbc 全角转半角
func DbcToSbc(str string) string {
	numConv := unicode.SpecialCase{
		unicode.CaseRange{
			Lo: 0x3002, // Lo 全角句号
			Hi: 0x3002, // Hi 全角句号
			Delta: [unicode.MaxCase]rune{
				0,               // UpperCase
				0x002e - 0x3002, // LowerCase 转成半角句号
				0,               // TitleCase
			},
		},
		//
		unicode.CaseRange{
			Lo: 0xFF01, // 从全角！
			Hi: 0xFF19, // 到全角 9
			Delta: [unicode.MaxCase]rune{
				0,               // UpperCase
				0x0021 - 0xFF01, // LowerCase 转成半角
				0,               // TitleCase
			},
		},
		unicode.CaseRange{
			Lo: 0xff21, // Lo: 全角 Ａ
			Hi: 0xFF5A, // Hi:到全角 ｚ
			Delta: [unicode.MaxCase]rune{
				0,               // UpperCase
				0x0041 - 0xff21, // LowerCase 转成半角
				0,               // TitleCase
			},
		},
	}

	return strings.ToLowerSpecial(numConv, str)
}
