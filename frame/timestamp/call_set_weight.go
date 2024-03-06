// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: `2024-03-05 21:26:04.021039 +0200 EET m=+34.317516876`, STEPS: `50`, REPEAT: `20`, DBCACHE: `1024`, HEAPPAGES: `4096`, HOSTNAME: `MacBook-Pro.local`, CPU: `Apple M2 Pro(10 cores, 3504 mhz)`, GC: ``, TINYGO VERSION: ``, TARGET: ``

// Summary:
// BaseExtrinsicTime: 655600000, BaseReads: 2, BaseWrites: 1, SlopesExtrinsicTime: [], SlopesReads: [], SlopesWrites: [], MinExtrinsicTime: 655600, MinReads: 2, MinWrites: 1

package timestamp

import (
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callSetWeight(dbWeight primitives.RuntimeDbWeight) primitives.Weight {
	return primitives.WeightFromParts(655600000, 0).
		SaturatingAdd(dbWeight.Reads(2)).
		SaturatingAdd(dbWeight.Writes(1))
}
