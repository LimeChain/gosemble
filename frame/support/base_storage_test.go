package support

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	defaultValue = sc.U32(1)
	storageValue = sc.U32(5)
	decodeFunc   = sc.DecodeU32
	key          = []byte("key")

	mockStorage *mocks.IoStorage
)

func Test_baseStorage_get(t *testing.T) {
	target := setupBaseStorage()

	mockStorage.On("Get", key).Return(sc.NewOption[sc.Sequence[sc.U8]](nil))

	result := target.get(key)

	assert.Equal(t, defaultValue, result)
	mockStorage.AssertCalled(t, "Get", key)
}

func Test_baseStorage_get_Nil(t *testing.T) {
	target := setupBaseStorage()
	target.defaultValue = nil

	mockStorage.On("Get", key).Return(
		sc.NewOption[sc.Sequence[sc.U8]](
			sc.BytesToSequenceU8(storageValue.Bytes()),
		),
	)

	result := target.get(key)

	assert.Equal(t, storageValue, result)
	mockStorage.AssertCalled(t, "Get", key)
}

func Test_baseStorage_getBytes(t *testing.T) {
	target := setupBaseStorage()
	expect := sc.NewOption[sc.Sequence[sc.U8]](nil)

	mockStorage.On("Get", key).Return(expect)

	result := target.getBytes(key)

	assert.Equal(t, expect, result)
	mockStorage.AssertCalled(t, "Get", key)
}

func Test_baseStorage_exists(t *testing.T) {
	target := setupBaseStorage()

	mockStorage.On("Exists", key).Return(true)

	result := target.exists(key)

	assert.True(t, result)
	mockStorage.AssertCalled(t, "Exists", key)
}

func Test_baseStorage_put(t *testing.T) {
	target := setupBaseStorage()

	mockStorage.On("Set", key, storageValue.Bytes())

	target.put(key, storageValue)

	mockStorage.AssertCalled(t, "Set", key, storageValue.Bytes())
}

func Test_baseStorage_clear(t *testing.T) {
	target := setupBaseStorage()

	mockStorage.On("Clear", key)

	target.clear(key)

	mockStorage.AssertCalled(t, "Clear", key)
}

func Test_baseStorage_append(t *testing.T) {
	target := setupBaseStorage()

	mockStorage.On("Append", key, storageValue.Bytes())

	target.append(key, storageValue)

	mockStorage.AssertCalled(t, "Append", key, storageValue.Bytes())
}

func Test_baseStorage_take(t *testing.T) {
	target := setupBaseStorage()

	mockStorage.On("Get", key).Return(
		sc.NewOption[sc.Sequence[sc.U8]](
			sc.BytesToSequenceU8(storageValue.Bytes()),
		),
	)
	mockStorage.On("Clear", key).Return()

	result := target.take(key)

	assert.Equal(t, storageValue, result)
	mockStorage.AssertCalled(t, "Get", key)
	mockStorage.AssertCalled(t, "Clear", key)
}

func Test_baseStorage_decodeLen(t *testing.T) {
	target := setupBaseStorage()
	compactBytes := [5]byte{}
	offset := int32(0)

	mockStorage.On("Read", key, compactBytes[:], offset).Return(sc.NewOption[sc.U32](sc.U32(4)))

	result := target.decodeLen(key)

	assert.Equal(t, sc.NewOption[sc.U64](sc.U64(0)), result)
	mockStorage.AssertCalled(t, "Read", key, compactBytes[:], offset)
}

func Test_baseStorage_decodeLen_Empty(t *testing.T) {
	target := setupBaseStorage()
	compactBytes := [5]byte{}
	offset := int32(0)

	mockStorage.On("Read", key, compactBytes[:], offset).Return(sc.NewOption[sc.U32](nil))

	result := target.decodeLen(key)

	assert.Equal(t, sc.NewOption[sc.U64](nil), result)
	mockStorage.AssertCalled(t, "Read", key, compactBytes[:], offset)
}

