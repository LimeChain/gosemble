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
	transferAllArgsBytes = sc.NewVaryingData(primitives.MultiAddress{}, sc.Bool(true)).Bytes()
)

func Test_Call_TransferAll_new(t *testing.T) {
	target := setupCallTransferAll()
	expected := callTransferAll{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionTransferAllIndex,
			Arguments:  sc.NewVaryingData(primitives.MultiAddress{}, sc.Bool(true)),
		},
		transfer: transfer{
			moduleId:       moduleId,
			storedMap:      mockStoredMap,
			constants:      testConstants,
			accountMutator: mockMutator,
		},
		logger: logger,
	}

	assert.Equal(t, expected, target)
}

func Test_Call_TransferAll_DecodeArgs(t *testing.T) {
	keepAlive := sc.Bool(true)
	buf := bytes.NewBuffer(append(targetAddress.Bytes(), keepAlive.Bytes()...))

	target := setupCallTransferAll()
	call, err := target.DecodeArgs(buf)
	assert.Nil(t, err)

	assert.Equal(t, sc.NewVaryingData(targetAddress, keepAlive), call.Args())
}

func Test_Call_TransferAll_Encode(t *testing.T) {
	target := setupCallTransferAll()
	expectedBuffer := bytes.NewBuffer(append([]byte{moduleId, functionTransferAllIndex}, transferAllArgsBytes...))
	buf := &bytes.Buffer{}

	err := target.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_TransferAll_Bytes(t *testing.T) {
	expected := append([]byte{moduleId, functionTransferAllIndex}, transferAllArgsBytes...)

	target := setupCallTransferAll()

	assert.Equal(t, expected, target.Bytes())
}

func Test_Call_TransferAll_ModuleIndex(t *testing.T) {
	target := setupCallTransferAll()

	assert.Equal(t, sc.U8(moduleId), target.ModuleIndex())
}

func Test_Call_TransferAll_FunctionIndex(t *testing.T) {
	target := setupCallTransferAll()

	assert.Equal(t, sc.U8(functionTransferAllIndex), target.FunctionIndex())
}

func Test_Call_TransferAll_BaseWeight(t *testing.T) {
	target := setupCallTransferAll()

	assert.Equal(t, callTransferAllWeight(dbWeight), target.BaseWeight())
}

func Test_Call_TransferAll_WeighData(t *testing.T) {
	target := setupCallTransferAll()
	assert.Equal(t, primitives.WeightFromParts(124, 0), target.WeighData(baseWeight))
}

func Test_Call_TransferAll_ClassifyDispatch(t *testing.T) {
	target := setupCallTransferAll()

	assert.Equal(t, primitives.NewDispatchClassNormal(), target.ClassifyDispatch(baseWeight))
}

func Test_Call_TransferAll_PaysFee(t *testing.T) {
	target := setupCallTransferAll()

	assert.Equal(t, primitives.PaysYes, target.PaysFee(baseWeight))
}

func Test_Call_TransferAll_Dispatch_Success(t *testing.T) {
	target := setupCallTransferAll()

	fromAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)

	toAddressId, err := toAddress.AsAccountId()
	assert.Nil(t, err)

	mockStoredMap.On("Get", fromAddressId).Return(accountInfo, nil)
	mockStoredMap.On("CanDecProviders", fromAddressId).Return(true, nil)
	mockMutator.On("tryMutateAccountWithDust",
		toAddressId,
		mockTypeMutateAccountDataBool,
	).Return(sc.Empty{}, nil)
	mockStoredMap.On(
		"DepositEvent",
		newEventTransfer(
			moduleId,
			fromAddressId,
			toAddressId,
			accountInfo.Data.Free.Sub(sc.NewU128(1)),
		),
	).
		Return()

	_, dispatchErr := target.Dispatch(
		primitives.NewRawOriginSigned(fromAddressId),
		sc.NewVaryingData(toAddress, sc.Bool(true)),
	)

	assert.Nil(t, dispatchErr)
	mockStoredMap.AssertCalled(t, "Get", fromAddressId)
	mockStoredMap.AssertCalled(t, "CanDecProviders", fromAddressId)
	mockMutator.AssertCalled(t,
		"tryMutateAccountWithDust",
		toAddressId,
		mockTypeMutateAccountDataBool,
	)
	mockStoredMap.AssertCalled(t,
		"DepositEvent",
		newEventTransfer(
			moduleId,
			fromAddressId,
			toAddressId,
			accountInfo.Data.Free.Sub(sc.NewU128(1)),
		),
	)
}

func Test_Call_TransferAll_Dispatch_BadOrigin(t *testing.T) {
	target := setupCallTransferAll()

	_, dispatchErr := target.Dispatch(
		primitives.NewRawOriginNone(),
		sc.NewVaryingData(fromAddress, sc.Bool(true)),
	)

	assert.Equal(t, primitives.NewDispatchErrorBadOrigin(), dispatchErr)
	mockStoredMap.AssertNotCalled(t, "Get", mock.Anything)
	mockStoredMap.AssertNotCalled(t, "CanDecProviders", mock.Anything)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_TransferAll_Dispatch_CannotLookup(t *testing.T) {
	target := setupCallTransferAll()

	fromAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)
	mockStoredMap.On("Get", fromAddressId).Return(accountInfo, nil)
	mockStoredMap.On("CanDecProviders", fromAddressId).Return(true, nil)

	_, dispatchErr := target.Dispatch(
		primitives.NewRawOriginSigned(fromAddressId),
		sc.NewVaryingData(primitives.NewMultiAddress20(primitives.Address20{}), sc.Bool(true)),
	)

	assert.Equal(t, primitives.NewDispatchErrorCannotLookup(), dispatchErr)
	mockStoredMap.AssertCalled(t, "Get", fromAddressId)
	mockStoredMap.AssertCalled(t, "CanDecProviders", fromAddressId)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_TransferAll_Dispatch_AllowDeath(t *testing.T) {
	target := setupCallTransferAll()

	fromAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)

	toAddressId, err := toAddress.AsAccountId()
	assert.Nil(t, err)

	mockStoredMap.On("Get", fromAddressId).Return(accountInfo, nil)
	mockStoredMap.On("CanDecProviders", fromAddressId).Return(true, nil)
	mockMutator.On(
		"tryMutateAccountWithDust",
		toAddressId,
		mockTypeMutateAccountDataBool,
	).Return(sc.Empty{}, nil)
	mockStoredMap.On(
		"DepositEvent",
		newEventTransfer(
			moduleId,
			fromAddressId,
			toAddressId,
			accountInfo.Data.Free,
		),
	).Return()

	_, dispatchErr := target.Dispatch(
		primitives.NewRawOriginSigned(fromAddressId),
		sc.NewVaryingData(toAddress, sc.Bool(false)))

	assert.Nil(t, dispatchErr)
	mockStoredMap.AssertCalled(t, "Get", fromAddressId)
	mockStoredMap.AssertCalled(t, "CanDecProviders", fromAddressId)
	mockMutator.AssertCalled(t,
		"tryMutateAccountWithDust",
		toAddressId,
		mockTypeMutateAccountDataBool,
	)
	mockStoredMap.AssertCalled(t,
		"DepositEvent",
		newEventTransfer(
			moduleId,
			fromAddressId,
			toAddressId,
			accountInfo.Data.Free,
		),
	)
}

func setupCallTransferAll() primitives.Call {
	mockStoredMap = new(mocks.StoredMap)
	mockMutator = new(mockAccountMutator)

	return newCallTransferAll(moduleId, functionTransferAllIndex, mockStoredMap, testConstants, mockMutator, logger)
}
