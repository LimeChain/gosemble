package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_DispatchOutcome_New_Panics(t *testing.T) {
	assert.PanicsWithValue(t, errDispatchOutcomeInvalid, func() {
		NewDispatchOutcome(sc.U8(5))
	})
}

func Test_DispatchOutcome_Encode(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       DispatchOutcome
		expectation []byte
	}{
		{label: "Encode DispatchOutcome(None)", input: NewDispatchOutcome(nil), expectation: []byte{0x00}},
		{label: "Encode  DispatchOutcome(DispatchErrorBadOrigin)", input: NewDispatchOutcome(NewDispatchErrorBadOrigin()), expectation: []byte{0x01, 0x02}},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			testExample.input.Encode(buffer)

			assert.Equal(t, testExample.expectation, buffer.Bytes())
		})
	}
}

func Test_DispatchOutcome_Encode_Panics(t *testing.T) {
	assert.PanicsWithValue(t, errDispatchOutcomeInvalid, func() {
		buffer := &bytes.Buffer{}
		dispatchOutcome := DispatchOutcome{
			sc.U32(5),
		}

		dispatchOutcome.Encode(buffer)
	})
}

func Test_DispatchOutcome_Bytes(t *testing.T) {
	var testExamples = []struct {
		label  string
		input  DispatchOutcome
		expect []byte
	}{
		{label: "Encode DispatchOutcome(None)", input: NewDispatchOutcome(nil), expect: []byte{0x00}},
		{label: "Encode  DispatchOutcome(DispatchErrorBadOrigin)", input: NewDispatchOutcome(NewDispatchErrorBadOrigin()), expect: []byte{0x01, 0x02}},
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
		{label: "0x00", input: []byte{0x00}, expectation: NewDispatchOutcome(nil)},
		{label: "0x01, 0x02", input: []byte{0x01, 0x02}, expectation: NewDispatchOutcome(NewDispatchErrorBadOrigin())},
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

func Test_DispatchOutcome_Decode_Panics(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{0x3})

	assert.PanicsWithValue(t, errDispatchOutcomeInvalid, func() {
		DecodeDispatchOutcome(buffer)
	})
}
