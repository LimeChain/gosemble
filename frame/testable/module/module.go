package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/hooks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	functionTestIndex = 255
)

type TestableModule struct {
	primitives.DefaultProvideInherent
	hooks.DefaultDispatchModule[sc.U32]
	Index     sc.U8
	functions map[sc.U8]primitives.Call
}

func NewTestingModule(index sc.U8) TestableModule {
	functions := make(map[sc.U8]primitives.Call)
	functions[functionTestIndex] = newTestCall(index, functionTestIndex)

	return TestableModule{
		Index:     index,
		functions: functions,
	}
}

func (tm TestableModule) GetIndex() sc.U8 {
	return tm.Index
}

func (tm TestableModule) Functions() map[sc.U8]primitives.Call {
	return tm.functions
}

func (tm TestableModule) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (tm TestableModule) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (tm TestableModule) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	return tm.metadataTypes(), primitives.MetadataModule{
		Name:      "Testable",
		Storage:   sc.Option[primitives.MetadataModuleStorage]{},
		Call:      sc.NewOption[sc.Compact](sc.ToCompact(metadata.TestableCalls)),
		Event:     sc.NewOption[sc.Compact](nil),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{},
		Error:     sc.NewOption[sc.Compact](nil),
		Index:     tm.Index,
	}
}

func (tm TestableModule) metadataTypes() sc.Sequence[primitives.MetadataType] {
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
