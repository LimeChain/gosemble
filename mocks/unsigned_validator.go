package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type UnsignedValidator struct {
	mock.Mock
}

func (m *UnsignedValidator) PreDispatch(call types.Call) (sc.Empty, error) {
	args := m.Called(call)

	if args.Get(1) == nil {
		return args.Get(0).(sc.Empty), nil
	}

	return args.Get(0).(sc.Empty), args.Get(1).(error)
}

func (m *UnsignedValidator) ValidateUnsigned(source types.TransactionSource, call types.Call) (types.ValidTransaction, error) {
	args := m.Called(source, call)

	if args.Get(1) == nil {
		return args.Get(0).(types.ValidTransaction), nil
	}

	return args.Get(0).(types.ValidTransaction), args.Get(1).(error)
}
