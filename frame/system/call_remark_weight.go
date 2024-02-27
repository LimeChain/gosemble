// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: 2024-02-27 09:04:21.843885 +0200 EET m=+92.346653293, STEPS: `5`, REPEAT: `2`, DBCACHE: `1024`, HEAPPAGES: 4096, HOSTNAME: MacBook-Pro.local, CPU: Apple M2 Pro(10 cores, 3504 mhz), GC: , TINYGO VERSION: , TARGET:

// Summary:
// BaseExtrinsicTime: 70000000, BaseReads: 0, BaseWrites: 0, SlopesExtrinsicTime: [11], SlopesReads: [0], SlopesWrites: [0], MinExtrinsicTime: 62500, MinReads: 0, MinWrites: 0

package system

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callRemarkWeight(dbWeight primitives.RuntimeDbWeight, size sc.U64) primitives.Weight {
	return primitives.WeightFromParts(70000000, 0).
		SaturatingAdd(primitives.WeightFromParts(11, 0).SaturatingMul(size))
}
