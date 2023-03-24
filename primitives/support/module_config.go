package support

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

type ModuleMetadata interface {
	Functions() map[sc.U8]FunctionMetadata
	PreDispatch(call Call) (sc.Empty, primitives.TransactionValidityError)
	ValidateUnsigned(source primitives.TransactionSource, call Call) (primitives.ValidTransaction, primitives.TransactionValidityError)
}

type FunctionMetadata interface {
	BaseWeight(...any) primitives.Weight
	ClassifyDispatch(baseWeight primitives.Weight) primitives.DispatchClass
	Decode(buffer *bytes.Buffer) sc.VaryingData
	Dispatch(origin primitives.RuntimeOrigin, args sc.VaryingData) primitives.DispatchResultWithPostInfo[primitives.PostDispatchInfo]
	IsInherent() bool
	PaysFee(baseWeight primitives.Weight) primitives.Pays
	WeightInfo(baseWeight primitives.Weight) primitives.Weight

	// WeightFee        types.Pays
	// LengthFee        types.Pays
}

type Call interface {
	Function() FunctionMetadata
}
