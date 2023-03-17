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
	WeightInfo(baseWeight types.Weight) types.Weight
	ClassifyDispatch(baseWeight types.Weight) types.DispatchClass
	PaysFee(baseWeight types.Weight) types.Pays
	Dispatch(origin types.RuntimeOrigin, args ...sc.Encodable) types.DispatchResultWithPostInfo[types.PostDispatchInfo]

	// WeightFee        types.Pays
	// LengthFee        types.Pays
}
