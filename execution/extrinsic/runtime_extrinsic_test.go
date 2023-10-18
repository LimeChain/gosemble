package extrinsic

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/execution/types"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	blockNumber                 = sc.U64(5)
	bytesExtrinsicFormatVersion = sc.U8(types.ExtrinsicFormatVersion).Bytes()
	inherentIdentifier          = [8]byte{'t', 'e', 's', 't', 't', 'i', 'n', 'g'}

	moduleOneId = sc.U8(2)
	moduleTwoId = sc.U8(3)

	weightOne = primitives.WeightFromParts(1, 2)
	weightTwo = primitives.WeightFromParts(3, 4)
)

var (
	mockModuleOne *mocks.Module
	mockModuleTwo *mocks.Module

	mockCallOne *mocks.Call
	mockCallTwo *mocks.Call

	mockBlock              *mocks.Block
	mockUncheckedExtrinsic *mocks.UncheckedExtrinsic
)

var (
	metadataTypes = sc.Sequence[primitives.MetadataType]{}
	metadataOne   = primitives.MetadataModule{
		ModuleV14: primitives.MetadataModuleV14{
			Name:      "moduleOne",
			Storage:   sc.Option[primitives.MetadataModuleStorage]{},
			Call:      sc.Option[sc.Compact]{},
			CallDef:   sc.NewOption[primitives.MetadataDefinitionVariant](nil),
			Event:     sc.Option[sc.Compact]{},
			EventDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
			Constants: nil,
			Error:     sc.Option[sc.Compact]{},
			ErrorDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
			Index:     moduleOneId,
		},
		ModuleV15: primitives.MetadataModuleV15{
			Name:      "moduleOne",
			Storage:   sc.Option[primitives.MetadataModuleStorage]{},
			Call:      sc.Option[sc.Compact]{},
			CallDef:   sc.Option[primitives.MetadataDefinitionVariant]{},
			Event:     sc.Option[sc.Compact]{},
			EventDef:  sc.Option[primitives.MetadataDefinitionVariant]{},
			Constants: nil,
			Error:     sc.Option[sc.Compact]{},
			ErrorDef:  sc.Option[primitives.MetadataDefinitionVariant]{},
			Index:     moduleOneId,
			Docs:      nil,
		},
	}
	metadataTwo = primitives.MetadataModule{
		ModuleV14: primitives.MetadataModuleV14{
			Name:    "moduleTwo",
			Storage: sc.Option[primitives.MetadataModuleStorage]{},
			Call:    sc.Option[sc.Compact]{},
			CallDef: sc.NewOption[primitives.MetadataDefinitionVariant](
				primitives.NewMetadataDefinitionVariantStr(
					"moduleTwo",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(7, "self::ModuleTwo"),
					},
					moduleTwoId,
					"Call.ModuleTwo")),
			Event:     sc.Option[sc.Compact]{},
			EventDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
			Constants: nil,
			Error:     sc.Option[sc.Compact]{},
			ErrorDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
			Index:     moduleTwoId,
		},
		ModuleV15: primitives.MetadataModuleV15{
			Name:    "moduleTwo",
			Storage: sc.Option[primitives.MetadataModuleStorage]{},
			Call:    sc.Option[sc.Compact]{},
			CallDef: sc.Option[primitives.MetadataDefinitionVariant]{},
			Event:   sc.Option[sc.Compact]{},
			EventDef: sc.NewOption[primitives.MetadataDefinitionVariant](
				primitives.NewMetadataDefinitionVariantStr(
					"moduleTwo",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithName(7, "self::ModuleTwo"),
					},
					moduleTwoId,
					"EventDef.ModuleTwo")),
			Constants: nil,
			Error:     sc.Option[sc.Compact]{},
			ErrorDef:  sc.Option[primitives.MetadataDefinitionVariant]{},
			Index:     moduleTwoId,
			Docs:      nil,
		},
	}
	signedExtensions = sc.Sequence[primitives.MetadataSignedExtension]{}
)

func Test_RuntimeExtrinsic_Module(t *testing.T) {
	target := setupRuntimeExtrinsic()

	mockModuleOne.On("GetIndex").Return(moduleOneId)

	result, ok := target.Module(moduleOneId)

	assert.Equal(t, mockModuleOne, result)
	assert.Equal(t, true, ok)
	mockModuleOne.AssertCalled(t, "GetIndex")
}

