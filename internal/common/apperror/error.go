package apperror

import "go.uber.org/zap"

type Code string

const (
	CodeValidation Code = "VALIDATION"
	CodeNotFound   Code = "NOT_FOUND"
	CodeConflict   Code = "CONFLICT"
	CodeInternal   Code = "INTERNAL"
)

type Error struct {
	Code    Code
	Message string
	Err     error
}

func (e *Error) Error() string { return e.Message }
func (e *Error) Unwrap() error { return e.Err }
func (e *Error) Log(logger *zap.Logger) *Error {
	logger.Error(e.Message, zap.Error(e.Err))
	return e
}

func Validation(msg string) *Error {
	return &Error{Code: CodeValidation, Message: msg}
}

func NotFound(msg string) *Error {
	return &Error{Code: CodeNotFound, Message: msg}
}

func Conflict(msg string) *Error {
	return &Error{Code: CodeConflict, Message: msg}
}

func Internal(err error) *Error {
	return &Error{Code: CodeInternal, Message: "internal error", Err: err}
}
