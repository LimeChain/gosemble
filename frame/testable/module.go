package testable

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/hooks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	functionTestIndex = iota
)

type Module struct {
	primitives.DefaultInherentProvider
	hooks.DefaultDispatchModule
	Index     sc.U8
	functions map[sc.U8]primitives.Call
}

func New(index sc.U8) Module {
	functions := make(map[sc.U8]primitives.Call)
	functions[functionTestIndex] = newCallTest(index, functionTestIndex)

	return Module{
		Index:     index,
		functions: functions,
	}
}

func (m Module) GetIndex() sc.U8 {
	return m.Index
}

func (m Module) name() sc.Str {
	return "Testable"
}

func (m Module) Functions() map[sc.U8]primitives.Call {
	return m.functions
}

func (m Module) PreDispatch(_ primitives.Call) (sc.Empty, error) {
	return sc.Empty{}, nil
}

func (m Module) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, error) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (m Module) Metadata(mdGenerator *primitives.MetadataTypeGenerator) primitives.MetadataModule {
	testableCallsMetadataId := mdGenerator.BuildCallsMetadata("Testable", m.functions, &sc.Sequence[primitives.MetadataTypeParameter]{primitives.NewMetadataEmptyTypeParameter("T")})

	dataV14 := primitives.MetadataModuleV14{
		Name:    m.name(),
		Storage: sc.Option[primitives.MetadataModuleStorage]{},
		Call:    sc.NewOption[sc.Compact](sc.ToCompact(testableCallsMetadataId)),
		CallDef: sc.NewOption[primitives.MetadataDefinitionVariant](
			primitives.NewMetadataDefinitionVariantStr(
				m.name(),
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(testableCallsMetadataId, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<Testable, Runtime>"),
				},
				m.Index,
				"Call.Testable"),
		),
		Event:     sc.NewOption[sc.Compact](nil),
		EventDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{},
		Error:     sc.NewOption[sc.Compact](nil),
		ErrorDef:  sc.NewOption[primitives.MetadataDefinitionVariant](nil),
		Index:     m.Index,
	}

	mdGenerator.AppendMetadataTypes(m.metadataTypes())

	return primitives.MetadataModule{
		Version:   primitives.ModuleVersion14,
		ModuleV14: dataV14,
	}
}

func (m Module) metadataTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{}
}
