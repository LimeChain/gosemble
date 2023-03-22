package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/frame/balances/dispatchables"
	"github.com/LimeChain/gosemble/primitives/support"
)

type BalancesModule struct {
	functions map[sc.U8]support.FunctionMetadata
}

func NewBalancesModule() BalancesModule {
	functions := make(map[sc.U8]support.FunctionMetadata)
	functions[balances.FunctionTransferIndex] = dispatchables.FnTransfer{}
	functions[balances.FunctionSetBalanceIndex] = dispatchables.FnSetBalance{}
	functions[balances.FunctionForceTransferIndex] = dispatchables.FnForceTransfer{}
	functions[balances.FunctionTransferKeepAliveIndex] = dispatchables.FnTransferKeepAlive{}
	functions[balances.FunctionTransferAllIndex] = dispatchables.FnTransferAll{}
	functions[balances.FunctionForceFreeIndex] = dispatchables.FnForceFree{}

	return BalancesModule{
		functions: functions,
	}
}

func (bm BalancesModule) Functions() map[sc.U8]support.FunctionMetadata {
	return bm.functions
}