package balances

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
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
	source, err := primitives.DecodeMultiAddress[T](buffer)
	if err != nil {
		return nil, err
	}
	dest, err := primitives.DecodeMultiAddress[T](buffer)
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

func (c callForceTransfer[T]) BaseWeight() primitives.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `135`
	//  Estimated: `6196`
	// Minimum execution time: 39_713 nanoseconds.
	r := c.constants.DbWeight.Reads(2)
	w := c.constants.DbWeight.Writes(2)
	e := primitives.WeightFromParts(0, 6196)
	return primitives.WeightFromParts(40_360_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ callForceTransfer[T]) WeighData(baseWeight primitives.Weight) primitives.Weight {
	return primitives.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callForceTransfer[T]) ClassifyDispatch(baseWeight primitives.Weight) primitives.DispatchClass {
	return primitives.NewDispatchClassNormal()
}

func (_ callForceTransfer[T]) PaysFee(baseWeight primitives.Weight) primitives.Pays {
	return primitives.NewPaysYes()
}

func (c callForceTransfer[T]) Dispatch(origin primitives.RuntimeOrigin, args sc.VaryingData) primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo] {
	value := sc.U128(args[2].(sc.Compact))

	err := c.forceTransfer(origin, args[0].(primitives.MultiAddress), args[1].(primitives.MultiAddress), value)
	if err.VaryingData != nil {
		return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
			HasError: true,
			Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
				Error: err,
			},
		}
	}

	return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: false,
		Ok:       primitives.PostDispatchInfo{},
	}
}

// forceTransfer transfers liquid free balance from `source` to `dest`.
// Can only be called by ROOT.
func (c callForceTransfer[T]) forceTransfer(origin primitives.RawOrigin, source primitives.MultiAddress, dest primitives.MultiAddress, value sc.U128) primitives.DispatchError {
	if !origin.IsRootOrigin() {
		return primitives.NewDispatchErrorBadOrigin()
	}

	sourceAddress, err := primitives.Lookup(source)
	if err != nil {
		return primitives.NewDispatchErrorCannotLookup()
	}
	destinationAddress, err := primitives.Lookup(dest)
	if err != nil {
		return primitives.NewDispatchErrorCannotLookup()
	}

	return c.transfer.trans(sourceAddress, destinationAddress, value, primitives.ExistenceRequirementAllowDeath)
}
