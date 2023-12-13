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

type ArithmeticError sc.VaryingData

func NewArithmeticErrorUnderflow() ArithmeticError {
	return ArithmeticError(sc.NewVaryingData(ArithmeticErrorUnderflow))
}

func NewArithmeticErrorOverflow() ArithmeticError {
	return ArithmeticError(sc.NewVaryingData(ArithmeticErrorOverflow))
}

func NewArithmeticErrorDivisionByZero() ArithmeticError {
	return ArithmeticError(sc.NewVaryingData(ArithmeticErrorDivisionByZero))
}

func (err ArithmeticError) Encode(buffer *bytes.Buffer) error {
	return err[0].Encode(buffer)
}

func (err ArithmeticError) Error() string {
	if len(err) == 0 {
		return newTypeError("ArithmeticError").Error()
	}

	switch err[0] {
	case ArithmeticErrorUnderflow:
		return "An underflow would occur"
	case ArithmeticErrorOverflow:
		return "An overflow would occur"
	case ArithmeticErrorDivisionByZero:
		return "Division by zero"
	default:
		return newTypeError("ArithmeticError").Error()
	}
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
		return ArithmeticError{}, newTypeError("ArithmeticError")
	}
}

func (err ArithmeticError) Bytes() []byte {
	return sc.EncodedBytes(err)
}
