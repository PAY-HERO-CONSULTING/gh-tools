package apperrs

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

var (
	ErrConnectionRefused           = errors.New("connection refused")
	ErrConnectionResetByPeer       = errors.New("connection reset by peer")
	ErrContextCanceled             = errors.New("context canceled")
	ErrContextDeadlineExceeded     = errors.New("context deadline exceeded")
	ErrDialTcpIOTimeout            = errors.New("dial tcp io timeout")
	ErrLockTimeout                 = errors.New("pq: canceling statement due to lock timeout")
	ErrNetwork                     = errors.New("network request failed")
	ErrNoCountryCode               = errors.New("no country code for phone number")
	ErrNoSurveyResponseOptionMatch = errors.New("no matching survey response option")
	ErrRequestCancelledByUser      = errors.New("pq: canceling statement due to user request")
)

// ErrorHandler handles an error.
type ErrorHandler interface {
	Handle(err error)
	HandleContext(ctx context.Context, err error)
}

// NoopErrorHandler is an error handler that discards every error.
type NoopErrorHandler struct{}

func (NoopErrorHandler) Handle(_ error)                           {}
func (NoopErrorHandler) HandleContext(_ context.Context, _ error) {}

func IsErrConnectionRefused(err error) bool {
	return isSpecificError(err, ErrConnectionRefused)
}

func IsErrConnectionResetByPeer(err error) bool {
	return isSpecificError(err, ErrConnectionResetByPeer)
}

func IsErrContextDeadlineExceeded(err error) bool {
	return isSpecificError(err, ErrContextDeadlineExceeded)
}

func IsErrDialTcpIOTimeout(err error) bool {
	return isSpecificError(err, ErrDialTcpIOTimeout)
}

func IsErrLockTimeout(err error) bool {
	return isSpecificError(err, ErrLockTimeout)
}

func IsErrNoCountryCode(err error) bool {
	return isSpecificError(err, ErrNoCountryCode)
}

func IsErrNetwork(err error) bool {
	return isSpecificError(err, ErrNetwork)
}

func IsErrContextCancelled(err error) bool {
	return isSpecificError(err, ErrContextCanceled)
}

func IsErrRequestCancelledByUser(err error) bool {
	return isSpecificError(err, ErrRequestCancelledByUser)
}

func isSpecificError(err error, specificErr error) bool {
	appError := Wrap(err)
	return appError.Err().Error() == specificErr.Error()
}

func IsNoRowsErr(err error) bool {
	appError := CastError(err, Internal)
	return appError.Error() == sql.ErrNoRows.Error()
}

func IsTokenExpired(err error) bool {
	appError := CastError(err, Authorization)
	return strings.Contains(appError.Error(), "token expired")
}
