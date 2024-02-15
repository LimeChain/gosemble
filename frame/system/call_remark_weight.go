// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: 2024-02-08 12:21:26.168771 +0200 EET m=+34.595230251, STEPS: 50, REPEAT: 20, DBCACHE: 1024, HEAPPAGES: 4096, HOSTNAME: MacBook-Pro.local, CPU: Apple M2 Pro(10 cores, 3504 mhz), GC: extalloc, TINYGO VERSION: 0.31.0-dev, TARGET: polkawasm

// Summary:
// BaseExtrinsicTime: 70088888, BaseReads: 0, BaseWrites: 0, SlopesExtrinsicTime: [0], SlopesReads: [0], SlopesWrites: [0], MinExtrinsicTime: 65950, MinReads: 0, MinWrites: 0
package system

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callRemarkWeight(dbWeight primitives.RuntimeDbWeight) primitives.Weight {
	const refTime sc.U64 = 70088888
	const reads sc.U64 = 0
	const writes sc.U64 = 0
	const proofSize sc.U64 = 0

	return primitives.WeightFromParts(refTime, 0).
		SaturatingAdd(primitives.WeightFromParts(0, proofSize)).
		SaturatingAdd(dbWeight.Reads(reads)).
		SaturatingAdd(dbWeight.Writes(writes))
}
