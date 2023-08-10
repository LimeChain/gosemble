package module

import (
	"math/big"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/balances"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/frame/balances/errors"
	"github.com/LimeChain/gosemble/frame/balances/events"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	functionTransferIndex          = 0
	functionSetBalanceIndex        = 1
	functionForceTransferIndex     = 2
	functionTransferKeepAliveIndex = 3
	functionTransferAllIndex       = 4
	functionForceFreeIndex         = 5
)

type BalancesModule struct {
	Index     sc.U8
	Config    *Config
	Constants *consts
	functions map[sc.U8]primitives.Call
}

func NewBalancesModule(index sc.U8, config *Config) BalancesModule {
	constants := newConstants(config.MaxLocks, config.MaxReserves, config.ExistentialDeposit)

	balancesModule := BalancesModule{
		Index:     index,
		Config:    config,
		Constants: constants,
	}
	functions := make(map[sc.U8]primitives.Call)
	functions[functionTransferIndex] = newTransferCall(index, functionTransferIndex, config.StoredMap, constants, balancesModule)
	functions[functionSetBalanceIndex] = newSetBalanceCall(index, functionSetBalanceIndex, config.StoredMap, constants, balancesModule)
	functions[functionForceTransferIndex] = newForceTransferCall(index, functionForceTransferIndex, config.StoredMap, constants, balancesModule)
	functions[functionTransferKeepAliveIndex] = newTransferKeepAliveCall(index, functionTransferKeepAliveIndex, config.StoredMap, constants, balancesModule)
	functions[functionTransferAllIndex] = newTransferAllCall(index, functionTransferAllIndex, config.StoredMap, constants, balancesModule)
	functions[functionForceFreeIndex] = newForceFreeCall(index, functionForceFreeIndex, config.StoredMap)

	balancesModule.functions = functions

	return balancesModule
}

func (bm BalancesModule) Functions() map[sc.U8]primitives.Call {
	return bm.functions
}

func (bm BalancesModule) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (bm BalancesModule) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

// DepositIntoExisting deposits `value` into the free balance of an existing target account `who`.
// If `value` is 0, it does nothing.
func (bm BalancesModule) DepositIntoExisting(who primitives.Address32, value sc.U128) (primitives.Balance, primitives.DispatchError) {
	if value.ToBigInt().Cmp(constants.Zero) == 0 {
		return sc.NewU128FromUint64(uint64(0)), nil
	}

	result := bm.tryMutateAccount(who, func(from *primitives.AccountData, isNew bool) sc.Result[sc.Encodable] {
		if isNew {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
					Index:   balances.ModuleIndex,
					Error:   sc.U32(errors.ErrorDeadAccount),
					Message: sc.NewOption[sc.Str](nil),
				}),
			}
		}

		sum := new(big.Int).Add(from.Free.ToBigInt(), value.ToBigInt())

		from.Free = sc.NewU128FromBigInt(sum)

		bm.Config.StoredMap.DepositEvent(events.NewEventDeposit(who.FixedSequence, value))

		return sc.Result[sc.Encodable]{}
	})

	if result.HasError {
		return primitives.Balance{}, result.Value.(primitives.DispatchError)
	}

	return value, nil
}

// Withdraw withdraws `value` free balance from `who`, respecting existence requirements.
// Does not do anything if value is 0.
func (bm BalancesModule) Withdraw(who primitives.Address32, value sc.U128, reasons sc.U8, liveness primitives.ExistenceRequirement) (primitives.Balance, primitives.DispatchError) {
	if value.ToBigInt().Cmp(constants.Zero) == 0 {
		return sc.NewU128FromUint64(uint64(0)), nil
	}

	result := bm.tryMutateAccount(who, func(account *primitives.AccountData, _ bool) sc.Result[sc.Encodable] {
		newFromAccountFree := new(big.Int).Sub(account.Free.ToBigInt(), value.ToBigInt())

		if newFromAccountFree.Cmp(constants.Zero) < 0 {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
					Index:   balances.ModuleIndex,
					Error:   sc.U32(errors.ErrorInsufficientBalance),
					Message: sc.NewOption[sc.Str](nil),
				}),
			}
		}

		existentialDeposit := bm.Constants.ExistentialDeposit
		sumNewFreeReserved := new(big.Int).Add(newFromAccountFree, account.Reserved.ToBigInt())
		sumFreeReserved := new(big.Int).Add(account.Free.ToBigInt(), account.Reserved.ToBigInt())

		wouldBeDead := sumNewFreeReserved.Cmp(existentialDeposit) < 0
		wouldKill := wouldBeDead && (sumFreeReserved.Cmp(existentialDeposit) >= 0)

		if !(liveness == primitives.ExistenceRequirementAllowDeath || !wouldKill) {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
					Index:   balances.ModuleIndex,
					Error:   sc.U32(errors.ErrorKeepAlive),
					Message: sc.NewOption[sc.Str](nil),
				}),
			}
		}

		err := bm.ensureCanWithdraw(who, value.ToBigInt(), primitives.Reasons(reasons), newFromAccountFree)
		if err != nil {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value:    err,
			}
		}

		account.Free = sc.NewU128FromBigInt(newFromAccountFree)

		bm.Config.StoredMap.DepositEvent(events.NewEventWithdraw(who.FixedSequence, value))

		return sc.Result[sc.Encodable]{
			HasError: false,
			Value:    value,
		}
	})

	if result.HasError {
		return primitives.Balance{}, result.Value.(primitives.DispatchError)
	}

	return value, nil
}

