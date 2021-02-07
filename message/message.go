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

func init() {
	_tryRegisteryCode(NormalMessageID)
	_tryRegisteryCode(NotFindError)
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
func GetSuccessMsg(codes ...ErrCode) (msg MessageBody) {
	code := NormalMessageID
	if len(codes) > 0 {
		code = codes[0]
	}
	_tryRegisteryCode(code)

	msg.State = true
	msg.Code = int(code)
	msg.Error = code.String()
	return
}

// GetError 获取错误信息
func GetError(code ErrCode) error {
	_tryRegisteryCode(code)
	return fmt.Errorf(code.String())
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
