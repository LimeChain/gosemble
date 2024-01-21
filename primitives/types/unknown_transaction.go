package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
)

const (
	// Could not lookup some information that is required to validate the transaction. Reject
	UnknownTransactionCannotLookup sc.U8 = iota

	// No validator found for the given unsigned transaction. Reject
	UnknownTransactionNoUnsignedValidator

	// Any other custom unknown validity that is not covered by this type. Reject
	UnknownTransactionCustomUnknownTransaction // + sc.U8
)

type UnknownTransaction struct {
	sc.VaryingData
}

func NewUnknownTransactionCannotLookup() UnknownTransaction {
	return UnknownTransaction{sc.NewVaryingData(UnknownTransactionCannotLookup)}
}

func NewUnknownTransactionNoUnsignedValidator() UnknownTransaction {
	return UnknownTransaction{sc.NewVaryingData(UnknownTransactionNoUnsignedValidator)}
}

func NewUnknownTransactionCustomUnknownTransaction(unknown sc.U8) UnknownTransaction {
	return UnknownTransaction{sc.NewVaryingData(UnknownTransactionCustomUnknownTransaction, unknown)}
}

func (err UnknownTransaction) Error() string {
	if len(err.VaryingData) == 0 {
		return newTypeError("UnknownError").Error()
	}

	switch err.VaryingData[0] {
	case UnknownTransactionCannotLookup:
		return "Could not lookup information required to validate the transaction"
	case UnknownTransactionNoUnsignedValidator:
		return "Could not find an unsigned validator for the unsigned transaction"
	case UnknownTransactionCustomUnknownTransaction:
		return "UnknownTransaction custom error"
	default:
		return newTypeError("TransactionalError").Error()
	}
}

func DecodeUnknownTransaction(buffer *bytes.Buffer) (UnknownTransaction, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return UnknownTransaction{}, err
	}

	switch b {
	case UnknownTransactionCannotLookup:
		return NewUnknownTransactionCannotLookup(), nil
	case UnknownTransactionNoUnsignedValidator:
		return NewUnknownTransactionNoUnsignedValidator(), nil
	case UnknownTransactionCustomUnknownTransaction:
		v, err := sc.DecodeU8(buffer)
		if err != nil {
			return UnknownTransaction{}, err
		}
		return NewUnknownTransactionCustomUnknownTransaction(v), nil
	default:
		return UnknownTransaction{}, newTypeError("UnknownTransaction")
	}
}

func (err UnknownTransaction) MetadataDefinition() *MetadataTypeDefinition {
	def := NewMetadataTypeDefinitionVariant(
		sc.Sequence[MetadataDefinitionVariant]{
			NewMetadataDefinitionVariant(
				"CannotLookup",
				sc.Sequence[MetadataTypeDefinitionField]{},
				UnknownTransactionCannotLookup,
				""),
			NewMetadataDefinitionVariant(
				"NoUnsignedValidator",
				sc.Sequence[MetadataTypeDefinitionField]{},
				UnknownTransactionNoUnsignedValidator,
				""),
			NewMetadataDefinitionVariant(
				"Custom",
				sc.Sequence[MetadataTypeDefinitionField]{
					NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU8),
				},
				UnknownTransactionCustomUnknownTransaction,
				""),
		},
	)

	return &def
}
