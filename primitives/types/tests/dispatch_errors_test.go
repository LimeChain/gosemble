package types

import (
	"bytes"
	"testing"

	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/require"
)

func Test_EncodeDispatchError(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       types.DispatchError
		expectation []byte
	}{
		{label: "Encode(DispatchError('unknwon error'))", input: types.NewDispatchError(types.UnknownError("unknown error")), expectation: []byte{0x00, 0x34, 0x75, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x20, 0x65, 0x72, 0x72, 0x6f, 0x72}},
		{label: "Encode(DispatchError(DataLookupError))", input: types.NewDispatchError(types.DataLookupError{}), expectation: []byte{0x01}},
		{label: "Encode(DispatchError(BadOriginError))", input: types.NewDispatchError(types.BadOriginError{}), expectation: []byte{0x02}},
		{label: "Encode(DispatchError(CustomModuleError))", input: types.NewDispatchError(types.CustomModuleError{}), expectation: []byte{0x03, 0x00, 0x00, 0x00}},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			testExample.input.Encode(buffer)

			require.Equal(t, testExample.expectation, buffer.Bytes())
		})
	}
}

func Test_DecodeDispatchError(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation types.DispatchError
	}{

		{label: "DecodeDispatchError(0x00, 0x34, 0x75, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x20, 0x65, 0x72, 0x72, 0x6f, 0x72)", input: []byte{0x00, 0x34, 0x75, 0x6e, 0x6b, 0x6e, 0x6f, 0x77, 0x6e, 0x20, 0x65, 0x72, 0x72, 0x6f, 0x72}, expectation: types.NewDispatchError(types.UnknownError("unknown error"))},
		{label: "DecodeDispatchError(0x01)", input: []byte{0x01}, expectation: types.NewDispatchError(types.DataLookupError{})},
		{label: "DecodeDispatchError(0x02)", input: []byte{0x02}, expectation: types.NewDispatchError(types.BadOriginError{})},
		{label: "DecodeDispatchError(0x03, 0x00, 0x00, 0x00)", input: []byte{0x03, 0x00, 0x00, 0x00}, expectation: types.NewDispatchError(types.CustomModuleError{})},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result := types.DecodeDispatchError(buffer)

			require.Equal(t, testExample.expectation, result)
		})
	}
}
