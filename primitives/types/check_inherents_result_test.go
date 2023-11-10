package types

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedCheckInherentsResultOkBytes, _  = hex.DecodeString("010000")
	expectedCheckInherentsResultErrBytes, _ = hex.DecodeString("00010474657374696e67300402")
	invalidCheckInherentsResultBytes, _     = hex.DecodeString("01000101")
)

var (
	inherentIdentifier = [8]byte{'t', 'e', 's', 't', 'i', 'n', 'g', '0'}
	fatalError         = NewTimestampErrorInvalid()

	inherentDataOk = InherentData{
		data: map[[8]byte]sc.Sequence[sc.U8]{},
	}

	inherentDataOkWithIdent = InherentData{
		data: map[[8]byte]sc.Sequence[sc.U8]{
			inherentIdentifier: {1},
		},
	}

	inherentDataErr = InherentData{
		data: map[[8]byte]sc.Sequence[sc.U8]{
			inherentIdentifier: {2},
		},
	}
)

var (
	targetCheckInherentsResultOk = CheckInherentsResult{
		Okay:       true,
		FatalError: false,
		Errors:     inherentDataOk,
	}

	targetCheckInherentsResultOkWithIdent = CheckInherentsResult{
		Okay:       true,
		FatalError: false,
		Errors:     inherentDataOkWithIdent,
	}

	targetCheckInherentsResultErr = CheckInherentsResult{
		Okay:       false,
		FatalError: true,
		Errors:     inherentDataErr,
	}
)

func Test_NewCheckInherentsResult(t *testing.T) {
	assert.Equal(t, targetCheckInherentsResultOk, NewCheckInherentsResult())
}

func Test_CheckInherentsResult_Encode_Ok(t *testing.T) {
	buffer := &bytes.Buffer{}

	targetCheckInherentsResultOk.Encode(buffer)

	assert.Equal(t, expectedCheckInherentsResultOkBytes, buffer.Bytes())
}

func Test_CheckInherentsResult_Encode_Err(t *testing.T) {
	buffer := &bytes.Buffer{}

	err := targetCheckInherentsResultErr.Encode(buffer)

	assert.NoError(t, err)
	assert.Equal(t, expectedCheckInherentsResultErrBytes, buffer.Bytes())
}

func Test_CheckInherentsResult_Bytes_Ok(t *testing.T) {
	assert.Equal(t, expectedCheckInherentsResultOkBytes, targetCheckInherentsResultOk.Bytes())
}

func Test_CheckInherentsResult_Bytes_Err(t *testing.T) {
	assert.Equal(t, expectedCheckInherentsResultErrBytes, targetCheckInherentsResultErr.Bytes())
}

func Test_CheckInherentsResult_PutError_AlreadyWithError(t *testing.T) {
	err := targetCheckInherentsResultErr.PutError(inherentIdentifier, fatalError)

	assert.Equal(t, NewInherentErrorFatalErrorReported(), err)
}

func Test_CheckInherentsResult_PutError(t *testing.T) {
	err := targetCheckInherentsResultOk.PutError(inherentIdentifier, fatalError)

	assert.NoError(t, err)
	assert.Equal(t, sc.Bool(false), targetCheckInherentsResultOk.Okay)
	assert.Equal(t, sc.Bool(true), targetCheckInherentsResultOk.FatalError)
	assert.Equal(t, inherentDataErr, targetCheckInherentsResultOk.Errors)
}

func Test_CheckInherentsResult_PutError_ExistingIdentifier(t *testing.T) {
	err := targetCheckInherentsResultOkWithIdent.PutError(inherentIdentifier, NewInherentErrorApplication())

	assert.Equal(t, NewInherentErrorInherentDataExists(sc.BytesToSequenceU8(inherentIdentifier[:])), err)
}

func Test_DecodeCheckInherentsResult(t *testing.T) {
	buffer := bytes.NewBuffer(expectedCheckInherentsResultErrBytes)

	result, err := DecodeCheckInherentsResult(buffer)
	assert.NoError(t, err)

	assert.Equal(t, targetCheckInherentsResultErr, result)
}

func Test_DecodeCheckInherentsResult_DecodeError(t *testing.T) {
	buffer := bytes.NewBuffer(invalidCheckInherentsResultBytes)

	result, err := DecodeCheckInherentsResult(buffer)

	assert.Error(t, err)
	assert.Equal(t, io.EOF, err)
	assert.Equal(t, CheckInherentsResult{}, result)
}
