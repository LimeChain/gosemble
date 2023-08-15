package hooks

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type DispatchModule[T sc.Encodable] interface {
	OnInitialize[T]
	OnRuntimeUpgrade
	OnFinalize(n T)
	OnIdle(n T, remainingWeight primitives.Weight) primitives.Weight
	OffchainWorker(n T)
}

type DefaultDispatchModule[T sc.Encodable] struct{}

func (dmh DefaultDispatchModule[T]) OnInitialize(n T) primitives.Weight {
	return primitives.WeightZero()
}

func (dmh DefaultDispatchModule[T]) OnRuntimeUpgrade() primitives.Weight {
	return primitives.WeightZero()
}

func (dmh DefaultDispatchModule[T]) OnFinalize(n T) {}

func (dmh DefaultDispatchModule[T]) OnIdle(n T, remainingWeight primitives.Weight) primitives.Weight {
	return primitives.WeightZero()
}

func (dmh DefaultDispatchModule[T]) OffchainWorker(n T) {}
