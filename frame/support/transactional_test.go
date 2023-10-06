package support

import (
	"testing"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/mocks"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	transactionLevel      = sc.U32(10)
	mockStorageValue      *mocks.StorageValue[sc.U32]
	mockTransactionBroker *mocks.IoTransactionBroker
)

func Test_Transactional_GetTransactionLevel(t *testing.T) {
	target := setupTransactional()

	mockStorageValue.On("Get").Return(transactionLevel)

	assert.Equal(t, transactionLevel, target.GetTransactionLevel())
	mockStorageValue.AssertCalled(t, "Get")
}

func Test_Transactional_SetTransactionLevel(t *testing.T) {
	target := setupTransactional()

	mockStorageValue.On("Put", transactionLevel).Return()

	target.SetTransactionLevel(transactionLevel)

	mockStorageValue.AssertCalled(t, "Put", transactionLevel)
}

func Test_Transactional_KillTransactionLevel(t *testing.T) {
	target := setupTransactional()

	mockStorageValue.On("Clear").Return()

	target.KillTransactionLevel()

	mockStorageValue.AssertCalled(t, "Clear")
}

func Test_Transactional_IncTransactionLevel(t *testing.T) {
	target := setupTransactional()

	mockStorageValue.On("Get").Return(transactionLevel)
	mockStorageValue.On("Put", transactionLevel+1).Return()

	err := target.IncTransactionLevel()

	assert.Nil(t, err)
	mockStorageValue.AssertCalled(t, "Get")
	mockStorageValue.AssertCalled(t, "Put", transactionLevel+1)
}

func Test_Transactional_IncTransactionLevel_MaxLimit(t *testing.T) {
	target := setupTransactional()

	mockStorageValue.On("Get").Return(TransactionalLimit)

	result := target.IncTransactionLevel()

	assert.Equal(t, errTransactionalLimitReached, result)
	mockStorageValue.AssertCalled(t, "Get")
	mockStorageValue.AssertNotCalled(t, "Put", mock.Anything)
}

func Test_Transactional_DecTransactionLevel_Zero(t *testing.T) {
	target := setupTransactional()

	mockStorageValue.On("Get").Return(sc.U32(0))

	target.DecTransactionLevel()

	mockStorageValue.AssertCalled(t, "Get")
	mockStorageValue.AssertNotCalled(t, "Clear")
	mockStorageValue.AssertNotCalled(t, "Put", mock.Anything)
}

func Test_Transactional_DecTransactionLevel_One(t *testing.T) {
	target := setupTransactional()

	mockStorageValue.On("Get").Return(sc.U32(1))
	mockStorageValue.On("Clear").Return()

	target.DecTransactionLevel()

	mockStorageValue.AssertCalled(t, "Get")
	mockStorageValue.AssertCalled(t, "Clear")
	mockStorageValue.AssertNotCalled(t, "Put", mock.Anything)
}

func Test_Transactional_DecTransactionLevel(t *testing.T) {
	target := setupTransactional()

	mockStorageValue.On("Get").Return(transactionLevel)
	mockStorageValue.On("Put", transactionLevel-1).Return()

	target.DecTransactionLevel()

	mockStorageValue.AssertCalled(t, "Get")
	mockStorageValue.AssertCalled(t, "Put", transactionLevel-1)
}

func Test_Transactional_WithTransaction_ErrorLimitReached(t *testing.T) {
	target := setupTransactional()
	expect := primitives.NewDispatchErrorTransactional(primitives.NewTransactionalErrorLimitReached())

	mockStorageValue.On("Get").Return(TransactionalLimit)

	res, err := target.WithTransaction(func() primitives.TransactionOutcome {
		return primitives.NewTransactionOutcomeCommit(sc.U32(2))
	})

	assert.Equal(t, sc.U32(0), res)
	assert.Equal(t, expect, err)

	mockStorageValue.AssertCalled(t, "Get")
}

func Test_Transactional_WithTransaction_Commit(t *testing.T) {
	target := setupTransactional()
	expect := sc.U32(2)

	mockStorageValue.On("Get").Return(transactionLevel).Once()
	mockStorageValue.On("Put", transactionLevel+1).Once()
	mockTransactionBroker.On("Start").Return()
	mockTransactionBroker.On("Commit").Return()
	mockStorageValue.On("Get").Return(transactionLevel + 1).Once()
	mockStorageValue.On("Put", transactionLevel).Once()

	res, err := target.WithTransaction(func() primitives.TransactionOutcome {
		return primitives.NewTransactionOutcomeCommit(expect)
	})

	assert.Equal(t, expect, res)
	assert.Nil(t, err)

	mockStorageValue.AssertNumberOfCalls(t, "Get", 2)
	mockStorageValue.AssertCalled(t, "Put", transactionLevel+1)
	mockTransactionBroker.AssertCalled(t, "Start")
	mockTransactionBroker.AssertCalled(t, "Commit")
	mockStorageValue.AssertCalled(t, "Put", transactionLevel)
}

