package log

import (
	"fmt"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

const (
	accessLogTemplate = ""
)

var (
	// access logger
	Access func(...interface{})

	accessMap sync.Map
)

func init() {
	fmt.Println("init normal logger ...")
	accessLogger := newAccessLogger()
	Access = accessLogger.Access
	fmt.Println("init access logger done.")
}

type AccessLogger interface {
	Access(...interface{})
	Complete(...interface{})
}

func (zl *zapLogger) Access(args ...interface{}) {
	zl.zapSugar.Infof(accessLogTemplate, args...)
}

func (zl *zapLogger) Complete(args ...interface{}) {
	zl.zapSugar.Infof(accessLogTemplate, args...)
}

func newAccessLogger() AccessLogger {
	format := zapcore.EncoderConfig{
		MessageKey: "message",
		TimeKey:    "time",
		EncodeTime: zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
	}
	fileRotate := &lumberjack.Logger{
		Filename: "logs/access.log",
		Compress: false,
	}
	cores := []zapcore.Core{}
	cores = append(cores, zapcore.NewCore(
		zapcore.NewConsoleEncoder(format),
		zapcore.AddSync(fileRotate),
		zapcore.DebugLevel,
	))
	return &zapLogger{
		zapSugar: zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1)).Sugar(),
	}
}
