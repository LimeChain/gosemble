package dispatchables

import (
	"math/big"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/frame/balances/errors"
	"github.com/LimeChain/gosemble/frame/balances/events"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
)

// TODO: Remove once withdraw is attached to a struct
var (
	existentialDeposit = 1 * constants.Dollar
	ExistentialDeposit = big.NewInt(0).SetUint64(existentialDeposit)
)

// Withdraw withdraws `value` free balance from `who`, respecting existence requirements.
// Does not do anything if value is 0.
// TODO: Refactor
func Withdraw(who types.Address32, value sc.U128, reasons sc.U8, liveness types.ExistenceRequirement) (types.Balance, types.DispatchError) {
	if value.ToBigInt().Cmp(constants.Zero) == 0 {
		return sc.NewU128FromUint64(uint64(0)), nil
	}

	result := tryMutateAccount(who, func(account *types.AccountData, _ bool) sc.Result[sc.Encodable] {
		newFromAccountFree := new(big.Int).Sub(account.Free.ToBigInt(), value.ToBigInt())

		if newFromAccountFree.Cmp(constants.Zero) < 0 {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value: types.NewDispatchErrorModule(types.CustomModuleError{
					Index:   balances.ModuleIndex,
					Error:   sc.U32(errors.ErrorInsufficientBalance),
					Message: sc.NewOption[sc.Str](nil),
				}),
			}
		}

		existentialDeposit := ExistentialDeposit
		sumNewFreeReserved := new(big.Int).Add(newFromAccountFree, account.Reserved.ToBigInt())
		sumFreeReserved := new(big.Int).Add(account.Free.ToBigInt(), account.Reserved.ToBigInt())

		wouldBeDead := sumNewFreeReserved.Cmp(existentialDeposit) < 0
		wouldKill := wouldBeDead && (sumFreeReserved.Cmp(existentialDeposit) >= 0)

		if !(liveness == types.ExistenceRequirementAllowDeath || !wouldKill) {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value: types.NewDispatchErrorModule(types.CustomModuleError{
					Index:   balances.ModuleIndex,
					Error:   sc.U32(errors.ErrorKeepAlive),
					Message: sc.NewOption[sc.Str](nil),
				}),
			}
		}

		err := ensureCanWithdraw(who, value.ToBigInt(), types.Reasons(reasons), newFromAccountFree)
		if err != nil {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value:    err,
			}
		}

		account.Free = sc.NewU128FromBigInt(newFromAccountFree)

		system.DepositEvent(events.NewEventWithdraw(who.FixedSequence, value))

		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    value,
		}
	})

	if result.HasError {
		return types.Balance{}, result.Value.(types.DispatchError)
	}

	return value, nil
}

// ensureCanWithdraw checks that an account can withdraw from their balance given any existing withdraw restrictions.
func ensureCanWithdraw(who types.Address32, amount *big.Int, reasons types.Reasons, newBalance *big.Int) types.DispatchError {
	if amount.Cmp(constants.Zero) == 0 {
		return nil
	}

	accountInfo := system.StorageGetAccount(who.FixedSequence)
	minBalance := accountInfo.Frozen(reasons)
	if minBalance.Cmp(newBalance) > 0 {
		return types.NewDispatchErrorModule(types.CustomModuleError{
			Index:   balances.ModuleIndex,
			Error:   sc.U32(errors.ErrorLiquidityRestrictions),
			Message: sc.NewOption[sc.Str](nil),
		})
	}

	return nil
}

// tryMutateAccount mutates an account based on argument `f`. Does not change total issuance.
// Does not do anything if `f` returns an error.
func tryMutateAccount(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	result := tryMutateAccountWithDust(who, f)
	if result.HasError {
		return result
	}

	r := result.Value.(sc.VaryingData)

	// TODO: Convert this to an Option and uncomment it.
	// Check Substrate implementation for reference.
	//dustCleaner := r[1].(DustCleanerValue)
	//dustCleaner.Drop()

	return sc.Result[sc.Encodable]{HasError: false, Value: r[0].(sc.Encodable)}
}

func tryMutateAccountWithDust(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	result := system.TryMutateExists(who, func(maybeAccount *types.AccountData) sc.Result[sc.Encodable] {
		account := &types.AccountData{}
		isNew := true
		if !reflect.DeepEqual(maybeAccount, types.AccountData{}) {
			account = maybeAccount
			isNew = false
		}

		result := f(account, isNew)
		if result.HasError {
			return result
		}

		maybeEndowed := sc.NewOption[types.Balance](nil)
		if isNew {
			maybeEndowed = sc.NewOption[types.Balance](account.Free)
		}
		maybeAccountWithDust, imbalance := postMutation(*account)
		if !maybeAccountWithDust.HasValue {
			maybeAccount = &types.AccountData{}
		} else {
			maybeAccount.Free = maybeAccountWithDust.Value.Free
			maybeAccount.MiscFrozen = maybeAccountWithDust.Value.MiscFrozen
			maybeAccount.FeeFrozen = maybeAccountWithDust.Value.FeeFrozen
			maybeAccount.Reserved = maybeAccountWithDust.Value.Reserved
		}

		r := sc.NewVaryingData(maybeEndowed, imbalance, result)

		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    r,
		}
	})
	if result.HasError {
		return result
	}

	resultValue := result.Value.(sc.VaryingData)
	maybeEndowed := resultValue[0].(sc.Option[types.Balance])
	if maybeEndowed.HasValue {
		system.DepositEvent(events.NewEventEndowed(who.FixedSequence, maybeEndowed.Value))
	}
	maybeDust := resultValue[1].(sc.Option[NegativeImbalance])
	dustCleaner := DustCleanerValue{
		AccountId:         who,
		NegativeImbalance: maybeDust.Value,
	}

	r := sc.NewVaryingData(resultValue[2], dustCleaner)

	return sc.Result[sc.Encodable]{HasError: false, Value: r}
}

func postMutation(
	new types.AccountData) (sc.Option[types.AccountData], sc.Option[NegativeImbalance]) {
	total := new.Total()

	if total.Cmp(ExistentialDeposit) < 0 {
		if total.Cmp(constants.Zero) == 0 {
			return sc.NewOption[types.AccountData](nil), sc.NewOption[NegativeImbalance](nil)
		} else {
			return sc.NewOption[types.AccountData](nil), sc.NewOption[NegativeImbalance](NewNegativeImbalance(sc.NewU128FromBigInt(total)))
		}
	}

	return sc.NewOption[types.AccountData](new), sc.NewOption[NegativeImbalance](nil)
}
