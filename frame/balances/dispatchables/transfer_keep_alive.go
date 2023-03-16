package dispatchables

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	bc "github.com/LimeChain/gosemble/frame/balances/constants"
	"github.com/LimeChain/gosemble/primitives/types"
)

type FnTransferKeepAlive struct{}

func (_ FnTransferKeepAlive) Index() sc.U8 {
	return bc.FunctionTransferKeepAliveIndex
}

func (_ FnTransferKeepAlive) BaseWeight(b ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `0`
	//  Estimated: `3593`
	// Minimum execution time: 28_184 nanoseconds.
	r := constants.DbWeight.Reads(1)
	w := constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(49_250_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ FnTransferKeepAlive) WeightInfo(baseWeight types.Weight, target []byte) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ FnTransferKeepAlive) ClassifyDispatch(baseWeight types.Weight, target []byte) types.DispatchClass {
	return types.NewDispatchClassMandatory()
}

func (_ FnTransferKeepAlive) PaysFee(baseWeight types.Weight, target []byte) types.Pays {
	return types.NewPaysYes()
}

func (fn FnTransferKeepAlive) Dispatch(origin types.RuntimeOrigin, dest types.MultiAddress, value sc.U128) (ok sc.Empty, err types.DispatchError) {
	return transferKeepAlive(origin, dest, value)
}

func transferKeepAlive(origin types.RawOrigin, dest types.MultiAddress, value sc.U128) (sc.Empty, types.DispatchError) {
	if !origin.IsSignedOrigin() {
		return sc.Empty{}, types.NewDispatchErrorBadOrigin()
	}
	transactor := origin.AsSigned()

	address, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		return sc.Empty{}, types.NewDispatchErrorCannotLookup()
	}

	return sc.Empty{}, trans(transactor, address, value, types.ExistenceRequirementKeepAlive)
}
