package errors

import (
	"errors"
	"fmt"
)

// ErrUnknownCode is a string code representing an unknown error
// This will be used when no error code is sent by the handler
const ErrUnknownCode = "PIT-000"

// ErrInternalCode is a string code representing an internal Pitaya error
const ErrInternalCode = "PIT-500"

// ErrNotFoundCode is a string code representing a not found related error
const ErrNotFoundCode = "PIT-404"

// ErrBadRequestCode is a string code representing a bad request related error
const ErrBadRequestCode = "PIT-400"

// ErrClientClosedRequest is a string code representing the client closed request error
const ErrClientClosedRequest = "PIT-499"

// Error is an error with a code, message and metadata
type Error struct {
	Code     string
	Message  string
	Metadata map[string]string
}

//NewError ctor
func NewError(err error, code string, metadata ...map[string]string) *Error {
	if cstmErr, ok := err.(*Error); ok {
		if len(metadata) > 0 {
			mergeMetadatas(cstmErr, metadata[0])
		}
		return cstmErr
	}

	e := &Error{
		Code:    code,
		Message: err.Error(),
	}
	if len(metadata) > 0 {
		e.Metadata = metadata[0]
	}
	return e

}

func (e *Error) Error() string {
	return e.Message
}

func mergeMetadatas(pitayaErr *Error, metadata map[string]string) {
	if pitayaErr.Metadata == nil {
		pitayaErr.Metadata = metadata
		return
	}

	for key, value := range metadata {
		pitayaErr.Metadata[key] = value
	}
}

// CodeFromError returns the code of error.
// If error is nil, return empty string.
// If error is not a pitaya error, returns unkown code
func CodeFromError(err error) string {
	if err == nil {
		return ""
	}

	pitayaErr, ok := err.(*Error)
	if !ok {
		return ErrUnknownCode
	}

	if pitayaErr == nil {
		return ""
	}

	return pitayaErr.Code
}

func NewErr(format string, args ...interface{}) error {
	return errors.New(fmt.Sprintf(format, args...))
}
