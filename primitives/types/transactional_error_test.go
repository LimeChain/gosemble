package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewTransactionalErrorLimitReached(t *testing.T) {
	assert.Equal(t, TransactionalError{sc.NewVaryingData(TransactionalErrorLimitReached)}, NewTransactionalErrorLimitReached())
}

func Test_NewTransactionalErrorNoLayer(t *testing.T) {
	assert.Equal(t, TransactionalError{sc.NewVaryingData(TransactionalErrorNoLayer)}, NewTransactionalErrorNoLayer())
}

func Test_DecodeTransactionalError_LimitReached(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(0)
	assert.Equal(t, NewTransactionalErrorLimitReached(), DecodeTransactionalError(buffer))
}

func Test_DecodeTransactionalError_NoLayer(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(1)
	assert.Equal(t, NewTransactionalErrorNoLayer(), DecodeTransactionalError(buffer))
}

func Test_DecodeTransactionalError_TypeError(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(5)
	assert.Equal(t, "not a valid 'TransactionalError' type", DecodeTransactionalError(buffer).Error())
}
