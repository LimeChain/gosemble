package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	EcdsaVerifyErrorBadRS sc.U8 = iota
	EcdsaVerifyErrorBadV
	EcdsaVerifyErrorBadSignature
)

func NewEcdsaVerifyErrorBadRS() EcdsaVerifyError {
	return EcdsaVerifyError(sc.NewVaryingData(EcdsaVerifyErrorBadRS))
}

func NewEcdsaVerifyErrorBadV() EcdsaVerifyError {
	return EcdsaVerifyError(sc.NewVaryingData(EcdsaVerifyErrorBadV))
}

func NewEcdsaVerifyErrorBadSignature() EcdsaVerifyError {
	return EcdsaVerifyError(sc.NewVaryingData(EcdsaVerifyErrorBadSignature))
}

type EcdsaVerifyError sc.VaryingData

func (err EcdsaVerifyError) Encode(buffer *bytes.Buffer) error {
	if len(err) == 0 {
		return newTypeError("EcdsaVerifyError")
	}
	return err[0].Encode(buffer)
}

func DecodeEcdsaVerifyError(buffer *bytes.Buffer) (EcdsaVerifyError, error) {
	b, err := sc.DecodeU8(buffer)
	if err != nil {
		return nil, err
	}

	switch b {
	case EcdsaVerifyErrorBadRS:
		return NewEcdsaVerifyErrorBadRS(), nil
	case EcdsaVerifyErrorBadV:
		return NewEcdsaVerifyErrorBadV(), nil
	case EcdsaVerifyErrorBadSignature:
		return NewEcdsaVerifyErrorBadSignature(), nil
	default:
		return nil, newTypeError("EcdsaVerifyError")
	}
}

func (err EcdsaVerifyError) Bytes() []byte {
	return sc.EncodedBytes(err)
}

func (err EcdsaVerifyError) Error() string {
	if len(err) == 0 {
		return newTypeError("EcdsaVerifyError").Error()
	}

	switch err[0] {
	case EcdsaVerifyErrorBadRS:
		return "Bad RS"
	case EcdsaVerifyErrorBadV:
		return "Bad V"
	case EcdsaVerifyErrorBadSignature:
		return "Bad signature"
	default:
		return newTypeError("EcdsaVerifyError").Error()
	}
}
