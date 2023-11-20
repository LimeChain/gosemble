package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	DispatchErrorOther sc.U8 = iota
	DispatchErrorCannotLookup
	DispatchErrorBadOrigin
	DispatchErrorModule
	DispatchErrorConsumerRemaining
	DispatchErrorNoProviders
	DispatchErrorTooManyConsumers
	DispatchErrorToken
	DispatchErrorArithmetic
	DispatchErrorTransactional
	DispatchErrorExhausted
	DispatchErrorCorruption
	DispatchErrorUnavailable
)

type DispatchError struct {
	sc.VaryingData
}

func NewDispatchErrorOther(str sc.Str) DispatchError {
	return DispatchError{sc.NewVaryingData(DispatchErrorOther, str)}
}

func NewDispatchErrorCannotLookup() DispatchError {
	return DispatchError{sc.NewVaryingData(DispatchErrorCannotLookup)}
}

func NewDispatchErrorBadOrigin() DispatchError {
	return DispatchError{sc.NewVaryingData(DispatchErrorBadOrigin)}
}

func NewDispatchErrorModule(customModuleError CustomModuleError) DispatchError {
	return DispatchError{sc.NewVaryingData(DispatchErrorModule, customModuleError)}
}

func NewDispatchErrorConsumerRemaining() DispatchError {
	return DispatchError{sc.NewVaryingData(DispatchErrorConsumerRemaining)}
}

func NewDispatchErrorNoProviders() DispatchError {
	return DispatchError{sc.NewVaryingData(DispatchErrorNoProviders)}
}

func NewDispatchErrorTooManyConsumers() DispatchError {
	return DispatchError{sc.NewVaryingData(DispatchErrorTooManyConsumers)}
}

func NewDispatchErrorToken(tokenError TokenError) DispatchError {
	// TODO: type safety
	return DispatchError{sc.NewVaryingData(DispatchErrorToken, tokenError)}
}

func NewDispatchErrorArithmetic(arithmeticError ArithmeticError) DispatchError {
	// TODO: type safety
	return DispatchError{sc.NewVaryingData(DispatchErrorArithmetic, arithmeticError)}
}

func NewDispatchErrorTransactional(transactionalError TransactionalError) DispatchError {
	// TODO: type safety
	return DispatchError{sc.NewVaryingData(DispatchErrorTransactional, transactionalError)}
}

func NewDispatchErrorExhausted() DispatchError {
	return DispatchError{sc.NewVaryingData(DispatchErrorExhausted)}
}

func NewDispatchErrorCorruption() DispatchError {
	return DispatchError{sc.NewVaryingData(DispatchErrorCorruption)}
}

func NewDispatchErrorUnavailable() DispatchError {
	return DispatchError{sc.NewVaryingData(DispatchErrorUnavailable)}
}

func (err DispatchError) Error() string {
	if len(err.VaryingData) == 0 {
		return ""
	}

	switch dispatchErr := err.VaryingData[0]; dispatchErr {
	case DispatchErrorOther:
		return "Some unknown error occurred"
	case DispatchErrorCannotLookup:
		return "Cannot lookup"
	case DispatchErrorBadOrigin:
		return "Bad origin"
	case DispatchErrorModule:
		return dispatchErr.(CustomModuleError).Error()
	case DispatchErrorConsumerRemaining:
		return "Consumer remaining"
	case DispatchErrorNoProviders:
		return "No providers"
	case DispatchErrorTooManyConsumers:
		return "Too many consumers"
	case DispatchErrorToken:
		return dispatchErr.(TokenError).Error()
	case DispatchErrorArithmetic:
		return dispatchErr.(ArithmeticError).Error()
	case DispatchErrorExhausted:
		return "Resources exhausted"
	case DispatchErrorCorruption:
		return "State corrupt"
	case DispatchErrorUnavailable:
		return "Resource unavailable"
	default:
		return ""
	}
} // e[0].(InvalidTransaction).Error()

