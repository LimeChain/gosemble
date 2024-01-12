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
	transferKeepAliveArgsBytes = sc.NewVaryingData(primitives.MultiAddress{}, sc.Compact{Number: sc.U128{}}).Bytes()
)

func Test_Call_TransferKeepAlive_new(t *testing.T) {
	target := setupCallTransferKeepAlive()
	expected := callTransferKeepAlive{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionTransferKeepAliveIndex,
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

func Test_Call_TransferKeepAlive_DecodeArgs(t *testing.T) {
	amount := sc.ToCompact(sc.NewU128(1))
	buf := bytes.NewBuffer(append(targetAddress.Bytes(), amount.Bytes()...))

	target := setupCallTransferKeepAlive()
	call, err := target.DecodeArgs(buf)
	assert.Nil(t, err)

	assert.Equal(t, sc.NewVaryingData(targetAddress, amount), call.Args())
}

func Test_Call_TransferKeepAlive_Encode(t *testing.T) {
	target := setupCallTransferKeepAlive()
	expectedBuffer := bytes.NewBuffer(append([]byte{moduleId, functionTransferKeepAliveIndex}, transferKeepAliveArgsBytes...))
	buf := &bytes.Buffer{}

	err := target.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_TransferKeepAlive_Bytes(t *testing.T) {
	expected := append([]byte{moduleId, functionTransferKeepAliveIndex}, transferKeepAliveArgsBytes...)

	target := setupCallTransferKeepAlive()

	assert.Equal(t, expected, target.Bytes())
}

func Test_Call_TransferKeepAlive_ModuleIndex(t *testing.T) {
	target := setupCallTransferKeepAlive()

	assert.Equal(t, sc.U8(moduleId), target.ModuleIndex())
}

func Test_Call_TransferKeepAlive_FunctionIndex(t *testing.T) {
	target := setupCallTransferKeepAlive()

	assert.Equal(t, sc.U8(functionTransferKeepAliveIndex), target.FunctionIndex())
}

func Test_Call_TransferKeepAlive_BaseWeight(t *testing.T) {
	target := setupCallTransferKeepAlive()

	assert.Equal(t, primitives.WeightFromParts(49_250_003, 3593), target.BaseWeight())
}

func Test_Call_TransferKeepAlive_WeighData(t *testing.T) {
	target := setupCallTransferKeepAlive()
	assert.Equal(t, primitives.WeightFromParts(124, 0), target.WeighData(baseWeight))
}

func Test_Call_TransferKeepAlive_ClassifyDispatch(t *testing.T) {
	target := setupCallTransferKeepAlive()

	assert.Equal(t, primitives.NewDispatchClassNormal(), target.ClassifyDispatch(baseWeight))
}

func Test_Call_TransferKeepAlive_PaysFee(t *testing.T) {
	target := setupCallTransferKeepAlive()

	assert.Equal(t, primitives.PaysYes, target.PaysFee(baseWeight))
}

func Test_Call_TransferKeepAlive_Dispatch_Success(t *testing.T) {
	target := setupCallTransferKeepAlive()

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
		newEventTransfer(
			moduleId,
			fromAddressId,
			toAddressId,
			targetValue,
		),
	).Return()

	_, dispatchErr := target.Dispatch(primitives.NewRawOriginSigned(fromAddressId), sc.NewVaryingData(toAddress, sc.ToCompact(targetValue)))

	assert.Nil(t, dispatchErr)
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
			targetValue,
		),
	)
}

func Test_Call_TransferKeepAlive_Dispatch_BadOrigin(t *testing.T) {
	target := setupCallTransferKeepAlive()

	_, dispatchErr := target.Dispatch(
		primitives.NewRawOriginNone(),
		sc.NewVaryingData(fromAddress, sc.ToCompact(targetValue)),
	)

	assert.Equal(t, primitives.NewDispatchErrorBadOrigin(), dispatchErr)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_TransferKeepAlive_Dispatch_InvalidArgs(t *testing.T) {
	target := setupCallTransferKeepAlive()

	_, dispatchErr := target.Dispatch(
		primitives.NewRawOriginNone(),
		sc.NewVaryingData(fromAddress, sc.NewU64(0)),
	)

	assert.Equal(t, errors.New("invalid compact value when dispatching call transfer keep alive"), dispatchErr)

	_, dispatchErr = target.Dispatch(
		primitives.NewRawOriginNone(),
		sc.NewVaryingData(fromAddress, sc.Compact{}),
	)

	assert.Equal(t, errors.New("invalid compact number field when dispatching call transfer keep alive"), dispatchErr)

	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_TransferKeepAlive_Dispatch_CannotLookup(t *testing.T) {
	target := setupCallTransferKeepAlive()

	fromAddressId, err := fromAddress.AsAccountId()
	assert.Nil(t, err)

	_, dispatchErr := target.
		Dispatch(
			primitives.NewRawOriginSigned(fromAddressId),
			sc.NewVaryingData(primitives.NewMultiAddress20(primitives.Address20{}), sc.ToCompact(targetValue)),
		)

	assert.Equal(t, primitives.NewDispatchErrorCannotLookup(), dispatchErr)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func setupCallTransferKeepAlive() primitives.Call {
	mockStoredMap = new(mocks.StoredMap)
	mockMutator = new(mockAccountMutator)

	return newCallTransferKeepAlive(moduleId, functionTransferKeepAliveIndex, mockStoredMap, testConstants, mockMutator)
}
