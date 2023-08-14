package extrinsic

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type RuntimeExtrinsic struct {
	modules map[sc.U8]types.Module
}

func New(modules map[sc.U8]types.Module) RuntimeExtrinsic {
	return RuntimeExtrinsic{modules: modules}
}

func (re RuntimeExtrinsic) Module(index sc.U8) (module types.Module, isFound bool) {
	m, ok := re.modules[index]
	return m, ok
}
