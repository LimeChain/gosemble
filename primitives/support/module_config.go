package support

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type ModuleMetadata interface {
	Index() sc.U8
	Functions() []FunctionMetadata
}

type FunctionMetadata interface {
	Index() sc.U8
	BaseWeight(...any) types.Weight
	WeightInfo(baseWeight types.Weight, target sc.Sequence[sc.U8]) types.Weight
	ClassifyDispatch(baseWeight types.Weight, target sc.Sequence[sc.U8]) types.DispatchClass
	PaysFee(baseWeight types.Weight, target sc.Sequence[sc.U8]) types.Pays

	// WeightFee        types.Pays
	// LengthFee        types.Pays
}
