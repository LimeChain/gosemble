package types

import (
	"bytes"
	"encoding/hex"
	"math"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/assert"
)

var (
	expectedSignedExtraOkBytes, _ = hex.DecodeString("010002030004")

	expectedValidTransaction = ValidTransaction{
		Priority:  TransactionPriority(2),
		Requires:  sc.Sequence[TransactionTag]{},
		Provides:  sc.Sequence[TransactionTag]{},
		Longevity: TransactionLongevity(math.MaxUint64),
		Propagate: true,
	}

	expectedTransactionValidityError = NewTransactionValidityError(
		NewUnknownTransactionCustomUnknownTransaction(sc.U8(0)),
	)
)

var (
	pre            = sc.Option[sc.Sequence[Pre]]{}
	who            = Address32{}
	call           = testCall{}
	info           = &DispatchInfo{}
	length         = sc.Compact{}
	postInfo       = &PostDispatchInfo{}
	dispatchResult = &DispatchResult{}

	extraCheckOk1  = NewTestExtraCheck(false, sc.U16(1), sc.U8(2))
	extraCheckOk2  = NewTestExtraCheck(false, sc.U16(3), sc.U8(4))
	extraCheckErr1 = NewTestExtraCheck(true, sc.U16(5), sc.U8(6))

	extraChecksWithOk = []SignedExtension{
		extraCheckOk1,
		extraCheckOk2,
	}

	extraChecksWithErr = []SignedExtension{
		extraCheckOk1,
		extraCheckErr1,
		extraCheckOk2,
	}

	targetSignedExtraOk  = NewSignedExtra(extraChecksWithOk)
	targetSignedExtraErr = NewSignedExtra(extraChecksWithErr)
)

func Test_NewSignedExtra(t *testing.T) {
	assert.Equal(t, signedExtra{extras: extraChecksWithOk}, targetSignedExtraOk)
}

func Test_SignedExtra_Encode(t *testing.T) {
	buf := &bytes.Buffer{}

	targetSignedExtraOk.Encode(buf)

	assert.Equal(t, expectedSignedExtraOkBytes, buf.Bytes())
}

func Test_SignedExtra_Bytes(t *testing.T) {
	assert.Equal(t, expectedSignedExtraOkBytes, targetSignedExtraOk.Bytes())
}

func Test_SignedExtra_Decode(t *testing.T) {
	buf := bytes.NewBuffer(expectedSignedExtraOkBytes)

	targetSignedExtraOk.Decode(buf)

	assert.Equal(t, signedExtra{extras: extraChecksWithOk}, targetSignedExtraOk)
}

func Test_SignedExtra_AdditionalSigned_Ok(t *testing.T) {
	result, err := targetSignedExtraOk.AdditionalSigned()

	assert.Equal(t, AdditionalSigned{sc.U16(1), sc.U8(2), sc.U16(3), sc.U8(4)}, result)
	assert.Nil(t, err)
}

func Test_SignedExtra_AdditionalSigned_Err(t *testing.T) {
	result, err := targetSignedExtraErr.AdditionalSigned()

	assert.Nil(t, result)
	assert.Equal(t, expectedTransactionValidityError, err)
}

func Test_SignedExtra_Validate_Ok(t *testing.T) {
	result, err := targetSignedExtraOk.Validate(who, call, info, length)

	assert.Equal(t, expectedValidTransaction, result)
	assert.Nil(t, err)
}

func Test_SignedExtra_Validate_Err(t *testing.T) {
	result, err := targetSignedExtraErr.Validate(who, call, info, length)

	assert.Equal(t, ValidTransaction{}, result)
	assert.Equal(t, expectedTransactionValidityError, err)
}

func Test_SignedExtra_ValidateUnsigned_Ok(t *testing.T) {
	result, err := targetSignedExtraOk.ValidateUnsigned(call, info, length)

	assert.Equal(t, expectedValidTransaction, result)
	assert.Nil(t, err)
}

func Test_SignedExtra_ValidateUnsigned_Err(t *testing.T) {
	result, err := targetSignedExtraErr.ValidateUnsigned(call, info, length)

	assert.Equal(t, ValidTransaction{}, result)
	assert.Equal(t, expectedTransactionValidityError, err)
}

func Test_SignedExtra_PreDispatch_Ok(t *testing.T) {
	result, err := targetSignedExtraOk.PreDispatch(who, call, info, length)

	assert.Equal(t, sc.Sequence[sc.VaryingData]{Pre{}, Pre{}}, result)
	assert.Nil(t, err)
}

func Test_SignedExtra_PreDispatch_Err(t *testing.T) {
	result, err := targetSignedExtraErr.PreDispatch(who, call, info, length)

	assert.Equal(t, sc.Sequence[sc.VaryingData](nil), result)
	assert.Equal(t, expectedTransactionValidityError, err)
}

func Test_SignedExtra_PreDispatchUnsigned_Ok(t *testing.T) {
	err := targetSignedExtraOk.PreDispatchUnsigned(call, info, length)

	assert.Nil(t, err)
}

func Test_SignedExtra_PreDispatchUnsigned_Err(t *testing.T) {
	err := targetSignedExtraErr.PreDispatchUnsigned(call, info, length)

	assert.Equal(t, expectedTransactionValidityError, err)
}

func Test_SignedExtra_PostDispatch_Ok(t *testing.T) {
	err := targetSignedExtraOk.PostDispatch(pre, info, postInfo, length, dispatchResult)

	assert.Nil(t, err)
}

func Test_SignedExtra_PostDispatch_Err(t *testing.T) {
	err := targetSignedExtraErr.PostDispatch(pre, info, postInfo, length, dispatchResult)

	assert.Equal(t, expectedTransactionValidityError, err)
}
