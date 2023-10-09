package mocks

import (
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type DefaultOnRuntimeUpgrade struct {
	mock.Mock
}

func (doru *DefaultOnRuntimeUpgrade) OnRuntimeUpgrade() types.Weight {
	args := doru.Called()
	return args.Get(0).(types.Weight)
}
