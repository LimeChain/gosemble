package types

import (
	"bytes"
	"errors"
	sc "github.com/LimeChain/goscale"
)

const (
	InherentDataExists = sc.U8(iota)
	DecodingFailed
	FatalErrorReported
	Application
)

const (
	InherentError_ValidAtTimestamp = sc.U8(iota)
	InherentError_TooFarInTheFuture
)

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
		return errors.New(FatalErrorReported.String())
	}

	if error.IsFatalError() {
		cir.errors.Clear()
	}

	err := cir.errors.Put(inherentIdentifier, sc.Str(error.Error()))
	if err != nil {
		return err
	}

	cir.Okay = false
	cir.FatalError = error.IsFatalError()

	return nil
}
