package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/mock"
)

type MockStorageValue[T sc.Encodable] struct {
	mock.Mock
}

func (m *MockStorageValue[T]) Get() T {
	args := m.Called()

	return args[0].(T)
}

func (m *MockStorageValue[T]) GetBytes() sc.Option[sc.Sequence[sc.U8]] {
	args := m.Called()
	return args[0].(sc.Option[sc.Sequence[sc.U8]])
}

func (m *MockStorageValue[T]) Exists() bool {
	args := m.Called()

	return args.Bool(0)
}

func (m *MockStorageValue[T]) Put(value T) {
	m.Called(value)
}

func (m *MockStorageValue[T]) Clear() {
	m.Called()
}

func (m *MockStorageValue[T]) Append(value T) {
	m.Called(value)
}

func (m *MockStorageValue[T]) Take() T {
	args := m.Called()

	return args[0].(T)
}

func (m *MockStorageValue[T]) TakeBytes() []byte {
	args := m.Called()

	return args[0].([]byte)
}

func (m *MockStorageValue[T]) DecodeLen() sc.Option[sc.U64] {
	args := m.Called()

	return args[0].(sc.Option[sc.U64])
}
