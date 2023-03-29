package module

import (
	sc "github.com/LimeChain/goscale"
	cs "github.com/LimeChain/gosemble/constants/system"
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
