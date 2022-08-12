package mylog

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/xxjwxc/public/dev"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapLog struct {
	logger *zap.Logger
	errDeal
}

// GetDefaultZap 获取默认zap日志库
func GetDefaultZap() *zapLog {
	dir := fmt.Sprintf("%v/log/%v.log", getCurrentDirectory(), dev.GetService())
	os.MkdirAll(path.Dir(dir), os.ModePerm) //生成多级目录
	hook := lumberjack.Logger{
		Filename:   dir,   // 日志文件路径
		MaxSize:    100,   // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: 10,    // 日志文件最多保存多少个备份
		MaxAge:     30,    // 文件最多保存多少天
		Compress:   false, // 是否压缩
	}
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	level := zap.InfoLevel
	if dev.IsDev() {
		level = zap.DebugLevel
	}
	atomicLevel.SetLevel(level)
	core := zapcore.NewCore(
		//zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		zapcore.NewJSONEncoder(encoderConfig),                                           // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)
	// 开启开发模式，堆栈跟踪
	caller := zap.AddCaller()
	// 开启文件及行号
	development := zap.Development()
	// 设置初始化字段
	filed := zap.Fields(zap.String("service", dev.GetService()))
	return &zapLog{
		logger: zap.New(core, caller, development, filed, zap.AddCallerSkip(2)), // 构造日志
	}
}

// SetLogger set net logger
func (z *zapLog) SetLogger(logger *zap.Logger) {
	z.logger = logger
}

// Info level info msg
func (z *zapLog) Info(a ...interface{}) {
	z.logger.Info(getStr(a...))
}

// Info level info msg
func (z *zapLog) Infof(msg string, a ...interface{}) {
	z.logger.Info(fmt.Sprintf(msg, a...))
}

// Error 记录错误信息
func (z *zapLog) Error(a ...interface{}) {
	// err = errors.Cause(err) //获取原始对象
	z.logger.Error(getStr(a...))
	z.SaveError(getStr(a...), "err")
}

// Errorf 记录错误信息
func (z *zapLog) Errorf(msg string, a ...interface{}) {
	z.logger.Error(fmt.Sprintf(msg, a...))
	z.SaveError(fmt.Sprintf(msg, a...), "err")
}

// ErrorString 打印错误信息
func (z *zapLog) ErrorString(a ...interface{}) {
	z.logger.Error(getStr(a...))
}

// Debug level info msg
func (z *zapLog) Debug(a ...interface{}) {
	if dev.IsDev() {
		z.logger.Debug(getStr(a...))
	}
}

// Debug level info msg
func (z *zapLog) Debugf(msg string, a ...interface{}) {
	if dev.IsDev() {
		z.logger.Debug(fmt.Sprintf(msg, a...))
	}
}

//Fatal 系统级错误
func (z *zapLog) Fatal(a ...interface{}) {
	z.logger.Fatal(getStr(a...))
	os.Exit(1)
}

//Fatalf 系统级错误
func (z *zapLog) Fatalf(msg string, a ...interface{}) {
	z.logger.Fatal(fmt.Sprintf(msg, a...))
	os.Exit(1)
}

//JSON json输出
func (z *zapLog) JSON(a ...interface{}) {
	for _, v := range a {
		b, _ := json.MarshalIndent(v, "", "     ")
		z.logger.Info(string(b))
	}
}

// TraceError return trace of error
func (z *zapLog) TraceError(err error) error {
	e := errors.Cause(err) //获取原始对象
	z.logger.Error("Cause", zap.Error(e))
	return errors.WithStack(err)
}

// Close close the logger
func (z *zapLog) Close() {
}
