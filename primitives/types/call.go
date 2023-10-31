package types

import (
	"bytes"

	sc "github.com/LimeChain/goscale"
)

type Call interface {
	sc.Encodable

	ModuleIndex() sc.U8
	FunctionIndex() sc.U8
	Args() sc.VaryingData
	Dispatch(origin RuntimeOrigin, args sc.VaryingData) DispatchResultWithPostInfo[PostDispatchInfo]
	BaseWeight() Weight
	ClassifyDispatch(baseWeight Weight) DispatchClass
	PaysFee(baseWeight Weight) Pays
	WeighData(baseWeight Weight) Weight
	DecodeArgs(buffer *bytes.Buffer) (Call, error)
}
