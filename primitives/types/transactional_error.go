package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	// Too many transactional layers have been spawned.
	TransactionalErrorLimitReached sc.U8 = iota
	// A transactional layer was expected, but does not exist.
	TransactionalErrorNoLayer
)

type TransactionalError = sc.VaryingData

func NewTransactionalErrorLimitReached() TransactionalError {
	return sc.NewVaryingData(TransactionalErrorLimitReached)
}

func NewTransactionalErrorNoLayer() TransactionalError {
	return sc.NewVaryingData(TransactionalErrorNoLayer)
}

func DecodeTransactionalError(buffer *bytes.Buffer) (TransactionalError, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return TransactionalError{}, err
	}

	switch b {
	case TransactionalErrorLimitReached:
		return NewTransactionalErrorLimitReached(), nil
	case TransactionalErrorNoLayer:
		return NewTransactionalErrorNoLayer(), nil
	default:
		return nil, NewTypeError("TransactionalError")
	}
}
