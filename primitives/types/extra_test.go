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
	mdGenerator = NewMetadataTypeGenerator()

	lastIndex = mdGenerator.GetLastAvailableIndex()

	mdGeneratorAll = NewMetadataTypeGenerator()

	mdGeneratorSome = NewMetadataTypeGenerator()

	expectedExtraCheckMetadataId = lastIndex + 1

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
		Id:     sc.ToCompact(expectedExtraCheckMetadataId),
		Path:   sc.Sequence[sc.Str]{"primitives_types", "extensions", "test_extra_check", "testExtraCheck"},
		Params: sc.Sequence[MetadataTypeParameter]{},
		Definition: NewMetadataTypeDefinitionComposite(
			sc.Sequence[MetadataTypeDefinitionField]{
				NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesBool, "hasError"),
				NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesU32, "value"),
			},
		),
		Docs: sc.Sequence[sc.Str]{"testExtraCheck"},
	}

	signedExtraMetadataType = MetadataType{
		Id:         sc.ToCompact(metadata.SignedExtra),
		Path:       sc.Sequence[sc.Str]{},
		Params:     sc.Sequence[MetadataTypeParameter]{},
		Definition: MetadataTypeDefinition{sc.VaryingData{sc.U8(4), sc.Sequence[sc.Compact]{sc.ToCompact(expectedExtraCheckMetadataId), sc.ToCompact(expectedExtraCheckMetadataId)}}},
		Docs:       sc.Sequence[sc.Str]{"SignedExtra"},
	}

	expectedMetadataTypes = sc.Sequence[MetadataType]{ // TODO: Confirm that in such a case where we have repeating extras (so their metadata types will be equal), we don't expect to have duplicates in the returned metadataTypes
		testExtraCheckMetadataType,
		signedExtraMetadataType,
	}

	metadataSignedExtension = MetadataSignedExtension{
		Identifier:       "testExtraCheck",
		Type:             sc.ToCompact(expectedExtraCheckMetadataId),
		AdditionalSigned: sc.ToCompact(metadata.PrimitiveTypesU32),
	}

	expectedMetadataSignedExtensions = sc.Sequence[MetadataSignedExtension]{
		metadataSignedExtension,
		metadataSignedExtension,
	}
)

var (
	pre          = sc.Option[sc.Sequence[Pre]]{}
	preWithValue = sc.NewOption[sc.Sequence[Pre]](sc.Sequence[Pre]{Pre{}, Pre{}})
	who          = AccountId{}
	call         = testCall{}
	info         = &DispatchInfo{}
	length       = sc.Compact{}
	postInfo     = &PostDispatchInfo{}

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

	targetSignedExtraOk  = NewSignedExtra(extraChecksWithOk, mdGenerator)
	targetSignedExtraErr = NewSignedExtra(extraChecksWithErr, mdGenerator)
)

// Some different extra checks that contain a more complex structure (fields and additional signed types)
var (
	extraCheckEmpty   = newTestExtraCheckEmpty()
	extraCheckEra     = newtTestExtraCheckEra()
	extraCheckComplex = newtTestExtraCheckComplex()

	extraChecks = []SignedExtension{
		extraCheckEmpty,
		extraCheckEra,
		extraCheckComplex,
	}

	targetSignedExtraAll  = NewSignedExtra(extraChecks, mdGeneratorAll)
	targetSignedExtraSome = NewSignedExtra(extraChecks, mdGeneratorSome)
)

func Test_NewSignedExtra(t *testing.T) {
	assert.Equal(t, signedExtra{extras: extraChecksWithOk, mdGenerator: mdGenerator}, targetSignedExtraOk)
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

	assert.Equal(t, signedExtra{extras: extraChecksWithOk, mdGenerator: mdGenerator}, targetSignedExtraOk)
}

func Test_SignedExtra_DeepCopy(t *testing.T) {
	target := NewSignedExtra(extraChecks, mdGenerator)

	result := target.DeepCopy()
	assert.Equal(t, result, target)
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
	err := targetSignedExtraOk.PostDispatch(pre, info, postInfo, length, nil)

	assert.Nil(t, err)
}

func Test_SignedExtra_PostDispatch_PreWithValue_Ok(t *testing.T) {
	err := targetSignedExtraOk.PostDispatch(preWithValue, info, postInfo, length, nil)

	assert.Nil(t, err)
}

func Test_SignedExtra_PostDispatch_PreWithValue_Err(t *testing.T) {
	err := targetSignedExtraErr.PostDispatch(preWithValue, info, postInfo, length, nil)

	assert.Equal(t, expectedTransactionValidityError, err)
}

