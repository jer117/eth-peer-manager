package utils

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap/zapcore"
	"time"
)

var Client = resty.New().SetTimeout(5 * time.Second)

func SetLogLevel(logLevel string) zapcore.Level {
	switch logLevel {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	}
	return zapcore.InfoLevel
}
