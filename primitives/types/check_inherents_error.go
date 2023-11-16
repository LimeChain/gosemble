package types

import (
	"fmt"

	sc "github.com/LimeChain/goscale"
)

const (
	InherentErrorInherentDataExists sc.U8 = iota
	InherentErrorDecodingFailed
	InherentErrorFatalErrorReported
	InherentErrorApplication
)

type InherentError struct {
	sc.VaryingData
}

func NewInherentErrorInherentDataExists(inherentIdentifier sc.FixedSequence[sc.U8]) InherentError {
	return InherentError{sc.NewVaryingData(InherentErrorInherentDataExists, inherentIdentifier)}
}

func NewInherentErrorDecodingFailed(error sc.Str, inherentIdentifier sc.FixedSequence[sc.U8]) InherentError {
	return InherentError{sc.NewVaryingData(InherentErrorDecodingFailed, error, inherentIdentifier)}
}

func NewInherentErrorFatalErrorReported() InherentError {
	return InherentError{sc.NewVaryingData(InherentErrorFatalErrorReported)}
}

func NewInherentErrorApplication(error sc.Str) InherentError {
	return InherentError{sc.NewVaryingData(InherentErrorApplication, error)}
}

func (ie InherentError) IsFatal() sc.Bool {
	switch ie.VaryingData[0] {
	case InherentErrorFatalErrorReported:
		return true
	default:
		return false
	}
}

func (ie InherentError) Error() string {
	switch ie.VaryingData[0] {
	case InherentErrorInherentDataExists:
		return fmt.Sprintf("Inherent data already exists for identifier: %v", ie.VaryingData[1].(sc.FixedSequence[sc.U8]))
	case InherentErrorDecodingFailed:
		return fmt.Sprintf("Failed to decode inherent data for identifier: %v", ie.VaryingData[2].(sc.FixedSequence[sc.U8]))
	case InherentErrorFatalErrorReported:
		return "There was already a fatal error reported and no other errors are allowed"
	case InherentErrorApplication:
		return fmt.Sprintf("Inherent error application: [%s]", ie.VaryingData[1].(sc.Str))
	default:
		return newTypeError("InherentError").Error()
	}
}
