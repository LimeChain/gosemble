package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	ArithmeticErrorUnderflow sc.U8 = iota
	ArithmeticErrorOverflow
	ArithmeticErrorDivisionByZero
)

type ArithmeticError struct {
	sc.VaryingData
}

func NewArithmeticErrorUnderflow() ArithmeticError {
	return ArithmeticError{sc.NewVaryingData(ArithmeticErrorUnderflow)}
}

func NewArithmeticErrorOverflow() ArithmeticError {
	return ArithmeticError{sc.NewVaryingData(ArithmeticErrorOverflow)}
}

func NewArithmeticErrorDivisionByZero() ArithmeticError {
	return ArithmeticError{sc.NewVaryingData(ArithmeticErrorDivisionByZero)}
}

func (err ArithmeticError) Error() string {
	if len(err.VaryingData) == 0 {
		return ""
	}

	switch err.VaryingData[0] {
	case ArithmeticErrorUnderflow:
		return "An underflow would occur"
	case ArithmeticErrorOverflow:
		return "An overflow would occur"
	case ArithmeticErrorDivisionByZero:
		return "Division by zero"
	default:
		return ""
	}
}

func DecodeArithmeticError(buffer *bytes.Buffer) error {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return err
	}

	switch b {
	case ArithmeticErrorUnderflow:
		return NewArithmeticErrorUnderflow()
	case ArithmeticErrorOverflow:
		return NewArithmeticErrorOverflow()
	case ArithmeticErrorDivisionByZero:
		return NewArithmeticErrorDivisionByZero()
	default:
		return newTypeError("ArithmeticError")
	}
}
