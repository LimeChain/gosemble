package types

import (
	"bytes"
	"fmt"

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

type DispatchError sc.VaryingData

func NewDispatchErrorOther(str sc.Str) DispatchError {
	return DispatchError(sc.NewVaryingData(DispatchErrorOther, str))
}

func NewDispatchErrorCannotLookup() DispatchError {
	return DispatchError(sc.NewVaryingData(DispatchErrorCannotLookup))
}

func NewDispatchErrorBadOrigin() DispatchError {
	return DispatchError(sc.NewVaryingData(DispatchErrorBadOrigin))
}

func NewDispatchErrorModule(customModuleError CustomModuleError) DispatchError {
	return DispatchError(sc.NewVaryingData(DispatchErrorModule, customModuleError))
}

func NewDispatchErrorConsumerRemaining() DispatchError {
	return DispatchError(sc.NewVaryingData(DispatchErrorConsumerRemaining))
}

func NewDispatchErrorNoProviders() DispatchError {
	return DispatchError(sc.NewVaryingData(DispatchErrorNoProviders))
}

func NewDispatchErrorTooManyConsumers() DispatchError {
	return DispatchError(sc.NewVaryingData(DispatchErrorTooManyConsumers))
}

func NewDispatchErrorToken(tokenError TokenError) DispatchError {
	return DispatchError(sc.NewVaryingData(DispatchErrorToken, tokenError))
}

func NewDispatchErrorArithmetic(arithmeticError ArithmeticError) DispatchError {
	return DispatchError(sc.NewVaryingData(DispatchErrorArithmetic, arithmeticError))
}

func NewDispatchErrorTransactional(transactionalError TransactionalError) DispatchError {
	return DispatchError(sc.NewVaryingData(DispatchErrorTransactional, transactionalError))
}

func NewDispatchErrorExhausted() DispatchError {
	return DispatchError(sc.NewVaryingData(DispatchErrorExhausted))
}

func NewDispatchErrorCorruption() DispatchError {
	return DispatchError(sc.NewVaryingData(DispatchErrorCorruption))
}

func NewDispatchErrorUnavailable() DispatchError {
	return DispatchError(sc.NewVaryingData(DispatchErrorUnavailable))
}

func (err DispatchError) Error() string {
	if len(err) == 0 {
		return newTypeError("DispatchError").Error()
	}
	switch err[0] {
	case DispatchErrorOther:
		return "Some unknown error occurred"
	case DispatchErrorCannotLookup:
		return "Cannot lookup"
	case DispatchErrorBadOrigin:
		return "Bad origin"
	case DispatchErrorModule:
		return err[1].(CustomModuleError).Error()
	case DispatchErrorConsumerRemaining:
		return "Consumer remaining"
	case DispatchErrorNoProviders:
		return "No providers"
	case DispatchErrorTooManyConsumers:
		return "Too many consumers"
	case DispatchErrorToken:
		return err[1].(TokenError).Error()
	case DispatchErrorArithmetic:
		return err[1].(ArithmeticError).Error()
	case DispatchErrorExhausted:
		return "Resources exhausted"
	case DispatchErrorCorruption:
		return "State corrupt"
	case DispatchErrorUnavailable:
		return "Resource unavailable"
	default:
		return newTypeError("DispatchError").Error()
	}
}
func (err DispatchError) Encode(buffer *bytes.Buffer) error {
	switch err[0] {
	case DispatchErrorCannotLookup, DispatchErrorBadOrigin, DispatchErrorConsumerRemaining, DispatchErrorNoProviders, DispatchErrorTooManyConsumers, DispatchErrorExhausted, DispatchErrorCorruption, DispatchErrorUnavailable:
		return err[0].Encode(buffer)
	case DispatchErrorOther, DispatchErrorModule, DispatchErrorToken, DispatchErrorArithmetic, DispatchErrorTransactional:
		return sc.EncodeEach(buffer, err[0], err[1])
	default:
		return newTypeError("DispatchError")
	}
}
func DecodeDispatchError(buffer *bytes.Buffer) (DispatchError, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return DispatchError{}, err
	}

	switch b {
	case DispatchErrorOther:
		value, err := sc.DecodeStr(buffer)
		if err != nil {
			return DispatchError{}, err
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
			return DispatchError{}, err
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

func (err DispatchError) Bytes() []byte {
	return sc.EncodedBytes(err)
}

// CustomModuleError A custom error in a module.
type CustomModuleError struct {
	Index   sc.U8             // Module index matching the metadata module index.
	Err     sc.U32            // Module specific error value.
	Message sc.Option[sc.Str] // Varying data type Option (Definition 190). The optional value is a SCALE encoded byte array containing a valid UTF-8 sequence.
}

func (err CustomModuleError) Error() string {
	return fmt.Sprintf("custom module error: index [%d], err [%v], message [%s]", err.Index, err.Err, err.Message.Value)
}

func (err CustomModuleError) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
		err.Index,
		err.Err,
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
	e.Err = decodedErr
	//e.Message = sc.DecodeOption[sc.Str](buffer) // Skipped in codec
	return e, nil
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
	Error DispatchError
}

func (e DispatchErrorWithPostInfo[PostDispatchInfo]) Encode(buffer *bytes.Buffer) error {
	return sc.EncodeEach(buffer,
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
