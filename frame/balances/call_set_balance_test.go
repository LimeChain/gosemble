package balances

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	newFree     = sc.NewU128(5)
	newReserved = sc.NewU128(6)
	oldFree     = sc.NewU128(4)
	oldReserved = sc.NewU128(3)
)

func Test_Call_SetBalance_new(t *testing.T) {
	target := setupCallSetBalance()
	expected := callSetBalance{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionSetBalanceIndex,
		},
		constants:      testConstants,
		storedMap:      mockStoredMap,
		accountMutator: mockMutator,
		issuance:       mockStorageTotalIssuance,
	}

	assert.Equal(t, expected, target)
}

func Test_Call_SetBalance_DecodeArgs(t *testing.T) {
	freeAmount := sc.ToCompact(sc.NewU128(1))
	reserveAmount := sc.ToCompact(sc.NewU128(1))
	buf := &bytes.Buffer{}
	buf.Write(targetAddress.Bytes())
	buf.Write(freeAmount.Bytes())
	buf.Write(reserveAmount.Bytes())

	target := setupCallSetBalance()
	call := target.DecodeArgs(buf)

	assert.Equal(t, sc.NewVaryingData(targetAddress, freeAmount, reserveAmount), call.Args())
}

func Test_Call_SetBalance_Encode(t *testing.T) {
	target := setupCallSetBalance()
	expectedBuffer := bytes.NewBuffer([]byte{moduleId, functionSetBalanceIndex})
	buf := &bytes.Buffer{}

	target.Encode(buf)

	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_SetBalance_Bytes(t *testing.T) {
	expected := []byte{moduleId, functionSetBalanceIndex}

	target := setupCallSetBalance()

	assert.Equal(t, expected, target.Bytes())
}

func Test_Call_SetBalance_ModuleIndex(t *testing.T) {
	target := setupCallSetBalance()

	assert.Equal(t, sc.U8(moduleId), target.ModuleIndex())
}

func Test_Call_SetBalance_FunctionIndex(t *testing.T) {
	target := setupCallSetBalance()

	assert.Equal(t, sc.U8(functionSetBalanceIndex), target.FunctionIndex())
}

func Test_Call_SetBalance_BaseWeight(t *testing.T) {
	target := setupCallSetBalance()

	assert.Equal(t, primitives.WeightFromParts(17_777_003, 3593), target.BaseWeight())
}

func Test_Call_SetBalance_IsInherent(t *testing.T) {
	assert.Equal(t, false, setupCallSetBalance().IsInherent())
}

func Test_Call_SetBalance_WeighData(t *testing.T) {
	target := setupCallSetBalance()
	assert.Equal(t, primitives.WeightFromParts(124, 0), target.WeighData(baseWeight))
}

func Test_Call_SetBalance_ClassifyDispatch(t *testing.T) {
	target := setupCallSetBalance()

	assert.Equal(t, primitives.NewDispatchClassNormal(), target.ClassifyDispatch(baseWeight))
}

func Test_Call_SetBalance_PaysFee(t *testing.T) {
	target := setupCallSetBalance()

	assert.Equal(t, primitives.NewPaysYes(), target.PaysFee(baseWeight))
}

func Test_Call_SetBalance_Dispatch_Success(t *testing.T) {
	newFree := sc.NewU128(0)
	newReserved := sc.NewU128(0)
	target := setupCallSetBalance()
	expect := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: false,
		Ok:       primitives.PostDispatchInfo{},
	}

	mockResult := sc.Result[sc.Encodable]{
		Value: sc.NewVaryingData(sc.NewU128(0), sc.NewU128(0)),
	}

	mockMutator.On(
		"tryMutateAccount",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	).
		Return(mockResult)
	mockStoredMap.On("DepositEvent", newEventBalanceSet(moduleId, targetAddress.AsAddress32().FixedSequence, newFree, newReserved))

	result := target.Dispatch(primitives.NewRawOriginRoot(), sc.NewVaryingData(targetAddress, sc.ToCompact(newFree), sc.ToCompact(newReserved)))

	assert.Equal(t, expect, result)
	mockStorageTotalIssuance.AssertNotCalled(t, "Get")
	mockStorageTotalIssuance.AssertNotCalled(t, "Put", mock.Anything)
	mockMutator.AssertCalled(t,
		"tryMutateAccount",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	)
	mockStoredMap.AssertCalled(t,
		"DepositEvent",
		newEventBalanceSet(moduleId, targetAddress.AsAddress32().FixedSequence, newFree, newReserved),
	)
}

