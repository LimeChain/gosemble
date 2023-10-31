package aura

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	moduleId                          = 13
	weightRefTimePerNanos      sc.U64 = 1_000
	timestampMinimumPeriod            = 2_000
	maxAuthorites                     = 10
	allowMultipleBlocksPerSlot        = false
	keyType                           = types.PublicKeySr25519
	blockNumber                sc.U64 = 0
)

var (
	dbWeight = types.RuntimeDbWeight{
		Read:  3 * weightRefTimePerNanos,
		Write: 7 * weightRefTimePerNanos,
	}
	module                 Module
	mockStorageDigest      *mocks.StorageValue[types.Digest]
	mockStorageCurrentSlot *mocks.StorageValue[sc.U64]
	mockStorageAuthorities *mocks.StorageValue[sc.Sequence[sc.U8]]
)

var (
	expectedMetadataTypes = sc.Sequence[types.MetadataType]{
		types.NewMetadataTypeWithParams(
			metadata.TypesAuraStorageAuthorities,
			"BoundedVec<T::AuthorityId, T::MaxAuthorities>",
			sc.Sequence[sc.Str]{"bounded_collection", "bounded_vec", "BoundedVec"},
			types.NewMetadataTypeDefinitionComposite(
				sc.Sequence[types.MetadataTypeDefinitionField]{
					types.NewMetadataTypeDefinitionField(metadata.TypesSequencePubKeys),
				}), sc.Sequence[types.MetadataTypeParameter]{
				types.NewMetadataTypeParameter(metadata.TypesAuthorityId, "T"),
				types.NewMetadataEmptyTypeParameter("S"),
			}),

		types.NewMetadataTypeWithPath(metadata.TypesAuthorityId,
			"sp_consensus_aura sr25519 app_sr25519 Public",
			sc.Sequence[sc.Str]{"sp_consensus_aura", "sr25519", "app_sr25519", "Public"},
			types.NewMetadataTypeDefinitionComposite(
				sc.Sequence[types.MetadataTypeDefinitionField]{types.NewMetadataTypeDefinitionField(metadata.TypesSr25519PubKey)})),

		types.NewMetadataTypeWithPath(metadata.TypesSr25519PubKey,
			"sp_core sr25519 Public",
			sc.Sequence[sc.Str]{"sp_core", "sr25519", "Public"},
			types.NewMetadataTypeDefinitionComposite(
				sc.Sequence[types.MetadataTypeDefinitionField]{types.NewMetadataTypeDefinitionField(metadata.TypesFixedSequence32U8)})),

		types.NewMetadataType(metadata.TypesSequencePubKeys,
			"[]PublicKey",
			types.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesAuthorityId))),

		types.NewMetadataTypeWithPath(metadata.TypesAuraSlot,
			"sp_consensus_slots Slot",
			sc.Sequence[sc.Str]{"sp_consensus_slots", "Slot"},
			types.NewMetadataTypeDefinitionComposite(
				sc.Sequence[types.MetadataTypeDefinitionField]{
					types.NewMetadataTypeDefinitionField(metadata.PrimitiveTypesU64),
				})),

		// type 924
		types.NewMetadataType(metadata.TypesTupleSequenceU8KeyTypeId, "(Seq<U8>, KeyTypeId)",
			types.NewMetadataTypeDefinitionTuple(sc.Sequence[sc.Compact]{sc.ToCompact(metadata.TypesSequenceU8), sc.ToCompact(metadata.TypesKeyTypeId)})),

		// type 923
		types.NewMetadataType(metadata.TypesSequenceTupleSequenceU8KeyTypeId, "[]byte TupleSequenceU8KeyTypeId", types.NewMetadataTypeDefinitionSequence(sc.ToCompact(metadata.TypesTupleSequenceU8KeyTypeId))),

		// type 922
		types.NewMetadataTypeWithParam(metadata.TypesOptionTupleSequenceU8KeyTypeId, "Option<TupleSequenceU8KeyTypeId>", sc.Sequence[sc.Str]{"Option"}, types.NewMetadataTypeDefinitionVariant(
			sc.Sequence[types.MetadataDefinitionVariant]{
				types.NewMetadataDefinitionVariant(
					"None",
					sc.Sequence[types.MetadataTypeDefinitionField]{},
					0,
					""),
				types.NewMetadataDefinitionVariant(
					"Some",
					sc.Sequence[types.MetadataTypeDefinitionField]{
						types.NewMetadataTypeDefinitionField(metadata.TypesSequenceTupleSequenceU8KeyTypeId),
					},
					1,
					""),
			}),
			types.NewMetadataTypeParameter(metadata.TypesSequenceTupleSequenceU8KeyTypeId, "T")),
	}

	moduleV14 = types.MetadataModuleV14{
		Name: "Aura",
		Storage: sc.NewOption[types.MetadataModuleStorage](types.MetadataModuleStorage{
			Prefix: "Aura",
			Items: sc.Sequence[types.MetadataModuleStorageEntry]{
				types.NewMetadataModuleStorageEntry(
					"Authorities",
					types.MetadataModuleStorageEntryModifierDefault,
					types.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.TypesAuraStorageAuthorities)),
					"The current authority set."),
				types.NewMetadataModuleStorageEntry(
					"CurrentSlot",
					types.MetadataModuleStorageEntryModifierDefault,
					types.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.TypesAuraSlot)),
					"The current slot of this block.   This will be set in `on_initialize`."),
			},
		}),
		Call:      sc.NewOption[sc.Compact](nil),
		CallDef:   sc.NewOption[types.MetadataDefinitionVariant](nil),
		Event:     sc.NewOption[sc.Compact](nil),
		EventDef:  sc.NewOption[types.MetadataDefinitionVariant](nil),
		Constants: sc.Sequence[types.MetadataModuleConstant]{},
		Error:     sc.NewOption[sc.Compact](nil),
		ErrorDef:  sc.NewOption[types.MetadataDefinitionVariant](nil),
		Index:     sc.U8(13),
	}

	expectedMetadataModule = types.MetadataModule{
		Version:   types.ModuleVersion14,
		ModuleV14: moduleV14,
	}
)

