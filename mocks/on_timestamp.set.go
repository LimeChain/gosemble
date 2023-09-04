package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/mock"
)

type MockOnTimestampSet struct {
	mock.Mock
}

func (m *MockOnTimestampSet) OnTimestampSet(n sc.U64) {
	m.Called(n)
}
