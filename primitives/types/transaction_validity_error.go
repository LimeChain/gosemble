package types

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
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

func NewTransactionValidityError(value sc.Encodable) (TransactionValidityError, error) {
	// InvalidTransaction = 0 - Transaction is invalid.
	// UnknownTransaction = 1 - Transaction validity can’t be determined.
	switch value.(type) {
	case InvalidTransaction, UnknownTransaction:
	default:
		return TransactionValidityError{}, errInvalidTransactionValidityErrorType
	}
	return TransactionValidityError(sc.NewVaryingData(value)), nil
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
		log.Critical(errInvalidTransactionValidityErrorType.Error())
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
		return NewTransactionValidityError(value)
	case TransactionValidityErrorUnknownTransaction:
		value, err := DecodeUnknownTransaction(buffer)
		if err != nil {
			return TransactionValidityError{}, err
		}
		return NewTransactionValidityError(value)
	default:
		return TransactionValidityError{}, errInvalidTransactionValidityErrorType
	}
}

func (e TransactionValidityError) Bytes() []byte {
	return sc.EncodedBytes(e)
}
