package types

import (
	sc "github.com/LimeChain/goscale"
)

type StoredMap interface {
	EventDepositor
	Get(key AccountId[PublicKey]) (AccountInfo, error)
	Put(key AccountId[PublicKey], accInfo AccountInfo)
	CanDecProviders(who AccountId[PublicKey]) (bool, error)
	TryMutateExists(who AccountId[PublicKey], f func(who *AccountData) sc.Result[sc.Encodable]) (sc.Result[sc.Encodable], error)
	IncProviders(who AccountId[PublicKey]) (IncRefStatus, error)
}
