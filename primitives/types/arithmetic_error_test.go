package types

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

func Test_NewArithmeticErrorUnderflow(t *testing.T) {
	result := NewArithmeticErrorUnderflow()

	assert.Equal(t, sc.U8(0), result.VaryingData[0])
}

func Test_NewArithmeticErrorOverflow(t *testing.T) {
	result := NewArithmeticErrorOverflow()

	assert.Equal(t, sc.U8(1), result.VaryingData[0])
}

func Test_NewArithmeticErrorDivisionByZero(t *testing.T) {
	result := NewArithmeticErrorDivisionByZero()

	assert.Equal(t, sc.U8(2), result.VaryingData[0])
}

func Test_ArithmeticError_Decode(t *testing.T) {
	var testExamples = []struct {
		label       string
		input       []byte
		expectation ArithmeticError
	}{
		{label: "ArithmeticErrorUnderflow", input: []byte{0x00}, expectation: NewArithmeticErrorUnderflow()},
		{label: "ArithmeticErrorOverflow", input: []byte{0x01}, expectation: NewArithmeticErrorOverflow()},
		{label: "ArithmeticErrorDivisionByZero", input: []byte{0x02}, expectation: NewArithmeticErrorDivisionByZero()},
	}

	for _, testExample := range testExamples {
		t.Run(testExample.label, func(t *testing.T) {
			result, err := DecodeArithmeticError(bytes.NewBuffer(testExample.input))
			assert.NoError(t, err)

			assert.Equal(t, result, testExample.expectation)
		})
	}
}

func Test_Decode_TypeError(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{0x03})

	res, err := DecodeArithmeticError(buffer)

	assert.Error(t, err)
	assert.Equal(t, "not a valid 'ArithmeticError' type", err.Error())
	assert.Nil(t, res.VaryingData)
}
