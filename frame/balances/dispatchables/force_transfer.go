package dispatchables

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/primitives/types"
)

type FnForceTransfer struct{}

func (_ FnForceTransfer) Index() sc.U8 {
	return balances.FunctionForceTransferIndex
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

func (_ FnForceTransfer) Decode(buffer *bytes.Buffer) []sc.Encodable {
	return []sc.Encodable{
		types.DecodeMultiAddress(buffer),
		types.DecodeMultiAddress(buffer),
		sc.U128(sc.DecodeCompact(buffer)),
	}
}

func (_ FnForceTransfer) IsInherent() bool {
	return false
}

func (_ FnForceTransfer) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ FnForceTransfer) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassMandatory()
}

func (_ FnForceTransfer) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (fn FnForceTransfer) Dispatch(origin types.RuntimeOrigin, args ...sc.Encodable) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	err := forceTransfer(origin, args[0].(types.MultiAddress), args[1].(types.MultiAddress), args[2].(sc.U128))
	if err != nil {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				Error: err,
			},
		}
	}

	return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		HasError: false,
		Ok:       types.PostDispatchInfo{},
	}
}

func forceTransfer(origin types.RawOrigin, source types.MultiAddress, dest types.MultiAddress, value sc.U128) types.DispatchError {
	if !origin.IsRootOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	sourceAddress, err := types.DefaultAccountIdLookup().Lookup(source)
	if err != nil {
		return types.NewDispatchErrorCannotLookup()
	}
	destinationAddress, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		return types.NewDispatchErrorCannotLookup()
	}

	return trans(sourceAddress, destinationAddress, value, types.ExistenceRequirementAllowDeath)
}
