package hooks

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type OnInitialize[N sc.Numeric] interface {
	OnInitialize(N) primitives.Weight
}
