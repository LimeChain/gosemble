package balances

import (
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/constants/metadata"
	"github.com/LimeChain/gosemble/frame/balances/errors"
	"github.com/LimeChain/gosemble/hooks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

const (
	functionTransferIndex = iota
	functionSetBalanceIndex
	functionForceTransferIndex
	functionTransferKeepAliveIndex
	functionTransferAllIndex
	functionForceFreeIndex
)

type Module[N sc.Numeric] struct {
	primitives.DefaultProvideInherent
	hooks.DefaultDispatchModule[N]
	Index     sc.U8
	Config    *Config
	Constants *consts
	functions map[sc.U8]primitives.Call
}

func New[N sc.Numeric](index sc.U8, config *Config) Module[N] {
	constants := newConstants(config.DbWeight, config.MaxLocks, config.MaxReserves, config.ExistentialDeposit)

	module := Module[N]{
		Index:     index,
		Config:    config,
		Constants: constants,
	}
	functions := make(map[sc.U8]primitives.Call)
	functions[functionTransferIndex] = newCallTransfer(index, functionTransferIndex, config.StoredMap, constants, module)
	functions[functionSetBalanceIndex] = newCallSetBalance(index, functionSetBalanceIndex, config.StoredMap, constants, module)
	functions[functionForceTransferIndex] = newCallForceTransfer(index, functionForceTransferIndex, config.StoredMap, constants, module)
	functions[functionTransferKeepAliveIndex] = newCallTransferKeepAlive(index, functionTransferKeepAliveIndex, config.StoredMap, constants, module)
	functions[functionTransferAllIndex] = newCallTransferAll(index, functionTransferAllIndex, config.StoredMap, constants, module)
	functions[functionForceFreeIndex] = newCallForceFree(index, functionForceFreeIndex, config.DbWeight, config.StoredMap)

	module.functions = functions

	return module
}

func (m Module[N]) GetIndex() sc.U8 {
	return m.Index
}

func (m Module[N]) name() sc.Str {
	return "Balances"
}

func (m Module[N]) Functions() map[sc.U8]primitives.Call {
	return m.functions
}

func (m Module[N]) PreDispatch(_ primitives.Call) (sc.Empty, primitives.TransactionValidityError) {
	return sc.Empty{}, nil
}

func (m Module[N]) ValidateUnsigned(_ primitives.TransactionSource, _ primitives.Call) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	return primitives.ValidTransaction{}, primitives.NewTransactionValidityError(primitives.NewUnknownTransactionNoUnsignedValidator())
}

// DepositIntoExisting deposits `value` into the free balance of an existing target account `who`.
// If `value` is 0, it does nothing.
func (m Module[N]) DepositIntoExisting(who primitives.Address32, value sc.U128) (primitives.Balance, primitives.DispatchError) {
	if value.Eq(constants.Zero) {
		return sc.NewU128FromUint64(0), nil
	}

	result := m.tryMutateAccount(who, func(from *primitives.AccountData, isNew bool) sc.Result[sc.Encodable] {
		if isNew {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
					Index:   m.Index,
					Error:   sc.U32(errors.ErrorDeadAccount),
					Message: sc.NewOption[sc.Str](nil),
				}),
			}
		}

		from.Free = from.Free.Add(value).(sc.U128)

		m.Config.StoredMap.DepositEvent(newEventDeposit(m.Index, who.FixedSequence, value))

		return sc.Result[sc.Encodable]{}
	})

	if result.HasError {
		return primitives.Balance{}, result.Value.(primitives.DispatchError)
	}

	return value, nil
}

