package types

import (
	"fmt"

	sc "github.com/LimeChain/goscale"
)

type InherentError struct {
	sc.VaryingData
}

func NewInherentErrorInherentDataExists(inherentIdentifier sc.Sequence[sc.U8]) InherentError {
	return InherentError{sc.NewVaryingData(InherentErrorInherentDataExists, inherentIdentifier)}
}

func NewInherentErrorDecodingFailed(inherentIdentifier sc.Sequence[sc.U8]) InherentError {
	return InherentError{sc.NewVaryingData(InherentErrorDecodingFailed, inherentIdentifier)}
}

func NewInherentErrorFatalErrorReported() InherentError {
	return InherentError{sc.NewVaryingData(InherentErrorFatalErrorReported)}
}

func NewInherentErrorApplication() InherentError {
	// TODO: encode additional value
	return InherentError{sc.NewVaryingData(InherentErrorApplication)}
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
		return fmt.Sprintf("Inherent data already exists for identifier: %v", ie.VaryingData[1].(sc.Sequence[sc.U8]))
	case InherentErrorDecodingFailed:
		return fmt.Sprintf("Failed to decode inherent data for identifier: %v", ie.VaryingData[1].(sc.Sequence[sc.U8]))
	case InherentErrorFatalErrorReported:
		return "There was already a fatal error reported and no other errors are allowed"
	case InherentErrorApplication:
		return "Inherent error application"
	default:
		return newTypeError("InherentError").Error()
	}
}
