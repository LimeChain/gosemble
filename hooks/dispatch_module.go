package hooks

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type DispatchModule[N sc.Numeric] interface {
	OnInitialize[N]
	OnRuntimeUpgrade
	OnFinalize(n N)
	OnIdle(n N, remainingWeight primitives.Weight) primitives.Weight
	OffchainWorker(n N)
}

type DefaultDispatchModule[N sc.Numeric] struct{}

func (dmh DefaultDispatchModule[N]) OnInitialize(n N) primitives.Weight {
	return primitives.WeightZero()
}

func (dmh DefaultDispatchModule[N]) OnRuntimeUpgrade() primitives.Weight {
	return primitives.WeightZero()
}

func (dmh DefaultDispatchModule[N]) OnFinalize(n N) {}

func (dmh DefaultDispatchModule[N]) OnIdle(n N, remainingWeight primitives.Weight) primitives.Weight {
	return primitives.WeightZero()
}

func (dmh DefaultDispatchModule[N]) OffchainWorker(n N) {}
