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

func Test_Call_ForceTransfer_new(t *testing.T) {
	target := setupCallForceTransfer()
	expected := callForceTransfer{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionForceTransferIndex,
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
	call := target.DecodeArgs(buf)

	assert.Equal(t, sc.NewVaryingData(fromAddress, toAddress, amount), call.Args())
}

func Test_Call_ForceTransfer_Encode(t *testing.T) {
	target := setupCallForceTransfer()
	expectedBuffer := bytes.NewBuffer([]byte{moduleId, functionForceTransferIndex})
	buf := &bytes.Buffer{}

	target.Encode(buf)

	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_ForceTransfer_Bytes(t *testing.T) {
	expected := []byte{moduleId, functionForceTransferIndex}

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

	assert.Equal(t, primitives.NewPaysYes(), target.PaysFee(baseWeight))
}

func Test_Call_ForceTransfer_Dispatch_Success(t *testing.T) {
	target := setupCallForceTransfer()
	expect := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: false,
		Ok:       primitives.PostDispatchInfo{},
	}

	mockMutator.On(
		"tryMutateAccountWithDust",
		toAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	).
		Return(sc.Result[sc.Encodable]{})
	mockStoredMap.On(
		"DepositEvent",
		newEventTransfer(moduleId, fromAddress.AsAddress32().FixedSequence, toAddress.AsAddress32().FixedSequence, targetValue),
	).
		Return()

	result := target.Dispatch(primitives.NewRawOriginRoot(), sc.NewVaryingData(fromAddress, toAddress, sc.ToCompact(targetValue)))

	assert.Equal(t, expect, result)
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

func Test_Call_ForceTransfer_Dispatch_Fails(t *testing.T) {
	target := setupCallForceTransfer()
	expect := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: true,
		Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
			Error: primitives.NewDispatchErrorBadOrigin(),
		},
	}

	result := target.Dispatch(primitives.NewRawOriginNone(), sc.NewVaryingData(fromAddress, toAddress, sc.ToCompact(targetValue)))

	assert.Equal(t, expect, result)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceTransfer_forceTransfer_Success(t *testing.T) {
	target := setupCallForceTransfer()

	mockMutator.On(
		"tryMutateAccountWithDust",
		toAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	).
		Return(sc.Result[sc.Encodable]{})
	mockStoredMap.On(
		"DepositEvent",
		newEventTransfer(moduleId, fromAddress.AsAddress32().FixedSequence, toAddress.AsAddress32().FixedSequence, targetValue),
	).
		Return()

	result := target.forceTransfer(primitives.NewRawOriginRoot(), fromAddress, toAddress, targetValue)

	assert.Equal(t, sc.VaryingData(nil), result)
	mockMutator.AssertCalled(t,
		"tryMutateAccountWithDust",
		toAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountData, bool) goscale.Result[github.com/LimeChain/goscale.Encodable]"))
	mockStoredMap.AssertCalled(t,
		"DepositEvent",
		newEventTransfer(moduleId, fromAddress.AsAddress32().FixedSequence, toAddress.AsAddress32().FixedSequence, targetValue))

}

func Test_Call_ForceTransfer_forceTransfer_InvalidOrigin(t *testing.T) {
	target := setupCallForceTransfer()
	expect := primitives.NewDispatchErrorBadOrigin()

	result := target.forceTransfer(primitives.NewRawOriginNone(), fromAddress, toAddress, targetValue)

	assert.Equal(t, expect, result)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceTransfer_forceTransfer_From_Lookup(t *testing.T) {
	target := setupCallForceTransfer()

	result := target.forceTransfer(primitives.NewRawOriginRoot(), primitives.NewMultiAddress20(primitives.Address20{}), toAddress, targetValue)

	assert.Equal(t, primitives.NewDispatchErrorCannotLookup(), result)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceTransfer_forceTransfer_Dest_Lookup(t *testing.T) {
	target := setupCallForceTransfer()

	result := target.forceTransfer(primitives.NewRawOriginRoot(), fromAddress, primitives.NewMultiAddress20(primitives.Address20{}), targetValue)

	assert.Equal(t, primitives.NewDispatchErrorCannotLookup(), result)
	mockMutator.AssertNotCalled(t, "tryMutateAccountWithDust", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func setupCallForceTransfer() callForceTransfer {
	mockStoredMap = new(mocks.MockStoredMap)
	mockMutator = new(mockAccountMutator)

	return newCallForceTransfer(moduleId, functionForceTransferIndex, mockStoredMap, testConstants, mockMutator).(callForceTransfer)
}
