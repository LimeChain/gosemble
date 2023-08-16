package balances

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type transferKeepAliveCall struct {
	primitives.Callable
	transfer
}

func newTransferKeepAliveCall(moduleId sc.U8, functionId sc.U8, storedMap primitives.StoredMap, constants *consts, mutator accountMutator) primitives.Call {
	call := transferKeepAliveCall{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		transfer: newTransfer(moduleId, storedMap, constants, mutator),
	}

	return call
}

func (c transferKeepAliveCall) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(
		types.DecodeMultiAddress(buffer),
		sc.DecodeCompact(buffer),
	)
	return c
}

func (c transferKeepAliveCall) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c transferKeepAliveCall) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c transferKeepAliveCall) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c transferKeepAliveCall) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c transferKeepAliveCall) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (_ transferKeepAliveCall) BaseWeight(b ...any) types.Weight {
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

func (_ transferKeepAliveCall) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ transferKeepAliveCall) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ transferKeepAliveCall) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (c transferKeepAliveCall) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	value := sc.U128(args[1].(sc.Compact))

	err := c.transferKeepAlive(origin, args[0].(types.MultiAddress), value)
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

// transferKeepAlive is similar to transfer, but includes a check that the origin transactor will not be "killed".
func (c transferKeepAliveCall) transferKeepAlive(origin types.RawOrigin, dest types.MultiAddress, value sc.U128) types.DispatchError {
	if !origin.IsSignedOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}
	transactor := origin.AsSigned()

	address, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		return types.NewDispatchErrorCannotLookup()
	}

	return c.transfer.trans(transactor, address, value, types.ExistenceRequirementKeepAlive)
}
