// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: 2024-02-20 09:10:00.456063 +0200 EET m=+35.863787001, STEPS: 50, REPEAT: 20, DBCACHE: 1024, HEAPPAGES: 4096, HOSTNAME: Rados-MBP.lan, CPU: Apple M1 Pro(8 cores, 3228 mhz), GC: , TINYGO VERSION: , TARGET:
// Summary:
// BaseExtrinsicTime: 106700000, BaseReads: 1, BaseWrites: 2, SlopesExtrinsicTime: [], SlopesReads: [], SlopesWrites: [], MinExtrinsicTime: 106700, MinReads: 1, MinWrites: 2
package system

import (
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callAuthorizeUpgradeWithoutChecksWeight(dbWeight primitives.RuntimeDbWeight) primitives.Weight {
	return primitives.WeightFromParts(0, 0).
		SaturatingAdd(primitives.WeightFromParts(0, 0)).
		SaturatingAdd(dbWeight.Reads(0)).
		SaturatingAdd(dbWeight.Writes(0))
}
