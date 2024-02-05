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

func NewTransactionValidityError(value sc.Encodable) error {
	// InvalidTransaction = 0 - Transaction is invalid.
	// UnknownTransaction = 1 - Transaction validity can’t be determined.
	switch value.(type) {
	case InvalidTransaction, UnknownTransaction:
	default:
		return errInvalidTransactionValidityErrorType
	}
	return TransactionValidityError(sc.NewVaryingData(value))
}

func (err TransactionValidityError) Error() string {
	if len(err) == 0 {
		return errInvalidTransactionValidityErrorType.Error()
	}

	switch err[0] {
	case TransactionValidityErrorUnknownTransaction:
		return err[1].(UnknownTransaction).Error()
	case TransactionValidityErrorInvalidTransaction:
		return err[1].(InvalidTransaction).Error()
	default:
		return errInvalidTransactionValidityErrorType.Error()
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
		err = NewTransactionValidityError(value)
		if txErr, ok := err.(TransactionValidityError); ok {
			return txErr, nil
		}
		return TransactionValidityError{}, err
	case TransactionValidityErrorUnknownTransaction:
		value, err := DecodeUnknownTransaction(buffer)
		if err != nil {
			return TransactionValidityError{}, err
		}
		err = NewTransactionValidityError(value)
		if txErr, ok := err.(TransactionValidityError); ok {
			return txErr, nil
		}
		return TransactionValidityError{}, err
	default:
		return TransactionValidityError{}, errInvalidTransactionValidityErrorType
	}
}

func (e TransactionValidityError) Bytes() []byte {
	return sc.EncodedBytes(e)
}

func (e TransactionValidityError) MetadataDefinition(typesInvalidTxId int, typesUnknownTxId int) *MetadataTypeDefinition {
	def := NewMetadataTypeDefinitionVariant(
		sc.Sequence[MetadataDefinitionVariant]{
			NewMetadataDefinitionVariant(
				"Invalid",
				sc.Sequence[MetadataTypeDefinitionField]{
					NewMetadataTypeDefinitionField(typesInvalidTxId),
				},
				TransactionValidityErrorInvalidTransaction,
				""),
			NewMetadataDefinitionVariant(
				"Unknown",
				sc.Sequence[MetadataTypeDefinitionField]{
					NewMetadataTypeDefinitionField(typesUnknownTxId),
				},
				TransactionValidityErrorUnknownTransaction,
				""),
		},
	)
	return &def
}
