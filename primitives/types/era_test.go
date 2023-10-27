package types

import (
	"bytes"
	"math"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewMortalEra(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []sc.U64
		expectation Era
	}{
		{
			label:       "New(64, 42)",
			input:       []sc.U64{64, 42},
			expectation: Era{EraPeriod: 64, EraPhase: 42},
		},
		{
			label:       "New(32768, 20000)",
			input:       []sc.U64{32768, 20000},
			expectation: Era{EraPeriod: 32768, EraPhase: 20000},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			eraResult := NewMortalEra(testExample.input[0], testExample.input[1])
			assert.Equal(t, testExample.expectation, eraResult)
		})
	}
}

func Test_NewImmortalEra(t *testing.T) {
	assert.Equal(t, Era{IsImmortal: true}, NewImmortalEra())
}

func Test_Era_Encode(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       Era
		expectation []byte
	}{
		{
			label:       "ImmortalEra",
			input:       Era{IsImmortal: true},
			expectation: []byte{0x00},
		},
		{
			label:       "MortalEra(64, 42)",
			input:       Era{IsImmortal: false, EraPeriod: 64, EraPhase: 42},
			expectation: []byte{165, 2},
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
			label:       "0x00",
			input:       []byte{0x00},
			expectation: Era{IsImmortal: true},
		},
		{
			label:       "0xa5, 0x02",
			input:       []byte{0xa5, 0x02},
			expectation: Era{IsImmortal: false, EraPeriod: 64, EraPhase: 42},
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

func Test_Era_Bytes(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       Era
		expectation []byte
	}{
		{
			label:       "ImmortalEra",
			input:       Era{IsImmortal: true},
			expectation: []byte{0x00},
		},
		{
			label:       "MortalEra(64, 42)",
			input:       Era{IsImmortal: false, EraPeriod: 64, EraPhase: 42},
			expectation: []byte{165, 2},
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			assert.Equal(t, testExample.expectation, testExample.input.Bytes())
		})
	}
}

func Test_Era_Birth(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       Era
		expectation sc.U64
	}{
		{
			label:       "ImmortalEra",
			input:       Era{IsImmortal: true},
			expectation: sc.U64(0),
		},
		{
			label:       "MortalEra(30, 20)",
			input:       Era{IsImmortal: false, EraPeriod: 30, EraPhase: 20},
			expectation: sc.U64(20),
		},
		{
			label:       "MortalEra(20, 10)",
			input:       Era{IsImmortal: false, EraPeriod: 20, EraPhase: 10},
			expectation: sc.U64(10),
		},
		{
			label:       "MortalEra(10, 5)",
			input:       Era{IsImmortal: false, EraPeriod: 10, EraPhase: 5},
			expectation: sc.U64(15),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			current := sc.U64(15)

			assert.Equal(t, testExample.expectation, testExample.input.Birth(current))
		})
	}
}

func Test_Era_Death(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       Era
		expectation sc.U64
	}{
		{
			label:       "ImmortalEra",
			input:       Era{IsImmortal: true},
			expectation: sc.U64(math.MaxUint64),
		},
		{
			label:       "MortalEra(30, 20)",
			input:       Era{IsImmortal: false, EraPeriod: 30, EraPhase: 20},
			expectation: sc.U64(50),
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			current := sc.U64(15)

			assert.Equal(t, testExample.expectation, testExample.input.Death(current))
		})
	}
}

func Test_EraTypeDefinition(t *testing.T) {
	// TODO
}
