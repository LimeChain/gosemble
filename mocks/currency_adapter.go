package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type MockCurrencyAdapter struct {
	mock.Mock
}

func (m *MockCurrencyAdapter) DepositIntoExisting(who types.Address32, value sc.U128) (types.Balance, types.DispatchError) {
	args := m.Called(who, value)

	if args[1] != nil {
		return args[0].(types.Balance), args[1].(types.DispatchError)
	}

	return args[0].(types.Balance), nil
}

func (m *MockCurrencyAdapter) Withdraw(who types.Address32, value sc.U128, reasons sc.U8, liveness types.ExistenceRequirement) (types.Balance, types.DispatchError) {
	args := m.Called(who, value, reasons, liveness)

	if args[1] != nil {
		return args[0].(types.Balance), args[1].(types.DispatchError)
	}

	return args[0].(types.Balance), nil
}
