package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

// Some unknown error occurred.
type UnknownError = sc.Str

// Failed to lookup some data.
type DataLookupError struct {
	sc.Empty
}

// A bad origin.
type BadOriginError struct {
	sc.Empty
}

// // At least one consumer is remaining so the account cannot be destroyed.
// ConsumerRemainingError

// // There are no providers so the account cannot be created.
// NoProvidersError

// // There are too many consumers so the account cannot be created.
// TooManyConsumersError

// // An error to do with tokens.
// Token(TokenError)

// // An arithmetic error.
// Arithmetic(ArithmeticError)

// // The number of transactional layers has been reached, or we are not in a transactional
// // layer.
// Transactional(TransactionalError)

// // Resources exhausted, e.g. attempt to read/write data which is too large to manipulate.
// ExhaustedError

// // The state is corrupt; this is generally not going to fix itself.
// CorruptionError

// // Some resource (e.g. a preimage) is unavailable right now. This might fix itself later.
// UnavailableError

type DispatchError sc.VaryingData

func NewDispatchError(value sc.Encodable) DispatchError {
	switch value.(type) {
	case UnknownError, DataLookupError, BadOriginError, CustomModuleError:
		return DispatchError(sc.NewVaryingData(value))
	default:
		panic("invalid DispatchError type")
	}
}

func (e DispatchError) Encode(buffer *bytes.Buffer) {
	switch e[0].(type) {
	case UnknownError:
		sc.U8(0).Encode(buffer)
		e[0].Encode(buffer)
	case DataLookupError:
		sc.U8(1).Encode(buffer)
	case BadOriginError:
		sc.U8(2).Encode(buffer)
	case CustomModuleError:
		sc.U8(3).Encode(buffer)
		e[0].Encode(buffer)
	default:
		panic("invalid DispatchError type")
	}
}

func DecodeDispatchError(buffer *bytes.Buffer) DispatchError {
	b := sc.DecodeU8(buffer)

	switch b {
	case 0:
		value := sc.DecodeStr(buffer)
		return NewDispatchError(value)
	case 1:
		return NewDispatchError(DataLookupError{})
	case 2:
		return NewDispatchError(BadOriginError{})
	case 3:
		value := DecodeCustomModuleError(buffer)
		return NewDispatchError(value)
	default:
		panic("invalid DispatchError type")
	}
}

func (e DispatchError) Bytes() []byte {
	return sc.EncodedBytes(e)
}

// A custom error in a module.
type CustomModuleError struct {
	Index   sc.U8             // Module index matching the metadata module index.
	Error   sc.U8             // Module specific error value.
	Message sc.Option[sc.Str] // Varying data type Option (Definition 190). The optional value is a SCALE encoded byte array containing a valid UTF-8 sequence.
}

func (e CustomModuleError) Encode(buffer *bytes.Buffer) {
	e.Index.Encode(buffer)
	e.Error.Encode(buffer)
	e.Message.Encode(buffer)
}

func DecodeCustomModuleError(buffer *bytes.Buffer) CustomModuleError {
	e := CustomModuleError{}
	e.Index = sc.DecodeU8(buffer)
	e.Error = sc.DecodeU8(buffer)
	e.Message = sc.DecodeOption[sc.Str](buffer)
	return e
}

func (e CustomModuleError) Bytes() []byte {
	return sc.EncodedBytes(e)
}

// Result of a `Dispatchable` which contains the `DispatchResult` and additional information about
// the `Dispatchable` that is only known post dispatch.
type DispatchErrorWithPostInfo = TransactionValidityError

// type DispatchErrorWithPostInfo struct {
// 	// Additional information about the `Dispatchable` which is only known post dispatch.
// 	PostInfo DispatchInfo

// 	// The actual `DispatchResult` indicating whether the dispatch was successful.
// 	Error DispatchError
// }
