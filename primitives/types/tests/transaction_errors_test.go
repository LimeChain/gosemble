package types

import (
	"bytes"
	"testing"

	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/require"
)

// TODO: add more test cases

func Test_EncodeTransactionValidityError(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       types.TransactionValidityError
		expectation []byte
	}{
		{
			label:       "Encode(TransactionValidityError(InvalidTransaction(PaymentError)))",
			input:       types.NewTransactionValidityError(types.NewInvalidTransaction(types.PaymentError)),
			expectation: []byte{0x00, 0x01},
		},
		{
			label:       "Encode(TransactionValidityError(UnknownTransaction(0)))",
			input:       types.NewTransactionValidityError(types.NewUnknownTransaction(types.CannotLookupError)),
			expectation: []byte{0x01, 0x00},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			testExample.input.Encode(buffer)

			require.Equal(t, testExample.expectation, buffer.Bytes())
		})
	}
}

func Test_DecodeTransactionValidityError(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation types.TransactionValidityError
	}{
		{
			label:       "Encode(TransactionValidityError(InvalidTransaction(PaymentError)))",
			input:       []byte{0x00, 0x01},
			expectation: types.NewTransactionValidityError(types.NewInvalidTransaction(types.PaymentError)),
		},
		{
			label:       "Encode(TransactionValidityError(UnknownTransaction(0)))",
			input:       []byte{0x01, 0x00},
			expectation: types.NewTransactionValidityError(types.NewUnknownTransaction(types.CannotLookupError)),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result := types.DecodeTransactionValidityError(buffer)

			require.Equal(t, testExample.expectation, result)
		})
	}
}
