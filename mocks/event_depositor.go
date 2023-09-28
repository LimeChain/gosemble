package mocks

import (
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type MockEventDepositor struct {
	mock.Mock
}

func (m *MockEventDepositor) DepositEvent(event types.Event) {
	_ = m.Called(event)
}