func Test_RuntimeExtrinsic_Module_NotFound(t *testing.T) {
	target := setupRuntimeExtrinsic()

	invalidModuleId := sc.U8(6)

	mockModuleOne.On("GetIndex").Return(moduleOneId)
	mockModuleTwo.On("GetIndex").Return(moduleTwoId)

	result, ok := target.Module(invalidModuleId)

	assert.Equal(t, nil, result)
	assert.Equal(t, false, ok)
	mockModuleOne.AssertCalled(t, "GetIndex")
	mockModuleTwo.AssertCalled(t, "GetIndex")
}

func Test_RuntimeExtrinsic_CreateInherents(t *testing.T) {
	target := setupRuntimeExtrinsic()

	inherentData := *primitives.NewInherentData()
	buffer := bytes.NewBuffer(bytesExtrinsicFormatVersion)
	bytesExtrinsic := append(sc.ToCompact(buffer.Len()).Bytes(), buffer.Bytes()...)
	expect := append(sc.ToCompact(2).Bytes(), append(bytesExtrinsic, bytesExtrinsic...)...) // two modules, including each extrinsic's bytes

	mockModuleOne.On("CreateInherent", inherentData).Return(sc.NewOption[primitives.Call](mockCallOne))
	mockCallOne.On("Encode", buffer).Return()
	mockModuleTwo.On("CreateInherent", inherentData).Return(sc.NewOption[primitives.Call](mockCallTwo))
	mockCallTwo.On("Encode", buffer).Return()

	result := target.CreateInherents(inherentData)

	assert.Equal(t, expect, result)

	mockModuleOne.AssertCalled(t, "CreateInherent", inherentData)
	mockCallOne.AssertCalled(t, "Encode", buffer)
	mockModuleTwo.AssertCalled(t, "CreateInherent", inherentData)
	mockCallTwo.AssertCalled(t, "Encode", buffer)
}

func Test_RuntimeExtrinsic_CreateInherents_Empty(t *testing.T) {
	target := setupRuntimeExtrinsic()

	inherentData := *primitives.NewInherentData()
	emptyCall := sc.NewOption[primitives.Call](nil)

	mockModuleOne.On("CreateInherent", inherentData).Return(emptyCall)
	mockModuleTwo.On("CreateInherent", inherentData).Return(emptyCall)

	result := target.CreateInherents(inherentData)

	assert.Equal(t, []byte{}, result)

	mockModuleOne.AssertCalled(t, "CreateInherent", inherentData)
	mockCallOne.AssertNotCalled(t, "Encode", mock.Anything)
	mockModuleTwo.AssertCalled(t, "CreateInherent", inherentData)
	mockCallTwo.AssertNotCalled(t, "Encode", mock.Anything)
}

func Test_RuntimeExtrinsic_CheckInherents(t *testing.T) {
	target := setupRuntimeExtrinsic()

	inherentData := *primitives.NewInherentData()
	expect := primitives.NewCheckInherentsResult()

	mockBlock.On("Extrinsics").Return(sc.Sequence[primitives.UncheckedExtrinsic]{mockUncheckedExtrinsic})
	mockUncheckedExtrinsic.On("IsSigned").Return(false)
	mockUncheckedExtrinsic.On("Function").Return(mockCallOne)
	mockModuleOne.On("IsInherent", mockCallOne).Return(true)
	mockModuleOne.On("CheckInherent", mockCallOne, inherentData).Return(nil)
	mockModuleTwo.On("IsInherent", mockCallOne).Return(false)

	result := target.CheckInherents(inherentData, mockBlock)

	assert.Equal(t, expect, result)
	mockBlock.AssertCalled(t, "Extrinsics")
	mockUncheckedExtrinsic.AssertCalled(t, "IsSigned")
	mockUncheckedExtrinsic.AssertCalled(t, "Function")
	mockModuleOne.AssertCalled(t, "IsInherent", mockCallOne)
	mockModuleOne.AssertCalled(t, "CheckInherent", mockCallOne, inherentData)
	mockModuleTwo.AssertCalled(t, "IsInherent", mockCallOne)
}

