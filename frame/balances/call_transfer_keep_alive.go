package balances

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type callTransferKeepAlive struct {
	primitives.Callable
	transfer
}

func newCallTransferKeepAlive(moduleId sc.U8, functionId sc.U8, storedMap primitives.StoredMap, constants *consts, mutator accountMutator) primitives.Call {
	call := callTransferKeepAlive{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
			Arguments:  sc.NewVaryingData(types.MultiAddress{}, sc.Compact[sc.U128]{}),
		},
		transfer: newTransfer(moduleId, storedMap, constants, mutator),
	}

	return call
}

func (c callTransferKeepAlive) DecodeArgs(buffer *bytes.Buffer) (primitives.Call, error) {
	dest, err := types.DecodeMultiAddress(buffer)
	if err != nil {
		return nil, err
	}
	value, err := sc.DecodeCompact[sc.Numeric](buffer)
	if err != nil {
		return nil, err
	}
	c.Arguments = sc.NewVaryingData(
		dest,
		value,
	)
	return c, nil
}

func (c callTransferKeepAlive) Encode(buffer *bytes.Buffer) error {
	return c.Callable.Encode(buffer)
}

func (c callTransferKeepAlive) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callTransferKeepAlive) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callTransferKeepAlive) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callTransferKeepAlive) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (c callTransferKeepAlive) BaseWeight() types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `0`
	//  Estimated: `3593`
	// Minimum execution time: 28_184 nanoseconds.
	r := c.constants.DbWeight.Reads(1)
	w := c.constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(49_250_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ callTransferKeepAlive) WeighData(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callTransferKeepAlive) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ callTransferKeepAlive) PaysFee(baseWeight types.Weight) types.Pays {
	return types.PaysYes
}

func (c callTransferKeepAlive) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	valueCompact, _ := args[1].(sc.Compact[sc.Numeric])
	value := valueCompact.Number.(sc.U128)

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

func (_ callTransferKeepAlive) Docs() string {
	return "Same as the [`transfer`] call, but with a check that the transfer will not kill the origin account."
}

// transferKeepAlive is similar to transfer, but includes a check that the origin transactor will not be "killed".
func (c callTransferKeepAlive) transferKeepAlive(origin types.RawOrigin, dest types.MultiAddress, value sc.U128) types.DispatchError {
	if !origin.IsSignedOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}
	transactor, originErr := origin.AsSigned()
	if originErr != nil {
		return primitives.NewDispatchErrorOther(sc.Str(originErr.Error()))
	}

	address, err := types.Lookup(dest)
	if err != nil {
		return types.NewDispatchErrorCannotLookup()
	}

	return c.transfer.trans(transactor, address, value, types.ExistenceRequirementKeepAlive)
}
