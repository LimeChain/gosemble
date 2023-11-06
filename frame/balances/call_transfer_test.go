package balances

import (
	"bytes"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
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
			NewMultiAddressId(constants.OneAddressAccountId)
	toAddress = primitives.
			NewMultiAddressId(constants.TwoAddressAccountId)
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
	call, err := target.DecodeArgs(buf)
	assert.Nil(t, err)

	assert.Equal(t, sc.NewVaryingData(targetAddress, amount), call.Args())
}

func Test_Call_Transfer_Encode(t *testing.T) {
	target := setupCallTransfer()
	expectedBuffer := bytes.NewBuffer([]byte{moduleId, functionTransferIndex})
	buf := &bytes.Buffer{}

	err := target.Encode(buf)

	assert.NoError(t, err)
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
		Dispatch(primitives.NewRawOriginSigned(fromAddress.AsAccountId()), sc.NewVaryingData(fromAddress, sc.ToCompact(targetValue)))

	assert.Equal(t, expected, result)
}

func Test_Call_Transfer_Dispatch_BadOrigin(t *testing.T) {
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

func Test_Call_Transfer_Dispatch_CannotLookup(t *testing.T) {
	target := setupCallTransfer()
	expected := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: true,
		Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
			Error: primitives.NewDispatchErrorCannotLookup(),
		},
	}

	result := target.Dispatch(
		primitives.NewRawOriginSigned(fromAddress.AsAccountId()),
		sc.NewVaryingData(primitives.NewMultiAddress20(primitives.Address20{}), sc.ToCompact(targetValue)),
	)

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

	result := target.transfer(primitives.NewRawOriginSigned(fromAddress.AsAccountId()), fromAddress, targetValue)

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
		transfer(primitives.NewRawOriginSigned(fromAddress.AsAccountId()), primitives.NewMultiAddress20(primitives.Address20{}), targetValue)

	assert.Equal(t, primitives.NewDispatchErrorCannotLookup(), result)
}

func Test_transfer_trans_Success(t *testing.T) {
	target := setupTransfer()

	mockMutator.On(
		"tryMutateAccountWithDust",
		toAddress.AsAccountId(),
		mockTypeMutateAccountDataBool,
	).Return(sc.Result[sc.Encodable]{})
	mockStoredMap.On(
		"DepositEvent",
		newEventTransfer(moduleId, fromAddress.AsAccountId(), toAddress.AsAccountId(), targetValue),
	).Return()

	result := target.trans(fromAddress.AsAccountId(), toAddress.AsAccountId(), targetValue, primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, sc.VaryingData(nil), result)
	mockMutator.AssertCalled(t,
		"tryMutateAccountWithDust",
		toAddress.AsAccountId(),
		mockTypeMutateAccountDataBool,
	)
	mockStoredMap.AssertCalled(t,
		"DepositEvent",
		newEventTransfer(moduleId, fromAddress.AsAccountId(), toAddress.AsAccountId(), targetValue),
	)
}

