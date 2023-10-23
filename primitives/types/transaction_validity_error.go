package types

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	errInvalidTransactionValidityErrorType = "invalid TransactionValidityError type"
)

const (
	TransactionValidityErrorInvalidTransaction sc.U8 = iota
	TransactionValidityErrorUnknownTransaction
)

// TransactionValidityError Errors that can occur while checking the validity of a transaction.
type TransactionValidityError sc.VaryingData

func NewTransactionValidityError(value sc.Encodable) TransactionValidityError {
	// InvalidTransaction = 0 - Transaction is invalid.
	// UnknownTransaction = 1 - Transaction validity can’t be determined.
	switch value.(type) {
	case InvalidTransaction, UnknownTransaction:
	default:
		log.Critical(errInvalidTransactionValidityErrorType)
	}

	return TransactionValidityError(sc.NewVaryingData(value))
}

func (e TransactionValidityError) Encode(buffer *bytes.Buffer) {
	value := e[0]

	switch reflect.TypeOf(value) {
	case reflect.TypeOf(*new(InvalidTransaction)):
		TransactionValidityErrorInvalidTransaction.Encode(buffer)
	case reflect.TypeOf(*new(UnknownTransaction)):
		TransactionValidityErrorUnknownTransaction.Encode(buffer)
	default:
		log.Critical(errInvalidTransactionValidityErrorType)
	}

	value.Encode(buffer)
}

func DecodeTransactionValidityError(buffer *bytes.Buffer) TransactionValidityError {
	b := sc.DecodeU8(buffer)

	switch b {
	case TransactionValidityErrorInvalidTransaction:
		value := DecodeInvalidTransaction(buffer)
		return NewTransactionValidityError(value)
	case TransactionValidityErrorUnknownTransaction:
		value := DecodeUnknownTransaction(buffer)
		return NewTransactionValidityError(value)
	default:
		log.Critical(errInvalidTransactionValidityErrorType)
	}

	panic("unreachable")
}

func (e TransactionValidityError) Bytes() []byte {
	return sc.EncodedBytes(e)
}
