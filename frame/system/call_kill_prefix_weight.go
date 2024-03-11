// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: `2024-03-11 12:30:02.448959 +0200 EET m=+16.455047917`, STEPS: `50`, REPEAT: `20`, DBCACHE: `1024`, HEAPPAGES: `4096`, HOSTNAME: `Rados-MBP.lan`, CPU: `Apple M1 Pro(8 cores, 3228 mhz)`, GC: ``, TINYGO VERSION: ``, TARGET: ``

// Summary:
// BaseExtrinsicTime: 115634662, BaseReads: 0, BaseWrites: 0, SlopesExtrinsicTime: [3413650], SlopesReads: [0], SlopesWrites: [1], MinExtrinsicTime: 88350, MinReads: 0, MinWrites: 1

package system

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callKillPrefixWeight(dbWeight primitives.RuntimeDbWeight, size sc.U64) primitives.Weight {
	return primitives.WeightFromParts(115634662, 0).
		SaturatingAdd(primitives.WeightFromParts(3413650, 0).SaturatingMul(size)).
		SaturatingAdd(dbWeight.Reads(0)).
		SaturatingAdd(dbWeight.Writes(0)).
		SaturatingAdd(dbWeight.Writes(1).SaturatingMul(size))
}