func Test_transfer_trans_ZeroValue(t *testing.T) {
	target := setupTransfer()

	result := target.trans(fromAddress.AsAccountId(), toAddress.AsAccountId(), sc.NewU128(0), primitives.ExistenceRequirementAllowDeath)

	assert.Equal(t, sc.VaryingData(nil), result)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_transfer_trans_EqualFromTo(t *testing.T) {
	target := setupTransfer()

	result := target.trans(fromAddress.AsAccountId(), fromAddress.AsAccountId(), targetValue, primitives.ExistenceRequirementAllowDeath)

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
		toAddress.AsAccountId(),
		mockTypeMutateAccountDataBool,
	).Return(error)

	result := target.trans(fromAddress.AsAccountId(), toAddress.AsAccountId(), targetValue, primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, expectdError, result)
	mockMutator.AssertCalled(t,
		"tryMutateAccountWithDust",
		toAddress.AsAccountId(),
		mockTypeMutateAccountDataBool,
	)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_transfer_sanityChecks_Success(t *testing.T) {
	target := setupTransfer()
	expected := sc.Result[sc.Encodable]{}

	mockMutator.On("ensureCanWithdraw", targetAddress.AsAccountId(), targetValue, primitives.ReasonsAll, sc.NewU128(0)).Return(sc.VaryingData(nil))
	mockStoredMap.On("CanDecProviders", targetAddress.AsAccountId()).Return(true, nil)

	result := target.sanityChecks(targetAddress.AsAccountId(), fromAccountData, toAccountData, targetValue, primitives.ExistenceRequirementAllowDeath)

	assert.Equal(t, expected, result)
	assert.Equal(t, sc.NewU128(0), fromAccountData.Free)
	assert.Equal(t, sc.NewU128(6), toAccountData.Free)
	mockMutator.AssertCalled(t, "ensureCanWithdraw", targetAddress.AsAccountId(), targetValue, primitives.ReasonsAll, sc.NewU128(0))
	mockStoredMap.AssertCalled(t, "CanDecProviders", targetAddress.AsAccountId())
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

	result := target.sanityChecks(targetAddress.AsAccountId(), fromAccountData, toAccountData, sc.NewU128(6), primitives.ExistenceRequirementKeepAlive)

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

	result := target.sanityChecks(targetAddress.AsAccountId(), fromAccountData, toAccountData, sc.NewU128(1), primitives.ExistenceRequirementKeepAlive)

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

	result := target.sanityChecks(targetAddress.AsAccountId(), fromAccountData, toAccountData, sc.NewU128(0), primitives.ExistenceRequirementKeepAlive)

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
	mockMutator.On("ensureCanWithdraw", targetAddress.AsAccountId(), targetValue, primitives.ReasonsAll, sc.NewU128(0)).Return(expectedError)

	result := target.sanityChecks(targetAddress.AsAccountId(), fromAccountData, toAccountData, targetValue, primitives.ExistenceRequirementAllowDeath)

	assert.Equal(t, expected, result)
	assert.Equal(t, sc.NewU128(0), fromAccountData.Free)
	assert.Equal(t, sc.NewU128(6), toAccountData.Free)
	mockMutator.AssertCalled(t, "ensureCanWithdraw", targetAddress.AsAccountId(), targetValue, primitives.ReasonsAll, sc.NewU128(0))
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
	mockMutator.On("ensureCanWithdraw", targetAddress.AsAccountId(), targetValue, primitives.ReasonsAll, sc.NewU128(0)).Return(sc.VaryingData(nil))
	mockStoredMap.On("CanDecProviders", targetAddress.AsAccountId()).Return(false, nil)

	result := target.sanityChecks(targetAddress.AsAccountId(), fromAccountData, toAccountData, targetValue, primitives.ExistenceRequirementAllowDeath)

	assert.Equal(t, expected, result)
	assert.Equal(t, sc.NewU128(0), fromAccountData.Free)
	assert.Equal(t, sc.NewU128(6), toAccountData.Free)
	mockMutator.AssertCalled(t, "ensureCanWithdraw", targetAddress.AsAccountId(), targetValue, primitives.ReasonsAll, sc.NewU128(0))
	mockStoredMap.AssertCalled(t, "CanDecProviders", targetAddress.AsAccountId())
}

func Test_transfer_reducibleBalance_NotKeepAlive(t *testing.T) {
	target := setupTransfer()
	mockStoredMap.On("Get", targetAddress.AsAccountId()).Return(accountInfo, nil)
	mockStoredMap.On("CanDecProviders", targetAddress.AsAccountId()).Return(true, nil)

	result, err := target.reducibleBalance(targetAddress.AsAccountId(), false)
	assert.Nil(t, err)

	assert.Equal(t, accountInfo.Data.Free, result)
	mockStoredMap.AssertCalled(t, "Get", targetAddress.AsAccountId())
	mockStoredMap.AssertCalled(t, "CanDecProviders", targetAddress.AsAccountId())
}

func Test_transfer_reducibleBalance_KeepAlive(t *testing.T) {
	target := setupTransfer()
	mockStoredMap.On("Get", targetAddress.AsAccountId()).Return(accountInfo, nil)
	mockStoredMap.On("CanDecProviders", targetAddress.AsAccountId()).Return(false, nil)

	result, err := target.reducibleBalance(targetAddress.AsAccountId(), true)
	assert.Nil(t, err)

	assert.Equal(t, accountInfo.Data.Free.Sub(existentialDeposit), result)
	mockStoredMap.AssertCalled(t, "Get", targetAddress.AsAccountId())
	mockStoredMap.AssertCalled(t, "CanDecProviders", targetAddress.AsAccountId())
}

func setupCallTransfer() callTransfer {
	mockStoredMap = new(mocks.StoredMap)
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
	mockStoredMap = new(mocks.StoredMap)
	mockMutator = new(mockAccountMutator)

	fromAccountData = &primitives.AccountData{
		Free: sc.NewU128(5),
	}

	toAccountData = &primitives.AccountData{
		Free: sc.NewU128(1),
	}

	return newTransfer(moduleId, mockStoredMap, testConstants, mockMutator)
}
