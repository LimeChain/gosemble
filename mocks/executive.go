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
	m.Called(header)
	return nil
}

func (m *Executive) ExecuteBlock(block types.Block) error {
	m.Called(block)
	return nil
}

func (m *Executive) ApplyExtrinsic(uxt types.UncheckedExtrinsic) (primitives.DispatchOutcome, primitives.TransactionValidityError) {
	args := m.Called(uxt)

	if args.Get(1) != nil {
		return args.Get(0).(primitives.DispatchOutcome), args.Get(1).(primitives.TransactionValidityError)
	}

	return args.Get(0).(primitives.DispatchOutcome), nil
}

func (m *Executive) FinalizeBlock() (primitives.Header, error) {
	args := m.Called()

	return args.Get(0).(primitives.Header), nil
}

func (m *Executive) ValidateTransaction(source primitives.TransactionSource, uxt types.UncheckedExtrinsic, blockHash primitives.Blake2bHash) (primitives.ValidTransaction, primitives.TransactionValidityError) {
	args := m.Called(source, uxt, blockHash)

	if args.Get(1) != nil {
		return args.Get(0).(primitives.ValidTransaction), args.Get(1).(primitives.TransactionValidityError)
	}

	return args.Get(0).(primitives.ValidTransaction), nil
}

func (m *Executive) OffchainWorker(header primitives.Header) {
	m.Called(header)
}
