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
	m.Called(event)
}

func (m *StoredMap) Get(key types.AccountId) (types.AccountInfo, error) {
	args := m.Called(key)

	if args.Get(1) == nil {
		return args.Get(0).(types.AccountInfo), nil
	}

	return args.Get(0).(types.AccountInfo), args.Get(1).(error)
}

func (m *StoredMap) Put(key types.AccountId, accInfo types.AccountInfo) {
	m.Called(key, accInfo)
}

func (m *StoredMap) CanDecProviders(who types.AccountId) (bool, error) {
	args := m.Called(who)

	if args.Get(1) == nil {
		return args.Get(0).(bool), nil
	}

	return args.Get(0).(bool), args.Get(1).(error)
}

func (m *StoredMap) TryMutateExists(who types.AccountId, f func(who *types.AccountData) sc.Result[sc.Encodable]) (sc.Result[sc.Encodable], error) {
	args := m.Called(who, f)

	if args.Get(1) == nil {
		return args.Get(0).(sc.Result[sc.Encodable]), nil
	}

	return args.Get(0).(sc.Result[sc.Encodable]), args.Get(1).(error)
}

func (m *StoredMap) incProviders(who types.AccountId) (types.IncRefStatus, error) {
	args := m.Called(who)
	if args.Get(1) == nil {
		return args.Get(0).(types.IncRefStatus), nil
	}
	return args.Get(0).(types.IncRefStatus), args.Get(1).(error)
}
