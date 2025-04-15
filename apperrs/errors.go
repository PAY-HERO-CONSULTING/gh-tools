package apperrs

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PAY-HERO-CONSULTING/gh-tools/ctxutils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Type holds a type string and integer code for the error
type Type string

// "Set" of valid errorTypes
const (
	Authorization        Type = "AUTHORIZATION" // Authentication Failures -
	BadRequest           Type = "BAD_REQUEST"   // Validation errors / BadInput
	Conflict             Type = "CONFLICT"      // Already exists (eg, create account with existent email) - 409
	Internal             Type = "INTERNAL"      // Server (500) and fallback errors
	Permission           Type = "PERMISSION_DENIED"
	NotFound             Type = "NOT_FOUND"              // For not finding resource
	PayloadTooLarge      Type = "PAYLOAD_TOO_LARGE"      // for uploading tons of JSON, or an image over the limit - 413
	ServiceUnavailable   Type = "SERVICE_UNAVAILABLE"    // For long running handlers
	UnsupportedMediaType Type = "UNSUPPORTED_MEDIA_TYPE" // for http 415
	UnexpextedError      Type = "UNEXPECTED_ERROR"
	DatabaseError        Type = "DATABASE_ERROR"
	InvalidArgument      Type = "INVALID_ARGUMENT"
)

// Error holds a custom error for the application
// which is helpful in returning a consistent
// error type/message from API endpoints
type Error struct {
	err         error
	Type        Type   `json:"type"`
	Service     string `json:"service"`
	Message     string `json:"message"`
	RequestID   string `json:"request_id"`
	IPAddress   string `json:"ip_address"`
	MACAddress  string `json:"mac_address"`
	Payload     any    `json:"payload"`
	logMessages []string
	variables   []any
	notify      bool
}

// Error satisfies standard error interface
// we can return errors from this package as
// a regular old go _error_
func (e *Error) Err() error {
	return e.err
}

func (e *Error) Error() string {
	errorStr := e.Message
	if errorStr != "" {
		if e.err == nil {
			errorStr += ", err: " + "nil error"
		} else {
			errorStr += ", err: " + e.err.Error()
		}
	} else {
		if e == nil || e.err == nil {
			errorStr = "nil error"
		} else {
			errorStr = e.err.Error()
		}
	}

	if e.notify {
		errorStr = fmt.Sprintf("[ERROR] %s", errorStr)
	}

	return errorStr
}

func Wrap(err error) *Error {
	return New(err, Internal)
}

// message
func (e *Error) GetMessage() string {
	return e.Message
}

func (e *Error) SendNotification() bool {
	return e.notify
}

func (e *Error) Notify() *Error {
	e.notify = true
	return e
}

func (e *Error) LogMessages() string {
	messages := []string{}

	if len(e.logMessages) > 0 {
		for _, logMessage := range e.logMessages {
			messages = append(messages, "\t"+logMessage)
		}
	}

	return strings.Join(messages, "\n")
}

func (e *Error) GetLogMessageStack() []string {
	return e.logMessages
}

func (e *Error) SetErrorType(errorType Type) *Error {
	e.Type = errorType
	return e
}

func (e *Error) AddLogMessage(message string) *Error {
	e.logMessages = append([]string{message}, e.logMessages...)
	return e
}

func (e *Error) AddLogMessagef(format string, args ...interface{}) *Error {
	return e.AddLogMessage(fmt.Sprintf(format, args...))
}

func (e *Error) SetErrorMessage(message string) *Error {
	if len(message) > 0 {
		e.Message = message
	}

	return e
}

func (e *Error) SetErrorMessagef(format string, args ...interface{}) *Error {
	return e.SetErrorMessage(fmt.Sprintf(format, args...))
}

func (e *Error) SetVariables(a ...interface{}) *Error {
	e.variables = a
	return e
}

// Status is a mapping errors to status codes
// Of course, this is somewhat redundant since
// our errors already map http status codes
func (e *Error) Status() int {
	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	case PayloadTooLarge:
		return http.StatusRequestEntityTooLarge
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	case UnsupportedMediaType:
		return http.StatusUnsupportedMediaType
	case UnexpextedError:
		return http.StatusExpectationFailed
	case DatabaseError:
		return http.StatusBadGateway
	case InvalidArgument:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func (e *Error) GRPCStatus() codes.Code {
	switch e.Type {
	case Authorization:
		return codes.Unauthenticated
	case BadRequest:
		return codes.InvalidArgument
	case Conflict:
		return codes.AlreadyExists
	case Internal:
		return codes.Internal
	case NotFound:
		return codes.NotFound
	case PayloadTooLarge:
		return codes.OutOfRange
	case Permission:
		return codes.PermissionDenied
	case ServiceUnavailable:
		return codes.Unimplemented
	case UnsupportedMediaType:
		return codes.Unknown
	case UnexpextedError:
		return codes.Unavailable
	case DatabaseError:
		return codes.Internal
	default:
		return codes.Internal
	}
}

func (e *Error) AddLogMessages(message string) *Error {
	e.logMessages = append([]string{message}, e.logMessages...)
	return e
}

func (e *Error) AddIPAddress(ipAddress string) *Error {
	e.IPAddress = ipAddress
	return e
}

func (e *Error) AddRequestID(requestID string) *Error {
	e.RequestID = requestID
	return e
}

func (e *Error) AddPayload(payload any) *Error {
	e.Payload = payload
	return e
}

func (e *Error) AddMacAddress(macAddress string) *Error {
	e.MACAddress = macAddress
	return e
}

func (e *Error) SetMessage(message string) *Error {
	e.Message = message
	return e
}

func (e *Error) SetContextVariables(ctx context.Context) *Error {
	return e.AddIPAddress(
		ctxutils.IPAddress(ctx),
	).AddMacAddress(
		ctxutils.MacAddress(ctx),
	).AddRequestID(
		ctxutils.RequestId(ctx),
	)
}

// Status checks the runtime type
// of the error and returns an http
// status code if the error is model.Error
func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

/*
* Error "Factories"
 */

// NewAuthorization to create a 401
func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

// NewBadRequest to create 400 errors (validation, for example)
func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: reason,
		err:     fmt.Errorf(reason),
	}
}

