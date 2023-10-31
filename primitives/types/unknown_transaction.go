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
		log.Critical("invalid UnknownTransaction type")
	}

	panic("unreachable")
}
