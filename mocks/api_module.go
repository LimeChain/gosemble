package mocks

import (
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type ApiModule struct {
	mock.Mock
}

func (m *ApiModule) Name() string {
	args := m.Called()

	return args.Get(0).(string)
}

func (m *ApiModule) Item() types.ApiItem {
	args := m.Called()

	return args.Get(0).(types.ApiItem)
}
