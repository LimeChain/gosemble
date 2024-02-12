// DO NOT EDIT - TEMPLATE FILE USED FOR GENERATING WEIGHT FILES (see weight_file_generator.go)
package benchmarking

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func extrinsicWeightFn(dbWeight primitives.RuntimeDbWeight) primitives.Weight {
	const refTime sc.U64 = 0
	const reads sc.U64 = 0
	const writes sc.U64 = 0
	const proofSize sc.U64 = 0

	return primitives.WeightFromParts(refTime, 0).
		SaturatingAdd(primitives.WeightFromParts(0, proofSize)).
		SaturatingAdd(dbWeight.Reads(reads)).
		SaturatingAdd(dbWeight.Writes(writes))
}
