package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
)

const (
	InherentErrorInherentDataExists sc.U8 = iota
	InherentErrorDecodingFailed
	InherentErrorFatalErrorReported
	InherentErrorApplication
)

const (
	errDecodeInherentData = "failed to decode InherentData"
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

func (cir CheckInherentsResult) Encode(buffer *bytes.Buffer) {
	cir.Okay.Encode(buffer)
	cir.FatalError.Encode(buffer)
	cir.Errors.Encode(buffer)
}

func (cir CheckInherentsResult) Bytes() []byte {
	return sc.EncodedBytes(cir)
}

func (cir *CheckInherentsResult) PutError(inherentIdentifier [8]byte, er FatalError) error {
	if cir.FatalError {
		return NewInherentErrorFatalErrorReported()
	}

	if er.IsFatal() {
		cir.Errors.Clear()
	}

	err := cir.Errors.Put(inherentIdentifier, er)
	if err != nil {
		return err
	}

	cir.Okay = false
	cir.FatalError = er.IsFatal()

	return nil
}

func DecodeCheckInherentsResult(buffer *bytes.Buffer) CheckInherentsResult {
	okay := sc.DecodeBool(buffer)
	fatalError := sc.DecodeBool(buffer)
	errors, err := DecodeInherentData(buffer)
	if err != nil {
		log.Critical(errDecodeInherentData)
	}

	return CheckInherentsResult{
		Okay:       okay,
		FatalError: fatalError,
		Errors:     *errors,
	}
}
