package types

// typeError is a custom type error specific to
// the Runtime and expected to be returned to the Host
type typeError struct {
	message string
}

func newTypeError(msg string) error {
	return typeError{message: msg}
}

func (err typeError) Error() string {
	return "not a valid '" + err.message + "' type"
}
