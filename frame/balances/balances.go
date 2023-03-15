package balances

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/LimeChain/gosemble/primitives/log"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/frame/system"
	"github.com/LimeChain/gosemble/primitives/support"
	"github.com/LimeChain/gosemble/primitives/types"
)

var Module = support.ModuleMetadata{
	Index: balances.ModuleIndex,
	Functions: map[sc.U8]support.FunctionMetadata{
		balances.FunctionTransferIndex:          {Func: Transfer},
		balances.FunctionSetBalanceIndex:        {Func: SetBalance},
		balances.FunctionForceTransferIndex:     {Func: ForceTransfer},
		balances.FunctionTransferKeepAliveIndex: {Func: TransferKeepAlive},
		balances.FunctionTransferAllIndex:       {Func: TransferAll},
		balances.FunctionForceFreeIndex:         {Func: ForceFree},
	},
}

var (
	existentialDeposit = big.NewInt(0).SetUint64(balances.ExistentialDeposit)
)

func Transfer(origin types.RawOrigin, dest types.MultiAddress, value sc.U128) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	if origin[0] != types.RawOriginSigned {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				PostInfo: types.PostDispatchInfo{
					ActualWeight: sc.Option[types.Weight]{
						HasValue: true,
						Value:    types.WeightFromParts(38_109_000, 0),
					},
					PaysFee: 0,
				},
				Error: types.NewDispatchErrorBadOrigin(),
			},
		}
	}

	to, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				// TODO: weight
				PostInfo: types.PostDispatchInfo{
					ActualWeight: sc.Option[types.Weight]{
						HasValue: true,
						Value:    types.WeightFromParts(38_109_000, 0),
					},
					PaysFee: 0,
				},
				Error: types.NewDispatchErrorCannotLookup(),
			},
		}
	}

	transactor := origin[1].(types.Address32)

	e := transfer(transactor, to, value, types.ExistenceRequirementAllowDeath)
	if e != nil {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				PostInfo: types.PostDispatchInfo{
					ActualWeight: sc.Option[types.Weight]{
						HasValue: false,
					},
					PaysFee: 0,
				},
				Error: e,
			},
		}
	}

	return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		HasError: false,
	}
}

func SetBalance(origin types.RawOrigin, who types.MultiAddress, newFree *big.Int, newReserved *big.Int) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	if origin[0] != types.RawOriginRoot {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				PostInfo: types.PostDispatchInfo{
					ActualWeight: sc.Option[types.Weight]{
						HasValue: true,
						Value:    types.WeightFromParts(38_109_000, 0),
					},
					PaysFee: 0,
				},
				Error: types.NewDispatchErrorBadOrigin(),
			},
		}
	}

	address, err := types.DefaultAccountIdLookup().Lookup(who)
	if err != nil {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				PostInfo: types.PostDispatchInfo{
					ActualWeight: sc.Option[types.Weight]{
						// TODO: weight
						HasValue: true,
						Value:    types.WeightFromParts(38_109_000, 0),
					},
					PaysFee: 0,
				},
				Error: types.NewDispatchErrorCannotLookup(),
			},
		}
	}

	existentialDeposit := big.NewInt(0).SetUint64(balances.ExistentialDeposit)
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
		NewEventBalanceSet(who.AsAddress32().FixedSequence, sc.NewU128FromBigInt(newFree), sc.NewU128FromBigInt(newReserved)),
	)
	return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{}
}

func ForceTransfer(origin types.RawOrigin, source types.MultiAddress, dest types.MultiAddress, value sc.U128) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	if origin[0] != types.RawOriginRoot {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				PostInfo: types.PostDispatchInfo{
					ActualWeight: sc.Option[types.Weight]{
						HasValue: true,
						Value:    types.WeightFromParts(38_109_000, 0),
					},
					PaysFee: 0,
				},
				Error: types.NewDispatchErrorBadOrigin(),
			},
		}
	}

	sourceAddress, err := types.DefaultAccountIdLookup().Lookup(source)
	if err != nil {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				// TODO: weight
				PostInfo: types.PostDispatchInfo{
					ActualWeight: sc.Option[types.Weight]{
						HasValue: true,
						Value:    types.WeightFromParts(38_109_000, 0),
					},
					PaysFee: 0,
				},
				Error: types.NewDispatchErrorCannotLookup(),
			},
		}
	}
	destinationAddress, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				// TODO: weight
				PostInfo: types.PostDispatchInfo{
					ActualWeight: sc.Option[types.Weight]{
						HasValue: true,
						Value:    types.WeightFromParts(38_109_000, 0),
					},
					PaysFee: 0,
				},
				Error: types.NewDispatchErrorCannotLookup(),
			},
		}
	}

	e := transfer(sourceAddress, destinationAddress, value, types.ExistenceRequirementAllowDeath)
	if e != nil {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				PostInfo: types.PostDispatchInfo{
					ActualWeight: sc.Option[types.Weight]{
						HasValue: false,
					},
					PaysFee: 0,
				},
				Error: e,
			},
		}
	}

	return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		HasError: false,
	}
}

