package balances

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/balances/constants"
	"github.com/LimeChain/gosemble/frame/balances/dispatchables"
	"github.com/LimeChain/gosemble/primitives/support"
)

var Module = BalancesModule{}

type BalancesModule struct {
	Transfer          dispatchables.FnTransfer
	SetBalance        dispatchables.FnSetBalance
	ForceTransfer     dispatchables.FnForceTransfer
	TransferKeepAlive dispatchables.FnTransferKeepAlive
	TransferAll       dispatchables.FnTransferAll
	ForceFree         dispatchables.FnForceFree
}

func (bm BalancesModule) Functions() []support.FunctionMetadata {
	return []support.FunctionMetadata{
		bm.Transfer,
		bm.SetBalance,
		bm.ForceTransfer,
		bm.TransferKeepAlive,
		bm.TransferAll,
		bm.ForceFree,
	}
}

func (bm BalancesModule) Index() sc.U8 {
	return constants.ModuleIndex
}
