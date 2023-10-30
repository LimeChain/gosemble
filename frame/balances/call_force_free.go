package balances

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type callForceFree struct {
	primitives.Callable
	transfer
}

func newCallForceFree(moduleId sc.U8, functionId sc.U8, storedMap primitives.StoredMap, constants *consts, mutator accountMutator) primitives.Call {
	call := callForceFree{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		transfer: newTransfer(moduleId, storedMap, constants, mutator),
	}

	return call
}

func (c callForceFree) DecodeArgs(buffer *bytes.Buffer) (primitives.Call, error) {
	who, err := types.DecodeMultiAddress(buffer)
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

func (c callForceFree) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c callForceFree) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callForceFree) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callForceFree) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callForceFree) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (c callForceFree) BaseWeight() types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `206`
	//  Estimated: `3593`
	// Minimum execution time: 16_790 nanoseconds.
	r := c.constants.DbWeight.Reads(1)
	w := c.constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(17_029_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ callForceFree) WeighData(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callForceFree) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ callForceFree) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (c callForceFree) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	amount := args[1].(sc.U128)

	err := c.forceFree(origin, args[0].(types.MultiAddress), amount)
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

// forceFree frees some balance from a user by force.
// Can only be called by ROOT.
// Consider Substrate fn force_unreserve
func (c callForceFree) forceFree(origin types.RawOrigin, who types.MultiAddress, amount sc.U128) types.DispatchError {
	if !origin.IsRootOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	target, err := types.DefaultAccountIdLookup().Lookup(who)
	if err != nil {
		// TODO: there is an issue with fmt.Sprintf when compiled with the "custom gc"
		// log.Debug(fmt.Sprintf("Failed to lookup [%s]", who.Bytes()))
		log.Debug("Failed to lookup [" + string(who.Bytes()) + "]")
		return types.NewDispatchErrorCannotLookup()
	}

	c.force(target, amount)

	return nil
}

// forceFree frees funds, returning the amount that has not been freed.
func (c callForceFree) force(who types.Address32, value sc.U128) (sc.U128, error) {
	if value.Eq(constants.Zero) {
		return constants.Zero, nil
	}

	seq, err := c.storedMap.Get(who.FixedSequence)
	if err != nil {
		return sc.U128{}, err
	}

	totalBalance := seq.Data.Total()
	if totalBalance.Eq(constants.Zero) {
		return value, nil
	}

	result := c.accountMutator.tryMutateAccount(
		who,
		func(account *types.AccountData, _ bool) sc.Result[sc.Encodable] {
			return removeReserveAndFree(account, value)
		},
	)

	if result.HasError {
		return value, nil
	}

	actual := result.Value.(sc.U128)
	c.storedMap.DepositEvent(newEventUnreserved(c.ModuleId, who.FixedSequence, actual))

	return value.Sub(actual), nil
}

// removeReserveAndFree frees reserved value from the account.
func removeReserveAndFree(account *types.AccountData, value sc.U128) sc.Result[sc.Encodable] {
	actual := sc.Min128(account.Reserved, value)
	account.Reserved = account.Reserved.Sub(actual)

	// TODO: defensive_saturating_add
	account.Free = account.Free.Add(actual)

	return sc.Result[sc.Encodable]{
		HasError: false,
		Value:    actual,
	}
}