func DecodeDispatchError(buffer *bytes.Buffer) error {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return err
	}

	switch b {
	case DispatchErrorOther:
		value, err := sc.DecodeStr(buffer)
		if err != nil {
			return err
		}
		return NewDispatchErrorOther(value)
	case DispatchErrorCannotLookup:
		return NewDispatchErrorCannotLookup()
	case DispatchErrorBadOrigin:
		return NewDispatchErrorBadOrigin()
	case DispatchErrorModule:
		customErr := DecodeCustomModuleError(buffer)
		if _, ok := customErr.(CustomModuleError); !ok {
			return err
		}
		return NewDispatchErrorModule(customErr.(CustomModuleError))
	case DispatchErrorConsumerRemaining:
		return NewDispatchErrorConsumerRemaining()
	case DispatchErrorNoProviders:
		return NewDispatchErrorNoProviders()
	case DispatchErrorTooManyConsumers:
		return NewDispatchErrorTooManyConsumers()
	case DispatchErrorToken:
		tokenError := DecodeTokenError(buffer)
		if _, ok := tokenError.(TokenError); !ok {
			return err
		}
		return NewDispatchErrorToken(tokenError.(TokenError))
	case DispatchErrorArithmetic:
		arithmeticErr := DecodeArithmeticError(buffer)
		if _, ok := arithmeticErr.(ArithmeticError); !ok {
			return err
		}
		return NewDispatchErrorArithmetic(arithmeticErr.(ArithmeticError))
	case DispatchErrorTransactional:
		txErr := DecodeTransactionalError(buffer)
		if _, ok := txErr.(TransactionalError); !ok {
			return err
		}
		return NewDispatchErrorTransactional(txErr.(TransactionalError))
	case DispatchErrorExhausted:
		return NewDispatchErrorExhausted()
	case DispatchErrorCorruption:
		return NewDispatchErrorCorruption()
	case DispatchErrorUnavailable:
		return NewDispatchErrorUnavailable()
	default:
		return newTypeError("DispatchError")
	}
}

// CustomModuleError A custom error in a module.
type CustomModuleError struct {
	Index   sc.U8             // Module index matching the metadata module index.
	Err     sc.U32            // Module specific error value.
	Message sc.Option[sc.Str] // Varying data type Option (Definition 190). The optional value is a SCALE encoded byte array containing a valid UTF-8 sequence.
}

func (err CustomModuleError) Error() string {
	// todo
	return "todo"
}

func (err CustomModuleError) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		err.Index,
		err.Err,
	) // e.Message is skipped in codec
}

func DecodeCustomModuleError(buffer *bytes.Buffer) error {
	cErr := CustomModuleError{}
	idx, err := sc.DecodeU8(buffer)
	if err != nil {
		return err
	}
	cErr.Index = idx
	decodedErr, err := sc.DecodeU32(buffer)
	if err != nil {
		return err
	}
	cErr.Err = decodedErr
	//e.Message = sc.DecodeOption[sc.Str](buffer) // Skipped in codec
	return cErr
}

func (err CustomModuleError) Bytes() []byte {
	return sc.EncodedBytes(err)
}

// DispatchErrorWithPostInfo Result of a `Dispatchable` which contains the `DispatchResult` and additional information about
// the `Dispatchable` that is only known post dispatch.
type DispatchErrorWithPostInfo[T sc.Encodable] struct {
	// Additional information about the `Dispatchable` which is only known post dispatch.
	PostInfo T

	// The actual `DispatchResult` indicating whether the dispatch was successful.
	Err error
}

func (err DispatchErrorWithPostInfo[PostDispatchInfo]) Error() string {
	return err.Err.Error()
}

func (err DispatchErrorWithPostInfo[PostDispatchInfo]) Encode(buffer *bytes.Buffer) error {
	if err := err.PostInfo.Encode(buffer); err != nil {
		return err
	}
	if _, ok := err.Err.(DispatchError); ok {
		return err.Err.(DispatchError).Encode(buffer)
	}

	return nil
}

func DecodeErrorWithPostInfo(buffer *bytes.Buffer) error {
	postInfo, err := DecodePostDispatchInfo(buffer)
	if err != nil {
		return err
	}

	return DispatchErrorWithPostInfo[PostDispatchInfo]{PostInfo: postInfo, Err: DecodeDispatchError(buffer)}
}

func (err DispatchErrorWithPostInfo[PostDispatchInfo]) Bytes() []byte {
	return sc.EncodedBytes(err)
}
