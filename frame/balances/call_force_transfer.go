package balances

import (
	"bytes"
	"errors"

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
			Arguments:  sc.NewVaryingData(types.MultiAddress{}, types.MultiAddress{}, sc.Compact{Number: sc.U128{}}),
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
	value, err := sc.DecodeCompact[sc.U128](buffer)
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

func (c callForceTransfer) Encode(buffer *bytes.Buffer) error {
	return c.Callable.Encode(buffer)
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
	return types.PaysYes
}

func (c callForceTransfer) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) (types.PostDispatchInfo, error) {
	valueCompact, ok := args[2].(sc.Compact)
	if !ok {
		return types.PostDispatchInfo{}, errors.New("invalid Compact value when dispatching call_force_transfer")
	}
	value, ok := valueCompact.Number.(sc.U128)
	if !ok {
		return types.PostDispatchInfo{}, errors.New("invalid Compact field number when dispatching call_force_transfer")
	}
	return types.PostDispatchInfo{}, c.forceTransfer(origin, args[0].(types.MultiAddress), args[1].(types.MultiAddress), value)
}

func (_ callForceTransfer) Docs() string {
	return "Exactly as `transfer`, except the origin must be root and the source account may be specified."
}

// forceTransfer transfers liquid free balance from `source` to `dest`.
// Can only be called by ROOT.
func (c callForceTransfer) forceTransfer(origin types.RawOrigin, source types.MultiAddress, dest types.MultiAddress, value sc.U128) error {
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
