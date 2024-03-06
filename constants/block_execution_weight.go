// THIS FILE WAS GENERATED USING GOSEMBLE BENCHMARKING PACKAGE
// DATE: `2024-03-05 21:29:51.382184 +0200 EET m=+0.844070335`, STEPS: `50`, REPEAT: `20`, DBCACHE: `1024`, HEAPPAGES: `4096`, HOSTNAME: `MacBook-Pro.local`, CPU: `Apple M2 Pro(10 cores, 3504 mhz)`, GC: ``, TINYGO VERSION: ``, TARGET: ``

// Summary:
// Total: 591897000.000000, Min: 2891000.000000, Max: 17777000.000000, Average: 5918970.000000, Median: 4490500.000000, Stddev: 3545938.477907, Percentiles 99th, 95th, 75th: 17626500.000000, 14773000.000000, 5270000.000000

package constants

import (
	sc "github.com/LimeChain/goscale"
	primitives "github.com/LimeChain/gosemble/primitives/types"
)

func blockExecutionWeight(multiplier sc.U64) primitives.Weight {
	return primitives.WeightFromParts(
		sc.SaturatingMulU64(multiplier, 5918970),
		0,
	)
}
