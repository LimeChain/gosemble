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

type TokenError = sc.VaryingData

func NewTokenErrorNoFounds() TokenError {
	return sc.NewVaryingData(TokenErrorNoFunds)
}

func NewTokenErrorWouldDie() TokenError {
	return sc.NewVaryingData(TokenErrorWouldDie)
}

func NewTokenErrorBelowMinimum() TokenError {
	return sc.NewVaryingData(TokenErrorBelowMinimum)
}

func NewTokenErrorCannotCreate() TokenError {
	return sc.NewVaryingData(TokenErrorCannotCreate)
}

func NewTokenErrorUnknownAsset() TokenError {
	return sc.NewVaryingData(TokenErrorUnknownAsset)
}

func NewTokenErrorFrozen() TokenError {
	return sc.NewVaryingData(TokenErrorFrozen)
}

func NewTokenErrorUnsupported() TokenError {
	return sc.NewVaryingData(TokenErrorUnsupported)
}

func DecodeTokenError(buffer *bytes.Buffer) (TokenError, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return nil, err
	}

	switch b {
	case TokenErrorNoFunds:
		return NewTokenErrorNoFounds(), nil
	case TokenErrorWouldDie:
		return NewTokenErrorWouldDie(), nil
	case TokenErrorBelowMinimum:
		return NewTokenErrorBelowMinimum(), nil
	case TokenErrorCannotCreate:
		return NewTokenErrorCannotCreate(), nil
	case TokenErrorUnknownAsset:
		return NewTokenErrorUnknownAsset(), nil
	case TokenErrorFrozen:
		return NewTokenErrorFrozen(), nil
	case TokenErrorUnsupported:
		return NewTokenErrorUnsupported(), nil
	default:
		return nil, newTypeError("TokenError")
	}
}
