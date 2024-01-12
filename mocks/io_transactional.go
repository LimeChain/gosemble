package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/mock"
)

type IoTransactional[T sc.Encodable] struct {
	mock.Mock
}

func (m *IoTransactional[T]) WithStorageLayer(fn func() (T, error)) (T, error) {
	args := m.Called(fn)

	if args.Get(1) == nil {
		return args.Get(0).(T), nil
	}

	return args.Get(0).(T), args.Get(1).(error)
}
