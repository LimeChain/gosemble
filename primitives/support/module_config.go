package support

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type ModuleMetadata interface {
	Functions() map[sc.U8]FunctionMetadata
}

type FunctionMetadata interface {
	BaseWeight(...any) types.Weight
	ClassifyDispatch(baseWeight types.Weight) types.DispatchClass
	Decode(buffer *bytes.Buffer) sc.VaryingData
	Dispatch(origin types.RuntimeOrigin, args sc.VaryingData) types.DispatchResultWithPostInfo[types.PostDispatchInfo]
	IsInherent() bool
	PaysFee(baseWeight types.Weight) types.Pays
	WeightInfo(baseWeight types.Weight) types.Weight

	// WeightFee        types.Pays
	// LengthFee        types.Pays
}