// Withdraw withdraws `value` free balance from `who`, respecting existence requirements.
// Does not do anything if value is 0.
func (m Module[N]) Withdraw(who primitives.Address32, value sc.U128, reasons sc.U8, liveness primitives.ExistenceRequirement) (primitives.Balance, primitives.DispatchError) {
	if value.Eq(constants.Zero) {
		return sc.NewU128FromUint64(0), nil
	}

	result := m.tryMutateAccount(who, func(account *primitives.AccountData, _ bool) sc.Result[sc.Encodable] {
		newFromAccountFree := account.Free.Sub(value)

		if newFromAccountFree.Lt(constants.Zero) {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
					Index:   m.Index,
					Error:   sc.U32(errors.ErrorInsufficientBalance),
					Message: sc.NewOption[sc.Str](nil),
				}),
			}
		}

		existentialDeposit := m.Constants.ExistentialDeposit
		sumNewFreeReserved := newFromAccountFree.Add(account.Reserved)
		sumFreeReserved := account.Free.Add(account.Reserved)

		wouldBeDead := sumNewFreeReserved.Lt(existentialDeposit)
		wouldKill := wouldBeDead && (sumFreeReserved.Gte(existentialDeposit))

		if !(liveness == primitives.ExistenceRequirementAllowDeath || !wouldKill) {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
					Index:   m.Index,
					Error:   sc.U32(errors.ErrorKeepAlive),
					Message: sc.NewOption[sc.Str](nil),
				}),
			}
		}

		err := m.ensureCanWithdraw(who, value, primitives.Reasons(reasons), newFromAccountFree.(sc.U128))
		if err != nil {
			return sc.Result[sc.Encodable]{
				HasError: true,
				Value:    err,
			}
		}

		account.Free = newFromAccountFree.(sc.U128)

		m.Config.StoredMap.DepositEvent(newEventWithdraw(m.Index, who.FixedSequence, value))

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
func (m Module[N]) ensureCanWithdraw(who primitives.Address32, amount sc.U128, reasons primitives.Reasons, newBalance sc.U128) primitives.DispatchError {
	if amount.Eq(constants.Zero) {
		return nil
	}

	accountInfo := m.Config.StoredMap.Get(who.FixedSequence)
	minBalance := accountInfo.Frozen(reasons)
	if minBalance.Gt(newBalance) {
		return primitives.NewDispatchErrorModule(primitives.CustomModuleError{
			Index:   m.Index,
			Error:   sc.U32(errors.ErrorLiquidityRestrictions),
			Message: sc.NewOption[sc.Str](nil),
		})
	}

	return nil
}

