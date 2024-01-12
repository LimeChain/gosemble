package balances

import (
	"bytes"
	"errors"
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
			NewMultiAddressId(constants.OneAccountId)
	toAddress = primitives.
			NewMultiAddressId(constants.TwoAccountId)
	argsBytes = sc.NewVaryingData(primitives.MultiAddress{}, sc.Compact{Number: sc.U128{}}).Bytes()

	callTransferArgsBytes = sc.NewVaryingData(primitives.MultiAddress{}, sc.Compact{Number: sc.U128{}}).Bytes()
)

func Test_Call_Transfer_New(t *testing.T) {
	target := setupCallTransfer()
	expected := callTransfer{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionTransferIndex,
			Arguments:  sc.NewVaryingData(primitives.MultiAddress{}, sc.Compact{Number: sc.U128{}}),
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
	expectedBuffer := bytes.NewBuffer(append([]byte{moduleId, functionTransferIndex}, callTransferArgsBytes...))
	buf := &bytes.Buffer{}

	err := target.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_Transfer_Bytes(t *testing.T) {
	expected := append([]byte{moduleId, functionTransferIndex}, callTransferArgsBytes...)

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

	assert.Equal(t, primitives.PaysYes, target.PaysFee(baseWeight))
}

func Test_Call_Transfer_Dispatch_Success(t *testing.T) {
	target := setupCallTransfer()

	fromAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)

	_, dispatchErr := target.
		Dispatch(primitives.NewRawOriginSigned(fromAddressId), sc.NewVaryingData(fromAddress, sc.ToCompact(targetValue)))

	assert.Nil(t, dispatchErr)
}

func Test_Call_Transfer_Dispatch_BadOrigin(t *testing.T) {
	target := setupCallTransfer()

	_, dispatchErr := target.Dispatch(primitives.NewRawOriginNone(), sc.NewVaryingData(toAddress, sc.ToCompact(targetValue)))

	assert.Equal(t, primitives.NewDispatchErrorBadOrigin(), dispatchErr)
}

func Test_Call_Transfer_Dispatch_InvalidArgs(t *testing.T) {
	target := setupCallTransfer()

	_, dispatchErr := target.Dispatch(primitives.NewRawOriginNone(), sc.NewVaryingData(toAddress, sc.NewU64(0)))

	assert.Equal(t, errors.New("invalid compact value when dispatching call transfer"), dispatchErr)

	_, dispatchErr = target.Dispatch(primitives.NewRawOriginNone(), sc.NewVaryingData(toAddress, sc.Compact{}))

	assert.Equal(t, errors.New("invalid compact number field when dispatching call transfer"), dispatchErr)
}

func Test_Call_Transfer_Dispatch_CannotLookup(t *testing.T) {
	target := setupCallTransfer()

	fromAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)

	_, dispatchErr := target.Dispatch(
		primitives.NewRawOriginSigned(fromAddressId),
		sc.NewVaryingData(primitives.NewMultiAddress20(primitives.Address20{}), sc.ToCompact(targetValue)),
	)

	assert.Equal(t, primitives.NewDispatchErrorCannotLookup(), dispatchErr)
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

	fromAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)

	result := target.transfer(primitives.NewRawOriginSigned(fromAddressId), fromAddress, targetValue)

	assert.Nil(t, result)
}

func Test_transfer_InvalidOrigin(t *testing.T) {
	target := setupTransfer()

	result := target.transfer(primitives.NewRawOriginRoot(), toAddress, targetValue)

	assert.Equal(t, primitives.NewDispatchErrorBadOrigin(), result)
}

func Test_transfer_InvalidLookup(t *testing.T) {
	target := setupTransfer()

	fromAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)

	result := target.
		transfer(primitives.NewRawOriginSigned(fromAddressId), primitives.NewMultiAddress20(primitives.Address20{}), targetValue)

	assert.Equal(t, primitives.NewDispatchErrorCannotLookup(), result)
}

func Test_transfer_trans_Success(t *testing.T) {
	target := setupTransfer()

	fromAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)

	toAddressId, err := toAddress.AsAccountId()
	assert.Nil(t, err)

	mockMutator.On(
		"tryMutateAccountWithDust",
		toAddressId,
		mockTypeMutateAccountDataBool,
	).Return(sc.Empty{}, nil)
	mockStoredMap.On(
		"DepositEvent",
		newEventTransfer(moduleId, fromAddressId, toAddressId, targetValue),
	).Return()

	result := target.trans(fromAddressId, toAddressId, targetValue, primitives.ExistenceRequirementKeepAlive)

	assert.Nil(t, result)
	mockMutator.AssertCalled(t,
		"tryMutateAccountWithDust",
		toAddressId,
		mockTypeMutateAccountDataBool,
	)
	mockStoredMap.AssertCalled(t,
		"DepositEvent",
		newEventTransfer(moduleId, fromAddressId, toAddressId, targetValue),
	)
}

