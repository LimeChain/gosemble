package grandpa

import (
	"github.com/LimeChain/gosemble/frame/grandpa"
	"github.com/LimeChain/gosemble/primitives/hashing"
	"github.com/LimeChain/gosemble/primitives/log"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "GrandpaApi"
	apiVersion    = 3
)

type Module struct {
	grandpa  grandpa.GrandpaModule
	memUtils utils.WasmMemoryTranslator
}

func New(grandpa grandpa.GrandpaModule) Module {
	return Module{
		grandpa:  grandpa,
		memUtils: utils.NewMemoryTranslator(),
	}
}

func (m Module) Name() string {
	return ApiModuleName
}

func (m Module) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return primitives.NewApiItem(hash, apiVersion)
}

func (m Module) Authorities() int64 {
	authorities, err := m.grandpa.Authorities()
	if err != nil {
		log.Critical(err.Error())
		return 0
	}
	return m.memUtils.BytesToOffsetAndSize(authorities.Bytes())
}
