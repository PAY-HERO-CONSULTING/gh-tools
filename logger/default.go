package logger

import (
	"fmt"
	"sync"

	"github.com/PAY-HERO-CONSULTING/gh-tools/null"
)

type environment string

const (
	EnvProduction  environment = "production"
	EnvDevelopment environment = "development"
)

var loggers []logger

type logConfig struct {
	DisableFileLogger  bool
	OutputFile         *string
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
	wg                 *sync.WaitGroup
}

func (e environment) String() string {
	return string(e)
}

func NewLoggerConfig(
	wg *sync.WaitGroup,
) logConfig {
	return logConfig{
		Environment: EnvDevelopment.String(),
		wg:          wg,
	}
}

func (l *logConfig) Initialize() error {
	if l.isProduction() {
		if len(null.ValueFromNull(l.OutputFile)) < 1 {
			return fmt.Errorf("missing output file")
		}
	}

	return l.setupLoggers()
}

func (l *logConfig) isProduction() bool {
	return environment(l.Environment) == EnvProduction
}

func (l logConfig) setupLoggers() error {
	logger := NewAppLogger(l.wg, l)
	loggers = append(loggers, logger)

	if l.isProduction() {
		sentryLogger := l.newSentryLogger()
		loggers = append(loggers, sentryLogger)
	}

	return nil
}
