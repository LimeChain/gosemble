package types

import (
	"bytes"
	"encoding/hex"
	"math"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/stretchr/testify/assert"
)

var (
	metadataIds = map[string]int{
		"Bool":             metadata.PrimitiveTypesBool,
		"String":           metadata.PrimitiveTypesString,
		"U8":               metadata.PrimitiveTypesU8,
		"U16":              metadata.PrimitiveTypesU16,
		"U32":              metadata.PrimitiveTypesU32,
		"U64":              metadata.PrimitiveTypesU64,
		"U128":             metadata.PrimitiveTypesU128,
		"U256":             metadata.PrimitiveTypesU256,
		"I8":               metadata.PrimitiveTypesI8,
		"I16":              metadata.PrimitiveTypesI16,
		"I32":              metadata.PrimitiveTypesI32,
		"I64":              metadata.PrimitiveTypesI64,
		"I128":             metadata.PrimitiveTypesI128,
		"H256":             metadata.TypesH256,
		"H512":             15,
		"Ed25519PublicKey": 16,
	}

	lastIndex = len(metadataIds)
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

	expectedTransactionValidityError = NewTransactionValidityError(NewUnknownTransactionCustomUnknownTransaction(sc.U8(0)))

	testExtraCheckMetadataType = MetadataType{
		Id:     sc.ToCompact(lastIndex + 1),
		Path:   sc.Sequence[sc.Str]{"extensions", "test_extra_check", "testExtraCheck"},
		Params: sc.Sequence[MetadataTypeParameter]{},
		Definition: NewMetadataTypeDefinitionComposite(
			sc.Sequence[MetadataTypeDefinitionField]{
				NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesBool, "Bool"),
				NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesU32, "U32"),
			},
		),
		Docs: sc.Sequence[sc.Str]{"testExtraCheck"},
	}

	testExtraCheckEmptyMetadataType = MetadataType{
		Id:     sc.ToCompact(lastIndex + 1),
		Path:   sc.Sequence[sc.Str]{"extensions", "test_extra_check_empty", "testExtraCheckEmpty"},
		Params: sc.Sequence[MetadataTypeParameter]{},
		Definition: NewMetadataTypeDefinitionComposite(
			nil,
		),
		Docs: sc.Sequence[sc.Str]{"testExtraCheckEmpty"},
	}

	testExtraCheckEraMetadataType = MetadataType{
		Id:     sc.ToCompact(lastIndex + 2),
		Path:   sc.Sequence[sc.Str]{"extensions", "test_extra_check_era", "testExtraCheckEra"},
		Params: sc.Sequence[MetadataTypeParameter]{},
		Definition: NewMetadataTypeDefinitionComposite(
			sc.Sequence[MetadataTypeDefinitionField]{
				NewMetadataTypeDefinitionFieldWithName(lastIndex+3, "Era"),
			},
		),
		Docs: sc.Sequence[sc.Str]{"testExtraCheckEra"},
	}

	testExtraCheckComplexMetadataType = MetadataType{
		Id:     sc.ToCompact(lastIndex + 4),
		Path:   sc.Sequence[sc.Str]{"extensions", "test_extra_check_complex", "testExtraCheckComplex"},
		Params: sc.Sequence[MetadataTypeParameter]{},
		Definition: NewMetadataTypeDefinitionComposite(
			sc.Sequence[MetadataTypeDefinitionField]{
				NewMetadataTypeDefinitionFieldWithName(lastIndex+3, "Era"),
				NewMetadataTypeDefinitionFieldWithName(metadata.TypesH256, "H256"),
				NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesU64, "U64"),
			},
		),
		Docs: sc.Sequence[sc.Str]{"testExtraCheckComplex"},
	}

	eraMetadataType = NewMetadataType(
		lastIndex+3,
		"Era",
		NewMetadataTypeDefinitionComposite(sc.Sequence[MetadataTypeDefinitionField]{}),
	)

	signedExtraMetadataType = MetadataType{
		Id:         sc.ToCompact(97),
		Path:       sc.Sequence[sc.Str]{},
		Params:     sc.Sequence[MetadataTypeParameter]{},
		Definition: MetadataTypeDefinition{sc.VaryingData{sc.U8(4), sc.Sequence[sc.Compact]{sc.ToCompact(lastIndex + 1), sc.ToCompact(lastIndex + 1)}}},
		Docs:       sc.Sequence[sc.Str]{"SignedExtra"},
	}

	signedExtraMetadataTypeDifferent = MetadataType{
		Id:         sc.ToCompact(97),
		Path:       sc.Sequence[sc.Str]{},
		Params:     sc.Sequence[MetadataTypeParameter]{},
		Definition: MetadataTypeDefinition{sc.VaryingData{sc.U8(4), sc.Sequence[sc.Compact]{sc.ToCompact(lastIndex + 1), sc.ToCompact(lastIndex + 2), sc.ToCompact(lastIndex + 4)}}},
		Docs:       sc.Sequence[sc.Str]{"SignedExtra"},
	}

	expectedMetadataTypes = sc.Sequence[MetadataType]{
		testExtraCheckMetadataType,
		testExtraCheckMetadataType,
		signedExtraMetadataType,
	}

	expectedMetadataTypesDifferent = sc.Sequence[MetadataType]{
		testExtraCheckEmptyMetadataType,
		eraMetadataType, // during the process of generating the metadata of testExtraCheckEra that has a field "Era", this metadata type was generated
		testExtraCheckEraMetadataType,
		NewMetadataType(lastIndex+5, "H256U32U64H512Ed25519PublicKey", NewMetadataTypeDefinitionTuple(sc.Sequence[sc.Compact]{sc.ToCompact(metadata.TypesH256), sc.ToCompact(metadata.PrimitiveTypesU32), sc.ToCompact(metadata.PrimitiveTypesU64), sc.ToCompact(15), sc.ToCompact(16)})),
		testExtraCheckComplexMetadataType,
		signedExtraMetadataTypeDifferent,
	}

	metadataSignedExtension = MetadataSignedExtension{
		Identifier:       "testExtraCheck",
		Type:             sc.ToCompact(lastIndex + 1),
		AdditionalSigned: sc.ToCompact(metadata.PrimitiveTypesU32),
	}

	metadataSignedExtensionEmpty = MetadataSignedExtension{
		Identifier:       "testExtraCheckEmpty",
		Type:             sc.ToCompact(lastIndex + 1),
		AdditionalSigned: sc.ToCompact(int(metadata.TypesEmptyTuple)),
	}

	metadataSignedExtensionEra = MetadataSignedExtension{
		Identifier:       "testExtraCheckEra",
		Type:             sc.ToCompact(lastIndex + 2),
		AdditionalSigned: sc.ToCompact(metadata.TypesH256),
	}

	metadataSignedExtensionComplex = MetadataSignedExtension{
		Identifier:       "testExtraCheckComplex",
		Type:             sc.ToCompact(lastIndex + 4),
		AdditionalSigned: sc.ToCompact(lastIndex + 5),
	}

	expectedMetadataSignedExtensions = sc.Sequence[MetadataSignedExtension]{
		metadataSignedExtension,
		metadataSignedExtension,
	}

	expectedMetadataSignedExtensionsDifferent = sc.Sequence[MetadataSignedExtension]{
		metadataSignedExtensionEmpty,
		metadataSignedExtensionEra,
		metadataSignedExtensionComplex,
	}
)

