package balances

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
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
	dest, err := primitives.DecodeMultiAddress[T](buffer)
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

func (c callTransferKeepAlive[T]) BaseWeight() primitives.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `0`
	//  Estimated: `3593`
	// Minimum execution time: 28_184 nanoseconds.
	r := c.constants.DbWeight.Reads(1)
	w := c.constants.DbWeight.Writes(1)
	e := primitives.WeightFromParts(0, 3593)
	return primitives.WeightFromParts(49_250_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ callTransferKeepAlive[T]) WeighData(baseWeight primitives.Weight) primitives.Weight {
	return primitives.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callTransferKeepAlive[T]) ClassifyDispatch(baseWeight primitives.Weight) primitives.DispatchClass {
	return primitives.NewDispatchClassNormal()
}

func (_ callTransferKeepAlive[T]) PaysFee(baseWeight primitives.Weight) primitives.Pays {
	return primitives.NewPaysYes()
}

func (c callTransferKeepAlive[T]) Dispatch(origin primitives.RuntimeOrigin, args sc.VaryingData) primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo] {
	value := sc.U128(args[1].(sc.Compact))

	err := c.transferKeepAlive(origin, args[0].(primitives.MultiAddress), value)
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

// transferKeepAlive is similar to transfer, but includes a check that the origin transactor will not be "killed".
func (c callTransferKeepAlive[T]) transferKeepAlive(origin primitives.RawOrigin, dest primitives.MultiAddress, value sc.U128) primitives.DispatchError {
	if !origin.IsSignedOrigin() {
		return primitives.NewDispatchErrorBadOrigin()
	}
	transactor, originErr := origin.AsSigned()
	if originErr != nil {
		log.Critical(originErr.Error())
	}

	address, err := primitives.Lookup(dest)
	if err != nil {
		return primitives.NewDispatchErrorCannotLookup()
	}

	return c.transfer.trans(transactor, address, value, primitives.ExistenceRequirementKeepAlive)
}
