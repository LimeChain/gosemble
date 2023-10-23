package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewTransactionValidityError_Panics(t *testing.T) {
	assert.PanicsWithValue(t, errInvalidTransactionValidityErrorType, func() {
		NewTransactionValidityError(sc.U8(6))
	})
}

func Test_TransactionValidityError_Encode(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       TransactionValidityError
		expectation []byte
	}{
		{
			label:       "Encode(TransactionValidityError(InvalidTransaction(PaymentError)))",
			input:       NewTransactionValidityError(NewInvalidTransactionPayment()),
			expectation: []byte{0x00, 0x01},
		},
		{
			label:       "Encode(TransactionValidityError(UnknownTransaction(0)))",
			input:       NewTransactionValidityError(NewUnknownTransactionCannotLookup()),
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

	assert.PanicsWithValue(t, errInvalidTransactionValidityErrorType, func() {
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
			expectation: NewTransactionValidityError(NewInvalidTransactionPayment()),
		},
		{
			label:       "Encode(TransactionValidityError(UnknownTransaction(0)))",
			input:       []byte{0x01, 0x00},
			expectation: NewTransactionValidityError(NewUnknownTransactionCannotLookup()),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result := DecodeTransactionValidityError(buffer)

			assert.Equal(t, testExample.expectation, result)
		})
	}
}

func Test_DecodeTransactionValidityError_Panics(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.WriteByte(6)

	assert.PanicsWithValue(t, errInvalidTransactionValidityErrorType, func() {
		DecodeTransactionValidityError(buffer)
	})
}

func Test_TransactionValidityError_Bytes(t *testing.T) {
	expect, _ := hex.DecodeString("0001")
	tve := NewTransactionValidityError(NewInvalidTransactionPayment())

	assert.Equal(t, expect, tve.Bytes())
}
