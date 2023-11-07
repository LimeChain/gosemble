package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_DispatchOutcome_New_TypeError(t *testing.T) {
	result, err := NewDispatchOutcome(sc.U8(5))

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'DispatchOutcome' type", err.Error())
	assert.Equal(t, DispatchOutcome{}, result)
}

func Test_DispatchOutcome_Encode(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       DispatchOutcome
		expectation []byte
	}{
		{label: "Encode DispatchOutcome(None)", input: dispatchOutcome, expectation: []byte{0x00}},
		{label: "Encode  DispatchOutcome(DispatchErrorBadOrigin)", input: dispatchOutcomeBadOriginErr, expectation: []byte{0x01, 0x02}},
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

func Test_DispatchOutcome_Encode_TypeError(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := DispatchOutcome{sc.U32(5)}.Encode(buffer)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'DispatchOutcome' type", err.Error())
}

func Test_DispatchOutcome_Bytes(t *testing.T) {
	var testExamples = []struct {
		label  string
		input  DispatchOutcome
		expect []byte
	}{
		{label: "Encode DispatchOutcome(None)", input: dispatchOutcome, expect: []byte{0x00}},
		{label: "Encode  DispatchOutcome(DispatchErrorBadOrigin)", input: dispatchOutcomeBadOriginErr, expect: []byte{0x01, 0x02}},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			result := testExample.input.Bytes()

			assert.Equal(t, testExample.expect, result)
		})
	}
}

func Test_DispatchOutcome_Decode(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation DispatchOutcome
	}{
		{label: "0x00", input: []byte{0x00}, expectation: dispatchOutcome},
		{label: "0x01, 0x02", input: []byte{0x01, 0x02}, expectation: dispatchOutcomeBadOriginErr},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result, err := DecodeDispatchOutcome(buffer)
			assert.NoError(t, err)

			assert.Equal(t, testExample.expectation, result)
		})
	}
}

func Test_DispatchOutcome_Decode_TypeError(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{0x3})

	result, err := DecodeDispatchOutcome(buffer)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'DispatchOutcome' type", err.Error())
	assert.Equal(t, DispatchOutcome{}, result)
}
