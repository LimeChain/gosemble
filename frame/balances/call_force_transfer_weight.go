// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: `2024-03-05 21:25:30.174313 +0200 EET m=+0.471085168`, STEPS: `50`, REPEAT: `20`, DBCACHE: `1024`, HEAPPAGES: `4096`, HOSTNAME: `MacBook-Pro.local`, CPU: `Apple M2 Pro(10 cores, 3504 mhz)`, GC: ``, TINYGO VERSION: ``, TARGET: ``

// Summary:
// BaseExtrinsicTime: 3394200000, BaseReads: 2, BaseWrites: 2, SlopesExtrinsicTime: [], SlopesReads: [], SlopesWrites: [], MinExtrinsicTime: 3394200, MinReads: 2, MinWrites: 2

package balances

import (
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callForceTransferWeight(dbWeight primitives.RuntimeDbWeight) primitives.Weight {
	return primitives.WeightFromParts(3394200000, 0).
		SaturatingAdd(dbWeight.Reads(2)).
		SaturatingAdd(dbWeight.Writes(2))
}
