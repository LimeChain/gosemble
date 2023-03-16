package dispatchables

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	bc "github.com/LimeChain/gosemble/frame/balances/constants"
	"github.com/LimeChain/gosemble/primitives/types"
)

type FnForceTransfer struct{}

func (_ FnForceTransfer) Index() sc.U8 {
	return bc.FunctionForceTransferIndex
}

func (_ FnForceTransfer) BaseWeight(b ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `135`
	//  Estimated: `6196`
	// Minimum execution time: 39_713 nanoseconds.
	r := constants.DbWeight.Reads(2)
	w := constants.DbWeight.Writes(2)
	e := types.WeightFromParts(0, 6196)
	return types.WeightFromParts(40_360_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ FnForceTransfer) WeightInfo(baseWeight types.Weight, target []byte) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ FnForceTransfer) ClassifyDispatch(baseWeight types.Weight, target []byte) types.DispatchClass {
	return types.NewDispatchClassMandatory()
}

func (_ FnForceTransfer) PaysFee(baseWeight types.Weight, target []byte) types.Pays {
	return types.NewPaysYes()
}

func (fn FnForceTransfer) Dispatch(origin types.RuntimeOrigin, source types.MultiAddress, dest types.MultiAddress, value sc.U128) (ok sc.Empty, err types.DispatchError) {
	return forceTransfer(origin, source, dest, value)
}

func forceTransfer(origin types.RawOrigin, source types.MultiAddress, dest types.MultiAddress, value sc.U128) (sc.Empty, types.DispatchError) {
	if !origin.IsRootOrigin() {
		return sc.Empty{}, types.NewDispatchErrorBadOrigin()
	}

	sourceAddress, err := types.DefaultAccountIdLookup().Lookup(source)
	if err != nil {
		return sc.Empty{}, types.NewDispatchErrorCannotLookup()
	}
	destinationAddress, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		return sc.Empty{}, types.NewDispatchErrorCannotLookup()
	}

	e := trans(sourceAddress, destinationAddress, value, types.ExistenceRequirementAllowDeath)
	if e != nil {
		return sc.Empty{}, e
	}

	return sc.Empty{}, nil
}
