package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/grandpa"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type GrandpaModule struct {
}

func NewGrandpaModule() GrandpaModule {
	return GrandpaModule{}
}

func (gm GrandpaModule) Functions() map[sc.U8]primitives.Call {
	return map[sc.U8]primitives.Call{}
}

func (gm GrandpaModule) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (gm GrandpaModule) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

func (gm GrandpaModule) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	// TODO: types
	return sc.Sequence[primitives.MetadataType]{}, primitives.MetadataModule{
		Name:      "Grandpa",
		Storage:   sc.Option[primitives.MetadataModuleStorage]{},
		Call:      sc.NewOption[sc.Compact](nil),
		Event:     sc.NewOption[sc.Compact](nil),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{},
		Error:     sc.NewOption[sc.Compact](nil),
		Index:     grandpa.ModuleIndex,
	}
}
