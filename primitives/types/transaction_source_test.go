package types

import (
	"bytes"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewTransactionSourceInBlock(t *testing.T) {
	assert.Equal(t, TransactionSource(sc.NewVaryingData(TransactionSourceInBlock)), NewTransactionSourceInBlock())
}

func Test_NewTransactionSourceLocal(t *testing.T) {
	assert.Equal(t, TransactionSource(sc.NewVaryingData(TransactionSourceLocal)), NewTransactionSourceLocal())
}

func Test_NewTransactionSourceExternal(t *testing.T) {
	assert.Equal(t, TransactionSource(sc.NewVaryingData(TransactionSourceExternal)), NewTransactionSourceExternal())
}

func Test_TransactionSource_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := NewTransactionSourceInBlock().Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, []byte{0}, buffer.Bytes())
}

func Test_TransactionSource_Encode_Empty(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := TransactionSource{}.Encode(buffer)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'TransactionSource' type", err.Error())
}

func Test_TransactionSource_Bytes(t *testing.T) {
	result := NewTransactionSourceExternal().Bytes()

	assert.Equal(t, []byte{2}, result)
}

func Test_DecodeTransactionSource_InBlock(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(0)

	txSource, err := DecodeTransactionSource(buffer)
	assert.NoError(t, err)
	assert.Equal(t, NewTransactionSourceInBlock(), txSource)
}

func Test_DecodeTransactionSource_Local(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(1)

	txSource, err := DecodeTransactionSource(buffer)
	assert.NoError(t, err)
	assert.Equal(t, NewTransactionSourceLocal(), txSource)
}

func Test_DecodeTransactionSource_External(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(2)

	txSource, err := DecodeTransactionSource(buffer)
	assert.NoError(t, err)
	assert.Equal(t, NewTransactionSourceExternal(), txSource)
}

func Test_DecodeTransactionSource_TypeError(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(3)

	result, err := DecodeTransactionSource(buffer)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'TransactionSource' type", err.Error())
	assert.Equal(t, TransactionSource{}, result)
}

func Test_DecodeTransactionSource_Empty(t *testing.T) {
	buffer := &bytes.Buffer{}

	result, err := DecodeTransactionSource(buffer)

	assert.Equal(t, io.EOF, err)
	assert.Equal(t, TransactionSource{}, result)
}
