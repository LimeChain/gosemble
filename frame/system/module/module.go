package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	cs "github.com/LimeChain/gosemble/constants/system"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/frame/system/dispatchables"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type SystemModule struct {
	functions map[sc.U8]primitives.Call
	// TODO: add more dispatchables
}

func NewSystemModule() SystemModule {
	functions := make(map[sc.U8]primitives.Call)
	functions[cs.FunctionRemarkIndex] = dispatchables.NewRemarkCall(nil)

	return SystemModule{
		functions: functions,
	}
}

func (sm SystemModule) Functions() map[sc.U8]primitives.Call {
	return sm.functions
}

func (sm SystemModule) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (sm SystemModule) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (sm SystemModule) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	metadataModule := primitives.MetadataModule{
		Name: "System",
		Storage: sc.NewOption[primitives.MetadataModuleStorage](primitives.MetadataModuleStorage{
			Prefix: "System",
			Items: sc.Sequence[primitives.MetadataModuleStorageEntry]{
				primitives.NewMetadataModuleStorageEntry(
					"Account",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionMap(
						sc.Sequence[primitives.MetadataModuleStorageHashFunc]{primitives.MetadataModuleStorageHashFuncMultiBlake128Concat},
						sc.ToCompact(metadata.TypesAddress32),
						sc.ToCompact(metadata.TypesAccountInfo)),
					"The full account information for a particular account ID."),
				primitives.NewMetadataModuleStorageEntry(
					"ExtrinsicCount",
					primitives.MetadataModuleStorageEntryModifierOptional,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(
						sc.ToCompact(metadata.PrimitiveTypesU32)),
					"Total extrinsics count for the current block."),
				primitives.NewMetadataModuleStorageEntry(
					"BlockWeight",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(
						sc.ToCompact(metadata.TypesPerDispatchClassWeight)),
					"The current weight for the block."),
				primitives.NewMetadataModuleStorageEntry(
					"AllExtrinsicsLen",
					primitives.MetadataModuleStorageEntryModifierOptional,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(
						sc.ToCompact(metadata.PrimitiveTypesU32)),
					"Total length (in bytes) for all extrinsics put together, for the current block."),
				primitives.NewMetadataModuleStorageEntry(
					"BlockHash",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionMap(
						sc.Sequence[primitives.MetadataModuleStorageHashFunc]{primitives.MetadataModuleStorageHashFuncMultiXX64},
						sc.ToCompact(metadata.PrimitiveTypesU32),
						sc.ToCompact(metadata.TypesFixedSequence32U8)),
					"Map of block numbers to block hashes."),
				primitives.NewMetadataModuleStorageEntry(
					"ExtrinsicData",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionMap(
						sc.Sequence[primitives.MetadataModuleStorageHashFunc]{primitives.MetadataModuleStorageHashFuncMultiXX64},
						sc.ToCompact(metadata.PrimitiveTypesU32),
						sc.ToCompact(metadata.TypesSequenceU8)),
					"Extrinsics data for the current block (maps an extrinsic's index to its data)."),
				primitives.NewMetadataModuleStorageEntry(
					"Number",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(
						sc.ToCompact(metadata.PrimitiveTypesU32)),
					"The current block number being processed. Set by `execute_block`."),
				primitives.NewMetadataModuleStorageEntry(
					"ParentHash",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(
						sc.ToCompact(metadata.TypesFixedSequence32U8)),
					"Hash of the previous block."),
				primitives.NewMetadataModuleStorageEntry(
					"Digest",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(
						sc.ToCompact(metadata.TypesSliceDigestItem)),
					"Digest of the current block, also part of the block header."),
				// TODO: Events
				// TODO: EventTopics
				primitives.NewMetadataModuleStorageEntry(
					"EventCount",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(
						sc.ToCompact(metadata.PrimitiveTypesU32)),
					"The number of events in the `Events<T>` list."),
				// TODO: LastRuntimeUpgrade
				// TODO: ExecutionPhase
			},
		}),
		Call:  sc.NewOption[sc.Compact](sc.ToCompact(metadata.SystemCalls)),
		Event: sc.NewOption[sc.Compact](nil), // TODO:
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{
			primitives.NewMetadataModuleConstant(
				"BlockWeights",
				sc.ToCompact(metadata.TypesBlockWeights),
				sc.BytesToSequenceU8(system.DefaultBlockWeights().Bytes()),
				"Block & extrinsics weights: base values and limits.",
			),
		}, // TODO:
		Error: sc.NewOption[sc.Compact](nil), // TODO:
		Index: cs.ModuleIndex,
	}

	return sm.metadataTypes(), metadataModule
}

func (sm SystemModule) metadataTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataType(metadata.SystemCalls, "System calls", primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"remark",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceU8),
					},
					cs.FunctionRemarkIndex,
					"Make some on-chain remark."),
			})),
	}
}
