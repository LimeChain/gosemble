package types

import sc "github.com/LimeChain/goscale"

type OnRuntimeUpgrade interface {
	OnRuntimeUpgrade() Weight
}

type OnInitialize interface {
	OnInitialize(n sc.U64) (Weight, error)
}

type DispatchModule interface {
	OnInitialize
	OnRuntimeUpgrade
	OnFinalize(n sc.U64) error
	OnIdle(n sc.U64, remainingWeight Weight) Weight
	OffchainWorker(n sc.U64)
}
