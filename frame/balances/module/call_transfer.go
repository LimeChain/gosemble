package module

import (
	"bytes"
	"math/big"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/frame/balances/errors"
	"github.com/LimeChain/gosemble/frame/balances/events"
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type transferCall struct {
	primitives.Callable
	transfer
}

func newTransferCall(moduleId sc.U8, functionId sc.U8, storedMap primitives.StoredMap, constants *consts,
	mutator accountMutator) primitives.Call {
	call := transferCall{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		transfer: newTransfer(moduleId, storedMap, constants, mutator),
	}

	return call
}

func (c transferCall) DecodeArgs(buffer *bytes.Buffer) primitives.Call {
	c.Arguments = sc.NewVaryingData(
		types.DecodeMultiAddress(buffer),
		sc.DecodeCompact(buffer),
	)
	return c
}

func (c transferCall) Encode(buffer *bytes.Buffer) {
	c.Callable.Encode(buffer)
}

func (c transferCall) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c transferCall) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c transferCall) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c transferCall) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (_ transferCall) BaseWeight(b ...any) types.Weight {
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

func (_ transferCall) WeightInfo(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ transferCall) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ transferCall) PaysFee(baseWeight types.Weight) types.Pays {
	return types.NewPaysYes()
}

func (c transferCall) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
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

type transfer struct {
	moduleId       sc.U8
	storedMap      primitives.StoredMap
	constants      *consts
	accountMutator accountMutator
}

func newTransfer(moduleId sc.U8, storedMap primitives.StoredMap, constants *consts, mutator accountMutator) transfer {
	return transfer{
		moduleId:       moduleId,
		storedMap:      storedMap,
		constants:      constants,
		accountMutator: mutator,
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

	result := t.accountMutator.tryMutateAccountWithDust(to, func(toAccount *types.AccountData, _ bool) sc.Result[sc.Encodable] {
		return t.accountMutator.tryMutateAccountWithDust(from, func(fromAccount *types.AccountData, _ bool) sc.Result[sc.Encodable] {
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

			err := t.accountMutator.ensureCanWithdraw(from, value.ToBigInt(), types.ReasonsAll, fromAccount.Free.ToBigInt())
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
