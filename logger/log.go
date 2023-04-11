package logger

import (
	"log"

	"go.uber.org/zap"
)

var logger *zap.Logger

func InitLogger() error {
	_logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	logger = _logger
	return nil
}

func L() *zap.Logger {
	if logger == nil {
		log.Fatal("logger not initial")
	}
	return logger
}

func S() *zap.SugaredLogger {
	if logger == nil {
		log.Fatal("logger not initial")
	}
	return logger.Sugar()
}
