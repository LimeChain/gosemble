package types

import (
	"bytes"
	sc "github.com/LimeChain/goscale"
)

const (
	InherentError_InherentDataExists = sc.U8(iota)
	InherentError_DecodingFailed
	InherentError_FatalErrorReported
	InherentError_Application
)

type InherentError sc.VaryingData

func NewInherentError(values ...sc.Encodable) InherentError {
	return InherentError{sc.NewVaryingData(values...)}
}

func (ie InherentError) IsFatal() sc.Bool {
	switch ie[0] {
	case InherentError_FatalErrorReported:
		return true
	default:
		return false
	}
}

func (ie InherentError) Encode(buffer *bytes.Buffer) {
	switch ie[0] {
	case InherentError_InherentDataExists:
		ie[0].Encode(buffer)
		ie[1].Encode(buffer)
	case InherentError_DecodingFailed:
		ie[0].Encode(buffer)
		ie[1].Encode(buffer)
	case InherentError_FatalErrorReported:
		ie[0].Encode(buffer)
	case InherentError_Application:
		ie[0].Encode(buffer)
		// TODO:
	default:
		panic("Invalid InherentError Type")
	}
}

func (ie InherentError) Bytes() []byte {
	return sc.EncodedBytes(ie)
}

type CheckInherentsResult struct {
	Okay       sc.Bool
	FatalError sc.Bool
	errors     InherentData
}

func NewCheckInherentsResult() CheckInherentsResult {
	return CheckInherentsResult{
		Okay:       true,
		FatalError: false,
		errors:     *NewInherentData(),
	}
}

func (cir CheckInherentsResult) Encode(buffer *bytes.Buffer) {
	cir.Okay.Encode(buffer)
	cir.FatalError.Encode(buffer)
	cir.errors.Encode(buffer)
}

func (cir CheckInherentsResult) PutError(inherentIdentifier [8]byte, error IsFatalError) IsFatalError {
	if cir.FatalError {
		return InherentError{}
	}

	if error.IsFatal() {
		cir.errors.Clear()
	}

	err := cir.errors.Put(inherentIdentifier, error)
	if err != nil {
		return err
	}

	cir.Okay = false
	cir.FatalError = error.IsFatal()

	return nil
}
