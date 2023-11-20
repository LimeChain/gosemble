package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewTokenErrorNoFunds(t *testing.T) {
	assert.Equal(t, TokenError{sc.NewVaryingData(TokenErrorNoFunds)}, NewTokenErrorNoFounds())
}

func Test_NewTokenErrorWouldDie(t *testing.T) {
	assert.Equal(t, TokenError{sc.NewVaryingData(TokenErrorWouldDie)}, NewTokenErrorWouldDie())
}

func Test_NewTokenErrorBelowMinimum(t *testing.T) {
	assert.Equal(t, TokenError{sc.NewVaryingData(TokenErrorBelowMinimum)}, NewTokenErrorBelowMinimum())
}

func Test_NewTokenErrorCannotCreate(t *testing.T) {
	assert.Equal(t, TokenError{sc.NewVaryingData(TokenErrorCannotCreate)}, NewTokenErrorCannotCreate())
}

func Test_NewTokenErrorUnknownAsset(t *testing.T) {
	assert.Equal(t, TokenError{sc.NewVaryingData(TokenErrorUnknownAsset)}, NewTokenErrorUnknownAsset())
}

func Test_NewTokenErrorFrozen(t *testing.T) {
	assert.Equal(t, TokenError{sc.NewVaryingData(TokenErrorFrozen)}, NewTokenErrorFrozen())
}

func Test_NewTokenErrorUnsupported(t *testing.T) {
	assert.Equal(t, TokenError{sc.NewVaryingData(TokenErrorUnsupported)}, NewTokenErrorUnsupported())
}

func Test_DecodeTokenError_NoFunds(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(0)
	assert.Equal(t, NewTokenErrorNoFounds(), DecodeTokenError(buffer))
}

func Test_DecodeTokenError_WouldDie(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(1)
	assert.Equal(t, NewTokenErrorWouldDie(), DecodeTokenError(buffer))
}

func Test_DecodeTokenError_BelowMinimum(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(2)
	assert.Equal(t, NewTokenErrorBelowMinimum(), DecodeTokenError(buffer))
}

func Test_DecodeTokenError_CannotCreate(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(3)
	assert.Equal(t, NewTokenErrorCannotCreate(), DecodeTokenError(buffer))
}

func Test_DecodeTokenError_UnknownAsset(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(4)
	assert.Equal(t, NewTokenErrorUnknownAsset(), DecodeTokenError(buffer))
}

func Test_DecodeTokenError_Frozen(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(5)
	assert.Equal(t, NewTokenErrorFrozen(), DecodeTokenError(buffer))
}

func Test_DecodeTokenError_Unsupported(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(6)
	assert.Equal(t, NewTokenErrorUnsupported(), DecodeTokenError(buffer))
}

func Test_DecodeTokenError_TypeError(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(7)
	assert.Equal(t, "not a valid 'TokenError' type", DecodeTokenError(buffer).Error())
}
