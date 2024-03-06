// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: `2024-03-05 21:25:30.382205 +0200 EET m=+0.678975126`, STEPS: `50`, REPEAT: `20`, DBCACHE: `1024`, HEAPPAGES: `4096`, HOSTNAME: `MacBook-Pro.local`, CPU: `Apple M2 Pro(10 cores, 3504 mhz)`, GC: ``, TINYGO VERSION: ``, TARGET: ``

// Summary:
// BaseExtrinsicTime: 1236050000, BaseReads: 2, BaseWrites: 2, SlopesExtrinsicTime: [], SlopesReads: [], SlopesWrites: [], MinExtrinsicTime: 1236050, MinReads: 2, MinWrites: 2

package balances

import (
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callSetBalanceCreatingWeight(dbWeight primitives.RuntimeDbWeight) primitives.Weight {
	return primitives.WeightFromParts(1236050000, 0).
		SaturatingAdd(dbWeight.Reads(2)).
		SaturatingAdd(dbWeight.Writes(2))
}