func Test_RuntimeExtrinsic_CheckInherents_FatalError(t *testing.T) {
	target := setupRuntimeExtrinsic()

	inherentData := *primitives.NewInherentData()
	err := primitives.NewTimestampErrorInvalid()
	inherentData.Put(inherentIdentifier, err)
	expect := primitives.CheckInherentsResult{
		Okay:       false,
		FatalError: true,
		Errors:     inherentData,
	}

	mockBlock.On("Extrinsics").Return(sc.Sequence[primitives.UncheckedExtrinsic]{mockUncheckedExtrinsic})
	mockUncheckedExtrinsic.On("IsSigned").Return(false)
	mockUncheckedExtrinsic.On("Function").Return(mockCallOne)
	mockModuleOne.On("IsInherent", mockCallOne).Return(true)
	mockModuleOne.On("CheckInherent", mockCallOne, inherentData).Return(err)
	mockModuleOne.On("InherentIdentifier").Return(inherentIdentifier)

	result := target.CheckInherents(inherentData, mockBlock)

	assert.Equal(t, expect, result)
	mockBlock.AssertCalled(t, "Extrinsics")
	mockUncheckedExtrinsic.AssertCalled(t, "IsSigned")
	mockUncheckedExtrinsic.AssertCalled(t, "Function")
	mockModuleOne.AssertCalled(t, "IsInherent", mockCallOne)
	mockModuleOne.AssertCalled(t, "CheckInherent", mockCallOne, inherentData)
	mockModuleOne.AssertCalled(t, "InherentIdentifier")
}

func Test_RuntimeExtrinsic_CheckInherents_NoInherents(t *testing.T) {
	target := setupRuntimeExtrinsic()

	inherentData := *primitives.NewInherentData()
	expect := primitives.NewCheckInherentsResult()

	mockBlock.On("Extrinsics").Return(sc.Sequence[primitives.UncheckedExtrinsic]{mockUncheckedExtrinsic})
	mockUncheckedExtrinsic.On("IsSigned").Return(false)
	mockUncheckedExtrinsic.On("Function").Return(mockCallOne)
	mockModuleOne.On("IsInherent", mockCallOne).Return(false)
	mockModuleTwo.On("IsInherent", mockCallOne).Return(false)

	result := target.CheckInherents(inherentData, mockBlock)

	assert.Equal(t, expect, result)
	mockBlock.AssertCalled(t, "Extrinsics")
	mockUncheckedExtrinsic.AssertCalled(t, "IsSigned")
	mockUncheckedExtrinsic.AssertCalled(t, "Function")
	mockModuleOne.AssertCalled(t, "IsInherent", mockCallOne)
	mockModuleTwo.AssertCalled(t, "IsInherent", mockCallOne)
}

func Test_RuntimeExtrinsic_CheckInherents_Signed(t *testing.T) {
	target := setupRuntimeExtrinsic()

	inherentData := *primitives.NewInherentData()
	expect := primitives.NewCheckInherentsResult()

	mockBlock.On("Extrinsics").Return(sc.Sequence[primitives.UncheckedExtrinsic]{mockUncheckedExtrinsic})
	mockUncheckedExtrinsic.On("IsSigned").Return(true)

	result := target.CheckInherents(inherentData, mockBlock)

	assert.Equal(t, expect, result)

	mockBlock.AssertCalled(t, "Extrinsics")
	mockUncheckedExtrinsic.AssertCalled(t, "IsSigned")
}

func Test_RuntimeExtrinsic_EnsureInherentsAreFirst_Signed(t *testing.T) {
	target := setupRuntimeExtrinsic()

	mockBlock.On("Extrinsics").Return(sc.Sequence[primitives.UncheckedExtrinsic]{mockUncheckedExtrinsic})
	mockUncheckedExtrinsic.On("IsSigned").Return(true)

	result := target.EnsureInherentsAreFirst(mockBlock)

	assert.Equal(t, -1, result)
	mockBlock.AssertCalled(t, "Extrinsics")
	mockUncheckedExtrinsic.AssertCalled(t, "IsSigned")
}