func TransferKeepAlive(origin types.RawOrigin, dest types.MultiAddress, value sc.U128) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	if origin[0] != types.RawOriginSigned {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				PostInfo: types.PostDispatchInfo{
					ActualWeight: sc.Option[types.Weight]{
						HasValue: false,
					},
					PaysFee: 0,
				},
				Error: types.NewDispatchErrorBadOrigin(),
			},
		}
	}
	transactor := origin[1].(types.Address32)

	address, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				PostInfo: types.PostDispatchInfo{
					ActualWeight: sc.Option[types.Weight]{
						HasValue: false,
					},
					PaysFee: 0,
				},
				Error: types.NewDispatchErrorCannotLookup(),
			},
		}
	}

	e := transfer(transactor, address, value, types.ExistenceRequirementKeepAlive)
	if e != nil {
		return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
			HasError: true,
			Err: types.DispatchErrorWithPostInfo[types.PostDispatchInfo]{
				PostInfo: types.PostDispatchInfo{
					ActualWeight: sc.Option[types.Weight]{
						HasValue: false,
					},
					PaysFee: 0,
				},
				Error: e,
			},
		}
	}

	return types.DispatchResultWithPostInfo[types.PostDispatchInfo]{
		HasError: false,
	}
}

func TransferAll(origin types.RawOrigin, dest types.MultiAddress, keepAlive bool) types.DispatchError {
	if origin[0] != types.RawOriginSigned {
		return types.NewDispatchErrorBadOrigin()
	}

	transactor := origin[1].(types.Address32)
	reducibleBalance := reducibleBalance(transactor, keepAlive)

	to, err := types.DefaultAccountIdLookup().Lookup(dest)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to lookup [%s]", dest.Bytes()))
		return types.NewDispatchErrorCannotLookup()
	}

	keep := types.ExistenceRequirementKeepAlive
	if !keepAlive {
		keep = types.ExistenceRequirementAllowDeath
	}

	return transfer(transactor, to, reducibleBalance, keep)
}

// ForceFree
// Consider Substrate fn force_unreserve
func ForceFree(origin types.RawOrigin, who types.MultiAddress, amount *big.Int) types.DispatchError {
	if origin[0] != types.RawOriginRoot {
		return types.NewDispatchErrorBadOrigin()
	}

	target, err := types.DefaultAccountIdLookup().Lookup(who)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to lookup [%s]", who.Bytes()))
		return types.NewDispatchErrorCannotLookup()
	}

	forceFree(target, amount)

	return nil
}

func transfer(from types.Address32, to types.Address32, value sc.U128, existenceRequirement types.ExistenceRequirement) types.DispatchError {
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
						Error:   sc.U32(ErrorLiquidityRestrictions),
						Message: sc.NewOption[sc.Str](nil),
					}),
				}
			}
			fromAccount.Free = sc.NewU128FromBigInt(newFromAccountFree)

			newToAccountFree := new(big.Int).Add(toAccount.Free.ToBigInt(), value.ToBigInt())
			toAccount.Free = sc.NewU128FromBigInt(newToAccountFree)

			existentialDeposit := big.NewInt(0).SetUint64(balances.ExistentialDeposit)
			if toAccount.Total().Cmp(existentialDeposit) < 0 {
				return sc.Result[sc.Encodable]{
					HasError: true,
					Value: types.NewDispatchErrorModule(types.CustomModuleError{
						Index:   balances.ModuleIndex,
						Error:   sc.U32(ErrorExistentialDeposit),
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
						Error:   sc.U32(ErrorKeepAlive),
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

	system.DepositEvent(NewEventTransfer(from.FixedSequence, to.FixedSequence, value))
	return nil
}

// forceFree frees some funds, returning the amount that has not been freed.
func forceFree(who types.Address32, value *big.Int) *big.Int {
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

	system.DepositEvent(NewEventUnreserved(who.FixedSequence, actual))

	return new(big.Int).Sub(value, actual.ToBigInt())
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
			Error:   sc.U32(ErrorLiquidityRestrictions),
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
		system.DepositEvent(NewEventEndowed(who.FixedSequence, maybeEndowed.Value))
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

	if total.Cmp(existentialDeposit) < 0 {
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

	existentialDeposit := big.NewInt(0).SetUint64(balances.ExistentialDeposit)
	diff := new(big.Int).Sub(accountData.Total(), liquid)

	mustRemainToExist := new(big.Int).Sub(existentialDeposit, diff)

	result := new(big.Int).Sub(liquid, mustRemainToExist)
	if result.Cmp(liquid) > 0 {
		return sc.NewU128FromBigInt(big.NewInt(0))
	}

	return sc.NewU128FromBigInt(result)
}