// NewConflict to create an error for 409
func NewConflict(name string, value string) *Error {
	return &Error{
		Type:    Conflict,
		Message: fmt.Sprintf("resource: %v with value: %v already exists", name, value),
		err:     fmt.Errorf("resource: %v with value: %v already exists", name, value),
	}
}

// NewInternal for 500 errors and unknown errors
func NewInternal(message string) *Error {
	return &Error{
		Type:    Internal,
		Message: message,
		err:     fmt.Errorf("%v", message),
	}
}

// NewInvalidArgument to create an error for 400
func NewInvalidArgument(reason string) *Error {
	return &Error{
		Type:    InvalidArgument,
		Message: reason,
		err:     fmt.Errorf("%v", reason),
	}
}

func NewPermission(reason string) *Error {
	return &Error{
		Type:    Permission,
		Message: reason,
		err:     fmt.Errorf("%v", reason),
	}
}

// NewNotFound to create an error for 404
func NewNotFound(name string, value string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("%v not found", name),
		err:     fmt.Errorf("%v not found", name),
	}
}

// NewPayloadTooLarge to create an error for 413
func NewPayloadTooLarge(maxBodySize int64, contentLength int64) *Error {
	return &Error{
		Type:    PayloadTooLarge,
		Message: fmt.Sprintf("Max payload size of %v exceeded. Actual payload size: %v", maxBodySize, contentLength),
		err:     fmt.Errorf("max payload size of %v exceeded. actual payload size: %v", maxBodySize, contentLength),
	}
}

// NewServiceUnavailable to create an error for 503
func NewServiceUnavailable(name string) *Error {
	return &Error{
		Type:    ServiceUnavailable,
		Message: fmt.Sprintf("service %s unavailable or timed out", name),
		err:     fmt.Errorf("service %s unavailable or timed out", name),
	}
}

// NewUnsupportedMediaType to create an error for 415
func NewUnsupportedMediaType(reason string) *Error {
	return &Error{
		Type:    UnsupportedMediaType,
		Message: reason,
		err:     fmt.Errorf("%v", reason),
	}
}

func NewUnexpectedError(reason string) *Error {
	return &Error{
		Type:    UnexpextedError,
		Message: reason,
		err:     fmt.Errorf("%v", reason),
	}
}

func NewError(err error) *Error {
	appError, ok := err.(*Error)
	if !ok {
		return NewErrorWithType(
			err,
			Internal,
		)
	}

	return appError
}

func NewErrorWithType(
	err error,
	errorType Type,
) *Error {
	if errorType == "" {
		errorType = Internal
	}

	appError, ok := err.(*Error)
	if !ok {
		appError = &Error{
			Message: err.Error(),
			Type:    errorType,
		}
	}

	appError.Type = errorType

	return appError
}

func (e *Error) LogErrorMessage(message string, values ...interface{}) error {
	// logger.Errorf(message, values...)
	return e
}

func (e *Error) LogErrorMessagePayload(err error, payload interface{}) error {
	_, ok := payload.(*[]byte)
	if !ok {
		// logger.Errorf(e.Error())
		return e
	} else {
		// logger.ErrorWithPayload(e.Error(), *bytesPayload)
		return e
	}
}

func HTTPError(
	err error,
) *Error {
	return getHTTPError(status.Convert(err))
}

func ToHTTPError(
	err error,
) *Error {
	appError, ok := err.(*Error)
	if !ok {
		appError = &Error{
			Message: err.Error(),
			Type:    Internal,
		}
	}

	return appError
}

func getHTTPError(
	status *status.Status,
) *Error {
	appErr := New(
		errors.New(status.Message()),
		Internal,
	)

	switch status.Code() {
	case codes.Unauthenticated:
		appErr.Type = Permission
	case codes.InvalidArgument:
		appErr.Type = BadRequest
	case codes.AlreadyExists:
		appErr.Type = Conflict
	case codes.Internal:
		appErr.Type = Internal
	case codes.NotFound:
		appErr.Type = NotFound
	case codes.OutOfRange:
		appErr.Type = PayloadTooLarge
	case codes.PermissionDenied:
		appErr.Type = Permission
	case codes.Unimplemented:
		appErr.Type = ServiceUnavailable
	case codes.Unknown:
		appErr.Type = UnsupportedMediaType
	case codes.Unavailable:
		appErr.Type = UnexpextedError
	default:
		appErr.Type = Internal
	}

	appErr.Message = status.Message()

	return appErr
}

func (e *Error) LogInfoMessage(service, message string, err error) error {
	// logger.Infof(service, message, err)
	return e
}

func LogErrorMessage(service, message string, err error) {
	// logger.Errorf(service, message, err)
}

func (e *Error) JsonResponse() map[string]interface{} {
	if e.Message == "" {
		e.Message = "Failed to perform request. Please try again."
	}

	jsonResponse := map[string]interface{}{
		"error_code":    e.Type,
		"error_message": e.Message,
		"status_code":   e.Status(),
	}

	return jsonResponse
}
