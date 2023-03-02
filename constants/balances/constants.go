package balances

import "github.com/LimeChain/gosemble/constants"

const (
	ModuleIndex                    = 5
	FunctionTransferIndex          = 0
	FunctionSetBalanceIndex        = 1
	FunctionForceTransferIndex     = 2
	FunctionTransferKeepAliceIndex = 3
	FunctionTransferAllIndex       = 4
	FunctionForceFreeIndex         = 5
)

const (
	ExistentialDeposit = 1 * constants.Dollar
	MaxLocks           = 50
	MaxReserves        = 50
)
