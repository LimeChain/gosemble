package hooks

import primitives "github.com/LimeChain/gosemble/primitives/types"

type OnRuntimeUpgrade interface {
	OnRuntimeUpgrade() primitives.Weight
}

type DefaultOnRuntimeUpgrade struct{}

func (doru DefaultOnRuntimeUpgrade) OnRuntimeUpgrade() primitives.Weight {
	return primitives.WeightZero()
}
