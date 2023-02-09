package types

import (
	"bytes"
	"testing"

	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/require"
)

func Test_EncodeApplyExtrinsicResult(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       types.ApplyExtrinsicResult
		expectation []byte
	}{
		{
			label:       "Encode ApplyExtrinsicResult(NewDispatchOutcome(None))",
			input:       types.NewApplyExtrinsicResult(types.NewDispatchOutcome(nil)),
			expectation: []byte{0x00, 0x00},
		},
		{
			label:       "Encode ApplyExtrinsicResult(NewDispatchOutcome(NewDispatchError(BadOriginError)))",
			input:       types.NewApplyExtrinsicResult(types.NewDispatchOutcome(types.NewDispatchError(types.BadOriginError{}))),
			expectation: []byte{0x00, 0x01, 0x02},
		},
		{
			label:       "Encode ApplyExtrinsicResult(NewTransactionValidityError(NewInvalidTransaction(CallError)))",
			input:       types.NewApplyExtrinsicResult(types.NewTransactionValidityError(types.NewInvalidTransaction(types.CallError))),
			expectation: []byte{0x01, 0x00, 0x00},
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

func Test_DecodeApplyExtrinsicResult(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation types.ApplyExtrinsicResult
	}{
		{
			label:       "Decode ApplyExtrinsicResult(NewDispatchOutcome(None))",
			expectation: types.NewApplyExtrinsicResult(types.NewDispatchOutcome(nil)),
			input:       []byte{0x00, 0x00},
		},
		{
			label:       "Decode ApplyExtrinsicResult(NewDispatchOutcome(NewDispatchError(BadOriginError)))",
			expectation: types.NewApplyExtrinsicResult(types.NewDispatchOutcome(types.NewDispatchError(types.BadOriginError{}))),
			input:       []byte{0x00, 0x01, 0x02},
		},
		{
			label:       "Decode ApplyExtrinsicResult(NewTransactionValidityError(NewInvalidTransaction(CallError)))",
			expectation: types.NewApplyExtrinsicResult(types.NewTransactionValidityError(types.NewInvalidTransaction(types.CallError))),
			input:       []byte{0x01, 0x00, 0x00},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result := types.DecodeApplyExtrinsicResult(buffer)

			require.Equal(t, testExample.expectation, result)
		})
	}
}
