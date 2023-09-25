package balances

import (
	"bytes"
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type callForceFree struct {
	primitives.Callable
	dbWeight  primitives.RuntimeDbWeight
	storedMap primitives.StoredMap
}

func newCallForceFree(moduleId sc.U8, functionId sc.U8, dbWeight primitives.RuntimeDbWeight, storedMap primitives.StoredMap) primitives.Call {
	call := callForceFree{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		dbWeight:  dbWeight,
		storedMap: storedMap,
	}

	return call
}

func (c callForceFree) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(
		types.DecodeMultiAddress(buffer),
		sc.DecodeU128(buffer),
	)
	return c
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
	r := c.dbWeight.Reads(1)
	w := c.dbWeight.Writes(1)
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

// forceFree frees some funds, returning the amount that has not been freed.
func (c callForceFree) force(who types.Address32, value sc.U128) sc.U128 {
	if value.Eq(sc.NewU128FromBigInt(constants.Zero)) {
		return sc.NewU128FromBigInt(constants.Zero)
	}

	totalBalance := c.storedMap.Get(who.FixedSequence).Data.Total()
	if totalBalance.Eq(sc.NewU128FromBigInt(constants.Zero)) {
		return value
	}

	result := c.storedMap.Mutate(who, func(accountData *types.AccountInfo) sc.Result[sc.Encodable] {
		actual := accountData.Data.Reserved
		if value.Lt(actual) {
			actual = value
		}

		newReserved := new(big.Int).Sub(accountData.Data.Reserved.ToBigInt(), actual.ToBigInt())
		accountData.Data.Reserved = sc.NewU128FromBigInt(newReserved)

		// TODO: defensive_saturating_add
		newFree := new(big.Int).Add(accountData.Data.Free.ToBigInt(), actual.ToBigInt())
		accountData.Data.Free = sc.NewU128FromBigInt(newFree)

		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    actual,
		}
	})

	actual := result.Value.(sc.U128)

	if result.HasError {
		return value
	}

	c.storedMap.DepositEvent(newEventUnreserved(c.ModuleId, who.FixedSequence, actual))

	return sc.NewU128FromBigInt(new(big.Int).Sub(value.ToBigInt(), actual.ToBigInt()))
}