func Test_transfer_trans_ZeroValue(t *testing.T) {
	target := setupTransfer()

	fromAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)

	toAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)

	result := target.trans(fromAddressId, toAddressId, sc.NewU128(0), primitives.ExistenceRequirementAllowDeath)

	assert.Nil(t, result)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_transfer_trans_EqualFromTo(t *testing.T) {
	target := setupTransfer()

	fromAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)

	result := target.trans(fromAddressId, fromAddressId, targetValue, primitives.ExistenceRequirementAllowDeath)

	assert.Nil(t, result)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_transfer_trans_MutateAccountWithDust_Fails(t *testing.T) {
	target := setupTransfer()
	expectedErr := primitives.NewDispatchErrorBadOrigin()

	fromAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)

	toAddressId, err := toAddress.AsAccountId()
	assert.Nil(t, err)

	mockMutator.On(
		"tryMutateAccountWithDust",
		toAddressId,
		mockTypeMutateAccountDataBool,
	).Return(sc.Empty{}, expectedErr)

	result := target.trans(fromAddressId, toAddressId, targetValue, primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, expectedErr, result)
	mockMutator.AssertCalled(t,
		"tryMutateAccountWithDust",
		toAddressId,
		mockTypeMutateAccountDataBool,
	)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_transfer_sanityChecks_Success(t *testing.T) {
	target := setupTransfer()

	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	mockMutator.On("ensureCanWithdraw", targetAddressId, targetValue, primitives.ReasonsAll, sc.NewU128(0)).Return(nil)
	mockStoredMap.On("CanDecProviders", targetAddressId).Return(true, nil)

	result, err := target.sanityChecks(targetAddressId, fromAccountData, toAccountData, targetValue, primitives.ExistenceRequirementAllowDeath)

	assert.Nil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, sc.NewU128(0), fromAccountData.Free)
	assert.Equal(t, sc.NewU128(6), toAccountData.Free)
	mockMutator.AssertCalled(t, "ensureCanWithdraw", targetAddressId, targetValue, primitives.ReasonsAll, sc.NewU128(0))
	mockStoredMap.AssertCalled(t, "CanDecProviders", targetAddressId)
}

func Test_transfer_sanityChecks_InsufficientBalance(t *testing.T) {
	target := setupTransfer()
	expectedErr := primitives.NewDispatchErrorModule(primitives.CustomModuleError{
		Index:   moduleId,
		Err:     sc.U32(ErrorInsufficientBalance),
		Message: sc.NewOption[sc.Str](nil),
	})

	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	_, err = target.sanityChecks(targetAddressId, fromAccountData, toAccountData, sc.NewU128(6), primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, sc.NewU128(5), fromAccountData.Free)
	assert.Equal(t, sc.NewU128(1), toAccountData.Free)
	mockMutator.AssertNotCalled(t, "ensureCanWithdraw", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "CanDecProviders", mock.Anything)
}

func Test_transfer_sanityChecks_ArithmeticOverflow(t *testing.T) {
	target := setupTransfer()
	expectedErr := primitives.NewDispatchErrorArithmetic(primitives.NewArithmeticErrorOverflow())
	toAccountData.Free = sc.MaxU128()

	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	_, err = target.sanityChecks(targetAddressId, fromAccountData, toAccountData, sc.NewU128(1), primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, sc.NewU128(4), fromAccountData.Free)
	assert.Equal(t, sc.MaxU128(), toAccountData.Free)
	mockMutator.AssertNotCalled(t, "ensureCanWithdraw", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "CanDecProviders", mock.Anything)
}

func Test_transfer_sanityChecks_ExistentialDeposit(t *testing.T) {
	target := setupTransfer()
	expectedErr := primitives.NewDispatchErrorModule(primitives.CustomModuleError{
		Index:   moduleId,
		Err:     sc.U32(ErrorExistentialDeposit),
		Message: sc.NewOption[sc.Str](nil),
	})
	toAccountData.Free = sc.NewU128(0)

	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	_, err = target.sanityChecks(targetAddressId, fromAccountData, toAccountData, sc.NewU128(0), primitives.ExistenceRequirementKeepAlive)

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, sc.NewU128(5), fromAccountData.Free)
	assert.Equal(t, sc.NewU128(0), toAccountData.Free)
	mockMutator.AssertNotCalled(t, "ensureCanWithdraw", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "CanDecProviders", mock.Anything)
}

