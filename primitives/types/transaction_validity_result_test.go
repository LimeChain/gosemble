package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	invalidTransactionPayment, _ = NewTransactionValidityError(NewInvalidTransactionPayment())

	transactionValidityResultTransactionPayment, _ = NewTransactionValidityResult(invalidTransactionPayment)
	transactionValidityResultDefaultValid, _       = NewTransactionValidityResult(DefaultValidTransaction())
)

func Test_NewTransactionValidityResult_TypeError(t *testing.T) {
	result, err := NewTransactionValidityResult(sc.U8(6))

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'TransactionValidityResult' type", err.Error())
	assert.Equal(t, TransactionValidityResult{}, result)
}

func Test_TransactionValidityResult_Encode(t *testing.T) {
	var testExamples = []struct {
		label  string
		input  TransactionValidityResult
		expect []byte
	}{
		{
			label:  "Encode(TransactionValidityResult(ValidTransaction))",
			input:  transactionValidityResultDefaultValid,
			expect: append(TransactionValidityResultValid.Bytes(), DefaultValidTransaction().Bytes()...),
		},
		{
			label:  "Encode(TransactionValidityResult(TransactionValidityError))",
			input:  transactionValidityResultTransactionPayment,
			expect: append(TransactionValidityResultError.Bytes(), invalidTransactionPayment.Bytes()...),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			err := testExample.input.Encode(buffer)

			assert.NoError(t, err)
			assert.Equal(t, testExample.expect, buffer.Bytes())
		})
	}
}

func Test_TransactionValidityResult_Encode_TypeError(t *testing.T) {
	buffer := &bytes.Buffer{}
	tve := TransactionValidityResult(sc.NewVaryingData(sc.U8(6)))

	err := tve.Encode(buffer)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'TransactionValidityResult' type", err.Error())
	assert.Equal(t, &bytes.Buffer{}, buffer)
}

func Test_DecodeTransactionValidityResult(t *testing.T) {
	var testExamples = []struct {
		label  string
		input  []byte
		expect TransactionValidityResult
	}{
		{
			label:  "Encode(TransactionValidityResult(ValidTransaction))",
			input:  append(TransactionValidityResultValid.Bytes(), DefaultValidTransaction().Bytes()...),
			expect: transactionValidityResultDefaultValid,
		},
		{
			label:  "Encode(TransactionValidityResult(TransactionValidityError))",
			input:  []byte{0x01, 0x00, 0x01},
			expect: transactionValidityResultTransactionPayment,
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result, err := DecodeTransactionValidityResult(buffer)
			assert.NoError(t, err)

			assert.Equal(t, testExample.expect, result)
		})
	}
}

func Test_DecodeTransactionValidityResult_TypeError(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(6)

	result, err := DecodeTransactionValidityResult(buffer)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'TransactionValidityResult' type", err.Error())
	assert.Equal(t, TransactionValidityResult{}, result)
}

func Test_TransactionValidityResult_Bytes(t *testing.T) {
	expect, _ := hex.DecodeString("010001")

	assert.Equal(t, expect, transactionValidityResultTransactionPayment.Bytes())
}

func Test_TransactionValidityResult_IsValidTransaction_True(t *testing.T) {
	assert.Equal(t, true, transactionValidityResultDefaultValid.IsValidTransaction())
}

func Test_TransactionValidityResult_IsValidTransaction_False(t *testing.T) {
	assert.Equal(t, false, transactionValidityResultTransactionPayment.IsValidTransaction())
}

func Test_TransactionValidityResult_AsValidTransaction(t *testing.T) {
	validTx := ValidTransaction{
		Priority:  1,
		Requires:  sc.Sequence[TransactionTag]{},
		Provides:  sc.Sequence[TransactionTag]{},
		Longevity: 4,
		Propagate: true,
	}
	target, _ := NewTransactionValidityResult(validTx)

	validTransaction, err := target.AsValidTransaction()
	assert.NoError(t, err)
	assert.Equal(t, validTx, validTransaction)
}

func Test_TransactionValidityResult_AsValidTransaction_TypeError(t *testing.T) {
	result, err := transactionValidityResultTransactionPayment.AsValidTransaction()

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'ValidTransaction' type", err.Error())
	assert.Equal(t, ValidTransaction{}, result)
}
