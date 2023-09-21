package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path/filepath"
	"sync"
	"time"
)

var (
	Logger   *zap.Logger
	lock     sync.Mutex
	logLevel string
)

func InitLogger(level string) {
	logLevel = level
	buildLogger()
	go rotateLogs()
}

func buildLogger() {
	lock.Lock()
	defer lock.Unlock()
	var level zapcore.Level
	switch logLevel {
	case "error":
		level = zap.ErrorLevel
	case "debug":
		level = zap.DebugLevel
	case "warn":
		level = zap.WarnLevel
	case "info":
		level = zap.InfoLevel
	default:
		level = zap.InfoLevel
	}
	zapConfig := zap.NewProductionConfig()
	zapConfig.Encoding = "console"
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	zapConfig.OutputPaths = []string{"stdout", getLogFileName("info")}
	zapConfig.ErrorOutputPaths = []string{"stderr", getLogFileName("error")}
	zapConfig.Level = zap.NewAtomicLevelAt(level)
	var err error
	Logger, err = zapConfig.Build()
	if err != nil {
		panic("log failed" + err.Error())
	}
}

// 根据日期获得日志文件
func getLogFileName(name string) string {
	return filepath.Join("logs", time.Now().Format("2006-01-02")+"-"+name+".log")
}

func rotateLogs() {
	for {
		now := time.Now()
		nextMidnight := time.Date(
			now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location(),
		).Add(24 * time.Hour)
		durationUntilMidnight := nextMidnight.Sub(now)

		time.Sleep(durationUntilMidnight)
		buildLogger()
	}
}
