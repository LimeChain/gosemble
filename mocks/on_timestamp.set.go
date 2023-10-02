package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/mock"
)

type OnTimestampSet struct {
	mock.Mock
}

func (m *OnTimestampSet) OnTimestampSet(n sc.U64) {
	m.Called(n)
}
