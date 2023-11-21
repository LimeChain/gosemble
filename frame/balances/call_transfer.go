package balances

import (
	"bytes"
	"reflect"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type callTransfer[T primitives.PublicKey] struct {
	primitives.Callable
	transfer
}

func newCallTransfer[T primitives.PublicKey](moduleId sc.U8, functionId sc.U8, storedMap primitives.StoredMap, constants *consts,
	mutator accountMutator) primitives.Call {
	call := callTransfer[T]{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
		},
		transfer: newTransfer(moduleId, storedMap, constants, mutator),
	}

	return call
}

func (c callTransfer[T]) DecodeArgs(buffer *bytes.Buffer) (primitives.Call, error) {
	dest, err := primitives.DecodeMultiAddress[testPublicKeyType](buffer)
	if err != nil {
		return nil, err
	}
	balance, err := sc.DecodeCompact(buffer)
	if err != nil {
		return nil, err
	}
	c.Arguments = sc.NewVaryingData(
		dest,
		balance,
	)
	return c, nil
}

func (c callTransfer[T]) Encode(buffer *bytes.Buffer) error {
	return c.Callable.Encode(buffer)
}

func (c callTransfer[T]) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callTransfer[T]) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callTransfer[T]) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callTransfer[T]) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (c callTransfer[T]) BaseWeight() primitives.Weight {
	// Proof Size summary in bytes:
	//  Measured:  `0`
	//  Estimated: `3593`
	// Minimum execution time: 37_815 nanoseconds.
	r := c.constants.DbWeight.Reads(1)
	w := c.constants.DbWeight.Writes(1)
	e := primitives.WeightFromParts(0, 3593)
	return primitives.WeightFromParts(38_109_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
}

func (_ callTransfer[T]) WeighData(baseWeight primitives.Weight) primitives.Weight {
	return primitives.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callTransfer[T]) ClassifyDispatch(baseWeight primitives.Weight) primitives.DispatchClass {
	return primitives.NewDispatchClassNormal()
}

func (_ callTransfer[T]) PaysFee(baseWeight primitives.Weight) primitives.Pays {
	return primitives.NewPaysYes()
}

func (c callTransfer[T]) Dispatch(origin primitives.RuntimeOrigin, args sc.VaryingData) primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo] {
	value := sc.U128(args[1].(sc.Compact))

	err := c.transfer.transfer(origin, args[0].(primitives.MultiAddress), value)
	if err.VaryingData != nil {
		return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
			HasError: true,
			Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
				Error: err,
			},
		}
	}

	return primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: false,
		Ok:       primitives.PostDispatchInfo{},
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
func (t transfer) transfer(origin primitives.RawOrigin, dest primitives.MultiAddress, value sc.U128) primitives.DispatchError {
	if !origin.IsSignedOrigin() {
		return primitives.NewDispatchErrorBadOrigin()
	}

	to, err := primitives.Lookup(dest)
	if err != nil {
		return primitives.NewDispatchErrorCannotLookup()
	}

	transactor, originErr := origin.AsSigned()
	if err != nil {
		log.Critical(originErr.Error())
	}

	return t.trans(transactor, to, value, primitives.ExistenceRequirementAllowDeath)
}

// trans transfers `value` free balance from `from` to `to`.
// Does not do anything if value is 0 or `from` and `to` are the same.
func (t transfer) trans(from primitives.AccountId[primitives.PublicKey], to primitives.AccountId[primitives.PublicKey], value sc.U128, existenceRequirement primitives.ExistenceRequirement) primitives.DispatchError {
	if value.Eq(constants.Zero) || reflect.DeepEqual(from, to) {
		return primitives.DispatchError{VaryingData: nil}
	}

	result := t.accountMutator.tryMutateAccountWithDust(to, func(toAccount *primitives.AccountData, _ bool) sc.Result[sc.Encodable] {
		return t.accountMutator.tryMutateAccountWithDust(from, func(fromAccount *primitives.AccountData, _ bool) sc.Result[sc.Encodable] {
			return t.sanityChecks(from, fromAccount, toAccount, value, existenceRequirement)
		})
	})
	if result.HasError {
		return result.Value.(primitives.DispatchError)
	}

	t.storedMap.DepositEvent(newEventTransfer(t.moduleId, from, to, value))
	return primitives.DispatchError{VaryingData: nil}
}

