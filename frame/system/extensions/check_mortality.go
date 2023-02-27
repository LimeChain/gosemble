package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
)

// to be able to provide a custom implementation of the Validate function
type CheckMortality types.Era

func (e CheckMortality) Validate(_who *types.Address32, _call *types.Call, _info *types.DispatchInfo, _length sc.Compact) (ok types.ValidTransaction, err types.TransactionValidityError) {
	currentU64 := sc.U64(system.StorageGetBlockNumber()) // TDOO: per module implementation

	validTill := types.Era(e).Death(currentU64)

	ok = types.DefaultValidTransaction()
	ok.Longevity = validTill.SaturatingSub(currentU64)

	return ok, err
}

func (e CheckMortality) PreDispatch(who *types.Address32, call *types.Call, info *types.DispatchInfo, length sc.Compact) (ok types.Pre, err types.TransactionValidityError) {
	_, err = e.Validate(who, call, info, length)
	return ok, err
}
