package types

import (
	sc "github.com/LimeChain/goscale"
)

type StoredMap interface {
	EventDepositor
	Get(key PublicKey) (AccountInfo, error)
	CanDecProviders(who Address32) (bool, error)
	TryMutateExists(who Address32, f func(who *AccountData) sc.Result[sc.Encodable]) (sc.Result[sc.Encodable], error)
}
