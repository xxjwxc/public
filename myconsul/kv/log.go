package kv

import (
	"encoding/json"
	"log"
	"os"

	"github.com/xxjwxc/public/mylog"
)

// Logger ...
type Logger interface {
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

type kvLog struct {
	*log.Logger
}

// NewLogger ...
func NewLogger() Logger {
	return &kvLog{
		log.New(os.Stdout, "", log.LstdFlags),
	}
}

// Info ...
func (l *kvLog) Info(msg string, args ...interface{}) {
	l.output("INFO", msg, args...)
}

// Warn ...
func (l *kvLog) Warn(msg string, args ...interface{}) {
	l.output("WARN", msg, args...)
}

// Error ...
func (l *kvLog) Error(msg string, args ...interface{}) {
	l.output("ERROR", msg, args...)
}

func (l *kvLog) output(level, msg string, args ...interface{}) {
	length := len(args)
	if length%2 != 0 {
		l.Println(level, msg, "invalid key-value pairs")
		return
	}

	m := make(map[string]interface{})
	for i := 0; i < length; i = i + 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		m[key] = args[i+1]
	}

	b, _ := json.Marshal(m)
	switch level {
	case "INFO":
		mylog.Info(level, msg, string(b))
	case "WARN":
		mylog.Info(level, msg, string(b))
	case "ERROR":
		mylog.Error(level, msg, string(b))
	default:
		mylog.Info(level, msg, string(b))
	}
}
