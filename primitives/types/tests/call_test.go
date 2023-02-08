package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/require"
)

func Test_NewCall(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       types.Call
		expectation types.Call
	}{
		{
			label: "Encode(Call(System.remark(0xab, 0xcd)))",
			input: types.NewCall("System", "remark", sc.Sequence[sc.U8]{0xab, 0xcd}),
			expectation: types.Call{
				CallIndex: types.CallIndex{
					ModuleIndex:   0,
					FunctionIndex: 0,
				},
				Args: sc.Sequence[sc.U8]{0xab, 0xcd},
			},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			require.Equal(t, testExample.input.CallIndex.ModuleIndex, testExample.expectation.CallIndex.ModuleIndex)
			require.Equal(t, testExample.input.CallIndex.FunctionIndex, testExample.expectation.CallIndex.FunctionIndex)
			require.Equal(t, testExample.input.Args, testExample.expectation.Args)
		})
	}
}

func Test_EncodeCall(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       types.Call
		expectation []byte
	}{
		{
			label:       "Encode(Call(System.remark(0xab, 0xcd)))",
			input:       types.NewCall("System", "remark", sc.Sequence[sc.U8]{0xab, 0xcd}),
			expectation: []byte{0x0, 0x0, 0x8, 0xab, 0xcd},
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

func Test_DecodeCall(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation types.Call
	}{
		{
			label:       "Decode(0x0, 0x0, 0x8, 0xab, 0xcd)",
			input:       []byte{0x0, 0x0, 0x8, 0xab, 0xcd},
			expectation: types.NewCall("System", "remark", sc.Sequence[sc.U8]{0xab, 0xcd}),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result := types.DecodeCall(buffer)

			require.Equal(t, testExample.expectation, result)
		})
	}
}
