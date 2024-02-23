// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: `2024-03-07 15:53:46.286418 +0200 EET m=+10.890857876`, STEPS: `50`, REPEAT: `20`, DBCACHE: `1024`, HEAPPAGES: `4096`, HOSTNAME: `Rados-MacBook-Pro.local`, CPU: `Apple M1 Pro(8 cores, 3228 mhz)`, GC: ``, TINYGO VERSION: ``, TARGET: ``

// Summary:
// BaseExtrinsicTime: 114571052, BaseReads: 0, BaseWrites: 0, SlopesExtrinsicTime: [3434962], SlopesReads: [0], SlopesWrites: [0], MinExtrinsicTime: 92400, MinReads: 0, MinWrites: 0

package system

import (sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callKillPrefixWeight(dbWeight primitives.RuntimeDbWeight, size sc.U64) primitives.Weight {
	return primitives.WeightFromParts(114571052, 0).
			SaturatingAdd(primitives.WeightFromParts(3434962, 0).SaturatingMul(size)).
		SaturatingAdd(dbWeight.Reads(0)).
		SaturatingAdd(dbWeight.Writes(0))
}
