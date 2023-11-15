package types

import sc "github.com/LimeChain/goscale"

// typeError is a custom type error specific to
// the Runtime and expected to be returned to the Host
type typeError struct {
	message string
}

func newTypeError(msg string) error {
	return typeError{message: msg}
}

func (e typeError) Error() string {
	return "not a valid '" + e.message + "' type"
}

// UnexpectedError is a error type that is used to wrap unexpected errors that should panic
type UnexpectedError struct {
	sc.Str
}

func NewUnexpectedError(err error) UnexpectedError {
	return UnexpectedError{sc.Str(err.Error())}
}

func (e UnexpectedError) Error() string {
	return string(e.Str)
}
