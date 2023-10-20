package types

import (
	"bytes"
	"encoding/hex"
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

	targetCheckInherentsResultErr.Encode(buffer)

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

	assert.Nil(t, err)
	assert.Equal(t, sc.Bool(false), targetCheckInherentsResultOk.Okay)
	assert.Equal(t, sc.Bool(true), targetCheckInherentsResultOk.FatalError)
	assert.Equal(t, inherentDataErr, targetCheckInherentsResultOk.Errors)
}

func Test_DecodeCheckInherentsResult(t *testing.T) {
	buffer := bytes.NewBuffer(expectedCheckInherentsResultErrBytes)

	result := DecodeCheckInherentsResult(buffer)

	assert.Equal(t, targetCheckInherentsResultErr, result)
}

func Test_DecodeCheckInherentsResult_Panics(t *testing.T) {
	buffer := bytes.NewBuffer(invalidCheckInherentsResultBytes)

	assert.PanicsWithValue(t, "failed to decode InherentData", func() {
		DecodeCheckInherentsResult(buffer)
	})

}
