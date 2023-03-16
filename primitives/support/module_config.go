package support

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type ModuleMetadata struct {
	Index     sc.U8
	Functions map[sc.U8]FunctionMetadata
}

type FunctionMetadata struct {
	Func       interface{}
	BaseWeight types.Weight // converted to fee
	WeightFee  types.Pays
	LengthFee  types.Pays
}

// TODO: move elsewhere

func WeighData(baseWeight types.Weight, args []byte) types.Weight {
	// TODO:
	return types.WeightFromRefTime(sc.U64(len(args))) // + baseWeight
}

func ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	// TODO:
	return types.NewDispatchClassNormal()
}

func PaysFee(baseWeight types.Weight) types.Pays {
	// TODO:
	return types.NewPaysYes()
}
