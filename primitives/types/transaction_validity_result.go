package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	errInvalidTransactionValidityResultType = "invalid TransactionValidityResult type"
)

const (
	TransactionValidityResultValid sc.U8 = iota
	TransactionValidityResultError
)

// TransactionValidityResult Information on a transaction's validity and, if valid, on how it relates to other transactions.
type TransactionValidityResult sc.VaryingData

func NewTransactionValidityResult(value sc.Encodable) (TransactionValidityResult, error) {
	switch value.(type) {
	case ValidTransaction, TransactionValidityError:
		return TransactionValidityResult(sc.NewVaryingData(value)), nil
	default:
		return TransactionValidityResult{}, newTypeError("TransactionValidityResult")
	}
}

func (r TransactionValidityResult) Encode(buffer *bytes.Buffer) error {
	switch r[0].(type) {
	case ValidTransaction:
		err := TransactionValidityResultValid.Encode(buffer)
		if err != nil {
			return err
		}
	case TransactionValidityError:
		err := TransactionValidityResultError.Encode(buffer)
		if err != nil {
			return err
		}
	default:
		log.Critical(errInvalidTransactionValidityResultType)
	}

	return r[0].Encode(buffer)
}

func DecodeTransactionValidityResult(buffer *bytes.Buffer) (TransactionValidityResult, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return TransactionValidityResult{}, err
	}

	switch b {
	case TransactionValidityResultValid:
		val, err := DecodeValidTransaction(buffer)
		if err != nil {
			return TransactionValidityResult{}, err
		}
		return NewTransactionValidityResult(val)
	case TransactionValidityResultError:
		val, err := DecodeTransactionValidityError(buffer)
		if err != nil {
			return TransactionValidityResult{}, err
		}
		return NewTransactionValidityResult(val)
	default:
		return TransactionValidityResult{}, newTypeError("TransactionValidityResult")
	}
}

func (r TransactionValidityResult) Bytes() []byte {
	return sc.EncodedBytes(r)
}

func (r TransactionValidityResult) IsValidTransaction() bool {
	switch r[0].(type) {
	case ValidTransaction:
		return true
	default:
		return false
	}
}

func (r TransactionValidityResult) AsValidTransaction() (ValidTransaction, error) {
	if r.IsValidTransaction() {
		return r[0].(ValidTransaction), nil
	} else {
		return ValidTransaction{}, newTypeError("ValidTransaction")
	}
}
