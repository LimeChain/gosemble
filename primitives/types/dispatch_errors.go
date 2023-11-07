package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/utils"
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

type DispatchError = sc.VaryingData

func NewDispatchErrorOther(str sc.Str) DispatchError {
	return sc.NewVaryingData(DispatchErrorOther, str)
}

func NewDispatchErrorCannotLookup() DispatchError {
	return sc.NewVaryingData(DispatchErrorCannotLookup)
}

func NewDispatchErrorBadOrigin() DispatchError {
	return sc.NewVaryingData(DispatchErrorBadOrigin)
}

func NewDispatchErrorModule(customModuleError CustomModuleError) DispatchError {
	return sc.NewVaryingData(DispatchErrorModule, customModuleError)
}

func NewDispatchErrorConsumerRemaining() DispatchError {
	return sc.NewVaryingData(DispatchErrorConsumerRemaining)
}

func NewDispatchErrorNoProviders() DispatchError {
	return sc.NewVaryingData(DispatchErrorNoProviders)
}

func NewDispatchErrorTooManyConsumers() DispatchError {
	return sc.NewVaryingData(DispatchErrorTooManyConsumers)
}

func NewDispatchErrorToken(tokenError TokenError) DispatchError {
	// TODO: type safety
	return sc.NewVaryingData(DispatchErrorToken, tokenError)
}

func NewDispatchErrorArithmetic(arithmeticError ArithmeticError) DispatchError {
	// TODO: type safety
	return sc.NewVaryingData(DispatchErrorArithmetic, arithmeticError)
}

func NewDispatchErrorTransactional(transactionalError TransactionalError) DispatchError {
	// TODO: type safety
	return sc.NewVaryingData(DispatchErrorTransactional, transactionalError)
}

func NewDispatchErrorExhausted() DispatchError {
	return sc.NewVaryingData(DispatchErrorExhausted)
}

func NewDispatchErrorCorruption() DispatchError {
	return sc.NewVaryingData(DispatchErrorCorruption)
}

func NewDispatchErrorUnavailable() DispatchError {
	return sc.NewVaryingData(DispatchErrorUnavailable)
}

func DecodeDispatchError(buffer *bytes.Buffer) (DispatchError, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return nil, err
	}

	switch b {
	case DispatchErrorOther:
		value, err := sc.DecodeStr(buffer)
		if err != nil {
			return nil, err
		}
		return NewDispatchErrorOther(value), nil
	case DispatchErrorCannotLookup:
		return NewDispatchErrorCannotLookup(), nil
	case DispatchErrorBadOrigin:
		return NewDispatchErrorBadOrigin(), nil
	case DispatchErrorModule:
		module, err := DecodeCustomModuleError(buffer)
		if err != nil {
			return DispatchError{}, err
		}
		return NewDispatchErrorModule(module), nil
	case DispatchErrorConsumerRemaining:
		return NewDispatchErrorConsumerRemaining(), nil
	case DispatchErrorNoProviders:
		return NewDispatchErrorNoProviders(), nil
	case DispatchErrorTooManyConsumers:
		return NewDispatchErrorTooManyConsumers(), nil
	case DispatchErrorToken:
		tokenError, err := DecodeTokenError(buffer)
		if err != nil {
			return DispatchError{}, err
		}
		return NewDispatchErrorToken(tokenError), nil
	case DispatchErrorArithmetic:
		arithmeticError, err := DecodeArithmeticError(buffer)
		if err != nil {
			return nil, err
		}
		return NewDispatchErrorArithmetic(arithmeticError), nil
	case DispatchErrorTransactional:
		transactionalError, err := DecodeTransactionalError(buffer)
		if err != nil {
			return DispatchError{}, err
		}
		return NewDispatchErrorTransactional(transactionalError), nil
	case DispatchErrorExhausted:
		return NewDispatchErrorExhausted(), nil
	case DispatchErrorCorruption:
		return NewDispatchErrorCorruption(), nil
	case DispatchErrorUnavailable:
		return NewDispatchErrorUnavailable(), nil
	default:
		return DispatchError{}, newTypeError("DispatchError")
	}
}

// CustomModuleError A custom error in a module.
type CustomModuleError struct {
	Index   sc.U8             // Module index matching the metadata module index.
	Error   sc.U32            // Module specific error value.
	Message sc.Option[sc.Str] // Varying data type Option (Definition 190). The optional value is a SCALE encoded byte array containing a valid UTF-8 sequence.
}

func (e CustomModuleError) Encode(buffer *bytes.Buffer) error {
	return utils.EncodeEach(buffer,
		e.Index,
		e.Error,
	) // e.Message is skipped in codec
}

func DecodeCustomModuleError(buffer *bytes.Buffer) (CustomModuleError, error) {
	e := CustomModuleError{}
	idx, err := sc.DecodeU8(buffer)
	if err != nil {
		return CustomModuleError{}, err
	}
	e.Index = idx
	decodedErr, err := sc.DecodeU32(buffer)
	if err != nil {
		return CustomModuleError{}, err
	}
	e.Error = decodedErr
	//e.Message = sc.DecodeOption[sc.Str](buffer) // Skipped in codec
	return e, nil
}

func (e CustomModuleError) Bytes() []byte {
	return sc.EncodedBytes(e)
}

// DispatchErrorWithPostInfo Result of a `Dispatchable` which contains the `DispatchResult` and additional information about
// the `Dispatchable` that is only known post dispatch.
type DispatchErrorWithPostInfo[T sc.Encodable] struct {
	// Additional information about the `Dispatchable` which is only known post dispatch.
	PostInfo T

	// The actual `DispatchResult` indicating whether the dispatch was successful.
	Error DispatchError
}

func (e DispatchErrorWithPostInfo[PostDispatchInfo]) Encode(buffer *bytes.Buffer) error {
	return utils.EncodeEach(buffer,
		e.PostInfo,
		e.Error,
	)
}

func DecodeErrorWithPostInfo(buffer *bytes.Buffer) (DispatchErrorWithPostInfo[PostDispatchInfo], error) {
	e := DispatchErrorWithPostInfo[PostDispatchInfo]{}
	postInfo, err := DecodePostDispatchInfo(buffer)
	if err != nil {
		return DispatchErrorWithPostInfo[PostDispatchInfo]{}, err
	}
	e.PostInfo = postInfo
	dispatchError, err := DecodeDispatchError(buffer)
	if err != nil {
		return DispatchErrorWithPostInfo[PostDispatchInfo]{}, err
	}
	e.Error = dispatchError
	return e, nil
}

func (e DispatchErrorWithPostInfo[PostDispatchInfo]) Bytes() []byte {
	return sc.EncodedBytes(e)
}
