package balances

import (
	"bytes"

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
	newFree, err := sc.DecodeCompact[sc.Numeric](buffer)
	if err != nil {
		return nil, err
	}
	newReserved, err := sc.DecodeCompact[sc.Numeric](buffer)
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
	// Proof Size summary in bytes:
	//  Measured:  `206`
	//  Estimated: `3593`
	// Minimum execution time: 17_474 nanoseconds.
	r := c.constants.DbWeight.Reads(1)
	w := c.constants.DbWeight.Writes(1)
	e := types.WeightFromParts(0, 3593)
	return types.WeightFromParts(17_777_000, 0).
		SaturatingAdd(e).
		SaturatingAdd(r).
		SaturatingAdd(w)
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

func (_ callSetBalance) Docs() string {
	return "Set the balances of a given account."
}

func (c callSetBalance) Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo] {
	compactFree, _ := args[1].(sc.Compact)
	newFree := compactFree.Number.(sc.U128)

	compactReserved, _ := args[2].(sc.Compact)
	newReserved := compactReserved.Number.(sc.U128)

	err := c.setBalance(origin, args[0].(types.MultiAddress), newFree, newReserved)
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

// setBalance sets the balance of a given account.
// Changes free and reserve balance of `who`,
// including the total issuance.
// Can only be called by ROOT.
func (c callSetBalance) setBalance(origin types.RawOrigin, who types.MultiAddress, newFree sc.U128, newReserved sc.U128) types.DispatchError {
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

	result := c.accountMutator.tryMutateAccount(
		address,
		func(account *types.AccountData, _ bool) sc.Result[sc.Encodable] {
			return updateAccount(account, newFree, newReserved)
		},
	)
	if result.HasError {
		return result.Value.(types.DispatchError)
	}

	parsedResult := result.Value.(sc.VaryingData)
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
func updateAccount(account *types.AccountData, newFree, newReserved sc.U128) sc.Result[sc.Encodable] {
	oldFree := account.Free
	oldReserved := account.Reserved

	account.Free = newFree
	account.Reserved = newReserved

	return sc.Result[sc.Encodable]{
		HasError: false,
		Value:    sc.NewVaryingData(oldFree, oldReserved),
	}
}
