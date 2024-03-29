package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewTransactionalErrorLimitReached(t *testing.T) {
	err := NewTransactionalErrorLimitReached()
	assert.Equal(t, TransactionalError(sc.NewVaryingData(TransactionalErrorLimitReached)), err)
	assert.NotEmpty(t, err.Error())
}

func Test_NewTransactionalErrorNoLayer(t *testing.T) {
	err := NewTransactionalErrorNoLayer()
	assert.Equal(t, TransactionalError(sc.NewVaryingData(TransactionalErrorNoLayer)), err)
	assert.NotEmpty(t, err.Error())
}

func Test_DecodeTransactionalError_LimitReached(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(0)

	result, err := DecodeTransactionalError(buffer)
	assert.NoError(t, err)

	assert.Equal(t, NewTransactionalErrorLimitReached(), result)
}

func Test_DecodeTransactionalError_NoLayer(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(1)

	result, err := DecodeTransactionalError(buffer)
	assert.NoError(t, err)

	assert.Equal(t, NewTransactionalErrorNoLayer(), result)
}

func Test_DecodeTransactionalError_TypeError(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(5)

	_, err := DecodeTransactionalError(buffer)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'TransactionalError' type", err.Error())
}
