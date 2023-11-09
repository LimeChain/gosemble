package types

import (
	sc "github.com/LimeChain/goscale"
)

type StoredMap interface {
	EventDepositor
	Get(key AccountId[SignerAddress]) (AccountInfo, error)
	CanDecProviders(who AccountId[SignerAddress]) (bool, error)
	TryMutateExists(who AccountId[SignerAddress], f func(who *AccountData) sc.Result[sc.Encodable]) (sc.Result[sc.Encodable], error)
}
