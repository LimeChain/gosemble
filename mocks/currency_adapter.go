package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type CurrencyAdapter struct {
	mock.Mock
}

func (m *CurrencyAdapter) DepositIntoExisting(who types.AccountId[types.PublicKey], value sc.U128) (types.Balance, types.DispatchError) {
	args := m.Called(who, value)

	if args.Get(1) != nil {
		return args.Get(0).(types.Balance), args.Get(1).(types.DispatchError)
	}

	return args.Get(0).(types.Balance), types.DispatchError{VaryingData: nil}
}

func (m *CurrencyAdapter) Withdraw(who types.AccountId[types.PublicKey], value sc.U128, reasons sc.U8, liveness types.ExistenceRequirement) (types.Balance, types.DispatchError) {
	args := m.Called(who, value, reasons, liveness)

	if args.Get(1) != nil {
		return args.Get(0).(types.Balance), args.Get(1).(types.DispatchError)
	}

	return args.Get(0).(types.Balance), types.DispatchError{VaryingData: nil}
}
