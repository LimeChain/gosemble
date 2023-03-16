package system

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/system"
	dispatchables "github.com/LimeChain/gosemble/frame/system/dispatchables"
	"github.com/LimeChain/gosemble/primitives/support"
)

var Module = SystemModule{}

type SystemModule struct {
	Remark dispatchables.FnRemark
	// TODO: add more dispatchables
}

func (m SystemModule) Functions() []support.FunctionMetadata {
	return []support.FunctionMetadata{
		m.Remark,
		// TODO: add more dispatchables
	}
}

func (m SystemModule) Index() sc.U8 {
	return system.ModuleIndex
}
