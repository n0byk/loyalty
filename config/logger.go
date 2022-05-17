package config

import (
	"go.uber.org/zap"
)

func LoggerInit() *zap.Logger {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return logger
}
