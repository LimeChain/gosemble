package timestamp

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
)

const (
	moduleId = 0
)

var (
	c          = newConstants(constants.RocksDbWeight, 5)
	baseWeight = primitives.WeightFromParts(124, 123)
	origin     = primitives.NewRawOriginNone()
	now        = sc.U64(time.Unix(1, 0).UnixMilli())

	argsBytesCallSet = sc.NewVaryingData(sc.Compact{Number: sc.NewU64(0)}).Bytes()

	mockOnTimestampSet   *mocks.OnTimestampSet
	mockStorageNow       *mocks.StorageValue[sc.U64]
	mockStorageDidUpdate *mocks.StorageValue[sc.Bool]
	mockStorage          *storage
)

func Test_Call_Set_NewSetCall(t *testing.T) {
	target := setUpCallSet()
	expected := callSet{
		storage:        mockStorage,
		onTimestampSet: mockOnTimestampSet,
		constants:      c,
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionSetIndex,
			Arguments:  sc.NewVaryingData(sc.Compact{Number: sc.NewU64(0)}),
		},
	}

	assert.Equal(t, expected, target)
}

func Test_Call_Set_NewSetCallWithArgs(t *testing.T) {
	expected := callSet{
		storage:        nil,
		onTimestampSet: nil,
		constants:      nil,
		Callable: primitives.Callable{
			ModuleId:   moduleId,
			FunctionId: functionSetIndex,
			Arguments:  sc.NewVaryingData(),
		},
	}

	res := newCallSetWithArgs(moduleId, functionSetIndex, sc.NewVaryingData())

	assert.Equal(t, expected, res)
}

func Test_Call_Set_DecodeArgs(t *testing.T) {
	target := setUpCallSet()
	compact := sc.Compact{Number: sc.U64(5)}
	buf := bytes.NewBuffer(compact.Bytes())

	call, err := target.DecodeArgs(buf)
	assert.Nil(t, err)

	assert.Equal(t, sc.NewVaryingData(compact), call.Args())
}

func Test_Call_Set_DecodeArgs_Fails(t *testing.T) {
	target := setUpCallSet()

	call, err := target.DecodeArgs(&bytes.Buffer{})
	assert.Equal(t, io.EOF, err)
	assert.Nil(t, call)
}

