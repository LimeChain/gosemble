package dispatchables

import (
	"fmt"
	"math/big"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	bc "github.com/LimeChain/gosemble/frame/balances/constants"
	"github.com/LimeChain/gosemble/frame/balances/events"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

type FnForceFree struct{}

func (_ FnForceFree) Index() sc.U8 {
	return bc.FunctionForceFreeIndex
}

func (_ FnForceFree) BaseWeight(b ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `206`
	//  Estimated: `3593`
	// Minimum execution time: 16_790 nanoseconds.
	r := constants.DbWeight.Reads(1)
	w := constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(17_029_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ FnForceFree) WeightInfo(baseWeight types.Weight, target []byte) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ FnForceFree) ClassifyDispatch(baseWeight types.Weight, target []byte) types.DispatchClass {
	return types.NewDispatchClassMandatory()
}

func (_ FnForceFree) PaysFee(baseWeight types.Weight, target []byte) types.Pays {
	return types.NewPaysYes()
}

func (fn FnForceFree) Dispatch(origin types.RuntimeOrigin, who types.MultiAddress, amount *big.Int) (ok sc.Empty, err types.DispatchError) {
	return forceFree(origin, who, amount)
}

// ForceFree
// Consider Substrate fn force_unreserve
func forceFree(origin types.RawOrigin, who types.MultiAddress, amount *big.Int) (sc.Empty, types.DispatchError) {
	if !origin.IsRootOrigin() {
		return sc.Empty{}, types.NewDispatchErrorBadOrigin()
	}

	target, err := types.DefaultAccountIdLookup().Lookup(who)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to lookup [%s]", who.Bytes()))
		return sc.Empty{}, types.NewDispatchErrorCannotLookup()
	}

	force(target, amount)

	return sc.Empty{}, nil
}

// forceFree frees some funds, returning the amount that has not been freed.
func force(who types.Address32, value *big.Int) *big.Int {
	if value.Cmp(constants.Zero) == 0 {
		return big.NewInt(0)
	}

	if totalBalance(who).Cmp(constants.Zero) == 0 {
		return value
	}

	result := system.Mutate(who, func(accountData *types.AccountInfo) sc.Result[sc.Encodable] {
		actual := accountData.Data.Reserved.ToBigInt()
		if value.Cmp(actual) < 0 {
			actual = value
		}

		newReserved := new(big.Int).Sub(accountData.Data.Reserved.ToBigInt(), actual)
		accountData.Data.Reserved = sc.NewU128FromBigInt(newReserved)

		// TODO: defensive_saturating_add
		newFree := new(big.Int).Add(accountData.Data.Free.ToBigInt(), actual)
		accountData.Data.Free = sc.NewU128FromBigInt(newFree)

		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    sc.NewU128FromBigInt(actual),
		}
	})

	actual := result.Value.(sc.U128)

	if result.HasError {
		return value
	}

	system.DepositEvent(events.NewEventUnreserved(who.FixedSequence, actual))

	return new(big.Int).Sub(value, actual.ToBigInt())
}
