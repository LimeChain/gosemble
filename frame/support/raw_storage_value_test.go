package support

import (
	"errors"
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/stretchr/testify/assert"
)

var (
	rawKey   = []byte("rawkey1")
	rawValue = []byte("rawvalue1")
)

func Test_RawStorageValue_Get(t *testing.T) {
	target := setupRawStorageValue()

	mockStorage.On("Get", rawKey).Return(sc.NewOption[sc.Sequence[sc.U8]](sc.BytesToSequenceU8(rawValue)), nil)

	result, err := target.Get()
	assert.NoError(t, err)

	assert.Equal(t, sc.BytesToSequenceU8(rawValue), result)
	mockStorage.AssertCalled(t, "Get", rawKey)
}

func Test_RawStorageValue_Get_Nil(t *testing.T) {
	target := setupRawStorageValue()

	mockStorage.On("Get", rawKey).Return(sc.NewOption[sc.Sequence[sc.U8]](nil), nil)

	result, err := target.Get()
	assert.NoError(t, err)

	assert.Equal(t, sc.Sequence[sc.U8]{}, result)
	mockStorage.AssertCalled(t, "Get", rawKey)
}

func Test_RawStorageValue_Get_Error(t *testing.T) {
	target := setupRawStorageValue()

	expectedError := errors.New("decode option error")

	mockStorage.On("Get", rawKey).Return(sc.NewOption[sc.Sequence[sc.U8]](nil), expectedError)

	_, err := target.Get()

	assert.Equal(t, expectedError, err)
	mockStorage.AssertCalled(t, "Get", rawKey)
}

func Test_RawStorageValue_Put(t *testing.T) {
	target := setupRawStorageValue()

	mockStorage.On("Set", rawKey, rawValue).Return()

	target.Put(sc.BytesToSequenceU8(rawValue))

	mockStorage.AssertCalled(t, "Set", rawKey, rawValue)
}

func Test_RawStorageValue_Clear(t *testing.T) {
	target := setupRawStorageValue()

	mockStorage.On("Clear", rawKey).Return()

	target.Clear()

	mockStorage.AssertCalled(t, "Clear", rawKey)
}

func Test_RawStorageValue_ClearPrefix(t *testing.T) {
	target := setupRawStorageValue()

	limit := sc.U32(1)
	mockStorage.On("ClearPrefix", rawKey, sc.NewOption[sc.U32](limit).Bytes()).Return()

	target.ClearPrefix(limit)

	mockStorage.AssertCalled(t, "ClearPrefix", rawKey, sc.NewOption[sc.U32](limit).Bytes())
}

func setupRawStorageValue() RawStorageValue {
	mockStorage = new(mocks.IoStorage)
	return NewRawStorageValueFrom(mockStorage, rawKey).(RawStorageValue)
}