// tryMutateAccount mutates an account based on argument `f`. Does not change total issuance.
// Does not do anything if `f` returns an error.
func (m Module[N]) tryMutateAccount(who primitives.Address32, f func(who *primitives.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	result := m.tryMutateAccountWithDust(who, f)
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

func (m Module[N]) tryMutateAccountWithDust(who primitives.Address32, f func(who *primitives.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	result := m.Config.StoredMap.TryMutateExists(who, func(maybeAccount *primitives.AccountData) sc.Result[sc.Encodable] {
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
		maybeAccountWithDust, imbalance := m.postMutation(*account)
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
		m.Config.StoredMap.DepositEvent(newEventEndowed(m.Index, who.FixedSequence, maybeEndowed.Value))
	}
	maybeDust := resultValue[1].(sc.Option[negativeImbalance])
	dustCleaner := newDustCleanerValue(m.Index, who, maybeDust.Value, m.Config.StoredMap)

	r := sc.NewVaryingData(resultValue[2], dustCleaner)

	return sc.Result[sc.Encodable]{HasError: false, Value: r}
}

func (m Module[N]) postMutation(new primitives.AccountData) (sc.Option[primitives.AccountData], sc.Option[negativeImbalance]) {
	total := new.Total()

	if total.Lt(m.Constants.ExistentialDeposit) {
		if total.Eq(constants.Zero) {
			return sc.NewOption[primitives.AccountData](nil), sc.NewOption[negativeImbalance](nil)
		} else {
			return sc.NewOption[primitives.AccountData](nil), sc.NewOption[negativeImbalance](newNegativeImbalance(total))
		}
	}

	return sc.NewOption[primitives.AccountData](new), sc.NewOption[negativeImbalance](nil)
}

func (m Module[N]) Metadata() (sc.Sequence[primitives.MetadataType], primitives.MetadataModule) {
	return m.metadataTypes(), primitives.MetadataModule{
		Name:    m.name(),
		Storage: m.metadataStorage(),
		Call:    sc.NewOption[sc.Compact](sc.ToCompact(metadata.BalancesCalls)),
		CallDef: sc.NewOption[primitives.MetadataDefinitionVariant](
			primitives.NewMetadataDefinitionVariantStr(
				m.name(),
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.BalancesCalls, "self::sp_api_hidden_includes_construct_runtime::hidden_include::dispatch\n::CallableCallFor<Balances, Runtime>"),
				},
				m.Index,
				"Call.Balances"),
		),
		Event: sc.NewOption[sc.Compact](sc.ToCompact(metadata.TypesBalancesEvent)),
		EventDef: sc.NewOption[primitives.MetadataDefinitionVariant](
			primitives.NewMetadataDefinitionVariantStr(
				m.name(),
				sc.Sequence[primitives.MetadataTypeDefinitionField]{
					primitives.NewMetadataTypeDefinitionFieldWithName(metadata.TypesBalancesEvent, "pallet_balances::Event<Runtime>"),
				},
				m.Index,
				"Events.Balances"),
		),
		Constants: m.metadataConstants(),
		Error:     sc.NewOption[sc.Compact](sc.ToCompact(metadata.TypesBalancesErrors)),
		Index:     m.Index,
	}
}

func (m Module[N]) metadataTypes() sc.Sequence[primitives.MetadataType] {
	return sc.Sequence[primitives.MetadataType]{
		primitives.NewMetadataTypeWithPath(metadata.TypesBalancesEvent, "pallet_balances pallet Event", sc.Sequence[sc.Str]{"pallet_balances", "pallet", "Event"}, primitives.NewMetadataTypeDefinitionVariant(
			sc.Sequence[primitives.MetadataDefinitionVariant]{
				primitives.NewMetadataDefinitionVariant(
					"Endowed",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "account", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "free_balance", "T::Balance"),
					},
					EventEndowed,
					"Event.Endowed"),
				primitives.NewMetadataDefinitionVariant(
					"DustLost",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "account", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					EventDustLost,
					"Events.DustLost"),
				primitives.NewMetadataDefinitionVariant(
					"Transfer",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "from", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "to", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					EventTransfer,
					"Events.Transfer"),
				primitives.NewMetadataDefinitionVariant(
					"BalanceSet",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "free", "T::Balance"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "reserved", "T::Balance"),
					},
					EventBalanceSet,
					"Events.BalanceSet"),
				primitives.NewMetadataDefinitionVariant(
					"Reserved",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					EventReserved,
					"Events.Reserved"),
				primitives.NewMetadataDefinitionVariant(
					"Unreserved",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					EventUnreserved,
					"Events.Unreserved"),
				primitives.NewMetadataDefinitionVariant(
					"ReserveRepatriated",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "from", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "to", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesBalanceStatus, "destination_status", "Status"),
					},
					EventReserveRepatriated,
					"Events.ReserveRepatriated"),
				primitives.NewMetadataDefinitionVariant(
					"Deposit",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					EventDeposit,
					"Event.Deposit"),
				primitives.NewMetadataDefinitionVariant(
					"Withdraw",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					EventWithdraw,
					"Event.Withdraw"),
				primitives.NewMetadataDefinitionVariant(
					"Slashed",
					sc.Sequence[primitives.MetadataTypeDefinitionField]{
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.TypesAddress32, "who", "T::AccountId"),
						primitives.NewMetadataTypeDefinitionFieldWithNames(metadata.PrimitiveTypesU128, "amount", "T::Balance"),
					},
					EventSlashed,
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

func (m Module[N]) metadataStorage() sc.Option[primitives.MetadataModuleStorage] {
	return sc.NewOption[primitives.MetadataModuleStorage](primitives.MetadataModuleStorage{
		Prefix: m.name(),
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
	})
}

func (m Module[N]) metadataConstants() sc.Sequence[primitives.MetadataModuleConstant] {
	return sc.Sequence[primitives.MetadataModuleConstant]{
		primitives.NewMetadataModuleConstant(
			"ExistentialDeposit",
			sc.ToCompact(metadata.PrimitiveTypesU128),
			sc.BytesToSequenceU8(m.Constants.ExistentialDeposit.Bytes()),
			"The minimum amount required to keep an account open. MUST BE GREATER THAN ZERO!",
		),
		primitives.NewMetadataModuleConstant(
			"MaxLocks",
			sc.ToCompact(metadata.PrimitiveTypesU32),
			sc.BytesToSequenceU8(m.Constants.MaxLocks.Bytes()),
			"The maximum number of locks that should exist on an account.  Not strictly enforced, but used for weight estimation.",
		),
		primitives.NewMetadataModuleConstant(
			"MaxReserves",
			sc.ToCompact(metadata.PrimitiveTypesU32),
			sc.BytesToSequenceU8(m.Constants.MaxReserves.Bytes()),
			"The maximum number of named reserves that can exist on an account.",
		),
	} // TODO: add more
}
