// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: `2024-03-11 12:31:57.118419 +0200 EET m=+131.125246042`, STEPS: `50`, REPEAT: `20`, DBCACHE: `1024`, HEAPPAGES: `4096`, HOSTNAME: `Rados-MBP.lan`, CPU: `Apple M1 Pro(8 cores, 3228 mhz)`, GC: ``, TINYGO VERSION: ``, TARGET: ``

// Summary:
// BaseExtrinsicTime: 151600000, BaseReads: 2, BaseWrites: 1, SlopesExtrinsicTime: [], SlopesReads: [], SlopesWrites: [], MinExtrinsicTime: 151600, MinReads: 2, MinWrites: 1

package timestamp

import (
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callSetWeight(dbWeight primitives.RuntimeDbWeight) primitives.Weight {
	return primitives.WeightFromParts(151600000, 0).
		SaturatingAdd(dbWeight.Reads(2)).
		SaturatingAdd(dbWeight.Writes(1))
}
