// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: `2024-03-07 14:49:14.234637 +0200 EET m=+31.854662710`, STEPS: `50`, REPEAT: `20`, DBCACHE: `1024`, HEAPPAGES: `4096`, HOSTNAME: `Rados-MacBook-Pro.local`, CPU: `Apple M1 Pro(8 cores, 3228 mhz)`, GC: ``, TINYGO VERSION: ``, TARGET: ``

// Summary:
// BaseExtrinsicTime: 226523426, BaseReads: 0, BaseWrites: 0, SlopesExtrinsicTime: [8813986], SlopesReads: [0], SlopesWrites: [0], MinExtrinsicTime: 77150, MinReads: 0, MinWrites: 0

package system

import (sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func callKillStorageWeight(dbWeight primitives.RuntimeDbWeight, size sc.U64) primitives.Weight {
	return primitives.WeightFromParts(226523426, 0).
			SaturatingAdd(primitives.WeightFromParts(8813986, 0).SaturatingMul(size)).
		SaturatingAdd(dbWeight.Reads(0)).
		SaturatingAdd(dbWeight.Writes(0))
}
