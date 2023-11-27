package types

import (
	"bytes"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewEcdsaVerifyErrorBadRS(t *testing.T) {
	assert.Equal(t, EcdsaVerifyError(sc.NewVaryingData(EcdsaVerifyErrorBadRS)), NewEcdsaVerifyErrorBadRS())
}

func Test_NewEcdsaVerifyErrorBadV(t *testing.T) {
	assert.Equal(t, EcdsaVerifyError(sc.NewVaryingData(EcdsaVerifyErrorBadV)), NewEcdsaVerifyErrorBadV())
}

func Test_NewEcdsaVerifyErrorBadSignature(t *testing.T) {
	assert.Equal(t, EcdsaVerifyError(sc.NewVaryingData(EcdsaVerifyErrorBadSignature)), NewEcdsaVerifyErrorBadSignature())
}

func Test_EcdsaVerifyError_Encode(t *testing.T) {
	target := NewEcdsaVerifyErrorBadSignature()
	buffer := &bytes.Buffer{}

	err := target.Encode(buffer)

	assert.Nil(t, err)
	assert.Equal(t, []byte{byte(EcdsaVerifyErrorBadSignature)}, buffer.Bytes())
}

func Test_EcdsaVerifyError_Encode_Empty(t *testing.T) {
	target := EcdsaVerifyError(sc.NewVaryingData())
	buffer := &bytes.Buffer{}

	err := target.Encode(buffer)

	assert.Equal(t, newTypeError("EcdsaVerifyError"), err)
}

func Test_DecodeEcdsaVerifyError_BadRS(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(0)

	result, err := DecodeEcdsaVerifyError(buffer)

	assert.Nil(t, err)
	assert.Equal(t, NewEcdsaVerifyErrorBadRS(), result)
}

func Test_DecodeEcdsaVerifyError_BadV(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(1)

	result, err := DecodeEcdsaVerifyError(buffer)

	assert.Nil(t, err)
	assert.Equal(t, NewEcdsaVerifyErrorBadV(), result)
}

func Test_DecodeEcdsaVerifyError_BadSignature(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(2)

	result, err := DecodeEcdsaVerifyError(buffer)

	assert.Nil(t, err)
	assert.Equal(t, NewEcdsaVerifyErrorBadSignature(), result)
}

func Test_DecodeEcdsaVerifyError_TypeError(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(3)

	result, err := DecodeEcdsaVerifyError(buffer)

	assert.Nil(t, result)
	assert.Equal(t, newTypeError("EcdsaVerifyError"), err)
}

func Test_DecodeEcdsaVerifyError_EmptyBuffer(t *testing.T) {
	buffer := &bytes.Buffer{}

	result, err := DecodeEcdsaVerifyError(buffer)

	assert.Equal(t, io.EOF, err)
	assert.Nil(t, result)
}

func Test_EcdsaVerifyError_Bytes(t *testing.T) {
	target := NewEcdsaVerifyErrorBadSignature()

	result := target.Bytes()

	assert.Equal(t, []byte{byte(EcdsaVerifyErrorBadSignature)}, result)
}

func Test_EcdsaVerifyError_Error_BadRS(t *testing.T) {
	target := NewEcdsaVerifyErrorBadRS()

	assert.Equal(t, "Bad RS", target.Error())
}

func Test_EcdsaVerifyError_Error_BadV(t *testing.T) {
	target := NewEcdsaVerifyErrorBadV()

	assert.Equal(t, "Bad V", target.Error())
}

func Test_EcdsaVerifyError_Error_BadSignature(t *testing.T) {
	target := NewEcdsaVerifyErrorBadSignature()

	assert.Equal(t, "Bad signature", target.Error())
}

func Test_EcdsaVerifyError_Error_Empty(t *testing.T) {
	target := EcdsaVerifyError(sc.NewVaryingData())

	errorString := target.Error()

	assert.Equal(t, "not a valid 'EcdsaVerifyError' type", errorString)
}

func Test_EcdsaVerifyError_Error_InvalidValue(t *testing.T) {
	target := EcdsaVerifyError(sc.NewVaryingData(sc.U8(3)))

	errorString := target.Error()

	assert.Equal(t, "not a valid 'EcdsaVerifyError' type", errorString)
}