// sanityChecks checks the following:
// `fromAccount` has sufficient balance
// `toAccount` balance does not overflow
// `toAccount` total balance is more than the existential deposit
// `fromAccount` can withdraw `value`
// the existence requirements for `fromAccount`
// Updates the balances of `fromAccount` and `toAccount`.
func (t transfer) sanityChecks(from primitives.AccountId[primitives.PublicKey], fromAccount *primitives.AccountData, toAccount *primitives.AccountData, value sc.U128, existenceRequirement primitives.ExistenceRequirement) sc.Result[sc.Encodable] {
	fromFree, err := sc.CheckedSubU128(fromAccount.Free, value)
	if err != nil {
		return sc.Result[sc.Encodable]{
			HasError: true,
			Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
				Index:   t.moduleId,
				Err:     sc.U32(ErrorInsufficientBalance),
				Message: sc.NewOption[sc.Str](nil),
			}),
		}
	}
	fromAccount.Free = fromFree

	toFree, err := sc.CheckedAddU128(toAccount.Free, value)
	if err != nil {
		return sc.Result[sc.Encodable]{
			HasError: true,
			Value:    primitives.NewDispatchErrorArithmetic(primitives.NewArithmeticErrorOverflow()),
		}
	}
	toAccount.Free = toFree

	if toAccount.Total().Lt(t.constants.ExistentialDeposit) {
		return sc.Result[sc.Encodable]{
			HasError: true,
			Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
				Index:   t.moduleId,
				Err:     sc.U32(ErrorExistentialDeposit),
				Message: sc.NewOption[sc.Str](nil),
			}),
		}
	}

	dispatchErr := t.accountMutator.ensureCanWithdraw(from, value, primitives.ReasonsAll, fromAccount.Free)
	if dispatchErr.VaryingData != nil {
		return sc.Result[sc.Encodable]{
			HasError: true,
			Value:    dispatchErr,
		}
	}

	canDecProviders, err := t.storedMap.CanDecProviders(from)
	if err != nil {
		return sc.Result[sc.Encodable]{
			HasError: true,
			Value:    primitives.NewDispatchErrorOther(sc.Str(err.Error())),
		}
	}
	allowDeath := existenceRequirement == primitives.ExistenceRequirementAllowDeath
	allowDeath = allowDeath && canDecProviders

	if !(allowDeath || fromAccount.Total().Gt(t.constants.ExistentialDeposit)) {
		return sc.Result[sc.Encodable]{
			HasError: true,
			Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
				Index:   t.moduleId,
				Err:     sc.U32(ErrorKeepAlive),
				Message: sc.NewOption[sc.Str](nil),
			}),
		}
	}

	return sc.Result[sc.Encodable]{}
}

func (t transfer) reducibleBalance(who primitives.AccountId[primitives.PublicKey], keepAlive bool) (primitives.Balance, error) {
	account, err := t.storedMap.Get(who)
	if err != nil {
		return primitives.Balance{}, err
	}
	accountData := account.Data

	liquid := sc.SaturatingSubU128(accountData.Free, sc.Max128(accountData.FeeFrozen, accountData.MiscFrozen))
	canDecProviders, err := t.storedMap.CanDecProviders(who)
	if err != nil {
		return primitives.Balance{}, err
	}
	if canDecProviders && !keepAlive {
		return liquid, nil
	}

	mustRemainToExist := sc.SaturatingSubU128(t.constants.ExistentialDeposit, accountData.Total().Sub(liquid))
	return sc.SaturatingSubU128(liquid, mustRemainToExist), nil
}
