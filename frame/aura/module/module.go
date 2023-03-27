package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/support"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type AuraModule struct {
}

func NewAuraModule() AuraModule {
	return AuraModule{}
}

func (am AuraModule) Functions() map[sc.U8]support.FunctionMetadata {
	return map[sc.U8]support.FunctionMetadata{}
}

func (am AuraModule) PreDispatch(_ support.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (am AuraModule) ValidateUnsigned(_ primitives.TransactionSource, _ support.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}
