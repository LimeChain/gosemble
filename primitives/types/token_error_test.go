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

	result, err := DecodeTokenError(buffer)
	assert.NoError(t, err)

	assert.Equal(t, NewTokenErrorNoFounds(), result)
}

func Test_DecodeTokenError_WouldDie(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(1)

	result, err := DecodeTokenError(buffer)
	assert.NoError(t, err)

	assert.Equal(t, NewTokenErrorWouldDie(), result)
}

func Test_DecodeTokenError_BelowMinimum(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(2)

	result, err := DecodeTokenError(buffer)
	assert.NoError(t, err)

	assert.Equal(t, NewTokenErrorBelowMinimum(), result)
}

func Test_DecodeTokenError_CannotCreate(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(3)

	result, err := DecodeTokenError(buffer)
	assert.NoError(t, err)

	assert.Equal(t, NewTokenErrorCannotCreate(), result)
}

func Test_DecodeTokenError_UnknownAsset(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(4)

	result, err := DecodeTokenError(buffer)
	assert.NoError(t, err)

	assert.Equal(t, NewTokenErrorUnknownAsset(), result)
}

func Test_DecodeTokenError_Frozen(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(5)

	result, err := DecodeTokenError(buffer)
	assert.NoError(t, err)

	assert.Equal(t, NewTokenErrorFrozen(), result)
}

func Test_DecodeTokenError_Unsupported(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(6)

	result, err := DecodeTokenError(buffer)
	assert.NoError(t, err)

	assert.Equal(t, NewTokenErrorUnsupported(), result)
}

func Test_DecodeTokenError_TypeError(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(7)

	res, err := DecodeTokenError(buffer)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'TokenError' type", err.Error())
	assert.Nil(t, res.VaryingData)
}
