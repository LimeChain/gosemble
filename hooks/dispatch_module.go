package hooks

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type DefaultDispatchModule struct{}

func (dmh DefaultDispatchModule) OnInitialize(n sc.U64) (primitives.Weight, error) {
	return primitives.WeightZero(), nil
}

func (dmh DefaultDispatchModule) OnRuntimeUpgrade() primitives.Weight {
	return primitives.WeightZero()
}

func (dmh DefaultDispatchModule) OnFinalize(n sc.U64) error { return nil }

func (dmh DefaultDispatchModule) OnIdle(n sc.U64, remainingWeight primitives.Weight) primitives.Weight {
	return primitives.WeightZero()
}

func (dmh DefaultDispatchModule) OffchainWorker(n sc.U64) {}
