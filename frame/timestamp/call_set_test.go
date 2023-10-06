package timestamp

import (
	"bytes"
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
			Arguments:  nil,
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
	compact := sc.ToCompact(sc.U8(5))
	buf := bytes.NewBuffer(compact.Bytes())

	call := target.DecodeArgs(buf)

	assert.Equal(t, sc.NewVaryingData(compact), call.Args())
}

func Test_Call_Set_Encode(t *testing.T) {
	target := setUpCallSet()
	expectedBuffer := bytes.NewBuffer([]byte{moduleId, functionSetIndex})
	buf := &bytes.Buffer{}

	target.Encode(buf)

	assert.Equal(t, expectedBuffer, buf)
}

func Test_Call_Set_EncodeWithArgs(t *testing.T) {
	target := setUpCallSet()
	compact := sc.ToCompact(sc.U8(5))

	expectedBuf := bytes.NewBuffer(append([]byte{moduleId, functionSetIndex}, compact.Bytes()...))

	buf := bytes.NewBuffer(compact.Bytes())

	call := target.DecodeArgs(buf)

	buf.Reset()
	call.Encode(buf)

	assert.Equal(t, expectedBuf, buf)
}

func Test_Call_Set_Bytes(t *testing.T) {
	target := setUpCallSet()
	expected := []byte{moduleId, functionSetIndex}

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

	assert.Equal(t, primitives.WeightFromParts(159_258_000, 1006), target.BaseWeight())
}

func Test_Call_Set_WeighData(t *testing.T) {
	target := setUpCallSet()
	assert.Equal(t, primitives.WeightFromParts(124, 0), target.WeighData(baseWeight))
}

func Test_Call_Set_ClassifyDispatch(t *testing.T) {
	target := setUpCallSet()

	assert.Equal(t, primitives.NewDispatchClassMandatory(), target.ClassifyDispatch(baseWeight))
}

func Test_Call_Set_PaysFee(t *testing.T) {
	target := setUpCallSet()

	assert.Equal(t, primitives.NewPaysYes(), target.PaysFee(baseWeight))
}

func Test_Call_Set_Dispatch_Success(t *testing.T) {
	target := setUpCallSet()
	mockStorageDidUpdate.On("Exists").Return(false)
	mockStorageNow.On("Get").Return(sc.U64(0))
	mockStorageNow.On("Put", now).Return()
	mockStorageDidUpdate.On("Put", sc.Bool(true)).Return()
	mockOnTimestampSet.On("OnTimestampSet", now).Return()

	result := target.Dispatch(origin, sc.NewVaryingData(sc.ToCompact(now)))

	assert.Equal(t, primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{}, result)
	mockStorageDidUpdate.AssertCalled(t, "Exists")
	mockStorageNow.AssertCalled(t, "Get")
	mockStorageNow.AssertCalled(t, "Put", now)
	mockStorageDidUpdate.AssertCalled(t, "Put", sc.Bool(true))
	mockOnTimestampSet.AssertCalled(t, "OnTimestampSet", now)
}

func Test_Call_Set_set_Success_ZeroPreviousTimestamp(t *testing.T) {
	target := setUpCallSet()
	mockStorageDidUpdate.On("Exists").Return(false)
	mockStorageNow.On("Get").Return(sc.U64(0))
	mockStorageNow.On("Put", now).Return()
	mockStorageDidUpdate.On("Put", sc.Bool(true)).Return()
	mockOnTimestampSet.On("OnTimestampSet", now).Return()

	result := target.set(origin, now)

	assert.Equal(t, primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{}, result)
	mockStorageDidUpdate.AssertCalled(t, "Exists")
	mockStorageNow.AssertCalled(t, "Get")
	mockStorageNow.AssertCalled(t, "Put", now)
	mockStorageDidUpdate.AssertCalled(t, "Put", sc.Bool(true))
	mockOnTimestampSet.AssertCalled(t, "OnTimestampSet", now)
}

func Test_Call_Set_set_Success_ValidTimestamp(t *testing.T) {
	target := setUpCallSet()
	mockStorageDidUpdate.On("Exists").Return(false)
	mockStorageNow.On("Get").Return(now - c.MinimumPeriod)
	mockStorageNow.On("Put", now).Return()
	mockStorageDidUpdate.On("Put", sc.Bool(true)).Return()
	mockOnTimestampSet.On("OnTimestampSet", now).Return()

	result := target.set(origin, now)

	assert.Equal(t, primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{}, result)
	mockStorageDidUpdate.AssertCalled(t, "Exists")
	mockStorageNow.AssertCalled(t, "Get")
	mockStorageNow.AssertCalled(t, "Put", now)
	mockStorageDidUpdate.AssertCalled(t, "Put", sc.Bool(true))
	mockOnTimestampSet.AssertCalled(t, "OnTimestampSet", now)
}

func Test_Call_Set_set_InvalidOrigin(t *testing.T) {
	target := setUpCallSet()
	expected := primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]{
		HasError: true,
		Err: primitives.DispatchErrorWithPostInfo[primitives.PostDispatchInfo]{
			Error: primitives.NewDispatchErrorBadOrigin(),
		},
	}

	result := target.set(primitives.NewRawOriginRoot(), now)

	assert.Equal(t, expected, result)
	mockStorageDidUpdate.AssertNotCalled(t, "Exists")
	mockStorageNow.AssertNotCalled(t, "Get")
	mockStorageNow.AssertNotCalled(t, "Put")
	mockStorageDidUpdate.AssertNotCalled(t, "Put")
	mockOnTimestampSet.AssertNotCalled(t, "OnTimestampSet")
}

func Test_Call_Set_set_InvalidStorageDidUpdate(t *testing.T) {
	target := setUpCallSet()
	mockStorageDidUpdate.On("Exists").Return(true)

	assert.PanicsWithValue(
		t,
		errTimestampUpdatedOnce,
		func() {
			target.set(origin, now)
		})
	mockStorageNow.AssertNotCalled(t, "Get")
	mockStorageNow.AssertNotCalled(t, "Put")
	mockStorageDidUpdate.AssertNotCalled(t, "Put")
	mockOnTimestampSet.AssertNotCalled(t, "OnTimestampSet")
}

func Test_Call_Set_set_InvalidPreviousTimestamp(t *testing.T) {
	target := setUpCallSet()
	mockStorageDidUpdate.On("Exists").Return(false)
	mockStorageNow.On("Get").Return(sc.U64(1000))

	assert.PanicsWithValue(t,
		errTimestampMinimumPeriod,
		func() {
			target.set(origin, now)
		})

	mockStorageNow.AssertNotCalled(t, "Put")
	mockStorageDidUpdate.AssertNotCalled(t, "Put")
	mockOnTimestampSet.AssertNotCalled(t, "OnTimestampSet")
}

func Test_Call_Set_set_InvalidLessThanMinPeriod(t *testing.T) {
	target := setUpCallSet()
	mockStorageDidUpdate.On("Exists").Return(false)
	mockStorageNow.On("Get").Return(sc.U64(1001))

	assert.PanicsWithValue(t,
		errTimestampMinimumPeriod,
		func() {
			target.set(origin, now)
		})
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
