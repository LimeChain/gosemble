package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/stretchr/testify/mock"
)

type TransactionPaymentStorage struct {
	mock.Mock
}

func (s *TransactionPaymentStorage) GetNextFeeMultiplier() sc.U128 {
	args := s.Called()
	return args.Get(0).(sc.U128)
}
