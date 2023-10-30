package types

import (
	"bytes"
	"encoding/hex"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedSignedPayloadBytes, _ = hex.DecodeString("01020304000100000001000000")

	expectedAdditionalSigned = sc.NewVaryingData(sc.U32(1))

	expectedTransactionValidityErr = NewTransactionValidityError(NewUnknownTransactionCustomUnknownTransaction(sc.U8(0)))
)

var (
	dispatchCall = testCall{
		Callable: Callable{
			ModuleId:   1,
			FunctionId: 2,
			Arguments:  sc.NewVaryingData(sc.U8(3), sc.U8(4)),
		},
	}

	extraCheckOk  = newTestExtraCheck(false, sc.U32(1))
	extraCheckErr = newTestExtraCheck(true, sc.U32(5))

	extraChecksWithOk1 = []SignedExtension{
		extraCheckOk,
	}

	extraChecksWithErr1 = []SignedExtension{
		extraCheckOk,
		extraCheckErr,
	}

	signedExtraWithOk  = NewSignedExtra(extraChecksWithOk1)
	signedExtraWithErr = NewSignedExtra(extraChecksWithErr1)
)

var (
	targetSignedPayload = signedPayload{
		call:             dispatchCall,
		extra:            signedExtraWithOk,
		additionalSigned: expectedAdditionalSigned,
	}
)

func Test_NewSignedPayload_Ok(t *testing.T) {
	result, err := NewSignedPayload(dispatchCall, signedExtraWithOk)

	assert.Nil(t, err)
	assert.Equal(t, targetSignedPayload, result)
}

func Test_NewSignedPayload_Err(t *testing.T) {
	result, err := NewSignedPayload(dispatchCall, signedExtraWithErr)

	expectedSignedPayload := signedPayload{}

	assert.Equal(t, expectedTransactionValidityErr, err)
	assert.Equal(t, expectedSignedPayload, result)
}

func Test_SignedPayload_Encode(t *testing.T) {
	buf := &bytes.Buffer{}

	targetSignedPayload.Encode(buf)

	assert.Equal(t, expectedSignedPayloadBytes, buf.Bytes())
}

func Test_SignedPayload_Bytes(t *testing.T) {
	assert.Equal(t, expectedSignedPayloadBytes, targetSignedPayload.Bytes())
}

func Test_SignedPayload_AdditionalSigned(t *testing.T) {
	assert.Equal(t, expectedAdditionalSigned, targetSignedPayload.AdditionalSigned())
}

func Test_SignedPayload_Call(t *testing.T) {
	assert.Equal(t, dispatchCall, targetSignedPayload.Call())
}

func Test_SignedPayload_Extra(t *testing.T) {
	assert.Equal(t, signedExtraWithOk, targetSignedPayload.Extra())
}
