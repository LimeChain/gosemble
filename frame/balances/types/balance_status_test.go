package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_DecodeBalanceStatus(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation sc.U8
	}{
		{
			label:       "BalanceStatusFree",
			input:       []byte{0x00},
			expectation: BalanceStatusFree,
		},
		{
			label:       "BalanceStatusReserved",
			input:       []byte{0x01},
			expectation: BalanceStatusReserved,
		},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			result, err := DecodeBalanceStatus(bytes.NewBuffer(testExample.input))
			assert.NoError(t, err)

			assert.Equal(t, testExample.expectation, result)
		})
	}
}

func Test_DecodeBalanceStatus_Panics(t *testing.T) {
	assert.PanicsWithValue(t, "invalid balance status type", func() {
		DecodeBalanceStatus(bytes.NewBuffer([]byte{0x02}))
	})
}
