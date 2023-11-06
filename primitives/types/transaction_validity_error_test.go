package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	unknownTransactionCannotLookup, _ = NewTransactionValidityError(NewUnknownTransactionCannotLookup())
)

func Test_NewTransactionValidityError_TypeError(t *testing.T) {
	result, err := NewTransactionValidityError(sc.U8(6))

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'TransactionValidityError' type", err.Error())
	assert.Equal(t, TransactionValidityError{}, result)
}

func Test_TransactionValidityError_Encode(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       TransactionValidityError
		expectation []byte
	}{
		{
			label:       "Encode(TransactionValidityError(InvalidTransaction(PaymentError)))",
			input:       invalidTransactionPayment,
			expectation: []byte{0x00, 0x01},
		},
		{
			label:       "Encode(TransactionValidityError(UnknownTransaction(0)))",
			input:       unknownTransactionCannotLookup,
			expectation: []byte{0x01, 0x00},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			testExample.input.Encode(buffer)

			assert.Equal(t, testExample.expectation, buffer.Bytes())
		})
	}
}

func Test_TransactionValidityError_Encode_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}

	assert.PanicsWithValue(t, "not a valid 'TransactionValidityError' type", func() {
		tve := TransactionValidityError(sc.NewVaryingData(sc.U8(6)))

		tve.Encode(buffer)
	})

	assert.Equal(t, &bytes.Buffer{}, buffer)
}

func Test_DecodeTransactionValidityError(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation TransactionValidityError
	}{
		{
			label:       "Encode(TransactionValidityError(InvalidTransaction(PaymentError)))",
			input:       []byte{0x00, 0x01},
			expectation: invalidTransactionPayment,
		},
		{
			label:       "Encode(TransactionValidityError(UnknownTransaction(0)))",
			input:       []byte{0x01, 0x00},
			expectation: unknownTransactionCannotLookup,
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result, err := DecodeTransactionValidityError(buffer)
			assert.NoError(t, err)

			assert.Equal(t, testExample.expectation, result)
		})
	}
}

func Test_DecodeTransactionValidityError_TypeError(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(6)

	res, err := DecodeTransactionValidityError(buffer)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'TransactionValidityError' type", err.Error())
	assert.Equal(t, TransactionValidityError{}, res)

}

func Test_TransactionValidityError_Bytes(t *testing.T) {
	expect, _ := hex.DecodeString("0001")

	assert.Equal(t, expect, invalidTransactionPayment.Bytes())
}
