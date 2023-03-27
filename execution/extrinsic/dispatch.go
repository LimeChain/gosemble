package extrinsic

import (
	"github.com/LimeChain/gosemble/execution/types"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func GetDispatchInfo(xt types.CheckedExtrinsic) primitives.DispatchInfo {
	function := xt.Function.Function()
	baseWeight := function.BaseWeight(xt.Function.Args)

	return primitives.DispatchInfo{
		Weight:  function.WeightInfo(baseWeight),
		Class:   function.ClassifyDispatch(baseWeight),
		PaysFee: function.PaysFee(baseWeight),
	}
}
