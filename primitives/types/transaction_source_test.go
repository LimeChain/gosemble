package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewTransactionSourceInBlock(t *testing.T) {
	assert.Equal(t, sc.NewVaryingData(TransactionSourceInBlock), NewTransactionSourceInBlock())
}

func Test_NewTransactionSourceLocal(t *testing.T) {
	assert.Equal(t, sc.NewVaryingData(TransactionSourceLocal), NewTransactionSourceLocal())
}

func Test_NewTransactionSourceExternal(t *testing.T) {
	assert.Equal(t, sc.NewVaryingData(TransactionSourceExternal), NewTransactionSourceExternal())
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

func Test_DecodeTransactionSource_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(3)

	assert.PanicsWithValue(t, "invalid TransactionSource type", func() {
		DecodeTransactionSource(buffer)
	})
}
