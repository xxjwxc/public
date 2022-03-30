/*
	消息定义接口
*/
package message

import (
	"fmt"

	"github.com/xxjwxc/public/mylog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MessageBody 消息头
type MessageBody struct {
	State bool        `json:"state"`
	Code  int         `json:"code,omitempty"`
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// GetError 获取错误信息(grpc)
// func (m *MessageBody) GetError() error {
// 	return status.Errorf(codes.Code(m.Code), m.Error)
// }

// GetError 获取错误信息(grpc)
func (m MessageBody) GetError() error {
	return status.Errorf(codes.Code(m.Code), m.Error)
}

// func init() {
// 	_tryRegisteryCode(NormalMessageID)
// 	_tryRegisteryCode(NotFindError)
// }

// GetErrorMsg 获取错误消息 参数(int,string)
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

	msg.State = true
	msg.Code = int(code)
	msg.Error = code.String()
	return
}

// GetError 获取错误信息(grpc)
func GetError(code ErrCode) error {
	return status.Errorf(codes.Code(code), code.String())
}

//GetErrorStrMsg 获取错误消息 参数(int,string)
func GetErrorStrMsg(err error) (msg MessageBody) {
	msg.State = false
	gerr := status.Convert(err)
	if gerr != nil {
		msg.Code = int(gerr.Code())
		msg.Error = gerr.Message()
	} else {
		msg.Error = err.Error()
	}
	return
}

// var _mp sync.Map

// func _tryRegisteryCode(code ErrCode) {
// 	_mp.LoadOrStore(code.String(), int(code))
// }

// func _tryGetCodeID(codeStr string) int {
// 	v, ok := _mp.Load(codeStr)
// 	if ok {
// 		return v.(int)
// 	}

// 	return -1
// }
