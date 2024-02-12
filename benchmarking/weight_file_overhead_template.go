// DO NOT EDIT - TEMPLATE FILE USED FOR GENERATING WEIGHT FILES (see weight_file_generator.go)
package benchmarking

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func overheadWeightFn(multiplier sc.U64) primitives.Weight {
	const refTime sc.U64 = 0

	return primitives.WeightFromParts(
		sc.SaturatingMulU64(multiplier, refTime),
		0,
	)
}
