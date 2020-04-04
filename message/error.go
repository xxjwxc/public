package message

import "github.com/xxjwxc/public/errors"

// GetErrFromID 通过id返回错误信息
func GetErrFromID(codeID int) error {
	if _, ok := MessageMap[codeID]; ok {
		return errors.New(MessageMap[codeID])
	}

	// 返回默认值
	return errors.New(MessageMap[-1])
}
