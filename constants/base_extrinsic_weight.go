// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: 2024-02-08 12:20:22.21156 +0200 EET m=+216.996954293, STEPS: 50, REPEAT: 20, DBCACHE: 1024, HEAPPAGES: 4096, HOSTNAME: MacBook-Pro.local, CPU: Apple M2 Pro(10 cores, 3504 mhz), GC: extalloc, TINYGO VERSION: 0.31.0-dev, TARGET: polkawasm

// Summary:
// Total: 380342875.000000, Min: 3357385.000000, Max: 6924964.000000, Average: 3803428.750000, Median: 3552995.500000, Stddev: 784711.028755, Percentiles 99th, 95th, 75th: 6898816.000000, 5807199.000000, 3594239.000000
package constants

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func baseExtrinsicWeight(multiplier sc.U64) primitives.Weight {
	const refTime sc.U64 = 3803428

	return primitives.WeightFromParts(
		sc.SaturatingMulU64(multiplier, refTime),
		0,
	)
}
