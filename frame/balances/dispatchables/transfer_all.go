package dispatchables

import (
	"fmt"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	bc "github.com/LimeChain/gosemble/frame/balances/constants"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

type FnTransferAll struct{}

func (_ FnTransferAll) Index() sc.U8 {
	return bc.FunctionTransferAllIndex
}

func (_ FnTransferAll) BaseWeight(b ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `0`
	//  Estimated: `3593`
	// Minimum execution time: 34_878 nanoseconds.
	r := constants.DbWeight.Reads(1)
	w := constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(35_121_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ FnTransferAll) WeightInfo(baseWeight types.Weight, target []byte) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ FnTransferAll) ClassifyDispatch(baseWeight types.Weight, target []byte) types.DispatchClass {
	return types.NewDispatchClassMandatory()
}

func (_ FnTransferAll) PaysFee(baseWeight types.Weight, target []byte) types.Pays {
	return types.NewPaysYes()
}

func (fn FnTransferAll) Dispatch(origin types.RuntimeOrigin, dest types.MultiAddress, keepAlive bool) (ok sc.Empty, err types.DispatchError) {
	return transferAll(origin, dest, keepAlive)
}

func transferAll(origin types.RawOrigin, dest types.MultiAddress, keepAlive bool) (sc.Empty, types.DispatchError) {
	if !origin.IsSignedOrigin() {
		return sc.Empty{}, types.NewDispatchErrorBadOrigin()
	}

	transactor := origin.AsSigned()
	reducibleBalance := reducibleBalance(transactor, keepAlive)

	to, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to lookup [%s]", dest.Bytes()))
		return sc.Empty{}, types.NewDispatchErrorCannotLookup()
	}

	keep := types.ExistenceRequirementKeepAlive
	if !keepAlive {
		keep = types.ExistenceRequirementAllowDeath
	}

	return sc.Empty{}, trans(transactor, to, reducibleBalance, keep)
}
