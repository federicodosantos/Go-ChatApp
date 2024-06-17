package log

import (
	"strings"

	"go.uber.org/zap"
)

func NewLogger(env string) *zap.Logger {
	var logger *zap.Logger

	if strings.ToLower(env) == "production" {
		logger = zap.Must(zap.NewProduction())
	} else {
		logger = zap.Must(zap.NewDevelopment())
	}

	return logger
}