func Test_baseStorage_getDecode(t *testing.T) {
	target := setupBaseStorage()

	mockStorage.On("Get", key).Return(
		sc.NewOption[sc.Sequence[sc.U8]](
			sc.BytesToSequenceU8(storageValue.Bytes()),
		),
	)

	result := target.getDecode(key)

	assert.Equal(t, storageValue, result)
	mockStorage.AssertCalled(t, "Get", key)
}

func Test_baseStorage_getDecode_Default(t *testing.T) {
	target := setupBaseStorage()

	mockStorage.On("Get", key).Return(sc.NewOption[sc.Sequence[sc.U8]](nil))

	result := target.getDecode(key)

	assert.Equal(t, sc.U32(0), result)
	mockStorage.AssertCalled(t, "Get", key)
}

func Test_baseStorage_getDecodeOnEmpty_Default(t *testing.T) {
	target := setupBaseStorage()

	mockStorage.On("Get", key).Return(sc.NewOption[sc.Sequence[sc.U8]](nil))

	result := target.getDecodeOnEmpty(key)

	assert.Equal(t, defaultValue, result)
	mockStorage.AssertCalled(t, "Get", key)
}

func Test_baseStorage_getDecodeOnEmpty(t *testing.T) {
	target := setupBaseStorage()

	mockStorage.On("Get", key).Return(
		sc.NewOption[sc.Sequence[sc.U8]](
			sc.BytesToSequenceU8(storageValue.Bytes()),
		),
	)

	result := target.getDecodeOnEmpty(key)

	assert.Equal(t, storageValue, result)
	mockStorage.AssertCalled(t, "Get", key)
}

func Test_baseStorage_takeBytes(t *testing.T) {
	target := setupBaseStorage()

	mockStorage.On("Get", key).Return(
		sc.NewOption[sc.Sequence[sc.U8]](
			sc.BytesToSequenceU8(storageValue.Bytes()),
		),
	)
	mockStorage.On("Clear", key).Return()

	result := target.takeBytes(key)

	assert.Equal(t, storageValue.Bytes(), result)
	mockStorage.AssertCalled(t, "Get", key)
	mockStorage.AssertCalled(t, "Clear", key)
}

func Test_baseStorage_takeDecode_Nil(t *testing.T) {
	target := setupBaseStorage()

	mockStorage.On("Get", key).Return(sc.NewOption[sc.Sequence[sc.U8]](nil))

	result := target.takeBytes(key)

	assert.Equal(t, []byte(nil), result)
	mockStorage.AssertCalled(t, "Get", key)
	mockStorage.AssertNotCalled(t, "Clear", mock.Anything)
}

func Test_baseStorage_takeDecode(t *testing.T) {
	target := setupBaseStorage()

	mockStorage.On("Get", key).Return(
		sc.NewOption[sc.Sequence[sc.U8]](
			sc.BytesToSequenceU8(storageValue.Bytes()),
		),
	)
	mockStorage.On("Clear", key).Return()

	result := target.takeDecode(key)

	assert.Equal(t, storageValue, result)
	mockStorage.AssertCalled(t, "Get", key)
	mockStorage.AssertCalled(t, "Clear", key)
}

func Test_baseStorage_takeDecode_Default(t *testing.T) {
	target := setupBaseStorage()

	mockStorage.On("Get", key).Return(sc.NewOption[sc.Sequence[sc.U8]](nil))

	result := target.takeDecode(key)

	assert.Equal(t, sc.U32(0), result)
	mockStorage.AssertCalled(t, "Get", key)
	mockStorage.AssertNotCalled(t, "Clear", mock.Anything)
}

func setupBaseStorage() baseStorage[sc.U32] {
	mockStorage = new(mocks.IoStorage)
	target := newBaseStorage(decodeFunc, &defaultValue)

	target.storage = mockStorage

	return target
}
