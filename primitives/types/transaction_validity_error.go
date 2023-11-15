package types

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
)

var (
	errInvalidTransactionValidityErrorType = newTypeError("TransactionValidityError")
)

const (
	TransactionValidityErrorInvalidTransaction sc.U8 = iota
	TransactionValidityErrorUnknownTransaction
)

// TransactionValidityError Errors that can occur while checking the validity of a transaction.
type TransactionValidityError sc.VaryingData

// todo delete
func NewTransactionValidityErrorr(value sc.Encodable) (TransactionValidityError, error) {
	// InvalidTransaction = 0 - Transaction is invalid.
	// UnknownTransaction = 1 - Transaction validity can’t be determined.
	switch value.(type) {
	case InvalidTransaction, UnknownTransaction, UnexpectedError:
		return TransactionValidityError(sc.NewVaryingData(value)), nil
	default:
		return TransactionValidityError{NewUnexpectedError(errInvalidTransactionValidityErrorType)}, nil
	}
}

func NewTransactionValidityError(value sc.Encodable) TransactionValidityError {
	// InvalidTransaction = 0 - Transaction is invalid.
	// UnknownTransaction = 1 - Transaction validity can’t be determined.
	switch value.(type) {
	case InvalidTransaction, UnknownTransaction, UnexpectedError:
		return TransactionValidityError(sc.NewVaryingData(value))
	default:
		return TransactionValidityError{NewUnexpectedError(errInvalidTransactionValidityErrorType)}
	}
}

func (e TransactionValidityError) UnexpectedError() error {
	if len(e) == 0 {
		return nil
	}

	unexpectedErr, ok := e[0].(UnexpectedError)
	if !ok {
		return nil
	}

	return unexpectedErr
}

func (e TransactionValidityError) Error() string {
	if len(e) == 0 {
		return ""
	}

	switch e[0].(type) {
	case UnexpectedError:
		return e[0].(UnexpectedError).Error()
	case UnknownTransaction:
		return e[0].(UnknownTransaction).Error()
	case InvalidTransaction:
		return e[0].(InvalidTransaction).Error()
	default:
		return ""
	}
}

func (e TransactionValidityError) Encode(buffer *bytes.Buffer) error {
	value := e[0]

	switch reflect.TypeOf(value) {
	case reflect.TypeOf(*new(InvalidTransaction)):
		err := TransactionValidityErrorInvalidTransaction.Encode(buffer)
		if err != nil {
			return err
		}
	case reflect.TypeOf(*new(UnknownTransaction)):
		err := TransactionValidityErrorUnknownTransaction.Encode(buffer)
		if err != nil {
			return err
		}
	default:
		return errInvalidTransactionValidityErrorType
	}

	return value.Encode(buffer)
}

func DecodeTransactionValidityError(buffer *bytes.Buffer) (TransactionValidityError, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return TransactionValidityError{}, err
	}

	switch b {
	case TransactionValidityErrorInvalidTransaction:
		value, err := DecodeInvalidTransaction(buffer)
		if err != nil {
			return TransactionValidityError{}, err
		}
		tErr := NewTransactionValidityError(value)
		return tErr, tErr.UnexpectedError()
	case TransactionValidityErrorUnknownTransaction:
		value, err := DecodeUnknownTransaction(buffer)
		if err != nil {
			return TransactionValidityError{}, err
		}
		tErr := NewTransactionValidityError(value)
		return tErr, tErr.UnexpectedError()
	default:
		return TransactionValidityError{}, errInvalidTransactionValidityErrorType
	}
}

// todo delete
func DecodeTransactionValidityErrorr(buffer *bytes.Buffer) (TransactionValidityError, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return TransactionValidityError{}, err
	}

	switch b {
	case TransactionValidityErrorInvalidTransaction:
		value, err := DecodeInvalidTransaction(buffer)
		if err != nil {
			return TransactionValidityError{}, err
		}
		tErr, _ := NewTransactionValidityErrorr(value)
		return tErr, tErr.UnexpectedError()
	case TransactionValidityErrorUnknownTransaction:
		value, err := DecodeUnknownTransaction(buffer)
		if err != nil {
			return TransactionValidityError{}, err
		}
		tErr, _ := NewTransactionValidityErrorr(value)
		return tErr, tErr.UnexpectedError()
	default:
		return TransactionValidityError{}, errInvalidTransactionValidityErrorType
	}
}

func (e TransactionValidityError) Bytes() []byte {
	return sc.EncodedBytes(e)
}
