package mocks

import (
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type EventDepositor struct {
	mock.Mock
}

func (m *EventDepositor) DepositEvent(event types.Event) {
	_ = m.Called(event)
}