func Test_RuntimeExtrinsic_EnsureInherentsAreFirst_Unsigned(t *testing.T) {
	target := setupRuntimeExtrinsic()

	mockBlock.On("Extrinsics").Return(sc.Sequence[primitives.UncheckedExtrinsic]{mockUncheckedExtrinsic})
	mockUncheckedExtrinsic.On("IsSigned").Return(false)
	mockUncheckedExtrinsic.On("Function").Return(mockCallOne)
	mockModuleOne.On("IsInherent", mockCallOne).Return(true)
	mockModuleTwo.On("IsInherent", mockCallOne).Return(false)

	result := target.EnsureInherentsAreFirst(mockBlock)

	assert.Equal(t, -1, result)
	mockBlock.AssertCalled(t, "Extrinsics")
	mockUncheckedExtrinsic.AssertCalled(t, "IsSigned")
	mockUncheckedExtrinsic.AssertCalled(t, "Function")
	mockModuleOne.AssertCalled(t, "IsInherent", mockCallOne)
	mockModuleTwo.AssertCalled(t, "IsInherent", mockCallOne)
}

func Test_RuntimeExtrinsic_EnsureInherentsAreFirst_SignedBeforeUnsigned(t *testing.T) {
	target := setupRuntimeExtrinsic()

	mockSignedUncheckedExtrinsic := new(mocks.UncheckedExtrinsic)

	mockBlock.On("Extrinsics").
		Return(sc.Sequence[primitives.UncheckedExtrinsic]{
			mockSignedUncheckedExtrinsic,
			mockUncheckedExtrinsic,
		})
	mockSignedUncheckedExtrinsic.On("IsSigned").Return(true)
	mockUncheckedExtrinsic.On("IsSigned").Return(false)
	mockUncheckedExtrinsic.On("Function").Return(mockCallOne)
	mockModuleOne.On("IsInherent", mockCallOne).Return(true)
	mockModuleTwo.On("IsInherent", mockCallOne).Return(false)

	result := target.EnsureInherentsAreFirst(mockBlock)

	assert.Equal(t, 1, result)
	mockSignedUncheckedExtrinsic.AssertCalled(t, "IsSigned")
	mockUncheckedExtrinsic.AssertCalled(t, "IsSigned")
	mockUncheckedExtrinsic.AssertCalled(t, "Function")
	mockModuleOne.AssertCalled(t, "IsInherent", mockCallOne)
	mockModuleTwo.AssertCalled(t, "IsInherent", mockCallOne)
}

func Test_RuntimeExtrinsic_OnInitialize(t *testing.T) {
	target := setupRuntimeExtrinsic()

	expect := weightOne.Add(weightTwo)

	mockModuleOne.On("OnInitialize", blockNumber).Return(weightOne)
	mockModuleTwo.On("OnInitialize", blockNumber).Return(weightTwo)

	result := target.OnInitialize(blockNumber)

	assert.Equal(t, expect, result)
	mockModuleOne.AssertCalled(t, "OnInitialize", blockNumber)
	mockModuleTwo.AssertCalled(t, "OnInitialize", blockNumber)
}

func Test_RuntimeExtrinsic_OnRuntimeUpgrade(t *testing.T) {
	target := setupRuntimeExtrinsic()

	expect := weightOne.Add(weightTwo)

	mockModuleOne.On("OnRuntimeUpgrade").Return(weightOne)
	mockModuleTwo.On("OnRuntimeUpgrade").Return(weightTwo)

	result := target.OnRuntimeUpgrade()

	assert.Equal(t, expect, result)
	mockModuleOne.AssertCalled(t, "OnRuntimeUpgrade")
	mockModuleTwo.AssertCalled(t, "OnRuntimeUpgrade")
}

func Test_RuntimeExtrinsic_OnFinalize(t *testing.T) {
	target := setupRuntimeExtrinsic()

	mockModuleOne.On("OnFinalize", blockNumber).Return()
	mockModuleTwo.On("OnFinalize", blockNumber).Return()

	target.OnFinalize(blockNumber)

	mockModuleOne.AssertCalled(t, "OnFinalize", blockNumber)
	mockModuleTwo.AssertCalled(t, "OnFinalize", blockNumber)
}

