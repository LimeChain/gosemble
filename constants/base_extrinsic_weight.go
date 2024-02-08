// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: 2024-02-08 11:40:01.068305 +0200 EET m=+225.557151334, STEPS: 50, REPEAT: 20, DBCACHE: 1024, HEAPPAGES: 4096, HOSTNAME: MacBook-Pro.local, CPU: arm64, GC: extalloc, TINYGO VERSION: 0.31.0-dev, TARGET: polkawasm

// Summary:
// Total: 401009593.000000, Min: 3447291.000000, Max: 6853728.000000, Average: 4010095.930000, Median: 3774921.000000, Stddev: 540253.416748, Percentiles 99th, 95th, 75th: 6407273.000000, 5125738.000000, 3945403.000000
package constants

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func baseExtrinsicWeight(multiplier sc.U64) primitives.Weight {
	const refTime sc.U64 = 4010095

	return primitives.WeightFromParts(
		sc.SaturatingMulU64(multiplier, refTime),
		0,
	)
}
