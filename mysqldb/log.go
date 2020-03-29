package mysqldb

import (
	"fmt"

	"github.com/xxjwxc/public/mylog"

	"github.com/jinzhu/gorm"
)

var LogFormatter = func(values ...interface{}) (messages []interface{}) {
	if len(values) > 1 {
		var (
			currentTime = "\t[" + gorm.NowFunc().Format("2006-01-02 15:04:05") + "]"
			source      = fmt.Sprintf("(%v)\t", values[1])
		)
		messages = []interface{}{source, currentTime}
		messages = append(messages, "\t [")
		messages = append(messages, values[2:]...)
		messages = append(messages, "]")

	}

	return
}

//
type DbLog struct {
	gorm.Logger
}

//
func (db DbLog) Print(values ...interface{}) {
	msg := LogFormatter(values...)
	str := fmt.Sprint(msg...)
	mylog.SaveError(str, "sql")
}
