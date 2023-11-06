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

type ArithmeticError = sc.VaryingData

func NewArithmeticErrorUnderflow() ArithmeticError {
	return sc.NewVaryingData(ArithmeticErrorUnderflow)
}

func NewArithmeticErrorOverflow() ArithmeticError {
	return sc.NewVaryingData(ArithmeticErrorOverflow)
}

func NewArithmeticErrorDivisionByZero() ArithmeticError {
	return sc.NewVaryingData(ArithmeticErrorDivisionByZero)
}

func DecodeArithmeticError(buffer *bytes.Buffer) (ArithmeticError, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return ArithmeticError{}, err
	}

	switch b {
	case ArithmeticErrorUnderflow:
		return NewArithmeticErrorUnderflow(), nil
	case ArithmeticErrorOverflow:
		return NewArithmeticErrorOverflow(), nil
	case ArithmeticErrorDivisionByZero:
		return NewArithmeticErrorDivisionByZero(), nil
	default:
		return nil, NewTypeError("ArithmeticError")
	}
}
