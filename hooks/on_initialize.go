package hooks

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type OnInitialize interface {
	OnInitialize(n sc.U64) primitives.Weight
}
