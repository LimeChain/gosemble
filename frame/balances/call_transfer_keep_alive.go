package balances

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type callTransferKeepAlive[T primitives.PublicKey] struct {
	primitives.Callable
	transfer
}

func newCallTransferKeepAlive[T primitives.PublicKey](moduleId sc.U8, functionId sc.U8, storedMap primitives.StoredMap, constants *consts, mutator accountMutator) primitives.Call {
	call := callTransferKeepAlive[T]{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		transfer: newTransfer(moduleId, storedMap, constants, mutator),
	}

	return call
}

func (c callTransferKeepAlive[T]) DecodeArgs(buffer *bytes.Buffer) (primitives.Call, error) {
	dest, err := types.DecodeMultiAddress[T](buffer)
	if err != nil {
		return nil, err
	}
	value, err := sc.DecodeCompact(buffer)
	if err != nil {
		return nil, err
	}
	c.Arguments = sc.NewVaryingData(
		dest,
		value,
	)
	return c, nil
}

func (c callTransferKeepAlive[T]) Encode(buffer *bytes.Buffer) error {
	return c.Callable.Encode(buffer)
}

func (c callTransferKeepAlive[T]) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callTransferKeepAlive[T]) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callTransferKeepAlive[T]) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callTransferKeepAlive[T]) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (c callTransferKeepAlive[T]) BaseWeight() types.Weight {
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

func (_ callTransferKeepAlive[T]) WeighData(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callTransferKeepAlive[T]) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ callTransferKeepAlive[T]) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (c callTransferKeepAlive[T]) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	value := sc.U128(args[1].(sc.Compact))

	err := c.transferKeepAlive(origin, args[0].(types.MultiAddress), value)
	if err != nil {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				Err: err,
			},
		}
	}

	return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		HasError: false,
		Ok:       types.PostDispatchInfo{},
	}
}

// transferKeepAlive is similar to transfer, but includes a check that the origin transactor will not be "killed".
func (c callTransferKeepAlive[T]) transferKeepAlive(origin types.RawOrigin, dest types.MultiAddress, value sc.U128) error {
	if !origin.IsSignedOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}
	transactor, originErr := origin.AsSigned()
	if originErr != nil {
		log.Critical(originErr.Error())
	}

	address, err := types.Lookup(dest)
	if err != nil {
		return types.NewDispatchErrorCannotLookup()
	}

	return c.transfer.trans(transactor, address, value, types.ExistenceRequirementKeepAlive)
}
