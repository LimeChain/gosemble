package balances

import (
	"bytes"
	"fmt"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type callTransferAll[T primitives.PublicKey] struct {
	primitives.Callable
	transfer
}

func newCallTransferAll[T primitives.PublicKey](moduleId sc.U8, functionId sc.U8, storedMap primitives.StoredMap, constants *consts, mutator accountMutator) primitives.Call {
	call := callTransferAll[T]{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		transfer: newTransfer(moduleId, storedMap, constants, mutator),
	}

	return call
}

func (c callTransferAll[T]) DecodeArgs(buffer *bytes.Buffer) (primitives.Call, error) {
	dest, err := primitives.DecodeMultiAddress[T](buffer)
	if err != nil {
		return nil, err
	}
	keepAlive, err := sc.DecodeBool(buffer)
	if err != nil {
		return nil, err
	}
	c.Arguments = sc.NewVaryingData(
		dest,
		keepAlive,
	)
	return c, nil
}

func (c callTransferAll[T]) Encode(buffer *bytes.Buffer) error {
	return c.Callable.Encode(buffer)
}

func (c callTransferAll[T]) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callTransferAll[T]) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callTransferAll[T]) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callTransferAll[T]) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (c callTransferAll[T]) BaseWeight() primitives.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `0`
	//  Estimated: `3593`
	// Minimum execution time: 34_878 nanoseconds.
	r := c.constants.DbWeight.Reads(1)
	w := c.constants.DbWeight.Writes(1)
	e := primitives.WeightFromParts(0, 3593)
	return primitives.WeightFromParts(35_121_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ callTransferAll[T]) WeighData(baseWeight primitives.Weight) primitives.Weight {
	return primitives.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callTransferAll[T]) ClassifyDispatch(baseWeight primitives.Weight) primitives.DispatchClass {
	return primitives.NewDispatchClassNormal()
}

func (_ callTransferAll[T]) PaysFee(baseWeight primitives.Weight) primitives.Pays {
	return primitives.NewPaysYes()
}

func (c callTransferAll[T]) Dispatch(origin primitives.RuntimeOrigin, args sc.VaryingData) primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo] {
	err := c.transferAll(origin, args[0].(primitives.MultiAddress), bool(args[1].(sc.Bool)))
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

// transferAll transfers the entire transferable balance from `origin` to `dest`.
// By transferable it means that any locked or reserved amounts will not be transferred.
// `keepAlive`: A boolean to determine if the `transfer_all` operation should send all
// the funds the account has, causing the sender account to be killed (false), or
// transfer everything except at least the existential deposit, which will guarantee to
// keep the sender account alive (true).
func (c callTransferAll[T]) transferAll(origin primitives.RawOrigin, dest primitives.MultiAddress, keepAlive bool) primitives.DispatchError {
	if !origin.IsSignedOrigin() {
		return primitives.NewDispatchErrorBadOrigin()
	}

	transactor, err := origin.AsSigned()
	if err != nil {
		log.Critical(err.Error())
	}

	reducibleBalance, err := c.reducibleBalance(transactor, keepAlive)
	if err != nil {
		return primitives.NewDispatchErrorOther(sc.Str(err.Error()))
	}

	to, errLookup := primitives.Lookup(dest)
	if errLookup != nil {
		log.Debug(fmt.Sprintf("Failed to lookup [%s]", dest.Bytes()))
		return primitives.NewDispatchErrorCannotLookup()
	}

	keep := primitives.ExistenceRequirementKeepAlive
	if !keepAlive {
		keep = primitives.ExistenceRequirementAllowDeath
	}

	return c.transfer.trans(transactor, to, reducibleBalance, keep)
}
