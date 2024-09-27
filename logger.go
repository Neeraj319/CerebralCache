package main

import (
	"go.uber.org/zap"
	"os"
)

func GetLogger() *zap.Logger {
	logger := zap.Must(zap.NewProduction())
	if os.Getenv("ENV") == "develoment" {
		logger = zap.Must(zap.NewDevelopment())
	}
	zap.ReplaceGlobals(logger)
	return logger
}
