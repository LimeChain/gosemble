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

func NewTransactionValidityResult(value sc.Encodable) TransactionValidityResult {
	switch value.(type) {
	case ValidTransaction, TransactionValidityError:
		return TransactionValidityResult(sc.NewVaryingData(value))
	default:
		log.Critical(errInvalidTransactionValidityResultType)
	}

	panic("unreachable")
}

func (r TransactionValidityResult) Encode(buffer *bytes.Buffer) {
	switch r[0].(type) {
	case ValidTransaction:
		TransactionValidityResultValid.Encode(buffer)
	case TransactionValidityError:
		TransactionValidityResultError.Encode(buffer)
	default:
		log.Critical(errInvalidTransactionValidityResultType)
	}

	r[0].Encode(buffer)
}

func DecodeTransactionValidityResult(buffer *bytes.Buffer) TransactionValidityResult {
	b := sc.DecodeU8(buffer)

	switch b {
	case TransactionValidityResultValid:
		return NewTransactionValidityResult(DecodeValidTransaction(buffer))
	case TransactionValidityResultError:
		return NewTransactionValidityResult(DecodeTransactionValidityError(buffer))
	default:
		log.Critical(errInvalidTransactionValidityResultType)
	}

	panic("unreachable")
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

func (r TransactionValidityResult) AsValidTransaction() ValidTransaction {
	if r.IsValidTransaction() {
		return r[0].(ValidTransaction)
	} else {
		log.Critical("not a ValidTransaction type")
	}

	panic("unreachable")
}
