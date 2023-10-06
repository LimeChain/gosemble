package mocks

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
	"github.com/stretchr/testify/mock"
)

type StoredMap struct {
	mock.Mock
}

func (m *StoredMap) DepositEvent(event types.Event) {
	_ = m.Called(event)
}

func (m *StoredMap) Get(key types.PublicKey) types.AccountInfo {
	args := m.Called(key)

	return args.Get(0).(types.AccountInfo)
}

func (m *StoredMap) CanDecProviders(who types.Address32) bool {
	args := m.Called(who)

	return args.Get(0).(bool)
}

func (m *StoredMap) TryMutateExists(who types.Address32, f func(who *types.AccountData) sc.Result[sc.Encodable]) sc.Result[sc.Encodable] {
	args := m.Called(who, f)

	return args.Get(0).(sc.Result[sc.Encodable])
}
