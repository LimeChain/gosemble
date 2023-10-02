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
	maxLocks           = sc.U32(5)
	maxReserves        = sc.U32(6)
	existentialDeposit = sc.NewU128(1)
	mockMutator        *mockAccountMutator
	testConstants      = newConstants(dbWeight, maxLocks, maxReserves, existentialDeposit)

	fromAccountData *primitives.AccountData
	toAccountData   *primitives.AccountData

	fromAddress = primitives.
			NewMultiAddress32(primitives.
				NewAddress32(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1))
	toAddress = primitives.
			NewMultiAddress32(primitives.
				NewAddress32(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2))
)

func Test_Call_Transfer_New(t *testing.T) {
	target := setupCallTransfer()
	expected := callTransfer{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionTransferIndex,
		},
		transfer: transfer{
			moduleId:       moduleId,
			storedMap:      mockStoredMap,
			constants:      testConstants,
			accountMutator: mockMutator,
		},
	}

	assert.Equal(t, expected, target)
}

func Test_Call_Transfer_DecodeArgs(t *testing.T) {
	amount := sc.ToCompact(sc.NewU128(5))
	buf := bytes.NewBuffer(append(targetAddress.Bytes(), amount.Bytes()...))

	target := setupCallTransfer()
	call := target.DecodeArgs(buf)

	assert.Equal(t, sc.NewVaryingData(targetAddress, amount), call.Args())
}

func Test_Call_Transfer_Encode(t *testing.T) {
	target := setupCallTransfer()
	expectedBuffer := bytes.NewBuffer([]byte{moduleId, functionTransferIndex})
	buf := &bytes.Buffer{}

	target.Encode(buf)

	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_Transfer_Bytes(t *testing.T) {
	expected := []byte{moduleId, functionTransferIndex}

	target := setupCallTransfer()

	assert.Equal(t, expected, target.Bytes())
}

func Test_Call_Transfer_ModuleIndex(t *testing.T) {
	target := setupCallTransfer()

	assert.Equal(t, sc.U8(moduleId), target.ModuleIndex())
}

func Test_Call_Transfer_FunctionIndex(t *testing.T) {
	target := setupCallTransfer()

	assert.Equal(t, sc.U8(functionTransferIndex), target.FunctionIndex())
}

func Test_Call_Transfer_BaseWeight(t *testing.T) {
	target := setupCallTransfer()

	assert.Equal(t, primitives.WeightFromParts(38_109_003, 3593), target.BaseWeight())
}

func Test_Call_Transfer_WeighData(t *testing.T) {
	target := setupCallTransfer()
	assert.Equal(t, primitives.WeightFromParts(124, 0), target.WeighData(baseWeight))
}

func Test_Call_Transfer_ClassifyDispatch(t *testing.T) {
	target := setupCallTransfer()

	assert.Equal(t, primitives.NewDispatchClassNormal(), target.ClassifyDispatch(baseWeight))
}

func Test_Call_Transfer_PaysFee(t *testing.T) {
	target := setupCallTransfer()

	assert.Equal(t, primitives.NewPaysYes(), target.PaysFee(baseWeight))
}

func Test_Call_Transfer_Dispatch_Success(t *testing.T) {
	target := setupCallTransfer()
	expected := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: false,
		Ok:       primitives.PostDispatchInfo{},
	}

	result := target.
		Dispatch(primitives.NewRawOriginSigned(fromAddress.AsAddress32()), sc.NewVaryingData(fromAddress, sc.ToCompact(targetValue)))

	assert.Equal(t, expected, result)
}

func Test_Call_Transfer_Dispatch_Fails(t *testing.T) {
	target := setupCallTransfer()
	expected := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: true,
		Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
			Error: primitives.NewDispatchErrorBadOrigin(),
		},
	}

	result := target.Dispatch(primitives.NewRawOriginNone(), sc.NewVaryingData(toAddress, sc.ToCompact(targetValue)))

	assert.Equal(t, expected, result)
}

func Test_transfer_New(t *testing.T) {
	target := setupTransfer()
	expected := transfer{
		moduleId:       moduleId,
		storedMap:      mockStoredMap,
		constants:      testConstants,
		accountMutator: mockMutator,
	}

	assert.Equal(t, expected, target)
}

func Test_transfer_Success(t *testing.T) {
	target := setupTransfer()

	result := target.transfer(primitives.NewRawOriginSigned(fromAddress.AsAddress32()), fromAddress, targetValue)

	assert.Equal(t, sc.VaryingData(nil), result)
}

func Test_transfer_InvalidOrigin(t *testing.T) {
	target := setupTransfer()

	result := target.transfer(primitives.NewRawOriginRoot(), toAddress, targetValue)

	assert.Equal(t, primitives.NewDispatchErrorBadOrigin(), result)
}

func Test_transfer_InvalidLookup(t *testing.T) {
	target := setupTransfer()

	result := target.
		transfer(primitives.NewRawOriginSigned(fromAddress.AsAddress32()), primitives.NewMultiAddress20(primitives.Address20{}), targetValue)

	assert.Equal(t, primitives.NewDispatchErrorCannotLookup(), result)
}

