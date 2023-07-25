package hooks

import "github.com/LimeChain/goscale"

type OnTimestampSet[T goscale.Encodable] interface {
	OnTimestampSet(T)
}
