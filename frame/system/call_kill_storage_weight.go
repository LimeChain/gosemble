// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: `2024-03-11 12:30:16.88032 +0200 EET m=+30.886501876`, STEPS: `50`, REPEAT: `20`, DBCACHE: `1024`, HEAPPAGES: `4096`, HOSTNAME: `Rados-MBP.lan`, CPU: `Apple M1 Pro(8 cores, 3228 mhz)`, GC: ``, TINYGO VERSION: ``, TARGET: ``

// Summary:
// BaseExtrinsicTime: 73285788, BaseReads: 0, BaseWrites: 0, SlopesExtrinsicTime: [8733462], SlopesReads: [0], SlopesWrites: [1], MinExtrinsicTime: 78100, MinReads: 0, MinWrites: 0

package system

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callKillStorageWeight(dbWeight primitives.RuntimeDbWeight, size sc.U64) primitives.Weight {
	return primitives.WeightFromParts(73285788, 0).
		SaturatingAdd(primitives.WeightFromParts(8733462, 0).SaturatingMul(size)).
		SaturatingAdd(dbWeight.Reads(0)).
		SaturatingAdd(dbWeight.Writes(0)).
		SaturatingAdd(dbWeight.Writes(1).SaturatingMul(size))
}
