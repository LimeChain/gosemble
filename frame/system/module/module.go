package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/constants/system"
	dispatchables "github.com/LimeChain/gosemble/frame/system/dispatchables"
	"github.com/LimeChain/gosemble/primitives/support"
)

type SystemModule struct {
	functions map[sc.U8]support.FunctionMetadata
	// TODO: add more dispatchables
}

func NewSystemModule() SystemModule {
	functions := make(map[sc.U8]support.FunctionMetadata)
	functions[system.FunctionRemarkIndex] = dispatchables.FnRemark{}

	return SystemModule{
		functions: functions,
	}
}

func (sm SystemModule) Functions() map[sc.U8]support.FunctionMetadata {
	return sm.functions
}