func Test_Call_SetBalance_Dispatch_Fails(t *testing.T) {
	target := setupCallSetBalance()
	expect := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: true,
		Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
			Error: primitives.NewDispatchErrorBadOrigin(),
		},
	}

	result := target.Dispatch(
		primitives.NewRawOriginNone(),
		sc.NewVaryingData(targetAddress, sc.ToCompact(newFree), sc.ToCompact(newReserved)))

	assert.Equal(t, expect, result)
	mockMutator.AssertNotCalled(t, "tryMutateAccount", mock.Anything, mock.Anything)
	mockStorageTotalIssuance.AssertNotCalled(t, "Get")
	mockStorageTotalIssuance.AssertNotCalled(t, "Put", mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_SetBalance_setBalance_Success(t *testing.T) {
	target := setupCallSetBalance()
	mockResult := sc.Result[sc.Encodable]{
		Value: sc.NewVaryingData(oldFree, oldReserved),
	}

	mockMutator.On(
		"tryMutateAccount",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	).Return(mockResult)
	mockStorageTotalIssuance.On("Get").Return(sc.NewU128(1))                                            // positive imbalance
	mockStorageTotalIssuance.On("Put", newFree.Sub(oldFree).Add(sc.NewU128(1))).Return().Once()         // newFree positive imbalance
	mockStorageTotalIssuance.On("Put", newReserved.Sub(oldReserved).Add(sc.NewU128(1))).Return().Once() // newReserved positive imbalance
	mockStoredMap.On(
		"DepositEvent",
		newEventBalanceSet(moduleId, targetAddress.AsAddress32().FixedSequence, newFree, newReserved))

	result := target.setBalance(primitives.NewRawOriginRoot(), targetAddress, newFree, newReserved)

	assert.Equal(t, sc.VaryingData(nil), result)
	mockMutator.AssertCalled(t,
		"tryMutateAccount",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	)
	mockStorageTotalIssuance.AssertNumberOfCalls(t, "Get", 2)
	mockStorageTotalIssuance.AssertNumberOfCalls(t, "Put", 2)
	mockStorageTotalIssuance.AssertCalled(t, "Put", newFree.Sub(oldFree).Add(sc.NewU128(1)))
	mockStorageTotalIssuance.AssertCalled(t, "Put", newReserved.Sub(oldReserved).Add(sc.NewU128(1)))
	mockStoredMap.AssertCalled(t,
		"DepositEvent",
		newEventBalanceSet(moduleId, targetAddress.AsAddress32().FixedSequence, newFree, newReserved),
	)
}

func Test_Call_SetBalance_setBalance_Success_LessThanExistentialDeposit(t *testing.T) {
	newFree := sc.NewU128(0)
	newReserved := sc.NewU128(0)
	target := setupCallSetBalance()
	mockResult := sc.Result[sc.Encodable]{
		Value: sc.NewVaryingData(sc.NewU128(0), sc.NewU128(0)),
	}

	mockMutator.On(
		"tryMutateAccount",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	).Return(mockResult)
	mockStoredMap.On(
		"DepositEvent",
		newEventBalanceSet(moduleId, targetAddress.AsAddress32().FixedSequence, newFree, newReserved))

	result := target.setBalance(primitives.NewRawOriginRoot(), targetAddress, newFree, newReserved)

	assert.Equal(t, sc.VaryingData(nil), result)
	mockStorageTotalIssuance.AssertNotCalled(t, "Get")
	mockStorageTotalIssuance.AssertNotCalled(t, "Put", mock.Anything)
	mockMutator.AssertCalled(t,
		"tryMutateAccount",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	)
	mockStoredMap.AssertCalled(t,
		"DepositEvent",
		newEventBalanceSet(moduleId, targetAddress.AsAddress32().FixedSequence, newFree, newReserved),
	)
}

func Test_Call_SetBalance_setBalance_Success_NegativeImbalance(t *testing.T) {
	newFree := sc.NewU128(1)
	newReserved := sc.NewU128(1)
	target := setupCallSetBalance()
	mockResult := sc.Result[sc.Encodable]{
		Value: sc.NewVaryingData(oldFree, oldReserved),
	}

	mockMutator.On("tryMutateAccount",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	).Return(mockResult)
	mockStorageTotalIssuance.On("Get").Return(oldReserved.Add(oldFree)).Once() // newFree negative imbalance
	mockStorageTotalIssuance.On("Put", oldFree).Return().Once()                // newFree negative imbalance
	mockStorageTotalIssuance.On("Get").Return(sc.NewU128(4)).Once()            // newReserved negative imbalance
	mockStorageTotalIssuance.On("Put", sc.NewU128(2)).Return().Once()          // newReserved negative imbalance
	mockStoredMap.On("DepositEvent", newEventBalanceSet(moduleId, targetAddress.AsAddress32().FixedSequence, newFree, newReserved))

	result := target.setBalance(primitives.NewRawOriginRoot(), targetAddress, newFree, newReserved)

	assert.Equal(t, sc.VaryingData(nil), result)
	mockMutator.AssertCalled(t,
		"tryMutateAccount",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	)
	mockStorageTotalIssuance.AssertNumberOfCalls(t, "Get", 2)
	mockStorageTotalIssuance.AssertNumberOfCalls(t, "Put", 2)
	mockStorageTotalIssuance.AssertCalled(t, "Put", sc.NewU128(4))
	mockStorageTotalIssuance.AssertCalled(t, "Put", sc.NewU128(2))
	mockStoredMap.AssertCalled(t,
		"DepositEvent",
		newEventBalanceSet(moduleId, targetAddress.AsAddress32().FixedSequence, newFree, newReserved),
	)
}

func Test_Call_SetBalance_setBalance_InvalidOrigin(t *testing.T) {
	target := setupCallSetBalance()

	result := target.setBalance(primitives.NewRawOriginNone(), targetAddress, targetValue, targetValue)

	assert.Equal(t, primitives.NewDispatchErrorBadOrigin(), result)
	mockMutator.AssertNotCalled(t, "tryMutateAccount", mock.Anything, mock.Anything)
	mockStorageTotalIssuance.AssertNotCalled(t, "Get")
	mockStorageTotalIssuance.AssertNotCalled(t, "Put", mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_SetBalance_setBalance_Lookup(t *testing.T) {
	target := setupCallSetBalance()

	result := target.
		setBalance(primitives.NewRawOriginRoot(), primitives.NewMultiAddress20(primitives.Address20{}), targetValue, targetValue)

	assert.Equal(t, primitives.NewDispatchErrorCannotLookup(), result)
	mockMutator.AssertNotCalled(t, "tryMutateAccount", mock.Anything, mock.Anything)
	mockStorageTotalIssuance.AssertNotCalled(t, "Get")
	mockStorageTotalIssuance.AssertNotCalled(t, "Put", mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_SetBalance_setBalance_tryMutateAccount_Fails(t *testing.T) {
	target := setupCallSetBalance()
	err := primitives.NewDispatchErrorBadOrigin()
	mockResult := sc.Result[sc.Encodable]{
		HasError: true,
		Value:    err,
	}
	mockMutator.On(
		"tryMutateAccount",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	).Return(mockResult)

	result := target.setBalance(primitives.NewRawOriginRoot(), targetAddress, targetValue, targetValue)

	assert.Equal(t, err, result)
	mockMutator.AssertCalled(t,
		"tryMutateAccount",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	)
	mockStorageTotalIssuance.AssertNotCalled(t, "Get")
	mockStorageTotalIssuance.AssertNotCalled(t, "Put", mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_SetBalance_updateAccount(t *testing.T) {
	oldFree := sc.NewU128(1)
	oldReserved := sc.NewU128(2)
	newFree := sc.NewU128(5)
	newReserved := sc.NewU128(6)

	account := &primitives.AccountData{
		Free:       oldFree,
		Reserved:   oldReserved,
		MiscFrozen: sc.NewU128(3),
		FeeFrozen:  sc.NewU128(4),
	}
	expectAccount := &primitives.AccountData{
		Free:       newFree,
		Reserved:   newReserved,
		MiscFrozen: sc.NewU128(3),
		FeeFrozen:  sc.NewU128(4),
	}

	result := updateAccount(account, newFree, newReserved)

	assert.Equal(t, sc.Result[sc.Encodable]{Value: sc.NewVaryingData(oldFree, oldReserved)}, result)
	assert.Equal(t, expectAccount, account)
}

func setupCallSetBalance() callSetBalance {
	mockStoredMap = new(mocks.MockStoredMap)
	mockMutator = new(mockAccountMutator)
	mockStorageTotalIssuance = new(mocks.StorageValue[sc.U128])

	return newCallSetBalance(moduleId, functionSetBalanceIndex, mockStoredMap, testConstants, mockMutator, mockStorageTotalIssuance).(callSetBalance)
}
