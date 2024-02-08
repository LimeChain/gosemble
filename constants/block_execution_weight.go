// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: 2024-02-08 11:36:16.01503 +0200 EET m=+0.503037126, STEPS: 50, REPEAT: 20, DBCACHE: 1024, HEAPPAGES: 4096, HOSTNAME: MacBook-Pro.local, CPU: arm64, GC: extalloc, TINYGO VERSION: 0.31.0-dev, TARGET: polkawasm

// Summary:
// Total: 298361000.000000, Min: 1648000.000000, Max: 12352000.000000, Average: 2983610.000000, Median: 1749000.000000, Stddev: 3338646.141462, Percentiles 99th, 95th, 75th: 12284500.000000, 12030000.000000, 1809000.000000
package constants

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func blockExecutionWeight(multiplier sc.U64) primitives.Weight {
	const refTime sc.U64 = 2983610

	return primitives.WeightFromParts(
		sc.SaturatingMulU64(multiplier, refTime),
		0,
	)
}
