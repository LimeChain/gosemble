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
	baseWeight                    = primitives.WeightFromParts(124, 123)
	targetAddress                 = primitives.NewMultiAddressId(constants.ZeroAccountId)
	targetValue                   = sc.NewU128(5)
	mockTypeMutateAccountDataBool = mock.AnythingOfType("func(*types.AccountData, bool) (goscale.Encodable, error)")
	argsBytesCallForceFree        = sc.NewVaryingData(primitives.MultiAddress{}, sc.U128{}).Bytes()
	mockStoredMap                 *mocks.StoredMap
	errPanic                      = errors.New("panic")
)

func Test_Call_ForceFree_new(t *testing.T) {
	target := setupCallForceFree()
	expected := callForceFree{
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionForceFreeIndex,
			Arguments:  sc.NewVaryingData(primitives.MultiAddress{}, sc.U128{}),
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

func Test_Call_ForceFree_DecodeArgs(t *testing.T) {
	amount := sc.NewU128(5)
	buf := bytes.NewBuffer(append(targetAddress.Bytes(), amount.Bytes()...))

	target := setupCallForceFree()
	call, err := target.DecodeArgs(buf)
	assert.Nil(t, err)

	assert.Equal(t, sc.NewVaryingData(targetAddress, amount), call.Args())
}

func Test_Call_ForceFree_Encode(t *testing.T) {
	target := setupCallForceFree()
	expectedBuffer := bytes.NewBuffer(append([]byte{moduleId, functionForceFreeIndex}, argsBytesCallForceFree...))
	buf := &bytes.Buffer{}

	err := target.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_ForceFree_Bytes(t *testing.T) {
	expected := append([]byte{moduleId, functionForceFreeIndex}, argsBytesCallForceFree...)

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
	call, err := target.DecodeArgs(buf)
	assert.Nil(t, err)

	buf.Reset()
	call.Encode(buf)

	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_ForceFree_BaseWeight(t *testing.T) {
	target := setupCallForceFree()

	assert.Equal(t, callForceFreeWeight(dbWeight), target.BaseWeight())
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

	assert.Equal(t, primitives.PaysYes, target.PaysFee(baseWeight))
}

func Test_Call_ForceFree_Dispatch_Success(t *testing.T) {
	target := setupCallForceFree()
	actual := sc.NewU128(1)
	mutateResult := actual
	targetAddressAccId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)
	event := newEventUnreserved(moduleId, targetAddressAccId, actual)

	mockStoredMap.On("Get", targetAddressAccId).Return(accountInfo, nil)
	mockMutator.On("tryMutateAccount",
		targetAddressAccId,
		mockTypeMutateAccountDataBool).
		Return(mutateResult, nil)
	mockStoredMap.On("DepositEvent", event)

	_, err = target.Dispatch(primitives.NewRawOriginRoot(), sc.NewVaryingData(targetAddress, targetValue))
	assert.Nil(t, err)
	mockStoredMap.AssertCalled(t, "Get", targetAddressAccId)
	mockMutator.AssertCalled(t,
		"tryMutateAccount",
		targetAddressAccId,
		mockTypeMutateAccountDataBool,
	)
	mockStoredMap.AssertCalled(t, "DepositEvent", event)
}

func Test_Call_ForceFree_Dispatch_InvalidOrigin(t *testing.T) {
	target := setupCallForceFree()

	_, dispatchErr := target.Dispatch(primitives.NewRawOriginNone(), sc.NewVaryingData(targetAddress, targetValue))

	assert.Equal(t, primitives.NewDispatchErrorBadOrigin(), dispatchErr)
	mockStoredMap.AssertNotCalled(t, "Get", mock.Anything)
	mockMutator.AssertNotCalled(t, "tryMutateAccount", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceFree_Dispatch_InvalidArgs(t *testing.T) {
	target := setupCallForceFree()

	_, dispatchErr := target.Dispatch(primitives.NewRawOriginNone(), sc.NewVaryingData(targetAddress, sc.NewU64(0)))

	assert.Equal(t, errors.New("invalid amount value when dispatching call force free"), dispatchErr)
	mockStoredMap.AssertNotCalled(t, "Get", mock.Anything)
	mockMutator.AssertNotCalled(t, "tryMutateAccount", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceFree_Dispatch_InvalidLookup(t *testing.T) {
	target := setupCallForceFree()

	_, dispatchErr := target.Dispatch(primitives.NewRawOriginRoot(), sc.NewVaryingData(primitives.NewMultiAddress20(primitives.Address20{}), targetValue))

	assert.Equal(t, primitives.NewDispatchErrorCannotLookup(), dispatchErr)
	targetAddressAccId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)
	mockStoredMap.AssertNotCalled(t, "Get", targetAddressAccId)
	mockMutator.AssertNotCalled(t, "tryMutateAccount", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceFree_Dispatch_ZeroBalance(t *testing.T) {
	target := setupCallForceFree()

	_, dispatchErr := target.Dispatch(primitives.NewRawOriginRoot(), sc.NewVaryingData(targetAddress, constants.Zero))

	assert.Nil(t, dispatchErr)
	targetAddressAccId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)
	mockStoredMap.AssertNotCalled(t, "Get", targetAddressAccId)
	mockMutator.AssertNotCalled(t, "tryMutateAccount", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceFree_Dispatch_ZeroTotalStorageBalance(t *testing.T) {
	target := setupCallForceFree()
	accountInfo := primitives.AccountInfo{Data: primitives.AccountData{}}

	targetAddressAccId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)
	mockStoredMap.On("Get", targetAddressAccId).Return(accountInfo, nil)

	_, dispatchErr := target.Dispatch(primitives.NewRawOriginRoot(), sc.NewVaryingData(targetAddress, targetValue))

	assert.Nil(t, dispatchErr)
	mockStoredMap.AssertCalled(t, "Get", targetAddressAccId)
	mockMutator.AssertNotCalled(t, "tryMutateAccount", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceFree_Dispatch_Other(t *testing.T) {
	target := setupCallForceFree()
	accountInfo := primitives.AccountInfo{Data: primitives.AccountData{}}

	targetAddressAccId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)

	mockStoredMap.On("Get", targetAddressAccId).Return(accountInfo, errPanic)

	_, dispatchErr := target.Dispatch(primitives.NewRawOriginRoot(), sc.NewVaryingData(targetAddress, targetValue))

	assert.Equal(t, primitives.NewDispatchErrorOther(sc.Str(errPanic.Error())), dispatchErr)
	mockStoredMap.AssertCalled(t, "Get", targetAddressAccId)
	mockMutator.AssertNotCalled(t, "tryMutateAccount", mock.Anything, mock.Anything)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_Call_ForceFree_Dispatch_Mutation_Fails(t *testing.T) {
	target := setupCallForceFree()

	expectedErr := primitives.NewDispatchErrorOther(sc.Str(errPanic.Error()))
	targetAddressAccId, err := targetAddress.AsAccountId()
	assert.Nil(t, err)
	mockStoredMap.On("Get", targetAddressAccId).Return(accountInfo, nil)
	mockMutator.On("tryMutateAccount",
		targetAddressAccId,
		mockTypeMutateAccountDataBool,
	).Return(sc.Empty{}, errPanic)

	_, dispatchErr := target.Dispatch(primitives.NewRawOriginRoot(), sc.NewVaryingData(targetAddress, targetValue))

	assert.Equal(t, expectedErr, dispatchErr)
	mockStoredMap.AssertCalled(t, "Get", targetAddressAccId)
	mockMutator.AssertCalled(t,
		"tryMutateAccount",
		targetAddressAccId,
		mockTypeMutateAccountDataBool,
	)
	mockStoredMap.AssertNotCalled(t, "DepositEvent", mock.Anything)
}

func Test_removeReserveAndFree(t *testing.T) {
	value := sc.NewU128(4)
	accountData := &primitives.AccountData{
		Free:     sc.NewU128(1),
		Reserved: sc.NewU128(10),
	}
	expectedResult := value

	result := removeReserveAndFree(accountData, value)

	assert.Equal(t, expectedResult, result)
	assert.Equal(t, sc.NewU128(6), accountData.Reserved)
	assert.Equal(t, sc.NewU128(5), accountData.Free)
}

func setupCallForceFree() primitives.Call {
	mockStoredMap = new(mocks.StoredMap)
	mockMutator = new(mockAccountMutator)

	return newCallForceFree(moduleId, sc.U8(functionForceFreeIndex), mockStoredMap, testConstants, mockMutator, logger)
}
