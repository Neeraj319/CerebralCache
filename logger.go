package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

func GetLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	if os.Getenv("ENV") == "develoment" {
		config = zap.NewDevelopmentConfig()
	}
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	config.Encoding = "console"
	logger := zap.Must(config.Build())
	zap.ReplaceGlobals(logger)
	return logger
}
