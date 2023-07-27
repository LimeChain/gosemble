package module

import (
	"bytes"
	"math/big"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/balances/dispatchables"
	"github.com/LimeChain/gosemble/frame/balances/errors"
	"github.com/LimeChain/gosemble/frame/balances/events"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type TransferCall struct {
	primitives.Callable
	transfer
}

func NewTransferCall(moduleId sc.U8, functionId sc.U8, storedMap primitives.StoredMap, constants *consts) TransferCall {
	call := TransferCall{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		transfer: newTransfer(moduleId, storedMap, constants),
	}

	return call
}

func (c TransferCall) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(
		types.DecodeMultiAddress(buffer),
		sc.DecodeCompact(buffer),
	)
	return c
}

func (c TransferCall) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c TransferCall) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c TransferCall) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c TransferCall) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c TransferCall) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (_ TransferCall) BaseWeight(b ...any) types.Weight {
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

func (_ TransferCall) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ TransferCall) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ TransferCall) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (c TransferCall) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	value := sc.U128(args[1].(sc.Compact))

	err := c.transfer.transfer(origin, args[0].(types.MultiAddress), value)
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

func (_ TransferCall) IsInherent() bool {
	return false
}

type transfer struct {
	moduleId  sc.U8
	storedMap primitives.StoredMap
	constants *consts
}

func newTransfer(moduleId sc.U8, storedMap primitives.StoredMap, constants *consts) transfer {
	return transfer{
		moduleId:  moduleId,
		storedMap: storedMap,
		constants: constants,
	}
}

// transfer transfers liquid free balance from `source` to `dest`.
// Increases the free balance of `dest` and decreases the free balance of `origin` transactor.
// Must be signed by the transactor.
func (t transfer) transfer(origin types.RawOrigin, dest types.MultiAddress, value sc.U128) types.DispatchError {
	if !origin.IsSignedOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	to, e := types.DefaultAccountIdLookup().Lookup(dest)
	if e != nil {
		return types.NewDispatchErrorCannotLookup()
	}

	transactor := origin.AsSigned()

	return t.trans(transactor, to, value, types.ExistenceRequirementAllowDeath)
}

