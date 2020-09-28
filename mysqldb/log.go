package mysqldb

import (
	"log"
	"time"

	"github.com/xxjwxc/public/mylog"

	"gorm.io/gorm/logger"
)

// DbLog ...
type DbLog struct {
}

// Write ...
func (lg DbLog) Write(p []byte) (n int, err error) {
	mylog.SaveError(string(p), "sql")
	return len(p), err
}

// GetDBlog 获取默认logger
func GetDBlog() logger.Interface {
	newLogger := logger.New(
		log.New(DbLog{}, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,  // 慢 SQL 阈值
			LogLevel:      logger.Error, // Log level
			Colorful:      false,        // 禁用彩色打印
		},
	)
	return newLogger
}
