package types

import (
	sc "github.com/LimeChain/goscale"
)

type StoredMap interface {
	EventDepositor
	Get(key AccountId) (AccountInfo, error)
	CanDecProviders(who AccountId) (bool, error)
	TryMutateExists(who AccountId, f func(who *AccountData) (sc.Encodable, error)) (sc.Encodable, error)
}
