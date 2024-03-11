package system

import (
	"bytes"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_KeyValue_Encode(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       KeyValue
		expectation []byte
	}{
		{
			label:       "KeyValue()",
			input:       KeyValue{},
			expectation: []byte{0x0, 0x0},
		},
		{
			label: "KeyValue(abc:123)",
			input: KeyValue{
				Key:   sc.BytesToSequenceU8([]byte("abc")),
				Value: sc.BytesToSequenceU8([]byte("123")),
			},
			expectation: []byte{0xc, 0x61, 0x62, 0x63, 0xc, 0x31, 0x32, 0x33},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			err := testExample.input.Encode(buffer)

			assert.NoError(t, err)
			assert.Equal(t, testExample.expectation, buffer.Bytes())
			assert.Equal(t, testExample.expectation, testExample.input.Bytes())
		})
	}
}

func Test_DecodeKeyValue(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation KeyValue
	}{
		{
			label: "(00)",
			input: []byte{0x0, 0x0},
			expectation: KeyValue{
				Key:   sc.Sequence[sc.U8]{},
				Value: sc.Sequence[sc.U8]{},
			},
		},
		{
			label: "(0c6162630c313233)",
			input: []byte{0xc, 0x61, 0x62, 0x63, 0xc, 0x31, 0x32, 0x33},
			expectation: KeyValue{
				Key:   sc.BytesToSequenceU8([]byte("abc")),
				Value: sc.BytesToSequenceU8([]byte("123")),
			},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result, err := DecodeKeyValue(buffer)

			assert.NoError(t, err)
			assert.Equal(t, testExample.expectation, result)
		})
	}
}

func Test_DecodeKeyValue_Fails(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation KeyValue
	}{
		{
			label:       "()",
			input:       []byte{},
			expectation: KeyValue{},
		},
		{
			label:       "(0c6162630c)",
			input:       []byte{0xc, 0x61, 0x62, 0x63, 0xc},
			expectation: KeyValue{},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result, err := DecodeKeyValue(buffer)

			assert.Equal(t, io.EOF, err)
			assert.Equal(t, testExample.expectation, result)
		})
	}
}

func Test_CodeUpgradeAuthorization_Encode(t *testing.T) {
	buffer := &bytes.Buffer{}

	codeUpgradeAuthorization := CodeUpgradeAuthorization{codeHash, sc.Bool(true)}
	err := codeUpgradeAuthorization.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, append(hashBytes, 0x1), buffer.Bytes())
	assert.Equal(t, append(hashBytes, 0x1), codeUpgradeAuthorization.Bytes())
}

func Test_DecodeCodeUpgradeAuthorization(t *testing.T) {
	buffer := &bytes.Buffer{}
	buffer.Write(append(hashBytes, 0x1))

	result, err := DecodeCodeUpgradeAuthorization(buffer)

	expectation := CodeUpgradeAuthorization{codeHash, sc.Bool(true)}

	assert.NoError(t, err)
	assert.Equal(t, expectation, result)
}
