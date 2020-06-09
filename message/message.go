/*
	消息定义接口
*/
package message

import (
	"fmt"
	"sync"

	"github.com/xxjwxc/public/mylog"
)

// MessageBody 消息头
type MessageBody struct {
	State bool        `json:"state"`
	Code  int         `json:"code,omitempty"`
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

//GetErrorMsg 获取错误消息 参数(int,string)
func GetErrorMsg(errorCode ...interface{}) (msg MessageBody) {
	if len(errorCode) == 0 {
		mylog.ErrorString("unknow")
		msg.State = false
		msg.Code = -1
		return
	}
	msg.State = false
	for _, e := range errorCode {
		switch v := e.(type) {
		case int:
			msg.Code = int(v)
			msg.Error = ErrCode(v).String()
		case ErrCode:
			_tryRegisteryCode(v)
			msg.Code = int(v)
			msg.Error = v.String()
		case string:
			msg.Error = string(v)
		case error:
			msg.Error = v.Error()
		default:
			msg.Error = fmt.Sprintf("Unknow type:(%v)", e)
		}
	}
	return
}

// GetSuccessMsg 成功消息
func GetSuccessMsg(code ...ErrCode) (msg MessageBody) {
	if len(code) == 0 {
		code = append(code, NormalMessageID)
	}
	_tryRegisteryCode(code[0])

	msg.State = true
	msg.Code = int(code[0])
	msg.Error = code[0].String()
	return
}

//GetErrorStrMsg 获取错误消息 参数(int,string)
func GetErrorStrMsg(errorCode string) (msg MessageBody) {
	// if k, ok := _MessageMap[errorCode]; ok {
	// 	return GetErrorMsg(k)
	// }

	msg.State = false
	msg.Code = _tryGetCodeID(errorCode)
	msg.Error = errorCode
	return
}

var _mp sync.Map

func _tryRegisteryCode(code ErrCode) {
	_mp.LoadOrStore(code.String(), int(code))
}

func _tryGetCodeID(codeStr string) int {
	v, ok := _mp.Load(codeStr)
	if ok {
		return v.(int)
	}

	return -1
}
