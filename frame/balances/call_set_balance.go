package balances

import (
	"bytes"
	"errors"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/support"
	"github.com/LimeChain/gosemble/primitives/types"
)

type callSetBalance struct {
	types.Callable
	constants      *consts
	storedMap      types.StoredMap
	accountMutator accountMutator
	issuance       support.StorageValue[sc.U128]
}

func newCallSetBalance(moduleId sc.U8, functionId sc.U8, storedMap types.StoredMap, constants *consts, mutator accountMutator, issuance support.StorageValue[sc.U128]) types.Call {
	call := callSetBalance{
		Callable: types.Callable{
			ModuleId:   moduleId,
			FunctionId: functionId,
			Arguments:  sc.NewVaryingData(types.MultiAddress{}, sc.Compact{Number: sc.U128{}}, sc.Compact{Number: sc.U128{}}),
		},
		constants:      constants,
		storedMap:      storedMap,
		accountMutator: mutator,
		issuance:       issuance,
	}

	return call
}

func (c callSetBalance) DecodeArgs(buffer *bytes.Buffer) (types.Call, error) {
	targetAddress, err := types.DecodeMultiAddress(buffer)
	if err != nil {
		return nil, err
	}
	newFree, err := sc.DecodeCompact[sc.U128](buffer)
	if err != nil {
		return nil, err
	}
	newReserved, err := sc.DecodeCompact[sc.U128](buffer)
	if err != nil {
		return nil, err
	}

	c.Arguments = sc.NewVaryingData(
		targetAddress,
		newFree,
		newReserved,
	)
	return c, nil
}

func (c callSetBalance) Encode(buffer *bytes.Buffer) error {
	return c.Callable.Encode(buffer)
}

func (c callSetBalance) Bytes() []byte {
	return c.Callable.Bytes()
}

func (c callSetBalance) ModuleIndex() sc.U8 {
	return c.Callable.ModuleIndex()
}

func (c callSetBalance) FunctionIndex() sc.U8 {
	return c.Callable.FunctionIndex()
}

func (c callSetBalance) Args() sc.VaryingData {
	return c.Callable.Args()
}

func (c callSetBalance) BaseWeight() types.Weight {
	return callSetBalanceCreatingWeight(c.constants.DbWeight).Max(callSetBalanceKillingWeight(c.constants.DbWeight))
}

func (_ callSetBalance) IsInherent() bool {
	return false
}

func (_ callSetBalance) WeighData(baseWeight types.Weight) types.Weight {
	return types.WeightFromParts(baseWeight.RefTime, 0)
}

func (_ callSetBalance) ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	return types.NewDispatchClassNormal()
}

func (_ callSetBalance) PaysFee(baseWeight types.Weight) types.Pays {
	return types.PaysYes
}

func (c callSetBalance) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) (types.PostDispatchInfo, error) {
	compactFree, ok := args[1].(sc.Compact)
	if !ok {
		return types.PostDispatchInfo{}, errors.New("invalid free compact value when dispatching balance call set")
	}
	newFree, ok := compactFree.Number.(sc.U128)
	if !ok {
		return types.PostDispatchInfo{}, errors.New("invalid free compact number when dispatching balance call set")
	}

	compactReserved, ok := args[2].(sc.Compact)
	if !ok {
		return types.PostDispatchInfo{}, errors.New("invalid reserved compact value when dispatching balance call set")
	}
	newReserved, ok := compactReserved.Number.(sc.U128)
	if !ok {
		return types.PostDispatchInfo{}, errors.New("invalid reserved compact number when dispatching balance call set")
	}
	return types.PostDispatchInfo{}, c.setBalance(origin, args[0].(types.MultiAddress), newFree, newReserved)
}

func (_ callSetBalance) Docs() string {
	return "Set the balances of a given account."
}

// setBalance sets the balance of a given account.
// Changes free and reserve balance of `who`,
// including the total issuance.
// Can only be called by ROOT.
func (c callSetBalance) setBalance(origin types.RawOrigin, who types.MultiAddress, newFree sc.U128, newReserved sc.U128) error {
	if !origin.IsRootOrigin() {
		return types.NewDispatchErrorBadOrigin()
	}

	address, err := types.Lookup(who)
	if err != nil {
		return types.NewDispatchErrorCannotLookup()
	}

	sum := newFree.Add(newReserved)

	if sum.Lt(c.constants.ExistentialDeposit) {
		newFree = sc.NewU128(0)
		newReserved = sc.NewU128(0)
	}

	result, err := c.accountMutator.tryMutateAccount(
		address,
		func(account *types.AccountData, _ bool) (sc.Encodable, error) {
			oldFree, oldReserved := updateAccount(account, newFree, newReserved)
			return sc.NewVaryingData(oldFree, oldReserved), nil
		},
	)
	if err != nil {
		return err
	}

	parsedResult := result.(sc.VaryingData)
	oldFree := parsedResult[0].(types.Balance)
	oldReserved := parsedResult[1].(types.Balance)

	if newFree.Gt(oldFree) {
		if err := newPositiveImbalance(newFree.Sub(oldFree), c.issuance).Drop(); err != nil {
			return types.NewDispatchErrorOther(sc.Str(err.Error()))
		}

	} else if newFree.Lt(oldFree) {
		if err := newNegativeImbalance(oldFree.Sub(newFree), c.issuance).Drop(); err != nil {
			return types.NewDispatchErrorOther(sc.Str(err.Error()))
		}
	}

	if newReserved.Gt(oldReserved) {
		if err := newPositiveImbalance(newReserved.Sub(oldReserved), c.issuance).Drop(); err != nil {
			return types.NewDispatchErrorOther(sc.Str(err.Error()))
		}
	} else if newReserved.Lt(oldReserved) {
		if err := newNegativeImbalance(oldReserved.Sub(newReserved), c.issuance).Drop(); err != nil {
			return types.NewDispatchErrorOther(sc.Str(err.Error()))
		}

	}

	whoAccountId, errAccId := who.AsAccountId()
	if errAccId != nil {
		return types.NewDispatchErrorOther(sc.Str(errAccId.Error()))
	}

	c.storedMap.DepositEvent(
		newEventBalanceSet(
			c.ModuleId,
			whoAccountId,
			newFree,
			newReserved,
		),
	)
	return nil
}

// updateAccount updates the reserved and free amounts and returns the old amounts
func updateAccount(account *types.AccountData, newFree, newReserved sc.U128) (oldFree, oldReserved types.Balance) {
	oldFree = account.Free
	oldReserved = account.Reserved

	account.Free = newFree
	account.Reserved = newReserved

	return oldFree, oldReserved
}
