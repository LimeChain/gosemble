package types

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	identifier = sc.BytesToSequenceU8([]byte{0x05, 0x06})
)

func Test_NewInherentErrorInherentDataExists(t *testing.T) {
	assert.Equal(t,
		InherentError{sc.NewVaryingData(sc.U8(0), identifier)},
		NewInherentErrorInherentDataExists(identifier),
	)
}

func Test_NewInherentErrorDecodingFailed(t *testing.T) {
	assert.Equal(t,
		InherentError{sc.NewVaryingData(sc.U8(1), identifier)},
		NewInherentErrorDecodingFailed(identifier),
	)
}

func Test_NewInherentErrorFatalErrorReported(t *testing.T) {
	assert.Equal(t,
		InherentError{sc.NewVaryingData(sc.U8(2))},
		NewInherentErrorFatalErrorReported(),
	)
}

func Test_NewInherentErrorApplication(t *testing.T) {
	// TODO: encode additional value
	assert.Equal(t,
		InherentError{sc.NewVaryingData(sc.U8(3))},
		NewInherentErrorApplication(),
	)
}

func Test_InherentError_IsFatal_True(t *testing.T) {
	err := NewInherentErrorFatalErrorReported()

	assert.Equal(t, sc.Bool(true), err.IsFatal())
}

func Test_InherentError_IsFatal_False(t *testing.T) {
	err := NewInherentErrorInherentDataExists(identifier)

	assert.Equal(t, sc.Bool(false), err.IsFatal())
}

func Test_InherentError_Error_InherentDataExists(t *testing.T) {
	err := NewInherentErrorInherentDataExists(identifier)

	assert.Equal(t,
		"Inherent data already exists for identifier: [5 6]",
		err.Error(),
	)
}

func Test_InherentError_Error_DecodingFailed(t *testing.T) {
	err := NewInherentErrorDecodingFailed(identifier)

	assert.Equal(t,
		"Failed to decode inherent data for identifier: [5 6]",
		err.Error(),
	)
}

func Test_InherentError_Error_FatalErrorReported(t *testing.T) {
	err := NewInherentErrorFatalErrorReported()

	assert.Equal(t,
		"There was already a fatal error reported and no other errors are allowed",
		err.Error(),
	)
}

func Test_InherentError_Error_ErrorApplication(t *testing.T) {
	err := NewInherentErrorApplication()

	assert.Equal(t,
		"Inherent error application",
		err.Error(),
	)
}

func Test_InherentError_Error_TypeError(t *testing.T) {
	inherentError := InherentError{sc.NewVaryingData(sc.U8(4))}

	err := inherentError.Error()

	assert.Equal(t, "not a valid 'InherentError' type", err)
}
