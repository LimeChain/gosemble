// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: 2024-02-08 08:21:23.756162 +0200 EET m=+1.252778084, STEPS: 50, REPEAT: 20, DBCACHE: 1024, HOSTNAME: MacBook-Pro.local, CPU: arm64, GC: extalloc, TINYGO VERSION: 0.31.0-dev, TARGET: polkawasm
package balances

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callTransferKeepAliveWeight(dbWeight primitives.RuntimeDbWeight) primitives.Weight {
	const refTime sc.U64 = 1667450000
	const reads sc.U64 = 1
	const writes sc.U64 = 1
	const proofSize sc.U64 = 0

	return primitives.WeightFromParts(refTime, 0).
		SaturatingAdd(primitives.WeightFromParts(0, proofSize)).
		SaturatingAdd(dbWeight.Reads(reads)).
		SaturatingAdd(dbWeight.Writes(writes))
}
