package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type MockStoredMap struct {
	mock.Mock
}

func (m *MockStoredMap) DepositEvent(event types.Event) {
	_ = m.Called(event)
}

func (m *MockStoredMap) Get(key types.PublicKey) types.AccountInfo {
	args := m.Called(key)

	return args[0].(types.AccountInfo)
}

func (m *MockStoredMap) CanDecProviders(who types.Address32) bool {
	args := m.Called(who)

	return args[0].(bool)
}

func (m *MockStoredMap) Mutate(who types.Address32, f func(who *types.AccountInfo) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	args := m.Called(who, f)

	return args[0].(sc.Result[sc.Encodable])
}

func (m *MockStoredMap) TryMutateExists(who types.Address32, f func(who *types.AccountData) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	args := m.Called(who, f)

	return args[0].(sc.Result[sc.Encodable])
}
