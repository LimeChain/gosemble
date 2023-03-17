package dispatchables

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/primitives/types"
)

type FnTransferKeepAlive struct{}

func (_ FnTransferKeepAlive) Index() sc.U8 {
	return balances.FunctionTransferKeepAliveIndex
}

func (_ FnTransferKeepAlive) Decode(buffer *bytes.Buffer) []sc.Encodable {
	return []sc.Encodable{
		types.DecodeMultiAddress(buffer),
		sc.U128(sc.DecodeCompact(buffer)),
	}
}

func (_ FnTransferKeepAlive) IsInherent() bool {
	return false
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

func (_ FnTransferKeepAlive) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ FnTransferKeepAlive) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassMandatory()
}

func (_ FnTransferKeepAlive) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (fn FnTransferKeepAlive) Dispatch(origin types.RuntimeOrigin, args ...sc.Encodable) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	err := transferKeepAlive(origin, args[0].(types.MultiAddress), args[1].(sc.U128))
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

func transferKeepAlive(origin types.RawOrigin, dest types.MultiAddress, value sc.U128) types.DispatchError {
	if !origin.IsSignedOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}
	transactor := origin.AsSigned()

	address, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		return types.NewDispatchErrorCannotLookup()
	}

	return trans(transactor, address, value, types.ExistenceRequirementKeepAlive)
}
