package benchmarking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test taken from substrate:
// https://github.com/LimeChain/polkadot-sdk/blob/03841f6c0f51c6be6f491ce404e40d8323c994f1/substrate/frame/benchmarking/src/analysis.rs#L589
func TestMedianSlopesAnalysis(t *testing.T) {
	data := []benchmarkResult{
		{[]linear{{value: 1}, {value: 5}}, 11_500_000, 3, 10},
		{[]linear{{value: 2}, {value: 5}}, 12_500_000, 4, 10},
		{[]linear{{value: 3}, {value: 5}}, 13_500_000, 5, 10},
		{[]linear{{value: 4}, {value: 5}}, 14_500_000, 6, 10},
		{[]linear{{value: 3}, {value: 1}}, 13_100_000, 5, 2},
		{[]linear{{value: 3}, {value: 3}}, 13_300_000, 5, 6},
		{[]linear{{value: 3}, {value: 7}}, 13_700_000, 5, 14},
		{[]linear{{value: 3}, {value: 10}}, 14_000_000, 5, 20},
	}

	expectedAnalysis := analysis{
		baseExtrinsicTime:       10_000_000_000,
		slopesExtrinsicTime:     []uint64{1_000_000_000, 100_000_000},
		minimumExtrinsicTime:    11_500_000,
		baseReads:               2,
		slopesReads:             []uint64{1, 0},
		minimumReads:            3,
		baseWrites:              0,
		slopesWrites:            []uint64{0, 2},
		minimumWrites:           2,
		componentExtrinsicTimes: []componentSlope{{Slope: 1000000000}, {Slope: 100000000}},
		componentReads:          []componentSlope{{Slope: 1}},
		componentWrites:         []componentSlope{{Slope: 2}},
		componentNames:          []string{"", ""},
	}

	medianSlopesRes := medianSlopesAnalysis(data)
	assert.Equal(t, expectedAnalysis, medianSlopesRes)

	medianSlopesRes = medianSlopesAnalysis([]benchmarkResult{})
	assert.Equal(t, analysis{}, medianSlopesRes)
}

func TestMedianValuesAnalysis(t *testing.T) {
	data := []benchmarkResult{
		{[]linear{}, 11_500_000, 3, 10},
		{[]linear{}, 12_500_000, 4, 10},
		{[]linear{}, 13_500_000, 5, 10},
		{[]linear{}, 14_500_000, 6, 10},
		{[]linear{}, 13_100_000, 5, 2},
		{[]linear{}, 13_300_000, 5, 6},
		{[]linear{}, 13_700_000, 5, 14},
		{[]linear{}, 14_000_000, 5, 20},
	}

	expectedAnalysis := analysis{
		baseExtrinsicTime:    13_500_000_000,
		minimumExtrinsicTime: 11_500_000,
		baseReads:            5,
		minimumReads:         3,
		baseWrites:           10,
		minimumWrites:        2,
	}

	medianSlopesRes := medianSlopesAnalysis(data)
	assert.Equal(t, expectedAnalysis, medianSlopesRes)

	medianValuesRes := medianValuesAnalysis(data)
	assert.Equal(t, expectedAnalysis, medianValuesRes)

	medianValuesRes = medianValuesAnalysis([]benchmarkResult{})
	assert.Equal(t, analysis{}, medianValuesRes)
}
