package types

import (
	sc "github.com/LimeChain/goscale"
)

type StoredMap interface {
	EventDepositor
	Get(key AccountId) (AccountInfo, error)
	Put(key AccountId, accInfo AccountInfo)
	CanDecProviders(who AccountId) (bool, error)
	TryMutateExists(who AccountId, f func(who *AccountData) sc.Result[sc.Encodable]) (sc.Result[sc.Encodable], error)
	IncProviders(who AccountId) (IncRefStatus, error)
}
