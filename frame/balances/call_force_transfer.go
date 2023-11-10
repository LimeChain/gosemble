package balances

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type testPublicKeyType = primitives.Ed25519PublicKey

type callForceTransfer[T primitives.PublicKey] struct {
	primitives.Callable
	transfer
}

func newCallForceTransfer[T primitives.PublicKey](moduleId sc.U8, functionId sc.U8, storedMap primitives.StoredMap, constants *consts, mutator accountMutator) primitives.Call {
	call := callForceTransfer[T]{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		transfer: newTransfer(moduleId, storedMap, constants, mutator),
	}

	return call
}

func (c callForceTransfer[T]) DecodeArgs(buffer *bytes.Buffer) (primitives.Call, error) {
	source, err := types.DecodeMultiAddress[T](buffer)
	if err != nil {
		return nil, err
	}
	dest, err := types.DecodeMultiAddress[T](buffer)
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

func (c callForceTransfer[T]) Encode(buffer *bytes.Buffer) error {
	return c.Callable.Encode(buffer)
}

func (c callForceTransfer[T]) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callForceTransfer[T]) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callForceTransfer[T]) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callForceTransfer[T]) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (c callForceTransfer[T]) BaseWeight() types.Weight {
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

func (_ callForceTransfer[T]) WeighData(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callForceTransfer[T]) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ callForceTransfer[T]) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (c callForceTransfer[T]) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
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
func (c callForceTransfer[T]) forceTransfer(origin types.RawOrigin, source types.MultiAddress, dest types.MultiAddress, value sc.U128) types.DispatchError {
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
