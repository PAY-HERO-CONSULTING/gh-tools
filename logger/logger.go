package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
	"github.com/rs/zerolog"
)

type logger interface {
	AppError(ctx context.Context, err *apperrs.Error)
	Debug(msg string)
	Debugf(format string, args ...interface{})
	Error(msg string)
	Errorf(format string, args ...interface{})
	Fatal(msg string)
	Fatalf(format string, args ...interface{})
	// Flush() error
	Info(msg string)
	Infof(format string, args ...interface{})
	Panic(msg string)
	Panicf(format string, args ...interface{})
	// Valid() bool
	Warn(msg string)
	Warnf(format string, args ...interface{})
}

type appLogger struct {
	logger zerolog.Logger
	config logConfig
}

func NewAppLogger(wg *sync.WaitGroup, config logConfig) *appLogger {
	logger := config.newZeroLogger()

	// intialize the hook here
	logger = logger.Hook(NewSlackHook(wg))
	// logger.

	appLogger := &appLogger{
		logger: logger,
	}

	loggers = append(loggers, appLogger)

	return appLogger
}

func (l *appLogger) AppError(ctx context.Context, err *apperrs.Error) {
	if err.SendNotification() {
		l.logger.Level(zerolog.ErrorLevel)
		l.Error(err.LogMessages())
	} else {
		l.Warn(err.LogMessages())
	}
}

func (l *appLogger) Error(msg string) {
	l.logger.Error().Str("service", l.config.Service).Msg(msg)
}

func (l *appLogger) Errorf(msg string, values ...interface{}) {
	l.logger.Error().Str("service", l.config.Service).Msgf(msg, values...)
}

func (l *appLogger) ErrorWithPayload(msg string, payload any, values ...interface{}) {
	msg = msg + " body = [%v]"
	bytesBody := l.getBytes(payload)
	values = append(values, string(bytesBody))
	l.logger.Error().Str("service", l.config.Service).Bytes("Body", bytesBody).Msgf(msg, values...)
}

func (l *appLogger) Panic(msg string) {
	l.logger.Panic().Str("service", l.config.Service).Msg(msg)
}

func (l *appLogger) Panicf(msg string, values ...interface{}) {
	l.logger.Panic().Str("service", l.config.Service).Msgf(msg, values...)
}

func (l *appLogger) PanicWithPayload(msg string, payload any, values ...interface{}) {
	l.logger.Panic().Str("service", l.config.Service).Bytes("Body", l.getBytes(payload)).Msgf(msg, values...)
}

func (l *appLogger) Fatal(msg string) {
	l.logger.Fatal().Str("service", l.config.Service).Msg(msg)
}

func (l *appLogger) Fatalf(msg string, values ...interface{}) {
	l.logger.Fatal().Str("service", l.config.Service).Msgf(msg, values...)
}

func (l *appLogger) FatalWithPayload(msg string, payload any, values ...interface{}) {
	l.logger.Fatal().Str("service", l.config.Service).Bytes("Body", l.getBytes(payload)).Msgf(msg, values...)
}

func (l *appLogger) Info(msg string) {
	l.logger.Info().Str("service", l.config.Service).Msg(msg)
}

func (l *appLogger) Infof(msg string, values ...interface{}) {
	l.logger.Info().Str("service", l.config.Service).Msgf(msg, values...)
}

func (l *appLogger) InfoWithPayload(msg string, payload any, values ...interface{}) {
	msg = msg + " body = [%v]"
	bytesBody := l.getBytes(payload)
	values = append(values, string(bytesBody))
	l.logger.Info().Str("service", l.config.Service).Bytes("Body", l.getBytes(bytesBody)).Msgf(msg, values...)
}

func (l *appLogger) Warn(msg string) {
	l.logger.Warn().Str("service", l.config.Service).Msg(msg)
}

func (l *appLogger) Warnf(msg string, values ...interface{}) {
	l.logger.Warn().Str("service", l.config.Service).Msgf(msg, values...)
}

func (l *appLogger) WarnWithPayload(msg string, payload any, values ...interface{}) {
	l.logger.Warn().Str("service", l.config.Service).Bytes("Body", l.getBytes(payload)).Msgf(msg, values...)
}

func (l *appLogger) Debug(msg string) {
	l.logger.Debug().Str("service", l.config.Service).Msg(msg)
}

func (l *appLogger) Debugf(msg string, values ...interface{}) {
	l.logger.Debug().Str("service", l.config.Service).Msgf(msg, values...)
}

func (l *appLogger) DebugWithPayload(msg string, payload any, values ...interface{}) {
	l.logger.Debug().Str("service", l.config.Service).Bytes("Body", l.getBytes(payload)).Msgf(msg, values...)
}

func (l *appLogger) getBytes(data any) []byte {
	body, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		log.Printf("failed to get bytes for data: [%+v]", data)
		return nil
	}

	return body
}

func buildMessage(msg string, values ...interface{}) string {
	return fmt.Sprintf(msg, values...)
}
