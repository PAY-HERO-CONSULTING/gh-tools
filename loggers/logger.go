package logger

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
	"github.com/rs/zerolog"
)

type Logger interface {
	AppError(err *apperrs.Error)
	Error(msg string)
	Errorf(msg string, values ...interface{})
	ErrorWithPayload(service string, payload any, values ...interface{})
	Panic(msg string)
	Panicf(msg string, values ...interface{})
	PanicWithPayload(msg string, payload any, values ...interface{})
	Fatal(msg string)
	Fatalf(msg string, values ...interface{})
	FatalWithPayload(msg string, payload any, values ...interface{})
	Info(msg string)
	Infof(msg string, values ...interface{})
	InfoWithPayload(msg string, payload any, values ...interface{})
	Warn(msg string)
	Warnf(msg string, values ...interface{})
	WarnWithPayload(msg string, payload any, values ...interface{})
	Debug(msg string)
	Debugf(msg string, values ...interface{})
	DebugWithPayload(msg string, payload any, values ...interface{})
}

type appLogger struct {
	logger zerolog.Logger
	config *HookConfig
}

func NewAppLogger(wg *sync.WaitGroup, config *HookConfig) *appLogger {
	logger := newZeroLogger(config)

	// intialize the hook here
	logger = logger.Hook(NewTelegramHook(wg, config))
	// logger.

	return &appLogger{
		logger: logger,
		config: config,
	}
}

func (l *appLogger) AppError(err *apperrs.Error) {
	if err.SendNotification() {
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
