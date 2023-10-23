package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewUnknownTransactionCannotLookup(t *testing.T) {
	assert.Equal(t, UnknownTransaction{sc.NewVaryingData(UnknownTransactionCannotLookup)}, NewUnknownTransactionCannotLookup())
}

func Test_NewUnknownTransactionNoUnsignedValidator(t *testing.T) {
	assert.Equal(t, UnknownTransaction{sc.NewVaryingData(UnknownTransactionNoUnsignedValidator)}, NewUnknownTransactionNoUnsignedValidator())
}

func Test_NewUnknownTransactionCustom(t *testing.T) {
	unknown := sc.U8(5)

	assert.Equal(t, UnknownTransaction{sc.NewVaryingData(UnknownTransactionCustomUnknownTransaction, unknown)}, NewUnknownTransactionCustomUnknownTransaction(unknown))
}

func Test_DecodeUnknownTransaction_CannotLookup(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(0)

	result := DecodeUnknownTransaction(buffer)

	assert.Equal(t, NewUnknownTransactionCannotLookup(), result)
}

func Test_DecodeUnknownTransaction_NoUnsignedValidator(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(1)

	result := DecodeUnknownTransaction(buffer)

	assert.Equal(t, NewUnknownTransactionNoUnsignedValidator(), result)
}

func Test_DecodeUnknownTransaction_CustomUnknownTransaction(t *testing.T) {
	unknownTxId := sc.U8(5)

	buffer := &bytes.Buffer{}
	buffer.WriteByte(2)
	buffer.WriteByte(byte(unknownTxId))

	result := DecodeUnknownTransaction(buffer)

	assert.Equal(t, NewUnknownTransactionCustomUnknownTransaction(unknownTxId), result)
}

func Test_DecodeUnknownTransaction_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(5)

	assert.PanicsWithValue(t, "invalid UnknownTransaction type", func() {
		DecodeUnknownTransaction(buffer)
	})
}