func Test_Transactional_WithTransaction_Rollback(t *testing.T) {
	target := setupTransactional()
	expect := primitives.NewDispatchErrorCorruption()

	mockStorageValue.On("Get").Return(transactionLevel).Once()
	mockStorageValue.On("Put", transactionLevel+1).Once()
	mockTransactionBroker.On("Start").Return()
	mockTransactionBroker.On("Rollback").Return()
	mockStorageValue.On("Get").Return(transactionLevel + 1).Once()
	mockStorageValue.On("Put", transactionLevel).Once()

	res, err := target.WithTransaction(func() primitives.TransactionOutcome {
		return primitives.NewTransactionOutcomeRollback(expect)
	})

	assert.Equal(t, sc.U32(0), res)
	assert.Equal(t, expect, err)

	mockStorageValue.AssertNumberOfCalls(t, "Get", 2)
	mockStorageValue.AssertNumberOfCalls(t, "Put", 2)
	mockStorageValue.AssertCalled(t, "Put", transactionLevel+1)
	mockTransactionBroker.AssertCalled(t, "Start")
	mockTransactionBroker.AssertCalled(t, "Rollback")
	mockStorageValue.AssertCalled(t, "Put", transactionLevel)
}

func Test_Transactional_WithTransaction_InvalidTransactionOutcome(t *testing.T) {
	target := setupTransactional()

	mockStorageValue.On("Get").Return(transactionLevel).Once()
	mockStorageValue.On("Put", transactionLevel+1).Once()
	mockTransactionBroker.On("Start").Return()

	assert.PanicsWithValue(t, errInvalidTransactionOutcome, func() {
		target.WithTransaction(func() primitives.TransactionOutcome {
			return sc.NewVaryingData(sc.U32(3))
		})
	})

	mockStorageValue.AssertCalled(t, "Get")
	mockStorageValue.AssertCalled(t, "Put", transactionLevel+1)
	mockTransactionBroker.AssertCalled(t, "Start")
	mockTransactionBroker.AssertNotCalled(t, "Rollback")
}

func Test_Transactional_WithStorageLayer_Commit(t *testing.T) {
	target := setupTransactional()
	expect := sc.U32(2)

	mockStorageValue.On("Get").Return(transactionLevel).Once()
	mockStorageValue.On("Put", transactionLevel+1).Once()
	mockTransactionBroker.On("Start").Return()
	mockTransactionBroker.On("Commit").Return()
	mockStorageValue.On("Get").Return(transactionLevel + 1).Once()
	mockStorageValue.On("Put", transactionLevel).Once()

	res, err := target.WithStorageLayer(func() (sc.U32, primitives.DispatchError) {
		return expect, nil
	})

	assert.Equal(t, expect, res)
	assert.Nil(t, err)

	mockStorageValue.AssertNumberOfCalls(t, "Get", 2)
	mockStorageValue.AssertCalled(t, "Put", transactionLevel+1)
	mockTransactionBroker.AssertCalled(t, "Start")
	mockTransactionBroker.AssertCalled(t, "Commit")
	mockStorageValue.AssertCalled(t, "Put", transactionLevel)
}

func Test_Transactional_WithStorageLayer_Rollback(t *testing.T) {
	target := setupTransactional()
	expect := primitives.NewDispatchErrorCorruption()

	mockStorageValue.On("Get").Return(transactionLevel).Once()
	mockStorageValue.On("Put", transactionLevel+1).Once()
	mockTransactionBroker.On("Start").Return()
	mockTransactionBroker.On("Rollback").Return()
	mockStorageValue.On("Get").Return(transactionLevel + 1).Once()
	mockStorageValue.On("Put", transactionLevel).Once()

	res, err := target.WithStorageLayer(func() (sc.U32, primitives.DispatchError) {
		return sc.U32(0), expect
	})

	assert.Equal(t, sc.U32(0), res)
	assert.Equal(t, expect, err)

	mockStorageValue.AssertNumberOfCalls(t, "Get", 2)
	mockStorageValue.AssertNumberOfCalls(t, "Put", 2)
	mockStorageValue.AssertCalled(t, "Put", transactionLevel+1)
	mockTransactionBroker.AssertCalled(t, "Start")
	mockTransactionBroker.AssertCalled(t, "Rollback")
	mockStorageValue.AssertCalled(t, "Put", transactionLevel)
}

func setupTransactional() transactional[sc.U32, primitives.DispatchError] {
	mockStorageValue = new(mocks.StorageValue[sc.U32])
	mockTransactionBroker = new(mocks.IoTransactionBroker)

	target := NewTransactional[sc.U32, primitives.DispatchError]().(transactional[sc.U32, primitives.DispatchError])
	target.storage = mockStorageValue
	target.transactionBroker = mockTransactionBroker

	return target
}
