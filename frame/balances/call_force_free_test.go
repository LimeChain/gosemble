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

const (
	moduleId = 5
)

var (
	accountInfo = primitives.AccountInfo{
		Data: primitives.AccountData{
			Free:       sc.NewU128(4),
			Reserved:   primitives.Balance{},
			MiscFrozen: primitives.Balance{},
			FeeFrozen:  primitives.Balance{},
		},
	}
	dbWeight = primitives.RuntimeDbWeight{
		Read:  1,
		Write: 2,
	}
	baseWeight    = primitives.WeightFromParts(124, 123)
	targetAddress = primitives.NewMultiAddress32(constants.ZeroAddress)
	targetValue   = sc.NewU128(5)
	mockStoredMap *mocks.MockStoredMap
)

func Test_Call_ForceFree_new(t *testing.T) {
	target := setupCallForceFree()
	expected := callForceFree{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionForceFreeIndex,
		},
		dbWeight:  dbWeight,
		storedMap: mockStoredMap,
	}

	assert.Equal(t, expected, target)
}

func Test_Call_ForceFree_DecodeArgs(t *testing.T) {
	amount := sc.NewU128(5)
	buf := bytes.NewBuffer(append(targetAddress.Bytes(), amount.Bytes()...))

	target := setupCallForceFree()
	call := target.DecodeArgs(buf)

	assert.Equal(t, sc.NewVaryingData(targetAddress, amount), call.Args())
}

func Test_Call_ForceFree_Encode(t *testing.T) {
	target := setupCallForceFree()
	expectedBuffer := bytes.NewBuffer([]byte{moduleId, functionForceFreeIndex})
	buf := &bytes.Buffer{}

	target.Encode(buf)

	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_ForceFree_Bytes(t *testing.T) {
	expected := []byte{moduleId, functionForceFreeIndex}

	target := setupCallForceFree()

	assert.Equal(t, expected, target.Bytes())
}

func Test_Call_ForceFree_ModuleIndex(t *testing.T) {
	target := setupCallForceFree()

	assert.Equal(t, sc.U8(moduleId), target.ModuleIndex())
}

func TestCall_ForceFree_FunctionIndex(t *testing.T) {
	target := setupCallForceFree()

	assert.Equal(t, sc.U8(functionForceFreeIndex), target.FunctionIndex())
}

func Test_Call_ForceFree_EncodeWithArgs(t *testing.T) {
	expectedBuffer := bytes.NewBuffer([]byte{moduleId, functionForceFreeIndex})
	bArgs := append(targetAddress.Bytes(), targetValue.Bytes()...)
	expectedBuffer.Write(bArgs)

	buf := bytes.NewBuffer(bArgs)

	target := setupCallForceFree()
	call := target.DecodeArgs(buf)

	buf.Reset()
	call.Encode(buf)

	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_ForceFree_BaseWeight(t *testing.T) {
	target := setupCallForceFree()

	assert.Equal(t, primitives.WeightFromParts(17_029_003, 3593), target.BaseWeight())
}

func Test_Call_ForceFree_WeighData(t *testing.T) {
	target := setupCallForceFree()
	assert.Equal(t, primitives.WeightFromParts(124, 0), target.WeighData(baseWeight))
}

func Test_Call_ForceFree_ClassifyDispatch(t *testing.T) {
	target := setupCallForceFree()

	assert.Equal(t, primitives.NewDispatchClassNormal(), target.ClassifyDispatch(baseWeight))
}

func Test_Call_ForceFree_PaysFee(t *testing.T) {
	target := setupCallForceFree()

	assert.Equal(t, primitives.NewPaysYes(), target.PaysFee(baseWeight))
}

func Test_Call_ForceFree_Dispatch_Success(t *testing.T) {
	target := setupCallForceFree()
	actual := sc.NewU128(1)
	mutateResult := sc.Result[sc.Encodable]{HasError: false, Value: actual}
	event := newEventUnreserved(moduleId, targetAddress.AsAddress32().FixedSequence, actual)

	mockStoredMap.On("Get", targetAddress.AsAddress32().FixedSequence).Return(accountInfo)
	mockStoredMap.On("Mutate",
		targetAddress.
			AsAddress32(),
		mock.AnythingOfType("func(*types.AccountInfo) goscale.Result[github.com/LimeChain/goscale.Encodable]")).
		Return(mutateResult)
	mockStoredMap.On("DepositEvent", event)

	result := target.Dispatch(primitives.NewRawOriginRoot(), sc.NewVaryingData(targetAddress, targetValue))

	assert.Equal(t, primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{}, result)
	mockStoredMap.AssertCalled(t, "Get", targetAddress.AsAddress32().FixedSequence)
	mockStoredMap.AssertCalled(t,
		"Mutate",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountInfo) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	)
	mockStoredMap.AssertCalled(t, "DepositEvent", event)
}

func Test_Call_ForceFree_Dispatch_Fails(t *testing.T) {
	target := setupCallForceFree()
	expected := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: true,
		Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
			Error: primitives.NewDispatchErrorBadOrigin(),
		},
	}

	result := target.Dispatch(primitives.NewRawOriginNone(), sc.NewVaryingData(targetAddress, targetValue))

	assert.Equal(t, expected, result)
	mockStoredMap.AssertNotCalled(t, "Get", mock.Anything)
	mockStoredMap.AssertNotCalled(t, "Mutate", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceFree_forceFree_Success(t *testing.T) {
	target := setupCallForceFree()
	actual := sc.NewU128(1)
	mutateResult := sc.Result[sc.Encodable]{HasError: false, Value: actual}
	event := newEventUnreserved(moduleId, targetAddress.AsAddress32().FixedSequence, actual)

	mockStoredMap.On("Get", targetAddress.AsAddress32().FixedSequence).Return(accountInfo)
	mockStoredMap.On(
		"Mutate",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountInfo) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	).Return(mutateResult)
	mockStoredMap.On("DepositEvent", event)

	result := target.forceFree(primitives.NewRawOriginRoot(), targetAddress, targetValue)

	assert.Equal(t, sc.VaryingData(nil), result)
	mockStoredMap.AssertCalled(t, "Get", targetAddress.AsAddress32().FixedSequence)
	mockStoredMap.AssertCalled(t,
		"Mutate",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountInfo) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	)
	mockStoredMap.AssertCalled(t, "DepositEvent", event)
}

