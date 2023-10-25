package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewTransactionalErrorLimitReached(t *testing.T) {
	assert.Equal(t, sc.NewVaryingData(TransactionalErrorLimitReached), NewTransactionalErrorLimitReached())
}

func Test_NewTransactionalErrorNoLayer(t *testing.T) {
	assert.Equal(t, sc.NewVaryingData(TransactionalErrorNoLayer), NewTransactionalErrorNoLayer())
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

func Test_DecodeTransactionalError_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(5)

	assert.PanicsWithValue(t, "invalid TransactionalError type", func() {
		DecodeTransactionalError(buffer)
	})
}