func setup(minimumPeriod sc.U64) {
	mockStorageDigest = new(mocks.StorageValue[types.Digest])
	mockStorageCurrentSlot = new(mocks.StorageValue[sc.U64])
	mockStorageAuthorities = new(mocks.StorageValue[sc.Sequence[sc.U8]])

	config := NewConfig(
		keyType,
		dbWeight,
		minimumPeriod,
		maxAuthorites,
		allowMultipleBlocksPerSlot,
		mockStorageDigest.Get,
	)
	module = New(moduleId, config)
	module.storage.CurrentSlot = mockStorageCurrentSlot
	module.storage.Authorities = mockStorageAuthorities
}

func newPreRuntimeDigest(n sc.U64) *types.Digest {
	digest := types.Digest{}
	preRuntimeDigestItem := types.DigestItem{
		Engine:  sc.BytesToFixedSequenceU8(EngineId[:]),
		Payload: sc.BytesToSequenceU8(n.Bytes()),
	}
	digest[types.DigestTypePreRuntime] = append(digest[types.DigestTypePreRuntime], preRuntimeDigestItem)
	return &digest
}

func Test_Aura_GetIndex(t *testing.T) {
	setup(timestampMinimumPeriod)

	assert.Equal(t, sc.U8(13), module.GetIndex())
}

func Test_Aura_Functions(t *testing.T) {
	setup(timestampMinimumPeriod)

	assert.Equal(t, map[sc.U8]types.Call{}, module.Functions())
}

func Test_Module_PreDispatch(t *testing.T) {
	setup(timestampMinimumPeriod)

	result, err := module.PreDispatch(new(mocks.Call))

	assert.Nil(t, err)
	assert.Equal(t, sc.Empty{}, result)
}

