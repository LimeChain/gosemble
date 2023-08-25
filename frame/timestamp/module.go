package timestamp

import (
	"bytes"
	"errors"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/hooks"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	functionSetIndex = 0
)

var (
	inherentIdentifier = [8]byte{'t', 'i', 'm', 's', 't', 'a', 'p', '0'}
)

type Module[N sc.Numeric] struct {
	hooks.DefaultDispatchModule[N]
	Index     sc.U8
	Config    *Config
	Storage   *storage
	Constants *consts
	functions map[sc.U8]primitives.Call
}

func New[N sc.Numeric](index sc.U8, config *Config) Module[N] {
	functions := make(map[sc.U8]primitives.Call)
	storage := newStorage()
	constants := newConstants(config.MinimumPeriod)
	functions[functionSetIndex] = newSetCall(index, functionSetIndex, storage, constants, config.OnTimestampSet)

	return Module[N]{
		Index:     index,
		Config:    config,
		Storage:   storage,
		Constants: constants,
		functions: functions,
	}
}

func (m Module[N]) GetIndex() sc.U8 {
	return m.Index
}

func (m Module[N]) name() sc.Str {
	return "Timestamp"
}

func (m Module[N]) Functions() map[sc.U8]primitives.Call {
	return m.functions
}

func (m Module[N]) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (m Module[N]) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.DefaultValidTransaction(), nil
}

func (m Module[N]) OnFinalize(_ N) {
	value := m.Storage.DidUpdate.TakeBytes()
	if value == nil {
		log.Critical("Timestamp must be updated once in the block")
	}
}

func (m Module[N]) CreateInherent(inherent primitives.InherentData) sc.Option[primitives.Call] {
	inherentData := inherent.Data[inherentIdentifier]

	if inherentData == nil {
		log.Critical("Timestamp inherent must be provided.")
	}

	buffer := &bytes.Buffer{}
	buffer.Write(sc.SequenceU8ToBytes(inherentData))
	ts := sc.DecodeU64(buffer)
	// TODO: err if not able to parse it.

	nextTimestamp := m.Storage.Now.Get().Add(m.Constants.MinimumPeriod).(sc.U64)
	if ts.Gt(nextTimestamp) {
		nextTimestamp = ts
	}

	function := newSetCallWithArgs(m.Index, functionSetIndex, sc.NewVaryingData(sc.ToCompact(uint64(nextTimestamp))))

	return sc.NewOption[primitives.Call](function)
}

func (m Module[N]) CheckInherent(call primitives.Call, inherent primitives.InherentData) error {
	if !m.IsInherent(call) {
		return errors.New("invalid inherent check for timestamp module")
	}

	maxTimestampDriftMillis := sc.U64(30 * 1000)

	compactTs := call.Args()[0].(sc.Compact)
	t := sc.To[sc.U64](sc.U128(compactTs))

	inherentData := inherent.Data[inherentIdentifier]

	if inherentData == nil {
		log.Critical("Timestamp inherent must be provided.")
	}

	buffer := &bytes.Buffer{}
	buffer.Write(sc.SequenceU8ToBytes(inherentData))
	ts := sc.DecodeU64(buffer)
	// TODO: err if not able to parse it.

	systemNow := m.Storage.Now.Get()

	minimum := systemNow.Add(m.Constants.MinimumPeriod)
	if t.Gt(ts.Add(maxTimestampDriftMillis)) {
		return primitives.NewTimestampErrorTooFarInFuture()
	} else if t.Lt(minimum) {
		return primitives.NewTimestampErrorTooEarly()
	}

	return nil
}

func (m Module[N]) InherentIdentifier() [8]byte {
	return inherentIdentifier
}

func (m Module[N]) IsInherent(call primitives.Call) bool {
	return call.ModuleIndex().Eq(m.Index) && call.FunctionIndex().Eq(sc.U8(functionSetIndex))
}

func (m Module[N]) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	return m.metadataTypes(), primitives.MetadataModule{
		Name:    m.name(),
		Storage: m.metadataStorage(),
		Call:    sc.NewOption[sc.Compact](sc.ToCompact(metadata.TimestampCalls)),
		CallDef: sc.NewOption[primitives.MetadataDefinitionVariant](
			primitives.NewMetadataDefinitionVariantStr(
				m.name(),
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TimestampCalls, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<Timestamp, Runtime>"),
				},
				m.Index,
				"Call.Timestamp"),
		),
		Event:    sc.NewOption[sc.Compact](nil),
		EventDef: sc.NewOption[primitives.MetadataDefinitionVariant](nil),
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

func (m Module[N]) metadataTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithParam(metadata.TimestampCalls, "Timestamp calls", sc.Sequence[sc.Str]{"pallet_timestamp", "pallet", "Call"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"set",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU64, "now", "T::Moment"),
					},
					functionSetIndex,
					"Set the current time."),
			}), primitives.NewMetadataEmptyTypeParameter("T")),
	}
}

func (m Module[N]) metadataStorage() sc.Option[primitives.MetadataModuleStorage] {
	return sc.NewOption[primitives.MetadataModuleStorage](primitives.MetadataModuleStorage{
		Prefix: m.name(),
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
	})
}
