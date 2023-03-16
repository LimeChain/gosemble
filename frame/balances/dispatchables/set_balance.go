package dispatchables

import (
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	bc "github.com/LimeChain/gosemble/frame/balances/constants"
	"github.com/LimeChain/gosemble/frame/balances/events"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
)

type FnSetBalance struct{}

func (_ FnSetBalance) Index() sc.U8 {
	return bc.FunctionSetBalanceIndex
}

func (_ FnSetBalance) BaseWeight(b ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `206`
	//  Estimated: `3593`
	// Minimum execution time: 17_474 nanoseconds.
	r := constants.DbWeight.Reads(1)
	w := constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(17_777_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ FnSetBalance) WeightInfo(baseWeight types.Weight, target []byte) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ FnSetBalance) ClassifyDispatch(baseWeight types.Weight, target []byte) types.DispatchClass {
	return types.NewDispatchClassMandatory()
}

func (_ FnSetBalance) PaysFee(baseWeight types.Weight, target []byte) types.Pays {
	return types.NewPaysYes()
}

func (fn FnSetBalance) Dispatch(origin types.RuntimeOrigin, who types.MultiAddress, newFree *big.Int, newReserved *big.Int) (ok sc.Empty, err types.DispatchError) {
	return setBalance(origin, who, newFree, newReserved)
}

func setBalance(origin types.RawOrigin, who types.MultiAddress, newFree *big.Int, newReserved *big.Int) (sc.Empty, types.DispatchError) {
	if !origin.IsRootOrigin() {
		return sc.Empty{}, types.NewDispatchErrorBadOrigin()
	}

	address, err := types.DefaultAccountIdLookup().Lookup(who)
	if err != nil {
		return sc.Empty{}, types.NewDispatchErrorCannotLookup()
	}

	existentialDeposit := bc.ExistentialDeposit
	sum := new(big.Int).Add(newFree, newReserved)

	if sum.Cmp(existentialDeposit) < 0 {
		newFree = big.NewInt(0)
		newReserved = big.NewInt(0)
	}

	result := mutateAccount(address, func(acc *types.AccountData, bool bool) sc.Result[sc.Encodable] {
		oldFree := acc.Free
		oldReserved := acc.Reserved

		acc.Free = sc.NewU128FromBigInt(newFree)
		acc.Reserved = sc.NewU128FromBigInt(newReserved)

		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    sc.NewVaryingData(oldFree, oldReserved),
		}
	})
	parsedResult := result.Value.(sc.VaryingData)
	oldFree := parsedResult[0].(types.Balance)
	oldReserved := parsedResult[1].(types.Balance)

	if newFree.Cmp(oldFree.ToBigInt()) > 0 {
		diff := new(big.Int).Sub(newFree, oldFree.ToBigInt())

		NewPositiveImbalance(sc.NewU128FromBigInt(diff)).Drop()
	} else if newFree.Cmp(oldFree.ToBigInt()) < 0 {
		diff := new(big.Int).Sub(oldFree.ToBigInt(), newFree)

		NewNegativeImbalance(sc.NewU128FromBigInt(diff)).Drop()
	}

	if newReserved.Cmp(oldReserved.ToBigInt()) > 0 {
		diff := new(big.Int).Sub(newReserved, oldReserved.ToBigInt())

		NewPositiveImbalance(sc.NewU128FromBigInt(diff)).Drop()
	} else if newReserved.Cmp(oldReserved.ToBigInt()) < 0 {
		diff := new(big.Int).Sub(oldReserved.ToBigInt(), newReserved)

		NewNegativeImbalance(sc.NewU128FromBigInt(diff)).Drop()
	}

	system.DepositEvent(
		events.NewEventBalanceSet(
			who.AsAddress32().FixedSequence,
			sc.NewU128FromBigInt(newFree),
			sc.NewU128FromBigInt(newReserved),
		),
	)
	return sc.Empty{}, nil
}
