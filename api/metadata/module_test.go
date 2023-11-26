package metadata

import (
	"testing"

	"github.com/ChainSafe/gossamer/lib/common"
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

var (
	dataPtr    = int32(0)
	dataLen    = int32(1)
	ptrAndSize = int64(5)

	mdTypes = sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataType(1, "Test Metadata type", primitives.NewMetadataTypeDefinitionPrimitive(primitives.MetadataDefinitionPrimitiveBoolean)),
	}

	mdModules14 = sc.Sequence[primitives.MetadataModuleV14]{
		primitives.MetadataModuleV14{
			Name:      "Aura",
			Storage:   sc.Option[primitives.MetadataModuleStorage]{},
			Call:      sc.NewOption[sc.Compact](nil),
			CallDef:   sc.NewOption[primitives.MetadataDefinitionVariant](nil),
			Event:     sc.NewOption[sc.Compact](nil),
			EventDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
			Constants: sc.Sequence[primitives.MetadataModuleConstant]{},
			Error:     sc.NewOption[sc.Compact](nil),
			ErrorDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
			Index:     0,
		},
	}

	mdModules15 = sc.Sequence[primitives.MetadataModuleV15]{
		primitives.MetadataModuleV15{
			Name:      "Aura",
			Storage:   sc.Option[primitives.MetadataModuleStorage]{},
			Call:      sc.NewOption[sc.Compact](nil),
			CallDef:   sc.NewOption[primitives.MetadataDefinitionVariant](nil),
			Event:     sc.NewOption[sc.Compact](nil),
			EventDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
			Constants: sc.Sequence[primitives.MetadataModuleConstant]{},
			Error:     sc.NewOption[sc.Compact](nil),
			ErrorDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
			Index:     0,
			Docs:      sc.Sequence[sc.Str]{},
		},
	}

	signedExtensions = sc.Sequence[primitives.MetadataSignedExtension]{
		primitives.NewMetadataSignedExtension("ChargeTransactionPayment", metadata.ChargeTransactionPayment, metadata.TypesEmptyTuple),
	}

	mdExtrinsic = primitives.MetadataExtrinsicV14{
		Type:             sc.ToCompact(metadata.UncheckedExtrinsic),
		Version:          types.ExtrinsicFormatVersion,
		SignedExtensions: signedExtensions,
	}

	mdExtrinsic15 = primitives.MetadataExtrinsicV15{
		Version:          types.ExtrinsicFormatVersion,
		Address:          sc.ToCompact(metadata.TypesMultiAddress),
		Call:             sc.ToCompact(metadata.RuntimeCall),
		Signature:        sc.ToCompact(metadata.TypesMultiSignature),
		Extra:            sc.ToCompact(metadata.SignedExtra),
		SignedExtensions: signedExtensions,
	}

	mdRuntimeApi = sc.Sequence[primitives.RuntimeApiMetadata]{
		primitives.RuntimeApiMetadata{
			Name: ApiModuleName,
			Methods: sc.Sequence[primitives.RuntimeApiMethodMetadata]{
				primitives.RuntimeApiMethodMetadata{
					Name:   "metadata",
					Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{},
					Output: sc.ToCompact(metadata.TypesOpaqueMetadata),
					Docs:   sc.Sequence[sc.Str]{" Returns the metadata of a runtime."},
				},
				primitives.RuntimeApiMethodMetadata{
					Name: "metadata_at_version",
					Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{
						primitives.RuntimeApiMethodParamMetadata{
							Name: "version",
							Type: sc.ToCompact(metadata.PrimitiveTypesU32),
						},
					},
					Output: sc.ToCompact(metadata.TypeOptionOpaqueMetadata),
					Docs: sc.Sequence[sc.Str]{" Returns the metadata at a given version.",
						"",
						" If the given `version` isn't supported, this will return `None`.",
						" Use [`Self::metadata_versions`] to find out about supported metadata version of the runtime."},
				},
				primitives.RuntimeApiMethodMetadata{
					Name:   "metadata_versions",
					Inputs: sc.Sequence[primitives.RuntimeApiMethodParamMetadata]{},
					Output: sc.ToCompact(metadata.TypesSequenceU32),
					Docs: sc.Sequence[sc.Str]{" Returns the supported metadata versions.",
						"",
						" This can be used to call `metadata_at_version`."},
				},
			},
			Docs: sc.Sequence[sc.Str]{" The `Metadata` api trait that returns metadata for the runtime."},
		},
	}

	expectedSupportVersions = sc.Sequence[sc.U32]{
		sc.U32(primitives.MetadataVersion14), sc.U32(primitives.MetadataVersion15),
	}
)

var (
	mockRuntimeExtrinsic *mocks.RuntimeExtrinsic
	mockMemoryUtils      *mocks.MemoryTranslator
)

func Test_Module_Name(t *testing.T) {
	target := setup()

	result := target.Name()

	assert.Equal(t, ApiModuleName, result)
}

func Test_Module_Item(t *testing.T) {
	target := setup()

	hexName := common.MustBlake2b8([]byte(ApiModuleName))
	expect := primitives.NewApiItem(hexName, apiVersion)

	result := target.Item()

	assert.Equal(t, expect, result)
}

