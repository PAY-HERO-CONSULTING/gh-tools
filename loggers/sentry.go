package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/PAY-HERO-CONSULTING/gh-tools/apperrs"
	"github.com/PAY-HERO-CONSULTING/gh-tools/ctxutils"
	"github.com/getsentry/sentry-go"
)

type sentryLogger struct {
	ctx         context.Context
	environment string
	hub         *sentry.Hub
}

func newSentryLogger() *sentryLogger {
	return &sentryLogger{
		ctx:         context.Background(),
		environment: os.Getenv("ENVIRONMENT"),
		hub:         sentry.CurrentHub(),
	}
}

func (l *sentryLogger) AppError(ctx context.Context, err *apperrs.Error) {
	if !err.SendNotification() {
		return
	}

	l.withContext(ctx).captureError(err)
}

func (l *sentryLogger) Panic(msg string) {
	l.Error(msg)
}

func (l *sentryLogger) Panicf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *sentryLogger) Fatal(msg string) {
	l.Error(msg)
}

func (l *sentryLogger) Fatalf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *sentryLogger) Error(msg string) {
	l.captureMessage(msg)
}

func (l *sentryLogger) Errorf(format string, args ...interface{}) {
	l.Error(fmt.Sprintf(format, args...))
}

func (l *sentryLogger) Warn(msg string) {
}

func (l *sentryLogger) Warnf(format string, args ...interface{}) {
}

func (l *sentryLogger) Info(msg string) {
}

func (l *sentryLogger) Infof(format string, args ...interface{}) {
}

func (l *sentryLogger) Debug(msg string) {
}

func (l *sentryLogger) Debugf(format string, args ...interface{}) {
}

func (l *sentryLogger) Flush() error {
	return nil
}

func (l *sentryLogger) Valid() bool {
	return l.environment == "production" || l.environment == "staging"
}

func (l *sentryLogger) captureMessage(msg string) {
	event := sentry.NewEvent()
	event.Message = msg

	l.captureEvent(event)
}

func (l *sentryLogger) captureError(e *apperrs.Error) {
	event := sentry.NewEvent()
	event.Message = e.Err().Error()
	event.Exception = append(
		event.Exception,
		sentry.Exception{
			Stacktrace: sentry.ExtractStacktrace(e.Err()),
		},
	)

	for _, msg := range e.GetLogMessageStack() {
		event.Breadcrumbs = append(event.Breadcrumbs, &sentry.Breadcrumb{Message: msg})
	}

	l.captureEvent(event)
}

func (l *sentryLogger) captureEvent(event *sentry.Event) {
	event.Environment = l.environment
	ctx := l.ctx

	accountID := ctxutils.AccountID(ctx)
	if accountID != 0 {
		event.Extra["account_id"] = accountID
	}

	ipAddress := ctxutils.IPAddress(ctx)
	if ipAddress != "" {
		event.Extra["ip_address"] = ipAddress
	}

	requestID := ctxutils.RequestId(ctx)
	if requestID != "" {
		event.Extra["request_id"] = requestID
	}

	userID := ctxutils.TokenInfo(ctx).UserID
	if userID != "" {
		event.Extra["user_id"] = userID
	}

	sessionID := ctxutils.TokenInfo(ctx).SessionID
	if sessionID != "" {
		event.Extra["session_id"] = sessionID
	}

	userAgent := ctxutils.UserAgent(ctx)
	if userAgent != "" {
		event.Extra["user_agent"] = userAgent
	}

	l.hub.CaptureEvent(event)
}

func (l *sentryLogger) withContext(ctx context.Context) *sentryLogger {
	hub := l.hub
	if sentry.HasHubOnContext(ctx) {
		hub = sentry.GetHubFromContext(ctx)
	}

	return &sentryLogger{
		ctx:         ctx,
		environment: l.environment,
		hub:         hub,
	}
}
