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

func Test_Call_TransferKeepAlive_new(t *testing.T) {
	target := setupCallTransferKeepAlive()
	expected := callTransferKeepAlive{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionTransferKeepAliveIndex,
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
	expectedBuffer := bytes.NewBuffer([]byte{moduleId, functionTransferKeepAliveIndex})
	buf := &bytes.Buffer{}

	target.Encode(buf)

	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_TransferKeepAlive_Bytes(t *testing.T) {
	expected := []byte{moduleId, functionTransferKeepAliveIndex}

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

	assert.Equal(t, primitives.NewPaysYes(), target.PaysFee(baseWeight))
}

func Test_Call_TransferKeepAlive_Dispatch_Success(t *testing.T) {
	target := setupCallTransferKeepAlive()
	expect := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: false,
		Ok:       primitives.PostDispatchInfo{},
	}

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
			targetValue,
		),
	).Return()

	result := target.Dispatch(primitives.NewRawOriginSigned(fromAddress.AsAddress32()), sc.NewVaryingData(toAddress, sc.ToCompact(targetValue)))

	assert.Equal(t, expect, result)
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
			targetValue,
		),
	)
}

func Test_Call_TransferKeepAlive_Dispatch_BadOrigin(t *testing.T) {
	target := setupCallTransferKeepAlive()
	expect := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: true,
		Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
			Error: primitives.NewDispatchErrorBadOrigin(),
		},
	}

	result := target.Dispatch(
		primitives.NewRawOriginNone(),
		sc.NewVaryingData(fromAddress, sc.ToCompact(targetValue)),
	)

	assert.Equal(t, expect, result)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_TransferKeepAlive_Dispatch_CannotLookup(t *testing.T) {
	target := setupCallTransferKeepAlive()
	expect := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: true,
		Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
			Error: primitives.NewDispatchErrorCannotLookup(),
		},
	}

	result := target.
		Dispatch(
			primitives.NewRawOriginSigned(fromAddress.AsAddress32()),
			sc.NewVaryingData(primitives.NewMultiAddress20(primitives.Address20{}), sc.ToCompact(targetValue)),
		)

	assert.Equal(t, expect, result)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func setupCallTransferKeepAlive() callTransferKeepAlive {
	mockStoredMap = new(mocks.StoredMap)
	mockMutator = new(mockAccountMutator)

	return newCallTransferKeepAlive(moduleId, functionTransferKeepAliveIndex, mockStoredMap, testConstants, mockMutator).(callTransferKeepAlive)
}
