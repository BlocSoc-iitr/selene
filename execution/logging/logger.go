package logging

import (
    "go.uber.org/zap"
)

var logger *zap.Logger

func InitLogger() {
    logger, _ = zap.NewProduction()
}

func Trace(msg string, fields ...zap.Field) {
    logger.Debug(msg, fields...)
}