func Test_transfer_trans_Success(t *testing.T) {
	target := setupTransfer()

	mockMutator.On(
		"tryMutateAccountWithDust",
		toAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	).Return(sc.Result[sc.Encodable]{})
	mockStoredMap.On(
		"DepositEvent",
		newEventTransfer(moduleId, fromAddress.AsAddress32().FixedSequence, toAddress.AsAddress32().FixedSequence, targetValue),
	).Return()

	result := target.trans(fromAddress.AsAddress32(), toAddress.AsAddress32(), targetValue, primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, sc.VaryingData(nil), result)
	mockMutator.AssertCalled(t,
		"tryMutateAccountWithDust",
		toAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	)
	mockStoredMap.AssertCalled(t,
		"DepositEvent",
		newEventTransfer(moduleId, fromAddress.AsAddress32().FixedSequence, toAddress.AsAddress32().FixedSequence, targetValue),
	)
}

func Test_transfer_trans_ZeroValue(t *testing.T) {
	target := setupTransfer()

	result := target.trans(fromAddress.AsAddress32(), toAddress.AsAddress32(), sc.NewU128(0), primitives.ExistenceRequirementAllowDeath)

	assert.Equal(t, sc.VaryingData(nil), result)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_transfer_trans_EqualFromTo(t *testing.T) {
	target := setupTransfer()

	result := target.trans(fromAddress.AsAddress32(), fromAddress.AsAddress32(), targetValue, primitives.ExistenceRequirementAllowDeath)

	assert.Equal(t, sc.VaryingData(nil), result)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_transfer_trans_MutateAccountWithDust_Fails(t *testing.T) {
	target := setupTransfer()
	expectdError := primitives.NewDispatchErrorBadOrigin()
	error := sc.Result[sc.Encodable]{
		HasError: true,
		Value:    expectdError,
	}

	mockMutator.On(
		"tryMutateAccountWithDust",
		toAddress.AsAddress32(), mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	).Return(error)

	result := target.trans(fromAddress.AsAddress32(), toAddress.AsAddress32(), targetValue, primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, expectdError, result)
	mockMutator.AssertCalled(t,
		"tryMutateAccountWithDust",
		toAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_transfer_sanityChecks_Success(t *testing.T) {
	target := setupTransfer()
	expected := sc.Result[sc.Encodable]{}

	mockMutator.On("ensureCanWithdraw", targetAddress.AsAddress32(), targetValue, primitives.ReasonsAll, sc.NewU128(0)).Return(sc.VaryingData(nil))
	mockStoredMap.On("CanDecProviders", targetAddress.AsAddress32()).Return(true)

	result := target.sanityChecks(targetAddress.AsAddress32(), fromAccountData, toAccountData, targetValue, primitives.ExistenceRequirementAllowDeath)

	assert.Equal(t, expected, result)
	assert.Equal(t, sc.NewU128(0), fromAccountData.Free)
	assert.Equal(t, sc.NewU128(6), toAccountData.Free)
	mockMutator.AssertCalled(t, "ensureCanWithdraw", targetAddress.AsAddress32(), targetValue, primitives.ReasonsAll, sc.NewU128(0))
	mockStoredMap.AssertCalled(t, "CanDecProviders", targetAddress.AsAddress32())
}

func Test_transfer_sanityChecks_InsufficientBalance(t *testing.T) {
	target := setupTransfer()
	expected := sc.Result[sc.Encodable]{
		HasError: true,
		Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
			Index:   moduleId,
			Error:   sc.U32(ErrorInsufficientBalance),
			Message: sc.NewOption[sc.Str](nil),
		}),
	}

	result := target.sanityChecks(targetAddress.AsAddress32(), fromAccountData, toAccountData, sc.NewU128(6), primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, expected, result)
	assert.Equal(t, sc.NewU128(5), fromAccountData.Free)
	assert.Equal(t, sc.NewU128(1), toAccountData.Free)
	mockMutator.AssertNotCalled(t, "ensureCanWithdraw", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "CanDecProviders", mock.Anything)
}

func Test_transfer_sanityChecks_ArithmeticOverflow(t *testing.T) {
	target := setupTransfer()
	expected := sc.Result[sc.Encodable]{
		HasError: true,
		Value:    primitives.NewDispatchErrorArithmetic(primitives.NewArithmeticErrorOverflow()),
	}
	toAccountData.Free = sc.MaxU128()

	result := target.sanityChecks(targetAddress.AsAddress32(), fromAccountData, toAccountData, sc.NewU128(1), primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, expected, result)
	assert.Equal(t, sc.NewU128(4), fromAccountData.Free)
	assert.Equal(t, sc.MaxU128(), toAccountData.Free)
	mockMutator.AssertNotCalled(t, "ensureCanWithdraw", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "CanDecProviders", mock.Anything)
}

func Test_transfer_sanityChecks_ExistentialDeposit(t *testing.T) {
	target := setupTransfer()
	expected := sc.Result[sc.Encodable]{
		HasError: true,
		Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
			Index:   moduleId,
			Error:   sc.U32(ErrorExistentialDeposit),
			Message: sc.NewOption[sc.Str](nil),
		}),
	}
	toAccountData.Free = sc.NewU128(0)

	result := target.sanityChecks(targetAddress.AsAddress32(), fromAccountData, toAccountData, sc.NewU128(0), primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, expected, result)
	assert.Equal(t, sc.NewU128(5), fromAccountData.Free)
	assert.Equal(t, sc.NewU128(0), toAccountData.Free)
	mockMutator.AssertNotCalled(t, "ensureCanWithdraw", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "CanDecProviders", mock.Anything)
}

