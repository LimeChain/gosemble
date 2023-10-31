package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewTransactionValidityResult_Panics(t *testing.T) {
	assert.PanicsWithValue(t, errInvalidTransactionValidityResultType, func() {
		NewTransactionValidityResult(sc.U8(6))
	})
}

func Test_TransactionValidityResult_Encode(t *testing.T) {
	var testExamples = []struct {
		label  string
		input  TransactionValidityResult
		expect []byte
	}{
		{
			label:  "Encode(TransactionValidityResult(ValidTransaction))",
			input:  NewTransactionValidityResult(DefaultValidTransaction()),
			expect: append(TransactionValidityResultValid.Bytes(), DefaultValidTransaction().Bytes()...),
		},
		{
			label:  "Encode(TransactionValidityResult(TransactionValidityError))",
			input:  NewTransactionValidityResult(NewTransactionValidityError(NewInvalidTransactionPayment())),
			expect: append(TransactionValidityResultError.Bytes(), NewTransactionValidityError(NewInvalidTransactionPayment()).Bytes()...),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			testExample.input.Encode(buffer)

			assert.Equal(t, testExample.expect, buffer.Bytes())
		})
	}
}

func Test_TransactionValidityResult_Encode_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}

	assert.PanicsWithValue(t, errInvalidTransactionValidityResultType, func() {
		tve := TransactionValidityResult(sc.NewVaryingData(sc.U8(6)))

		tve.Encode(buffer)
	})

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
			expect: NewTransactionValidityResult(DefaultValidTransaction()),
		},
		{
			label:  "Encode(TransactionValidityResult(TransactionValidityError))",
			input:  []byte{0x01, 0x00, 0x01},
			expect: NewTransactionValidityResult(NewTransactionValidityError(NewInvalidTransactionPayment())),
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

func Test_DecodeTransactionValidityResult_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(6)

	assert.PanicsWithValue(t, errInvalidTransactionValidityResultType, func() {
		DecodeTransactionValidityResult(buffer)
	})
}

func Test_TransactionValidityResult_Bytes(t *testing.T) {
	expect, _ := hex.DecodeString("010001")
	tve := NewTransactionValidityResult(NewTransactionValidityError(NewInvalidTransactionPayment()))

	assert.Equal(t, expect, tve.Bytes())
}

func Test_TransactionValidityResult_IsValidTransaction_True(t *testing.T) {
	target := NewTransactionValidityResult(DefaultValidTransaction())

	assert.Equal(t, true, target.IsValidTransaction())
}

func Test_TransactionValidityResult_IsValidTransaction_False(t *testing.T) {
	target := NewTransactionValidityResult(NewTransactionValidityError(NewInvalidTransactionPayment()))

	assert.Equal(t, false, target.IsValidTransaction())
}

func Test_TransactionValidityResult_AsValidTransaction(t *testing.T) {
	validTx := ValidTransaction{
		Priority:  1,
		Requires:  sc.Sequence[TransactionTag]{},
		Provides:  sc.Sequence[TransactionTag]{},
		Longevity: 4,
		Propagate: true,
	}
	target := NewTransactionValidityResult(validTx)

	assert.Equal(t, validTx, target.AsValidTransaction())
}

func Test_TransactionValidityResult_AsValidTransaction_Panics(t *testing.T) {
	target := NewTransactionValidityResult(NewTransactionValidityError(NewInvalidTransactionPayment()))

	assert.PanicsWithValue(t, "not a ValidTransaction type", func() {
		target.AsValidTransaction()
	})
}
