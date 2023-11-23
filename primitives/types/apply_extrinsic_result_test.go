package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	invalidTransactionCall = NewTransactionValidityError(NewInvalidTransactionCall())

	dispatchOutcome, _             = NewDispatchOutcome(nil)
	dispatchOutcomeBadOriginErr, _ = NewDispatchOutcome(NewDispatchErrorBadOrigin())

	applyExtrinsicResultOutcome, _      = NewApplyExtrinsicResult(dispatchOutcome)
	applyExtrinsicResultBadOriginErr, _ = NewApplyExtrinsicResult(dispatchOutcomeBadOriginErr)
	applyExtrinsicResultInvalidCall, _  = NewApplyExtrinsicResult(invalidTransactionCall.(TransactionValidityError))
)

func Test_EncodeApplyExtrinsicResult(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       ApplyExtrinsicResult
		expectation []byte
	}{
		{
			label:       "Encode ApplyExtrinsicResult(NewDispatchOutcome(None))",
			input:       applyExtrinsicResultOutcome,
			expectation: []byte{0x00, 0x00},
		},
		{
			label:       "Encode ApplyExtrinsicResult(NewDispatchOutcome(NewDispatchErrorBadOrigin))",
			input:       applyExtrinsicResultBadOriginErr,
			expectation: []byte{0x00, 0x01, 0x02},
		},
		{
			label:       "Encode ApplyExtrinsicResult(NewTransactionValidityError(NewInvalidTransactionCall))",
			input:       applyExtrinsicResultInvalidCall,
			expectation: []byte{0x01, 0x00, 0x00},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			err := testExample.input.Encode(buffer)

			assert.NoError(t, err)
			assert.Equal(t, testExample.expectation, buffer.Bytes())
		})
	}
}

func Test_DecodeApplyExtrinsicResult(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation ApplyExtrinsicResult
	}{
		{
			label:       "Decode ApplyExtrinsicResult(NewDispatchOutcome(None))",
			expectation: applyExtrinsicResultOutcome,
			input:       []byte{0x00, 0x00},
		},
		{
			label:       "Decode ApplyExtrinsicResult(NewDispatchOutcome(NewDispatchErrorBadOrigin))",
			expectation: applyExtrinsicResultBadOriginErr,
			input:       []byte{0x00, 0x01, 0x02},
		},
		{
			label:       "Decode ApplyExtrinsicResult(NewTransactionValidityError(NewInvalidTransactionCall)",
			expectation: applyExtrinsicResultInvalidCall,
			input:       []byte{0x01, 0x00, 0x00},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result, err := DecodeApplyExtrinsicResult(buffer)
			assert.NoError(t, err)

			assert.Equal(t, testExample.expectation, result)
		})
	}
}
