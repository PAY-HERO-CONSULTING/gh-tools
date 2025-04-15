package apperrs

import "errors"

func CastError(err error, errorType Type) *Error {
	return New(err, errorType)
}

func New(err error, errorType Type) *Error {
	if errorType == "" {
		errorType = UnexpextedError
	}

	if err == nil {
		err = errors.New("nil error")
	}

	appError, ok := err.(*Error)
	if !ok {
		appError = &Error{
			err: err,
		}

		_ = appError.SetErrorType(errorType)
		_ = appError.AddLogMessage(err.Error())
	}

	return appError
}

func NewDatabaseError(err error) *Error {
	appError := CastError(err, Internal)

	if IsNoRowsErr(appError) ||
		IsErrContextCancelled(appError) ||
		IsErrRequestCancelledByUser(appError) {

		errorType := Internal

		if IsNoRowsErr(appError) {
			errorType = NotFound
		}

		_ = appError.SetErrorType(errorType)
	} else {
		_ = appError.SetErrorType(Internal)
		_ = appError.Notify()
	}

	return appError
}

func NewForbiddenError(err error, userID string) *Error {
	return Wrap(
		err,
	).SetErrorType(
		Permission,
	).AddLogMessagef(
		"user=[%v] does not have the required permission to do this action",
		userID,
	)
}

func NewInvalidArgumentError(errorMessage string) *Error {
	return New(
		errors.New(errorMessage),
		InvalidArgument,
	).SetErrorMessage(errorMessage)
}
