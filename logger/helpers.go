package logger

import (
	"context"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
)

func AppError(ctx context.Context, err *apperrs.Error) {
	for _, logger := range loggers {
		logger.AppError(ctx, err)
	}
}

func Error(msg string) {
	for _, logger := range loggers {
		logger.Error(msg)
	}
}

func Errorf(msg string, values ...interface{}) {
	for _, logger := range loggers {
		logger.Errorf(msg, values...)
	}
}

func Warn(msg string) {
	for _, logger := range loggers {
		logger.Warn(msg)
	}
}

func Warnf(msg string, values ...interface{}) {
	for _, logger := range loggers {
		logger.Warnf(msg, values...)
	}
}

func Info(msg string) {
	for _, logger := range loggers {
		logger.Info(msg)
	}
}

func Infof(msg string, values ...interface{}) {
	for _, logger := range loggers {
		logger.Infof(msg, values...)
	}
}

func Fatal(msg string) {
	for _, logger := range loggers {
		logger.Fatal(msg)
	}
}

func Fatalf(msg string, values ...interface{}) {
	for _, logger := range loggers {
		logger.Fatalf(msg, values...)
	}
}

func Panic(msg string) {
	for _, logger := range loggers {
		logger.Panic(msg)
	}
}

func Panicf(msg string, values ...interface{}) {
	for _, logger := range loggers {
		logger.Panicf(msg, values...)
	}
}

func Debug(msg string) {
	for _, logger := range loggers {
		logger.Debug(msg)
	}
}

func Debugf(msg string, values ...interface{}) {
	for _, logger := range loggers {
		logger.Debugf(msg, values...)
	}
}
