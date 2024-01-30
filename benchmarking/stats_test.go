package benchmarking

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_OverheadStats(t *testing.T) {
	input := []float64{1, 2, 3}
	expect := OverheadStats{
		Sum:    6,
		Min:    1,
		Max:    3,
		Mean:   2,
		Median: 2,
		Stddev: 0.816496580927726,
		P99:    2.5,
		P95:    2.5,
		P75:    2.5,
	}

	target, err := NewOverheadStats(input)
	assert.Nil(t, err)

	assert.Equal(t, expect, target)
}