func Test_Module_ValidateUnsigned(t *testing.T) {
	setup(timestampMinimumPeriod)

	result, err := module.ValidateUnsigned(types.NewTransactionSourceLocal(), new(mocks.Call))

	assert.Equal(t, types.NewTransactionValidityError(types.NewUnknownTransactionNoUnsignedValidator()), err)
	assert.Equal(t, types.ValidTransaction{}, result)
}

func Test_Aura_KeyType(t *testing.T) {
	setup(timestampMinimumPeriod)

	assert.Equal(t, keyType, module.KeyType())
}

func Test_Aura_KeyTypeId(t *testing.T) {
	setup(timestampMinimumPeriod)

	assert.Equal(t, [4]byte{'a', 'u', 'r', 'a'}, module.KeyTypeId())
}

func Test_Aura_Metadata(t *testing.T) {
	setup(timestampMinimumPeriod)

	metadataTypes, metadataModule := module.Metadata()
	assert.Equal(t, expectedMetadataTypes, metadataTypes)
	assert.Equal(t, expectedMetadataModule, metadataModule)
}

func Test_Aura_OnInitialize_EmptySlot(t *testing.T) {
	setup(timestampMinimumPeriod)
	mockStorageDigest.On("Get").Return(types.Digest{})

	onInit, err := module.OnInitialize(blockNumber)
	assert.Nil(t, err)

	assert.Equal(t, types.WeightFromParts(3000, 0), onInit)
	mockStorageDigest.AssertCalled(t, "Get")
	mockStorageCurrentSlot.AssertNotCalled(t, "Put", mock.Anything)
}

func Test_Aura_OnInitialize_CurrentSlotMustIncrease(t *testing.T) {
	setup(timestampMinimumPeriod)
	mockStorageDigest.On("Get").Return(*newPreRuntimeDigest(sc.U64(1)))
	mockStorageCurrentSlot.On("Get").Return(sc.U64(2))

	assert.PanicsWithValue(t, errSlotMustIncrease, func() {
		module.OnInitialize(blockNumber)
	})
	mockStorageDigest.AssertCalled(t, "Get")
	mockStorageCurrentSlot.AssertNotCalled(t, "Put", mock.Anything)
}

func Test_Aura_OnInitialize_CurrentSlotUpdate(t *testing.T) {
	setup(timestampMinimumPeriod)
	mockStorageDigest.On("Get").Return(*newPreRuntimeDigest(sc.U64(1)))
	mockStorageCurrentSlot.On("Get").Return(sc.U64(0))
	mockStorageCurrentSlot.On("Put", sc.U64(1)).Return()
	mockStorageAuthorities.On("DecodeLen").Return(sc.NewOption[sc.U64](sc.U64(3)), nil)

	onInit, err := module.OnInitialize(blockNumber)
	assert.Nil(t, err)

	assert.Equal(t, types.WeightFromParts(13_000, 0), onInit)
	mockStorageDigest.AssertCalled(t, "Get")
	mockStorageCurrentSlot.AssertCalled(t, "Put", sc.U64(1))
}

func Test_Aura_OnTimestampSet_DurationCannotBeZero(t *testing.T) {
	setup(0)
	mockStorageCurrentSlot.On("Get").Return(0)

	assert.PanicsWithValue(t, errSlotDurationZero, func() {
		module.OnTimestampSet(1)
	})
}

func Test_Aura_OnTimestampSet_TimestampSlotMismatch(t *testing.T) {
	setup(timestampMinimumPeriod)
	mockStorageCurrentSlot.On("Get").Return(sc.U64(2))

	assert.PanicsWithValue(t, errTimestampSlotMismatch, func() {
		module.OnTimestampSet(sc.U64(4_000))
	})
}

func Test_Aura_SlotDuration(t *testing.T) {
	setup(timestampMinimumPeriod)

	assert.Equal(t, sc.U64(4_000), module.SlotDuration())
}
