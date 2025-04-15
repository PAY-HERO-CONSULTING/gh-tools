package logger

import (
	"sync"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
)

var (
	logger      *appLogger
	docsMesaage = "has successfully uploaded  documents for account"
)

type HookConfig struct {
	InfoBotToken       string
	ErrorBotToken      string
	PanicBotToken      string
	ChannelID          int64
	UniversalChannelID int64
	DocsChatID         int64
	Service            string
	LogFile            string
	Environment        string
	LogLevel           string
}

func NewLogger(wg *sync.WaitGroup, config *HookConfig) {
	logger = NewAppLogger(wg, config)
}

func AppError(err *apperrs.Error) {
	logger.AppError(err)
}

func Error(msg string) {
	logger.Error(msg)
}

func Errorf(msg string, values ...interface{}) {
	logger.Errorf(msg, values...)
}

func ErrorWithPayload(msg string, payload any, values ...interface{}) {
	logger.ErrorWithPayload(msg, payload, values...)
}

func Warn(msg string) {
	logger.Warn(msg)
}

func Warnf(msg string, values ...interface{}) {
	logger.Warnf(msg, values...)
}

func WarnWithPayload(msg string, payload any, values ...interface{}) {
	logger.WarnWithPayload(msg, payload, values...)
}

func Info(msg string) {
	logger.Info(msg)
}

func Infof(msg string, values ...interface{}) {
	logger.Infof(msg, values...)
}

func InfoWithPayload(msg string, payload any, values ...interface{}) {
	logger.ErrorWithPayload(msg, payload, values...)
}

func Fatal(msg string) {
	logger.Fatal(msg)
}

func Fatalf(msg string, values ...interface{}) {
	logger.Fatalf(msg, values...)
}

func FatalWithPayload(msg string, payload any, values ...interface{}) {
	logger.FatalWithPayload(msg, payload, values...)
}

func Panic(msg string) {
	logger.Panic(msg)
}

func Panicf(msg string, values ...interface{}) {
	logger.Panicf(msg, values...)
}

func PanicWithPayload(msg string, payload any, values ...interface{}) {
	logger.PanicWithPayload(msg, payload, values...)
}

func Debug(msg string) {
	logger.Debug(msg)
}

func Debugf(msg string, values ...interface{}) {
	logger.Debugf(msg, values...)
}

func DebugWithPayload(msg string, payload any, values ...interface{}) {
	logger.DebugWithPayload(msg, payload, values...)
}

func SetService(service string) {
	logger.config.Service = service
}