func Test_transfer_sanityChecks_CannotWithdraw(t *testing.T) {
	target := setupTransfer()
	expectedError := primitives.NewDispatchErrorCannotLookup()
	expected := sc.Result[sc.Encodable]{
		HasError: true,
		Value:    expectedError,
	}
	mockMutator.On("ensureCanWithdraw", targetAddress.AsAddress32(), targetValue, primitives.ReasonsAll, sc.NewU128(0)).Return(expectedError)

	result := target.sanityChecks(targetAddress.AsAddress32(), fromAccountData, toAccountData, targetValue, primitives.ExistenceRequirementAllowDeath)

	assert.Equal(t, expected, result)
	assert.Equal(t, sc.NewU128(0), fromAccountData.Free)
	assert.Equal(t, sc.NewU128(6), toAccountData.Free)
	mockMutator.AssertCalled(t, "ensureCanWithdraw", targetAddress.AsAddress32(), targetValue, primitives.ReasonsAll, sc.NewU128(0))
}

func Test_transfer_sanityChecks_KeepAlive(t *testing.T) {
	target := setupTransfer()
	expected := sc.Result[sc.Encodable]{
		HasError: true,
		Value: primitives.NewDispatchErrorModule(primitives.CustomModuleError{
			Index:   moduleId,
			Error:   sc.U32(ErrorKeepAlive),
			Message: sc.NewOption[sc.Str](nil),
		}),
	}
	mockMutator.On("ensureCanWithdraw", targetAddress.AsAddress32(), targetValue, primitives.ReasonsAll, sc.NewU128(0)).Return(sc.VaryingData(nil))
	mockStoredMap.On("CanDecProviders", targetAddress.AsAddress32()).Return(false)

	result := target.sanityChecks(targetAddress.AsAddress32(), fromAccountData, toAccountData, targetValue, primitives.ExistenceRequirementAllowDeath)

	assert.Equal(t, expected, result)
	assert.Equal(t, sc.NewU128(0), fromAccountData.Free)
	assert.Equal(t, sc.NewU128(6), toAccountData.Free)
	mockMutator.AssertCalled(t, "ensureCanWithdraw", targetAddress.AsAddress32(), targetValue, primitives.ReasonsAll, sc.NewU128(0))
	mockStoredMap.AssertCalled(t, "CanDecProviders", targetAddress.AsAddress32())
}

func Test_transfer_reducibleBalance_NotKeepAlive(t *testing.T) {
	target := setupTransfer()
	mockStoredMap.On("Get", targetAddress.AsAddress32().FixedSequence).Return(accountInfo)
	mockStoredMap.On("CanDecProviders", targetAddress.AsAddress32()).Return(true)

	result := target.reducibleBalance(targetAddress.AsAddress32(), false)

	assert.Equal(t, accountInfo.Data.Free, result)
	mockStoredMap.AssertCalled(t, "Get", targetAddress.AsAddress32().FixedSequence)
	mockStoredMap.AssertCalled(t, "CanDecProviders", targetAddress.AsAddress32())
}

func Test_transfer_reducibleBalance_KeepAlive(t *testing.T) {
	target := setupTransfer()
	mockStoredMap.On("Get", targetAddress.AsAddress32().FixedSequence).Return(accountInfo)
	mockStoredMap.On("CanDecProviders", targetAddress.AsAddress32()).Return(false)

	result := target.reducibleBalance(targetAddress.AsAddress32(), true)

	assert.Equal(t, accountInfo.Data.Free.Sub(existentialDeposit), result)
	mockStoredMap.AssertCalled(t, "Get", targetAddress.AsAddress32().FixedSequence)
	mockStoredMap.AssertCalled(t, "CanDecProviders", targetAddress.AsAddress32())
}

func setupCallTransfer() callTransfer {
	mockStoredMap = new(mocks.MockStoredMap)
	mockMutator = new(mockAccountMutator)

	fromAccountData = &primitives.AccountData{
		Free: sc.NewU128(5),
	}

	toAccountData = &primitives.AccountData{
		Free: sc.NewU128(1),
	}

	return newCallTransfer(moduleId, functionTransferIndex, mockStoredMap, testConstants, mockMutator).(callTransfer)
}

func setupTransfer() transfer {
	mockStoredMap = new(mocks.MockStoredMap)
	mockMutator = new(mockAccountMutator)

	fromAccountData = &primitives.AccountData{
		Free: sc.NewU128(5),
	}

	toAccountData = &primitives.AccountData{
		Free: sc.NewU128(1),
	}

	return newTransfer(moduleId, mockStoredMap, testConstants, mockMutator)
}
