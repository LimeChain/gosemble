package balances

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type callForceTransfer struct {
	primitives.Callable
	transfer
}

func newCallForceTransfer(moduleId sc.U8, functionId sc.U8, storedMap primitives.StoredMap, constants *consts, mutator accountMutator) primitives.Call {
	call := callForceTransfer{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		transfer: newTransfer(moduleId, storedMap, constants, mutator),
	}

	return call
}

func (c callForceTransfer) DecodeArgs(buffer *bytes.Buffer) (primitives.Call, error) {
	source, err := types.DecodeMultiAddress(buffer)
	if err != nil {
		return nil, err
	}
	dest, err := types.DecodeMultiAddress(buffer)
	if err != nil {
		return nil, err
	}
	value, err := sc.DecodeCompact(buffer)
	if err != nil {
		return nil, err
	}
	c.Arguments = sc.NewVaryingData(
		source,
		dest,
		value,
	)
	return c, nil
}

func (c callForceTransfer) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c callForceTransfer) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callForceTransfer) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callForceTransfer) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callForceTransfer) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (c callForceTransfer) BaseWeight() types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `135`
	//  Estimated: `6196`
	// Minimum execution time: 39_713 nanoseconds.
	r := c.constants.DbWeight.Reads(2)
	w := c.constants.DbWeight.Writes(2)
	e := types.WeightFromParts(0, 6196)
	return types.WeightFromParts(40_360_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ callForceTransfer) WeighData(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callForceTransfer) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ callForceTransfer) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (c callForceTransfer) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	value := sc.U128(args[2].(sc.Compact))

	err := c.forceTransfer(origin, args[0].(types.MultiAddress), args[1].(types.MultiAddress), value)
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

// forceTransfer transfers liquid free balance from `source` to `dest`.
// Can only be called by ROOT.
func (c callForceTransfer) forceTransfer(origin types.RawOrigin, source types.MultiAddress, dest types.MultiAddress, value sc.U128) types.DispatchError {
	if !origin.IsRootOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	sourceAddress, err := types.Lookup(source)
	if err != nil {
		return types.NewDispatchErrorCannotLookup()
	}
	destinationAddress, err := types.Lookup(dest)
	if err != nil {
		return types.NewDispatchErrorCannotLookup()
	}

	return c.transfer.trans(sourceAddress, destinationAddress, value, types.ExistenceRequirementAllowDeath)
}
