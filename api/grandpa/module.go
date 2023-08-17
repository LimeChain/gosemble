package grandpa

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/frame/grandpa"
	"github.com/LimeChain/gosemble/primitives/hashing"
	primitives "github.com/LimeChain/gosemble/primitives/types"
	"github.com/LimeChain/gosemble/utils"
)

const (
	ApiModuleName = "GrandpaApi"
	apiVersion    = 3
)

type Module[N sc.Numeric] struct {
	grandpa grandpa.Module[N]
}

func New[N sc.Numeric](grandpa grandpa.Module[N]) Module[N] {
	return Module[N]{
		grandpa,
	}
}

func (m Module[N]) Name() string {
	return ApiModuleName
}

func (m Module[N]) Item() primitives.ApiItem {
	hash := hashing.MustBlake2b8([]byte(ApiModuleName))
	return primitives.NewApiItem(hash, apiVersion)
}

func (m Module[N]) Authorities() int64 {
	authorities := m.grandpa.Authorities()

	return utils.BytesToOffsetAndSize(authorities.Bytes())
}
