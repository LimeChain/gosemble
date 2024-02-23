package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/mock"
)

type RawStorageValue struct {
	mock.Mock
}

func (m *RawStorageValue) Get() (sc.Sequence[sc.U8], error) {
	args := m.Called()

	if args.Get(1) == nil {
		return args.Get(0).(sc.Sequence[sc.U8]), nil
	}

	return args.Get(0).(sc.Sequence[sc.U8]), args.Get(1).(error)
}

func (m *RawStorageValue) Put(value sc.Sequence[sc.U8]) {
	m.Called(value)
}

func (m *RawStorageValue) Clear() {
	m.Called()
}

func (m *RawStorageValue) ClearPrefix(limit sc.U32) {
	m.Called(limit)
}
