// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: 2024-02-08 08:21:56.693226 +0200 EET m=+34.189773876, STEPS: 50, REPEAT: 20, DBCACHE: 1024, HOSTNAME: MacBook-Pro.local, CPU: arm64, GC: extalloc, TINYGO VERSION: 0.31.0-dev, TARGET: polkawasm
package timestamp

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callSetWeight(dbWeight primitives.RuntimeDbWeight) primitives.Weight {
	const refTime sc.U64 = 131500000
	const reads sc.U64 = 2
	const writes sc.U64 = 1
	const proofSize sc.U64 = 0

	return primitives.WeightFromParts(refTime, 0).
		SaturatingAdd(primitives.WeightFromParts(0, proofSize)).
		SaturatingAdd(dbWeight.Reads(reads)).
		SaturatingAdd(dbWeight.Writes(writes))
}
