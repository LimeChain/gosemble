package grandpa

import (
	"github.com/LimeChain/gosemble/frame/grandpa"
	"github.com/LimeChain/gosemble/primitives/hashing"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "GrandpaApi"
	apiVersion    = 3
)

type Module struct {
	grandpa grandpa.Module
}

func New(grandpa grandpa.Module) Module {
	return Module{
		grandpa,
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
	authorities := m.grandpa.Authorities()

	return utils.BytesToOffsetAndSize(authorities.Bytes())
}