// ensureCanWithdraw checks that an account can withdraw from their balance given any existing withdraw restrictions.
func (bm BalancesModule) ensureCanWithdraw(who primitives.Address32, amount *big.Int, reasons primitives.Reasons, newBalance *big.Int) primitives.DispatchError {
	if amount.Cmp(constants.Zero) == 0 {
		return nil
	}

	accountInfo := bm.Config.StoredMap.Get(who.FixedSequence)
	minBalance := accountInfo.Frozen(reasons)
	if minBalance.Cmp(newBalance) > 0 {
		return primitives.NewDispatchErrorModule(primitives.CustomModuleError{
			Index:   balances.ModuleIndex,
			Error:   sc.U32(errors.ErrorLiquidityRestrictions),
			Message: sc.NewOption[sc.Str](nil),
		})
	}

	return nil
}

// tryMutateAccount mutates an account based on argument `f`. Does not change total issuance.
// Does not do anything if `f` returns an error.
func (bm BalancesModule) tryMutateAccount(who primitives.Address32, f func(who *primitives.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	result := bm.tryMutateAccountWithDust(who, f)
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

func (bm BalancesModule) tryMutateAccountWithDust(who primitives.Address32, f func(who *primitives.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	result := bm.Config.StoredMap.TryMutateExists(who, func(maybeAccount *primitives.AccountData) sc.Result[sc.Encodable] {
		account := &primitives.AccountData{}
		isNew := true
		if !reflect.DeepEqual(maybeAccount, primitives.AccountData{}) {
			account = maybeAccount
			isNew = false
		}

		result := f(account, isNew)
		if result.HasError {
			return result
		}

		maybeEndowed := sc.NewOption[primitives.Balance](nil)
		if isNew {
			maybeEndowed = sc.NewOption[primitives.Balance](account.Free)
		}
		maybeAccountWithDust, imbalance := bm.postMutation(*account)
		if !maybeAccountWithDust.HasValue {
			maybeAccount = &primitives.AccountData{}
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
	maybeEndowed := resultValue[0].(sc.Option[primitives.Balance])
	if maybeEndowed.HasValue {
		bm.Config.StoredMap.DepositEvent(events.NewEventEndowed(who.FixedSequence, maybeEndowed.Value))
	}
	maybeDust := resultValue[1].(sc.Option[negativeImbalance])
	dustCleaner := newDustCleanerValue(who, maybeDust.Value, bm.Config.StoredMap)

	r := sc.NewVaryingData(resultValue[2], dustCleaner)

	return sc.Result[sc.Encodable]{HasError: false, Value: r}
}

func (bm BalancesModule) postMutation(new primitives.AccountData) (sc.Option[primitives.AccountData], sc.Option[negativeImbalance]) {
	total := new.Total()

	if total.Cmp(bm.Constants.ExistentialDeposit) < 0 {
		if total.Cmp(constants.Zero) == 0 {
			return sc.NewOption[primitives.AccountData](nil), sc.NewOption[negativeImbalance](nil)
		} else {
			return sc.NewOption[primitives.AccountData](nil), sc.NewOption[negativeImbalance](newNegativeImbalance(sc.NewU128FromBigInt(total)))
		}
	}

	return sc.NewOption[primitives.AccountData](new), sc.NewOption[negativeImbalance](nil)
}

func (bm BalancesModule) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	return bm.metadataTypes(), primitives.MetadataModule{
		Name: "Balances",
		Storage: sc.NewOption[primitives.MetadataModuleStorage](primitives.MetadataModuleStorage{
			Prefix: "Balances",
			Items: sc.Sequence[primitives.MetadataModuleStorageEntry]{
				primitives.NewMetadataModuleStorageEntry(
					"TotalIssuance",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.PrimitiveTypesU128)),
					"The total units issued in the system."),
				primitives.NewMetadataModuleStorageEntry(
					"InactiveIssuance",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionPlain(sc.ToCompact(metadata.PrimitiveTypesU128)),
					"The total units of outstanding deactivated balance in the system."),
				primitives.NewMetadataModuleStorageEntry(
					"Account",
					primitives.MetadataModuleStorageEntryModifierDefault,
					primitives.NewMetadataModuleStorageEntryDefinitionMap(
						sc.Sequence[primitives.MetadataModuleStorageHashFunc]{primitives.MetadataModuleStorageHashFuncMultiBlake128Concat},
						sc.ToCompact(metadata.TypesAddress32),
						sc.ToCompact(metadata.TypesAccountData)),
					"The Balances pallet example of storing the balance of an account."),
				// TODO: Locks, Reserves, currently not used
			},
		}),
		Call:  sc.NewOption[sc.Compact](sc.ToCompact(metadata.BalancesCalls)),
		Event: sc.NewOption[sc.Compact](sc.ToCompact(metadata.TypesBalancesEvent)),
		Constants: sc.Sequence[primitives.MetadataModuleConstant]{
			primitives.NewMetadataModuleConstant(
				"ExistentialDeposit",
				sc.ToCompact(metadata.PrimitiveTypesU128),
				sc.BytesToSequenceU8(sc.NewU128FromBigInt(bm.Constants.ExistentialDeposit).Bytes()),
				"The minimum amount required to keep an account open. MUST BE GREATER THAN ZERO!",
			),
			primitives.NewMetadataModuleConstant(
				"MaxLocks",
				sc.ToCompact(metadata.PrimitiveTypesU32),
				sc.BytesToSequenceU8(bm.Constants.MaxLocks.Bytes()),
				"The maximum number of locks that should exist on an account.  Not strictly enforced, but used for weight estimation.",
			),
			primitives.NewMetadataModuleConstant(
				"MaxReserves",
				sc.ToCompact(metadata.PrimitiveTypesU32),
				sc.BytesToSequenceU8(bm.Constants.MaxReserves.Bytes()),
				"The maximum number of named reserves that can exist on an account.",
			),
		}, // TODO:
		Error: sc.NewOption[sc.Compact](sc.ToCompact(metadata.TypesBalancesErrors)),
		Index: bm.Index,
	}
}

func (bm BalancesModule) metadataTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithPath(metadata.TypesBalancesEvent, "pallet_balances pallet Event", sc.Sequence[sc.Str]{"pallet_balances", "pallet", "Event"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Endowed",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "account", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "free_balance", "T::Balance"),
					},
					events.EventEndowed,
					"Event.Endowed"),
				primitives.NewMetadataDefinitionVariant(
					"DustLost",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "account", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					events.EventDustLost,
					"Events.DustLost"),
				primitives.NewMetadataDefinitionVariant(
					"Transfer",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "from", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "to", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					events.EventTransfer,
					"Events.Transfer"),
				primitives.NewMetadataDefinitionVariant(
					"BalanceSet",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "free", "T::Balance"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "reserved", "T::Balance"),
					},
					events.EventBalanceSet,
					"Events.BalanceSet"),
				primitives.NewMetadataDefinitionVariant(
					"Reserved",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					events.EventReserved,
					"Events.Reserved"),
				primitives.NewMetadataDefinitionVariant(
					"Unreserved",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					events.EventUnreserved,
					"Events.Unreserved"),
				primitives.NewMetadataDefinitionVariant(
					"ReserveRepatriated",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "from", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "to", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesBalanceStatus, "destination_status", "Status"),
					},
					events.EventReserveRepatriated,
					"Events.ReserveRepatriated"),
				primitives.NewMetadataDefinitionVariant(
					"Deposit",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					events.EventDeposit,
					"Event.Deposit"),
				primitives.NewMetadataDefinitionVariant(
					"Withdraw",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					events.EventWithdraw,
					"Event.Withdraw"),
				primitives.NewMetadataDefinitionVariant(
					"Slashed",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					events.EventSlashed,
					"Event.Slashed"),
			},
		)),
		primitives.NewMetadataTypeWithPath(metadata.TypesBalanceStatus,
			"BalanceStatus",
			sc.Sequence[sc.Str]{"frame_support", "traits", "tokens", "misc", "BalanceStatus"}, primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"Free",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						primitives.BalanceStatusFree,
						"BalanceStatus.Free"),
					primitives.NewMetadataDefinitionVariant(
						"Reserved",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						primitives.BalanceStatusReserved,
						"BalanceStatus.Reserved"),
				})),

		primitives.NewMetadataTypeWithParams(metadata.TypesBalancesErrors,
			"pallet_balances pallet Error",
			sc.Sequence[sc.Str]{"pallet_balances", "pallet", "Error"},
			primitives.NewMetadataTypeDefinitionVariant(
				sc.Sequence[primitives.MetadataDefinitionVariant]{
					primitives.NewMetadataDefinitionVariant(
						"VestingBalance",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						errors.ErrorVestingBalance,
						"Vesting balance too high to send value"),
					primitives.NewMetadataDefinitionVariant(
						"LiquidityRestrictions",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						errors.ErrorLiquidityRestrictions,
						"Account liquidity restrictions prevent withdrawal"),
					primitives.NewMetadataDefinitionVariant(
						"InsufficientBalance",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						errors.ErrorInsufficientBalance,
						"Balance too low to send value."),
					primitives.NewMetadataDefinitionVariant(
						"ExistentialDeposit",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						errors.ErrorExistentialDeposit,
						"Value too low to create account due to existential deposit"),
					primitives.NewMetadataDefinitionVariant(
						"KeepAlive",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						errors.ErrorKeepAlive,
						"Transfer/payment would kill account"),
					primitives.NewMetadataDefinitionVariant(
						"ExistingVestingSchedule",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						errors.ErrorExistingVestingSchedule,
						"A vesting schedule already exists for this account"),
					primitives.NewMetadataDefinitionVariant(
						"DeadAccount",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						errors.ErrorDeadAccount,
						"Beneficiary account must pre-exist"),
					primitives.NewMetadataDefinitionVariant(
						"TooManyReserves",
						sc.Sequence[primitives.MetadataTypeDefinitionField]{},
						errors.ErrorTooManyReserves,
						"Number of named reserves exceed MaxReserves"),
				}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataEmptyTypeParameter("T"),
				primitives.NewMetadataEmptyTypeParameter("I"),
			}),

		primitives.NewMetadataTypeWithParams(metadata.BalancesCalls, "Balances calls", sc.Sequence[sc.Str]{"pallet_balances", "pallet", "Call"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"transfer",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesMultiAddress, "dest", "AccountIdLookupOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU128, "value", "T::Balance"),
					},
					functionTransferIndex,
					"Transfer some liquid free balance to another account."),
				primitives.NewMetadataDefinitionVariant(
					"set_balance",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesMultiAddress, "who", "AccountIdLookupOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU128, "new_free", "T::Balance"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU128, "new_reserved", "T::Balance"),
					},
					functionSetBalanceIndex,
					"Set the balances of a given account."),
				primitives.NewMetadataDefinitionVariant(
					"force_transfer",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesMultiAddress, "source", "AccountIdLookupOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesMultiAddress, "dest", "AccountIdLookupOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU128, "value", "T::Balance"),
					},
					functionForceTransferIndex,
					"Exactly as `transfer`, except the origin must be root and the source account may be specified."),
				primitives.NewMetadataDefinitionVariant(
					"transfer_keep_alive",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesMultiAddress, "dest", "AccountIdLookupOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesCompactU128, "value", "T::Balance"),
					},
					functionTransferKeepAliveIndex,
					"Same as the [`transfer`] call, but with a check that the transfer will not kill the origin account."),
				primitives.NewMetadataDefinitionVariant(
					"transfer_all",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesMultiAddress, "dest", "AccountIdLookupOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesBool, "keep_alive", "bool"),
					},
					functionTransferAllIndex,
					"Transfer the entire transferable balance from the caller account."),
				primitives.NewMetadataDefinitionVariant(
					"force_unreserve",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesMultiAddress, "who", "AccountIdLookupOf<T>"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					functionForceFreeIndex,
					"Unreserve some balance from a user by force."),
			}),
			sc.Sequence[primitives.MetadataTypeParameter]{
				primitives.NewMetadataEmptyTypeParameter("T"),
				primitives.NewMetadataEmptyTypeParameter("I"),
			}),
	}
}