func Test_SignedExtra_PostDispatch_Err(t *testing.T) {
	err := targetSignedExtraErr.PostDispatch(pre, info, postInfo, length, nil)

	assert.Equal(t, expectedTransactionValidityError, err)
}

func Test_SignedExtra_Metadata(t *testing.T) {
	metadataSignedExtensions := targetSignedExtraOk.Metadata()
	metadataTypes := mdGenerator.GetMetadataTypes()

	assert.Equal(t, expectedMetadataTypes, metadataTypes)
	assert.Equal(t, expectedMetadataSignedExtensions, metadataSignedExtensions)
}

func Test_SignedExtra_Metadata_Complex_All(t *testing.T) {
	expectedH512MetadataId := lastIndex + 1
	expectedEd25519PublicKeyMetadataId := lastIndex + 2
	expectedWeightId := lastIndex + 3

	// A map that contains the ids of all additional signed complex checks
	metadataIdsComplexAll := mdGeneratorAll.GetIdsMap()

	metadataIdsComplexAll["H512"] = expectedH512MetadataId
	metadataIdsComplexAll["Ed25519PublicKey"] = expectedEd25519PublicKeyMetadataId
	metadataIdsComplexAll["Weight"] = expectedWeightId

	lastIndexComplexChecks := mdGeneratorAll.GetLastAvailableIndex()

	// expected metadata type ids in the order they will be generated when Metadata() is invoked
	expectedEmptyCheckMetadataId := lastIndexComplexChecks + 1
	expectedCheckEraMetadataId := lastIndexComplexChecks + 2
	expectedEraMetadataId := lastIndexComplexChecks + 3
	expectedCheckComplexMetadataId := lastIndexComplexChecks + 4
	expectedTupleAdditionalSignedMetadataId := lastIndexComplexChecks + 5

	metadataSignedExtensionEmpty := MetadataSignedExtension{
		Identifier:       "testExtraCheckEmpty",
		Type:             sc.ToCompact(expectedEmptyCheckMetadataId),
		AdditionalSigned: sc.ToCompact(metadata.TypesEmptyTuple),
	}

	metadataSignedExtensionEra := MetadataSignedExtension{
		Identifier:       "testExtraCheckEra",
		Type:             sc.ToCompact(expectedCheckEraMetadataId),
		AdditionalSigned: sc.ToCompact(metadata.TypesH256),
	}

	metadataSignedExtensionComplex := MetadataSignedExtension{
		Identifier:       "testExtraCheckComplex",
		Type:             sc.ToCompact(expectedCheckComplexMetadataId),
		AdditionalSigned: sc.ToCompact(expectedTupleAdditionalSignedMetadataId),
	}

	expectedExtensionsEmptyEraComplex := sc.Sequence[MetadataSignedExtension]{
		metadataSignedExtensionEmpty,
		metadataSignedExtensionEra,
		metadataSignedExtensionComplex,
	}

	tupleAdditionalSignedMetadataType := NewMetadataType(expectedTupleAdditionalSignedMetadataId, "H256U32U64H512Ed25519PublicKeyWeight",
		NewMetadataTypeDefinitionTuple(sc.Sequence[sc.Compact]{
			sc.ToCompact(metadata.TypesH256),
			sc.ToCompact(metadata.PrimitiveTypesU32),
			sc.ToCompact(metadata.PrimitiveTypesU64),
			sc.ToCompact(expectedH512MetadataId),
			sc.ToCompact(expectedEd25519PublicKeyMetadataId),
			sc.ToCompact(expectedWeightId)}))

	signedExtraMdType := MetadataType{
		Id:         sc.ToCompact(metadata.SignedExtra),
		Path:       sc.Sequence[sc.Str]{},
		Params:     sc.Sequence[MetadataTypeParameter]{},
		Definition: MetadataTypeDefinition{sc.VaryingData{sc.U8(4), sc.Sequence[sc.Compact]{sc.ToCompact(expectedEmptyCheckMetadataId), sc.ToCompact(expectedCheckEraMetadataId), sc.ToCompact(expectedCheckComplexMetadataId)}}},
		Docs:       sc.Sequence[sc.Str]{"SignedExtra"},
	}

	eraMetadataType := NewMetadataType(
		expectedEraMetadataId,
		"Era",
		NewMetadataTypeDefinitionComposite(sc.Sequence[MetadataTypeDefinitionField]{
			NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesBool, "IsImmortal"),
			NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesU64, "EraPeriod"),
			NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesU64, "EraPhase")}),
	)

	testExtraCheckEmptyMetadataType := MetadataType{
		Id:     sc.ToCompact(expectedEmptyCheckMetadataId),
		Path:   sc.Sequence[sc.Str]{"primitives_types", "extensions", "test_extra_check_empty", "testExtraCheckEmpty"},
		Params: sc.Sequence[MetadataTypeParameter]{},
		Definition: NewMetadataTypeDefinitionComposite(
			sc.Sequence[MetadataTypeDefinitionField]{},
		),
		Docs: sc.Sequence[sc.Str]{"testExtraCheckEmpty"},
	}

	testExtraCheckEraMetadataType := MetadataType{
		Id:     sc.ToCompact(expectedCheckEraMetadataId),
		Path:   sc.Sequence[sc.Str]{"primitives_types", "extensions", "test_extra_check_era", "testExtraCheckEra"},
		Params: sc.Sequence[MetadataTypeParameter]{},
		Definition: NewMetadataTypeDefinitionComposite(
			sc.Sequence[MetadataTypeDefinitionField]{
				NewMetadataTypeDefinitionFieldWithName(expectedEraMetadataId, "era"),
			},
		),
		Docs: sc.Sequence[sc.Str]{"testExtraCheckEra"},
	}

	testExtraCheckComplexMetadataType := MetadataType{
		Id:     sc.ToCompact(expectedCheckComplexMetadataId),
		Path:   sc.Sequence[sc.Str]{"primitives_types", "extensions", "test_extra_check_complex", "testExtraCheckComplex"},
		Params: sc.Sequence[MetadataTypeParameter]{},
		Definition: NewMetadataTypeDefinitionComposite(
			sc.Sequence[MetadataTypeDefinitionField]{
				NewMetadataTypeDefinitionFieldWithName(expectedEraMetadataId, "era"),
				NewMetadataTypeDefinitionFieldWithName(metadata.TypesH256, "hash"),
				NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesU64, "value"),
			},
		),
		Docs: sc.Sequence[sc.Str]{"testExtraCheckComplex"},
	}

	expectedMetadataTypesEmptyEraComplexAll := sc.Sequence[MetadataType]{
		testExtraCheckEmptyMetadataType,
		eraMetadataType, // during the process of generating the metadata of testExtraCheckEra that has a field "Era", this metadata type was generated so we expect it included in the sequence of metadata types
		testExtraCheckEraMetadataType,
		testExtraCheckComplexMetadataType,
		tupleAdditionalSignedMetadataType, // same here
		signedExtraMdType,
	}

	metadataSignedExtensions := targetSignedExtraAll.Metadata()
	metadataTypes := mdGeneratorAll.GetMetadataTypes()

	assert.Equal(t, expectedMetadataTypesEmptyEraComplexAll, metadataTypes)
	assert.Equal(t, expectedExtensionsEmptyEraComplex, metadataSignedExtensions)
}