var (
	pre            = sc.Option[sc.Sequence[Pre]]{}
	preWithValue   = sc.NewOption[sc.Sequence[Pre]](sc.Sequence[Pre]{Pre{}, Pre{}})
	who            = AccountId{}
	call           = testCall{}
	info           = &DispatchInfo{}
	length         = sc.Compact{}
	postInfo       = &PostDispatchInfo{}
	dispatchResult = &DispatchResult{}

	extraCheckOk1  = newTestExtraCheck(false, sc.U32(1))
	extraCheckOk2  = newTestExtraCheck(false, sc.U32(3))
	extraCheckErr1 = newTestExtraCheck(true, sc.U32(5))

	extraCheckEmpty   = newTestExtraCheckEmpty()
	extraCheckEra     = newtTestExtraCheckEra()
	extraCheckComplex = newtTestExtraCheckComplex()

	extraChecksWithOk = []SignedExtension{
		extraCheckOk1,
		extraCheckOk2,
	}

	extraChecksWithErr = []SignedExtension{
		extraCheckOk1,
		extraCheckErr1,
		extraCheckOk2,
	}

	extraChecks = []SignedExtension{
		extraCheckEmpty,
		extraCheckEra,
		extraCheckComplex,
	}

	targetSignedExtraOk  = NewSignedExtra(extraChecksWithOk)
	targetSignedExtraErr = NewSignedExtra(extraChecksWithErr)

	targetSignedExtraDifferent = NewSignedExtra(extraChecks)
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
	metadataTypes, metadataSignedExtensions := targetSignedExtraOk.Metadata(metadataIds)

	assert.Equal(t, expectedMetadataTypes, metadataTypes)
	assert.Equal(t, expectedMetadataSignedExtensions, metadataSignedExtensions)
}

func Test_SignedExtra_Metadata_DifferentTypes(t *testing.T) {
	metadataTypes, metadataSignedExtensions := targetSignedExtraDifferent.Metadata(metadataIds)

	assert.Equal(t, expectedMetadataTypesDifferent, metadataTypes)
	assert.Equal(t, expectedMetadataSignedExtensionsDifferent, metadataSignedExtensions)
}
