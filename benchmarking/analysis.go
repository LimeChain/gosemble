package benchmarking

import benchmarkingtypes "github.com/LimeChain/gosemble/primitives/benchmarking"

// todo other analysis types
// todo enum for analysis types
// todo all enum values to implement Analysis interface

type benchmarkResult struct {
	extrinsicTime uint64
	reads, writes uint32
	components    []uint32
}

func newBenchmarkResult(benchmarkRes benchmarkingtypes.BenchmarkResult, componentValues []uint32) benchmarkResult {
	return benchmarkResult{
		extrinsicTime: benchmarkRes.ExtrinsicTime.ToBigInt().Uint64(),
		reads:         uint32(benchmarkRes.Reads),
		writes:        uint32(benchmarkRes.Writes),
		components:    componentValues,
	}
}

type analysis struct {
	base   uint64
	slopes []uint64
	names  []string
	// todo value_dists: Option<Vec<(Vec<u32>, u128, u128)>>
	minimum uint64
}

func minSquareAnalysis([]benchmarkResult) analysis {
	// todo
	return analysis{}
}
