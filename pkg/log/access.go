package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	// access logger
	Access func(AccessLog)
)

func init() {
	fmt.Println("init access logger ...")
	accessLogger := newAccessLogger()
	Access = accessLogger.Access
	fmt.Println("init access logger done.")
}

type AccessLogger interface {
	Access(AccessLog)
}

type AccessLog struct {
	RemoteAddr           string
	Method               string
	Proto                string
	RequestContentLength int64
	Host                 string
	RequestURI           string
	Status               int
	Url                  string
	UserAgent            string
	RequestTime          int64
}

func (zl *zapLogger) Access(accessLog AccessLog) {
	zl.zapSugar.Infof("%+v", accessLog)
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
