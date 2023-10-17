package mocks

import (
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type AccountIdLookup struct {
	mock.Mock
}

func (l *AccountIdLookup) Lookup(a types.MultiAddress) (types.Address32, types.TransactionValidityError) {
	args := l.Called(a)

	if args.Get(1) == nil {
		return args.Get(0).(types.Address32), nil
	}

	return args.Get(0).(types.Address32), args.Get(1).(types.TransactionValidityError)
}
