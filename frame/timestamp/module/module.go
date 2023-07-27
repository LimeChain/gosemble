package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	FunctionSetIndex = 0
)

type Module struct {
	Index     sc.U8
	Config    *Config
	Storage   *storage
	Constants *consts
	functions map[sc.U8]primitives.Call
}

func NewModule(index sc.U8, config *Config) Module {
	functions := make(map[sc.U8]primitives.Call)
	storage := newStorage()
	constants := newConstants(config.MinimumPeriod)
	functions[FunctionSetIndex] = NewSetCall(index, FunctionSetIndex, nil, storage, constants, config.OnTimestampSet)

	return Module{
		Index:     index,
		Config:    config,
		Storage:   storage,
		Constants: constants,
		functions: functions,
	}
}

func (m Module) Functions() map[sc.U8]primitives.Call {
	return m.functions
}

func (m Module) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (m Module) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (m Module) OnFinalize() {
	value := m.Storage.DidUpdate.Take()
	if value == nil {
		log.Critical("Timestamp must be updated once in the block")
	}
}

func (m Module) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	return m.metadataTypes(), primitives.MetadataModule{
		Name: "Timestamp",
		Storage: sc.NewOption[primitives.MetadataModuleStorage](primitives.MetadataModuleStorage{
			Prefix: "Timestamp",
			Items: sc.Sequence[primitives.MetadataModuleStorageEntry]{
				primitives.NewMetadataModuleStorageEntry(
					"Now",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.PrimitiveTypesU64)),
					"Current time for the current block."),
				primitives.NewMetadataModuleStorageEntry(
					"DidUpdate",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.PrimitiveTypesBool)),
					"Did the timestamp get updated in this block?"),
			},
		}),
		Call:  sc.NewOption[sc.Compact](sc.ToCompact(metadata.TimestampCalls)),
		Event: sc.NewOption[sc.Compact](nil),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{
			primitives.NewMetadataModuleConstant(
				"MinimumPeriod",
				sc.ToCompact(metadata.PrimitiveTypesU64),
				sc.BytesToSequenceU8(m.Constants.MinimumPeriod.Bytes()),
				"The minimum period between blocks. Beware that this is different to the *expected*  period that the block production apparatus provides.",
			),
		},
		Error: sc.NewOption[sc.Compact](nil),
		Index: m.Index,
	}
}

func (m Module) metadataTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithParam(metadata.TimestampCalls, "Timestamp calls", sc.Sequence[sc.Str]{"pallet_timestamp", "pallet", "Call"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"set",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU64, "now", "T::Moment"),
					},
					FunctionSetIndex,
					"Set the current time."),
			}), primitives.NewMetadataEmptyTypeParameter("T")),
	}
}
