package testable

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/hooks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	functionTestIndex = iota
)

type Module[N sc.Numeric] struct {
	primitives.DefaultProvideInherent
	hooks.DefaultDispatchModule[N]
	Index     sc.U8
	functions map[sc.U8]primitives.Call
}

func New[N sc.Numeric](index sc.U8) Module[N] {
	functions := make(map[sc.U8]primitives.Call)
	functions[functionTestIndex] = newCallTest(index, functionTestIndex)

	return Module[N]{
		Index:     index,
		functions: functions,
	}
}

func (m Module[N]) GetIndex() sc.U8 {
	return m.Index
}

func (m Module[N]) name() sc.Str {
	return "Testable"
}

func (m Module[N]) Functions() map[sc.U8]primitives.Call {
	return m.functions
}

func (m Module[N]) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (m Module[N]) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (m Module[N]) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	return m.metadataTypes(), primitives.MetadataModule{
		Name:    m.name(),
		Storage: sc.Option[primitives.MetadataModuleStorage]{},
		Call:    sc.NewOption[sc.Compact](sc.ToCompact(metadata.TestableCalls)),
		CallDef: sc.NewOption[primitives.MetadataDefinitionVariant](
			primitives.NewMetadataDefinitionVariantStr(
				m.name(),
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TestableCalls, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<Testable, Runtime>"),
				},
				m.Index,
				"Call.Testable"),
		),
		Event:     sc.NewOption[sc.Compact](nil),
		EventDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{},
		Error:     sc.NewOption[sc.Compact](nil),
		Index:     m.Index,
	}
}

func (m Module[N]) metadataTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithParam(metadata.TestableCalls,
			"Testable calls",
			sc.Sequence[sc.Str]{"frame_system", "testable", "Call"},
			primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"test",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{
							primitives.NewMetadataTypeDefinitionField(metadata.TypesSequenceU8),
						},
						functionTestIndex,
						"Make test"),
				}),
			primitives.NewMetadataEmptyTypeParameter("T")),
	}
}