// trans transfers `value` free balance from `from` to `to`.
// Does not do anything if value is 0 or `from` and `to` are the same.
func (t transfer) trans(from types.Address32, to types.Address32, value sc.U128, existenceRequirement types.ExistenceRequirement) types.DispatchError {
	bnInt := value.ToBigInt()
	if bnInt.Cmp(constants.Zero) == 0 || reflect.DeepEqual(from, to) {
		return nil
	}

	result := t.tryMutateAccountWithDust(to, func(toAccount *types.AccountData, _ bool) sc.Result[sc.Encodable] {
		return t.tryMutateAccountWithDust(from, func(fromAccount *types.AccountData, _ bool) sc.Result[sc.Encodable] {
			newFromAccountFree := new(big.Int).Sub(fromAccount.Free.ToBigInt(), value.ToBigInt())

			if newFromAccountFree.Cmp(constants.Zero) < 0 {
				return sc.Result[sc.Encodable]{
					HasError: true,
					Value: types.NewDispatchErrorModule(types.CustomModuleError{
						Index:   t.moduleId,
						Error:   sc.U32(errors.ErrorInsufficientBalance),
						Message: sc.NewOption[sc.Str](nil),
					}),
				}
			}
			fromAccount.Free = sc.NewU128FromBigInt(newFromAccountFree)

			newToAccountFree := new(big.Int).Add(toAccount.Free.ToBigInt(), value.ToBigInt())
			toAccount.Free = sc.NewU128FromBigInt(newToAccountFree)

			existentialDeposit := t.constants.ExistentialDeposit
			if toAccount.Total().Cmp(existentialDeposit) < 0 {
				return sc.Result[sc.Encodable]{
					HasError: true,
					Value: types.NewDispatchErrorModule(types.CustomModuleError{
						Index:   t.moduleId,
						Error:   sc.U32(errors.ErrorExistentialDeposit),
						Message: sc.NewOption[sc.Str](nil),
					}),
				}
			}

			err := t.ensureCanWithdraw(from, value.ToBigInt(), types.ReasonsAll, fromAccount.Free.ToBigInt())
			if err != nil {
				return sc.Result[sc.Encodable]{
					HasError: true,
					Value:    err,
				}
			}

			allowDeath := existenceRequirement == types.ExistenceRequirementAllowDeath
			allowDeath = allowDeath && t.storedMap.CanDecProviders(from)

			if !(allowDeath || fromAccount.Total().Cmp(existentialDeposit) > 0) {
				return sc.Result[sc.Encodable]{
					HasError: true,
					Value: types.NewDispatchErrorModule(types.CustomModuleError{
						Index:   t.moduleId,
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

	t.storedMap.DepositEvent(events.NewEventTransfer(from.FixedSequence, to.FixedSequence, value))
	return nil
}

// ensureCanWithdraw checks that an account can withdraw from their balance given any existing withdraw restrictions.
func (t transfer) ensureCanWithdraw(who types.Address32, amount *big.Int, reasons types.Reasons, newBalance *big.Int) types.DispatchError {
	if amount.Cmp(constants.Zero) == 0 {
		return nil
	}

	accountInfo := t.storedMap.Get(who.FixedSequence)
	minBalance := accountInfo.Frozen(reasons)
	if minBalance.Cmp(newBalance) > 0 {
		return types.NewDispatchErrorModule(types.CustomModuleError{
			Index:   t.moduleId,
			Error:   sc.U32(errors.ErrorLiquidityRestrictions),
			Message: sc.NewOption[sc.Str](nil),
		})
	}

	return nil
}

// mutateAccount mutates an account based on argument `f`. Does not change total issuance.
// Does not do anything if `f` returns an error.
func (t transfer) mutateAccount(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	return t.tryMutateAccount(who, f)
}

// tryMutateAccount mutates an account based on argument `f`. Does not change total issuance.
// Does not do anything if `f` returns an error.
func (t transfer) tryMutateAccount(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	result := t.tryMutateAccountWithDust(who, f)
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

func (t transfer) tryMutateAccountWithDust(who types.Address32, f func(who *types.AccountData, bool bool) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	result := t.storedMap.TryMutateExists(who, func(maybeAccount *types.AccountData) sc.Result[sc.Encodable] {
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
		maybeAccountWithDust, imbalance := t.postMutation(*account)
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
		t.storedMap.DepositEvent(events.NewEventEndowed(who.FixedSequence, maybeEndowed.Value))
	}
	maybeDust := resultValue[1].(sc.Option[dispatchables.NegativeImbalance])
	dustCleaner := dispatchables.DustCleanerValue{
		AccountId:         who,
		NegativeImbalance: maybeDust.Value,
	}

	r := sc.NewVaryingData(resultValue[2], dustCleaner)

	return sc.Result[sc.Encodable]{HasError: false, Value: r}
}

func (t transfer) postMutation(
	new types.AccountData) (sc.Option[types.AccountData], sc.Option[dispatchables.NegativeImbalance]) {
	total := new.Total()

	if total.Cmp(t.constants.ExistentialDeposit) < 0 {
		if total.Cmp(constants.Zero) == 0 {
			return sc.NewOption[types.AccountData](nil), sc.NewOption[dispatchables.NegativeImbalance](nil)
		} else {
			return sc.NewOption[types.AccountData](nil), sc.NewOption[dispatchables.NegativeImbalance](dispatchables.NewNegativeImbalance(sc.NewU128FromBigInt(total)))
		}
	}

	return sc.NewOption[types.AccountData](new), sc.NewOption[dispatchables.NegativeImbalance](nil)
}

// totalBalance returns the total storage balance of an account id.
func (t transfer) totalBalance(who types.Address32) *big.Int {
	return t.storedMap.Get(who.FixedSequence).Data.Total()
}

func (t transfer) reducibleBalance(who types.Address32, keepAlive bool) types.Balance {
	accountData := t.storedMap.Get(who.FixedSequence).Data

	lockedOrFrozen := accountData.FeeFrozen
	if accountData.FeeFrozen.ToBigInt().Cmp(accountData.MiscFrozen.ToBigInt()) < 0 {
		lockedOrFrozen = accountData.MiscFrozen
	}

	liquid := new(big.Int).Sub(accountData.Free.ToBigInt(), lockedOrFrozen.ToBigInt())
	if liquid.Cmp(accountData.Free.ToBigInt()) > 0 {
		liquid = big.NewInt(0)
	}

	if t.storedMap.CanDecProviders(who) && !keepAlive {
		return sc.NewU128FromBigInt(liquid)
	}

	existentialDeposit := t.constants.ExistentialDeposit
	diff := new(big.Int).Sub(accountData.Total(), liquid)

	mustRemainToExist := new(big.Int).Sub(existentialDeposit, diff)

	result := new(big.Int).Sub(liquid, mustRemainToExist)
	if result.Cmp(liquid) > 0 {
		return sc.NewU128FromBigInt(big.NewInt(0))
	}

	return sc.NewU128FromBigInt(result)
}
