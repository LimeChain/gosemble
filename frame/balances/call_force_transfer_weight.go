// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: 2024-02-08 12:20:51.981503 +0200 EET m=+0.407991418, STEPS: 50, REPEAT: 20, DBCACHE: 1024, HEAPPAGES: 4096, HOSTNAME: MacBook-Pro.local, CPU: Apple M2 Pro(10 cores, 3504 mhz), GC: extalloc, TINYGO VERSION: 0.31.0-dev, TARGET: polkawasm

// Summary:
// BaseExtrinsicTime: 1423850000, BaseReads: 2, BaseWrites: 2, SlopesExtrinsicTime: [], SlopesReads: [], SlopesWrites: [], MinExtrinsicTime: 1423850, MinReads: 2, MinWrites: 2
package balances

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callForceTransferWeight(dbWeight primitives.RuntimeDbWeight) primitives.Weight {
	const refTime sc.U64 = 1423850000
	const reads sc.U64 = 2
	const writes sc.U64 = 2
	const proofSize sc.U64 = 0

	return primitives.WeightFromParts(refTime, 0).
		SaturatingAdd(primitives.WeightFromParts(0, proofSize)).
		SaturatingAdd(dbWeight.Reads(reads)).
		SaturatingAdd(dbWeight.Writes(writes))
}
