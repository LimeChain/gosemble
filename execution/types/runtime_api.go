package types

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type RuntimeApi struct {
	apis map[string]primitives.ApiModule
}

func NewRuntimeApi(apis []primitives.ApiModule) RuntimeApi {
	result := make(map[string]primitives.ApiModule)
	for _, api := range apis {
		result[api.Name()] = api
	}

	return RuntimeApi{apis: result}
}

func (ra RuntimeApi) Items() sc.Sequence[primitives.ApiItem] {
	items := sc.Sequence[primitives.ApiItem]{}

	for _, api := range ra.apis {
		items = append(items, api.Item())
	}

	return items
}

func (ra RuntimeApi) Module(name string) primitives.ApiModule {
	module, ok := ra.apis[name]
	if !ok {
		log.Critical("runtime module [" + name + "] not found.")
	}

	return module
}
