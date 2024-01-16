package balances

import (
	"bytes"
	"errors"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	callForceTransferArgsBytes = sc.NewVaryingData(primitives.MultiAddress{}, primitives.MultiAddress{}, sc.Compact{Number: sc.U128{}}).Bytes()
)

func Test_Call_ForceTransfer_new(t *testing.T) {
	target := setupCallForceTransfer()
	expected := callForceTransfer{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionForceTransferIndex,
			Arguments:  sc.NewVaryingData(primitives.MultiAddress{}, primitives.MultiAddress{}, sc.Compact{Number: sc.U128{}}),
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

func Test_Call_ForceTransfer_DecodeArgs(t *testing.T) {
	amount := sc.ToCompact(sc.NewU128(1))
	buf := &bytes.Buffer{}
	buf.Write(fromAddress.Bytes())
	buf.Write(toAddress.Bytes())
	buf.Write(amount.Bytes())

	target := setupCallForceTransfer()
	call, err := target.DecodeArgs(buf)
	assert.Nil(t, err)

	assert.Equal(t, sc.NewVaryingData(fromAddress, toAddress, amount), call.Args())
}

func Test_Call_ForceTransfer_Encode(t *testing.T) {
	target := setupCallForceTransfer()
	expectedBuffer := bytes.NewBuffer(append([]byte{moduleId, functionForceTransferIndex}, callForceTransferArgsBytes...))
	buf := &bytes.Buffer{}

	err := target.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_ForceTransfer_Bytes(t *testing.T) {
	expected := append([]byte{moduleId, functionForceTransferIndex}, callForceTransferArgsBytes...)

	target := setupCallForceTransfer()

	assert.Equal(t, expected, target.Bytes())
}

func Test_Call_ForceTransfer_ModuleIndex(t *testing.T) {
	target := setupCallForceTransfer()

	assert.Equal(t, sc.U8(moduleId), target.ModuleIndex())
}

func Test_Call_ForceTransfer_FunctionIndex(t *testing.T) {
	target := setupCallForceTransfer()

	assert.Equal(t, sc.U8(functionForceTransferIndex), target.FunctionIndex())
}

func Test_Call_ForceTransfer_BaseWeight(t *testing.T) {
	target := setupCallForceTransfer()

	assert.Equal(t, primitives.WeightFromParts(40_360_006, 6196), target.BaseWeight())
}

func Test_Call_ForceTransfer_WeighData(t *testing.T) {
	target := setupCallForceTransfer()
	assert.Equal(t, primitives.WeightFromParts(124, 0), target.WeighData(baseWeight))
}

func Test_Call_ForceTransfer_ClassifyDispatch(t *testing.T) {
	target := setupCallForceTransfer()

	assert.Equal(t, primitives.NewDispatchClassNormal(), target.ClassifyDispatch(baseWeight))
}

func Test_Call_ForceTransfer_PaysFee(t *testing.T) {
	target := setupCallForceTransfer()

	assert.Equal(t, primitives.PaysYes, target.PaysFee(baseWeight))
}

func Test_Call_ForceTransfer_Dispatch_Success(t *testing.T) {
	target := setupCallForceTransfer()

	fromAddressAccId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)

	toAddressAccId, err := toAddress.AsAccountId()
	assert.Nil(t, err)

	mockMutator.On(
		"tryMutateAccountWithDust",
		toAddressAccId,
		mockTypeMutateAccountDataBool,
	).
		Return(sc.Empty{}, nil)
	mockStoredMap.On(
		"DepositEvent",
		newEventTransfer(moduleId, fromAddressAccId, toAddressAccId, targetValue),
	).
		Return()

	_, dispatchErr := target.Dispatch(primitives.NewRawOriginRoot(), sc.NewVaryingData(fromAddress, toAddress, sc.ToCompact(targetValue)))

	assert.Nil(t, dispatchErr)
	mockMutator.AssertCalled(t,
		"tryMutateAccountWithDust",
		toAddressAccId,
		mockTypeMutateAccountDataBool,
	)
	mockStoredMap.AssertCalled(t,
		"DepositEvent",
		newEventTransfer(moduleId, fromAddressAccId, toAddressAccId, targetValue),
	)
}

func Test_Call_ForceTransfer_Dispatch_InvalidBadOrigin(t *testing.T) {
	target := setupCallForceTransfer()

	_, dispatchErr := target.Dispatch(
		primitives.NewRawOriginNone(),
		sc.NewVaryingData(fromAddress, toAddress, sc.ToCompact(targetValue)))

	assert.Equal(t, primitives.NewDispatchErrorBadOrigin(), dispatchErr)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceTransfer_Dispatch_InvalidArg_InvalidCompact(t *testing.T) {
	target := setupCallForceTransfer()

	_, dispatchErr := target.Dispatch(
		primitives.NewRawOriginNone(),
		sc.NewVaryingData(fromAddress, toAddress, sc.NewU128(0)))

	assert.Equal(t, errors.New("invalid Compact value when dispatching call_force_transfer"), dispatchErr)

	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)

}

func Test_Call_ForceTransfer_Dispatch_InvalidArg_InvalidCompactNumber(t *testing.T) {
	target := setupCallForceTransfer()

	_, dispatchErr := target.Dispatch(
		primitives.NewRawOriginNone(),
		sc.NewVaryingData(fromAddress, toAddress, sc.Compact{}))

	assert.Equal(t, errors.New("invalid Compact field number when dispatching call_force_transfer"), dispatchErr)

	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceTransfer_Dispatch_CannotLookup_Source(t *testing.T) {
	target := setupCallForceTransfer()

	_, dispatchErr := target.Dispatch(
		primitives.NewRawOriginRoot(),
		sc.NewVaryingData(primitives.NewMultiAddress20(primitives.Address20{}), toAddress, sc.ToCompact(targetValue)),
	)

	assert.Equal(t, primitives.NewDispatchErrorCannotLookup(), dispatchErr)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceTransfer_Dispatch_CannotLookup_Dest(t *testing.T) {
	target := setupCallForceTransfer()

	_, dispatchErr := target.Dispatch(
		primitives.NewRawOriginRoot(),
		sc.NewVaryingData(fromAddress, primitives.NewMultiAddress20(primitives.Address20{}), sc.ToCompact(targetValue)),
	)

	assert.Equal(t, primitives.NewDispatchErrorCannotLookup(), dispatchErr)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func setupCallForceTransfer() primitives.Call {
	mockStoredMap = new(mocks.StoredMap)
	mockMutator = new(mockAccountMutator)

	return newCallForceTransfer(moduleId, functionForceTransferIndex, mockStoredMap, testConstants, mockMutator)
}
