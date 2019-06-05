package errorx

import "strings"

type Error interface {
	// Satisfy the generic error interface.
	error

	// Returns the short phrase depicting the classification of the error.
	Code() string

	// Returns the error details message.
	Message() string

	// Returns the original error if one was set.  Nil is returned if not set.
	OrigErr() error
}

type codeError struct {
	origError error
	code      string
	message   string
}

func (ce *codeError) Error() string {
	if ce.origError != nil {
		return ce.origError.Error()
	}
	return ""
}

func (ce *codeError) Code() string {
	return ce.code
}

func (ce *codeError) Message() string {
	if len(ce.message) > 0 {
		return ce.message
	}

	if ce.origError != nil {
		return ce.origError.Error()
	}
	return ""
}

func (ce *codeError) OrigErr() error {
	return ce.origError
}

//export
func NewError(code, message string) Error {
	return &codeError{
		code:    code,
		message: message,
	}
}

//export
func NewCodeError(err error, code string, message ...string) Error {
	ce := &codeError{
		origError: err,
		code:      code,
	}
	if len(message) > 0 {
		ce.message = strings.Join(message, "\n")
	}

	return ce
}
