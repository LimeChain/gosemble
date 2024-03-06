// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: `2024-03-05 21:25:31.291376 +0200 EET m=+1.588138584`, STEPS: `50`, REPEAT: `20`, DBCACHE: `1024`, HEAPPAGES: `4096`, HOSTNAME: `MacBook-Pro.local`, CPU: `Apple M2 Pro(10 cores, 3504 mhz)`, GC: ``, TINYGO VERSION: ``, TARGET: ``

// Summary:
// BaseExtrinsicTime: 3051650000, BaseReads: 1, BaseWrites: 1, SlopesExtrinsicTime: [], SlopesReads: [], SlopesWrites: [], MinExtrinsicTime: 3051650, MinReads: 1, MinWrites: 1

package balances

import (
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callTransferWeight(dbWeight primitives.RuntimeDbWeight) primitives.Weight {
	return primitives.WeightFromParts(3051650000, 0).
		SaturatingAdd(dbWeight.Reads(1)).
		SaturatingAdd(dbWeight.Writes(1))
}
