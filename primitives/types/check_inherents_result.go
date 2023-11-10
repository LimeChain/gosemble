package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

const (
	InherentErrorInherentDataExists sc.U8 = iota
	InherentErrorDecodingFailed
	InherentErrorFatalErrorReported
	InherentErrorApplication
)

type CheckInherentsResult struct {
	Okay       sc.Bool
	FatalError sc.Bool
	Errors     InherentData
}

func NewCheckInherentsResult() CheckInherentsResult {
	return CheckInherentsResult{
		Okay:       true,
		FatalError: false,
		Errors:     *NewInherentData(),
	}
}

func (cir CheckInherentsResult) Encode(buffer *bytes.Buffer) error {
	err := cir.Okay.Encode(buffer)
	if err != nil {
		return err
	}
	err = cir.FatalError.Encode(buffer)
	if err != nil {
		return err
	}
	return cir.Errors.Encode(buffer)
}

func (cir CheckInherentsResult) Bytes() []byte {
	return sc.EncodedBytes(cir)
}

func (cir *CheckInherentsResult) PutError(inherentIdentifier [8]byte, fatalError FatalError) error {
	if cir.FatalError {
		return NewInherentErrorFatalErrorReported()
	}

	if fatalError.IsFatal() {
		cir.Errors.Clear()
	}

	err := cir.Errors.Put(inherentIdentifier, fatalError)
	if err != nil {
		return err
	}

	cir.Okay = false
	cir.FatalError = fatalError.IsFatal()

	return nil
}

func DecodeCheckInherentsResult(buffer *bytes.Buffer) (CheckInherentsResult, error) {
	okay, err := sc.DecodeBool(buffer)
	if err != nil {
		return CheckInherentsResult{}, err
	}
	fatalError, err := sc.DecodeBool(buffer)
	if err != nil {
		return CheckInherentsResult{}, err
	}
	errors, err := DecodeInherentData(buffer)
	if err != nil {
		return CheckInherentsResult{}, err
	}

	return CheckInherentsResult{
		Okay:       okay,
		FatalError: fatalError,
		Errors:     *errors,
	}, nil
}
