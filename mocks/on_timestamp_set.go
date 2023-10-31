package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/mock"
)

type OnTimestampSet struct {
	mock.Mock
}

func (m *OnTimestampSet) OnTimestampSet(n sc.U64) error {
	args := m.Called(n)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}
