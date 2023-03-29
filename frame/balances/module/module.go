package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/frame/balances/dispatchables"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type BalancesModule struct {
	functions map[sc.U8]primitives.Call
}

func NewBalancesModule() BalancesModule {
	functions := make(map[sc.U8]primitives.Call)
	functions[balances.FunctionTransferIndex] = dispatchables.NewTransferCall(nil)
	functions[balances.FunctionSetBalanceIndex] = dispatchables.NewSetBalanceCall(nil)
	functions[balances.FunctionForceTransferIndex] = dispatchables.NewForceTransferCall(nil)
	functions[balances.FunctionTransferKeepAliveIndex] = dispatchables.NewTransferKeepAliveCall(nil)
	functions[balances.FunctionTransferAllIndex] = dispatchables.NewTransferAllCall(nil)
	functions[balances.FunctionForceFreeIndex] = dispatchables.NewForceFreeCall(nil)

	return BalancesModule{
		functions: functions,
	}
}

func (bm BalancesModule) Functions() map[sc.U8]primitives.Call {
	return bm.functions
}

func (bm BalancesModule) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (bm BalancesModule) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}
