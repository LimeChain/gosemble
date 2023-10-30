package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewInvalidTransactionCall(t *testing.T) {
	expect := InvalidTransaction{sc.NewVaryingData(InvalidTransactionCall)}

	assert.Equal(t, expect, NewInvalidTransactionCall())
}

func Test_NewInvalidTransactionPayment(t *testing.T) {
	expect := InvalidTransaction{sc.NewVaryingData(InvalidTransactionPayment)}

	assert.Equal(t, expect, NewInvalidTransactionPayment())
}

func Test_NewInvalidTransactionFuture(t *testing.T) {
	expect := InvalidTransaction{sc.NewVaryingData(InvalidTransactionFuture)}

	assert.Equal(t, expect, NewInvalidTransactionFuture())
}

func Test_NewInvalidTransactionStale(t *testing.T) {
	expect := InvalidTransaction{sc.NewVaryingData(InvalidTransactionStale)}

	assert.Equal(t, expect, NewInvalidTransactionStale())
}

func Test_NewInvalidTransactionBadProof(t *testing.T) {
	expect := InvalidTransaction{sc.NewVaryingData(InvalidTransactionBadProof)}

	assert.Equal(t, expect, NewInvalidTransactionBadProof())
}

func Test_NewInvalidTransactionAncientBirthBlock(t *testing.T) {
	expect := InvalidTransaction{sc.NewVaryingData(InvalidTransactionAncientBirthBlock)}

	assert.Equal(t, expect, NewInvalidTransactionAncientBirthBlock())
}

func Test_NewInvalidTransactionExhaustsResources(t *testing.T) {
	expect := InvalidTransaction{sc.NewVaryingData(InvalidTransactionExhaustsResources)}

	assert.Equal(t, expect, NewInvalidTransactionExhaustsResources())
}

func Test_NewInvalidTransactionCustom(t *testing.T) {
	customError := sc.U8(55)
	expect := InvalidTransaction{sc.NewVaryingData(InvalidTransactionCustom, customError)}

	assert.Equal(t, expect, NewInvalidTransactionCustom(customError))
}

func Test_NewInvalidTransactionBadMandatory(t *testing.T) {
	expect := InvalidTransaction{sc.NewVaryingData(InvalidTransactionBadMandatory)}

	assert.Equal(t, expect, NewInvalidTransactionBadMandatory())
}

func Test_NewInvalidTransactionMandatoryValidation(t *testing.T) {
	expect := InvalidTransaction{sc.NewVaryingData(InvalidTransactionMandatoryValidation)}

	assert.Equal(t, expect, NewInvalidTransactionMandatoryValidation())
}

func Test_NewInvalidTransactionBadSigner(t *testing.T) {
	expect := InvalidTransaction{sc.NewVaryingData(InvalidTransactionBadSigner)}

	assert.Equal(t, expect, NewInvalidTransactionBadSigner())
}

func Test_DecodeInvalidTransaction_Call(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(0)

	result, err := DecodeInvalidTransaction(buffer)
	assert.Nil(t, err)

	assert.Equal(t, NewInvalidTransactionCall(), result)
}

func Test_DecodeInvalidTransaction_Payment(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(1)

	result, err := DecodeInvalidTransaction(buffer)
	assert.Nil(t, err)

	assert.Equal(t, NewInvalidTransactionPayment(), result)
}

func Test_DecodeInvalidTransaction_Future(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(2)

	result, err := DecodeInvalidTransaction(buffer)
	assert.Nil(t, err)

	assert.Equal(t, NewInvalidTransactionFuture(), result)
}

func Test_DecodeInvalidTransaction_Stale(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(3)

	result, err := DecodeInvalidTransaction(buffer)
	assert.Nil(t, err)

	assert.Equal(t, NewInvalidTransactionStale(), result)
}

func Test_DecodeInvalidTransaction_BadProof(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(4)

	result, err := DecodeInvalidTransaction(buffer)
	assert.Nil(t, err)

	assert.Equal(t, NewInvalidTransactionBadProof(), result)
}

func Test_DecodeInvalidTransaction_AncientBirthBlock(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(5)

	result, err := DecodeInvalidTransaction(buffer)
	assert.Nil(t, err)

	assert.Equal(t, NewInvalidTransactionAncientBirthBlock(), result)
}

func Test_DecodeInvalidTransaction_ExhaustsResources(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(6)

	result, err := DecodeInvalidTransaction(buffer)
	assert.Nil(t, err)

	assert.Equal(t, NewInvalidTransactionExhaustsResources(), result)
}

func Test_DecodeInvalidTransaction_Custom(t *testing.T) {
	customError := sc.U8(50)
	buffer := &bytes.Buffer{}
	buffer.WriteByte(7)
	buffer.WriteByte(byte(customError))

	result, err := DecodeInvalidTransaction(buffer)
	assert.Nil(t, err)

	assert.Equal(t, NewInvalidTransactionCustom(customError), result)
}

func Test_DecodeInvalidTransaction_BadMandatory(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(8)

	result, err := DecodeInvalidTransaction(buffer)
	assert.Nil(t, err)

	assert.Equal(t, NewInvalidTransactionBadMandatory(), result)
}

func Test_DecodeInvalidTransaction_MandatoryValidation(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(9)

	result, err := DecodeInvalidTransaction(buffer)
	assert.Nil(t, err)

	assert.Equal(t, NewInvalidTransactionMandatoryValidation(), result)
}

func Test_DecodeInvalidTransaction_BadSigner(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(10)

	result, err := DecodeInvalidTransaction(buffer)
	assert.Nil(t, err)

	assert.Equal(t, NewInvalidTransactionBadSigner(), result)
}

func Test_DecodeInvalidTransaction_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(50)

	assert.PanicsWithValue(t, "invalid InvalidTransaction type", func() {
		DecodeInvalidTransaction(buffer)
	})
}
