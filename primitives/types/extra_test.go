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
	expectedSignedExtraOkBytes, _ = hex.DecodeString("00010000000003000000")

	expectedValidTransaction = ValidTransaction{
		Priority:  TransactionPriority(2),
		Requires:  sc.Sequence[TransactionTag]{},
		Provides:  sc.Sequence[TransactionTag]{},
		Longevity: TransactionLongevity(math.MaxUint64),
		Propagate: true,
	}

	expectedTransactionValidityError, _ = NewTransactionValidityError(NewUnknownTransactionCustomUnknownTransaction(sc.U8(0)))

	testExtraCheckMetadataType = MetadataType{
		Id:         sc.ToCompact(123456),
		Path:       sc.Sequence[sc.Str]{"frame_system", "extensions", "test_extra_check", "TestExtraCheck"},
		Params:     sc.Sequence[MetadataTypeParameter]{},
		Definition: sc.VaryingData{sc.U8(6), sc.ToCompact(123456)},
		Docs:       sc.Sequence[sc.Str]{"TestExtraCheck"},
	}

	signedExtraMetadataType = MetadataType{
		Id:         sc.ToCompact(97),
		Path:       sc.Sequence[sc.Str]{},
		Params:     sc.Sequence[MetadataTypeParameter]{},
		Definition: sc.VaryingData{sc.U8(4), sc.Sequence[sc.Compact]{sc.ToCompact(123456), sc.ToCompact(123456)}},
		Docs:       sc.Sequence[sc.Str]{"SignedExtra"},
	}

	expectedMetadataTypes = sc.Sequence[MetadataType]{
		testExtraCheckMetadataType,
		testExtraCheckMetadataType,
		signedExtraMetadataType,
	}

	metadataSignedExtension = MetadataSignedExtension{
		Identifier:       "TestExtraCheck",
		Type:             sc.ToCompact(123456),
		AdditionalSigned: sc.ToCompact(789),
	}

	expectedMetadataSignedExtensions = sc.Sequence[MetadataSignedExtension]{
		metadataSignedExtension,
		metadataSignedExtension,
	}
)

var (
	pre            = sc.Option[sc.Sequence[Pre]]{}
	preWithValue   = sc.NewOption[sc.Sequence[Pre]](sc.Sequence[Pre]{Pre{}, Pre{}})
	who            = AccountId[PublicKey]{}
	call           = testCall{}
	info           = &DispatchInfo{}
	length         = sc.Compact{}
	postInfo       = &PostDispatchInfo{}
	dispatchResult = &DispatchResult{}

	extraCheckOk1  = newTestExtraCheck(false, sc.U32(1))
	extraCheckOk2  = newTestExtraCheck(false, sc.U32(3))
	extraCheckErr1 = newTestExtraCheck(true, sc.U32(5))

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

	err := targetSignedExtraOk.Encode(buf)

	assert.NoError(t, err)
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

	assert.Equal(t, AdditionalSigned{sc.U32(1), sc.U32(3)}, result)
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

func Test_SignedExtra_PostDispatch_PreWithValue_Ok(t *testing.T) {
	err := targetSignedExtraOk.PostDispatch(preWithValue, info, postInfo, length, dispatchResult)

	assert.Nil(t, err)
}

func Test_SignedExtra_PostDispatch_PreWithValue_Err(t *testing.T) {
	err := targetSignedExtraErr.PostDispatch(preWithValue, info, postInfo, length, dispatchResult)

	assert.Equal(t, expectedTransactionValidityError, err)
}

func Test_SignedExtra_PostDispatch_Err(t *testing.T) {
	err := targetSignedExtraErr.PostDispatch(pre, info, postInfo, length, dispatchResult)

	assert.Equal(t, expectedTransactionValidityError, err)
}

func Test_SignedExtra_Metadata(t *testing.T) {
	metadataTypes, metadataSignedExtensions := targetSignedExtraOk.Metadata()

	assert.Equal(t, expectedMetadataTypes, metadataTypes)
	assert.Equal(t, expectedMetadataSignedExtensions, metadataSignedExtensions)
}
