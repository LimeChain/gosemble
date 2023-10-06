package mocks

import (
	"github.com/stretchr/testify/mock"
)

type IoHashing struct {
	mock.Mock
}

func (m *IoHashing) Blake128(value []byte) []byte {
	args := m.Called(value)

	return args.Get(0).([]byte)
}

func (m *IoHashing) Blake256(value []byte) []byte {
	args := m.Called(value)

	return args.Get(0).([]byte)
}

func (m *IoHashing) Twox128(value []byte) []byte {
	args := m.Called(value)

	return args.Get(0).([]byte)
}

func (m *IoHashing) Twox64(value []byte) []byte {
	args := m.Called(value)

	return args.Get(0).([]byte)
}
