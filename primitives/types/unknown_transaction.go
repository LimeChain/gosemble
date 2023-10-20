package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
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

func DecodeUnknownTransaction(buffer *bytes.Buffer) UnknownTransaction {
	b := sc.DecodeU8(buffer)

	switch b {
	case UnknownTransactionCannotLookup:
		return NewUnknownTransactionCannotLookup()
	case UnknownTransactionNoUnsignedValidator:
		return NewUnknownTransactionNoUnsignedValidator()
	case UnknownTransactionCustomUnknownTransaction:
		v := sc.DecodeU8(buffer)
		return NewUnknownTransactionCustomUnknownTransaction(v)
	default:
		log.Critical("invalid UnknownTransaction type")
	}

	panic("unreachable")
}
