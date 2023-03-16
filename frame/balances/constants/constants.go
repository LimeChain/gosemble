package constants

import (
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
)

const (
	ModuleIndex                    = sc.U8(5)
	FunctionTransferIndex          = 0
	FunctionSetBalanceIndex        = 1
	FunctionForceTransferIndex     = 2
	FunctionTransferKeepAliveIndex = 3
	FunctionTransferAllIndex       = 4
	FunctionForceFreeIndex         = 5
)

const (
	existentialDeposit = 1 * constants.Dollar
	MaxLocks           = 50
	MaxReserves        = 50
)

var (
	ExistentialDeposit = big.NewInt(0).SetUint64(existentialDeposit)
)
