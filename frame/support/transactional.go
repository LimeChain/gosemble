package support

import (
	"errors"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/io"
	"github.com/LimeChain/gosemble/primitives/log"
	"github.com/LimeChain/gosemble/primitives/types"
)

// Layer is the type that is being used to store the current number of active layers.
type Layer = sc.U32

// TransactionalLimit returns the maximum number of nested layers.
const TransactionalLimit Layer = 255

var keyTransactionLevel = []byte(":transaction_level:")

type Transactional[T sc.Encodable, E types.DispatchError] interface {
	WithStorageLayer(fn func() (T, types.DispatchError)) (T, E)
}

type transactional[T sc.Encodable, E types.DispatchError] struct {
	storage           StorageValue[sc.U32]
	transactionBroker io.TransactionBroker
}

func NewTransactional[T sc.Encodable, E types.DispatchError]() Transactional[T, E] {
	return transactional[T, E]{
		storage:           NewSimpleStorageValue(keyTransactionLevel, sc.DecodeU32),
		transactionBroker: io.NewTransactionBroker(),
	}
}

// GetTransactionLevel returns the current number of nested transactional layers.
func (t transactional[T, E]) GetTransactionLevel() Layer {
	return t.storage.Get()
}

// SetTransactionLevel Set the current number of nested transactional layers.
func (t transactional[T, E]) SetTransactionLevel(level Layer) {
	t.storage.Put(level)
}

// KillTransactionLevel kill the transactional layers storage.
func (t transactional[T, E]) KillTransactionLevel() {
	t.storage.Clear()
}

// IncTransactionLevel increments the transaction level. Returns an error if levels go past the limit.
//
// Returns a guard that when dropped decrements the transaction level automatically.
func (t transactional[T, E]) IncTransactionLevel() (ok sc.Empty, err error) {
	existingLevels := t.GetTransactionLevel()
	if existingLevels >= TransactionalLimit {
		return ok, errors.New("transactional error limit reached")
	}
	// Cannot overflow because of check above.
	t.SetTransactionLevel(existingLevels + 1)
	return sc.Empty{}, err
}

func (t transactional[T, E]) DecTransactionLevel() {
	existingLevels := t.GetTransactionLevel()
	if existingLevels == 0 {
		log.Warn("We are underflowing with calculating transactional levels. Not great, but let's not panic...")
	} else if existingLevels == 1 {
		// Don't leave any trace of this storage item.
		t.KillTransactionLevel()
	} else {
		// Cannot underflow because of checks above.
		t.SetTransactionLevel(existingLevels - 1)
	}
}

// WithTransaction executes the supplied function in a new storage transaction.
//
// All changes to storage performed by the supplied function are discarded if the returned
// outcome is `TransactionOutcome::Rollback`.
//
// Transactions can be nested up to `TRANSACTIONAL_LIMIT` times; more than that will result in an
// error.
//
// Commits happen to the parent transaction.
func (t transactional[T, E]) WithTransaction(fn func() types.TransactionOutcome) (ok T, err E) {
	// This needs to happen before `start_transaction` below.
	// Otherwise we may rollback the increase, then decrease as the guard goes out of scope
	// and then end in some bad state.
	_, e := t.IncTransactionLevel()
	if e != nil {
		return ok, E(types.NewDispatchErrorTransactional(types.NewTransactionalErrorLimitReached()))
	}

	t.transactionBroker.Start()

	res := fn()

	switch res[0] {
	case types.TransactionOutcomeCommit:
		t.transactionBroker.Commit()
		t.DecTransactionLevel()
		return res[1].(T), err
	case types.TransactionOutcomeRollback:
		t.transactionBroker.Rollback()
		t.DecTransactionLevel()
		return ok, res[1].(E)
	default:
		log.Critical("invalid transaction outcome")
		return ok, err
	}
}

// WithStorageLayer executes the supplied function, adding a new storage layer.
//
// This is the same as `with_transaction`, but assuming that any function returning an `Err` should
// rollback, and any function returning `Ok` should commit. This provides a cleaner API to the
// developer who wants this behavior.
func (t transactional[T, E]) WithStorageLayer(fn func() (T, types.DispatchError)) (T, E) {
	return t.WithTransaction(
		func() types.TransactionOutcome {
			ok, err := fn()

			if err != nil {
				return types.NewTransactionOutcomeRollback(err)
			} else {
				return types.NewTransactionOutcomeCommit(ok)
			}
		},
	)
}
