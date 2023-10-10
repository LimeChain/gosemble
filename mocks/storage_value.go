package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/mock"
)

type StorageValue[T sc.Encodable] struct {
	mock.Mock
}

func (m *StorageValue[T]) Get() T {
	args := m.Called()

	return args.Get(0).(T)
}

func (m *StorageValue[T]) GetBytes() sc.Option[sc.Sequence[sc.U8]] {
	args := m.Called()
	return args.Get(0).(sc.Option[sc.Sequence[sc.U8]])
}

func (m *StorageValue[T]) Exists() bool {
	args := m.Called()

	return args.Bool(0)
}

func (m *StorageValue[T]) Put(value T) {
	m.Called(value)
}

func (m *StorageValue[T]) Clear() {
	m.Called()
}

func (m *StorageValue[T]) Append(value T) {
	m.Called(value)
}

func (m *StorageValue[T]) Take() T {
	args := m.Called()

	return args.Get(0).(T)
}

func (m *StorageValue[T]) TakeBytes() []byte {
	args := m.Called()

	return args.Get(0).([]byte)
}

func (m *StorageValue[T]) DecodeLen() sc.Option[sc.U64] {
	args := m.Called()

	return args.Get(0).(sc.Option[sc.U64])
}
