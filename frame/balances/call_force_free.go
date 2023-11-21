package balances

import (
	"bytes"
	"fmt"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type callForceFree[T primitives.PublicKey] struct {
	primitives.Callable
	transfer
}

func newCallForceFree[T primitives.PublicKey](moduleId sc.U8, functionId sc.U8, storedMap primitives.StoredMap, constants *consts, mutator accountMutator) primitives.Call {
	call := callForceFree[T]{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		transfer: newTransfer(moduleId, storedMap, constants, mutator),
	}

	return call
}

func (c callForceFree[T]) DecodeArgs(buffer *bytes.Buffer) (primitives.Call, error) {
	who, err := primitives.DecodeMultiAddress[T](buffer)
	if err != nil {
		return nil, err
	}
	amount, err := sc.DecodeU128(buffer)
	if err != nil {
		return nil, err
	}
	c.Arguments = sc.NewVaryingData(
		who,
		amount,
	)
	return c, nil
}

func (c callForceFree[T]) Encode(buffer *bytes.Buffer) error {
	return c.Callable.Encode(buffer)
}

func (c callForceFree[T]) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callForceFree[T]) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callForceFree[T]) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callForceFree[T]) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (c callForceFree[T]) BaseWeight() primitives.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `206`
	//  Estimated: `3593`
	// Minimum execution time: 16_790 nanoseconds.
	r := c.constants.DbWeight.Reads(1)
	w := c.constants.DbWeight.Writes(1)
	e := primitives.WeightFromParts(0, 3593)
	return primitives.WeightFromParts(17_029_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ callForceFree[T]) WeighData(baseWeight primitives.Weight) primitives.Weight {
	return primitives.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callForceFree[T]) ClassifyDispatch(baseWeight primitives.Weight) primitives.DispatchClass {
	return primitives.NewDispatchClassNormal()
}

func (_ callForceFree[T]) PaysFee(baseWeight primitives.Weight) primitives.Pays {
	return primitives.NewPaysYes()
}

func (c callForceFree[T]) Dispatch(origin primitives.RuntimeOrigin, args sc.VaryingData) primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo] {
	amount := args[1].(sc.U128)

	err := c.forceFree(origin, args[0].(primitives.MultiAddress), amount)
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

// forceFree frees some balance from a user by force.
// Can only be called by ROOT.
// Consider Substrate fn force_unreserve
func (c callForceFree[T]) forceFree(origin primitives.RawOrigin, who primitives.MultiAddress, amount sc.U128) primitives.DispatchError {
	if !origin.IsRootOrigin() {
		return primitives.NewDispatchErrorBadOrigin()
	}

	target, err := primitives.Lookup(who)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to lookup [%s]", who.Bytes()))
		return primitives.NewDispatchErrorCannotLookup()
	}

	if _, err := c.force(target, amount); err != nil {
		return primitives.NewDispatchErrorOther(sc.Str(err.Error()))
	}

	return primitives.DispatchError{VaryingData: nil}
}

// forceFree frees funds, returning the amount that has not been freed.
func (c callForceFree[T]) force(who primitives.AccountId[primitives.PublicKey], value sc.U128) (sc.U128, error) {
	if value.Eq(constants.Zero) {
		return constants.Zero, nil
	}

	account, err := c.storedMap.Get(who)
	if err != nil {
		return sc.U128{}, err
	}

	totalBalance := account.Data.Total()
	if totalBalance.Eq(constants.Zero) {
		return value, nil
	}

	result := c.accountMutator.tryMutateAccount(
		who,
		func(account *primitives.AccountData, _ bool) sc.Result[sc.Encodable] {
			return removeReserveAndFree(account, value)
		},
	)

	if result.HasError {
		return value, nil
	}

	actual := result.Value.(sc.U128)
	c.storedMap.DepositEvent(newEventUnreserved(c.ModuleId, who, actual))

	return value.Sub(actual), nil
}

// removeReserveAndFree frees reserved value from the account.
func removeReserveAndFree(account *primitives.AccountData, value sc.U128) sc.Result[sc.Encodable] {
	actual := sc.Min128(account.Reserved, value)
	account.Reserved = account.Reserved.Sub(actual)

	account.Free = sc.SaturatingAddU128(account.Free, actual)

	return sc.Result[sc.Encodable]{
		HasError: false,
		Value:    actual,
	}
}