func Test_RuntimeExtrinsic_OnIdle(t *testing.T) {
	target := setupRuntimeExtrinsic()

	remainingWeight := primitives.WeightFromParts(10, 10)
	secondAdjustedWeight := remainingWeight.Sub(weightOne)

	mockModuleOne.On("OnIdle", blockNumber, remainingWeight).Return(weightOne)
	mockModuleTwo.On("OnIdle", blockNumber, secondAdjustedWeight).Return(weightTwo)

	result := target.OnIdle(blockNumber, remainingWeight)

	assert.Equal(t, weightOne.Add(weightTwo), result)
	mockModuleOne.AssertCalled(t, "OnIdle", blockNumber, remainingWeight)
	mockModuleTwo.AssertCalled(t, "OnIdle", blockNumber, secondAdjustedWeight)
}

func Test_RuntimeExtrinsic_OffchainWorker(t *testing.T) {
	target := setupRuntimeExtrinsic()

	mockModuleOne.On("OffchainWorker", blockNumber).Return()
	mockModuleTwo.On("OffchainWorker", blockNumber).Return()

	target.OffchainWorker(blockNumber)

	mockModuleOne.AssertCalled(t, "OffchainWorker", blockNumber)
	mockModuleTwo.AssertCalled(t, "OffchainWorker", blockNumber)
}

func Test_RuntimeExtrinsic_Metadata(t *testing.T) {
	target := setupRuntimeExtrinsic()

	expectTypes := sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithPath(
			metadata.TypesRuntimeEvent,
			"node_template_runtime RuntimeEvent",
			sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeEvent"},
			primitives.NewMetadataTypeDefinitionVariant(sc.Sequence[primitives.MetadataDefinitionVariant]{}),
		),
		primitives.NewMetadataTypeWithPath(
			metadata.RuntimeCall,
			"RuntimeCall",
			sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeCall"},
			primitives.NewMetadataTypeDefinitionVariant(sc.Sequence[primitives.MetadataDefinitionVariant]{
				metadataTwo.ModuleV14.CallDef.Value,
			}),
		),
		primitives.NewMetadataTypeWithPath(
			metadata.TypesRuntimeError,
			"node_template_runtime RuntimeError",
			sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeError"},
			primitives.NewMetadataTypeDefinitionVariant(sc.Sequence[primitives.MetadataDefinitionVariant]{}),
		),
		primitives.NewMetadataTypeWithParams(metadata.UncheckedExtrinsic, "UncheckedExtrinsic",
			sc.Sequence[sc.Str]{"sp_runtime", "generic", "unchecked_extrinsic", "UncheckedExtrinsic"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceU8),
				}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesMultiAddress, "Address"),
				primitives.NewMetadataTypeParameter(metadata.RuntimeCall, "Call"),
				primitives.NewMetadataTypeParameter(metadata.TypesMultiSignature, "Signature"),
				primitives.NewMetadataTypeParameter(metadata.SignedExtra, "Extra"),
			},
		),
	}
	expectExtrinsic := primitives.MetadataExtrinsicV14{
		Type:             sc.ToCompact(metadata.UncheckedExtrinsic),
		Version:          types.ExtrinsicFormatVersion,
		SignedExtensions: signedExtensions,
	}
	expectModules := sc.Sequence[primitives.MetadataModuleV14]{
		metadataOne.ModuleV14,
		metadataTwo.ModuleV14,
	}

	mockModuleOne.On("Metadata").Return(metadataTypes, metadataOne)
	mockModuleTwo.On("Metadata").Return(metadataTypes, metadataTwo)
	mockSignedExtra.On("Metadata").Return(metadataTypes, signedExtensions)

	resultTypes, resultModules, resultExtrinsic := target.Metadata()

	assert.Equal(t, expectTypes, resultTypes)
	assert.Equal(t, expectModules, resultModules)
	assert.Equal(t, expectExtrinsic, resultExtrinsic)
	mockModuleOne.AssertCalled(t, "Metadata")
	mockModuleTwo.AssertCalled(t, "Metadata")
	mockSignedExtra.AssertCalled(t, "Metadata")
}

