package types

import (
	"bytes"
	"testing"

	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/require"
)

// dispatchOutcome.Encode(buffer)

func Test_EncodeDispatchOutcome(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       types.DispatchOutcome
		expectation []byte
	}{
		{label: "Encode DispatchOutcome(None)", input: types.NewDispatchOutcome(nil), expectation: []byte{0x00}},
		{label: "Encode  DispatchOutcome(DispatchError(BadOriginError))", input: types.NewDispatchOutcome(types.NewDispatchError(types.BadOriginError{})), expectation: []byte{0x01, 0x02}},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			testExample.input.Encode(buffer)

			require.Equal(t, testExample.expectation, buffer.Bytes())
		})
	}
}

func Test_DecodeDispatchOutcome(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation types.DispatchOutcome
	}{
		{label: "0x00", input: []byte{0x00}, expectation: types.NewDispatchOutcome(nil)},
		{label: "0x01, 0x02", input: []byte{0x01, 0x02}, expectation: types.NewDispatchOutcome(types.NewDispatchError(types.BadOriginError{}))},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result := types.DecodeDispatchOutcome(buffer)

			require.Equal(t, testExample.expectation, result)
		})
	}
}
