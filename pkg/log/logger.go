package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var (
	Info  func(context.Context, string, ...interface{})
	Warn  func(context.Context, string, ...interface{})
	Error func(context.Context, string, ...interface{})
	Sync  func()
)

func init() {
	fmt.Println("init logger ...")
	defaultLogger := newLogger()
	Info = defaultLogger.Info
	Warn = defaultLogger.Warn
	Error = defaultLogger.Error
	Sync = defaultLogger.Sync
	fmt.Println("init logger done.")
}

type Logger interface {
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Sync()
}

type zapLogger struct {
	zapSugar *zap.SugaredLogger
}

func TraceId(ctx context.Context) string {
	traceId, _ := ctx.Value("trace_id").(string)
	return traceId
}

func (zl *zapLogger) Info(ctx context.Context, template string, args ...interface{}) {
	zl.zapSugar.Infof(template, args...)
}

func (zl *zapLogger) Warn(ctx context.Context, template string, args ...interface{}) {
	zl.zapSugar.Warnf(template, args...)
}

func (zl *zapLogger) Error(ctx context.Context, template string, args ...interface{}) {
	zl.zapSugar.Errorf(template, args...)
}

func (zl *zapLogger) Sync() {
	zl.zapSugar.Sync()
}

func newLogger() Logger {
	format := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "name",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000"),
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	fileRotate := &lumberjack.Logger{
		Filename: "logs/go2.log",
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