func Test_RuntimeExtrinsic_MetadataLatest(t *testing.T) {
	target := setupRuntimeExtrinsic()

	expectTypes := sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithPath(
			metadata.TypesRuntimeEvent,
			"node_template_runtime RuntimeEvent",
			sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeEvent"},
			primitives.NewMetadataTypeDefinitionVariant(sc.Sequence[primitives.MetadataDefinitionVariant]{
				metadataTwo.ModuleV15.EventDef.Value,
			}),
		),
		primitives.NewMetadataTypeWithPath(
			metadata.RuntimeCall,
			"RuntimeCall",
			sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeCall"},
			primitives.NewMetadataTypeDefinitionVariant(sc.Sequence[primitives.MetadataDefinitionVariant]{}),
		),
		primitives.NewMetadataTypeWithPath(
			metadata.TypesRuntimeError,
			"node_template_runtime RuntimeError",
			sc.Sequence[sc.Str]{"node_template_runtime", "RuntimeError"},
			primitives.NewMetadataTypeDefinitionVariant(sc.Sequence[primitives.MetadataDefinitionVariant]{}),
		),
		primitives.NewMetadataTypeWithParams(metadata.UncheckedExtrinsic, "UncheckedExtrinsic",
			sc.Sequence[sc.Str]{"sp_runtime", "generic", "unchecked_extrinsic", "UncheckedExtrinsic"},
			primitives.NewMetadataTypeDefinitionComposite(
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceU8),
				}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataTypeParameter(metadata.TypesMultiAddress, "Address"),
				primitives.NewMetadataTypeParameter(metadata.RuntimeCall, "Call"),
				primitives.NewMetadataTypeParameter(metadata.TypesMultiSignature, "Signature"),
				primitives.NewMetadataTypeParameter(metadata.SignedExtra, "Extra"),
			},
		),
	}
	expectModules := sc.Sequence[primitives.MetadataModuleV15]{
		metadataOne.ModuleV15,
		metadataTwo.ModuleV15,
	}
	expectOuterEnums := primitives.OuterEnums{
		CallEnumType:  sc.ToCompact(metadata.RuntimeCall),
		EventEnumType: sc.ToCompact(metadata.TypesRuntimeEvent),
		ErrorEnumType: sc.ToCompact(metadata.TypesRuntimeError),
	}
	expectCustom := primitives.CustomMetadata{
		Map: sc.Dictionary[sc.Str, primitives.CustomValueMetadata]{},
	}
	expectExtrinsic := primitives.MetadataExtrinsicV15{
		Version:          types.ExtrinsicFormatVersion,
		Address:          sc.ToCompact(metadata.TypesMultiAddress),
		Call:             sc.ToCompact(metadata.RuntimeCall),
		Signature:        sc.ToCompact(metadata.TypesMultiSignature),
		Extra:            sc.ToCompact(metadata.SignedExtra),
		SignedExtensions: signedExtensions,
	}

	mockModuleOne.On("Metadata").Return(metadataTypes, metadataOne)
	mockModuleTwo.On("Metadata").Return(metadataTypes, metadataTwo)
	mockSignedExtra.On("Metadata").Return(metadataTypes, signedExtensions)

	resultTypes, resultModules, resultExtrinsic, resultOuterEnums, resultCustom := target.MetadataLatest()

	assert.Equal(t, expectTypes, resultTypes)
	assert.Equal(t, expectModules, resultModules)
	assert.Equal(t, expectExtrinsic, resultExtrinsic)
	assert.Equal(t, expectOuterEnums, resultOuterEnums)
	assert.Equal(t, expectCustom, resultCustom)
	mockModuleOne.AssertCalled(t, "Metadata")
	mockModuleTwo.AssertCalled(t, "Metadata")
	mockSignedExtra.AssertCalled(t, "Metadata")
}

func setupRuntimeExtrinsic() RuntimeExtrinsic {
	mockSignedExtra = new(mocks.SignedExtra)

	mockModuleOne = new(mocks.Module)
	mockModuleTwo = new(mocks.Module)
	mockCallOne = new(mocks.Call)
	mockCallTwo = new(mocks.Call)
	mockBlock = new(mocks.Block)
	mockUncheckedExtrinsic = new(mocks.UncheckedExtrinsic)

	modules := []primitives.Module{
		mockModuleOne,
		mockModuleTwo,
	}

	return New(modules, mockSignedExtra)
}
