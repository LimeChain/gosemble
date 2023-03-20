package module

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/support"
)

type AuraModule struct {
}

func NewAuraModule() AuraModule {
	return AuraModule{}
}

func (am AuraModule) Functions() map[sc.U8]support.FunctionMetadata {
	return map[sc.U8]support.FunctionMetadata{}
}
