package support

import (
	sc "github.com/LimeChain/goscale"
	"github.com/LimeChain/gosemble/primitives/types"
)

type ModuleMetadata struct {
	Index     sc.U8
	Functions map[string]FunctionMetadata
}

type FunctionMetadata struct {
	Index      sc.U8
	Func       interface{}
	BaseWeight types.Weight // converted to fee
	WeightFee  types.Pays
	LengthFee  types.Pays
}

// TODO: move elsewhere

func WeighData(baseWeight types.Weight, args sc.Sequence[sc.U8]) types.Weight {
	// TODO:
	return types.WeightFromRefTime(sc.U64(len(args.Bytes()))) // + baseWeight
}

func ClassifyDispatch(baseWeight types.Weight) types.DispatchClass {
	// TODO:
	return types.NewDispatchClassNormal()
}

func PaysFee(baseWeight types.Weight) types.Pays {
	// TODO:
	return types.NewPaysYes()
}
