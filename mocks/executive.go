package mocks

import (
	"github.com/LimeChain/gosemble/primitives/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type Executive struct {
	mock.Mock
}

func (m *Executive) InitializeBlock(header primitives.Header) error {
	args := m.Called(header)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func (m *Executive) ExecuteBlock(block types.Block) error {
	args := m.Called(block)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}

func (m *Executive) ApplyExtrinsic(uxt types.UncheckedExtrinsic) error {
	args := m.Called(uxt)

	if args.Get(0) != nil {
		return args.Get(0).(error)
	}

	return nil
}

func (m *Executive) FinalizeBlock() (primitives.Header, error) {
	args := m.Called()

	if args.Get(1) == nil {
		return args.Get(0).(primitives.Header), nil
	}

	return args.Get(0).(primitives.Header), args.Get(1).(error)
}

func (m *Executive) ValidateTransaction(source primitives.TransactionSource, uxt types.UncheckedExtrinsic, blockHash primitives.Blake2bHash) (primitives.ValidTransaction, error) {
	args := m.Called(source, uxt, blockHash)

	if args.Get(1) != nil {
		return args.Get(0).(primitives.ValidTransaction), args.Get(1).(error)
	}

	return args.Get(0).(primitives.ValidTransaction), nil
}

func (m *Executive) OffchainWorker(header primitives.Header) error {
	args := m.Called(header)
	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(error)
}
