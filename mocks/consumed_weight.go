package mocks

import (
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type ConsumedWeight struct {
	mock.Mock
}

func (cw *ConsumedWeight) Total() types.Weight {
	args := cw.Called()
	return args.Get(0).(types.Weight)
}
