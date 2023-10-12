package hooks

import primitives "github.com/LimeChain/gosemble/primitives/types"

type DefaultOnRuntimeUpgrade struct{}

func (doru DefaultOnRuntimeUpgrade) OnRuntimeUpgrade() primitives.Weight {
	return primitives.WeightZero()
}
