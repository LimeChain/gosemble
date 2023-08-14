package hooks

import (
	"github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type OnInitialize[T goscale.Encodable] interface {
	OnInitialize(T) primitives.Weight
}
