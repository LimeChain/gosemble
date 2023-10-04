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

func Test_Call_TransferAll_new(t *testing.T) {
	target := setupCallTransferAll()
	expected := callTransferAll{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionTransferAllIndex,
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

func Test_Call_TransferAll_DecodeArgs(t *testing.T) {
	keepAlive := sc.Bool(true)
	buf := bytes.NewBuffer(append(targetAddress.Bytes(), keepAlive.Bytes()...))

	target := setupCallTransferAll()
	call := target.DecodeArgs(buf)

	assert.Equal(t, sc.NewVaryingData(targetAddress, keepAlive), call.Args())
}

func Test_Call_TransferAll_Encode(t *testing.T) {
	target := setupCallTransferAll()
	expectedBuffer := bytes.NewBuffer([]byte{moduleId, functionTransferAllIndex})
	buf := &bytes.Buffer{}

	target.Encode(buf)

	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_TransferAll_Bytes(t *testing.T) {
	expected := []byte{moduleId, functionTransferAllIndex}

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

	assert.Equal(t, primitives.WeightFromParts(35_121_003, 3593), target.BaseWeight())
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

	assert.Equal(t, primitives.NewPaysYes(), target.PaysFee(baseWeight))
}

func Test_Call_TransferAll_Dispatch_Success(t *testing.T) {
	target := setupCallTransferAll()
	expect := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: false,
		Ok:       primitives.PostDispatchInfo{},
	}

	mockStoredMap.On("Get", fromAddress.AsAddress32().FixedSequence).Return(accountInfo)
	mockStoredMap.On("CanDecProviders", fromAddress.AsAddress32()).Return(true)
	mockMutator.On("tryMutateAccountWithDust",
		toAddress.AsAddress32(),
		mockTypeMutateAccountDataBool,
	).Return(sc.Result[sc.Encodable]{})
	mockStoredMap.On(
		"DepositEvent",
		newEventTransfer(
			moduleId,
			fromAddress.AsAddress32().
				FixedSequence,
			toAddress.AsAddress32().FixedSequence,
			accountInfo.Data.Free.Sub(sc.NewU128(1)),
		),
	).
		Return()

	result := target.Dispatch(primitives.NewRawOriginSigned(fromAddress.AsAddress32()), sc.NewVaryingData(toAddress, sc.Bool(true)))

	assert.Equal(t, expect, result)
	mockStoredMap.AssertCalled(t, "Get", fromAddress.AsAddress32().FixedSequence)
	mockStoredMap.AssertCalled(t, "CanDecProviders", fromAddress.AsAddress32())
	mockMutator.AssertCalled(t,
		"tryMutateAccountWithDust",
		toAddress.AsAddress32(),
		mockTypeMutateAccountDataBool,
	)
	mockStoredMap.AssertCalled(t,
		"DepositEvent",
		newEventTransfer(
			moduleId,
			fromAddress.AsAddress32().FixedSequence,
			toAddress.AsAddress32().FixedSequence,
			accountInfo.Data.Free.Sub(sc.NewU128(1)),
		),
	)
}

func Test_Call_TransferAll_Dispatch_Fails(t *testing.T) {
	target := setupCallTransferAll()
	expect := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: true,
		Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
			Error: primitives.NewDispatchErrorBadOrigin(),
		},
	}

	result := target.Dispatch(primitives.NewRawOriginNone(), sc.NewVaryingData(fromAddress, sc.Bool(true)))

	assert.Equal(t, expect, result)
	mockStoredMap.AssertNotCalled(t, "Get", mock.Anything)
	mockStoredMap.AssertNotCalled(t, "CanDecProviders", mock.Anything)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_TransferAll_transferAll_Success(t *testing.T) {
	target := setupCallTransferAll()

	mockStoredMap.On("Get", fromAddress.AsAddress32().FixedSequence).Return(accountInfo)
	mockStoredMap.On("CanDecProviders", fromAddress.AsAddress32()).Return(true)
	mockMutator.On(
		"tryMutateAccountWithDust",
		toAddress.AsAddress32(),
		mockTypeMutateAccountDataBool,
	).Return(sc.Result[sc.Encodable]{})
	mockStoredMap.On(
		"DepositEvent",
		newEventTransfer(
			moduleId,
			fromAddress.AsAddress32().FixedSequence,
			toAddress.AsAddress32().FixedSequence,
			accountInfo.Data.Free.Sub(sc.NewU128(1)),
		),
	).Return()

	result := target.transferAll(primitives.NewRawOriginSigned(fromAddress.AsAddress32()), toAddress, true)

	assert.Equal(t, sc.VaryingData(nil), result)
	mockStoredMap.AssertCalled(t, "Get", fromAddress.AsAddress32().FixedSequence)
	mockStoredMap.AssertCalled(t, "CanDecProviders", fromAddress.AsAddress32())
	mockMutator.AssertCalled(t,
		"tryMutateAccountWithDust",
		toAddress.AsAddress32(),
		mockTypeMutateAccountDataBool,
	)
	mockStoredMap.AssertCalled(t,
		"DepositEvent",
		newEventTransfer(
			moduleId,
			fromAddress.AsAddress32().FixedSequence,
			toAddress.AsAddress32().FixedSequence,
			accountInfo.Data.Free.Sub(sc.NewU128(1)),
		),
	)
}

func Test_Call_TransferAll_transferAll_InvalidOrigin(t *testing.T) {
	target := setupCallTransferAll()
	expect := primitives.NewDispatchErrorBadOrigin()

	result := target.transferAll(primitives.NewRawOriginRoot(), toAddress, true)

	assert.Equal(t, expect, result)
	mockStoredMap.AssertNotCalled(t, "Get", mock.Anything)
	mockStoredMap.AssertNotCalled(t, "CanDecProviders", mock.Anything)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_TransferAll_transferAll_InvalidLookup(t *testing.T) {
	target := setupCallTransferAll()
	mockStoredMap.On("Get", fromAddress.AsAddress32().FixedSequence).Return(accountInfo)
	mockStoredMap.On("CanDecProviders", fromAddress.AsAddress32()).Return(true)

	result := target.
		transferAll(primitives.NewRawOriginSigned(fromAddress.AsAddress32()), primitives.NewMultiAddress20(primitives.Address20{}), true)

	assert.Equal(t, primitives.NewDispatchErrorCannotLookup(), result)
	mockStoredMap.AssertCalled(t, "Get", fromAddress.AsAddress32().FixedSequence)
	mockStoredMap.AssertCalled(t, "CanDecProviders", fromAddress.AsAddress32())
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_TransferAll_transferAll_AllowDeath(t *testing.T) {
	target := setupCallTransferAll()

	mockStoredMap.On("Get", fromAddress.AsAddress32().FixedSequence).Return(accountInfo)
	mockStoredMap.On("CanDecProviders", fromAddress.AsAddress32()).Return(true)
	mockMutator.On(
		"tryMutateAccountWithDust",
		toAddress.AsAddress32(),
		mockTypeMutateAccountDataBool,
	).Return(sc.Result[sc.Encodable]{})
	mockStoredMap.On(
		"DepositEvent",
		newEventTransfer(
			moduleId,
			fromAddress.AsAddress32().FixedSequence,
			toAddress.AsAddress32().FixedSequence,
			accountInfo.Data.Free,
		),
	).Return()

	result := target.transferAll(primitives.NewRawOriginSigned(fromAddress.AsAddress32()), toAddress, false)

	assert.Equal(t, sc.VaryingData(nil), result)
	mockStoredMap.AssertCalled(t, "Get", fromAddress.AsAddress32().FixedSequence)
	mockStoredMap.AssertCalled(t, "CanDecProviders", fromAddress.AsAddress32())
	mockMutator.AssertCalled(t,
		"tryMutateAccountWithDust",
		toAddress.AsAddress32(),
		mockTypeMutateAccountDataBool,
	)
	mockStoredMap.AssertCalled(t,
		"DepositEvent",
		newEventTransfer(
			moduleId,
			fromAddress.AsAddress32().FixedSequence,
			toAddress.AsAddress32().FixedSequence,
			accountInfo.Data.Free,
		),
	)
}

func setupCallTransferAll() callTransferAll {
	mockStoredMap = new(mocks.StoredMap)
	mockMutator = new(mockAccountMutator)

	return newCallTransferAll(moduleId, functionTransferAllIndex, mockStoredMap, testConstants, mockMutator).(callTransferAll)
}