func Test_SignedExtra_Metadata_Complex_Some(t *testing.T) {
	// A map that does not contain the ids of types H512 and Ed25519PublicKey hence they need to be generated
	metadataIdsComplexSome := mdGeneratorSome.GetIdsMap()

	metadataIdsComplexSome["H512"] = metadata.TypesFixedSequence64U8
	metadataIdsComplexSome["Ed25519PublicKey"] = metadata.TypesFixedSequence64U8

	lastIndexComplexChecksSome := mdGeneratorSome.GetLastAvailableIndex()

	expectedEmptyCheckMetadataIdSome := lastIndexComplexChecksSome + 1
	expectedCheckEraMetadataIdSome := lastIndexComplexChecksSome + 2
	expectedEraMetadataIdSome := lastIndexComplexChecksSome + 3
	expectedExtraCheckComplexIdSome := lastIndexComplexChecksSome + 4
	expectedTupleAdditionalSignedMetadataIdSome := lastIndexComplexChecksSome + 5

	testExtraCheckEmptyMetadataTypeSome := MetadataType{
		Id:         sc.ToCompact(expectedEmptyCheckMetadataIdSome),
		Path:       sc.Sequence[sc.Str]{"primitives_types", "extensions", "test_extra_check_empty", "testExtraCheckEmpty"},
		Params:     sc.Sequence[MetadataTypeParameter]{},
		Definition: NewMetadataTypeDefinitionComposite(sc.Sequence[MetadataTypeDefinitionField]{}),
		Docs:       sc.Sequence[sc.Str]{"testExtraCheckEmpty"},
	}

	testExtraCheckEraMetadataTypeSome := MetadataType{
		Id:     sc.ToCompact(expectedCheckEraMetadataIdSome),
		Path:   sc.Sequence[sc.Str]{"primitives_types", "extensions", "test_extra_check_era", "testExtraCheckEra"},
		Params: sc.Sequence[MetadataTypeParameter]{},
		Definition: NewMetadataTypeDefinitionComposite(
			sc.Sequence[MetadataTypeDefinitionField]{
				NewMetadataTypeDefinitionFieldWithName(expectedEraMetadataIdSome, "era"),
			},
		),
		Docs: sc.Sequence[sc.Str]{"testExtraCheckEra"},
	}

	eraMetadataTypeSome := NewMetadataType(
		expectedEraMetadataIdSome,
		"Era",
		NewMetadataTypeDefinitionComposite(sc.Sequence[MetadataTypeDefinitionField]{
			NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesBool, "IsImmortal"),
			NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesU64, "EraPeriod"),
			NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesU64, "EraPhase"),
		}),
	)

	tupleAdditionalSignedMetadataTypeSome := NewMetadataType(expectedTupleAdditionalSignedMetadataIdSome, "H256U32U64H512Ed25519PublicKeyWeight",
		NewMetadataTypeDefinitionTuple(sc.Sequence[sc.Compact]{
			sc.ToCompact(metadata.TypesH256),
			sc.ToCompact(metadata.PrimitiveTypesU32),
			sc.ToCompact(metadata.PrimitiveTypesU64),
			sc.ToCompact(metadata.TypesFixedSequence64U8),
			sc.ToCompact(metadata.TypesFixedSequence64U8),
			sc.ToCompact(metadata.TypesWeight)}))

	testExtraCheckComplexMetadataTypeSome := MetadataType{
		Id:     sc.ToCompact(expectedExtraCheckComplexIdSome),
		Path:   sc.Sequence[sc.Str]{"primitives_types", "extensions", "test_extra_check_complex", "testExtraCheckComplex"},
		Params: sc.Sequence[MetadataTypeParameter]{},
		Definition: NewMetadataTypeDefinitionComposite(
			sc.Sequence[MetadataTypeDefinitionField]{
				NewMetadataTypeDefinitionFieldWithName(expectedEraMetadataIdSome, "era"),
				NewMetadataTypeDefinitionFieldWithName(metadata.TypesH256, "hash"),
				NewMetadataTypeDefinitionFieldWithName(metadata.PrimitiveTypesU64, "value"),
			},
		),
		Docs: sc.Sequence[sc.Str]{"testExtraCheckComplex"},
	}

	signedExtraMdTypeSome := MetadataType{
		Id:         sc.ToCompact(metadata.SignedExtra),
		Path:       sc.Sequence[sc.Str]{},
		Params:     sc.Sequence[MetadataTypeParameter]{},
		Definition: MetadataTypeDefinition{sc.VaryingData{sc.U8(4), sc.Sequence[sc.Compact]{sc.ToCompact(expectedEmptyCheckMetadataIdSome), sc.ToCompact(expectedCheckEraMetadataIdSome), sc.ToCompact(expectedExtraCheckComplexIdSome)}}},
		Docs:       sc.Sequence[sc.Str]{"SignedExtra"},
	}

	expectedMetadataTypesEmptyEraComplexSome := sc.Sequence[MetadataType]{
		testExtraCheckEmptyMetadataTypeSome,
		eraMetadataTypeSome,
		testExtraCheckEraMetadataTypeSome,
		testExtraCheckComplexMetadataTypeSome,
		tupleAdditionalSignedMetadataTypeSome,
		signedExtraMdTypeSome,
	}

	metadataSignedExtensionEmptySome := MetadataSignedExtension{
		Identifier:       "testExtraCheckEmpty",
		Type:             sc.ToCompact(expectedEmptyCheckMetadataIdSome),
		AdditionalSigned: sc.ToCompact(metadata.TypesEmptyTuple),
	}

	metadataSignedExtensionEraSome := MetadataSignedExtension{
		Identifier:       "testExtraCheckEra",
		Type:             sc.ToCompact(expectedCheckEraMetadataIdSome),
		AdditionalSigned: sc.ToCompact(metadata.TypesH256),
	}

	metadataSignedExtensionComplexSome := MetadataSignedExtension{
		Identifier:       "testExtraCheckComplex",
		Type:             sc.ToCompact(expectedExtraCheckComplexIdSome),
		AdditionalSigned: sc.ToCompact(expectedTupleAdditionalSignedMetadataIdSome),
	}

	expectedExtensionsEmptyEraComplexSome := sc.Sequence[MetadataSignedExtension]{
		metadataSignedExtensionEmptySome,
		metadataSignedExtensionEraSome,
		metadataSignedExtensionComplexSome,
	}

	metadataSignedExtensions := targetSignedExtraSome.Metadata()
	metadataTypes := mdGeneratorSome.GetMetadataTypes()

	assert.Equal(t, expectedMetadataTypesEmptyEraComplexSome, metadataTypes)
	assert.Equal(t, expectedExtensionsEmptyEraComplexSome, metadataSignedExtensions)
}