func Test_transfer_sanityChecks_CannotWithdraw(t *testing.T) {
	target := setupTransfer()
	expectedErr := primitives.NewDispatchErrorCannotLookup()

	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	mockMutator.On("ensureCanWithdraw", targetAddressId, targetValue, primitives.ReasonsAll, sc.NewU128(0)).Return(expectedErr)

	_, err = target.sanityChecks(targetAddressId, fromAccountData, toAccountData, targetValue, primitives.ExistenceRequirementAllowDeath)

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, sc.NewU128(0), fromAccountData.Free)
	assert.Equal(t, sc.NewU128(6), toAccountData.Free)
	mockMutator.AssertCalled(t, "ensureCanWithdraw", targetAddressId, targetValue, primitives.ReasonsAll, sc.NewU128(0))
}

func Test_transfer_sanityChecks_KeepAlive(t *testing.T) {
	target := setupTransfer()
	expectedErr := primitives.NewDispatchErrorModule(primitives.CustomModuleError{
		Index:   moduleId,
		Err:     sc.U32(ErrorKeepAlive),
		Message: sc.NewOption[sc.Str](nil),
	})

	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	mockMutator.On("ensureCanWithdraw", targetAddressId, targetValue, primitives.ReasonsAll, sc.NewU128(0)).Return(nil)
	mockStoredMap.On("CanDecProviders", targetAddressId).Return(false, nil)

	_, err = target.sanityChecks(targetAddressId, fromAccountData, toAccountData, targetValue, primitives.ExistenceRequirementAllowDeath)

	assert.Equal(t, expectedErr, err)
	assert.Equal(t, sc.NewU128(0), fromAccountData.Free)
	assert.Equal(t, sc.NewU128(6), toAccountData.Free)
	mockMutator.AssertCalled(t, "ensureCanWithdraw", targetAddressId, targetValue, primitives.ReasonsAll, sc.NewU128(0))
	mockStoredMap.AssertCalled(t, "CanDecProviders", targetAddressId)
}

func Test_transfer_sanityChecks_CanDecProviders_Error(t *testing.T) {
	target := setupTransfer()

	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	mockErr := errors.New("err")
	expectedErr := primitives.NewDispatchErrorOther(sc.Str(mockErr.Error()))

	mockMutator.On("ensureCanWithdraw", targetAddressId, targetValue, primitives.ReasonsAll, sc.NewU128(0)).Return(nil)
	mockStoredMap.On("CanDecProviders", targetAddressId).Return(true, mockErr)

	_, err = target.sanityChecks(targetAddressId, fromAccountData, toAccountData, targetValue, primitives.ExistenceRequirementAllowDeath)

	assert.Equal(t, expectedErr, err)
	mockMutator.AssertCalled(t, "ensureCanWithdraw", targetAddressId, targetValue, primitives.ReasonsAll, sc.NewU128(0))
	mockStoredMap.AssertCalled(t, "CanDecProviders", targetAddressId)
}

func Test_transfer_reducibleBalance_NotKeepAlive(t *testing.T) {
	target := setupTransfer()

	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	mockStoredMap.On("Get", targetAddressId).Return(accountInfo, nil)
	mockStoredMap.On("CanDecProviders", targetAddressId).Return(true, nil)

	result, err := target.reducibleBalance(targetAddressId, false)
	assert.Nil(t, err)

	assert.Equal(t, accountInfo.Data.Free, result)
	mockStoredMap.AssertCalled(t, "Get", targetAddressId)
	mockStoredMap.AssertCalled(t, "CanDecProviders", targetAddressId)
}

func Test_transfer_reducibleBalance_KeepAlive(t *testing.T) {
	target := setupTransfer()

	targetAddressId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	mockStoredMap.On("Get", targetAddressId).Return(accountInfo, nil)
	mockStoredMap.On("CanDecProviders", targetAddressId).Return(false, nil)

	result, err := target.reducibleBalance(targetAddressId, true)
	assert.Nil(t, err)

	assert.Equal(t, accountInfo.Data.Free.Sub(existentialDeposit), result)
	mockStoredMap.AssertCalled(t, "Get", targetAddressId)
	mockStoredMap.AssertCalled(t, "CanDecProviders", targetAddressId)
}

func setupCallTransfer() primitives.Call {
	mockStoredMap = new(mocks.StoredMap)
	mockMutator = new(mockAccountMutator)

	fromAccountData = &primitives.AccountData{
		Free: sc.NewU128(5),
	}

	toAccountData = &primitives.AccountData{
		Free: sc.NewU128(1),
	}

	return newCallTransfer(moduleId, functionTransferIndex, mockStoredMap, testConstants, mockMutator)
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
