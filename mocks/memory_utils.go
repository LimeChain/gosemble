package mocks

import "github.com/stretchr/testify/mock"

type MemoryTranslator struct {
	mock.Mock
}

func (m *MemoryTranslator) Int64ToOffsetAndSize(offsetAndSize int64) (offset int32, size int32) {
	args := m.Called(offsetAndSize)
	return args.Get(0).(int32), args.Get(1).(int32)
}

func (m *MemoryTranslator) Offset32(data []byte) int32 {
	args := m.Called(data)
	return args.Get(0).(int32)
}

func (m *MemoryTranslator) BytesToOffsetAndSize(data []byte) int64 {
	args := m.Called(data)
	return args.Get(0).(int64)
}

func (m *MemoryTranslator) GetWasmMemorySlice(offset int32, size int32) []byte {
	args := m.Called(offset, size)
	return args.Get(0).([]byte)
}
