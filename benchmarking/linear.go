package benchmarking

import (
	"errors"
	"math"
)

var (
	errInvalidValues = errors.New("failed to initialize new linear component: linear component min value must be less than or equal to max value")
)

type linear struct {
	min, max uint32
	value    uint32
}

// Used by benchmark tests to specify that a benchmarking variable is linear
// over some specified range, i.e. `NewLinear(0, 1000)` means that the corresponding variable
// is allowed to range from `0` to `1000`, inclusive.
func NewLinear(min, max uint32) (*linear, error) {
	if max < min {
		return nil, errInvalidValues
	}

	// max is the default linear value
	return &linear{min: min, max: max, value: max}, nil
}

// Component value, modified for each step iteration
func (l *linear) Value() uint32 {
	return l.value
}

// Internal function used to set the linear value for each step iteration
func (l *linear) setValue(value uint32) {
	l.value = value
}

// Internal function that calculates linear values for given steps
func (l *linear) values(steps int) []uint32 {
	stepSize := math.Max(float64(l.max-l.min)/float64(steps-1), 0)

	values := make([]uint32, steps)
	for step := 0; step < steps; step++ {
		stepValue := uint32(float64(l.min) + stepSize*float64(step))
		values[step] = stepValue
	}

	return values
}

// Internal function that returns all component values
func componentValues(linearComponents []*linear) []uint32 {
	values := make([]uint32, len(linearComponents))
	for i, l := range linearComponents {
		values[i] = l.Value()
	}

	return values
}
