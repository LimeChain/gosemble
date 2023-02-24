package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
)

// to be able to provide a custom implementation of the Validate function
type CheckEra types.Era

func (e CheckEra) Validate(_who *types.Address32, _call *types.Call, _info *types.DispatchInfo, _length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	currentU64 := sc.U64(system.StorageGetBlockNumber()) // TDOO: module's implementation

	validTill := types.Era(e).Death(currentU64)

	ok = types.DefaultValidTransaction()
	ok.Longevity = validTill.SaturatingSub(currentU64)

	return ok, err
}
