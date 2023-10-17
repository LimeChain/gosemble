package types

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type RuntimeApi struct {
	apis []primitives.ApiModule
}

func NewRuntimeApi(apis []primitives.ApiModule) RuntimeApi {
	return RuntimeApi{apis: apis}
}

func (ra RuntimeApi) Items() sc.Sequence[primitives.ApiItem] {
	items := sc.Sequence[primitives.ApiItem]{}

	for _, api := range ra.apis {
		items = append(items, api.Item())
	}

	return items
}

func (ra RuntimeApi) Module(name string) primitives.ApiModule {
	for _, module := range ra.apis {
		if module.Name() == name {
			return module
		}
	}
	log.Critical("runtime module [" + name + "] not found.")

	panic("unreachable")
}
