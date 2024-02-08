// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: 2024-02-08 08:21:56.524289 +0200 EET m=+34.020837167, STEPS: 50, REPEAT: 20, DBCACHE: 1024, HOSTNAME: MacBook-Pro.local, CPU: arm64, GC: extalloc, TINYGO VERSION: 0.31.0-dev, TARGET: polkawasm
package system

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callRemarkWeight(dbWeight primitives.RuntimeDbWeight) primitives.Weight {
	const refTime sc.U64 = 63625001
	const reads sc.U64 = 0
	const writes sc.U64 = 0
	const proofSize sc.U64 = 0

	return primitives.WeightFromParts(refTime, 0).
		SaturatingAdd(primitives.WeightFromParts(0, proofSize)).
		SaturatingAdd(dbWeight.Reads(reads)).
		SaturatingAdd(dbWeight.Writes(writes))
}
