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

type TransactionalError struct {
	sc.VaryingData
}

func NewTransactionalErrorLimitReached() TransactionalError {
	return TransactionalError{sc.NewVaryingData(TransactionalErrorLimitReached)}
}

func NewTransactionalErrorNoLayer() TransactionalError {
	return TransactionalError{sc.NewVaryingData(TransactionalErrorNoLayer)}
}

func (err TransactionalError) Error() string {
	if len(err.VaryingData) == 0 {
		return ""
	}

	switch err.VaryingData[0] {
	case TransactionalErrorLimitReached:
		return "Too many transactional layers have been spawned"
	case TransactionalErrorNoLayer:
		return "A transactional layer was expected, but does not exist"
	default:
		return ""
	}
}

func DecodeTransactionalError(buffer *bytes.Buffer) error {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return err
	}

	switch b {
	case TransactionalErrorLimitReached:
		return NewTransactionalErrorLimitReached()
	case TransactionalErrorNoLayer:
		return NewTransactionalErrorNoLayer()
	default:
		return newTypeError("TransactionalError")
	}
}
