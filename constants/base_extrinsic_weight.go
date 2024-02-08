// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: 2024-02-08 10:24:49.715009 +0200 EET m=+217.285868876, STEPS: 50, REPEAT: 20, DBCACHE: 1024, HOSTNAME: MacBook-Pro.local, CPU: arm64, GC: extalloc, TINYGO VERSION: 0.31.0-dev, TARGET: polkawasm
package constants

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func baseExtrinsicWeight(multiplier sc.U64) primitives.Weight {
	const refTime sc.U64 = 3671032

	return primitives.WeightFromParts(
		sc.SaturatingMulU64(multiplier, refTime),
		0,
	)
}
