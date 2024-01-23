package benchmarking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// test taken from substrate:
// https://github.com/LimeChain/polkadot-sdk/blob/03841f6c0f51c6be6f491ce404e40d8323c994f1/substrate/frame/benchmarking/src/analysis.rs#L589
func TestMedianSlopesAnalysis(t *testing.T) {
	data := []benchmarkResult{
		{[]uint32{1, 5}, 11_500_000, 3, 10},
		{[]uint32{2, 5}, 12_500_000, 4, 10},
		{[]uint32{3, 5}, 13_500_000, 5, 10},
		{[]uint32{4, 5}, 14_500_000, 6, 10},
		{[]uint32{3, 1}, 13_100_000, 5, 2},
		{[]uint32{3, 3}, 13_300_000, 5, 6},
		{[]uint32{3, 7}, 13_700_000, 5, 14},
		{[]uint32{3, 10}, 14_000_000, 5, 20},
	}

	res, err := medianSlopesAnalysis(data)
	assert.NoError(t, err)

	expectedAnalysis := analysis{
		baseExtrinsicTime:    10_000_000_000,
		slopesExtrinsicTime:  []uint64{1_000_000_000, 100_000_000},
		minimumExtrinsicTime: 11_500_000,
		baseReads:            2,
		slopesReads:          []uint64{1, 0},
		minimumReads:         3,
		baseWrites:           0,
		slopesWrites:         []uint64{0, 2},
		minimumWrites:        2,
	}

	assert.Equal(t, expectedAnalysis, res)
}
