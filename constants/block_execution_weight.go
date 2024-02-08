// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: 2024-02-08 12:16:45.718365 +0200 EET m=+0.503945251, STEPS: 50, REPEAT: 20, DBCACHE: 1024, HEAPPAGES: 4096, HOSTNAME: MacBook-Pro.local, CPU: Apple M2 Pro(10 cores, 3504 mhz), GC: extalloc, TINYGO VERSION: 0.31.0-dev, TARGET: polkawasm

// Summary:
// Total: 287408000.000000, Min: 1584000.000000, Max: 11921000.000000, Average: 2874080.000000, Median: 1668000.000000, Stddev: 3222238.779731, Percentiles 99th, 95th, 75th: 11898500.000000, 11539000.000000, 1740000.000000
package constants

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func blockExecutionWeight(multiplier sc.U64) primitives.Weight {
	const refTime sc.U64 = 2874080

	return primitives.WeightFromParts(
		sc.SaturatingMulU64(multiplier, refTime),
		0,
	)
}
