package ecode

import (
	"errors"
	"fmt"
)

const (
	// UnknownCode is unknown code for error info.
	UnknownCode = -1
	// UnknownHttpCode is unknown reason for error info.
	UnknownHttpCode = 500
)

// Error is a status error.
type Error struct {
	Code     int // 业务错误码
	HttpCode int
	Message  string
	cause    error
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code = %d httpcode = %d message = %s cause = %v", e.Code, e.HttpCode,
		e.Message, e.cause)
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *Error) Unwrap() error { return e.cause }

// Is matches each error in the chain with the target value.
func (e *Error) Is(err error) bool {
	if se := new(Error); errors.As(err, &se) {
		return se.Code == e.Code && se.HttpCode == e.HttpCode
	}
	return false
}

// WithCause with the underlying cause of the error.
func (e *Error) WithCause(cause error) *Error {
	newErr := New(e.Code, e.HttpCode, e.Message)
	newErr.cause = cause
	return newErr
}

// New returns an error object for the code, message.
func New(code, httpCode int, message string) *Error {
	return &Error{
		Code:     code,
		HttpCode: httpCode,
		Message:  message,
	}
}

// Newf New(code fmt.Sprintf(format, a...))
func Newf(code, httpCode int, format string, a ...interface{}) *Error {
	return New(code, httpCode, fmt.Sprintf(format, a...))
}

// Errorf returns an error object for the code, message and error info.
func Errorf(code, httpCode int, format string, a ...interface{}) error {
	return New(code, httpCode, fmt.Sprintf(format, a...))
}

// Code returns the http code for an error.
// It supports wrapped errors.
func Code(err error) int {
	if err == nil {
		return UnknownCode
	}
	return FromError(err).Code
}

func HttpCode(err error) int {
	if err == nil {
		return UnknownHttpCode
	}
	return FromError(err).HttpCode
}

// FromError try to convert an error to *Error.
// It supports wrapped errors.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if se := new(Error); errors.As(err, &se) {
		return se
	}
	return New(UnknownCode, UnknownHttpCode, err.Error())
}
