package dispatchables

import (
	"bytes"
	"fmt"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

type FnTransferAll struct{}

func (_ FnTransferAll) Index() sc.U8 {
	return balances.FunctionTransferAllIndex
}

func (_ FnTransferAll) Decode(buffer *bytes.Buffer) sc.VaryingData {
	return sc.NewVaryingData(
		types.DecodeMultiAddress(buffer),
		sc.DecodeBool(buffer),
	)
}

func (_ FnTransferAll) IsInherent() bool {
	return false
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

func (_ FnTransferAll) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ FnTransferAll) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ FnTransferAll) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (fn FnTransferAll) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	err := transferAll(origin, args[0].(types.MultiAddress), bool(args[1].(sc.Bool)))
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

func transferAll(origin types.RawOrigin, dest types.MultiAddress, keepAlive bool) types.DispatchError {
	if !origin.IsSignedOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	transactor := origin.AsSigned()
	reducibleBalance := reducibleBalance(transactor, keepAlive)

	to, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to lookup [%s]", dest.Bytes()))
		return types.NewDispatchErrorCannotLookup()
	}

	keep := types.ExistenceRequirementKeepAlive
	if !keepAlive {
		keep = types.ExistenceRequirementAllowDeath
	}

	return trans(transactor, to, reducibleBalance, keep)
}
