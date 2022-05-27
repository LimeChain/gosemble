package helpers

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ChainSafe/gossamer/lib/common"
	"github.com/ChainSafe/gossamer/lib/keystore"
	"github.com/ChainSafe/gossamer/lib/runtime"
	"github.com/ChainSafe/gossamer/lib/runtime/storage"
	"github.com/ChainSafe/gossamer/lib/runtime/wasmer"
	"github.com/ChainSafe/gossamer/lib/transaction"
	"github.com/ChainSafe/gossamer/lib/trie"
	"github.com/ChainSafe/log15"
)

const WASM_RUNTIME = "../build/substrate_runtime.wasm"

type TransactionState struct {
	mock.Mock
}

func (_m *TransactionState) AddToPool(vt *transaction.ValidTransaction) common.Hash {
	ret := _m.Called(vt)

	var r0 common.Hash
	if rf, ok := ret.Get(0).(func(*transaction.ValidTransaction) common.Hash); ok {
		r0 = rf(vt)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Hash)
		}
	}

	return r0
}

func newTransactionStateMock() *TransactionState {
	m := new(TransactionState)
	m.On("AddToPool", mock.AnythingOfType("*transaction.ValidTransaction")).Return(common.BytesToHash([]byte("test")))
	return m
}

func setupConfig(t *testing.T, tt *trie.Trie, role byte) *wasmer.Config {
	t.Helper()

	s, err := storage.NewTrieState(tt)
	require.NoError(t, err)

	ns := runtime.NodeStorage{
		LocalStorage:      runtime.NewInMemoryDB(t),
		PersistentStorage: runtime.NewInMemoryDB(t),
		// BaseDB:            runtime.NewInMemoryDB(t),
	}
	cfg := &wasmer.Config{
		Imports: wasmer.ImportsNodeRuntime,
	}
	cfg.Storage = s
	cfg.Keystore = keystore.NewGlobalKeystore()
	cfg.LogLvl = log15.LvlInfo
	cfg.NodeStorage = ns
	cfg.Network = new(runtime.TestRuntimeNetwork)
	cfg.Transaction = newTransactionStateMock()
	cfg.Role = role
	return cfg
}

func NewTestInstanceWithTrie(t *testing.T, tt *trie.Trie) *wasmer.Instance {
	t.Helper()

	cfg := setupConfig(t, tt, 0)
	runtimeFilepath, err := filepath.Abs(WASM_RUNTIME)
	require.NoError(t, err)

	r, err := wasmer.NewInstanceFromFile(runtimeFilepath, cfg)
	require.NoError(t, err, "Got error when trying to create new VM")
	require.NotNil(t, r, "Could not create new VM instance")

	return r
}
