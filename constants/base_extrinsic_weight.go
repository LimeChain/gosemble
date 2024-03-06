// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: `2024-03-06 11:11:17.896577 +0200 EET m=+97.928613168`, STEPS: `50`, REPEAT: `20`, DBCACHE: `1024`, HEAPPAGES: `4096`, HOSTNAME: `MacBook-Pro.local`, CPU: `Apple M2 Pro(10 cores, 3504 mhz)`, GC: ``, TINYGO VERSION: ``, TARGET: ``

// Summary:
// Total: 428894879.000000, Min: 4109391.000000, Max: 4579784.000000, Average: 4288948.790000, Median: 4282162.000000, Stddev: 80080.072830, Percentiles 99th, 95th, 75th: 4566836.000000, 4436177.000000, 4322729.000000

package constants

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func baseExtrinsicWeight(multiplier sc.U64) primitives.Weight {
	return primitives.WeightFromParts(
		sc.SaturatingMulU64(multiplier, 4288948),
		0,
	)
}
