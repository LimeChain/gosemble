package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type IoTransactional[T sc.Encodable, E types.DispatchError] struct {
	mock.Mock
}

func (m *IoTransactional[T, E]) WithStorageLayer(fn func() (T, types.DispatchError)) (T, E) {
	args := m.Called(fn)

	if args.Get(1) == nil {
		return args.Get(0).(T), E(types.DispatchError{VaryingData: nil})
	}

	return args.Get(0).(T), args.Get(1).(E)
}
