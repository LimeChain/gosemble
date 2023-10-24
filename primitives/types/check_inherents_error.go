package types

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
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
	// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
	switch ie.VaryingData[0] {
	case InherentErrorInherentDataExists:
		// return fmt.Sprintf("Inherent data already exists for identifier: [%v]", ie.VaryingData[1])
		identifier := sc.SequenceU8ToBytes(ie.VaryingData[1].(sc.Sequence[sc.U8]))
		return "Inherent data already exists for identifier: [" + string(identifier) + "]"
	case InherentErrorDecodingFailed:
		// return fmt.Sprintf("Failed to decode inherent data for identifier: [%v]", ie.VaryingData[1])
		identifier := sc.SequenceU8ToBytes(ie.VaryingData[1].(sc.Sequence[sc.U8]))
		return "Failed to decode inherent data for identifier: [" + string(identifier) + "]"
	case InherentErrorFatalErrorReported:
		return "There was already a fatal error reported and no other errors are allowed"
	case InherentErrorApplication:
		return "Inherent error application"
	default:
		log.Critical("invalid inherent error")
	}

	panic("unreachable")
}
