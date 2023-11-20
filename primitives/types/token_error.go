package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	TokenErrorNoFunds sc.U8 = iota
	TokenErrorWouldDie
	TokenErrorBelowMinimum
	TokenErrorCannotCreate
	TokenErrorUnknownAsset
	TokenErrorFrozen
	TokenErrorUnsupported
)

type TokenError struct {
	sc.VaryingData
}

func NewTokenErrorNoFounds() TokenError {
	return TokenError{sc.NewVaryingData(TokenErrorNoFunds)}
}

func NewTokenErrorWouldDie() TokenError {
	return TokenError{sc.NewVaryingData(TokenErrorWouldDie)}
}

func NewTokenErrorBelowMinimum() TokenError {
	return TokenError{sc.NewVaryingData(TokenErrorBelowMinimum)}
}

func NewTokenErrorCannotCreate() TokenError {
	return TokenError{sc.NewVaryingData(TokenErrorCannotCreate)}
}

func NewTokenErrorUnknownAsset() TokenError {
	return TokenError{sc.NewVaryingData(TokenErrorUnknownAsset)}
}

func NewTokenErrorFrozen() TokenError {
	return TokenError{sc.NewVaryingData(TokenErrorFrozen)}
}

func NewTokenErrorUnsupported() TokenError {
	return TokenError{sc.NewVaryingData(TokenErrorUnsupported)}
}

func (err TokenError) Error() string {
	if len(err.VaryingData) == 0 {
		return ""
	}

	switch err.VaryingData[0] {
	case TokenErrorNoFunds:
		return "Funds are unavailable"
	case TokenErrorWouldDie:
		return "Account that must exist would die"
	case TokenErrorBelowMinimum:
		return "Account cannot exist with the funds that would be given"
	case TokenErrorCannotCreate:
		return "Account cannot be created"
	case TokenErrorUnknownAsset:
		return "The asset in question is unknown"
	case TokenErrorFrozen:
		return "Funds exist but are frozen"
	case TokenErrorUnsupported:
		return "Operation is not supported by the asset"
	default:
		return ""
	}
}

func DecodeTokenError(buffer *bytes.Buffer) error {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return err
	}

	switch b {
	case TokenErrorNoFunds:
		return NewTokenErrorNoFounds()
	case TokenErrorWouldDie:
		return NewTokenErrorWouldDie()
	case TokenErrorBelowMinimum:
		return NewTokenErrorBelowMinimum()
	case TokenErrorCannotCreate:
		return NewTokenErrorCannotCreate()
	case TokenErrorUnknownAsset:
		return NewTokenErrorUnknownAsset()
	case TokenErrorFrozen:
		return NewTokenErrorFrozen()
	case TokenErrorUnsupported:
		return NewTokenErrorUnsupported()
	default:
		return newTypeError("TokenError")
	}
}
