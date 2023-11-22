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

type TransactionalError sc.VaryingData

func NewTransactionalErrorLimitReached() TransactionalError {
	return TransactionalError(sc.NewVaryingData(TransactionalErrorLimitReached))
}

func NewTransactionalErrorNoLayer() TransactionalError {
	return TransactionalError(sc.NewVaryingData(TransactionalErrorNoLayer))
}

func (err TransactionalError) Encode(buffer *bytes.Buffer) error {
	return err[0].Encode(buffer)
}
func (err TransactionalError) Error() string {
	if len(err) == 0 {
		return newTypeError("TransactionalError").Error()
	}

	switch err[0] {
	case TransactionalErrorLimitReached:
		return "Too many transactional layers have been spawned"
	case TransactionalErrorNoLayer:
		return "A transactional layer was expected, but does not exist"
	default:
		return newTypeError("TransactionalError").Error()
	}
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
		return TransactionalError{}, newTypeError("TransactionalError")
	}
}

func (err TransactionalError) Bytes() []byte {
	return sc.EncodedBytes(err)
}