func Test_Call_Set_Encode(t *testing.T) {
	target := setUpCallSet()
	expectedBuffer := bytes.NewBuffer(append([]byte{moduleId, functionSetIndex}, argsBytesCallSet...))
	buf := &bytes.Buffer{}

	err := target.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_Set_EncodeWithArgs(t *testing.T) {
	target := setUpCallSet()

	expectedBuf := bytes.NewBuffer(append([]byte{moduleId, functionSetIndex}, argsBytesCallSet...))

	buf := bytes.NewBuffer(argsBytesCallSet)

	call, err := target.DecodeArgs(buf)
	assert.Nil(t, err)

	buf.Reset()
	err = call.Encode(buf)

	assert.NoError(t, err)
	assert.Equal(t, expectedBuf, buf)
}

func Test_Call_Set_Bytes(t *testing.T) {
	target := setUpCallSet()
	expected := append([]byte{moduleId, functionSetIndex}, argsBytesCallSet...)

	assert.Equal(t, expected, target.Bytes())
}

func Test_Call_Set_ModuleIndex(t *testing.T) {
	testCases := []sc.U8{
		moduleId,
		1,
		2,
		3,
	}

	for _, tc := range testCases {
		call := newCallSetWithArgs(tc, functionSetIndex, sc.NewVaryingData())

		assert.Equal(t, tc, call.ModuleIndex())
	}
}

func Test_Call_Set_FunctionIndex(t *testing.T) {
	testCases := []sc.U8{
		functionSetIndex,
		1,
		2,
		3,
	}

	for _, tc := range testCases {
		call := newCallSetWithArgs(moduleId, tc, sc.NewVaryingData())

		assert.Equal(t, tc, call.FunctionIndex())
	}
}

func Test_Call_Set_BaseWeight(t *testing.T) {
	target := setUpCallSet()

	assert.Equal(t, callSetWeight(c.DbWeight), target.BaseWeight())
}

func Test_Call_Set_WeighData(t *testing.T) {
	target := setUpCallSet()
	baseWeight := target.BaseWeight()
	assert.Equal(t, primitives.WeightFromParts(baseWeight.RefTime, 0), target.WeighData(baseWeight))
}

func Test_Call_Set_ClassifyDispatch(t *testing.T) {
	target := setUpCallSet()

	assert.Equal(t, primitives.NewDispatchClassMandatory(), target.ClassifyDispatch(baseWeight))
}

func Test_Call_Set_PaysFee(t *testing.T) {
	target := setUpCallSet()

	assert.Equal(t, primitives.PaysYes, target.PaysFee(baseWeight))
}

func Test_Call_Set_Docs(t *testing.T) {
	target := setUpCallSet()

	assert.Equal(t, "Set the current time.", target.Docs())
}

func Test_Call_Set_Dispatch_Success(t *testing.T) {
	target := setUpCallSet()
	mockStorageDidUpdate.On("Exists").Return(false)
	mockStorageNow.On("Get").Return(sc.U64(0), nil)
	mockStorageNow.On("Put", now).Return()
	mockStorageDidUpdate.On("Put", sc.Bool(true)).Return()
	mockOnTimestampSet.On("OnTimestampSet", now).Return(nil)

	_, dispatchErr := target.Dispatch(origin, sc.NewVaryingData(sc.ToCompact(now)))

	assert.Nil(t, dispatchErr)
	mockStorageDidUpdate.AssertCalled(t, "Exists")
	mockStorageNow.AssertCalled(t, "Get")
	mockStorageNow.AssertCalled(t, "Put", now)
	mockStorageDidUpdate.AssertCalled(t, "Put", sc.Bool(true))
	mockOnTimestampSet.AssertCalled(t, "OnTimestampSet", now)
}

func Test_Call_Set_Dispatch_Success_ValidTimestamp(t *testing.T) {
	target := setUpCallSet()
	mockStorageDidUpdate.On("Exists").Return(false)
	mockStorageNow.On("Get").Return(now-c.MinimumPeriod, nil)
	mockStorageNow.On("Put", now).Return()
	mockStorageDidUpdate.On("Put", sc.Bool(true)).Return()
	mockOnTimestampSet.On("OnTimestampSet", now).Return(nil)

	result, dispatchErr := target.Dispatch(origin, sc.NewVaryingData(sc.ToCompact(now)))

	assert.Equal(t, primitives.PostDispatchInfo{}, result)
	assert.Nil(t, dispatchErr)
	mockStorageDidUpdate.AssertCalled(t, "Exists")
	mockStorageNow.AssertCalled(t, "Get")
	mockStorageNow.AssertCalled(t, "Put", now)
	mockStorageDidUpdate.AssertCalled(t, "Put", sc.Bool(true))
	mockOnTimestampSet.AssertCalled(t, "OnTimestampSet", now)
}

func Test_Call_Set_Dispatch_InvalidOrigin(t *testing.T) {
	target := setUpCallSet()

	_, dispatchErr := target.Dispatch(primitives.NewRawOriginRoot(), sc.NewVaryingData(sc.ToCompact(now)))

	assert.Equal(t, primitives.NewDispatchErrorBadOrigin(), dispatchErr)
	mockStorageDidUpdate.AssertNotCalled(t, "Exists")
	mockStorageNow.AssertNotCalled(t, "Get")
	mockStorageNow.AssertNotCalled(t, "Put")
	mockStorageDidUpdate.AssertNotCalled(t, "Put")
	mockOnTimestampSet.AssertNotCalled(t, "OnTimestampSet")
}

func Test_Call_Set_Dispatch_InvalidArgs(t *testing.T) {
	target := setUpCallSet()

	_, dispatchErr := target.Dispatch(primitives.NewRawOriginRoot(), sc.NewVaryingData(sc.NewU64(0)))

	assert.Equal(t, errors.New("couldn't dispatch call set timestamp compact value"), dispatchErr)
	mockStorageDidUpdate.AssertNotCalled(t, "Exists")
	mockStorageNow.AssertNotCalled(t, "Get")
	mockStorageNow.AssertNotCalled(t, "Put")
	mockStorageDidUpdate.AssertNotCalled(t, "Put")
	mockOnTimestampSet.AssertNotCalled(t, "OnTimestampSet")
}

func Test_Call_Set_Dispatch_InvalidStorageDidUpdate(t *testing.T) {
	target := setUpCallSet()
	mockStorageDidUpdate.On("Exists").Return(true)

	_, err := target.Dispatch(origin, sc.NewVaryingData(sc.ToCompact(now)))
	assert.Equal(t, primitives.NewDispatchErrorOther(sc.Str(errTimestampUpdatedOnce.Error())), err)

	mockStorageNow.AssertNotCalled(t, "Get")
	mockStorageNow.AssertNotCalled(t, "Put")
	mockStorageDidUpdate.AssertNotCalled(t, "Put")
	mockOnTimestampSet.AssertNotCalled(t, "OnTimestampSet")
}

func Test_Call_Set_Dispatch_InvalidPreviousTimestamp(t *testing.T) {
	target := setUpCallSet()
	mockStorageDidUpdate.On("Exists").Return(false)
	mockStorageNow.On("Get").Return(sc.U64(1000), nil)

	_, err := target.Dispatch(origin, sc.NewVaryingData(sc.ToCompact(now)))
	assert.Equal(t, primitives.NewDispatchErrorOther(sc.Str(errTimestampMinimumPeriod.Error())), err)

	mockStorageNow.AssertNotCalled(t, "Put")
	mockStorageDidUpdate.AssertNotCalled(t, "Put")
	mockOnTimestampSet.AssertNotCalled(t, "OnTimestampSet")
}

func Test_Call_Set_Dispatch_InvalidLessThanMinPeriod(t *testing.T) {
	target := setUpCallSet()
	mockStorageDidUpdate.On("Exists").Return(false)
	mockStorageNow.On("Get").Return(sc.U64(1001), nil)

	_, err := target.Dispatch(origin, sc.NewVaryingData(sc.ToCompact(now)))
	assert.Equal(t, primitives.NewDispatchErrorOther(sc.Str(errTimestampMinimumPeriod.Error())), err)

	mockStorageNow.AssertNotCalled(t, "Put")
	mockStorageDidUpdate.AssertNotCalled(t, "Put")
	mockOnTimestampSet.AssertNotCalled(t, "OnTimestampSet")
}

func Test_Call_Set_Dispatch_NowGet_Error(t *testing.T) {
	target := setUpCallSet()

	mockErr := errors.New("err")
	expectedErr := primitives.NewDispatchErrorOther(sc.Str(mockErr.Error()))

	mockStorageDidUpdate.On("Exists").Return(false)
	mockStorageNow.On("Get").Return(sc.U64(1001), mockErr)

	_, err := target.Dispatch(origin, sc.NewVaryingData(sc.ToCompact(now)))

	assert.Equal(t, expectedErr, err)
	mockStorageNow.AssertNotCalled(t, "Put")
	mockStorageDidUpdate.AssertNotCalled(t, "Put")
	mockOnTimestampSet.AssertNotCalled(t, "OnTimestampSet")
}

func setUpCallSet() callSet {
	mockOnTimestampSet = new(mocks.OnTimestampSet)
	mockStorageNow = new(mocks.StorageValue[sc.U64])
	mockStorageDidUpdate = new(mocks.StorageValue[sc.Bool])
	mockStorage = &storage{
		Now:       mockStorageNow,
		DidUpdate: mockStorageDidUpdate,
	}

	return newCallSet(0, functionSetIndex, mockStorage, c, mockOnTimestampSet).(callSet)
}
