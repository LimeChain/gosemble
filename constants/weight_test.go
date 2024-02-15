package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeights(t *testing.T) {
	assert.Equal(t, baseExtrinsicWeight(WeightRefTimePerNanos), ExtrinsicBaseWeight)
	assert.Equal(t, blockExecutionWeight(WeightRefTimePerNanos), BlockExecutionWeight)
}
