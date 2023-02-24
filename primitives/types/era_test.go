package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_EncodeEra(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       Era
		expectation []byte
	}{
		{
			label:       "Encode Era(ImmortalEra)",
			input:       NewEra(ImmortalEra{}),
			expectation: []byte{0x00},
		},
		{
			label:       "Encode Era(MortalEra)",
			input:       NewEra(MortalEra{EraPeriod: 1, EraPhase: 2}),
			expectation: []byte{0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
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

func Test_DecodeEra(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation Era
	}{
		{
			label:       "Decode Era(ImmortalEra)",
			input:       []byte{0x00},
			expectation: NewEra(ImmortalEra{}),
		},
		{
			label:       "Encode Era(MortalEra)",
			input:       []byte{0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			expectation: NewEra(MortalEra{EraPeriod: 1, EraPhase: 2}),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			buffer := &bytes.Buffer{}
			buffer.Write(testExample.input)

			result := DecodeEra(buffer)

			assert.Equal(t, testExample.expectation, result)
		})
	}
}