func Test_Module_Metadata(t *testing.T) {
	target := setup()

	constantsMap := buildConstantsMap()

	mockRuntimeExtrinsic.On("Metadata", constantsMap).Return(mdTypes, mdModules14, mdExtrinsic)

	builtMeta := target.buildMetadata()

	mdV14 := primitives.NewMetadataV14(builtMeta.DataV14)

	bMetadata := sc.BytesToSequenceU8(mdV14.Bytes())

	mockMemoryUtils.On("BytesToOffsetAndSize", bMetadata.Bytes()).Return(ptrAndSize)

	result := target.Metadata()

	assert.Equal(t, ptrAndSize, result)

	mockRuntimeExtrinsic.AssertCalled(t, "Metadata", constantsMap)
	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", bMetadata.Bytes())
}

func Test_Module_Metadata_Versions(t *testing.T) {
	target := setup()

	mockMemoryUtils.On("BytesToOffsetAndSize", expectedSupportVersions.Bytes()).Return(ptrAndSize)

	resultVersions := target.MetadataVersions()

	assert.Equal(t, ptrAndSize, resultVersions)

	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", expectedSupportVersions.Bytes())
}

func Test_Module_Metadata_AtVersion_14(t *testing.T) {
	target := setup()

	metadataTypes := getAllMetadataTypes(&target)

	version14 := sc.U32(primitives.MetadataVersion14)

	constantsMap := buildConstantsMap()

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(version14.Bytes())

	mockRuntimeExtrinsic.On("Metadata", constantsMap).Return(mdTypes, mdModules14, mdExtrinsic)

	metadataV14 := primitives.RuntimeMetadataV14{
		Types:     metadataTypes,
		Modules:   mdModules14,
		Extrinsic: mdExtrinsic,
		Type:      sc.ToCompact(metadata.Runtime),
	}

	bMetadataV14 := sc.BytesToSequenceU8(primitives.NewMetadataV14(metadataV14).Bytes())
	optionMd14 := sc.Option[sc.Sequence[sc.U8]]{
		HasValue: sc.Bool(true),
		Value:    bMetadataV14,
	}

	mockMemoryUtils.On("BytesToOffsetAndSize", optionMd14.Bytes()).Return(ptrAndSize)

	resultVersion14 := target.MetadataAtVersion(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, resultVersion14)

	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", optionMd14.Bytes())
}

func Test_Module_Metadata_AtVersion_15(t *testing.T) {
	target := setup()

	constantsMap := buildConstantsMap()

	metadataTypes := getAllMetadataTypes(&target)

	version15 := sc.U32(primitives.MetadataVersion15)

	outerEnums := primitives.OuterEnums{
		CallEnumType:  sc.ToCompact(metadata.RuntimeCall),
		EventEnumType: sc.ToCompact(metadata.TypesRuntimeEvent),
		ErrorEnumType: sc.ToCompact(metadata.TypesRuntimeError),
	}

	custom := primitives.CustomMetadata{
		Map: sc.Dictionary[sc.Str, primitives.CustomValueMetadata]{},
	}

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(version15.Bytes())

	mockRuntimeExtrinsic.On("MetadataLatest", constantsMap).Return(mdTypes, mdModules15, mdExtrinsic15, outerEnums, custom)

	metadataV15 := primitives.RuntimeMetadataV15{
		Types:      metadataTypes,
		Modules:    mdModules15,
		Extrinsic:  mdExtrinsic15,
		Type:       sc.ToCompact(metadata.Runtime),
		Apis:       mdRuntimeApi,
		OuterEnums: outerEnums,
		Custom:     custom,
	}

	bMetadataV15 := sc.BytesToSequenceU8(primitives.NewMetadataV15(metadataV15).Bytes())
	optionMd15 := sc.Option[sc.Sequence[sc.U8]]{
		HasValue: sc.Bool(true),
		Value:    bMetadataV15,
	}

	mockMemoryUtils.On("BytesToOffsetAndSize", optionMd15.Bytes()).Return(ptrAndSize)

	resultVersion15 := target.MetadataAtVersion(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, resultVersion15)

	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", optionMd15.Bytes())
}

func Test_Module_Metadata_AtVersion_Unsupported(t *testing.T) {
	target := setup()

	version10 := sc.U32(10)

	mockMemoryUtils.On("GetWasmMemorySlice", dataPtr, dataLen).Return(version10.Bytes())

	optionUnsupported := sc.Option[sc.Sequence[sc.U8]]{
		HasValue: sc.Bool(false),
	}

	mockMemoryUtils.On("BytesToOffsetAndSize", optionUnsupported.Bytes()).Return(ptrAndSize)

	resultUnsupported := target.MetadataAtVersion(dataPtr, dataLen)

	assert.Equal(t, ptrAndSize, resultUnsupported)

	mockMemoryUtils.AssertCalled(t, "BytesToOffsetAndSize", optionUnsupported.Bytes())
}

func setup() Module {
	mockRuntimeExtrinsic = new(mocks.RuntimeExtrinsic)
	mockMemoryUtils = new(mocks.MemoryTranslator)

	target := New(mockRuntimeExtrinsic, []primitives.RuntimeApiModule{})
	target.memUtils = mockMemoryUtils

	return target
}

func getAllMetadataTypes(target *Module) sc.Sequence[primitives.MetadataType] {
	metadataTypes := append(primitiveTypes(), basicTypes()...)

	metadataTypes = append(metadataTypes, target.runtimeTypes()...)

	metadataTypes = append(metadataTypes, mdTypes...)

	return metadataTypes
}
