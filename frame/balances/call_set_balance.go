package balances

import (
	"bytes"
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type callSetBalance struct {
	primitives.Callable
	constants      *consts
	storedMap      primitives.StoredMap
	accountMutator accountMutator
}

func newCallSetBalance(moduleId sc.U8, functionId sc.U8, storedMap primitives.StoredMap, constants *consts, mutator accountMutator) primitives.Call {
	call := callSetBalance{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		constants:      constants,
		storedMap:      storedMap,
		accountMutator: mutator,
	}

	return call
}

func (c callSetBalance) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(
		types.DecodeMultiAddress(buffer),
		sc.DecodeCompact(buffer),
		sc.DecodeCompact(buffer),
	)
	return c
}

func (c callSetBalance) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c callSetBalance) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callSetBalance) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callSetBalance) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callSetBalance) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (c callSetBalance) BaseWeight() types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `206`
	//  Estimated: `3593`
	// Minimum execution time: 17_474 nanoseconds.
	r := c.constants.DbWeight.Reads(1)
	w := c.constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(17_777_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ callSetBalance) IsInherent() bool {
	return false
}

func (_ callSetBalance) WeighData(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callSetBalance) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ callSetBalance) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (c callSetBalance) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	newFree := sc.U128(args[1].(sc.Compact))
	newReserved := sc.U128(args[2].(sc.Compact))

	err := c.setBalance(origin, args[0].(types.MultiAddress), newFree, newReserved)
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

// setBalance sets the balance of a given account.
// Changes free and reserve balance of `who`,
// including the total issuance.
// Can only be called by ROOT.
func (c callSetBalance) setBalance(origin types.RawOrigin, who types.MultiAddress, newFree sc.U128, newReserved sc.U128) types.DispatchError {
	if !origin.IsRootOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	address, err := types.DefaultAccountIdLookup().Lookup(who)
	if err != nil {
		return types.NewDispatchErrorCannotLookup()
	}

	existentialDeposit := sc.NewU128FromBigInt(c.constants.ExistentialDeposit)
	sum := newFree.Add(newReserved)

	if sum.Lt(existentialDeposit) {
		newFree = sc.NewU128FromBigInt(big.NewInt(0))
		newReserved = sc.NewU128FromBigInt(big.NewInt(0))
	}

	result := c.accountMutator.tryMutateAccount(address, func(acc *types.AccountData, bool bool) sc.Result[sc.Encodable] {
		oldFree := acc.Free
		oldReserved := acc.Reserved

		acc.Free = newFree
		acc.Reserved = newReserved

		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    sc.NewVaryingData(oldFree, oldReserved),
		}
	})
	parsedResult := result.Value.(sc.VaryingData)
	oldFree := parsedResult[0].(types.Balance)
	oldReserved := parsedResult[1].(types.Balance)

	if newFree.Gt(oldFree) {
		diff := newFree.Sub(oldFree).(sc.U128)

		newPositiveImbalance(diff).Drop()
	} else if newFree.Lt(oldFree) {
		diff := oldFree.Sub(newFree).(sc.U128)

		newNegativeImbalance(diff).Drop()
	}

	if newReserved.Gt(oldReserved) {
		diff := newReserved.Sub(oldReserved).(sc.U128)

		newPositiveImbalance(diff).Drop()
	} else if newReserved.Lt(oldReserved) {
		diff := oldReserved.Sub(newReserved).(sc.U128)

		newNegativeImbalance(diff).Drop()
	}

	c.storedMap.DepositEvent(
		newEventBalanceSet(
			c.ModuleId,
			who.AsAddress32().FixedSequence,
			newFree,
			newReserved,
		),
	)
	return nil
}