func Test_Call_ForceFree_forceFree_InvalidOrigin(t *testing.T) {
	target := setupCallForceFree()

	result := target.forceFree(primitives.NewRawOriginNone(), targetAddress, targetValue)

	assert.Equal(t, primitives.NewDispatchErrorBadOrigin(), result)
	mockStoredMap.AssertNotCalled(t, "Get", targetAddress.AsAddress32().FixedSequence)
	mockStoredMap.AssertNotCalled(t, "Mutate", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceFree_forceFree_InvalidLookup(t *testing.T) {
	target := setupCallForceFree()

	result := target.forceFree(primitives.NewRawOriginRoot(), primitives.NewMultiAddress20(primitives.Address20{}), targetValue)

	assert.Equal(t, primitives.NewDispatchErrorCannotLookup(), result)
	mockStoredMap.AssertNotCalled(t, "Get", targetAddress.AsAddress32().FixedSequence)
	mockStoredMap.AssertNotCalled(t, "Mutate", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceFree_force_Success(t *testing.T) {
	target := setupCallForceFree()
	actual := sc.NewU128(1)
	expectedResult := targetValue.Sub(actual)
	mutateResult := sc.Result[sc.Encodable]{HasError: false, Value: actual}
	event := newEventUnreserved(moduleId, targetAddress.AsAddress32().FixedSequence, actual)

	mockStoredMap.On("Get", targetAddress.AsAddress32().FixedSequence).Return(accountInfo)
	mockStoredMap.On(
		"Mutate",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountInfo) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	).Return(mutateResult)
	mockStoredMap.On("DepositEvent", event)

	result := target.force(targetAddress.AsAddress32(), targetValue)

	assert.Equal(t, expectedResult, result)
	mockStoredMap.AssertCalled(t, "Get", targetAddress.AsAddress32().FixedSequence)
	mockStoredMap.AssertCalled(t, "Mutate",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountInfo) goscale.Result[github.com/LimeChain/goscale.Encodable]"))
	mockStoredMap.AssertCalled(t, "DepositEvent", event)
}

func Test_Call_ForceFree_force_ZeroBalance(t *testing.T) {
	target := setupCallForceFree()

	result := target.force(targetAddress.AsAddress32(), constants.Zero)

	assert.Equal(t, constants.Zero, result)
	mockStoredMap.AssertNotCalled(t, "Get", targetAddress.AsAddress32().FixedSequence)
	mockStoredMap.AssertNotCalled(t, "Mutate", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceFree_force_FullBalance(t *testing.T) {
	target := setupCallForceFree()
	accountInfo := primitives.AccountInfo{Data: primitives.AccountData{}}

	mockStoredMap.On("Get", targetAddress.AsAddress32().FixedSequence).Return(accountInfo)

	result := target.force(targetAddress.AsAddress32(), targetValue)

	assert.Equal(t, targetValue, result)
	mockStoredMap.AssertCalled(t, "Get", targetAddress.AsAddress32().FixedSequence)
	mockStoredMap.AssertNotCalled(t, "Mutate", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceFree_force_Mutate_Fails(t *testing.T) {
	target := setupCallForceFree()
	mutateResult := sc.Result[sc.Encodable]{HasError: true}

	mockStoredMap.On("Get", targetAddress.AsAddress32().FixedSequence).Return(accountInfo)
	mockStoredMap.On("Mutate",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountInfo) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	).Return(mutateResult)

	result := target.force(targetAddress.AsAddress32(), targetValue)

	assert.Equal(t, targetValue, result)
	mockStoredMap.AssertCalled(t, "Get", targetAddress.AsAddress32().FixedSequence)
	mockStoredMap.AssertCalled(t,
		"Mutate",
		targetAddress.AsAddress32(),
		mock.AnythingOfType("func(*types.AccountInfo) goscale.Result[github.com/LimeChain/goscale.Encodable]"),
	)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_removeReserveAndFree(t *testing.T) {
	value := sc.NewU128(4)
	accountInfo := &primitives.AccountInfo{
		Data: primitives.AccountData{
			Free:     sc.NewU128(1),
			Reserved: sc.NewU128(10),
		},
	}
	expectedResult := sc.Result[sc.Encodable]{HasError: false, Value: value}

	result := removeReserveAndFree(accountInfo, value)

	assert.Equal(t, expectedResult, result)
	assert.Equal(t, sc.NewU128(6), accountInfo.Data.Reserved)
	assert.Equal(t, sc.NewU128(5), accountInfo.Data.Free)
}

func setupCallForceFree() callForceFree {
	mockStoredMap = new(mocks.MockStoredMap)
	return newCallForceFree(moduleId, sc.U8(functionForceFreeIndex), dbWeight, mockStoredMap).(callForceFree)
}
