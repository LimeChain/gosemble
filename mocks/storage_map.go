package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/mock"
)

type StorageMap[K, V sc.Encodable] struct {
	mock.Mock
}

func (m *StorageMap[K, V]) Get(k K) (V, error) {
	args := m.Called(k)

	return args.Get(0).(V), nil
}

func (m *StorageMap[K, V]) Exists(k K) bool {
	args := m.Called(k)

	return args.Get(0).(bool)
}

func (m *StorageMap[K, V]) Put(k K, value V) {
	m.Called(k, value)
}

func (m *StorageMap[K, V]) Append(k K, v V) {
	m.Called(k, v)
}

func (m *StorageMap[K, V]) TakeBytes(k K) ([]byte, error) {
	args := m.Called(k)

	return args.Get(0).([]byte), nil
}

func (m *StorageMap[K, V]) Remove(k K) {
	m.Called(k)
}

func (m *StorageMap[K, V]) Clear(limit sc.U32) {
	m.Called(limit)
}

func (m *StorageMap[K, V]) Mutate(k K, f func(value *V) sc.Result[sc.Encodable]) (sc.Result[sc.Encodable], error) {
	args := m.Called(k, f)

	return args.Get(0).(sc.Result[sc.Encodable]), nil
}

func (m *StorageMap[K, V]) TryMutateExists(k K, f func(option *sc.Option[V]) sc.Result[sc.Encodable]) (sc.Result[sc.Encodable], error) {
	args := m.Called(k, f)

	return args.Get(0).(sc.Result[sc.Encodable]), nil
}
