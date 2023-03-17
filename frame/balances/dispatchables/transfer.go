package dispatchables

import (
	"math/big"
	"reflect"

	"github.com/LimeChain/gosemble/constants/balances"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/balances/errors"
	"github.com/LimeChain/gosemble/frame/balances/events"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/types"
)

type FnTransfer struct{}

func (_ FnTransfer) Index() sc.U8 {
	return balances.FunctionTransferIndex
}

func (_ FnTransfer) BaseWeight(b ...any) types.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `0`
	//  Estimated: `3593`
	// Minimum execution time: 37_815 nanoseconds.
	r := constants.DbWeight.Reads(1)
	w := constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(38_109_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ FnTransfer) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ FnTransfer) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassMandatory()
}

func (_ FnTransfer) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (fn FnTransfer) Dispatch(origin types.RuntimeOrigin, args ...sc.Encodable) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	err := transfer(origin, args[0].(types.MultiAddress), args[1].(sc.U128))
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

func transfer(origin types.RawOrigin, dest types.MultiAddress, value sc.U128) types.DispatchError {
	if !origin.IsSignedOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	to, e := types.DefaultAccountIdLookup().Lookup(dest)
	if e != nil {
		return types.NewDispatchErrorCannotLookup()
	}

	transactor := origin.AsSigned()

	return trans(transactor, to, value, types.ExistenceRequirementAllowDeath)
}

func trans(from types.Address32, to types.Address32, value sc.U128, existenceRequirement types.ExistenceRequirement) types.DispatchError {
	bnInt := value.ToBigInt()
	if bnInt.Cmp(constants.Zero) == 0 || reflect.DeepEqual(from, to) {
		return nil
	}

	result := tryMutateAccountWithDust(to, func(toAccount *types.AccountData, _ bool) sc.Result[sc.Encodable] {
		return tryMutateAccountWithDust(from, func(fromAccount *types.AccountData, _ bool) sc.Result[sc.Encodable] {
			newFromAccountFree := new(big.Int).Sub(fromAccount.Free.ToBigInt(), value.ToBigInt())

			if newFromAccountFree.Cmp(constants.Zero) < 0 {
				return sc.Result[sc.Encodable]{
					HasError: true,
					Value: types.NewDispatchErrorModule(types.CustomModuleError{
						Index:   balances.ModuleIndex,
						Error:   sc.U32(errors.ErrorLiquidityRestrictions),
						Message: sc.NewOption[sc.Str](nil),
					}),
				}
			}
			fromAccount.Free = sc.NewU128FromBigInt(newFromAccountFree)

			newToAccountFree := new(big.Int).Add(toAccount.Free.ToBigInt(), value.ToBigInt())
			toAccount.Free = sc.NewU128FromBigInt(newToAccountFree)

			existentialDeposit := balances.ExistentialDeposit
			if toAccount.Total().Cmp(existentialDeposit) < 0 {
				return sc.Result[sc.Encodable]{
					HasError: true,
					Value: types.NewDispatchErrorModule(types.CustomModuleError{
						Index:   balances.ModuleIndex,
						Error:   sc.U32(errors.ErrorExistentialDeposit),
						Message: sc.NewOption[sc.Str](nil),
					}),
				}
			}

			dispatchResult := ensureCanWithdraw(from, value.ToBigInt(), types.ReasonsAll, fromAccount.Free.ToBigInt())
			if !reflect.DeepEqual(dispatchResult[0], sc.Empty{}) {
				return sc.Result[sc.Encodable]{
					HasError: true,
					Value:    dispatchResult,
				}
			}

			allowDeath := existenceRequirement == types.ExistenceRequirementAllowDeath
			allowDeath = allowDeath && system.CanDecProviders(from)

			if !(allowDeath || fromAccount.Total().Cmp(existentialDeposit) > 0) {
				return sc.Result[sc.Encodable]{
					HasError: true,
					Value: types.NewDispatchErrorModule(types.CustomModuleError{
						Index:   balances.ModuleIndex,
						Error:   sc.U32(errors.ErrorKeepAlive),
						Message: sc.NewOption[sc.Str](nil),
					}),
				}
			}

			return sc.Result[sc.Encodable]{}
		})
	})

	if result.HasError {
		return result.Value.(types.DispatchError)
	}

	system.DepositEvent(events.NewEventTransfer(from.FixedSequence, to.FixedSequence, value))
	return nil
}

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

func mutateAccount(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	return tryMutateAccount(who, f)
}

func tryMutateAccount(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	result := tryMutateAccountWithDust(who, f)

	r := result.Value.(sc.VaryingData)

	dustCleaner := r[1].(DustCleanerValue)
	dustCleaner.Drop()

	return sc.Result[sc.Encodable]{HasError: false, Value: r[0].(sc.Encodable)}
}

func tryMutateAccountWithDust(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	result := system.TryMutateExists(who, func(maybeAccount *types.AccountData) sc.Result[sc.Encodable] {
		account := &types.AccountData{}
		isNew := true
		if maybeAccount != nil {
			account = maybeAccount
			isNew = false
		}

		result := f(account, isNew)

		maybeEndowed := sc.NewOption[types.Balance](nil)
		if isNew {
			maybeEndowed = sc.NewOption[types.Balance](account.Free)
		}
		maybeAccountWithDust, imbalance := postMutation(*account)
		maybeAccount = &maybeAccountWithDust.Value

		r := sc.NewVaryingData(maybeEndowed, imbalance, result)

		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    r,
		}
	})
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

	if total.Cmp(balances.ExistentialDeposit) < 0 {
		if total.Cmp(constants.Zero) == 0 {
			return sc.NewOption[types.AccountData](nil), sc.NewOption[NegativeImbalance](nil)
		} else {
			return sc.NewOption[types.AccountData](nil), sc.NewOption[NegativeImbalance](NewNegativeImbalance(sc.NewU128FromBigInt(total)))
		}
	}

	return sc.NewOption[types.AccountData](new), sc.NewOption[NegativeImbalance](nil)
}

func totalBalance(who types.Address32) *big.Int {
	return system.StorageGetAccount(who.FixedSequence).Data.Total()
}

func reducibleBalance(who types.Address32, keepAlive bool) types.Balance {
	accountData := system.StorageGetAccount(who.FixedSequence).Data

	lockedOrFrozen := accountData.FeeFrozen
	if accountData.FeeFrozen.ToBigInt().Cmp(accountData.MiscFrozen.ToBigInt()) < 0 {
		lockedOrFrozen = accountData.MiscFrozen
	}

	liquid := new(big.Int).Sub(accountData.Free.ToBigInt(), lockedOrFrozen.ToBigInt())
	if liquid.Cmp(accountData.Free.ToBigInt()) > 0 {
		liquid = big.NewInt(0)
	}

	if system.CanDecProviders(who) && !keepAlive {
		return sc.NewU128FromBigInt(liquid)
	}

	existentialDeposit := balances.ExistentialDeposit
	diff := new(big.Int).Sub(accountData.Total(), liquid)

	mustRemainToExist := new(big.Int).Sub(existentialDeposit, diff)

	result := new(big.Int).Sub(liquid, mustRemainToExist)
	if result.Cmp(liquid) > 0 {
		return sc.NewU128FromBigInt(big.NewInt(0))
	}

	return sc.NewU128FromBigInt(result)
}
