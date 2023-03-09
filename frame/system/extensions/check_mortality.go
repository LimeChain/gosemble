package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
)

func (e CheckMortality) AdditionalSigned() (ok types.H256, err types.TransactionValidityError) {
	current := sc.U64(system.StorageGetBlockNumber()) // TODO: impl saturated_into::<u64>()
	n := sc.U32(types.Era(e).Birth(current))          // TODO: impl saturated_into::<T::BlockNumber>()

	if !system.StorageExistsBlockHash(n) {
		err = types.NewTransactionValidityError(types.NewInvalidTransactionAncientBirthBlock())
		return ok, err
	} else {
		ok = types.H256(system.StorageGetBlockHash(n))
	}

	return ok, err
}

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
